package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type storeRedis struct {
	redisClient *redis.Client
	ctx         context.Context
	ttl         time.Duration
}

var _ Store = (*storeRedis)(nil)

func NewRedisStore(redisClient *redis.Client, ctx context.Context) Store {
	sr := storeRedis{
		redisClient: redisClient,
		ctx:         ctx,
	}
	return &sr
}
func (s *storeRedis) SetTTL(duration time.Duration) {
	s.ttl = duration
}

func (s *storeRedis) GetValue(key string) ([]byte, error) {
	return s.redisClient.Get(context.TODO(), key).Bytes()
}

func (s *storeRedis) SetValue(key string, value []byte) error {
	return s.redisClient.Set(context.TODO(), key, value, s.ttl).Err()
}
