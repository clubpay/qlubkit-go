package datadog

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	qkit "github.com/clubpay/qlubkit-go"
	"github.com/clubpay/qlubkit-go/telemetry/log"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

//type writeFunc func(lvl log.Level, buf *buffer.Buffer) error

type logEntry struct {
	l   log.Level
	buf *buffer.Buffer
}

type core struct {
	cfg config
	zapcore.LevelEnabler
	client *datadog.LogsApiService
	enc    log.Encoder
	f      *qkit.FlusherPool
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
	ddClient := datadog.NewAPIClient(ddConfig).LogsApi

	c := &core{
		cfg:          cfg,
		client:       ddClient,
		LevelEnabler: cfg.lvl,
		enc: log.EncoderBuilder().
			JsonEncoder(),
	}

	c.f = qkit.NewFlusherPool(10, 100, c.flushFuncAPI)

	return c
}

func NewAgent(hostPort string, opts ...Option) log.Core {
	if hostPort == "" {
		return zapcore.NewNopCore()
	}

	cfg := config{
		flushTimeout:  time.Second * 5,
		agentHostPort: hostPort,
		tags:          map[string]string{},
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	cfg.tagsToStr()

	_ = os.Setenv("DD_SITE", cfg.site)
	_ = os.Setenv("DD_API_KEY", cfg.apiKey)

	ddConfig := datadog.NewConfiguration()
	ddClient := datadog.NewAPIClient(ddConfig).LogsApi

	c := &core{
		cfg:          cfg,
		client:       ddClient,
		LevelEnabler: cfg.lvl,
		enc: log.EncoderBuilder().
			JsonEncoder(),
	}

	c.f = qkit.NewFlusherPool(10, 100, c.flushFuncAPI)

	for k, v := range cfg.tags {
		c.enc.AddString(fmt.Sprintf("tags.%s", k), v)
	}

	return c
}

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
	buf, err := c.enc.EncodeEntry(ent, fs)
	if err != nil {
		return err
	}

	c.f.Enter("api",
		qkit.NewEntry(logEntry{
			l:   ent.Level,
			buf: buf,
		}),
	)

	return nil
}

func (c *core) flushFuncAPI(_ string, entries []qkit.FlushEntry) {
	body := make([]datadog.HTTPLogItem, len(entries))
	for idx, e := range entries {
		ent := e.Value().(logEntry)
		body[idx] = datadog.HTTPLogItem{
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
		*datadog.NewSubmitLogOptionalParameters().
			WithContentEncoding(datadog.CONTENTENCODING_DEFLATE),
	)

	for _, e := range entries {
		ent := e.Value().(logEntry)
		ent.buf.Free()
	}
}

/*
func (c *core) flushFuncAgent(_ string, entries []qkit.FlushEntry) {
	conn, err := net.Dial("tcp", c.cfg.agentHostPort)
	if err != nil {
		return
	}

	for _, e := range entries {
		ent := e.Value().(logEntry)
		_, _ = conn.Write(ent.buf.Bytes())
		ent.buf.Free()
	}

	_ = conn.Close()
}
*/

func (c *core) Sync() error {
	return nil
}
