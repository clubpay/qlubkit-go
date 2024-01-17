package qkit_test

import (
	"github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestStruct struct {
	A int   `json:"a"`
	B string `json:"b"`
}

type TestStruct2 struct {
	A string `json:"b"`
	B int   `json:"a"`
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
		Convey("CastJSON", func(c C) {
			c.SoMsg("map to struct", qkit.CastJSON[TestStruct](map[string]interface{}{"a": 1, "b": "2"}), ShouldResemble, TestStruct{1, "2"})
			c.SoMsg("struct to struct", qkit.CastJSON[TestStruct](TestStruct2{"2", 1}), ShouldResemble, TestStruct{1, "2"})
		})
		Convey("ToJSON", func(c C) {
			c.SoMsg("struct to json byte array", string(qkit.ToJSON(TestStruct{1, "2"})), ShouldEqual, `{"a":1,"b":"2"}`)
		})
		Convey("FromJSON", func(c C) {
			c.SoMsg("json byte array to struct", qkit.FromJSON[TestStruct]([]byte(`{"a":1,"b":"2"}`)), ShouldResemble, TestStruct{1, "2"})
		})
	})
	Convey("ToMap", t, func(c C) {
		converted := qkit.ToMap(TestStruct{1, "2"})
		c.SoMsg("A", converted["a"], ShouldEqual, 1)
		c.SoMsg("B", converted["b"], ShouldEqual, "2")
	})
}
