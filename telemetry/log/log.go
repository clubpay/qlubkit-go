package log

import (
	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2021 - Sep - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
*/

var (
	DefaultLogger *logger
	NopLogger     *logger
)

type (
	Level           = zapcore.Level
	Field           = zapcore.Field
	Entry           = zapcore.Entry
	FieldType       = zapcore.FieldType
	CheckedEntry    = zapcore.CheckedEntry
	DurationEncoder = zapcore.DurationEncoder
	CallerEncoder   = zapcore.CallerEncoder
	LevelEncoder    = zapcore.LevelEncoder
	TimeEncoder     = zapcore.TimeEncoder
	Encoder         = zapcore.Encoder
	Core            = zapcore.Core
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Check(Level, string) *CheckedEntry
	SetLevel(level Level)
	With(name string) Logger
	WithCore(core Core) Logger
}

type SugaredFormatLogger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Printf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type SugaredLogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

func init() {
	DefaultLogger = New(
		WithSkipCaller(1),
	)

	NopLogger = newNOP()
}

func With(name string) Logger {
	return DefaultLogger.With(name)
}
