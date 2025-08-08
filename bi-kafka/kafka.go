package bikafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.AsyncProducer
	config   *Config
}

// Config holds the configuration for the Kafka producer
type Config struct {
	// BootstrapServers are the initial Kafka brokers for metadata discovery
	BootstrapServers []string
	// Brokers is an alias for BootstrapServers (for backward compatibility)
	Brokers []string
	// Optional configurations
	Compression sarama.CompressionCodec
	Timeout     time.Duration
	// SASL configuration
	Username string
	Password string
	// TLS configuration
	EnableTLS bool
	// ClientID for identification
	ClientID string
	// Required acks configuration
	RequiredAcks sarama.RequiredAcks
	// Max message bytes
	MaxMessageBytes int
}

// Message represents a Kafka message
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   []sarama.RecordHeader
	Timestamp time.Time
}

// NewProducer creates a new Kafka producer with the given configuration
func NewProducer(config *Config) (*Producer, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Use BootstrapServers if provided, otherwise fall back to Brokers
	brokers := config.BootstrapServers
	if len(brokers) == 0 {
		brokers = config.Brokers
	}

	if len(brokers) == 0 {
		return nil, fmt.Errorf("at least one bootstrap server must be specified")
	}

	// Set default values
	if config.Timeout == 0 {
		config.Timeout = 3 * time.Second
	}

	if config.ClientID == "" {
		config.ClientID = "bi-kafka-producer"
	}

	if config.RequiredAcks == 0 {
		config.RequiredAcks = sarama.WaitForLocal
	}

	if config.MaxMessageBytes == 0 {
		config.MaxMessageBytes = 1000000 // 1MB default
	}

	// Create Sarama configuration
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = config.RequiredAcks
	saramaConfig.Producer.Timeout = config.Timeout
	saramaConfig.Producer.Compression = config.Compression
	saramaConfig.Producer.MaxMessageBytes = config.MaxMessageBytes
	saramaConfig.ClientID = config.ClientID

	// Configure SASL if credentials are provided
	if config.Username != "" && config.Password != "" {
		saramaConfig.Net.SASL.Enable = true
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		saramaConfig.Net.SASL.User = config.Username
		saramaConfig.Net.SASL.Password = config.Password
	}

	// Configure TLS if enabled
	if config.EnableTLS {
		saramaConfig.Net.TLS.Enable = true
	}

	// Create producer
	producer, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		config:   config,
	}, nil
}

// Produce sends a message to Kafka
func (p *Producer) Produce(ctx context.Context, msg *Message, version string) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	if msg.Topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}
	msg.Headers = append(msg.Headers, sarama.RecordHeader{
		Key:   []byte("version"),
		Value: []byte(version),
	})
	// Create Sarama message
	saramaMsg := &sarama.ProducerMessage{
		Topic:   msg.Topic,
		Key:     sarama.ByteEncoder(msg.Key),
		Value:   sarama.ByteEncoder(msg.Value),
		Headers: msg.Headers,
	}

	if !msg.Timestamp.IsZero() {
		saramaMsg.Timestamp = msg.Timestamp
	}

	// Send message with context
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.producer.Input() <- saramaMsg:
		// Message sent to async producer
		return nil
	}
}

// ProduceBatch sends multiple messages to Kafka
func (p *Producer) ProduceBatch(ctx context.Context, messages []*Message, version string) error {
	if len(messages) == 0 {
		return nil
	}

	saramaMessages := make([]*sarama.ProducerMessage, len(messages))
	for i, msg := range messages {
		if msg == nil {
			return fmt.Errorf("message at index %d is nil", i)
		}

		if msg.Topic == "" {
			return fmt.Errorf("topic cannot be empty for message at index %d", i)
		}
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte("version"),
			Value: []byte(version),
		})
		saramaMessages[i] = &sarama.ProducerMessage{
			Topic:   msg.Topic,
			Key:     sarama.ByteEncoder(msg.Key),
			Value:   sarama.ByteEncoder(msg.Value),
			Headers: msg.Headers,
		}

		if !msg.Timestamp.IsZero() {
			saramaMessages[i].Timestamp = msg.Timestamp
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Send all messages to async producer
		for _, saramaMsg := range saramaMessages {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case p.producer.Input() <- saramaMsg:
				// Message sent
			}
		}
		return nil
	}
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
}

// IsConnected checks if the producer is connected to Kafka
func (p *Producer) IsConnected() bool {
	return p.producer != nil
}

// GetBootstrapServers returns the bootstrap servers being used
func (p *Producer) GetBootstrapServers() []string {
	if len(p.config.BootstrapServers) > 0 {
		return p.config.BootstrapServers
	}
	return p.config.Brokers
}
