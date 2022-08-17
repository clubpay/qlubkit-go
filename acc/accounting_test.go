package qacc_test

import (
	"testing"

	qacc "github.com/clubpay/qlubkit-go/acc"
	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	in       string
	prec     int
	outInt   int
	outFloat float64
}

func TestConvert(t *testing.T) {
	Convey("Converts", t, func(c C) {
		Convey("ToInt and ToFloat", func(c C) {
			testCases := []testCase{
				{"10.32", 2, 1032, 10.32},
				{"65.01", 2, 6501, 65.01},
				{"65.10", 2, 6510, 65.10},
				{"20", 2, 2000, 20.00},
				{"0.1", 3, 100, 0.100},
			}

			for _, tc := range testCases {
				c.So(qacc.ToIntX(qacc.ToPrecision(tc.in, tc.prec)), ShouldEqual, tc.outInt)
				c.So(qacc.ToFloatX(qacc.ToPrecision(tc.in, tc.prec)), ShouldEqual, tc.outFloat)
			}
		})

		Convey("FromInt", func(c C) {
			c.So(qacc.FromInt(1201, 2), ShouldEqual, "12.01")
			c.So(qacc.FromInt(2010, 2), ShouldEqual, "20.10")
			c.So(qacc.FromInt(201, 2), ShouldEqual, "2.01")
			c.So(qacc.FromInt(2, 2), ShouldEqual, "0.02")
			c.So(qacc.FromInt(2, 1), ShouldEqual, "0.2")
			c.So(qacc.FromInt(2, 0), ShouldEqual, "2")
		})
	})
}
