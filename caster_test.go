package qkit_test

import (
	"github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestStruct struct {
	A int
	B string
}

func TestCaster(t *testing.T) {
	Convey("Cast", t, func(c C) {
		Convey("Normal", func(c C) {
			var any interface{} = 1
			c.SoMsg("any to int", qkit.Cast[int](any), ShouldEqual, 1)
		})
		Convey("Zero value conversions", func(c C) {
			c.SoMsg("int to string", qkit.Cast[string](1), ShouldEqual, "")
			c.SoMsg("string to int", qkit.Cast[int]("1"), ShouldEqual, 0)
		})
		Convey("CastPartial", func(c C) {
			c.SoMsg("map to struct", qkit.CastPartial[TestStruct](map[string]interface{}{"A": 1, "B": "2"}), ShouldResemble, TestStruct{1, "2"})
		})
	})
	Convey("ToMap", t, func(c C) {
		converted := qkit.ToMap(TestStruct{1, "2"})
		c.SoMsg("A", converted["A"], ShouldEqual, 1)
		c.SoMsg("B", converted["B"], ShouldEqual, "2")
	})
	Convey("FromBytes", t, func(c C) {
		c.SoMsg("FromBytes", qkit.FromBytes[TestStruct]([]byte(`{"A":1,"B":"2"}`)), ShouldResemble, TestStruct{1, "2"})
	})
}
