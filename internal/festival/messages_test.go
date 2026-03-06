package festival

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func roadmapFixture(t *testing.T, name string) ShowRoadmapResponse {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", "fest", name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	var resp ShowRoadmapResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		t.Fatalf("unmarshal fixture %s: %v", name, err)
	}
	return resp
}

func TestBuildExecutionPlan(t *testing.T) {
	resp := roadmapFixture(t, "show_roadmap_valid.json")
	plan := BuildExecutionPlan(resp)

	if plan.FestivalID != "FR0001" {
		t.Fatalf("festival_id = %q, want FR0001", plan.FestivalID)
	}
	if len(plan.Sequences) != 1 {
		t.Fatalf("sequences = %d, want 1", len(plan.Sequences))
	}
	if len(plan.Sequences[0].Tasks) != 3 {
		t.Fatalf("tasks = %d, want 3", len(plan.Sequences[0].Tasks))
	}
	if plan.Sequences[0].Tasks[1].Status != "active" {
		t.Fatalf("task status = %q, want active", plan.Sequences[0].Tasks[1].Status)
	}
	if !plan.Sequences[0].Tasks[2].IsGate {
		t.Fatalf("expected third task to be gate")
	}
}

func TestBuildProgressSnapshot(t *testing.T) {
	resp := roadmapFixture(t, "show_roadmap_valid.json")
	snap := BuildProgressSnapshot(resp, "fest-ready", 30)

	if snap.Source != "fest" {
		t.Fatalf("source = %q, want fest", snap.Source)
	}
	if snap.Selector != "fest-ready" {
		t.Fatalf("selector = %q, want fest-ready", snap.Selector)
	}
	if snap.FestivalProgress.FestivalID != "FR0001" {
		t.Fatalf("festival id = %q, want FR0001", snap.FestivalProgress.FestivalID)
	}
	if snap.FestivalProgress.OverallCompletionPercent != 42 {
		t.Fatalf("overall progress = %d, want 42", snap.FestivalProgress.OverallCompletionPercent)
	}
	if len(snap.FestivalProgress.Phases) != 1 {
		t.Fatalf("phases = %d, want 1", len(snap.FestivalProgress.Phases))
	}
	seq := snap.FestivalProgress.Phases[0].Sequences[0]
	if seq.CompletionPercent != 33 {
		t.Fatalf("sequence completion = %d, want 33", seq.CompletionPercent)
	}
	if seq.Tasks[0].Status != "completed" || seq.Tasks[1].Status != "active" {
		t.Fatalf("unexpected task status values: %+v", seq.Tasks)
	}
}

func TestNormalizeStatus(t *testing.T) {
	cases := map[string]string{
		"in_progress": "active",
		"running":     "active",
		"done":        "completed",
		"completed":   "completed",
		"blocked":     "blocked",
		"failed":      "failed",
		"":            "pending",
		"unknown":     "pending",
	}
	for in, want := range cases {
		got := normalizeStatus(in)
		if got != want {
			t.Fatalf("normalizeStatus(%q) = %q, want %q", in, got, want)
		}
	}
}
