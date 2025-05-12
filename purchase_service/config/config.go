package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=postgres password=Eroha100! dbname=lottery_db port=5432 sslmode=disable")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
