package qkit_test

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	qkit "github.com/clubpay/qlubkit-go"
)

var (
	errInternal  = errors.New("internal error")
	errRuntime   = errors.New("runtime error")
	errEndOfFile = errors.New("end of file")

	errWrappedInternal = fmt.Errorf("wrapped error: %w", errInternal)
)

func TestWrapError(t *testing.T) {
	Convey("WrapError", t, func(c C) {
		we := qkit.WrapError(errRuntime, errEndOfFile)
		c.So(errors.Is(we, errRuntime), ShouldBeTrue)
		c.So(errors.Is(we, errEndOfFile), ShouldBeTrue)
		c.So(we.Error(), ShouldEqual, "runtime error: end of file")

		wwe := qkit.WrapError(errInternal, we)
		c.So(errors.Is(wwe, errInternal), ShouldBeTrue)
		c.So(errors.Is(wwe, errRuntime), ShouldBeTrue)
		c.So(errors.Is(wwe, errEndOfFile), ShouldBeTrue)
		c.So(wwe.Error(), ShouldEqual, "internal error: runtime error: end of file")

		wwi := qkit.WrapError(errWrappedInternal, we)
		c.So(errors.Is(wwi, errWrappedInternal), ShouldBeTrue)
		c.So(errors.Is(wwi, errInternal), ShouldBeTrue)
		c.So(errors.Is(wwi, errRuntime), ShouldBeTrue)
		c.So(errors.Is(wwi, errEndOfFile), ShouldBeTrue)
		c.So(wwi.Error(), ShouldEqual, "wrapped error: internal error: runtime error: end of file")
	})
}
