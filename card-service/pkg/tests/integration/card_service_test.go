package integration

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	pb "card-service/pkg/api"
	model "card-service/pkg/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var client pb.CardServiceClient
var mongoClient *mongo.Client
var cardCollection *mongo.Collection

func TestMain(m *testing.M) {
	// Подключение к gRPC
	conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC: %v", err)
	}
	client = pb.NewCardServiceClient(conn)

	// Подключение к MongoDB
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27018"))
	if err != nil {
		log.Fatalf("Ошибка создания клиента Mongo: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatalf("Ошибка подключения к MongoDB: %v", err)
	}

	// Коллекция карточек
	cardCollection = mongoClient.Database("card_service").Collection("cards")

	// Запуск всех тестов
	code := m.Run()

	// Очистка
	_ = mongoClient.Disconnect(ctx)
	os.Exit(code)
}

func TestCreateAndGetCardIntegration(t *testing.T) {
	md := metadata.Pairs("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMGRmY2M3YWItZTYyMC00ZTY3LWJhMzEtYjE2ZjQzNmIyMzI2IiwiaXNzIjoidXNlci1zZXJ2aWNlIiwiZXhwIjoxNzQ3Njc3OTg2fQ.uBpmbHLN-VJoYMw4WWqU2ChiIA5_5ZU5He_jW7litQw")

	baseCtx := context.Background()

	// === 1. Создание карточки через gRPC ===
	ctxCreate := metadata.NewOutgoingContext(baseCtx, md)
	ctxCreate, cancelCreate := context.WithTimeout(ctxCreate, 10*time.Second)
	defer cancelCreate()

	createResp, err := client.CreateCard(ctxCreate, &pb.CreateCardRequest{
		CardNumber:     "4400439988776655",
		CardholderName: "Maxat Integration",
		ExpirationDate: "12/30",
		Cvv:            "321",
		UserId:         "user-789",
		CardType:       "Kaspi Gold",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, createResp.CardId)

	// === 2. Получение карточки через gRPC ===
	ctxGet := metadata.NewOutgoingContext(baseCtx, md)
	ctxGet, cancelGet := context.WithTimeout(ctxGet, 10*time.Second)
	defer cancelGet()

	getResp, err := client.GetCard(ctxGet, &pb.GetCardRequest{
		CardId: createResp.CardId,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Maxat Integration", getResp.Card.CardholderName)

	// === 3. Проверка наличия записи в MongoDB ===
	var card model.Card
	id, err := primitive.ObjectIDFromHex(createResp.CardId)
	if err != nil {
		t.Fatalf("Неверный формат CardId: %v", err)
	}

	err = cardCollection.FindOne(ctxGet, bson.M{"_id": id}).Decode(&card)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			t.Fatalf("Карточка с ID %s не найдена в MongoDB", createResp.CardId)
		} else {
			t.Fatalf("Ошибка при поиске карточки в MongoDB: %v", err)
		}
	}

	// Дополнительная проверка полей
	assert.Equal(t, "Maxat Integration", card.CardHolder)
}
