package qkit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPhone(t *testing.T) {
	Convey("Phone", t, func(c C) {
		validPhone := "+905323214567"
		testPhone1 := "905323214567"
		testPhone2 := "00905323214567"
		testPhone3 := "+90 (532) 321-45-67"
		testPhone4 := "٩٠٥٣٢٣٢١٤٥٦٧"

		phone1, err := SanitizePhoneNumber(testPhone1)
		c.So(phone1, ShouldEqual, validPhone)
		c.So(err, ShouldBeNil)

		phone2, err := SanitizePhoneNumber(testPhone2)
		c.So(phone2, ShouldEqual, validPhone)
		c.So(err, ShouldBeNil)

		phone3, err := SanitizePhoneNumber(testPhone3)
		c.So(phone3, ShouldEqual, validPhone)
		c.So(err, ShouldBeNil)

		phone4, err := SanitizePhoneNumber(testPhone4)
		c.So(phone4, ShouldEqual, validPhone)
		c.So(err, ShouldBeNil)
	})
}
