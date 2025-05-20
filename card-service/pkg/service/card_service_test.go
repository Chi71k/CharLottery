package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	pb "card-service/pkg/api"
	model "card-service/pkg/db/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/stretchr/testify/assert"
)

// mockPublisher and mockSubscriber simulate NATS
type mockPublisher struct{}
type mockSubscriber struct{}

func (m *mockPublisher) Publish(subject string, data []byte) error { return nil }
func (m *mockSubscriber) Subscribe(subject string, handler func(data []byte)) error {
	return nil
}

// mockRedis implements Redis interface
type mockRedis struct {
	data map[string]string
}

func (m *mockRedis) Get(key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}
func (m *mockRedis) Set(key string, value string, ttl time.Duration) error {
	m.data[key] = value
	return nil
}
func (m *mockRedis) Del(key string) error {
	delete(m.data, key)
	return nil
}

// TestCardValidation checks model.Card validation
func TestCardValidation(t *testing.T) {
	card := &model.Card{
		CardNumber: "1234567890123456",
		CardHolder: "Maxat",
		ExpiryDate: "12/30",
		CVV:        "123",
		UserID:     "user-123",
		CardType:   "VISA",
	}

	err := card.Validate()
	assert.Nil(t, err)
}

// TestGetCardFromRedis checks if cached card can be retrieved
func TestGetCardFromRedis(t *testing.T) {
	redis := &mockRedis{data: make(map[string]string)}
	srv := &CardServiceServer{
		natsPub: &mockPublisher{},
		natsSub: &mockSubscriber{},
		cache:   redis,
	}

	cardID := primitive.NewObjectID() // генерируем ObjectID
	card := &model.Card{
		ID:         cardID,
		CardNumber: "4400430011223344",
		CardHolder: "Nurbibi",
		ExpiryDate: "11/30",
		CVV:        "123",
		UserID:     "user-456",
		CardType:   "Kaspi Gold",
	}

	cardJSON, err := json.Marshal(card)
	assert.NoError(t, err)

	redis.Set(cardID.Hex(), string(cardJSON), 5*time.Minute)

	resp, err := srv.GetCard(context.Background(), &pb.GetCardRequest{
		CardId: cardID.Hex(),
	})
	assert.NoError(t, err)
	assert.Equal(t, "4400430011223344", resp.Card.CardNumber)
	assert.Equal(t, "Nurbibi", resp.Card.CardholderName)
}
