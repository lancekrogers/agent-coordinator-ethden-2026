package schedule

import (
	"context"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// ScheduleCreator handles creating scheduled transactions on Hedera.
type ScheduleCreator interface {
	// CreateSchedule creates a new scheduled transaction and returns its ID.
	CreateSchedule(ctx context.Context, innerTx hiero.TransactionInterface, memo string) (hiero.ScheduleID, error)

	// ScheduleInfo retrieves information about an existing scheduled transaction.
	ScheduleInfo(ctx context.Context, scheduleID hiero.ScheduleID) (*ScheduleMetadata, error)
}

// HeartbeatRunner manages a periodic heartbeat using the Hedera Schedule Service.
// The heartbeat proves agent liveness by submitting scheduled transactions at
// a configurable interval.
type HeartbeatRunner interface {
	// Start begins the heartbeat loop. It blocks until the context is cancelled
	// or an unrecoverable error occurs. Non-fatal errors are sent to the
	// returned channel. The channel is closed when the runner stops.
	Start(ctx context.Context) <-chan error

	// LastHeartbeat returns the timestamp of the most recent successful heartbeat,
	// or zero time if none has been sent.
	LastHeartbeat() time.Time
}

// ScheduleMetadata holds information about a scheduled transaction.
type ScheduleMetadata struct {
	ScheduleID hiero.ScheduleID
	Memo       string
	Executed   bool
	Deleted    bool
}
