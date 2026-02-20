package coordinator

import "fmt"

// TaskStatus represents a task's position in the lifecycle state machine.
type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusAssigned   TaskStatus = "assigned"
	StatusInProgress TaskStatus = "in_progress"
	StatusReview     TaskStatus = "review"
	StatusComplete   TaskStatus = "complete"
	StatusPaid       TaskStatus = "paid"
	StatusFailed     TaskStatus = "failed"
)

// PaymentState represents the payment status for a task.
type PaymentState string

const (
	PaymentPending   PaymentState = "pending"
	PaymentProcessed PaymentState = "processed"
	PaymentFailed    PaymentState = "failed"
)

// validTransitions defines the allowed state transitions.
var validTransitions = map[TaskStatus][]TaskStatus{
	StatusPending:    {StatusAssigned, StatusFailed},
	StatusAssigned:   {StatusInProgress, StatusFailed},
	StatusInProgress: {StatusReview, StatusFailed},
	StatusReview:     {StatusComplete, StatusInProgress, StatusFailed},
	StatusComplete:   {StatusPaid, StatusFailed},
	StatusPaid:       {},
	StatusFailed:     {StatusPending},
}

// CanTransition checks if a state transition is valid.
func CanTransition(from, to TaskStatus) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition validates a state transition and returns an error if invalid.
func Transition(from, to TaskStatus) error {
	if !CanTransition(from, to) {
		return fmt.Errorf("invalid state transition from %s to %s", from, to)
	}
	return nil
}

// IsTerminal returns true if the status is a terminal state.
func IsTerminal(status TaskStatus) bool {
	return status == StatusPaid
}
