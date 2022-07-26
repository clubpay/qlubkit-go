package qkit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPhone(t *testing.T) {
	Convey("Phone", t, func(c C) {
		validPhoneWithPlus := "+905323214567"
		validPhone := "905323214567"

		testPhones := make(map[string]string)

		testPhones["+905323214567"] = validPhoneWithPlus
		testPhones["00905323214567"] = validPhoneWithPlus
		testPhones["+90 (532) 321-45-67"] = validPhoneWithPlus
		testPhones["0905323214567"] = validPhone
		testPhones["٩٠٥٣٢٣٢١٤٥٦٧"] = validPhone
		testPhones["+00905323214567"] = validPhoneWithPlus
		testPhones["00+905323214567"] = validPhoneWithPlus

		for k, v := range testPhones {
			phone, err := SanitizePhoneNumber(k)
			c.So(phone, ShouldEqual, v)
			c.So(err, ShouldBeNil)
		}
	})
}
