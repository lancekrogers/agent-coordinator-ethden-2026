package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// Heartbeat implements the HeartbeatRunner interface.
// It periodically creates a scheduled transaction as a liveness proof.
type Heartbeat struct {
	client    *hiero.Client
	scheduler ScheduleCreator
	config    HeartbeatConfig

	mu            sync.RWMutex
	lastHeartbeat time.Time
}

// NewHeartbeat creates a new Heartbeat runner. Returns an error if config is invalid.
func NewHeartbeat(client *hiero.Client, scheduler ScheduleCreator, config HeartbeatConfig) (*Heartbeat, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid heartbeat config: %w", err)
	}

	return &Heartbeat{
		client:    client,
		scheduler: scheduler,
		config:    config,
	}, nil
}

// Start begins the heartbeat loop. It blocks until the context is cancelled.
// Non-fatal errors are sent to the returned channel.
func (h *Heartbeat) Start(ctx context.Context) <-chan error {
	errCh := make(chan error, 10)
	go h.run(ctx, errCh)
	return errCh
}

// LastHeartbeat returns the timestamp of the most recent successful heartbeat.
func (h *Heartbeat) LastHeartbeat() time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.lastHeartbeat
}

func (h *Heartbeat) run(ctx context.Context, errCh chan<- error) {
	defer close(errCh)

	ticker := time.NewTicker(h.config.Interval)
	defer ticker.Stop()

	h.sendHeartbeat(ctx, errCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.sendHeartbeat(ctx, errCh)
		}
	}
}

func (h *Heartbeat) sendHeartbeat(ctx context.Context, errCh chan<- error) {
	if ctx.Err() != nil {
		return
	}

	innerTx := hiero.NewTransferTransaction().
		AddHbarTransfer(h.config.AccountID, hiero.NewHbar(0))

	memo := fmt.Sprintf("%s:%s:%d", h.config.Memo, h.config.AgentID, time.Now().Unix())

	_, err := h.scheduler.CreateSchedule(ctx, innerTx, memo)
	if err != nil {
		select {
		case errCh <- fmt.Errorf("heartbeat for agent %s: %w", h.config.AgentID, err):
		default:
		}
		return
	}

	h.mu.Lock()
	h.lastHeartbeat = time.Now()
	h.mu.Unlock()
}

// Compile-time interface compliance check.
var _ HeartbeatRunner = (*Heartbeat)(nil)
