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
	_, err = collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}
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
