package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongo(uri string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("Ошибка подключения к MongoDB: %v", err)
    }

    // Проверка соединения
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("Ошибка пинга MongoDB: %v", err)
    }

    MongoClient = client
    log.Println("Успешное подключение к MongoDB")
}
