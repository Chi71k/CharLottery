package postgres

import (
	"github.com/CharLottery/lottery_service/internal/model"
	"gorm.io/gorm"
)

type LotteryRepository struct {
	db *gorm.DB
}

func NewLotteryRepository(db *gorm.DB) *LotteryRepository {
	return &LotteryRepository{db: db}
}

func (r *LotteryRepository) CreateLottery(lottery *model.Lottery) error {
	return r.db.Create(lottery).Error
}

func (r *LotteryRepository) GetLottery(id int64) (*model.Lottery, error) {
	var lottery model.Lottery
	if err := r.db.First(&lottery, id).Error; err != nil {
		return nil, err
	}
	return &lottery, nil
}

func (r *LotteryRepository) ListLotteries() ([]model.Lottery, error) {
	var lotteries []model.Lottery
	if err := r.db.Find(&lotteries).Error; err != nil {
		return nil, err
	}
	return lotteries, nil
}

func (r *LotteryRepository) DecreaseAvailableTickets(id int64) error {
	return r.db.Model(&model.Lottery{}).Where("id = ? AND available_tickets > 0", id).
		UpdateColumn("available_tickets", gorm.Expr("available_tickets - 1")).Error
}
