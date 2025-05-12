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
