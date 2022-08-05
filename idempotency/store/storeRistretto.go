package store

import (
	"errors"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

var lock = &sync.Mutex{}

var (
	instance storeRistretto
)

type storeRistretto struct {
	c   *ristretto.Cache
	ttl time.Duration
}

var _ Store = (*storeRistretto)(nil)

func NewRistretto() Store {
	lock.Lock()
	defer lock.Unlock()

	if instance.c == nil {
		rc, _ := ristretto.NewCache(
			&ristretto.Config{
				MaxCost:     1 << 30, // 1GB
				NumCounters: 1e7,     // 10M
				BufferItems: 64,
			},
		)
		instance = storeRistretto{
			c: rc,
		}
	}

	return &instance
}

func (s *storeRistretto) SetTTL(duration time.Duration) {
	s.ttl = duration
}

func (s *storeRistretto) GetValue(key string) ([]byte, error) {
	oldData, found := s.c.Get(key)
	if !found {
		return nil, nil
	}
	return oldData.([]byte), nil
}

func (s *storeRistretto) SetValue(key string, value []byte) error {
	res := s.c.SetWithTTL(key, value, 1, s.ttl)
	s.c.Wait()
	if !res {
		return errors.New("Value cannot be set")
	}
	return nil
}
