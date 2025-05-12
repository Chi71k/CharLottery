package usecase

import (
	"github.com/CharLottery/purchase_service/internal/adapter/postgres"
	"github.com/CharLottery/purchase_service/internal/model"
	"github.com/lib/pq"
)

type PurchaseUsecase struct {
	repo      *postgres.PurchaseRepository
	publisher TicketEventPublisher
}

func NewPurchaseUsecase(repo *postgres.PurchaseRepository, publisher TicketEventPublisher) *PurchaseUsecase {
	return &PurchaseUsecase{repo: repo, publisher: publisher}
}

type TicketEventPublisher interface {
	PublishTicketBought(lotteryID int64)
}

func (uc *PurchaseUsecase) BuyTicket(userID, lotteryID int64, numbers []int32) (int64, error) {
	numbersArray := pq.Int32Array(numbers)

	purchase := &model.Purchase{
		UserID:    userID,
		LotteryID: lotteryID,
		Numbers:   numbersArray,
	}

	if err := uc.repo.Create(purchase); err != nil {
		return 0, err
	}

	uc.publisher.PublishTicketBought(lotteryID)

	return purchase.ID, nil
}
