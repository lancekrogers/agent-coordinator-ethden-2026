package hcs

import (
	"context"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TopicCreator handles HCS topic lifecycle operations.
// Topics are the communication channels for agent coordination.
type TopicCreator interface {
	// CreateTopic creates a new HCS topic and returns its ID.
	// The memo should describe the topic's purpose (e.g., "festival-tasks", "agent-status").
	CreateTopic(ctx context.Context, memo string) (hiero.TopicID, error)

	// DeleteTopic deletes an HCS topic by its ID.
	DeleteTopic(ctx context.Context, topicID hiero.TopicID) error

	// TopicInfo retrieves metadata about an existing topic.
	TopicInfo(ctx context.Context, topicID hiero.TopicID) (*TopicMetadata, error)
}

// MessagePublisher handles publishing messages to HCS topics.
type MessagePublisher interface {
	// Publish sends a message envelope to the specified HCS topic.
	// The message is serialized to JSON before publishing.
	Publish(ctx context.Context, topicID hiero.TopicID, msg Envelope) error
}

// MessageSubscriber handles subscribing to HCS topic messages.
type MessageSubscriber interface {
	// Subscribe starts a streaming subscription to an HCS topic.
	// Messages are delivered to the returned channel. The subscription
	// runs until the context is cancelled. The channel is closed when
	// the subscription ends.
	Subscribe(ctx context.Context, topicID hiero.TopicID) (<-chan Envelope, <-chan error)
}

// TopicMetadata holds information about an HCS topic.
type TopicMetadata struct {
	TopicID        hiero.TopicID
	Memo           string
	SequenceNumber uint64
}
