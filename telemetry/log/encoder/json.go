package encoder

import (
	"encoding/json"
	"io"
	"reflect"

	"go.uber.org/zap/zapcore"
)

var _ zapcore.ReflectedEncoder = (*jsonEncoder)(nil)

type jsonEncoder struct {
	enc *json.Encoder
}

func newJSONEncoder(
	w io.Writer,
) zapcore.ReflectedEncoder {
	return &jsonEncoder{enc: json.NewEncoder(w)}
}

func (j jsonEncoder) Encode(v any) error {
	rv := reflect.Indirect(reflect.ValueOf(v))

	return j.enc.Encode(mask(rv))
}

func mask(rv reflect.Value) any {
	if rv.Kind() != reflect.Struct {
		return rv.Interface()
	}

	if !rv.CanSet() {
		newRV := reflect.New(rv.Type())
		newRV.Elem().Set(rv)

		return mask(newRV.Elem())
	}

	rvt := rv.Type()
	for i := range rvt.NumField() {
		if rvt.Field(i).Type.Kind() != reflect.String {
			continue
		}

		f := rv.Field(i)
		if !f.CanAddr() {
			return rv.Interface()
		}

		switch rvt.Field(i).Tag.Get("sensitive") {
		default:
			if f.Len() < 4 {
				f.SetString("****")
			} else if f.Len() < 8 {
				s := f.String()
				f.SetString(s[:2] + "****")
			} else {
				s := f.String()
				f.SetString(s[:2] + "****" + s[len(s)-2:])
			}
		case "":
		}
	}

	return rv.Interface()
}
