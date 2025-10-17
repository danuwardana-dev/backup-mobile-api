package config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Address  string `envconfig:"REDIS_ADDRESS" required:"true"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB  int    `envconfig:"REDIS_DB" required:"true"`
}

func (rd Redis) RedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rd.Address,
		Password: rd.Password,
		DB:       rd.RedisDB,
	})
	if client.Ping(context.Background()).Err() != nil {
		return nil, client.Ping(context.Background()).Err()
	}
	return client, nil
}
func LoadRedis(redis Redis) *Redis {
	return &redis
}
