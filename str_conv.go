package qkit

import (
	"reflect"
	"strconv"
	"unsafe"
)

func StrToFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)

	return v
}

func StrToFloat32(s string) float32 {
	v, _ := strconv.ParseFloat(s, 32)

	return float32(v)
}

func StrToInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)

	return v
}

func StrToInt32(s string) int32 {
	v, _ := strconv.ParseInt(s, 10, 32)

	return int32(v)
}

func StrToUInt64(s string) uint64 {
	v, _ := strconv.ParseInt(s, 10, 64)

	return uint64(v)
}

func StrToUInt32(s string) uint32 {
	v, _ := strconv.ParseInt(s, 10, 32)

	return uint32(v)
}

func StrToInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 32)

	return int(v)
}

func StrToUInt(s string) uint {
	v, _ := strconv.ParseInt(s, 10, 32)

	return uint(v)
}

func Int64ToStr(x int64) string {
	return strconv.FormatInt(x, 10)
}

func Int32ToStr(x int32) string {
	return strconv.FormatInt(int64(x), 10)
}

func UInt64ToStr(x uint64) string {
	return strconv.FormatUint(x, 10)
}

func UInt32ToStr(x uint32) string {
	return strconv.FormatUint(uint64(x), 10)
}

func IntToStr(x int) string {
	return strconv.FormatUint(uint64(x), 10)
}

// ByteToStr converts byte slice to a string without memory allocation.
// Note it may break if string and/or slice header will change
// in the future go versions.
func ByteToStr(bts []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bts))

	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = bh.Data
	sh.Len = bh.Len

	return s
}

// B2S is alias for ByteToStr.
func B2S(bts []byte) string {
	return ByteToStr(bts)
}

// StrToByte converts string to a byte slice without memory allocation.
// Note it may break if string and/or slice header will change
// in the future go versions.
func StrToByte(str string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len

	return b
}

// S2B is alias for StrToByte.
func S2B(str string) []byte {
	return StrToByte(str)
}

func Float64ToStr(x float64) string {
	return strconv.FormatFloat(x, 'f', -1, 64)
}

func Float32ToStr(x float32) string {
	return strconv.FormatFloat(float64(x), 'f', -1, 32)
}
