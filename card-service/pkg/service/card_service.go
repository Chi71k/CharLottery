package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	pb "card-service/pkg/api"
	"card-service/pkg/cache"
	"card-service/pkg/db"
	model "card-service/pkg/db/models"
	"card-service/pkg/middleware/handler"
	"card-service/pkg/natswrap"

	// Импортируем Redis

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardServiceServer struct {
	pb.UnimplementedCardServiceServer
	natsPub natswrap.Publisher
	natsSub natswrap.Subscriber
	cache   *cache.RedisClient // Верное имя и тип
}

func NewCardServiceServer(pub natswrap.Publisher, sub natswrap.Subscriber, redisClient *cache.RedisClient) *CardServiceServer {
	return &CardServiceServer{
		natsPub: pub,
		natsSub: sub,
		cache:   redisClient, // Используем переданный аргумент
	}
}

func (s *CardServiceServer) CreateCard(ctx context.Context, req *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	var cardNumber string

	if req.CardType == "BONUS" {
		cardNumber = generateBonusCardNumber()
	} else {
		cardNumber = req.CardNumber
	}

	card := &model.Card{
		CardNumber: cardNumber,
		CardHolder: req.CardholderName,
		ExpiryDate: req.ExpirationDate,
		CVV:        req.Cvv,
		UserID:     req.UserId,
		CardType:   handler.GetCardType(cardNumber),
	}

	if err := card.Validate(); err != nil {
		return nil, err
	}

	collection := db.GetCollection("cards")
	result, err := collection.InsertOne(ctx, card)
	if err != nil {
		return nil, err
	}

	oid := result.InsertedID.(primitive.ObjectID)
	s.natsPub.Publish("card.created", []byte(fmt.Sprintf("Card created: %s, UserID: %s", card.CardNumber, card.UserID)))

	return &pb.CreateCardResponse{
		CardId:   oid.Hex(),
		CardType: card.CardType,
		Message:  "Card created successfully",
	}, nil
}

func (s *CardServiceServer) ListCards(ctx context.Context, req *pb.ListCardsRequest) (*pb.ListCardsResponse, error) {
	collection := db.GetCollection("cards")
	cursor, err := collection.Find(ctx, bson.M{"user_id": req.UserId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []*pb.Card
	for cursor.Next(ctx) {
		var card model.Card
		if err := cursor.Decode(&card); err != nil {
			log.Println("Decode error:", err)
			continue
		}
		cards = append(cards, &pb.Card{
			CardId:         card.ID.Hex(),
			CardNumber:     card.CardNumber,
			CardholderName: card.CardHolder,
			ExpirationDate: card.ExpiryDate,
			Cvv:            card.CVV,
			UserId:         card.UserID,
			CardType:       card.CardType,
		})
	}

	s.natsPub.Publish("cards.listed", []byte(fmt.Sprintf("Cards listed for UserID: %s", req.UserId)))

	return &pb.ListCardsResponse{Cards: cards}, nil
}

func (s *CardServiceServer) ChargeCard(ctx context.Context, req *pb.ChargeCardRequest) (*pb.ChargeCardResponse, error) {
	collection := db.GetCollection("cards")
	oid, err := primitive.ObjectIDFromHex(req.CardId)
	if err != nil {
		return &pb.ChargeCardResponse{
			Success: false,
			Message: "Invalid card ID",
		}, nil
	}

	var card model.Card
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&card)
	if err != nil {
		return &pb.ChargeCardResponse{
			Success: false,
			Message: "Card not found",
		}, nil
	}

	log.Printf("Charging %.2f KZT from card %s\n", req.Amount, card.CardNumber)
	s.natsPub.Publish("card.charge", []byte(fmt.Sprintf("Charged %.2f KZT from Card %s", req.Amount, card.CardNumber)))

	return &pb.ChargeCardResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully charged %.2f KZT", req.Amount),
	}, nil
}

func (s *CardServiceServer) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	// Проверим, есть ли карта в Redis
	cardID := req.CardId
	cachedCard, err := s.cache.Get(cardID)

	if err == nil && cachedCard != "" {
		var card model.Card
		if err := json.Unmarshal([]byte(cachedCard), &card); err == nil {
			// Карта найдена в Redis, возвращаем
			log.Printf("🔍 Карта получена из Redis-кэша: %s", cardID)
			return &pb.GetCardResponse{
				Card: &pb.Card{
					CardId:         card.ID.Hex(),
					CardNumber:     card.CardNumber,
					CardholderName: card.CardHolder,
					ExpirationDate: card.ExpiryDate,
					Cvv:            card.CVV,
					UserId:         card.UserID,
					CardType:       card.CardType,
				},
			}, nil
		} else {
			log.Printf("⚠️ Ошибка десериализации карты из Redis: %v", err)
		}
	}

	// Если карты нет в Redis, достаем из базы данных
	collection := db.GetCollection("cards")
	oid, err := primitive.ObjectIDFromHex(cardID)
	if err != nil {
		return nil, err
	}

	var card model.Card
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&card)
	if err != nil {
		return nil, err
	}

	// Сохраняем карту в Redis с TTL 5 минут
	cardJSON, err := json.Marshal(card)
	if err != nil {
		log.Printf("⚠️ Ошибка при сериализации карты: %v", err)
		return nil, err
	}
	if err := s.cache.Set(cardID, string(cardJSON), 5*time.Minute); err != nil {
		log.Printf("⚠️ Ошибка при сохранении карты в Redis: %v", err)
	}

	// Отправляем событие в NATS
	s.natsPub.Publish("card.retrieved", []byte(fmt.Sprintf("Card retrieved: %s, UserID: %s", card.CardNumber, card.UserID)))

	// Возвращаем ответ с картой
	return &pb.GetCardResponse{
		Card: &pb.Card{
			CardId:         card.ID.Hex(),
			CardNumber:     card.CardNumber,
			CardholderName: card.CardHolder,
			ExpirationDate: card.ExpiryDate,
			Cvv:            card.CVV,
			UserId:         card.UserID,
			CardType:       card.CardType,
		},
	}, nil
}

func (s *CardServiceServer) UpdateCard(ctx context.Context, req *pb.UpdateCardRequest) (*pb.CardResponse, error) {
	collection := db.GetCollection("cards")
	oid, err := primitive.ObjectIDFromHex(req.CardId)
	if err != nil {
		return nil, err
	}

	// Обновляем карту в базе данных
	update := bson.M{
		"$set": bson.M{
			"cardholder_name": req.CardholderName,
			"expiration_date": req.ExpirationDate,
			"card_type":       req.CardType,
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return nil, err
	}

	// Получаем обновленную карту
	var card model.Card
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&card)
	if err != nil {
		return nil, err
	}

	// Сохраняем обновленную карту в Redis
	cardJSON, err := json.Marshal(card)
	if err != nil {
		log.Printf("⚠️ Ошибка при сериализации карты: %v", err)
	}
	if err := s.cache.Set(req.CardId, string(cardJSON), 5*time.Minute); err != nil {
		log.Printf("⚠️ Ошибка при сохранении обновленной карты в Redis: %v", err)
	}

	// Возвращаем обновленную карту в ответе
	return &pb.CardResponse{
		Card: &pb.Card{
			CardId:         card.ID.Hex(),
			CardNumber:     card.CardNumber,
			CardholderName: card.CardHolder,
			ExpirationDate: card.ExpiryDate,
			Cvv:            card.CVV,
			UserId:         card.UserID,
			CardType:       card.CardType,
		},
		Success: true,
		Message: "Card updated successfully",
	}, nil
}

func (s *CardServiceServer) DeleteCard(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	collection := db.GetCollection("cards")
	oid, err := primitive.ObjectIDFromHex(req.CardId)
	if err != nil {
		return nil, err
	}

	// Удаление из базы данных
	res, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return &pb.DeleteCardResponse{
			Success: false,
			Message: "Card not found",
		}, nil
	}

	// Удаление из Redis
	if err := s.cache.Del(req.CardId); err != nil {
		log.Printf("⚠️ Не удалось удалить карту из Redis: %v", err)
	}

	// Публикация события в NATS
	s.natsPub.Publish("card.deleted", []byte(fmt.Sprintf("Card deleted: %s", req.CardId)))

	return &pb.DeleteCardResponse{
		Success: true,
		Message: "Card deleted",
	}, nil
}

func (s *CardServiceServer) StartNatsConsumers() {
	err := s.natsSub.Subscribe("card.created", func(data []byte) {
		log.Printf("[NATS] ➕ card.created: %s", string(data))
	})
	if err != nil {
		log.Printf("Error subscribing to card.created: %v", err)
	}

	err = s.natsSub.Subscribe("card.retrieved", func(data []byte) {
		log.Printf("[NATS] 🔍 card.retrieved: %s", string(data))
	})
	if err != nil {
		log.Printf("Error subscribing to card.retrieved: %v", err)
	}

	err = s.natsSub.Subscribe("cards.listed", func(data []byte) {
		log.Printf("[NATS] 📋 cards.listed: %s", string(data))
	})
	if err != nil {
		log.Printf("Error subscribing to cards.listed: %v", err)
	}

	err = s.natsSub.Subscribe("card.charge", func(data []byte) {
		log.Printf("[NATS] 💳 card.charge: %s", string(data))
	})
	if err != nil {
		log.Printf("Error subscribing to card.charge: %v", err)
	}

	err = s.natsSub.Subscribe("user.registered", func(data []byte) {
		log.Printf("[NATS] 👤 user.registered: %s", string(data))
		var event struct {
			UserId string `json:"user_id"`
		}
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("Failed to unmarshal user.registered: %v", err)
			return
		}
		card := &model.Card{
			CardNumber: generateBonusCardNumber(),
			CardHolder: "New User",
			ExpiryDate: "12/30",
			CVV:        "000",
			UserID:     event.UserId,
			CardType:   "BONUS",
		}
		collection := db.GetCollection("cards")
		_, err := collection.InsertOne(context.Background(), card)
		if err != nil {
			log.Printf("Failed to create bonus card: %v", err)
		}
	})
	if err != nil {
		log.Printf("Error subscribing to user.registered: %v", err)
	}
}

func generateBonusCardNumber() string {
	prefix := "777777"
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	sb.WriteString(prefix)
	for i := 0; i < 9; i++ {
		sb.WriteString(fmt.Sprintf("%d", rand.Intn(10)))
	}
	return sb.String()
}
