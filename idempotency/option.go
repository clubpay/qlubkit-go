package idempotency

import (
	"time"

	"github.com/clubpay/qlubkit-go/idempotency/store"
)

func WithTTL(ttl time.Duration) func(*Idempotency) {
	return func(i *Idempotency) {
		i.ttl = ttl
	}
}

func WithStore(store store.Store) func(*Idempotency) {
	return func(i *Idempotency) {
		i.store = store
	}
}
