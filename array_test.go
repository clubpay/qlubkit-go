package qkit_test

import (
	"testing"

	qkit "github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPaginate(t *testing.T) {
	Convey("Paginate", t, func(c C) {
		arr := []int{1, 2, 3, 4, 5, 6, 7, 8}
		idx := 0
		expected := [][]int{
			{0, 3},
			{3, 6},
			{6, 8},
		}
		_ = qkit.Paginate(arr, 3, func(start, end int) error {
			c.So(start, ShouldEqual, expected[idx][0])
			c.So(end, ShouldEqual, expected[idx][1])
			idx++
			return nil
		})

		idx = 0
		expected = [][]int{
			{0, 7},
			{7, 8},
		}
		_ = qkit.Paginate(arr, 7, func(start, end int) error {
			c.So(start, ShouldEqual, expected[idx][0])
			c.So(end, ShouldEqual, expected[idx][1])
			idx++
			return nil
		})
		idx = 0
		expected = [][]int{
			{0, 8},
		}
		_ = qkit.Paginate(arr, 8, func(start, end int) error {
			c.So(start, ShouldEqual, expected[idx][0])
			c.So(end, ShouldEqual, expected[idx][1])
			idx++
			return nil
		})

		idx = 0
		expected = [][]int{
			{0, 8},
		}
		_ = qkit.Paginate(arr, 20, func(start, end int) error {
			c.So(start, ShouldEqual, expected[idx][0])
			c.So(end, ShouldEqual, expected[idx][1])
			idx++
			return nil
		})
	})
}
