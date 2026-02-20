package coordinator

import (
	"context"
	"testing"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestAssigner_ContextCancellation(t *testing.T) {
	a := NewAssigner(nil, hiero.TopicID{}, []string{"agent-1"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := a.AssignTasks(ctx, Plan{})
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestAssignTask_ContextCancellation(t *testing.T) {
	a := NewAssigner(nil, hiero.TopicID{}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := a.AssignTask(ctx, "task-1", "agent-1")
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestAssigner_AssignmentTracking(t *testing.T) {
	a := NewAssigner(nil, hiero.TopicID{}, nil)

	if a.AssignmentCount() != 0 {
		t.Errorf("initial count = %d, want 0", a.AssignmentCount())
	}

	if got := a.Assignment("task-1"); got != "" {
		t.Errorf("Assignment(unassigned) = %q, want empty", got)
	}
}

func TestPlan_TaskCount(t *testing.T) {
	tests := []struct {
		name string
		plan Plan
		want int
	}{
		{
			name: "empty plan",
			plan: Plan{},
			want: 0,
		},
		{
			name: "single sequence with 3 tasks",
			plan: Plan{
				Sequences: []PlanSequence{
					{ID: "seq-1", Tasks: []PlanTask{{ID: "t1"}, {ID: "t2"}, {ID: "t3"}}},
				},
			},
			want: 3,
		},
		{
			name: "multiple sequences",
			plan: Plan{
				Sequences: []PlanSequence{
					{ID: "seq-1", Tasks: []PlanTask{{ID: "t1"}, {ID: "t2"}}},
					{ID: "seq-2", Tasks: []PlanTask{{ID: "t3"}}},
				},
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.plan.TaskCount(); got != tt.want {
				t.Errorf("TaskCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPlan_TaskByID(t *testing.T) {
	plan := Plan{
		Sequences: []PlanSequence{
			{
				ID: "seq-1",
				Tasks: []PlanTask{
					{ID: "task-1", Name: "First"},
					{ID: "task-2", Name: "Second"},
				},
			},
		},
	}

	tests := []struct {
		name   string
		taskID string
		found  bool
	}{
		{"existing task", "task-1", true},
		{"another existing", "task-2", true},
		{"non-existent", "task-99", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := plan.TaskByID(tt.taskID)
			if (result != nil) != tt.found {
				t.Errorf("TaskByID(%q) found = %v, want %v", tt.taskID, result != nil, tt.found)
			}
		})
	}
}

func TestCanTransition(t *testing.T) {
	tests := []struct {
		from, to TaskStatus
		want     bool
	}{
		{StatusPending, StatusAssigned, true},
		{StatusPending, StatusFailed, true},
		{StatusPending, StatusPaid, false},
		{StatusAssigned, StatusInProgress, true},
		{StatusInProgress, StatusReview, true},
		{StatusReview, StatusComplete, true},
		{StatusReview, StatusInProgress, true},
		{StatusComplete, StatusPaid, true},
		{StatusPaid, StatusPending, false},
		{StatusFailed, StatusPending, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			if got := CanTransition(tt.from, tt.to); got != tt.want {
				t.Errorf("CanTransition(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestTransition_Invalid(t *testing.T) {
	err := Transition(StatusPending, StatusPaid)
	if err == nil {
		t.Error("expected error for invalid transition")
	}
}

func TestIsTerminal(t *testing.T) {
	if !IsTerminal(StatusPaid) {
		t.Error("StatusPaid should be terminal")
	}
	if IsTerminal(StatusComplete) {
		t.Error("StatusComplete should not be terminal")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DefaultPaymentAmount != 100 {
		t.Errorf("DefaultPaymentAmount = %d, want 100", cfg.DefaultPaymentAmount)
	}
	if cfg.MonitorPollInterval != 5*time.Second {
		t.Errorf("MonitorPollInterval = %v, want 5s", cfg.MonitorPollInterval)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				TaskTopicID:          hiero.TopicID{Topic: 1},
				StatusTopicID:        hiero.TopicID{Topic: 2},
				PaymentTokenID:       hiero.TokenID{Token: 1},
				TreasuryAccountID:    hiero.AccountID{Account: 100},
				DefaultPaymentAmount: 100,
			},
			wantErr: false,
		},
		{
			name: "missing task topic",
			config: Config{
				StatusTopicID:        hiero.TopicID{Topic: 2},
				PaymentTokenID:       hiero.TokenID{Token: 1},
				TreasuryAccountID:    hiero.AccountID{Account: 100},
				DefaultPaymentAmount: 100,
			},
			wantErr: true,
		},
		{
			name: "zero payment amount",
			config: Config{
				TaskTopicID:          hiero.TopicID{Topic: 1},
				StatusTopicID:        hiero.TopicID{Topic: 2},
				PaymentTokenID:       hiero.TokenID{Token: 1},
				TreasuryAccountID:    hiero.AccountID{Account: 100},
				DefaultPaymentAmount: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
