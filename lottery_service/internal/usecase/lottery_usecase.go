package usecase

import (
	"github.com/CharLottery/lottery_service/internal/adapter/postgres"
	"github.com/CharLottery/lottery_service/internal/model"
)

type LotteryUsecase struct {
	repo      *postgres.LotteryRepository
	publisher LotteryEventPublisher
}

func NewLotteryUsecase(repo *postgres.LotteryRepository, publisher LotteryEventPublisher) *LotteryUsecase {
	return &LotteryUsecase{repo: repo, publisher: publisher}
}

func (uc *LotteryUsecase) CreateLottery(l *model.Lottery) error {
	if err := uc.repo.CreateLottery(l); err != nil {
		return err
	}

	uc.publisher.PublishLotteryCreated(l.ID, l.Prize, l.AvailableTickets)

	return nil
}

func (uc *LotteryUsecase) GetLottery(id int64) (*model.Lottery, error) {
	return uc.repo.GetLottery(id)
}

func (uc *LotteryUsecase) ListLotteries() ([]model.Lottery, error) {
	return uc.repo.ListLotteries()
}

func (uc *LotteryUsecase) DecreaseAvailableTickets(lotteryID int64) error {
	return uc.repo.DecreaseAvailableTickets(lotteryID)
}
