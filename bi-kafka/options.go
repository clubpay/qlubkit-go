package bikafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// Option is a function that configures a Producer
type Option func(*Config)

// WithBootstrapServers sets the Kafka bootstrap server addresses
func WithBootstrapServers(bootstrapServers ...string) Option {
	return func(c *Config) {
		c.BootstrapServers = bootstrapServers
	}
}

// WithBrokers sets the Kafka broker addresses (alias for WithBootstrapServers)
func WithBrokers(brokers ...string) Option {
	return func(c *Config) {
		c.Brokers = brokers
	}
}

// WithClientID sets the client ID for identification
func WithClientID(clientID string) Option {
	return func(c *Config) {
		c.ClientID = clientID
	}
}

// WithCompression sets the compression algorithm
func WithCompression(compression kafka.Compression) Option {
	return func(c *Config) {
		c.Compression = compression
	}
}

// WithTimeout sets the timeout for operations
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithSASL sets SASL authentication credentials
func WithSASL(username, password string) Option {
	return func(c *Config) {
		c.Username = username
		c.Password = password
	}
}

// WithTLS enables TLS encryption
func WithTLS() Option {
	return func(c *Config) {
		c.EnableTLS = true
	}
}

// WithMaxMessageBytes sets the maximum message size
func WithMaxMessageBytes(maxBytes int) Option {
	return func(c *Config) {
		c.MaxMessageBytes = maxBytes
	}
}

// Compression options
func WithSnappyCompression() Option {
	return func(c *Config) {
		c.Compression = kafka.Snappy
	}
}

func WithGzipCompression() Option {
	return func(c *Config) {
		c.Compression = kafka.Gzip
	}
}

func WithLz4Compression() Option {
	return func(c *Config) {
		c.Compression = kafka.Lz4
	}
}

func WithZstdCompression() Option {
	return func(c *Config) {
		c.Compression = kafka.Zstd
	}
}

// Required acks options
func WithWaitForAll() Option {
	return func(c *Config) {
		c.RequiredAcks = -1 // Wait for all in-sync replicas
	}
}

func WithWaitForLocal() Option {
	return func(c *Config) {
		c.RequiredAcks = 1 // Wait for local acknowledgment
	}
}

func WithNoResponse() Option {
	return func(c *Config) {
		c.RequiredAcks = 0 // No acknowledgment required
	}
}

// NewProducerWithOptions creates a new producer with the given options
func NewProducerWithOptions(options ...Option) (*Producer, error) {
	config := &Config{}

	for _, option := range options {
		option(config)
	}

	return NewProducer(config)
}
