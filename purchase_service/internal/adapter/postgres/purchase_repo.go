package postgres

import (
	"github.com/CharLottery/purchase_service/internal/model"
	"gorm.io/gorm"
)

type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(purchase *model.Purchase) error {
	return r.db.Create(purchase).Error
}

func (r *PurchaseRepository) ListByUser(userID int64) ([]model.Purchase, error) {
	var purchases []model.Purchase
	err := r.db.Where("user_id = ?", userID).Find(&purchases).Error
	return purchases, err
}
