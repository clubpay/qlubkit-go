package qacc

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	qkit "github.com/clubpay/qlubkit-go"
)

var currPrecision = map[string]int{
	"JPY": 0,
	"IRR": 0,
}

func Precision(curr string) int {
	p, ok := currPrecision[strings.ToUpper(curr)]
	if !ok {
		return 2
	}

	return p
}

func SumX(a ...string) string {
	s, err := Sum(a...)
	if err != nil {
		panic(err)
	}

	return s
}

func Sum(a ...string) (string, error) {
	var (
		r   = "0"
		err error
	)
	for idx := range a {
		r, err = Add(r, a[idx])
		if err != nil {
			return "", err
		}
	}

	return r, nil
}

func Add(a1, a2 string) (string, error) {
	if a1 == "" {
		a1 = "0"
	}
	if a2 == "" {
		a2 = "0"
	}
	d1, err := decimal(a1)
	if err != nil {
		return "", err
	}
	d2, err := decimal(a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)+qkit.StrToFloat64(a2), 'f', max(d1, d2), 64), nil
}

func AddX(a1, a2 string) string {
	s, err := Add(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Subtract(a1, a2 string) (string, error) {
	if a1 == "" {
		a1 = "0"
	}
	if a2 == "" {
		a2 = "0"
	}
	d1, err := decimal(a1)
	if err != nil {
		return "", err
	}
	d2, err := decimal(a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)-qkit.StrToFloat64(a2), 'f', max(d1, d2), 64), nil
}

func SubtractX(a1, a2 string) string {
	s, err := Subtract(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Multiply(a1, a2 string) (string, error) {
	if a1 == "" {
		a1 = "0"
	}
	if a2 == "" {
		a2 = "0"
	}
	d1, err := decimal(a1)
	if err != nil {
		return "", err
	}
	d2, err := decimal(a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)*qkit.StrToFloat64(a2), 'f', max(d1, d2), 64), nil
}

func MultiplyX(a1, a2 string) string {
	s, err := Multiply(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Divide(a1, a2 string) (string, error) {
	if a1 == "" {
		a1 = "0"
	}
	if a2 == "" {
		a2 = "0"
	}
	d1, err := decimal(a1)
	if err != nil {
		return "", err
	}
	d2, err := decimal(a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)/qkit.StrToFloat64(a2), 'f', max(d1, d2), 64), nil
}

func DivideX(a1, a2 string) string {
	s, err := Divide(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Equal(a1, a2 string) bool {
	return qkit.StrToFloat64(SubtractX(a1, a2)) == 0.0
}

func Ceil(a string) string {
	return qkit.Float64ToStr(math.Ceil(qkit.StrToFloat64(a)))
}

func Floor(a string) string {
	return qkit.Float64ToStr(math.Floor(qkit.StrToFloat64(a)))
}

func Round(a string) string {
	return qkit.Float64ToStr(math.Round(qkit.StrToFloat64(a)))
}

func max(x1, x2 int) int {
	if x1 > x2 {
		return x1
	}

	return x2
}

func decimal(a1 string) (int, error) {
	if a1 == "" {
		a1 = "0"
	}
	parts := strings.Split(a1, ".")
	switch len(parts) {
	case 1:
		return 0, nil
	case 2:
		return len(parts[1]), nil
	default:
		return 0, fmt.Errorf("invalid decimal number: %s", a1)
	}
}

func ToPrecision(a string, precision int) string {
	if a == "" {
		a = "0"
	}
	af, _ := strconv.ParseFloat(a, 64)

	return strconv.FormatFloat(af, 'f', precision, 64)
}

func FixPrecision(a string, currency string) string {
	return ToPrecision(a, Precision(currency))
}

func ToInt(a string) (int, error) {
	parts := strings.Split(a, ".")
	switch len(parts) {
	case 1:
		v, err := strconv.ParseInt(a, 10, 64)

		return int(v), err
	case 2:
		aa := fmt.Sprintf("%s%s", parts[0], parts[1])

		v, err := strconv.ParseInt(aa, 10, 64)

		return int(v), err
	default:
		panic("BUG!! invalid format")
	}
}

func ToIntX(a string) int {
	v, err := ToInt(a)
	if err != nil {
		panic(err)
	}

	return v
}

func ToUInt(a string) (uint, error) {
	parts := strings.Split(a, ".")
	switch len(parts) {
	case 1:
		v, err := strconv.ParseInt(a, 10, 64)

		return uint(v), err
	case 2:
		aa := fmt.Sprintf("%s%s", parts[0], parts[1])

		v, err := strconv.ParseInt(aa, 10, 64)

		return uint(v), err
	default:
		panic("BUG!! invalid format")
	}
}

func ToUIntX(a string) uint {
	v, err := ToUInt(a)
	if err != nil {
		panic(err)
	}

	return v
}

func ToFloat(a string) (float64, error) {
	return strconv.ParseFloat(a, 64)
}

func ToFloatX(a string) float64 {
	v, err := ToFloat(a)
	if err != nil {
		panic(err)
	}

	return v
}

// GT returns true if a > b
func GT(a, b string) bool {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	d1, _ := decimal(a)
	d2, _ := decimal(b)
	af, _ := strconv.ParseFloat(a, 64)
	bf, _ := strconv.ParseFloat(b, 64)
	d := max(d1, d2)

	return int64(af*math.Pow10(d)) > int64(bf*math.Pow10(d))
}

// GTE returns true if a >= b
func GTE(a, b string) bool {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	d1, _ := decimal(a)
	d2, _ := decimal(b)
	af, _ := strconv.ParseFloat(a, 64)
	bf, _ := strconv.ParseFloat(b, 64)
	d := max(d1, d2)

	return int64(af*math.Pow10(d)) >= int64(bf*math.Pow10(d))
}

// LT returns true if a < b
func LT(a, b string) bool {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	d1, _ := decimal(a)
	d2, _ := decimal(b)
	af, _ := strconv.ParseFloat(a, 64)
	bf, _ := strconv.ParseFloat(b, 64)
	d := max(d1, d2)

	return int64(af*math.Pow10(d)) < int64(bf*math.Pow10(d))
}

// LTE returns true if a <= b
func LTE(a, b string) bool {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	d1, _ := decimal(a)
	d2, _ := decimal(b)
	af, _ := strconv.ParseFloat(a, 64)
	bf, _ := strconv.ParseFloat(b, 64)
	d := max(d1, d2)

	return int64(af*math.Pow10(d)) <= int64(bf*math.Pow10(d))
}
