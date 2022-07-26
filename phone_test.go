package qkit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPhone(t *testing.T) {
	Convey("Phone", t, func(c C) {
		validPhoneWithPlus := "+905323214567"
		validPhone := "905323214567"

		testPhone1 := "+905323214567"
		testPhone2 := "00905323214567"
		testPhone3 := "+90 (532) 321-45-67"
		testPhone4 := "0905323214567"
		testPhone5 := "٩٠٥٣٢٣٢١٤٥٦٧"
		testPhone6 := "+00905323214567"
		testPhone7 := "00+905323214567"

		phone1, err := SanitizePhoneNumber(testPhone1)
		c.So(phone1, ShouldEqual, validPhoneWithPlus)
		c.So(err, ShouldBeNil)

		phone2, err := SanitizePhoneNumber(testPhone2)
		c.So(phone2, ShouldEqual, validPhoneWithPlus)
		c.So(err, ShouldBeNil)

		phone3, err := SanitizePhoneNumber(testPhone3)
		c.So(phone3, ShouldEqual, validPhoneWithPlus)
		c.So(err, ShouldBeNil)

		phone4, err := SanitizePhoneNumber(testPhone4)
		c.So(phone4, ShouldEqual, validPhone)
		c.So(err, ShouldBeNil)

		phone5, err := SanitizePhoneNumber(testPhone5)
		c.So(phone5, ShouldEqual, validPhone)
		c.So(err, ShouldBeNil)

		phone6, err := SanitizePhoneNumber(testPhone6)
		c.So(phone6, ShouldEqual, validPhoneWithPlus)
		c.So(err, ShouldBeNil)

		phone7, err := SanitizePhoneNumber(testPhone7)
		c.So(phone7, ShouldEqual, validPhoneWithPlus)
		c.So(err, ShouldBeNil)
	})
}
