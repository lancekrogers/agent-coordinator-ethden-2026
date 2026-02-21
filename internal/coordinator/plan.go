package coordinator

// Plan represents a parsed festival plan for coordinator execution.
type Plan struct {
	FestivalID string         `json:"festival_id"`
	Sequences  []PlanSequence `json:"sequences"`
}

// PlanSequence represents a sequence within the plan.
type PlanSequence struct {
	ID    string     `json:"id"`
	Tasks []PlanTask `json:"tasks"`
}

// PlanTask represents a single task in the plan.
type PlanTask struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	TaskType      string   `json:"task_type,omitempty"`
	AssignTo      string   `json:"assign_to,omitempty"`
	ModelID       string   `json:"model_id,omitempty"`
	Input         string   `json:"input,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	MaxTokens     int      `json:"max_tokens,omitempty"`
	PaymentAmount int64    `json:"payment_amount,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
}

// TaskCount returns the total number of tasks across all sequences.
func (p Plan) TaskCount() int {
	count := 0
	for _, seq := range p.Sequences {
		count += len(seq.Tasks)
	}
	return count
}

// TaskByID finds a task in the plan by its ID. Returns nil if not found.
func (p Plan) TaskByID(taskID string) *PlanTask {
	for i := range p.Sequences {
		for j := range p.Sequences[i].Tasks {
			if p.Sequences[i].Tasks[j].ID == taskID {
				return &p.Sequences[i].Tasks[j]
			}
		}
	}
	return nil
}
