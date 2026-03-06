package festival

import "time"

func normalizeStatus(status string) string {
	switch status {
	case "active", "in_progress", "running":
		return "active"
	case "completed", "complete", "done":
		return "completed"
	case "blocked":
		return "blocked"
	case "failed", "error":
		return "failed"
	default:
		return "pending"
	}
}

func computeCompletion(total, completed int) int {
	if total <= 0 {
		return 0
	}
	return int(float64(completed) * 100 / float64(total))
}

// BuildExecutionPlan converts roadmap response into a normalized plan.
func BuildExecutionPlan(resp ShowRoadmapResponse) ExecutionPlan {
	plan := ExecutionPlan{FestivalID: resp.Festival.MetadataID}
	if plan.FestivalID == "" {
		plan.FestivalID = resp.Festival.ID
	}

	for _, phase := range resp.Roadmap.Phases {
		for _, seq := range phase.Sequences {
			normSeq := ExecutionSequence{ID: seq.Name}
			for _, step := range seq.Steps {
				for _, task := range step.Tasks {
					normSeq.Tasks = append(normSeq.Tasks, ExecutionTask{
						ID:           task.ID,
						Name:         task.Name,
						Status:       normalizeStatus(task.Status),
						IsGate:       task.IsGate,
						Dependencies: task.Dependencies,
					})
				}
			}
			plan.Sequences = append(plan.Sequences, normSeq)
		}
	}

	return plan
}

// BuildProgressSnapshot maps roadmap response into a dashboard-friendly snapshot.
func BuildProgressSnapshot(resp ShowRoadmapResponse, selector string, staleAfterSeconds int) ProgressSnapshot {
	festivalID := resp.Festival.MetadataID
	if festivalID == "" {
		festivalID = resp.Festival.ID
	}

	festivalName := resp.Festival.MetadataName
	if festivalName == "" {
		festivalName = resp.Festival.Name
	}

	progress := FestivalProgress{
		FestivalID:   festivalID,
		FestivalName: festivalName,
	}

	totalPhaseTasks := 0
	totalCompletedTasks := 0

	for _, phase := range resp.Roadmap.Phases {
		normPhase := FestivalPhase{
			ID:   phase.Name,
			Name: phase.Name,
		}

		phaseTasks := 0
		phaseCompleted := 0

		for _, seq := range phase.Sequences {
			normSeq := FestivalSequence{
				ID:   seq.Name,
				Name: seq.Name,
			}

			seqTasks := 0
			seqCompleted := 0

			for _, step := range seq.Steps {
				for _, task := range step.Tasks {
					tStatus := normalizeStatus(task.Status)
					normSeq.Tasks = append(normSeq.Tasks, FestivalTask{
						ID:       task.Name,
						Name:     task.Name,
						Status:   tStatus,
						Autonomy: "medium",
					})
					seqTasks++
					if tStatus == "completed" {
						seqCompleted++
					}
				}
			}

			normSeq.CompletionPercent = computeCompletion(seqTasks, seqCompleted)
			normSeq.Status = normalizeStatus(seq.Status)
			if normSeq.Status == "pending" {
				switch {
				case seqCompleted == seqTasks && seqTasks > 0:
					normSeq.Status = "completed"
				case seqCompleted > 0:
					normSeq.Status = "active"
				}
			}

			normPhase.Sequences = append(normPhase.Sequences, normSeq)
			phaseTasks += seqTasks
			phaseCompleted += seqCompleted
		}

		normPhase.CompletionPercent = computeCompletion(phaseTasks, phaseCompleted)
		normPhase.Status = normalizeStatus(phase.Status)
		if normPhase.Status == "pending" {
			switch {
			case phaseCompleted == phaseTasks && phaseTasks > 0:
				normPhase.Status = "completed"
			case phaseCompleted > 0:
				normPhase.Status = "active"
			}
		}

		progress.Phases = append(progress.Phases, normPhase)
		totalPhaseTasks += phaseTasks
		totalCompletedTasks += phaseCompleted
	}

	if resp.Festival.Stats.Progress > 0 {
		progress.OverallCompletionPercent = resp.Festival.Stats.Progress
	} else {
		progress.OverallCompletionPercent = computeCompletion(totalPhaseTasks, totalCompletedTasks)
	}

	return ProgressSnapshot{
		Version:           "v1",
		Source:            "fest",
		Selector:          selector,
		SnapshotTime:      time.Now().UTC(),
		StaleAfterSeconds: staleAfterSeconds,
		FestivalProgress:  progress,
	}
}
