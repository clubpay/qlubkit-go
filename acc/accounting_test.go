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

		Convey("Equal", func(c C) {
			c.So(qacc.Equal("20.0", "20"), ShouldBeTrue)
			c.So(qacc.Equal("20.0", "20.00"), ShouldBeTrue)
			c.So(qacc.Equal("20.", "20.000"), ShouldBeTrue)
			c.So(qacc.Equal("20", "20"), ShouldBeTrue)
			c.So(qacc.Equal("10", "20"), ShouldBeFalse)
			c.So(qacc.Equal("123", "122"), ShouldBeFalse)
			c.So(qacc.Equal("123.001", "123"), ShouldBeFalse)
			c.So(qacc.Equal("240", "240."), ShouldBeTrue)
		})

		Convey("Multiply", func(c C) {
			c.So(qacc.MultiplyX("20.0", "20"), ShouldEqual, "400.0")
			c.So(qacc.MultiplyX("2.0", "2"), ShouldEqual, "4.0")
			c.So(qacc.MultiplyX("2", "2"), ShouldEqual, "4")
			c.So(qacc.MultiplyX("6.10", "3"), ShouldEqual, "18.30")
			c.So(qacc.MultiplyX("0", "1"), ShouldEqual, "0")
			c.So(qacc.MultiplyX("1.0", "0"), ShouldEqual, "0.0")
			c.So(qacc.MultiplyX("25", "4"), ShouldEqual, "100")
			c.So(qacc.MultiplyX("1", "23.32324"), ShouldEqual, "23.32324")
			c.So(qacc.SumX("", qacc.MultiplyX("2", "2")), ShouldEqual, "4")
		})
	})
}
