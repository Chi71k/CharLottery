package model

import "github.com/lib/pq"

type Purchase struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64
	LotteryID int64
	Numbers   pq.Int32Array `gorm:"type:int[]"`
}
