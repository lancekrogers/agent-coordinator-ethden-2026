package coordinator

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"

	"github.com/lancekrogers/agent-coordinator-ethden-2026/internal/hedera/hcs"
)

// TaskAssignmentPayload is the payload for a task assignment message.
type TaskAssignmentPayload struct {
	TaskID       string   `json:"task_id"`
	TaskName     string   `json:"task_name"`
	AgentID      string   `json:"agent_id"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// Assigner implements the TaskAssigner interface.
type Assigner struct {
	publisher hcs.MessagePublisher
	topicID   hiero.TopicID
	agentIDs  []string

	mu          sync.RWMutex
	assignments map[string]string // taskID -> agentID
	seqNum      uint64
}

// NewAssigner creates a new task assigner.
func NewAssigner(publisher hcs.MessagePublisher, topicID hiero.TopicID, agentIDs []string) *Assigner {
	return &Assigner{
		publisher:   publisher,
		topicID:     topicID,
		agentIDs:    agentIDs,
		assignments: make(map[string]string),
	}
}

// AssignTasks publishes task assignments for all tasks in the plan.
func (a *Assigner) AssignTasks(ctx context.Context, plan Plan) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("assign tasks for plan %s: %w", plan.FestivalID, err)
	}

	var assignedIDs []string
	agentIdx := 0

	for _, seq := range plan.Sequences {
		for _, task := range seq.Tasks {
			if err := ctx.Err(); err != nil {
				return assignedIDs, fmt.Errorf("assign tasks: cancelled during assignment: %w", err)
			}

			agentID := task.AssignTo
			if agentID == "" && len(a.agentIDs) > 0 {
				agentID = a.agentIDs[agentIdx%len(a.agentIDs)]
				agentIdx++
			}

			if err := a.AssignTask(ctx, task.ID, agentID); err != nil {
				return assignedIDs, fmt.Errorf("assign tasks: task %s: %w", task.ID, err)
			}

			assignedIDs = append(assignedIDs, task.ID)
		}
	}

	return assignedIDs, nil
}

// AssignTask assigns a single task to a specific agent via HCS.
func (a *Assigner) AssignTask(ctx context.Context, taskID string, agentID string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("assign task %s to %s: %w", taskID, agentID, err)
	}

	payload := TaskAssignmentPayload{
		TaskID:  taskID,
		AgentID: agentID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("assign task %s to %s: marshal payload: %w", taskID, agentID, err)
	}

	a.mu.Lock()
	a.seqNum++
	seqNum := a.seqNum
	a.mu.Unlock()

	env := hcs.Envelope{
		Type:        hcs.MessageTypeTaskAssignment,
		Sender:      "coordinator",
		Recipient:   agentID,
		TaskID:      taskID,
		SequenceNum: seqNum,
		Timestamp:   time.Now(),
		Payload:     payloadBytes,
	}

	if err := a.publisher.Publish(ctx, a.topicID, env); err != nil {
		return fmt.Errorf("assign task %s to %s: publish: %w", taskID, agentID, err)
	}

	a.mu.Lock()
	a.assignments[taskID] = agentID
	a.mu.Unlock()

	return nil
}

// Assignment returns the agent ID assigned to a task, or empty string if unassigned.
func (a *Assigner) Assignment(taskID string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.assignments[taskID]
}

// AssignmentCount returns the number of tasks that have been assigned.
func (a *Assigner) AssignmentCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.assignments)
}

// Compile-time interface compliance check.
var _ TaskAssigner = (*Assigner)(nil)
