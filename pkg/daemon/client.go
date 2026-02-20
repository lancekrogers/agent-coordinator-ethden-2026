package daemon

import (
	"context"
	"io"
	"time"
)

// DaemonClient defines the interface for communicating with the obey daemon.
// All three agents (coordinator, worker-1, worker-2) import and use this interface.
// The interface is intentionally small -- it covers the minimum surface needed
// for agent coordination via the daemon.
type DaemonClient interface {
	// Register registers this agent with the daemon, providing agent metadata.
	// Must be called before Execute or Heartbeat.
	// Returns the assigned agent ID from the daemon.
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)

	// Execute sends a task execution request to the daemon.
	// The daemon routes the task to the appropriate handler.
	Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error)

	// Heartbeat sends a liveness signal to the daemon.
	// Should be called periodically to maintain agent registration.
	Heartbeat(ctx context.Context, req HeartbeatRequest) error

	// Close cleanly shuts down the client connection.
	io.Closer
}

// RegisterRequest contains the data needed to register an agent with the daemon.
type RegisterRequest struct {
	// AgentName is the human-readable name for this agent.
	AgentName string

	// AgentType is the type of agent (e.g., "coordinator", "worker").
	AgentType string

	// Capabilities lists what this agent can do (e.g., "hcs", "hts", "schedule").
	Capabilities []string

	// HederaAccountID is the Hedera account this agent uses for transactions.
	HederaAccountID string
}

// RegisterResponse contains the daemon's response to a registration request.
type RegisterResponse struct {
	// AgentID is the unique identifier assigned by the daemon.
	AgentID string

	// SessionID is the session identifier for this registration.
	SessionID string

	// RegisteredAt is when the registration was accepted.
	RegisteredAt time.Time
}

// ExecuteRequest contains a task execution request.
type ExecuteRequest struct {
	// TaskID is the festival task identifier.
	TaskID string

	// TaskType describes what kind of work to perform.
	TaskType string

	// Payload is the task-specific data as JSON bytes.
	Payload []byte

	// Timeout is the maximum time the daemon should allow for execution.
	Timeout time.Duration
}

// ExecuteResponse contains the result of a task execution.
type ExecuteResponse struct {
	// TaskID echoes back the requested task.
	TaskID string

	// Status is the execution result status.
	Status string

	// Result is the task-specific result data as JSON bytes.
	Result []byte

	// Duration is how long the execution took.
	Duration time.Duration
}

// HeartbeatRequest contains the data for a heartbeat signal.
type HeartbeatRequest struct {
	// AgentID is the registered agent identifier.
	AgentID string

	// SessionID is the current session identifier.
	SessionID string

	// Timestamp is when this heartbeat was generated.
	Timestamp time.Time
}
