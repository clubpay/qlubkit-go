package qkit_test

import (
	"testing"

	qkit "github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStrConv(t *testing.T) {
	Convey("Accounting", t, func(c C) {
		Convey("StrToFloat64", func(c C) {
			type testCase struct {
				in       string
				outFloat float64
			}

			testCases := []testCase{
				{"10.32", 10.32},
				{"65.01", 65.01},
			}

			for _, tc := range testCases {
				c.SoMsg(tc.in, qkit.StrToFloat64(tc.in), ShouldEqual, tc.outFloat)
			}
		})
	})
}
