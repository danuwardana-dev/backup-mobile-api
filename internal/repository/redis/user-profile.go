package redis

import (
	"backend-mobile-api/model/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func (r *Redis) SetDataUpdateUserProfile(ctx context.Context, uuidKey string, req *dto.ProfileUpdateRequest, duration time.Duration) error {
	key := fmt.Sprintf("%s:PROFILE_UPDATE", uuidKey)
	jsonData, _ := json.Marshal(*req)
	return r.client.Set(ctx, key, jsonData, duration).Err()
}
func (r *Redis) GetDataUpdateProfile(ctx context.Context, uuidKey string) (*dto.ProfileUpdateRequest, error) {
	var (
		key   = fmt.Sprintf("%s:PROFILE_UPDATE", uuidKey)
		value dto.ProfileUpdateRequest
	)
	strJsonValue, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	err = json.Unmarshal([]byte(strJsonValue), &value)
	if err != nil {
		return nil, err
	}
	return &value, nil
}
