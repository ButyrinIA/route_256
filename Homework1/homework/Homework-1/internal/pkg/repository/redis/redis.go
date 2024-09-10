package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	return &Redis{
		redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "1234",
			DB:       0,
		}),
	}
}
func (r *Redis) Set(ctx context.Context, key string, value []byte) error {
	return r.client.Set(ctx, key, value, time.Minute*10).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	val, err := res.Bytes()
	if err != nil {
		return nil, err
	}
	return val, nil
}
