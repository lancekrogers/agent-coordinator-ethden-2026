package schedule

import (
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestHeartbeatConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  HeartbeatConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: HeartbeatConfig{
				Interval:  30 * time.Second,
				AgentID:   "agent-1",
				AccountID: hiero.AccountID{Account: 100},
			},
			wantErr: false,
		},
		{
			name: "interval below minimum",
			config: HeartbeatConfig{
				Interval:  1 * time.Second,
				AgentID:   "agent-1",
				AccountID: hiero.AccountID{Account: 100},
			},
			wantErr: true,
		},
		{
			name: "missing agent ID",
			config: HeartbeatConfig{
				Interval:  30 * time.Second,
				AccountID: hiero.AccountID{Account: 100},
			},
			wantErr: true,
		},
		{
			name: "missing account ID",
			config: HeartbeatConfig{
				Interval: 30 * time.Second,
				AgentID:  "agent-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHeartbeat_InvalidConfig(t *testing.T) {
	badConfig := HeartbeatConfig{Interval: 1 * time.Second}
	_, err := NewHeartbeat(nil, nil, badConfig)
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestNewHeartbeat_ValidConfig(t *testing.T) {
	cfg := HeartbeatConfig{
		Interval:  10 * time.Second,
		AgentID:   "agent-1",
		AccountID: hiero.AccountID{Account: 100},
	}
	hb, err := NewHeartbeat(nil, nil, cfg)
	if err != nil {
		t.Fatalf("NewHeartbeat() error = %v", err)
	}
	if hb == nil {
		t.Fatal("expected non-nil heartbeat")
	}
}

func TestHeartbeat_LastHeartbeat_ZeroBeforeStart(t *testing.T) {
	cfg := HeartbeatConfig{
		Interval:  10 * time.Second,
		AgentID:   "agent-1",
		AccountID: hiero.AccountID{Account: 100},
	}
	hb, err := NewHeartbeat(nil, nil, cfg)
	if err != nil {
		t.Fatalf("NewHeartbeat() error = %v", err)
	}

	if !hb.LastHeartbeat().IsZero() {
		t.Error("expected zero time before any heartbeat sent")
	}
}

func TestDefaultHeartbeatConfig(t *testing.T) {
	cfg := DefaultHeartbeatConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("Interval = %v, want 30s", cfg.Interval)
	}
	if cfg.Memo != "agent-heartbeat" {
		t.Errorf("Memo = %q, want %q", cfg.Memo, "agent-heartbeat")
	}
}

func TestHeartbeat_ImplementsInterface(t *testing.T) {
	var _ HeartbeatRunner = (*Heartbeat)(nil)
}
