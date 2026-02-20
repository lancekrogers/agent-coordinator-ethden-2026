package coordinator

import (
	"context"
	"testing"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestMonitor_InitAndQuery(t *testing.T) {
	m := NewMonitor(nil, hiero.TopicID{}, nil)
	m.InitTask("task-1")

	status, err := m.TaskState("task-1")
	if err != nil {
		t.Fatalf("TaskState() error = %v", err)
	}
	if status != StatusPending {
		t.Errorf("TaskState() = %s, want pending", status)
	}
}

func TestMonitor_TaskState_NotTracked(t *testing.T) {
	m := NewMonitor(nil, hiero.TopicID{}, nil)
	_, err := m.TaskState("nonexistent")
	if err == nil {
		t.Error("expected error for untracked task")
	}
}

func TestMonitor_AllTaskStates(t *testing.T) {
	m := NewMonitor(nil, hiero.TopicID{}, nil)
	m.InitTask("task-1")
	m.InitTask("task-2")

	states := m.AllTaskStates()
	if len(states) != 2 {
		t.Errorf("AllTaskStates() len = %d, want 2", len(states))
	}
	if states["task-1"] != StatusPending {
		t.Errorf("task-1 status = %s, want pending", states["task-1"])
	}
}

func TestMonitor_AllTaskStates_ReturnsCopy(t *testing.T) {
	m := NewMonitor(nil, hiero.TopicID{}, nil)
	m.InitTask("task-1")

	states := m.AllTaskStates()
	states["task-1"] = StatusPaid // mutate the copy

	status, _ := m.TaskState("task-1")
	if status != StatusPending {
		t.Error("AllTaskStates should return a copy, not a reference")
	}
}

func TestMonitor_Start_ContextCancellation(t *testing.T) {
	m := NewMonitor(nil, hiero.TopicID{}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := m.Start(ctx)
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestSimpleGateEnforcer_Evaluate(t *testing.T) {
	gate := NewSimpleGateEnforcer(nil)

	tests := []struct {
		name   string
		taskID string
		want   bool
	}{
		{"quality gate task passes", "06_testing", true},
		{"review gate passes", "07_review", true},
		{"implementation task passes", "03_implement_topic", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gate.Evaluate(context.Background(), tt.taskID)
			if err != nil {
				t.Fatalf("Evaluate() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.taskID, got, tt.want)
			}
		})
	}
}

func TestSimpleGateEnforcer_ContextCancellation(t *testing.T) {
	gate := NewSimpleGateEnforcer(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := gate.Evaluate(ctx, "task-1")
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}
