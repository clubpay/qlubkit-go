package log_test

import (
	"context"
	"io"
	"testing"

	"github.com/clubpay/qlubkit-go/telemetry/log"
	"go.uber.org/zap/zapcore"
)

type sampleData struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Card  string `json:"card" sensitive:"card"`
	Phone string `json:"phone" sensitive:"phone"`
	Email string `json:"email" sensitive:"email"`
}

const sampleRequest = `GET /payments/pay_ny2hdrfeub3utnz6jjb6alsaii HTTP/1.1
Host: api.sandbox.checkout.com
User-Agent: checkout-sdk-go/1.2.0
Accept: application/json
Authorization: Bearer sk_sbox_jhxwcevdkoxhi3qwlp7wtka5cu#
Content-Type: application/json
Accept-Encoding: gzip
`

func TestLog(t *testing.T) {
	t.Run("Log with Sensitive Data", func(t *testing.T) {
		l := log.New(log.WithLevel(log.DebugLevel), log.WithSensitive())
		x := sampleData{
			Name:  "ehsan",
			Age:   20,
			Card:  "5022291068612222",
			Phone: "905315802262",
			Email: "ehsan@ronak.com",
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

		var TestFullReqData = "GET / HTTP/1.1\r\nHost: localhost:8091\r\nConnection: keep-alive\r\nsec-ch-ua: \" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"97\", \"Chromium\";v=\"97\"\r\nsec-ch-ua-mobile: ?0\r\nsec-ch-ua-platform: \"Windows\"\r\nDNT: 1\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\r\nSec-Fetch-Site: none\r\nSec-Fetch-Mode: navigate\r\nSec-Fetch-User: ?1\r\nSec-Fetch-Dest: document\r\n\r\n"
		l.DebugCtx(
			context.Background(),
			"sample",
			log.String("str", TestFullReqData),
			log.Reflect("ref", TestFullReqData),
		)

		l.DebugCtx(
			context.Background(),
			"sample",
			log.String("str", sampleRequest),
			log.Reflect("ref", sampleRequest),
		)
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
