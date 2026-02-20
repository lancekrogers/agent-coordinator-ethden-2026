package schedule

import (
	"context"
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// ScheduleService implements the ScheduleCreator interface using the Hiero (Hedera) SDK.
type ScheduleService struct {
	client *hiero.Client
}

// NewScheduleService creates a new ScheduleService.
func NewScheduleService(client *hiero.Client) *ScheduleService {
	return &ScheduleService{client: client}
}

// CreateSchedule creates a new scheduled transaction and returns its ID.
func (s *ScheduleService) CreateSchedule(ctx context.Context, innerTx hiero.TransactionInterface, memo string) (hiero.ScheduleID, error) {
	if err := ctx.Err(); err != nil {
		return hiero.ScheduleID{}, fmt.Errorf("create schedule: %w", err)
	}

	scheduleTx, err := hiero.NewScheduleCreateTransaction().
		SetScheduledTransaction(innerTx)
	if err != nil {
		return hiero.ScheduleID{}, fmt.Errorf("create schedule with memo %q: set inner tx: %w", memo, err)
	}

	frozen, err := scheduleTx.
		SetScheduleMemo(memo).
		FreezeWith(s.client)
	if err != nil {
		return hiero.ScheduleID{}, fmt.Errorf("create schedule with memo %q: freeze: %w", memo, err)
	}

	resp, err := frozen.Execute(s.client)
	if err != nil {
		return hiero.ScheduleID{}, fmt.Errorf("create schedule with memo %q: execute: %w", memo, err)
	}

	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return hiero.ScheduleID{}, fmt.Errorf("create schedule with memo %q: receipt: %w", memo, err)
	}

	if receipt.ScheduleID == nil {
		return hiero.ScheduleID{}, fmt.Errorf("create schedule with memo %q: receipt contained nil schedule ID", memo)
	}

	return *receipt.ScheduleID, nil
}

// ScheduleInfo retrieves information about an existing scheduled transaction.
func (s *ScheduleService) ScheduleInfo(ctx context.Context, scheduleID hiero.ScheduleID) (*ScheduleMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("schedule info %s: %w", scheduleID, err)
	}

	info, err := hiero.NewScheduleInfoQuery().
		SetScheduleID(scheduleID).
		Execute(s.client)
	if err != nil {
		return nil, fmt.Errorf("schedule info %s: execute query: %w", scheduleID, err)
	}

	return &ScheduleMetadata{
		ScheduleID: scheduleID,
		Memo:       info.Memo,
		Executed:   info.ExecutedAt != nil,
		Deleted:    info.DeletedAt != nil,
	}, nil
}

// Compile-time interface compliance check.
var _ ScheduleCreator = (*ScheduleService)(nil)
