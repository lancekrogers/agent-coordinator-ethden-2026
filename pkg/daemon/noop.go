package daemon

import (
	"context"
	"time"
)

// noopClient satisfies DaemonClient without connecting to any daemon.
// Used for graceful degradation when no daemon is running.
type noopClient struct{}

// Noop returns a DaemonClient that silently succeeds on all operations.
// Agents use this to run standalone when no daemon is available.
func Noop() DaemonClient { return &noopClient{} }

func (n *noopClient) Register(_ context.Context, req RegisterRequest) (*RegisterResponse, error) {
	return &RegisterResponse{
		AgentID:      req.AgentName + "-standalone",
		SessionID:    "noop",
		RegisteredAt: time.Now(),
	}, nil
}

func (n *noopClient) Execute(_ context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	return &ExecuteResponse{
		TaskID:   req.TaskID,
		Status:   "skipped",
		Duration: 0,
	}, nil
}

func (n *noopClient) Heartbeat(_ context.Context, _ HeartbeatRequest) error {
	return nil
}

func (n *noopClient) Close() error {
	return nil
}

// Compile-time interface compliance check.
var _ DaemonClient = (*noopClient)(nil)
