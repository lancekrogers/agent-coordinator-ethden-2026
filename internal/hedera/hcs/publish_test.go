package hcs

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestPublish_ContextCancellation(t *testing.T) {
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

	pub := NewPublisher(nil, DefaultPublishConfig())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Envelope{
				Type:      MessageTypeHeartbeat,
				Sender:    "test-agent",
				Timestamp: time.Now(),
			}
			err := pub.Publish(tt.ctx, hiero.TopicID{}, env)
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestEnvelope_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		envelope Envelope
		wantType string
	}{
		{
			name: "task assignment",
			envelope: Envelope{
				Type:        MessageTypeTaskAssignment,
				Sender:      "coordinator",
				Recipient:   "agent-1",
				TaskID:      "01_link_project",
				SequenceNum: 1,
				Timestamp:   time.Date(2026, 2, 18, 14, 0, 0, 0, time.UTC),
			},
			wantType: "task_assignment",
		},
		{
			name: "heartbeat with no payload",
			envelope: Envelope{
				Type:        MessageTypeHeartbeat,
				Sender:      "agent-1",
				SequenceNum: 42,
				Timestamp:   time.Date(2026, 2, 18, 14, 0, 0, 0, time.UTC),
			},
			wantType: "heartbeat",
		},
		{
			name: "status update with payload",
			envelope: Envelope{
				Type:        MessageTypeStatusUpdate,
				Sender:      "agent-1",
				TaskID:      "03_implement",
				SequenceNum: 5,
				Timestamp:   time.Date(2026, 2, 18, 14, 0, 0, 0, time.UTC),
				Payload:     json.RawMessage(`{"progress": 50}`),
			},
			wantType: "status_update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.envelope.Marshal()
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("Unmarshal result: %v", err)
			}

			if got := result["type"]; got != tt.wantType {
				t.Errorf("type = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestEnvelope_UnmarshalRoundTrip(t *testing.T) {
	original := Envelope{
		Type:        MessageTypeTaskResult,
		Sender:      "agent-2",
		Recipient:   "coordinator",
		TaskID:      "05_build",
		SequenceNum: 10,
		Timestamp:   time.Date(2026, 2, 18, 15, 30, 0, 0, time.UTC),
		Payload:     json.RawMessage(`{"status":"success"}`),
	}

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got, err := UnmarshalEnvelope(data)
	if err != nil {
		t.Fatalf("UnmarshalEnvelope: %v", err)
	}

	if got.Type != original.Type {
		t.Errorf("Type = %v, want %v", got.Type, original.Type)
	}
	if got.Sender != original.Sender {
		t.Errorf("Sender = %v, want %v", got.Sender, original.Sender)
	}
	if got.TaskID != original.TaskID {
		t.Errorf("TaskID = %v, want %v", got.TaskID, original.TaskID)
	}
	if got.SequenceNum != original.SequenceNum {
		t.Errorf("SequenceNum = %v, want %v", got.SequenceNum, original.SequenceNum)
	}
}

func TestPublishConfig_Defaults(t *testing.T) {
	cfg := DefaultPublishConfig()
	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
	}
	if cfg.BaseBackoff != 500*time.Millisecond {
		t.Errorf("BaseBackoff = %v, want 500ms", cfg.BaseBackoff)
	}
	if cfg.MaxBackoff != 5*time.Second {
		t.Errorf("MaxBackoff = %v, want 5s", cfg.MaxBackoff)
	}
}

func TestPublisher_CalculateBackoff(t *testing.T) {
	pub := NewPublisher(nil, PublishConfig{
		BaseBackoff: 100 * time.Millisecond,
		MaxBackoff:  1 * time.Second,
	})

	tests := []struct {
		attempt int
		want    time.Duration
	}{
		{attempt: 0, want: 100 * time.Millisecond},
		{attempt: 1, want: 200 * time.Millisecond},
		{attempt: 2, want: 400 * time.Millisecond},
		{attempt: 3, want: 800 * time.Millisecond},
		{attempt: 4, want: 1 * time.Second}, // capped at MaxBackoff
	}

	for _, tt := range tests {
		got := pub.calculateBackoff(tt.attempt)
		if got != tt.want {
			t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
		}
	}
}
