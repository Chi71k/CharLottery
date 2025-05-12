package nats

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/CharLottery/purchase_service/internal/usecase"
	"github.com/nats-io/nats.go"
)

func generateRandomNumbers() []int32 {
	rand.Seed(time.Now().UnixNano())
	numSet := make(map[int32]struct{})
	for len(numSet) < 5 {
		n := int32(rand.Intn(49) + 1)
		numSet[n] = struct{}{}
	}
	numbers := make([]int32, 0, 5)
	for n := range numSet {
		numbers = append(numbers, n)
	}
	return numbers
}

func SubscribeToLotteryCreated(nc *nats.Conn, uc *usecase.PurchaseUsecase) {
	nc.Subscribe("lottery.created", func(msg *nats.Msg) {
		var data struct {
			LotteryID int64 `json:"lottery_id"`
		}
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Println("Failed to unmarshal:", err)
			return
		}
		log.Println("Received lottery.created:", data.LotteryID)

		numbers := generateRandomNumbers()
		if _, err := uc.BuyTicket(777, data.LotteryID, numbers); err != nil {
			log.Println("Failed to auto-buy ticket:", err)
		}
	})
}
