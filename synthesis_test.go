package qkit_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	qkit "github.com/clubpay/qlubkit-go"
)

func TestCoalesce(t *testing.T) {
	type testCase[T comparable] struct {
		in  []T
		out T
	}

	Convey("String Coalesce", t, func(c C) {
		cases := []testCase[string]{
			{in: []string{""}, out: ""},
			{in: []string{"1st"}, out: "1st"},
			{in: []string{"", ""}, out: ""},
			{in: []string{"1st", ""}, out: "1st"},
			{in: []string{"1st", "2nd"}, out: "1st"},
			{in: []string{"", "2nd"}, out: "2nd"},
			{in: []string{"", "", ""}, out: ""},
			{in: []string{"1st", "", ""}, out: "1st"},
			{in: []string{"1st", "2nd", ""}, out: "1st"},
			{in: []string{"1st", "2nd", "3rd"}, out: "1st"},
			{in: []string{"1st", "", "3rd"}, out: "1st"},
			{in: []string{"", "2nd", ""}, out: "2nd"},
			{in: []string{"", "2nd", "3rd"}, out: "2nd"},
			{in: []string{"", "", "3rd"}, out: "3rd"},
		}

		for _, tc := range cases {
			c.SoMsg(fmt.Sprintf("%+v", tc.in), qkit.Coalesce(tc.in...), ShouldEqual, tc.out)
		}
	})

	Convey("Int Coalesce", t, func(c C) {
		cases := []testCase[int]{
			{in: []int{0}, out: 0},
			{in: []int{1}, out: 1},
			{in: []int{0, 0}, out: 0},
			{in: []int{1, 0}, out: 1},
			{in: []int{1, 2}, out: 1},
			{in: []int{0, 2}, out: 2},
			{in: []int{0, 0, 0}, out: 0},
			{in: []int{1, 0, 0}, out: 1},
			{in: []int{1, 2, 0}, out: 1},
			{in: []int{1, 2, 3}, out: 1},
			{in: []int{1, 0, 3}, out: 1},
			{in: []int{0, 2, 0}, out: 2},
			{in: []int{0, 2, 3}, out: 2},
			{in: []int{0, 0, 3}, out: 3},
		}

		for _, tc := range cases {
			c.SoMsg(fmt.Sprintf("%+v", tc.in), qkit.Coalesce(tc.in...), ShouldEqual, tc.out)
		}
	})

	Convey("Duration Coalesce", t, func(c C) {
		cases := []testCase[time.Duration]{
			{in: []time.Duration{0 * time.Second}, out: 0 * time.Second},
			{in: []time.Duration{1 * time.Second}, out: 1 * time.Second},
			{in: []time.Duration{0 * time.Second, 0 * time.Second}, out: 0 * time.Second},
			{in: []time.Duration{1 * time.Second, 0 * time.Second}, out: 1 * time.Second},
			{in: []time.Duration{1 * time.Second, 2 * time.Second}, out: 1 * time.Second},
			{in: []time.Duration{0 * time.Second, 2 * time.Second}, out: 2 * time.Second},
		}

		for _, tc := range cases {
			c.SoMsg(fmt.Sprintf("%+v", tc.in), qkit.Coalesce(tc.in...), ShouldEqual, tc.out)
		}
	})
}
