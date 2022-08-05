package idempotency

import (
	"encoding/json"
	"time"

	"github.com/clubpay/qlubkit-go/global"
	"github.com/clubpay/qlubkit-go/idempotency/store"
	"golang.org/x/sync/singleflight"
)

type Data struct {
	Status int               `json:"status"`
	Body   []byte            `json:"body"`
	Header map[string]string `json:"hdr"`
}

type Idempotency struct {
	ttl   time.Duration
	store store.Store
}

var sf singleflight.Group

type Option func(*Idempotency)

func New(opts ...Option) *Idempotency {
	idm := &Idempotency{}
	for _, opt := range opts {
		opt(idm)
	}
	// After Options
	ttl := global.DefaultIdempotancyTtl
	if idm.ttl != time.Duration(0) {
		ttl = idm.ttl
	}
	idm.store.SetTTL(ttl)
	return idm
}

// Set Sets the data related with the key
func (i *Idempotency) Set(key string, data *Data) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, sfErr, _ := sf.Do(key, func() (interface{}, error) {
		errSetValue := i.store.SetValue(key, rawData)
		return nil, errSetValue
	})
	return sfErr
}

// Check Checks if idempotent and returns the data related with the key
func (i *Idempotency) Check(key string) (*Data, error) {
	sfV, sfErr, _ := sf.Do(key, func() (interface{}, error) {
		return i.store.GetValue(key)
	})
	if sfErr != nil {
		return nil, sfErr
	}
	rawData := sfV.([]byte)
	if len(rawData) == 0 {
		return nil, nil
	}
	data := &Data{}
	err := json.Unmarshal(rawData, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
