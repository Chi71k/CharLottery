package usecase

import (
	"github.com/CharLottery/lottery_service/internal/adapter/postgres"
	"github.com/CharLottery/lottery_service/internal/model"
)

type LotteryUsecase struct {
	repo *postgres.LotteryRepository
}

func NewLotteryUsecase(repo *postgres.LotteryRepository) *LotteryUsecase {
	return &LotteryUsecase{repo: repo}
}

func (uc *LotteryUsecase) CreateLottery(l *model.Lottery) error {
	return uc.repo.CreateLottery(l)
}

func (uc *LotteryUsecase) GetLottery(id int64) (*model.Lottery, error) {
	return uc.repo.GetLottery(id)
}

func (uc *LotteryUsecase) ListLotteries() ([]model.Lottery, error) {
	return uc.repo.ListLotteries()
}
