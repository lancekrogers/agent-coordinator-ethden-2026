package coordinator

import (
	"context"
)

// TaskAssigner reads a festival plan and assigns tasks to agents via HCS.
type TaskAssigner interface {
	// AssignTasks publishes task assignments for all tasks in the plan.
	AssignTasks(ctx context.Context, plan Plan) ([]string, error)

	// AssignTask assigns a single task to a specific agent.
	AssignTask(ctx context.Context, taskID string, agentID string) error
}

// ProgressMonitor listens for agent status updates and tracks task progress.
type ProgressMonitor interface {
	// Start begins monitoring the HCS topic for status updates.
	// Blocks until context is cancelled.
	Start(ctx context.Context) error

	// TaskState returns the current state of a task.
	TaskState(taskID string) (TaskStatus, error)

	// AllTaskStates returns a map of all tracked tasks and their states.
	AllTaskStates() map[string]TaskStatus
}

// QualityGateEnforcer checks whether a task has passed its quality gates.
type QualityGateEnforcer interface {
	// Evaluate checks quality gate criteria for a task.
	// Returns true if the gate is passed.
	Evaluate(ctx context.Context, taskID string) (bool, error)
}

// PaymentManager handles triggering HTS token payments for completed tasks.
type PaymentManager interface {
	// PayForTask triggers a token transfer to the agent that completed the task.
	PayForTask(ctx context.Context, taskID string, agentID string, amount int64) error

	// PaymentStatus returns the payment status for a task.
	PaymentStatus(taskID string) (PaymentState, error)
}
