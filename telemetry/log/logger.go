package log

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2019 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
*/

// zapLogger is a wrapper around zap.Logger and adds a good few features to it.
// It provides layered logs which could be used by separate packages, and could be turned off or on
// separately. Separate layers could also have independent log levels.
// Whenever you change log level it propagates through its children.
type zapLogger struct {
	prefix     string
	skipCaller int
	//encoder    zapcore.Encoder
	z *zap.Logger
	//sz         *zap.SugaredLogger
	lvl zap.AtomicLevel
}

func New(opts ...Option) *zapLogger {
	encodeBuilder := EncoderBuilder().
		WithTimeKey("ts").
		WithLevelKey("level").
		WithNameKey("name").
		WithCallerKey("caller").
		WithMessageKey("msg")

	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	l := &zapLogger{
		lvl:        zap.NewAtomicLevelAt(cfg.level),
		skipCaller: cfg.skipCaller,
	}

	var cores []Core
	switch cfg.encoder {
	case "json":
		cores = append(cores,
			zapcore.NewCore(encodeBuilder.JsonEncoder(), zapcore.Lock(os.Stdout), l.lvl),
		)
	case "console":
		cores = append(cores,
			zapcore.NewCore(encodeBuilder.ConsoleEncoder(), zapcore.Lock(os.Stdout), l.lvl),
		)
	}

	cores = append(cores, cfg.cores...)
	l.z = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(cfg.skipCaller),
	)

	return l
}

func newNOP() *zapLogger {
	l := &zapLogger{}
	l.z = zap.NewNop()

	return l
}

func (l *zapLogger) Sugared() *sugaredLogger {
	return &sugaredLogger{
		sz:     l.z.Sugar(),
		prefix: l.prefix,
	}
}

func (l *zapLogger) Sync() error {
	return l.z.Sync()
}

func (l *zapLogger) SetLevel(lvl Level) {
	l.lvl.SetLevel(lvl)
}

func (l *zapLogger) With(name string) Logger {
	return l.WithSkip(name, l.skipCaller)
}

func (l *zapLogger) WithSkip(name string, skipCaller int) Logger {
	return l.with(l.z.Core(), name, skipCaller)
}

func (l *zapLogger) WithCore(core Core) Logger {
	return l.with(
		zapcore.NewTee(
			l.z.Core(), core,
		),
		"",
		l.skipCaller,
	)
}

func (l *zapLogger) with(core zapcore.Core, name string, skip int) Logger {
	prefix := l.prefix
	if name != "" {
		prefix = fmt.Sprintf("%s[%s]", l.prefix, name)
	}
	childLogger := &zapLogger{
		prefix:     prefix,
		skipCaller: l.skipCaller,
		z: zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skip),
		),
		lvl: l.lvl,
	}

	return childLogger
}

func (l *zapLogger) WarnOnErr(guideTxt string, err error, fields ...Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Warn(guideTxt, fields...)
	}
}

func (l *zapLogger) ErrorOnErr(guideTxt string, err error, fields ...Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error(guideTxt, fields...)
	}
}

func (l *zapLogger) checkLevel(lvl Level) bool {
	if l == nil {
		return false
	}

	// Check the level first to reduce the cost of disabled log calls.
	// Since Panic and higher may exit, we skip the optimization for those levels.
	if lvl < zapcore.DPanicLevel && !l.z.Core().Enabled(lvl) {
		return false
	}

	return true
}

func (l *zapLogger) Check(lvl Level, msg string) *CheckedEntry {
	if !l.checkLevel(lvl) {
		return nil
	}

	return l.z.Check(lvl, addPrefix(l.prefix, msg))
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(DebugLevel) {
		return
	}
	if ce := l.z.Check(DebugLevel, addPrefix(l.prefix, msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *zapLogger) DebugCtx(ctx context.Context, msg string, fields ...Field) {
	addTraceEvent(ctx, msg, fields...)
	l.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(InfoLevel) {
		return
	}
	if ce := l.z.Check(InfoLevel, addPrefix(l.prefix, msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...Field) {
	addTraceEvent(ctx, msg, fields...)
	l.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(WarnLevel) {
		return
	}
	if ce := l.z.Check(WarnLevel, addPrefix(l.prefix, msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *zapLogger) WarnCtx(ctx context.Context, msg string, fields ...Field) {
	addTraceEvent(ctx, msg, fields...)
	l.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	if l == nil {
		return
	}
	if !l.checkLevel(ErrorLevel) {
		return
	}
	if ce := l.z.Check(ErrorLevel, addPrefix(l.prefix, msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *zapLogger) ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	addTraceEvent(ctx, msg, fields...)
	l.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	if l == nil {
		return
	}
	l.z.Fatal(addPrefix(l.prefix, msg), fields...)
}

func (l *zapLogger) FatalCtx(ctx context.Context, msg string, fields ...Field) {
	addTraceEvent(ctx, msg, fields...)
	l.Fatal(msg, fields...)
}

func (l *zapLogger) RecoverPanic(funcName string, extraInfo interface{}, compensationFunc func()) {
	if r := recover(); r != nil {
		l.Error("Panic Recovered",
			zap.String("Task", funcName),
			zap.Any("Info", extraInfo),
			zap.Any("Recover", r),
			zap.ByteString("StackTrace", debug.Stack()),
		)
		if compensationFunc != nil {
			go compensationFunc()
		}
	}
}

type sugaredLogger struct {
	sz     *zap.SugaredLogger
	prefix string
}

var (
	_ SugaredLogger        = (*sugaredLogger)(nil)
	_ SugaredContextLogger = (*sugaredLogger)(nil)
)

func (l sugaredLogger) Debug(template string, args ...interface{}) {
	l.sz.Debugf(addPrefix(l.prefix, template), args...)
}

func (l sugaredLogger) DebugCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Debug(template, args...)
}

func (l sugaredLogger) Info(template string, args ...interface{}) {
	l.sz.Infof(addPrefix(l.prefix, template), args...)
}

func (l sugaredLogger) InfoCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Info(template, args...)
}

func (l sugaredLogger) Warn(template string, args ...interface{}) {
	l.sz.Warnf(addPrefix(l.prefix, template), args...)
}

func (l sugaredLogger) WarnCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Warn(template, args...)
}

func (l sugaredLogger) Error(template string, args ...interface{}) {
	l.sz.Errorf(addPrefix(l.prefix, template), args...)
}

func (l sugaredLogger) ErrorCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Error(template, args...)
}

func (l sugaredLogger) Fatal(template string, args ...interface{}) {
	l.sz.Fatalf(addPrefix(l.prefix, template), args...)
}

func addPrefix(prefix, in string) (out string) {
	if prefix != "" {
		sb := &strings.Builder{}
		sb.WriteString(prefix)
		sb.WriteRune(' ')
		sb.WriteString(in)
		out = sb.String()

		return out
	}

	return in
}

func toTraceAttrs(fields ...Field) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(fields))

	e := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(e)
	}
	for k, v := range e.Fields {
		kk := attribute.Key(k)
		switch v := v.(type) {
		case string:
			attrs = append(attrs, kk.String(v))
		case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
			attrs = append(attrs, kk.String(fmt.Sprintf("%d", v)))
		case []byte:
			attrs = append(attrs, kk.String(string(v)))
		default:
			continue
		}
	}

	return attrs
}

func addTraceEvent(ctx context.Context, msg string, fields ...Field) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(
		msg,
		trace.WithAttributes(toTraceAttrs(fields...)...),
	)
}
