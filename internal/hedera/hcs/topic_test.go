package hcs

import (
	"context"
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestNewTopicService(t *testing.T) {
	t.Run("nil client returns non-nil service", func(t *testing.T) {
		svc := NewTopicService(nil)
		if svc == nil {
			t.Fatal("expected non-nil service")
		}
	})

	t.Run("stores client reference", func(t *testing.T) {
		svc := NewTopicService(nil)
		if svc.client != nil {
			t.Fatal("expected nil client in service")
		}
	})
}

func TestCreateTopic_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "cancelled context",
			ctx:  cancelledCtx(),
		},
		{
			name: "deadline exceeded context",
			ctx:  expiredCtx(),
		},
	}

	svc := NewTopicService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateTopic(tt.ctx, "test-memo")
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestDeleteTopic_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "cancelled context",
			ctx:  cancelledCtx(),
		},
		{
			name: "deadline exceeded context",
			ctx:  expiredCtx(),
		},
	}

	svc := NewTopicService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DeleteTopic(tt.ctx, hiero.TopicID{})
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestTopicInfo_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "cancelled context",
			ctx:  cancelledCtx(),
		},
		{
			name: "deadline exceeded context",
			ctx:  expiredCtx(),
		},
	}

	svc := NewTopicService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.TopicInfo(tt.ctx, hiero.TopicID{})
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestTopicService_ImplementsInterface(t *testing.T) {
	var _ TopicCreator = (*TopicService)(nil)
}

func cancelledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func expiredCtx() context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	return ctx
}
