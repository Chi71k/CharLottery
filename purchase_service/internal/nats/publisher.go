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

func (p *Publisher) PublishTicketBought(lotteryID int64) {
	data := map[string]interface{}{
		"lottery_id": lotteryID,
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal ticket bought message:", err)
		return
	}
	if err := p.conn.Publish("ticket.bought", bytes); err != nil {
		log.Println("Failed to publish ticket bought:", err)
	}
}
