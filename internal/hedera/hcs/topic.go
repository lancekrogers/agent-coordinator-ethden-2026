package hcs

import (
	"context"
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TopicService implements the TopicCreator interface using the Hiero (Hedera) SDK.
type TopicService struct {
	client *hiero.Client
}

// NewTopicService creates a new TopicService with the given Hiero client.
// The client must be configured with operator credentials for the target network.
func NewTopicService(client *hiero.Client) *TopicService {
	return &TopicService{client: client}
}

// CreateTopic creates a new HCS topic with the given memo and returns its ID.
func (s *TopicService) CreateTopic(ctx context.Context, memo string) (hiero.TopicID, error) {
	if err := ctx.Err(); err != nil {
		return hiero.TopicID{}, fmt.Errorf("create topic: %w", err)
	}

	tx, err := hiero.NewTopicCreateTransaction().
		SetTopicMemo(memo).
		FreezeWith(s.client)
	if err != nil {
		return hiero.TopicID{}, fmt.Errorf("create topic with memo %q: freeze: %w", memo, err)
	}

	resp, err := tx.Execute(s.client)
	if err != nil {
		return hiero.TopicID{}, fmt.Errorf("create topic with memo %q: execute: %w", memo, err)
	}

	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return hiero.TopicID{}, fmt.Errorf("create topic with memo %q: receipt: %w", memo, err)
	}

	if receipt.TopicID == nil {
		return hiero.TopicID{}, fmt.Errorf("create topic with memo %q: receipt contained nil topic ID", memo)
	}

	return *receipt.TopicID, nil
}

// DeleteTopic deletes an HCS topic by its ID.
func (s *TopicService) DeleteTopic(ctx context.Context, topicID hiero.TopicID) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("delete topic %s: %w", topicID, err)
	}

	tx, err := hiero.NewTopicDeleteTransaction().
		SetTopicID(topicID).
		FreezeWith(s.client)
	if err != nil {
		return fmt.Errorf("delete topic %s: freeze: %w", topicID, err)
	}

	resp, err := tx.Execute(s.client)
	if err != nil {
		return fmt.Errorf("delete topic %s: execute: %w", topicID, err)
	}

	_, err = resp.GetReceipt(s.client)
	if err != nil {
		return fmt.Errorf("delete topic %s: receipt: %w", topicID, err)
	}

	return nil
}

// TopicInfo retrieves metadata about an existing HCS topic.
func (s *TopicService) TopicInfo(ctx context.Context, topicID hiero.TopicID) (*TopicMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("topic info %s: %w", topicID, err)
	}

	info, err := hiero.NewTopicInfoQuery().
		SetTopicID(topicID).
		Execute(s.client)
	if err != nil {
		return nil, fmt.Errorf("topic info %s: execute query: %w", topicID, err)
	}

	return &TopicMetadata{
		TopicID:        topicID,
		Memo:           info.TopicMemo,
		SequenceNumber: info.SequenceNumber,
	}, nil
}

// Compile-time interface compliance check.
var _ TopicCreator = (*TopicService)(nil)
