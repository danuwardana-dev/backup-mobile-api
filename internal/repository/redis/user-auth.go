package redis

import (
	"backend-mobile-api/app/config"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
	config *config.Root
}

func NewRedis(client *redis.Client, config *config.Root) *Redis {
	return &Redis{
		client: client,
		config: config,
	}
}

type Tag string

func (r *Redis) SetOtp(ctx context.Context, otp string, uuidKey string, value string, duration time.Duration) error {
	key := fmt.Sprintf("%s:%s", uuidKey, otp)
	return r.client.Set(ctx, key, value, duration).Err()
}
func (r *Redis) GetOtp(ctx context.Context, otp string, uuidKey string) (string, error) {
	key := fmt.Sprintf("%s:%s", uuidKey, otp)
	strValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return strValue, nil
}
func (r *Redis) DeleteOtp(ctx context.Context, otp string, uuidKey string) error {
	key := fmt.Sprintf("%s:%s", uuidKey, otp)
	return r.client.Del(ctx, key).Err()
}
func (r *Redis) OtpIsExist(ctx context.Context, otp string, uuidKey string) (bool, error) {
	key := fmt.Sprintf("%s:%s", uuidKey, otp)
	result, err := r.client.Exists(ctx, key).Result()
	return result == 1, err
}
func (r *Redis) SetAccessKey(ctx context.Context, uuid string, value string, duration time.Duration) error {
	return r.client.Set(ctx, uuid, value, duration).Err()
}
func (r *Redis) GetAccessKey(ctx context.Context, uuid string) (string, error) {
	strValue, err := r.client.Get(ctx, uuid).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return strValue, nil
}
func (r *Redis) DeleteAccessKey(ctx context.Context, uuid string) error {
	return r.client.Del(ctx, uuid).Err()
}

func (r *Redis) SetBlaclistJwt(ctx context.Context, jwt string, duration time.Duration) error {
	return r.client.Set(ctx, jwt, "active", duration).Err()
}
func (r *Redis) GetBlaclistJwt(ctx context.Context, jwt string) (string, error) {
	strValue, err := r.client.Get(ctx, jwt).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return strValue, nil
}
func (r *Redis) GetXNonce(ctx context.Context, uuid string) (string, error) {
	key := fmt.Sprintf("%s:%s", uuid, "xsession")
	strValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return strValue, nil
}
func (r *Redis) SetXNONCE(ctx context.Context, uuid string) error {
	key := fmt.Sprintf("%s:%s", uuid, "xsession")
	duration := r.config.App.XsessionExpire
	return r.client.Set(ctx, key, "active", duration).Err()
}
