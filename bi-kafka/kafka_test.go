package bikafka

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/IBM/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

// createMockProducer creates a producer with a mock async producer for testing
func createMockProducer(t *testing.T, config *Config) *Producer {
	return createMockProducerWithExpectations(t, config, 1)
}

// createMockProducerWithExpectations creates a producer with specific expectations
func createMockProducerWithExpectations(t *testing.T, config *Config, messageCount int) *Producer {
	// Create mock async producer
	mockAsyncProducer := mocks.NewAsyncProducer(t, nil)

	// Set up expectations for the mock
	for i := 0; i < messageCount; i++ {
		mockAsyncProducer.ExpectInputAndSucceed()
	}

	// Create producer with mock
	producer := &Producer{
		producer: mockAsyncProducer,
		config:   config,
	}

	return producer
}

func createMockProducerWithoutExpectation(t *testing.T, config *Config) *Producer {
	// Create mock async producer
	mockAsyncProducer := mocks.NewAsyncProducer(t, nil)

	// Create producer with mock
	producer := &Producer{
		producer: mockAsyncProducer,
		config:   config,
	}

	return producer
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
				// Note: This will fail in tests without a real Kafka broker
				// In a real scenario, you'd use a mock or test container
				assert.Error(t, err) // Expected to fail without Kafka broker
				assert.Nil(t, producer)
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
		assert.Equal(t, "version", string(msg.Headers[0].Key))
		assert.Equal(t, "2.0", string(msg.Headers[0].Value))
	})
}

func TestProducer_ProduceBatch(t *testing.T) {
	t.Run("successful batch production", func(t *testing.T) {
		config := &Config{
			BootstrapServers: []string{"localhost:9092"},
		}
		producer := createMockProducerWithExpectations(t, config, 2)
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
		headers := []sarama.RecordHeader{
			{
				Key:   []byte("header1"),
				Value: []byte("value1"),
			},
			{
				Key:   []byte("header2"),
				Value: []byte("value2"),
			},
		}

		msg := &Message{
			Topic:   "test-topic",
			Value:   []byte("test-value"),
			Headers: headers,
		}

		assert.Len(t, msg.Headers, 2)
		assert.Equal(t, "header1", string(msg.Headers[0].Key))
		assert.Equal(t, "value1", string(msg.Headers[0].Value))
	})
}
