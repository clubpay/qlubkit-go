package datadog

import (
	"context"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	qkit "github.com/clubpay/qlubkit-go"
	"github.com/clubpay/qlubkit-go/telemetry/log"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type logEntry struct {
	buf *buffer.Buffer
}

type core struct {
	zapcore.LevelEnabler
	cfg    config
	client *datadogV2.LogsApi
	enc    log.Encoder
	f      *qkit.FlusherPool

	submitLogOpt datadogV2.SubmitLogOptionalParameters
}

func NewAPI(apiKey string, opts ...Option) log.Core {
	if apiKey == "" {
		return zapcore.NewNopCore()
	}

	cfg := config{
		flushTimeout: time.Second * 5,
		site:         "datadoghq.eu",
		apiKey:       apiKey,
		tags:         map[string]string{},
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	cfg.tagsToStr()

	_ = os.Setenv("DD_SITE", cfg.site)
	_ = os.Setenv("DD_API_KEY", cfg.apiKey)

	ddConfig := datadog.NewConfiguration()
	ddClient := datadog.NewAPIClient(ddConfig)

	c := &core{
		cfg:          cfg,
		client:       datadogV2.NewLogsApi(ddClient),
		LevelEnabler: cfg.lvl,
		enc: log.NewEncoderBuilder().
			WithMessageKey("msg").
			JsonEncoder(),
		submitLogOpt: *datadogV2.NewSubmitLogOptionalParameters().
			WithContentEncoding(datadogV2.CONTENTENCODING_DEFLATE),
	}

	c.f = qkit.NewFlusherPool(10, 100, c.flushFuncAPI)

	return c
}

var _ log.Core = (*core)(nil)

func (c *core) With(fs []log.Field) log.Core {
	return &core{
		cfg:          c.cfg,
		LevelEnabler: c.LevelEnabler,
		client:       c.client,
	}
}

func (c *core) Check(ent log.Entry, ce *log.CheckedEntry) *log.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c *core) Write(ent log.Entry, fs []log.Field) error {
	fs = append(fs, log.String("level", toDatadogLevel(ent.Level)))
	buf, err := c.enc.EncodeEntry(ent, fs)
	if err != nil {
		return err
	}

	c.f.Enter("api",
		qkit.NewEntry(logEntry{
			buf: buf,
		}),
	)

	return nil
}

func (c *core) flushFuncAPI(_ string, entries []qkit.FlushEntry) {
	body := make([]datadogV2.HTTPLogItem, len(entries))
	for idx, e := range entries {
		ent := e.Value().(logEntry)
		body[idx] = datadogV2.HTTPLogItem{
			Ddsource: c.cfg.source,
			Ddtags:   c.cfg.tagsStr,
			Hostname: c.cfg.hostname,
			Message:  ent.buf.String(),
			Service:  c.cfg.service,
		}
	}

	ctx, cf := context.WithTimeout(context.Background(), c.cfg.flushTimeout)
	defer cf()
	_, _, _ = c.client.SubmitLog(
		datadog.NewDefaultContext(ctx),
		body,
		c.submitLogOpt,
	)

	for _, e := range entries {
		ent := e.Value().(logEntry)
		ent.buf.Free()
	}
}

func (c *core) Sync() error {
	return nil
}

func toDatadogLevel(l log.Level) string {
	switch l {
	case log.DebugLevel:
		return "debug"
	case log.InfoLevel:
		return "info"
	case log.WarnLevel:
		return "warning"
	case log.ErrorLevel:
		return "err"
	case log.FatalLevel, log.PanicLevel:
		return "emergency"
	}

	return "info"
}
