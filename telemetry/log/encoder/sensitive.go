package encoder

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*Sensitive)(nil)

type SensitiveConfig struct {
	zapcore.EncoderConfig
}

type Sensitive struct {
	zapcore.Encoder
}

func NewSensitive(cfg SensitiveConfig) *Sensitive {
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
