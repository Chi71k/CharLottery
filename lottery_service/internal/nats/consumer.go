package nats

import (
	"encoding/json"
	"log"

	"github.com/CharLottery/lottery_service/internal/usecase"
	"github.com/nats-io/nats.go"
)

func SubscribeToTicketBought(nc *nats.Conn, uc *usecase.LotteryUsecase) {
	nc.Subscribe("ticket.bought", func(msg *nats.Msg) {
		var data struct {
			LotteryID int64 `json:"lottery_id"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Println("Failed to unmarshal ticket.bought:", err)
			return
		}

		log.Printf("Received ticket.bought for lottery %d\n", data.LotteryID)
		if err := uc.DecreaseAvailableTickets(data.LotteryID); err != nil {
			log.Println("Failed to decrease available tickets:", err)
		}
	})
}
