package hts

import (
	"context"
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestNewTokenService(t *testing.T) {
	svc := NewTokenService(nil)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestCreateFungibleToken_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "cancelled context",
			ctx:  cancelledCtx(),
		},
		{
			name: "deadline exceeded",
			ctx:  expiredCtx(),
		},
	}

	svc := NewTokenService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateFungibleToken(tt.ctx, DefaultTokenConfig())
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestTokenInfo_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{
			name: "cancelled context",
			ctx:  cancelledCtx(),
		},
		{
			name: "deadline exceeded",
			ctx:  expiredCtx(),
		},
	}

	svc := NewTokenService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.TokenInfo(tt.ctx, hiero.TokenID{})
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestDefaultTokenConfig(t *testing.T) {
	cfg := DefaultTokenConfig()

	if cfg.Name != "Agent Payment Token" {
		t.Errorf("Name = %q, want %q", cfg.Name, "Agent Payment Token")
	}
	if cfg.Symbol != "APT" {
		t.Errorf("Symbol = %q, want %q", cfg.Symbol, "APT")
	}
	if cfg.Decimals != 0 {
		t.Errorf("Decimals = %d, want 0", cfg.Decimals)
	}
	if cfg.InitialSupply != 1000000 {
		t.Errorf("InitialSupply = %d, want 1000000", cfg.InitialSupply)
	}
}

func TestTokenService_ImplementsInterface(t *testing.T) {
	var _ TokenCreator = (*TokenService)(nil)
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
