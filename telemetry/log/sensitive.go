package log

import (
	"encoding/json"
	"io"
	"reflect"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*Sensitive)(nil)

type sensitiveConfig struct {
	zapcore.EncoderConfig
}

type Sensitive struct {
	zapcore.Encoder
}

func newSensitive(cfg sensitiveConfig) *Sensitive {
	cfg.NewReflectedEncoder = newJSONEncoder
	return &Sensitive{
		Encoder: zapcore.NewJSONEncoder(cfg.EncoderConfig),
	}
}

func (s Sensitive) Clone() zapcore.Encoder {
	return Sensitive{Encoder: s.Encoder.Clone()}
}

func (s Sensitive) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return s.Encoder.EncodeEntry(entry, fields)
}

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

	switch rv.Kind() {
	case reflect.Struct:
		copiedRV := reflect.New(rv.Type()).Elem()
		copiedRV.Set(rv)

		return j.enc.Encode(maskStruct(copiedRV))
	default:
		return j.enc.Encode(v)
	}
}

func maskStruct(rv reflect.Value) any {
	if !rv.CanSet() {
		newRV := reflect.New(rv.Type())
		newRV.Elem().Set(rv)

		return maskStruct(newRV.Elem())
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
