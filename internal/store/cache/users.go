package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/llan0/go-social/internal/store"
	"github.com/redis/go-redis/v9"
)

const UserExpTime = time.Minute

type UserStore struct {
	rdb *redis.Client
}

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	// Theres nothing in cache intitially, so check for that Redis.Nil error
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEx(ctx, cacheKey, json, UserExpTime).Err()
}
