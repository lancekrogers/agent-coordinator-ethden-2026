package coordinator

import (
	"context"
	"testing"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestPayment_PayForTask_ContextCancellation(t *testing.T) {
	p := NewPayment(nil, nil, DefaultConfig())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := p.PayForTask(ctx, "task-1", "0.0.100", 100)
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestPayment_PayForTask_InvalidAmount(t *testing.T) {
	p := NewPayment(nil, nil, DefaultConfig())

	tests := []struct {
		name   string
		amount int64
	}{
		{"zero amount", 0},
		{"negative amount", -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.PayForTask(context.Background(), "task-1", "0.0.100", tt.amount)
			if err == nil {
				t.Error("expected error for invalid amount")
			}
		})
	}
}

func TestPayment_PayForTask_InvalidAgentID(t *testing.T) {
	p := NewPayment(nil, nil, DefaultConfig())

	err := p.PayForTask(context.Background(), "task-1", "not-a-valid-account", 100)
	if err == nil {
		t.Error("expected error for invalid agent account ID")
	}

	state, stateErr := p.PaymentStatus("task-1")
	if stateErr != nil {
		t.Fatalf("PaymentStatus() error = %v", stateErr)
	}
	if state != PaymentFailed {
		t.Errorf("payment state = %s, want failed", state)
	}
}

func TestPayment_PaymentStatus_NotTracked(t *testing.T) {
	p := NewPayment(nil, nil, DefaultConfig())

	_, err := p.PaymentStatus("nonexistent")
	if err == nil {
		t.Error("expected error for untracked task")
	}
}

func TestPayment_PaymentStatus_TracksState(t *testing.T) {
	p := NewPayment(nil, nil, DefaultConfig())

	// Trigger a payment that will fail at parsing to set the state.
	_ = p.PayForTask(context.Background(), "task-1", "invalid-id", 100)

	state, err := p.PaymentStatus("task-1")
	if err != nil {
		t.Fatalf("PaymentStatus() error = %v", err)
	}
	if state != PaymentFailed {
		t.Errorf("payment state = %s, want failed", state)
	}
}

func TestPayment_DoublePay_Rejected(t *testing.T) {
	p := NewPayment(nil, nil, DefaultConfig())

	// Manually set a task to processed.
	p.mu.Lock()
	p.payments["task-1"] = PaymentProcessed
	p.mu.Unlock()

	err := p.PayForTask(context.Background(), "task-1", "0.0.100", 100)
	if err == nil {
		t.Error("expected error for double payment")
	}
}

func TestNewPayment_InitializesMap(t *testing.T) {
	cfg := Config{
		PaymentTokenID:    hiero.TokenID{Token: 1},
		TreasuryAccountID: hiero.AccountID{Account: 100},
	}
	p := NewPayment(nil, nil, cfg)

	if p.payments == nil {
		t.Error("payments map should be initialized")
	}
	if p.config.PaymentTokenID.Token != 1 {
		t.Errorf("config token ID = %d, want 1", p.config.PaymentTokenID.Token)
	}
}

func TestPayment_InterfaceCompliance(t *testing.T) {
	var _ PaymentManager = (*Payment)(nil)
}
