package service_test

import (
	"context"
	"testing"
	"user-service/pkg/cache"
	"user-service/pkg/natswrap"
	"user-service/pkg/repository"
	"user-service/pkg/service"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupMongo(t *testing.T) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27018"))
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func setupRedis() *cache.RedisClient {
	return cache.NewRedisClient("localhost:6379")
}

func setupNATS(t *testing.T) natswrap.Publisher {
    natsClient, err := natswrap.NewNatsClient(nats.DefaultURL)
    if err != nil {
        t.Fatal(err)
    }
    return natsClient
}


func TestUserService_Register_Integration(t *testing.T) {
	mongoClient := setupMongo(t)
	defer mongoClient.Disconnect(context.Background())

	repo := repository.NewUserRepository(mongoClient)
	redisClient := setupRedis()
	natsClient := setupNATS(t)

	svc := service.NewUserService(repo, natsClient, redisClient)

	username := "intuser"
	email := "maksatkarzhaubaev91@gmail.com"
	password := "intpass"

	// Удалим пользователя, если он был с предыдущего теста
	repo.DeleteUser(username)

	userID, err := svc.Register(username, password, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	user, err := repo.GetUserByID(userID)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
}
