package cache

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
)
type CacheInterface interface {
    Get(key string) (string, error)
    Set(key, value string, ttl time.Duration) error
    Del(key string) error
}


type RedisClient struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisClient(addr string) *RedisClient {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr: addr,
        DB:   0,
    })

    return &RedisClient{
        client: rdb,
        ctx:    ctx,
    }
}

func (r *RedisClient) Get(key string) (string, error) {
    return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Set(key string, value string, ttl time.Duration) error {
    return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisClient) Del(key string) error {
    return r.client.Del(r.ctx, key).Err()
}
