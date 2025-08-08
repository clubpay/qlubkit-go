package bikafka

import (
	"context"
	"sync"

	"github.com/segmentio/kafka-go"
)

// MockWriter is a mock implementation of kafka.Writer for testing
type MockWriter struct {
	messages []kafka.Message
	mu       sync.RWMutex
	closed   bool
	errors   []error
}

// NewMockWriter creates a new mock writer
func NewMockWriter() *MockWriter {
	return &MockWriter{
		messages: make([]kafka.Message, 0),
		errors:   make([]error, 0),
	}
}

// WriteMessages implements the kafka.Writer interface
func (m *MockWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrWriterClosed
	}

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Add messages to the mock
	m.messages = append(m.messages, msgs...)
	return nil
}

// Close implements the kafka.Writer interface
func (m *MockWriter) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// GetMessages returns all messages written to the mock
func (m *MockWriter) GetMessages() []kafka.Message {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages := make([]kafka.Message, len(m.messages))
	copy(messages, m.messages)
	return messages
}

// GetMessageCount returns the number of messages written
func (m *MockWriter) GetMessageCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.messages)
}

// ClearMessages clears all stored messages
func (m *MockWriter) ClearMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = make([]kafka.Message, 0)
}

// IsClosed returns true if the writer is closed
func (m *MockWriter) IsClosed() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.closed
}

// MockProducer is a mock implementation of Producer for testing
type MockProducer struct {
	writer *MockWriter
	config *Config
}

// NewMockProducer creates a new mock producer
func NewMockProducer(config *Config) *MockProducer {
	return &MockProducer{
		writer: NewMockWriter(),
		config: config,
	}
}

// Produce implements the Producer interface
func (m *MockProducer) Produce(ctx context.Context, msg *Message, version string) error {
	if msg == nil {
		return ErrMessageNil
	}

	if msg.Topic == "" {
		return ErrTopicEmpty
	}

	// Add version header
	msg.Headers = append(msg.Headers, kafka.Header{
		Key:   "version",
		Value: []byte(version),
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

	// Write message using mock writer
	return m.writer.WriteMessages(ctx, kafkaMsg)
}

// ProduceBatch implements the Producer interface
func (m *MockProducer) ProduceBatch(ctx context.Context, messages []*Message, version string) error {
	if len(messages) == 0 {
		return nil
	}

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		if msg == nil {
			return ErrMessageNil
		}

		if msg.Topic == "" {
			return ErrTopicEmpty
		}

		// Add version header
		msg.Headers = append(msg.Headers, kafka.Header{
			Key:   "version",
			Value: []byte(version),
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

	// Write messages using mock writer
	return m.writer.WriteMessages(ctx, kafkaMessages...)
}

// Close implements the Producer interface
func (m *MockProducer) Close() error {
	return m.writer.Close()
}

// GetBootstrapServers returns the bootstrap servers
func (m *MockProducer) GetBootstrapServers() []string {
	if len(m.config.BootstrapServers) > 0 {
		return m.config.BootstrapServers
	}
	return m.config.Brokers
}

// GetMessages returns all messages written to the mock
func (m *MockProducer) GetMessages() []kafka.Message {
	return m.writer.GetMessages()
}

// GetMessageCount returns the number of messages written
func (m *MockProducer) GetMessageCount() int {
	return m.writer.GetMessageCount()
}

// ClearMessages clears all stored messages
func (m *MockProducer) ClearMessages() {
	m.writer.ClearMessages()
}

// IsClosed returns true if the producer is closed
func (m *MockProducer) IsClosed() bool {
	return m.writer.IsClosed()
}

// Error definitions
var (
	ErrWriterClosed = &MockError{"writer is closed"}
	ErrMessageNil   = &MockError{"message cannot be nil"}
	ErrTopicEmpty   = &MockError{"topic cannot be empty"}
)

// MockError represents a mock error
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}
