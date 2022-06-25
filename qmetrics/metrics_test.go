package qmetrics_test

import (
	"context"
	"testing"

	"github.com/clubpay/qlubkit-go/qmetrics"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetrics(t *testing.T) {
	Convey("Metrics", t, func(c C) {

		met, err := qmetrics.New(
			qmetrics.WithOTLP("localhost:4318"),
			//qmetrics.WithPrometheus(2222),
		)
		c.So(err, ShouldBeNil)
		c.So(met, ShouldNotBeNil)
		m := qmetrics.Meter("x")
		c.So(m, ShouldNotBeNil)
		cnt, err := m.SyncInt64().Counter("x")
		c.So(err, ShouldBeNil)
		c.So(cnt, ShouldNotBeNil)
		cnt.Add(context.TODO(), 1)
	})
}
