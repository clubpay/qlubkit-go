package qacc

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	qkit "github.com/clubpay/qlubkit-go"
)

const (
	Zero             = "0"
	defaultPrecision = 2
)

var currency = map[string]int{
	"BIF": 0,
	"CLF": 0,
	"DJF": 0,
	"GNF": 0,
	"ISK": 0,
	"IRR": 0,
	"JPY": 0,
	"KMF": 0,
	"KRW": 0,
	"PYG": 0,
	"RWF": 0,
	"UGX": 0,
	"VUV": 0,
	"VND": 0,
	"XAF": 0,
	"XOF": 0,
	"XPF": 0,
	"BHD": 3,
	"IQD": 3,
	"JOD": 3,
	"KWD": 3,
	"LYD": 3,
	"OMR": 3,
	"TND": 3,
}

func getCurrencyPrecision(curr string) int {
	p, ok := currency[curr]
	if !ok {
		return defaultPrecision
	}

	return p
}

func Precision(curr string) int {
	return getCurrencyPrecision(strings.ToUpper(curr))
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
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)+qkit.StrToFloat64(a2), 'f', d, 64), nil
}

func AddX(a1, a2 string) string {
	s, err := Add(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Subtract(a1, a2 string) (string, error) {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)-qkit.StrToFloat64(a2), 'f', d, 64), nil
}

func SubtractX(a1, a2 string) string {
	s, err := Subtract(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Multiply(a1, a2 string) (string, error) {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)*qkit.StrToFloat64(a2), 'f', d, 64), nil
}

func MultiplyX(a1, a2 string) string {
	s, err := Multiply(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Divide(a1, a2 string) (string, error) {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(qkit.StrToFloat64(a1)/qkit.StrToFloat64(a2), 'f', d, 64), nil
}

func DivideX(a1, a2 string) string {
	s, err := Divide(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Quotient(a1, a2 string) (string, error) {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return "", err
	}

	i1, i2 := ToIntX(ToPrecision(a1, d)), ToIntX(ToPrecision(a2, d))

	return qkit.IntToStr(i1 / i2), nil
}

func QuotientX(a1, a2 string) string {
	s, err := Quotient(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Remainder(a1, a2 string) (string, error) {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return "", err
	}

	i1, i2 := ToIntX(ToPrecision(a1, d)), ToIntX(ToPrecision(a2, d))

	return insignify(i1%i2, d), nil
}

func RemainderX(a1, a2 string) string {
	s, err := Remainder(a1, a2)
	if err != nil {
		panic(err)
	}

	return s
}

func Equal(a1, a2 string) bool {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return false
	}

	return (ToIntX(ToPrecision(a1, d)) - ToIntX(ToPrecision(a2, d))) == 0
}

func EqualIgnore(a1, a2, ignore string) bool {
	if GT(a1, a2) {
		return LTE(SubtractX(a1, a2), ignore)
	}
	return LTE(SubtractX(a2, a1), ignore)
}

func Abs(a string) string {
	if strings.HasPrefix(a, "-") {
		return a[1:]
	}
	return a
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

func sanitize(a string) string {
	a = strings.TrimSpace(a)
	if len(a) == 0 || a == "NaN" {
		return Zero
	}

	return a
}

func sanitizeAmount(a string) string {
	return strings.ReplaceAll(sanitize(a), ",", "")
}

func max(x1, x2 int) int {
	if x1 > x2 {
		return x1
	}

	return x2
}

func decimal(a string) (int, error) {
	parts := strings.Split(sanitize(a), ".")
	switch len(parts) {
	case 1:
		return 0, nil
	case 2:
		return len(parts[1]), nil
	default:
		return 0, fmt.Errorf("invalid decimal number: %s", a)
	}
}

func maxDecimal(a1, a2 string) (int, error) {
	d1, err := decimal(a1)
	if err != nil {
		return 0, err
	}

	d2, err := decimal(a2)
	if err != nil {
		return 0, err
	}

	return max(d1, d2), nil
}

func insignify(a, digits int) string {
	if digits == 0 {
		return qkit.IntToStr(a)
	}

	p := int(math.Pow10(digits))
	i, m := a/p, a%p
	leadingZs := strings.Repeat("0", length2(p)-1-length2(m))

	return fmt.Sprintf("%d.%s%d", i, leadingZs, m)
}

func length(a int) int {
	if a == 0 {
		return 1
	}

	return int(math.Floor(math.Log10(float64(a))) + 1)
}

func length2(a int) int {
	if a%10 == a {
		return 1
	}

	return length2(a/10) + 1
}

func zeroPrefix(a string, n int) string {
	if x := n - len(a); x > 0 {
		return fmt.Sprintf("%s%s", strings.Repeat("0", x), a)
	}

	return a
}

func ToPrecision(a string, precision int) string {
	af, _ := strconv.ParseFloat(sanitize(a), 64)

	return strconv.FormatFloat(af, 'f', precision, 64)
}

func FixPrecision(a string, currency string) string {
	return ToPrecision(a, Precision(currency))
}

// ToInt converts an amount to an integer by multiplying to 10 to power of
// floating points. "12.43" --> 1243
func ToInt(a string) (int, error) {
	a = sanitize(a)
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

func FromInt(a int, precision int) string {
	x := qkit.IntToStr(a)
	dotIndex := len(x) - precision

	// if precision is negative, then we just return the a untouched.
	if precision <= 0 {
		return x
	}

	// if precision is larger than the number of digits we use the number of the digits
	// as our new precision
	if dotIndex < 0 {
		dotIndex = 0
	}

	// this hacks to fix the problem of showing .002 -> 0.002
	if dotIndex == 0 {
		return fmt.Sprintf("0.%s", zeroPrefix(x[dotIndex:], precision))
	} else {
		return fmt.Sprintf("%s.%s", x[:dotIndex], zeroPrefix(x[dotIndex:], precision))
	}
}

func FromUInt(a uint, precision int) string {
	x := qkit.UIntToStr(a)
	dotIndex := len(x) - precision

	// if precision is negative, then we just return the a untouched.
	if precision <= 0 {
		return x
	}

	// if precision is larger than the number of digits we use the number of the digits
	// as our new precision
	if dotIndex < 0 {
		dotIndex = 0
	}

	// this hacks to fix the problem of showing .002 -> 0.002
	if dotIndex == 0 {
		return fmt.Sprintf("0.%s", zeroPrefix(x[dotIndex:], precision))
	} else {
		return fmt.Sprintf("%s.%s", x[:dotIndex], zeroPrefix(x[dotIndex:], precision))
	}
}

// ToUInt converts an amount to an integer by multiplying to 10 to power of
// floating points. "12.43" --> 1243
func ToUInt(a string) (uint, error) {
	a = sanitize(a)
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
	return strconv.ParseFloat(sanitize(a), 64)
}

func ToFloatX(a string) float64 {
	v, err := ToFloat(a)
	if err != nil {
		panic(err)
	}

	return v
}

// GT returns true if a > b
func GT(a1, a2 string) bool {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return false
	}

	a1f, _ := strconv.ParseFloat(a1, 64)
	a2f, _ := strconv.ParseFloat(a2, 64)

	return int64(a1f*math.Pow10(d)) > int64(a2f*math.Pow10(d))
}

// GTE returns true if a >= b
func GTE(a1, a2 string) bool {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return false
	}

	a1f, _ := strconv.ParseFloat(a1, 64)
	a2f, _ := strconv.ParseFloat(a2, 64)

	return int64(a1f*math.Pow10(d)) >= int64(a2f*math.Pow10(d))
}

// LT returns true if a < b
func LT(a1, a2 string) bool {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return false
	}

	a1f, _ := strconv.ParseFloat(a1, 64)
	a2f, _ := strconv.ParseFloat(a2, 64)

	return int64(a1f*math.Pow10(d)) < int64(a2f*math.Pow10(d))
}

// LTE returns true if a <= b
func LTE(a1, a2 string) bool {
	a1, a2 = sanitize(a1), sanitize(a2)
	d, err := maxDecimal(a1, a2)
	if err != nil {
		return false
	}

	a1f, _ := strconv.ParseFloat(a1, 64)
	a2f, _ := strconv.ParseFloat(a2, 64)

	return int64(a1f*math.Pow10(d)) <= int64(a2f*math.Pow10(d))
}

func EQ(a1, a2 string) bool {
	return Equal(a1, a2)
}

func NEQ(a1, a2 string) bool {
	return !EQ(a1, a2)
}
