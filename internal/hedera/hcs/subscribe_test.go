package hcs

import (
	"context"
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestSubscribe_ContextCancellation(t *testing.T) {
	sub := NewSubscriber(nil, DefaultSubscribeConfig())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	msgCh, errCh := sub.Subscribe(ctx, hiero.TopicID{})

	timeout := time.After(2 * time.Second)
	select {
	case _, ok := <-msgCh:
		if ok {
			t.Error("expected message channel to be closed")
		}
	case <-timeout:
		t.Error("message channel did not close within timeout")
	}

	select {
	case <-errCh:
	case <-timeout:
		t.Error("error channel did not close within timeout")
	}
}

func TestSubscribeConfig_Defaults(t *testing.T) {
	cfg := DefaultSubscribeConfig()

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"MessageBuffer", cfg.MessageBuffer, defaultMessageBuffer},
		{"ReconnectDelay", cfg.ReconnectDelay, defaultReconnectDelay},
		{"MaxReconnects", cfg.MaxReconnects, defaultMaxReconnects},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestUnmarshalEnvelope_InvalidJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "empty bytes",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			data:    []byte("not json"),
			wantErr: true,
		},
		{
			name:    "valid envelope",
			data:    []byte(`{"type":"heartbeat","sender":"agent-1","sequence_num":1,"timestamp":"2026-02-18T14:00:00Z"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := UnmarshalEnvelope(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalEnvelope() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && env == nil {
				t.Error("expected non-nil envelope on success")
			}
		})
	}
}

func TestSubscriber_ImplementsInterface(t *testing.T) {
	var _ MessageSubscriber = (*Subscriber)(nil)
}
