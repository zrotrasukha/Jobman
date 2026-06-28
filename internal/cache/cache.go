package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zrotrasukha/jobman/internal/data"
)

var (
	ErrCacheMiss = redis.Nil
)

type Cache struct {
	client *redis.Client
}

func NewCache(addr string, password string, db int) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *Cache) GetUserForToken(context context.Context, token string) (*data.User, error) {
	val, err := c.client.Get(context, token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var user data.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Cache) SetUserForToken(ctx context.Context, token string, user *data.User, expiry time.Duration) error {
	b, err := json.Marshal(user)
	val, err := c.client.Set(ctx, token, b, expiry).Result()
	if err != nil {
		return err
	}

	if val != "OK" {
		return errors.New("failed to set user for token in cache")
	}

	return nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}
