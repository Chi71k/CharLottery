package model

type Lottery struct {
	ID          int64 `gorm:"primaryKey;autoIncrement"`
	Title       string
	Description string
	Prize       string
	Status      string
}
