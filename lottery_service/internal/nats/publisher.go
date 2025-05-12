package nats

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(nc *nats.Conn) *Publisher {
	return &Publisher{conn: nc}
}

func (p *Publisher) PublishLotteryCreated(lotteryID int64, prize string, availableTickets int64) {
	data := map[string]interface{}{
		"lottery_id":        lotteryID,
		"prize":             prize,
		"available_tickets": availableTickets,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		return
	}
	if err := p.conn.Publish("lottery.created", bytes); err != nil {
		log.Println("Failed to publish message:", err)
	}
}
