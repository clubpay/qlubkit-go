package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

// ConsumerMessage is a single consumed record. It embeds the underlying kafka-go
// message so partition, offset, and commit metadata stay aligned with CommitMessages.
type ConsumerMessage struct {
	kafkago.Message
}

// Envelope is one delivery from Consume. Err is set when fetching stops (e.g.
// context cancelled, reader closed, or broker error). When Err is non-nil, Message is zero.
type Envelope struct {
	Message ConsumerMessage
	Err     error
}

// ConsumerConfig configures a group consumer backed by kafka-go Reader.
//
// Commit semantics: CommitInterval == 0 enables manual commit via CommitMessages.
// CommitInterval > 0 enables automatic commits on that interval (at-least-once; duplicates possible).
//
// IsolationLevel: zero value is ReadUncommitted; use ReadCommitted for transactional topics.
//
// ReadLagInterval: zero leaves kafka-go default (~1m). Set negative to disable lag updates.
//
// Throughput and fetch shape are controlled by the broker/client fetch settings below
// (MinBytes, MaxBytes, MaxWait, ReadBatchTimeout, ReadBackoffMin/Max), not by an
// application-side rate limiter. MaxBytes bounds how much data the broker may return per fetch;
// MaxWait bounds how long the broker waits to accumulate MinBytes.
//
// If MaxBytes is zero and both TargetBatchMessages and MessageBytesUpperBound are positive,
// MaxBytes is set to their product so fetches can hold roughly that many messages when messages
// are under the byte bound (Kafka does not guarantee an exact message count per response).
type ConsumerConfig struct {
	BootstrapServers []string
	Brokers          []string

	GroupID     string
	Topic       string
	GroupTopics []string

	CommitInterval time.Duration

	ClientID    string
	Username    string
	Password    string
	EnableTLS   bool
	DialTimeout time.Duration

	IsolationLevel kafkago.IsolationLevel

	MinBytes         int
	MaxBytes         int
	MaxWait          time.Duration
	ReadBatchTimeout time.Duration
	QueueCapacity    int

	// TargetBatchMessages and MessageBytesUpperBound optionally size MaxBytes when MaxBytes is 0.
	TargetBatchMessages    int
	MessageBytesUpperBound int

	HeartbeatInterval time.Duration
	SessionTimeout    time.Duration
	RebalanceTimeout  time.Duration
	RetentionTime     time.Duration
	StartOffset       int64

	GroupBalancers []kafkago.GroupBalancer

	ReadLagInterval time.Duration

	ReadBackoffMin time.Duration
	ReadBackoffMax time.Duration

	MaxAttempts           int
	Logger                kafkago.Logger
	ErrorLogger           kafkago.Logger
	OffsetOutOfRangeError bool
}

// Consumer wraps a kafka-go Reader. Call Consume at most once at a time — it runs a
// single pump that calls FetchMessage (Reader is not safe for concurrent fetches).
type Consumer struct {
	reader       *kafkago.Reader
	manualCommit bool

	pumpMu  sync.Mutex
	pumping bool
}

// NewConsumer validates cfg, opens a Reader, and returns a Consumer.
func NewConsumer(cfg *ConsumerConfig) (*Consumer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("kafka: ConsumerConfig cannot be nil")
	}
	if err := validateConsumerConfig(cfg); err != nil {
		return nil, err
	}

	clientID := cfg.ClientID
	if clientID == "" {
		clientID = "kafka-consumer"
	}
	dialTimeout := cfg.DialTimeout
	if dialTimeout == 0 {
		dialTimeout = 10 * time.Second
	}
	dialer := newDialer(dialerParams{
		ClientID:  clientID,
		Timeout:   dialTimeout,
		Username:  cfg.Username,
		Password:  cfg.Password,
		EnableTLS: cfg.EnableTLS,
	})

	rc, err := readerConfigFromConsumer(cfg, dialer)
	if err != nil {
		return nil, err
	}
	if err := rc.Validate(); err != nil {
		return nil, fmt.Errorf("kafka: invalid reader config: %w", err)
	}

	r := kafkago.NewReader(rc)
	return &Consumer{
		reader:       r,
		manualCommit: cfg.CommitInterval == 0,
	}, nil
}

func validateConsumerConfig(cfg *ConsumerConfig) error {
	brokers := cfg.brokers()
	if len(brokers) == 0 {
		return fmt.Errorf("kafka: at least one broker address is required")
	}
	if cfg.GroupID == "" {
		return fmt.Errorf("kafka: GroupID is required")
	}
	if cfg.Topic != "" && len(cfg.GroupTopics) > 0 {
		return fmt.Errorf("kafka: set either Topic or GroupTopics, not both")
	}
	if cfg.Topic == "" && len(cfg.GroupTopics) == 0 {
		return fmt.Errorf("kafka: Topic or GroupTopics is required when using GroupID")
	}
	if cfg.StartOffset != 0 && cfg.StartOffset != kafkago.FirstOffset && cfg.StartOffset != kafkago.LastOffset {
		return fmt.Errorf("kafka: StartOffset must be 0, FirstOffset, or LastOffset")
	}
	u, p := cfg.Username, cfg.Password
	if (u != "") != (p != "") {
		return fmt.Errorf("kafka: Username and Password must both be set for SASL PLAIN")
	}
	if cfg.MaxBytes <= 0 && (cfg.TargetBatchMessages > 0) != (cfg.MessageBytesUpperBound > 0) {
		return fmt.Errorf("kafka: TargetBatchMessages and MessageBytesUpperBound must both be set (or both zero) when MaxBytes is unset")
	}
	if cfg.TargetBatchMessages < 0 || cfg.MessageBytesUpperBound < 0 {
		return fmt.Errorf("kafka: TargetBatchMessages and MessageBytesUpperBound must not be negative")
	}
	if cfg.MinBytes > 0 && cfg.MaxBytes > 0 && cfg.MinBytes > cfg.MaxBytes {
		return fmt.Errorf("kafka: MinBytes must not exceed MaxBytes")
	}
	maxBytes := effectiveMaxBytes(cfg)
	if cfg.MinBytes > 0 && maxBytes > 0 && cfg.MinBytes > maxBytes {
		return fmt.Errorf("kafka: MinBytes must not exceed MaxBytes")
	}
	return nil
}

func effectiveMaxBytes(cfg *ConsumerConfig) int {
	if cfg.MaxBytes > 0 {
		return cfg.MaxBytes
	}
	if cfg.TargetBatchMessages > 0 && cfg.MessageBytesUpperBound > 0 {
		return cfg.TargetBatchMessages * cfg.MessageBytesUpperBound
	}
	return 0
}

func (c *ConsumerConfig) brokers() []string {
	if len(c.BootstrapServers) > 0 {
		return c.BootstrapServers
	}
	return c.Brokers
}

func readerConfigFromConsumer(cfg *ConsumerConfig, dialer *kafkago.Dialer) (kafkago.ReaderConfig, error) {
	rc := kafkago.ReaderConfig{
		Brokers:               cfg.brokers(),
		GroupID:               cfg.GroupID,
		Topic:                 cfg.Topic,
		GroupTopics:           append([]string(nil), cfg.GroupTopics...),
		Dialer:                dialer,
		CommitInterval:        cfg.CommitInterval,
		IsolationLevel:        cfg.IsolationLevel,
		HeartbeatInterval:     cfg.HeartbeatInterval,
		SessionTimeout:        cfg.SessionTimeout,
		RebalanceTimeout:      cfg.RebalanceTimeout,
		RetentionTime:         cfg.RetentionTime,
		StartOffset:           cfg.StartOffset,
		Logger:                cfg.Logger,
		ErrorLogger:           cfg.ErrorLogger,
		OffsetOutOfRangeError: cfg.OffsetOutOfRangeError,
	}
	if len(cfg.GroupBalancers) > 0 {
		rc.GroupBalancers = cfg.GroupBalancers
	}
	if cfg.MinBytes > 0 {
		rc.MinBytes = cfg.MinBytes
	}
	maxBytes := effectiveMaxBytes(cfg)
	if maxBytes > 0 {
		rc.MaxBytes = maxBytes
	}
	if cfg.MaxWait > 0 {
		rc.MaxWait = cfg.MaxWait
	}
	if cfg.ReadBatchTimeout > 0 {
		rc.ReadBatchTimeout = cfg.ReadBatchTimeout
	}
	if cfg.QueueCapacity > 0 {
		rc.QueueCapacity = cfg.QueueCapacity
	}
	if cfg.ReadLagInterval != 0 {
		rc.ReadLagInterval = cfg.ReadLagInterval
	}
	if cfg.ReadBackoffMin > 0 {
		rc.ReadBackoffMin = cfg.ReadBackoffMin
	}
	if cfg.ReadBackoffMax > 0 {
		rc.ReadBackoffMax = cfg.ReadBackoffMax
	}
	if cfg.MaxAttempts > 0 {
		rc.MaxAttempts = cfg.MaxAttempts
	}
	return rc, nil
}

// Consume starts a single goroutine that reads with FetchMessage and sends results on the
// returned channel. The channel is closed when the pump stops. Call Consume only once at a
// time for a given Consumer (a second call returns an error while the first pump is running).
//
// buffer is the channel capacity (minimum 1). Cancel ctx to stop the pump; the final Envelope
// may carry ctx.Err() or the error from closing the reader.
func (c *Consumer) Consume(ctx context.Context, buffer int) (<-chan Envelope, error) {
	if buffer < 1 {
		buffer = 1
	}
	c.pumpMu.Lock()
	if c.pumping {
		c.pumpMu.Unlock()
		return nil, fmt.Errorf("kafka: Consume already active for this consumer")
	}
	c.pumping = true
	c.pumpMu.Unlock()

	out := make(chan Envelope, buffer)
	go func() {
		defer close(out)
		defer func() {
			c.pumpMu.Lock()
			c.pumping = false
			c.pumpMu.Unlock()
		}()
		for {
			if ctx.Err() != nil {
				return
			}
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				select {
				case out <- Envelope{Err: err}:
				case <-ctx.Done():
				}
				return
			}
			select {
			case out <- Envelope{Message: ConsumerMessage{Message: m}}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

// Fetch reads one message synchronously (blocks until a message, error, or ctx cancellation).
// Do not call Fetch concurrently with an active Consume pump on the same Consumer.
func (c *Consumer) Fetch(ctx context.Context) (ConsumerMessage, error) {
	c.pumpMu.Lock()
	if c.pumping {
		c.pumpMu.Unlock()
		return ConsumerMessage{}, fmt.Errorf("kafka: Fetch cannot run while Consume is active")
	}
	c.pumpMu.Unlock()

	m, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return ConsumerMessage{}, err
	}
	return ConsumerMessage{Message: m}, nil
}

// CommitMessages commits offsets for the given messages. It returns an error if
// the consumer was created with CommitInterval > 0 (auto commit mode).
func (c *Consumer) CommitMessages(ctx context.Context, msgs ...ConsumerMessage) error {
	if !c.manualCommit {
		return fmt.Errorf("kafka: CommitMessages only allowed when CommitInterval is 0 (manual commit mode)")
	}
	if len(msgs) == 0 {
		return nil
	}
	raw := make([]kafkago.Message, len(msgs))
	for i := range msgs {
		raw[i] = msgs[i].Message
	}
	return c.reader.CommitMessages(ctx, raw...)
}

// Stats returns reader statistics including lag when ReadLagInterval is enabled.
func (c *Consumer) Stats() kafkago.ReaderStats {
	return c.reader.Stats()
}

// Close releases the underlying reader. If Consume is running, the pump exits when
// FetchMessage returns after the reader is closed.
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// IsEOF reports whether err indicates a closed consumer / end of stream.
func IsEOF(err error) bool {
	return errors.Is(err, io.EOF)
}
