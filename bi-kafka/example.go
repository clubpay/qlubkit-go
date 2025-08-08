package bikafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// Example demonstrates how to use the bi-kafka producer with traditional config
func Example() {
	// Create configuration
	config := &Config{
		BootstrapServers: []string{"localhost:9092"},
		// Optional: Set compression
		Compression: sarama.CompressionSnappy,
		// Optional: Set timeout
		Timeout: 30 * time.Second,
		// Optional: Set client ID
		ClientID: "my-app-producer",
		// Optional: SASL authentication
		// Username: "your-username",
		// Password: "your-password",
		// Optional: Enable TLS
		// EnableTLS: true,
	}

	// Create producer
	producer, err := NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create a message
	message := &Message{
		Topic: "my-topic",
		Key:   []byte("message-key"),
		Value: []byte("Hello, Kafka!"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("source"),
				Value: []byte("example-app"),
			},
		},
		Timestamp: time.Now(),
	}

	// Send message
	ctx := context.Background()
	err = producer.Produce(ctx, message, "1")
	if err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}

	fmt.Println("Message sent successfully!")
}

// ExampleWithOptions demonstrates how to use the bi-kafka producer with options pattern
func ExampleWithOptions() {
	// Create producer using options pattern
	producer, err := NewProducerWithOptions(
		WithBootstrapServers("localhost:9092", "localhost:9093"),
		WithSnappyCompression(),
		WithTimeout(30*time.Second),
		WithWaitForAll(),
		WithClientID("my-app-producer"),
		// WithSASL("username", "password"),
		// WithTLS(),
	)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create a message
	message := &Message{
		Topic: "my-topic",
		Key:   []byte("message-key"),
		Value: []byte("Hello, Kafka with Options!"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("source"),
				Value: []byte("example-app"),
			},
			{
				Key:   []byte("version"),
				Value: []byte("1.0.0"),
			},
		},
		Timestamp: time.Now(),
	}

	// Send message
	ctx := context.Background()
	err = producer.Produce(ctx, message, "1")
	if err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}

	fmt.Println("Message sent successfully with options!")
}

// ExampleBootstrapServers demonstrates different bootstrap server configurations
func ExampleBootstrapServers() {
	// Example 1: Single bootstrap server
	producer1, err := NewProducerWithOptions(
		WithBootstrapServers("kafka-broker-1:9092"),
		WithClientID("single-bootstrap-producer"),
	)
	if err != nil {
		log.Printf("Failed to create producer 1: %v", err)
	} else {
		defer producer1.Close()
		fmt.Printf("Connected to bootstrap servers: %v\n", producer1.GetBootstrapServers())
	}

	// Example 2: Multiple bootstrap servers for high availability
	producer2, err := NewProducerWithOptions(
		WithBootstrapServers(
			"kafka-broker-1:9092",
			"kafka-broker-2:9092",
			"kafka-broker-3:9092",
		),
		WithClientID("multi-bootstrap-producer"),
		WithSnappyCompression(),
	)
	if err != nil {
		log.Printf("Failed to create producer 2: %v", err)
	} else {
		defer producer2.Close()
		fmt.Printf("Connected to bootstrap servers: %v\n", producer2.GetBootstrapServers())
	}

	// Example 3: Bootstrap servers with SASL authentication
	producer3, err := NewProducerWithOptions(
		WithBootstrapServers("secure-kafka:9092"),
		WithSASL("username", "password"),
		WithTLS(),
		WithClientID("secure-producer"),
	)
	if err != nil {
		log.Printf("Failed to create producer 3: %v", err)
	} else {
		defer producer3.Close()
		fmt.Printf("Connected to secure bootstrap servers: %v\n", producer3.GetBootstrapServers())
	}
}

// ExampleBatch demonstrates how to send multiple messages in a batch
func ExampleBatch() {
	config := &Config{
		BootstrapServers: []string{"localhost:9092"},
		ClientID:         "batch-producer",
	}

	producer, err := NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Create multiple messages
	messages := []*Message{
		{
			Topic: "batch-topic",
			Key:   []byte("key1"),
			Value: []byte("Message 1"),
		},
		{
			Topic: "batch-topic",
			Key:   []byte("key2"),
			Value: []byte("Message 2"),
		},
		{
			Topic: "batch-topic",
			Key:   []byte("key3"),
			Value: []byte("Message 3"),
		},
	}

	// Send batch
	ctx := context.Background()
	err = producer.ProduceBatch(ctx, messages)
	if err != nil {
		log.Fatalf("Failed to produce batch: %v", err)
	}

	fmt.Println("Batch sent successfully!")
}

// ExampleDifferentCompression demonstrates different compression options
func ExampleDifferentCompression() {
	// Example with different compression types
	examples := []struct {
		name    string
		options []Option
	}{
		{
			name: "Snappy Compression",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithSnappyCompression(),
			},
		},
		{
			name: "Gzip Compression",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithGzipCompression(),
			},
		},
		{
			name: "LZ4 Compression",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithLz4Compression(),
			},
		},
		{
			name: "Zstd Compression",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithZstdCompression(),
			},
		},
	}

	for _, example := range examples {
		producer, err := NewProducerWithOptions(example.options...)
		if err != nil {
			log.Printf("Failed to create producer with %s: %v", example.name, err)
			continue
		}
		defer producer.Close()

		message := &Message{
			Topic: "compression-test",
			Key:   []byte("test-key"),
			Value: []byte(fmt.Sprintf("Testing %s", example.name)),
		}

		ctx := context.Background()
		err = producer.Produce(ctx, message, "1")
		if err != nil {
			log.Printf("Failed to produce message with %s: %v", example.name, err)
		} else {
			fmt.Printf("Message sent successfully with %s\n", example.name)
		}
	}
}

// ExampleDifferentAcks demonstrates different acknowledgment levels
func ExampleDifferentAcks() {
	// Example with different acknowledgment levels
	examples := []struct {
		name    string
		options []Option
	}{
		{
			name: "Wait for All",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithWaitForAll(),
			},
		},
		{
			name: "Wait for Local",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithWaitForLocal(),
			},
		},
		{
			name: "No Response",
			options: []Option{
				WithBootstrapServers("localhost:9092"),
				WithNoResponse(),
			},
		},
	}

	for _, example := range examples {
		producer, err := NewProducerWithOptions(example.options...)
		if err != nil {
			log.Printf("Failed to create producer with %s: %v", example.name, err)
			continue
		}
		defer producer.Close()

		message := &Message{
			Topic: "acks-test",
			Key:   []byte("test-key"),
			Value: []byte(fmt.Sprintf("Testing %s", example.name)),
		}

		ctx := context.Background()
		err = producer.Produce(ctx, message, "1")
		if err != nil {
			log.Printf("Failed to produce message with %s: %v", example.name, err)
		} else {
			fmt.Printf("Message sent successfully with %s\n", example.name)
		}
	}
}
