package hcs

import (
	"context"
	"fmt"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

const (
	defaultMaxRetries  = 3
	defaultBaseBackoff = 500 * time.Millisecond
	defaultMaxBackoff  = 5 * time.Second
)

// PublishConfig holds configuration for the publisher.
type PublishConfig struct {
	MaxRetries  int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
}

// DefaultPublishConfig returns sensible defaults for testnet usage.
func DefaultPublishConfig() PublishConfig {
	return PublishConfig{
		MaxRetries:  defaultMaxRetries,
		BaseBackoff: defaultBaseBackoff,
		MaxBackoff:  defaultMaxBackoff,
	}
}

// Publisher implements the MessagePublisher interface using the Hiero (Hedera) SDK.
type Publisher struct {
	client *hiero.Client
	config PublishConfig
}

// NewPublisher creates a new Publisher with the given Hiero client and config.
func NewPublisher(client *hiero.Client, config PublishConfig) *Publisher {
	return &Publisher{
		client: client,
		config: config,
	}
}

// Publish serializes the envelope to JSON and submits it to an HCS topic.
// Retries transient failures with exponential backoff.
func (p *Publisher) Publish(ctx context.Context, topicID hiero.TopicID, msg Envelope) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("publish to topic %s: %w", topicID, err)
	}

	data, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("publish to topic %s: marshal type %s: %w", topicID, msg.Type, err)
	}

	var lastErr error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("publish to topic %s attempt %d: %w", topicID, attempt+1, err)
		}

		lastErr = p.submitMessage(topicID, data)
		if lastErr == nil {
			return nil
		}

		if attempt < p.config.MaxRetries {
			backoff := p.calculateBackoff(attempt)
			select {
			case <-ctx.Done():
				return fmt.Errorf("publish to topic %s: cancelled during backoff: %w", topicID, ctx.Err())
			case <-time.After(backoff):
			}
		}
	}

	return fmt.Errorf("publish to topic %s type %s: exhausted %d attempts: %w",
		topicID, msg.Type, p.config.MaxRetries+1, lastErr)
}

func (p *Publisher) submitMessage(topicID hiero.TopicID, data []byte) error {
	tx, err := hiero.NewTopicMessageSubmitTransaction().
		SetTopicID(topicID).
		SetMessage(data).
		FreezeWith(p.client)
	if err != nil {
		return fmt.Errorf("freeze: %w", err)
	}

	resp, err := tx.Execute(p.client)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	_, err = resp.GetReceipt(p.client)
	if err != nil {
		return fmt.Errorf("receipt: %w", err)
	}

	return nil
}

func (p *Publisher) calculateBackoff(attempt int) time.Duration {
	backoff := p.config.BaseBackoff
	for i := 0; i < attempt; i++ {
		backoff *= 2
	}
	if backoff > p.config.MaxBackoff {
		backoff = p.config.MaxBackoff
	}
	return backoff
}

// Compile-time interface compliance check.
var _ MessagePublisher = (*Publisher)(nil)
