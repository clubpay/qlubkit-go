package kafka

import (
	"crypto/tls"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// dialerParams holds connection settings shared by consumer (and future producer wiring).
type dialerParams struct {
	ClientID  string
	Timeout   time.Duration
	Username  string
	Password  string
	EnableTLS bool
}

func newDialer(p dialerParams) *kafkago.Dialer {
	d := &kafkago.Dialer{
		ClientID: p.ClientID,
		Timeout:  p.Timeout,
	}
	if p.EnableTLS {
		d.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	if p.Username != "" || p.Password != "" {
		d.SASLMechanism = plain.Mechanism{
			Username: p.Username,
			Password: p.Password,
		}
	}
	return d
}
