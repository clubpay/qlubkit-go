package kafka

import (
	"context"
	"testing"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConsumer_nilConfig(t *testing.T) {
	c, err := NewConsumer(nil)
	assert.Nil(t, c)
	assert.Error(t, err)
}

func TestValidateConsumerConfig(t *testing.T) {
	base := func() *ConsumerConfig {
		return &ConsumerConfig{
			BootstrapServers: []string{"localhost:9092"},
			GroupID:          "g",
			Topic:            "orders",
		}
	}

	t.Run("missing brokers", func(t *testing.T) {
		cfg := base()
		cfg.BootstrapServers = nil
		cfg.Brokers = nil
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("missing group id", func(t *testing.T) {
		cfg := base()
		cfg.GroupID = ""
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("topic and group topics both set", func(t *testing.T) {
		cfg := base()
		cfg.GroupTopics = []string{"a"}
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("neither topic nor group topics", func(t *testing.T) {
		cfg := base()
		cfg.Topic = ""
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("invalid start offset", func(t *testing.T) {
		cfg := base()
		cfg.StartOffset = 123
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("SASL username without password", func(t *testing.T) {
		cfg := base()
		cfg.Username = "u"
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("target batch only one field set", func(t *testing.T) {
		cfg := base()
		cfg.TargetBatchMessages = 10
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("min bytes greater than max", func(t *testing.T) {
		cfg := base()
		cfg.MinBytes = 100
		cfg.MaxBytes = 50
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})

	t.Run("min bytes greater than derived max", func(t *testing.T) {
		cfg := base()
		cfg.MinBytes = 5000
		cfg.TargetBatchMessages = 10
		cfg.MessageBytesUpperBound = 100
		_, err := NewConsumer(cfg)
		assert.Error(t, err)
	})
}

func TestReaderConfigFromConsumer(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers:  []string{"127.0.0.1:9092"},
		GroupID:           "grp",
		Topic:             "t",
		CommitInterval:    2 * time.Second,
		MinBytes:          512,
		MaxBytes:          2e6,
		MaxWait:           3 * time.Second,
		ReadBatchTimeout:  4 * time.Second,
		QueueCapacity:     200,
		SessionTimeout:    25 * time.Second,
		RebalanceTimeout:  30 * time.Second,
		RetentionTime:     time.Hour,
		HeartbeatInterval: 2 * time.Second,
		StartOffset:       kafkago.FirstOffset,
		IsolationLevel:    kafkago.ReadCommitted,
		ReadLagInterval:   -1,
		ReadBackoffMin:    50 * time.Millisecond,
		ReadBackoffMax:    500 * time.Millisecond,
		MaxAttempts:       5,
		GroupBalancers:    []kafkago.GroupBalancer{kafkago.RoundRobinGroupBalancer{}},
	}
	d := newDialer(dialerParams{ClientID: "cid", Timeout: time.Second})
	rc, err := readerConfigFromConsumer(cfg, d)
	require.NoError(t, err)

	assert.Equal(t, cfg.brokers(), rc.Brokers)
	assert.Equal(t, cfg.GroupID, rc.GroupID)
	assert.Equal(t, cfg.Topic, rc.Topic)
	assert.Equal(t, cfg.CommitInterval, rc.CommitInterval)
	assert.Equal(t, 512, rc.MinBytes)
	assert.Equal(t, int(2e6), rc.MaxBytes)
	assert.Equal(t, 3*time.Second, rc.MaxWait)
	assert.Equal(t, 4*time.Second, rc.ReadBatchTimeout)
	assert.Equal(t, 200, rc.QueueCapacity)
	assert.Equal(t, 25*time.Second, rc.SessionTimeout)
	assert.Equal(t, 30*time.Second, rc.RebalanceTimeout)
	assert.Equal(t, time.Hour, rc.RetentionTime)
	assert.Equal(t, 2*time.Second, rc.HeartbeatInterval)
	assert.Equal(t, kafkago.FirstOffset, rc.StartOffset)
	assert.Equal(t, kafkago.ReadCommitted, rc.IsolationLevel)
	assert.Equal(t, time.Duration(-1), rc.ReadLagInterval)
	assert.Equal(t, 50*time.Millisecond, rc.ReadBackoffMin)
	assert.Equal(t, 500*time.Millisecond, rc.ReadBackoffMax)
	assert.Equal(t, 5, rc.MaxAttempts)
	require.Len(t, rc.GroupBalancers, 1)
	_, ok := rc.GroupBalancers[0].(kafkago.RoundRobinGroupBalancer)
	assert.True(t, ok)
	assert.Same(t, d, rc.Dialer)

	require.NoError(t, rc.Validate())
}

func TestReaderConfigFromConsumer_derivedMaxBytes(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers:       []string{"127.0.0.1:9092"},
		GroupID:                "g",
		Topic:                  "t",
		TargetBatchMessages:    50,
		MessageBytesUpperBound: 2048,
	}
	d := newDialer(dialerParams{ClientID: "c", Timeout: time.Second})
	rc, err := readerConfigFromConsumer(cfg, d)
	require.NoError(t, err)
	assert.Equal(t, 50*2048, rc.MaxBytes)
	require.NoError(t, rc.Validate())
}

func TestConsumer_CommitMessages_autoModeRejected(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers: []string{"localhost:9092"},
		GroupID:          "g",
		Topic:            "t",
		CommitInterval:   time.Second,
	}
	c, err := NewConsumer(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	err = c.CommitMessages(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "manual commit")
}

func TestConsumer_Stats(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers: []string{"localhost:9092"},
		GroupID:          "g",
		Topic:            "t",
	}
	c, err := NewConsumer(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	st := c.Stats()
	_ = st.Lag
	assert.NotEmpty(t, st.ClientID)
}

func TestConsume_doubleStart(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers: []string{"localhost:9092"},
		GroupID:          "g",
		Topic:            "t",
	}
	c, err := NewConsumer(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.Consume(ctx, 4)
	require.NoError(t, err)
	_, err = c.Consume(ctx, 4)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already active")

	cancel()
	for range ch {
	}
}

func TestConsume_ctxAlreadyCancelled(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers: []string{"localhost:9092"},
		GroupID:          "g",
		Topic:            "t",
	}
	c, err := NewConsumer(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ch, err := c.Consume(ctx, 2)
	require.NoError(t, err)
	_, ok := <-ch
	assert.False(t, ok, "channel should close without messages when ctx is done")
}

func TestFetch_whileConsumeActive(t *testing.T) {
	cfg := &ConsumerConfig{
		BootstrapServers: []string{"localhost:9092"},
		GroupID:          "g",
		Topic:            "t",
	}
	c, err := NewConsumer(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = c.Consume(ctx, 2)
	require.NoError(t, err)

	_, err = c.Fetch(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Consume is active")
}

func TestDialer_SASLAndTLS(t *testing.T) {
	d := newDialer(dialerParams{
		ClientID:  "c",
		Timeout:   time.Second,
		Username:  "u",
		Password:  "p",
		EnableTLS: true,
	})
	require.NotNil(t, d.TLS)
	require.NotNil(t, d.SASLMechanism)
}
