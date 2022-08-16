package log

import (
	"context"

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
	DefaultLogger *zapLogger
	NopLogger     *zapLogger
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

type ContextLogger interface {
	DebugCtx(ctx context.Context, msg string, fields ...Field)
	InfoCtx(ctx context.Context, msg string, fields ...Field)
	WarnCtx(ctx context.Context, msg string, fields ...Field)
	ErrorCtx(ctx context.Context, msg string, fields ...Field)
	FatalCtx(ctx context.Context, msg string, fields ...Field)
}

type SugaredLogger interface {
	Debug(template string, args ...interface{})
	Info(template string, args ...interface{})
	Warn(template string, args ...interface{})
	Error(template string, args ...interface{})
	Fatal(template string, args ...interface{})
}

type SugaredContextLogger interface {
	DebugCtx(ctx context.Context, template string, args ...interface{})
	InfoCtx(ctx context.Context, template string, args ...interface{})
	WarnCtx(ctx context.Context, template string, args ...interface{})
	ErrorCtx(ctx context.Context, template string, args ...interface{})
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
