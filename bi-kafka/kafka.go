package bikafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

const version = "version"

// Producer represents a Kafka producer
type Producer struct {
	writer *kafka.Writer
	config *Config
}

// Config holds the configuration for the Kafka producer
type Config struct {
	// BootstrapServers are the initial Kafka brokers for metadata discovery
	BootstrapServers []string
	// Brokers is an alias for BootstrapServers (for backward compatibility)
	Brokers []string
	// Optional configurations
	Compression kafka.Compression
	Timeout     time.Duration
	// SASL configuration
	Username string
	Password string
	// TLS configuration
	EnableTLS bool
	// Required acks
	RequiredAcks int
	// ClientID for identification
	ClientID string
	// Max message bytes
	MaxMessageBytes int
	// Async
	Async bool
	// Batch size
	BatchSize int
}

// Message represents a Kafka message
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   []kafka.Header
	Timestamp time.Time
	Version   string
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
		config.RequiredAcks = 1 // Wait for local acknowledgment
	}

	if config.BatchSize == 0 {
		config.BatchSize = 100
	}

	// Create kafka writer configuration
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		MaxAttempts:  3,
		BatchSize:    config.BatchSize,
		BatchTimeout: config.Timeout,
		Async:        config.Async,
	}

	return &Producer{
		writer: writer,
		config: config,
	}, nil
}

// Produce sends a message to Kafka
func (p *Producer) Produce(ctx context.Context, msg *Message) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	if msg.Topic == "" {
		return fmt.Errorf("topic cannot be empty")
	}
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   version,
		Value: []byte(msg.Version),
	})

	// Create kafka message
	kafkaMsg := kafka.Message{
		Topic:   msg.Topic,
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: msg.Headers,
	}

	if !msg.Timestamp.IsZero() {
		kafkaMsg.Time = msg.Timestamp
	}

	// Send message with context
	err := p.writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", msg.Topic, err)
	}

	return nil
}

// ProduceBatch sends multiple messages to Kafka
func (p *Producer) ProduceBatch(ctx context.Context, messages []*Message) error {
	if len(messages) == 0 {
		return nil
	}

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		if msg == nil {
			return fmt.Errorf("message at index %d is nil", i)
		}

		if msg.Topic == "" {
			return fmt.Errorf("topic cannot be empty for message at index %d", i)
		}
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   version,
			Value: []byte(msg.Version),
		})
		kafkaMessages[i] = kafka.Message{
			Topic:   msg.Topic,
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: msg.Headers,
		}

		if !msg.Timestamp.IsZero() {
			kafkaMessages[i].Time = msg.Timestamp
		}
	}

	// Send all messages with context
	err := p.writer.WriteMessages(ctx, kafkaMessages...)
	if err != nil {
		return fmt.Errorf("failed to send batch messages: %w", err)
	}

	return nil
}

// Close closes the producer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// GetBootstrapServers returns the bootstrap servers being used
func (p *Producer) GetBootstrapServers() []string {
	if len(p.config.BootstrapServers) > 0 {
		return p.config.BootstrapServers
	}
	return p.config.Brokers
}
