package db

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    client   *mongo.Client
    database *mongo.Database
)

func InitMongoDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    database = client.Database("card_service")
    log.Println("Connected to MongoDB")
}

func GetCollection(name string) *mongo.Collection {
    return database.Collection(name)
}
