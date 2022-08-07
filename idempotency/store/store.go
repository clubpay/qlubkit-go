package store

import "time"

type Store interface {
	SetTTL(duration time.Duration)
	GetValue(key string) ([]byte, error)
	SetValue(key string, value []byte) error
}
