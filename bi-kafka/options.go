package bikafka

import (
	"time"

	"github.com/IBM/sarama"
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
func WithCompression(compression sarama.CompressionCodec) Option {
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

// WithRequiredAcks sets the required acknowledgments
func WithRequiredAcks(acks sarama.RequiredAcks) Option {
	return func(c *Config) {
		c.RequiredAcks = acks
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
		c.Compression = sarama.CompressionSnappy
	}
}

func WithGzipCompression() Option {
	return func(c *Config) {
		c.Compression = sarama.CompressionGZIP
	}
}

func WithLz4Compression() Option {
	return func(c *Config) {
		c.Compression = sarama.CompressionLZ4
	}
}

func WithZstdCompression() Option {
	return func(c *Config) {
		c.Compression = sarama.CompressionZSTD
	}
}

// Required acks options
func WithWaitForAll() Option {
	return func(c *Config) {
		c.RequiredAcks = sarama.WaitForAll
	}
}

func WithWaitForLocal() Option {
	return func(c *Config) {
		c.RequiredAcks = sarama.WaitForLocal
	}
}

func WithNoResponse() Option {
	return func(c *Config) {
		c.RequiredAcks = sarama.NoResponse
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
