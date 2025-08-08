package bikafka

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

// createMockProducer creates a producer with a mock for testing
func createMockProducer(t *testing.T, config *Config) *MockProducer {
	return NewMockProducer(config)
}

func createMockProducerWithoutExpectation(t *testing.T, config *Config) *MockProducer {
	return NewMockProducer(config)
}
func TestNewProducer(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "empty bootstrap servers and brokers",
			config: &Config{
				BootstrapServers: []string{},
				Brokers:          []string{},
			},
			wantErr: true,
		},
		{
			name: "valid config with bootstrap servers",
			config: &Config{
				BootstrapServers: []string{"localhost:9092"},
			},
			wantErr: false,
		},
		{
			name: "valid config with brokers (backward compatibility)",
			config: &Config{
				Brokers: []string{"localhost:9092"},
			},
			wantErr: false,
		},
		{
			name: "bootstrap servers take precedence over brokers",
			config: &Config{
				BootstrapServers: []string{"localhost:9092"},
				Brokers:          []string{"localhost:9093"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			producer, err := NewProducer(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, producer)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestProducer_Produce(t *testing.T) {
	t.Run("nil message", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		producer := createMockProducerWithoutExpectation(t, config)
		defer producer.Close()

		err := producer.Produce(context.Background(), nil, "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message cannot be nil")
	})

	t.Run("empty topic", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		producer := createMockProducerWithoutExpectation(t, config)
		defer producer.Close()

		msg := &Message{
			Topic: "",
			Value: []byte("test message"),
		}

		err := producer.Produce(context.Background(), msg, "1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "topic cannot be empty")
	})

	t.Run("successful message production", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		producer := createMockProducer(t, config)
		defer producer.Close()

		msg := &Message{
			Topic: "test-topic",
			Key:   []byte("test-key"),
			Value: []byte("test message"),
		}

		err := producer.Produce(context.Background(), msg, "1")
		assert.NoError(t, err)
	})

	t.Run("message with version header", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		producer := createMockProducer(t, config)
		defer producer.Close()

		msg := &Message{
			Topic: "test-topic",
			Value: []byte("test message"),
		}

		err := producer.Produce(context.Background(), msg, "2.0")
		assert.NoError(t, err)

		// Check that version header was added
		assert.Len(t, msg.Headers, 1)
		assert.Equal(t, "version", msg.Headers[0].Key)
		assert.Equal(t, "2.0", string(msg.Headers[0].Value))
	})
}

func TestProducer_ProduceBatch(t *testing.T) {
	t.Run("successful batch production", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		producer := createMockProducer(t, config)
		defer producer.Close()

		messages := []*Message{
			{
				Topic: "test-topic",
				Key:   []byte("key1"),
				Value: []byte("test message 1"),
			},
			{
				Topic: "test-topic",
				Key:   []byte("key2"),
				Value: []byte("test message 2"),
			},
		}

		err := producer.ProduceBatch(context.Background(), messages, "1")
		assert.NoError(t, err)
	})
}

func TestProducer_GetBootstrapServers(t *testing.T) {
	t.Run("bootstrap servers take precedence", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"bootstrap1:9092", "bootstrap2:9092"},
			Brokers:          []string{"broker1:9092", "broker2:9092"},
		}
		producer := createMockProducerWithoutExpectation(t, config)
		defer producer.Close()

		servers := producer.GetBootstrapServers()
		assert.Equal(t, []string{"bootstrap1:9092", "bootstrap2:9092"}, servers)
	})

	t.Run("fallback to brokers when bootstrap servers empty", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{},
			Brokers:          []string{"broker1:9092", "broker2:9092"},
		}
		producer := createMockProducerWithoutExpectation(t, config)
		defer producer.Close()

		servers := producer.GetBootstrapServers()
		assert.Equal(t, []string{"broker1:9092", "broker2:9092"}, servers)
	})
}

func TestConfig_DefaultValues(t *testing.T) {

	t.Run("custom client ID", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
			ClientID:         "custom-client",
		}
		producer := createMockProducerWithoutExpectation(t, config)
		defer producer.Close()

		assert.Equal(t, "custom-client", producer.config.ClientID)
	})
}

func TestMessage_Validation(t *testing.T) {
	t.Run("valid message", func(t *testing.T) {
		msg := &Message{
			Topic:     "test-topic",
			Key:       []byte("test-key"),
			Value:     []byte("test-value"),
			Timestamp: time.Now(),
		}

		assert.NotEmpty(t, msg.Topic)
		assert.NotNil(t, msg.Key)
		assert.NotNil(t, msg.Value)
		assert.False(t, msg.Timestamp.IsZero())
	})

	t.Run("message with headers", func(t *testing.T) {
		headers := []kafka.Header{
			{
				Key:   "header1",
				Value: []byte("value1"),
			},
			{
				Key:   "header2",
				Value: []byte("value2"),
			},
		}

		msg := &Message{
			Topic:   "test-topic",
			Value:   []byte("test-value"),
			Headers: headers,
		}

		assert.Len(t, msg.Headers, 2)
		assert.Equal(t, "header1", msg.Headers[0].Key)
		assert.Equal(t, "value1", string(msg.Headers[0].Value))
	})
}

func TestMockProducer(t *testing.T) {
	t.Run("new mock producer", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		mock := NewMockProducer(config)
		defer mock.Close()

		assert.NotNil(t, mock)
		assert.Equal(t, 0, mock.GetMessageCount())
		assert.False(t, mock.IsClosed())
	})

	t.Run("produce single message", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		mock := NewMockProducer(config)
		defer mock.Close()

		msg := &Message{
			Topic: "test-topic",
			Key:   []byte("test-key"),
			Value: []byte("test-value"),
		}

		err := mock.Produce(context.Background(), msg, "1.0")
		assert.NoError(t, err)
		assert.Equal(t, 1, mock.GetMessageCount())

		messages := mock.GetMessages()
		assert.Len(t, messages, 1)
		assert.Equal(t, "test-topic", messages[0].Topic)
		assert.Equal(t, []byte("test-key"), messages[0].Key)
		assert.Equal(t, []byte("test-value"), messages[0].Value)
		assert.Equal(t, "1.0", string(messages[0].Headers[0].Value))
	})

	t.Run("produce batch messages", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		mock := NewMockProducer(config)
		defer mock.Close()

		messages := []*Message{
			{
				Topic: "test-topic",
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
			{
				Topic: "test-topic",
				Key:   []byte("key2"),
				Value: []byte("value2"),
			},
		}

		err := mock.ProduceBatch(context.Background(), messages, "2.0")
		assert.NoError(t, err)
		assert.Equal(t, 2, mock.GetMessageCount())

		storedMessages := mock.GetMessages()
		assert.Len(t, storedMessages, 2)
		assert.Equal(t, "2.0", string(storedMessages[0].Headers[0].Value))
		assert.Equal(t, "2.0", string(storedMessages[1].Headers[0].Value))
	})

	t.Run("clear messages", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		mock := NewMockProducer(config)
		defer mock.Close()

		msg := &Message{
			Topic: "test-topic",
			Value: []byte("test-value"),
		}

		err := mock.Produce(context.Background(), msg, "1.0")
		assert.NoError(t, err)
		assert.Equal(t, 1, mock.GetMessageCount())

		mock.ClearMessages()
		assert.Equal(t, 0, mock.GetMessageCount())
		assert.Len(t, mock.GetMessages(), 0)
	})

	t.Run("close producer", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		mock := NewMockProducer(config)

		assert.False(t, mock.IsClosed())
		err := mock.Close()
		assert.NoError(t, err)
		assert.True(t, mock.IsClosed())
	})

	t.Run("produce after close", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		mock := NewMockProducer(config)
		mock.Close()

		msg := &Message{
			Topic: "test-topic",
			Value: []byte("test-value"),
		}

		err := mock.Produce(context.Background(), msg, "1.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "writer is closed")
	})
}
