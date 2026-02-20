package daemon

import (
	"context"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "valid default config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "empty address",
			config:  Config{Address: "", DialTimeout: 10 * time.Second, CallTimeout: 30 * time.Second},
			wantErr: true,
		},
		{
			name:    "zero dial timeout",
			config:  Config{Address: "localhost:50051", DialTimeout: 0, CallTimeout: 30 * time.Second},
			wantErr: true,
		},
		{
			name:    "zero call timeout",
			config:  Config{Address: "localhost:50051", DialTimeout: 10 * time.Second, CallTimeout: 0},
			wantErr: true,
		},
		{
			name:    "TLS enabled without cert",
			config:  Config{Address: "localhost:50051", DialTimeout: 10 * time.Second, CallTimeout: 30 * time.Second, TLSEnabled: true},
			wantErr: true,
		},
		{
			name:    "TLS enabled with cert",
			config:  Config{Address: "localhost:50051", DialTimeout: 10 * time.Second, CallTimeout: 30 * time.Second, TLSEnabled: true, TLSCertPath: "/path/to/cert.pem"},
			wantErr: false,
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

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Address != "localhost:50051" {
		t.Errorf("Address = %q, want localhost:50051", cfg.Address)
	}
	if cfg.DialTimeout != 10*time.Second {
		t.Errorf("DialTimeout = %v, want 10s", cfg.DialTimeout)
	}
	if cfg.CallTimeout != 30*time.Second {
		t.Errorf("CallTimeout = %v, want 30s", cfg.CallTimeout)
	}
	if cfg.TLSEnabled {
		t.Error("TLSEnabled should be false by default")
	}
}

func TestNewGRPCClient_InvalidConfig(t *testing.T) {
	_, err := NewGRPCClient(context.Background(), Config{})
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestNewGRPCClient_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewGRPCClient(ctx, DefaultConfig())
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestGRPCClient_Close_Nil(t *testing.T) {
	c := &GRPCClient{}
	err := c.Close()
	if err != nil {
		t.Errorf("Close() on nil conn should not error, got %v", err)
	}
}

func TestGRPCClient_InterfaceCompliance(t *testing.T) {
	var _ DaemonClient = (*GRPCClient)(nil)
}
