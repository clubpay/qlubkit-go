package datadog

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/clubpay/qlubkit-go/telemetry/log"
	"go.uber.org/zap/buffer"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"go.uber.org/zap/zapcore"
)

type writeFunc func(lvl log.Level, buf *buffer.Buffer) error

type core struct {
	cfg config
	zapcore.LevelEnabler
	client *datadog.LogsApiService
	enc    log.Encoder
	wf     writeFunc
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

	c.wf = c.writeAPI

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

	c.wf = c.writeAgent
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

	return c.wf(ent.Level, buf)
}

func (c *core) writeAPI(_ log.Level, buf *buffer.Buffer) error {
	body := []datadog.HTTPLogItem{
		{
			Ddsource: c.cfg.source,
			Ddtags:   c.cfg.tagsStr,
			Hostname: c.cfg.hostname,
			Message:  buf.String(),
			Service:  c.cfg.service,
		},
	}

	ctx, cf := context.WithTimeout(context.Background(), c.cfg.flushTimeout)
	defer cf()
	_, _, err := c.client.SubmitLog(
		datadog.NewDefaultContext(ctx),
		body,
		*datadog.NewSubmitLogOptionalParameters().
			WithContentEncoding(datadog.CONTENTENCODING_DEFLATE),
	)

	buf.Free()

	return err
}

func (c *core) writeAgent(lvl log.Level, buf *buffer.Buffer) error {
	defer buf.Free()

	conn, err := net.Dial("tcp", c.cfg.agentHostPort)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf.Bytes())
	_ = conn.Close()

	return err
}

func (c *core) Sync() error {
	return nil
}
