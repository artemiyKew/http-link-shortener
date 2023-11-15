package redisdb

import "github.com/redis/go-redis/v9"

type RedisDB struct {
	RDB *redis.Client
}

func NewRDB() *RedisDB {
	return &RedisDB{
		RDB: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		}),
	}
}
