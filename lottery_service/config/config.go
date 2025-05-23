package config

import (
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=postgres password=1234 dbname=lottery_db port=5432 sslmode=disable")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
