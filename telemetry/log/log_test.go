package log_test

import (
	"context"
	"io"
	"testing"

	qkit "github.com/clubpay/qlubkit-go"
	"github.com/clubpay/qlubkit-go/telemetry/log"
	"go.uber.org/zap/zapcore"
)

type sampleData struct {
	Name       string            `json:"name"`
	Age        int               `json:"age"`
	Card       string            `json:"card" sensitive:"card"`
	Phone      string            `json:"phone" sensitive:"phone"`
	PhonePtr   *string           `json:"phone_ptr" sensitive:"phone"`
	Email      string            `json:"email" sensitive:"email"`
	EmailPtr   *string           `json:"email_ptr" sensitive:"email"`
	RandomData map[string]string `json:"random_data" sensitive:"-"`
}

func TestLog(t *testing.T) {
	t.Run("Log with Sensitive Data", func(t *testing.T) {
		l := log.New(log.WithLevel(log.DebugLevel), log.WithSensitive())
		x := sampleData{
			Name:     "ehsan",
			Age:      20,
			Card:     "5022291068612222",
			PhonePtr: qkit.StringPtr("905315802262"),
			Phone:    "905315802262",
			EmailPtr: qkit.StringPtr("ehsan@qlub.io"),
			Email:    "ehsan@qlub.io",
			RandomData: map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		}
		l.DebugCtx(
			context.Background(),
			"sample",
			log.Reflect("data", x))
		if x.Card != "5022291068612222" {
			t.Error("Card should be masked in log")
		}
		if x.Phone != "905315802262" {
			t.Error("Phone should be masked in log")
		}
		l.DebugCtx(
			context.Background(),
			"sample",
			log.Reflect("data", &x))
		if x.Card != "5022291068612222" {
			t.Error("Card should be masked in log")
		}
		if x.Phone != "905315802262" {
			t.Error("Phone should be masked in log")
		}
		if qkit.String(x.PhonePtr) != "905315802262" {
			t.Error("PhonePtr should be masked in log")
		}
		if qkit.String(x.EmailPtr) != "ehsan@qlub.io" {
			t.Error("EmailPtr should be masked in log")
		}
		if x.RandomData["key"] != "value" {
			t.Error("RandomData should be masked in log")
		}
		if x.RandomData["key2"] != "value2" {
		}
	})
}

func TestOTelEvent(t *testing.T) {
	t.Run("Trace with Sensitive Data", func(t *testing.T) {
		x := sampleData{
			Name:     "ehsan",
			Age:      20,
			Card:     "5022291068612222",
			PhonePtr: qkit.StringPtr("905315802262"),
			Phone:    "905315802262",
			EmailPtr: qkit.StringPtr("ehsan@qlub.io"),
			Email:    "ehsan@qlub.io",
			RandomData: map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		}
		var xi any = x
		attrs := log.AppendField(nil, log.Reflect("data", xi))
		t.Log(attrs)
	})
}

func BenchmarkLog(b *testing.B) {
	l := log.New(
		log.WithLevel(log.DebugLevel),
		log.WithJSON(),
		log.WithWriter(zapcore.AddSync(io.Discard)),
	)
	for i := 0; i < b.N; i++ {
		l.DebugCtx(
			context.Background(),
			"sample",
			log.Reflect(
				"data", sampleData{
					Name:  "ehsan",
					Age:   20,
					Card:  "5022291068612222",
					Phone: "905315802262",
					Email: "ehsan@ronak.com",
				}),
		)
	}
}

func BenchmarkLogWithSensitive(b *testing.B) {
	l := log.New(
		log.WithLevel(log.DebugLevel),
		log.WithSensitive(),
		log.WithWriter(zapcore.AddSync(io.Discard)),
	)

	for i := 0; i < b.N; i++ {
		l.DebugCtx(
			context.Background(),
			"sample",
			log.Reflect(
				"data", sampleData{
					Name:  "ehsan",
					Age:   20,
					Card:  "5022291068612222",
					Phone: "905315802262",
					Email: "ehsan@ronak.com",
				}),
		)
	}
}

func TestLogWithLevel(t *testing.T) {
	l := log.New(log.WithLevel(log.InfoLevel))
	l.DebugCtx(
		context.Background(),
		"sample",
		log.Reflect("data", sampleData{
			Name:  "ehsan",
			Age:   20,
			Card:  "5022291068612222",
			Phone: "905315802262",
			Email: "ehsan@ronak.com",
		}))

}
