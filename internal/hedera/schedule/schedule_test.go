package schedule

import (
	"context"
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestNewScheduleService(t *testing.T) {
	svc := NewScheduleService(nil)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestCreateSchedule_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{name: "cancelled context", ctx: cancelledCtx()},
		{name: "deadline exceeded", ctx: expiredCtx()},
	}

	svc := NewScheduleService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			innerTx := hiero.NewTransferTransaction()
			_, err := svc.CreateSchedule(tt.ctx, innerTx, "test-memo")
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestScheduleInfo_ContextCancellation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{name: "cancelled context", ctx: cancelledCtx()},
		{name: "deadline exceeded", ctx: expiredCtx()},
	}

	svc := NewScheduleService(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.ScheduleInfo(tt.ctx, hiero.ScheduleID{})
			if err == nil {
				t.Fatal("expected error for cancelled context")
			}
		})
	}
}

func TestScheduleService_ImplementsInterface(t *testing.T) {
	var _ ScheduleCreator = (*ScheduleService)(nil)
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
