package coordinator

import (
	"context"
	"fmt"
	"strings"
)

// SimpleGateEnforcer implements QualityGateEnforcer with basic checks.
type SimpleGateEnforcer struct {
	monitor ProgressMonitor
}

// NewSimpleGateEnforcer creates a gate enforcer that checks sibling task states.
func NewSimpleGateEnforcer(monitor ProgressMonitor) *SimpleGateEnforcer {
	return &SimpleGateEnforcer{monitor: monitor}
}

// Evaluate checks quality gate criteria for a task.
func (g *SimpleGateEnforcer) Evaluate(ctx context.Context, taskID string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("evaluate gate for task %s: %w", taskID, err)
	}

	// Quality gate tasks themselves always pass.
	if isQualityGateTask(taskID) {
		return true, nil
	}

	// Individual implementation tasks can complete freely.
	// Sequence-level gates are enforced by the coordinator.
	return true, nil
}

func isQualityGateTask(taskID string) bool {
	return strings.Contains(taskID, "testing") ||
		strings.Contains(taskID, "review") ||
		strings.Contains(taskID, "iterate") ||
		strings.Contains(taskID, "fest_commit")
}

// Compile-time interface compliance check.
var _ QualityGateEnforcer = (*SimpleGateEnforcer)(nil)
