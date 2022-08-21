package log

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SugaredLogger struct {
	sz     *zap.SugaredLogger
	prefix string
}

func (l SugaredLogger) Debug(template string, args ...interface{}) {
	l.sz.Debugf(addPrefix(l.prefix, template), args...)
}

func (l SugaredLogger) DebugCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Debug(template, args...)
}

func (l SugaredLogger) Info(template string, args ...interface{}) {
	l.sz.Infof(addPrefix(l.prefix, template), args...)
}

func (l SugaredLogger) InfoCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Info(template, args...)
}

func (l SugaredLogger) Warn(template string, args ...interface{}) {
	l.sz.Warnf(addPrefix(l.prefix, template), args...)
}

func (l SugaredLogger) WarnCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Warn(template, args...)
}

func (l SugaredLogger) Error(template string, args ...interface{}) {
	l.sz.Errorf(addPrefix(l.prefix, template), args...)
}

func (l SugaredLogger) ErrorCtx(ctx context.Context, template string, args ...interface{}) {
	addTraceEvent(ctx, fmt.Sprintf(template, args...))
	l.Error(template, args...)
}

func (l SugaredLogger) Fatal(template string, args ...interface{}) {
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
		traceKey := attribute.Key(k)
		switch v := v.(type) {
		case string:
			attrs = append(attrs, traceKey.String(v))
		case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
			attrs = append(attrs, traceKey.String(fmt.Sprintf("%d", v)))
		case []byte:
			attrs = append(attrs, traceKey.String(string(v)))
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
