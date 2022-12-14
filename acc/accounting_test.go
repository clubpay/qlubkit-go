package qacc_test

import (
	"fmt"
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

func TestAccounting(t *testing.T) {
	Convey("Accounting", t, func(c C) {
		Convey("ToInt and ToFloat", func(c C) {
			testCases := []testCase{
				{"10.32", 2, 1032, 10.32},
				{"65.01", 2, 6501, 65.01},
				{"65.10", 2, 6510, 65.10},
				{"20", 2, 2000, 20.00},
				{"0.1", 3, 100, 0.100},
				{"1.99", 2, 199, 1.99},
				{"2.01", 2, 201, 2.01},
				{"2.0001", 4, 20001, 2.0001},
				{"1.956874", 2, 196, 1.96},
				{"1.952874", 2, 195, 1.95},
				{"1.956874", 3, 1957, 1.957},
				{"3.1", 2, 310, 3.10},
				{"3.14", 3, 3140, 3.140},
				{"3", 1, 30, 3.0},
				{"0.3", 1, 3, 0.3},
				{"0.03", 2, 3, 0.03},
				{"0.003", 3, 3, 0.003},
				{"3.3", 2, 330, 3.30},
			}

			for _, tc := range testCases {
				c.So(qacc.ToIntX(qacc.ToPrecision(tc.in, tc.prec)), ShouldEqual, tc.outInt)
				c.So(qacc.ToFloatX(qacc.ToPrecision(tc.in, tc.prec)), ShouldEqual, tc.outFloat)
			}
		})

		Convey("Abs", func(c C) {
			c.So(qacc.Abs("20.0"), ShouldEqual, "20.0")
			c.So(qacc.Abs("2.00"), ShouldEqual, "2.00")
			c.So(qacc.Abs("2"), ShouldEqual, "2")
			c.So(qacc.Abs("-2"), ShouldEqual, "2")
			c.So(qacc.Abs("-2.00"), ShouldEqual, "2.00")
			c.So(qacc.Abs("-20.0"), ShouldEqual, "20.0")
			c.So(qacc.Abs("-2-0.0"), ShouldEqual, "2-0.0")
			c.So(qacc.Abs("-2-0-.0"), ShouldEqual, "2-0-.0")
			c.So(qacc.Abs("0"), ShouldEqual, "0")
			c.So(qacc.Abs(""), ShouldEqual, "")
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
			c.So(qacc.Equal("-0.1", "0.1"), ShouldBeFalse)
			c.So(qacc.Equal("", ""), ShouldBeTrue)
		})

		Convey("Equal ignore", func(c C) {
			c.So(qacc.EqualIgnore("20.0", "20", "0"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("20.0", "20.00", "0"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("20.", "20.000", "0"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("20", "20", "0"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("10", "20", "0"), ShouldBeFalse)
			c.So(qacc.EqualIgnore("123", "122", "0"), ShouldBeFalse)
			c.So(qacc.EqualIgnore("123.001", "123", "0"), ShouldBeFalse)
			c.So(qacc.EqualIgnore("240", "240.", "0"), ShouldBeTrue)

			c.So(qacc.EqualIgnore("20.02", "20", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("20.00", "20.02", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("20.00", "20.03", "0.02"), ShouldBeFalse)

			c.So(qacc.EqualIgnore("20.2", "20.00", "0.02"), ShouldBeFalse)
			c.So(qacc.EqualIgnore("20.", "20.002", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("20", "20", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("10", "20", "0.02"), ShouldBeFalse)
			c.So(qacc.EqualIgnore("123", "122", "0.02"), ShouldBeFalse)
			c.So(qacc.EqualIgnore("123.001", "123", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("240", "240.", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("1.01", "1.03", "0.02"), ShouldBeTrue)
			c.So(qacc.EqualIgnore("1.01", "1.04", "0.02"), ShouldBeFalse)
		})

		Convey("Multiply1", func(c C) {
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

		Convey("Multiply2", func(c C) {
			testCases := [][3]string{
				{"1.99", "2.01", "4.00"},
				{"2.01", "1.99", "4.00"},
				{"91.34375", "0.14453125", "13.20202637"},
			}

			for _, tc := range testCases {
				c.So(qacc.MultiplyX(tc[0], tc[1]), ShouldEqual, tc[2])
			}
		})

		Convey("Divide", func(c C) {
			testCases := [][3]string{
				{"1.99", "2.01", "0.99"},
				{"2.01", "1.99", "1.01"},
				{"91.34375", "0.14453125", "632.00000000"},
				{"1.57", "6.275", "0.250"},
			}

			for _, tc := range testCases {
				c.So(qacc.DivideX(tc[0], tc[1]), ShouldEqual, tc[2])
			}
		})

		Convey("Quotient", func(c C) {
			testCases := [][3]string{
				{"1.99", "2.01", "0"},
				{"2.01", "1.99", "1"},
				{"91.34375", "0.14453125", "632"},
				{"1.57", "6.275", "0"},
				{"25.1", "6.2", "4"},
			}

			for _, tc := range testCases {
				c.So(qacc.QuotientX(tc[0], tc[1]), ShouldEqual, tc[2])
			}
		})

		Convey("Sum", func(c C) {
			testCases := [][3]string{
				{"1.99", "2.01", "4.00"},
				{"2.01", "-1.99", "0.02"},
				{"-1.570", "6.275", "4.705"},
			}

			for _, tc := range testCases {
				c.So(qacc.SumX(tc[0], tc[1]), ShouldEqual, tc[2])
			}
		})

		Convey("Subtract", func(c C) {
			testCases := [][3]string{
				{"", "2.01", "-2.01"},
				{"x", "2.01", "-2.01"},
				{"xdw", "2.01", "-2.01"},
				{"xdw", "NaN", "0"},
				{"0", "2.31", "-2.31"},
				{"00", "2.01", "-2.01"},
				{"0.000", "2.01", "-2.010"},
				{"", "", "0"},
				{"2.01", "", "2.01"},
				{"2.01", "0", "2.01"},
				{"2.01", "00", "2.01"},
				{"2.01", "0.0", "2.01"},
				{"2.01", "0.00", "2.01"},
				{"-1.570", "6.275", "-7.845"},
			}

			for _, tc := range testCases {
				c.SoMsg(fmt.Sprintf("%s-%s=%s", tc[0], tc[1], tc[2]), qacc.SubtractX(tc[0], tc[1]), ShouldEqual, tc[2])
			}
		})

		Convey("Remainder", func(c C) {
			testCases := [][3]string{
				{"1.99", "2.01", "1.99"},
				{"2.01", "1.99", "0.02"},
				{"91.34375", "0.14453125", "0.00000000"},
				{"1.57", "6.275", "1.570"},
				{"25.1", "6.2", "0.3"},
			}

			for _, tc := range testCases {
				c.So(qacc.RemainderX(tc[0], tc[1]), ShouldEqual, tc[2])
			}
		})

		Convey("Insignify", func(c C) {
			testCases := []struct {
				input  int
				digits int
				output string
			}{
				{input: 195, digits: 2, output: "1.95"},
				{input: 2, digits: 0, output: "2"},
				{input: 2, digits: 1, output: "0.2"},
				{input: 2, digits: 2, output: "0.02"},
				{input: 2, digits: 3, output: "0.002"},
				{input: 2, digits: 4, output: "0.0002"},
				{input: 20, digits: 1, output: "2.0"},
				{input: 20, digits: 2, output: "0.20"},
				{input: 20, digits: 3, output: "0.020"},
			}

			for _, tc := range testCases {
				c.So(qacc.Insignify(tc.input, tc.digits), ShouldEqual, tc.output)
			}
		})
	})
}

func BenchmarkLength(b *testing.B) {
	for n := 0; n < b.N; n++ {
		qacc.Length(n)
	}
}

func BenchmarkLength2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		qacc.Length2(n)
	}
}
