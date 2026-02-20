package hts

import (
	"context"
	"testing"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestTransfer_ContextCancellation(t *testing.T) {
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

	svc := NewTransferService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := TransferRequest{
				TokenID:       hiero.TokenID{Token: 1},
				FromAccountID: hiero.AccountID{Account: 100},
				ToAccountID:   hiero.AccountID{Account: 200},
				Amount:        10,
			}
			_, err := svc.Transfer(tt.ctx, req)
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestTransfer_InvalidAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
	}{
		{name: "zero amount", amount: 0},
		{name: "negative amount", amount: -10},
	}

	svc := NewTransferService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := TransferRequest{
				TokenID:       hiero.TokenID{Token: 1},
				FromAccountID: hiero.AccountID{Account: 100},
				ToAccountID:   hiero.AccountID{Account: 200},
				Amount:        tt.amount,
			}
			_, err := svc.Transfer(context.Background(), req)
			if err == nil {
				t.Fatalf("expected error for amount %d", tt.amount)
			}
		})
	}
}

func TestAssociateToken_ContextCancellation(t *testing.T) {
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

	svc := NewTransferService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.AssociateToken(tt.ctx, hiero.TokenID{}, hiero.AccountID{})
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestTransferService_ImplementsInterface(t *testing.T) {
	var _ TokenTransfer = (*TransferService)(nil)
}
