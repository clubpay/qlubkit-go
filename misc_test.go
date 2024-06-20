package qkit_test

import (
	"github.com/clubpay/qlubkit-go"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

func TestMisc(t *testing.T) {
	Convey("Protect", t, func(c C) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go qkit.Protect(func() {
			defer wg.Done()
			panic("test")
		})
		wg.Wait()
		c.SoMsg("this line should execute", 1, ShouldEqual, 1)
	})
}
