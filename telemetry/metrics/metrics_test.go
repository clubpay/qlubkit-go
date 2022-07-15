package qmetrics_test

import (
	"context"
	"testing"
	"time"

	qmetrics "github.com/clubpay/qlubkit-go/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetrics(t *testing.T) {
	Convey("Metrics", t, func(c C) {

		met, err := qmetrics.New(
			qmetrics.WithOTLP("localhost:4318"),
		)
		c.So(err, ShouldBeNil)
		c.So(met, ShouldNotBeNil)
		m := qmetrics.Meter("x")
		c.So(m, ShouldNotBeNil)
		g, err := m.AsyncInt64().Gauge("g1")
		c.So(err, ShouldBeNil)
		c.So(g, ShouldNotBeNil)
		for i := 0; i < 10; i++ {
			g.Observe(context.TODO(), 40)
			time.Sleep(time.Second)
		}

	})
}
