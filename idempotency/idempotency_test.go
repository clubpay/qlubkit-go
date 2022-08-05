package idempotency_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/clubpay/qlubkit-go/idempotency"
	"github.com/clubpay/qlubkit-go/idempotency/store"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIdempotencyCheck(t *testing.T) {
	Convey("Idempotency check with ristretto store", t, func(c C) {
		data := idempotency.Data{
			Status: http.StatusOK,
			Body:   []byte("{\"Result\":\"Payment Received.\"}"),
			Header: map[string]string{"hdrKey1": "Value1", "hdrKey2": "Value2"},
		}
		key := "asd"
		idm := idempotency.New(
			idempotency.WithStore(store.NewRistretto()),
			idempotency.WithTTL(1*time.Minute),
		)
		res, err := idm.Check(key)
		c.So(res, ShouldBeNil)
		c.So(err, ShouldBeNil)
		err = idm.Set(key, &data)
		c.So(err, ShouldBeNil)
		idm = idempotency.New(
			idempotency.WithStore(store.NewRistretto()),
			idempotency.WithTTL(1*time.Minute),
		)
		res, err = idm.Check(key)
		c.So(res, ShouldNotBeNil)
		c.So(err, ShouldBeNil)
		c.So(res.Body, ShouldResemble, data.Body)
		c.So(res.Status, ShouldEqual, data.Status)
		c.So(res.Header, ShouldResemble, data.Header)
	})
}
