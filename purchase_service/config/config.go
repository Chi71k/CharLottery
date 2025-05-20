package config

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=postgres password=Eroha100! dbname=lottery_db port=5432 sslmode=disable")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
