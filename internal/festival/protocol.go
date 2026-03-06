package festival

import (
	"errors"
	"time"
)

var (
	ErrFestBinaryMissing = errors.New("fest binary missing")
	ErrShowAllFailed     = errors.New("fest show all failed")
	ErrShowAllParse      = errors.New("fest show all parse failed")
	ErrSelectorNotFound  = errors.New("fest selector unresolved")
	ErrRoadmapFailed     = errors.New("fest show roadmap failed")
	ErrRoadmapParse      = errors.New("fest show roadmap parse failed")
)

// ReaderConfig controls how fest commands are executed and parsed.
type ReaderConfig struct {
	RootDir        string
	Selector       string
	AllowCompleted bool
	CommandTimeout time.Duration
}

// ShowAllResponse models `fest show all --json`.
type ShowAllResponse struct {
	Active           FestivalBucket `json:"active"`
	Ready            FestivalBucket `json:"ready"`
	Planning         FestivalBucket `json:"planning"`
	DungeonSomeday   FestivalBucket `json:"dungeon/someday"`
	DungeonCompleted FestivalBucket `json:"dungeon/completed"`
	DungeonArchived  FestivalBucket `json:"dungeon/archived"`
}

type FestivalBucket struct {
	Count     int               `json:"count"`
	Festivals []FestivalSummary `json:"festivals"`
}

type FestivalSummary struct {
	ID           string        `json:"id"`
	MetadataID   string        `json:"metadata_id"`
	Name         string        `json:"name"`
	MetadataName string        `json:"metadata_name"`
	Status       string        `json:"status"`
	Path         string        `json:"path"`
	Stats        FestivalStats `json:"stats"`
}

type FestivalStats struct {
	Progress int `json:"progress"`
}

// ShowRoadmapResponse models `fest show --festival <selector> --json --roadmap`.
type ShowRoadmapResponse struct {
	Festival ShowRoadmapFestival `json:"festival"`
	Roadmap  FestivalRoadmap     `json:"roadmap"`
}

type ShowRoadmapFestival struct {
	ID           string        `json:"id"`
	MetadataID   string        `json:"metadata_id"`
	Name         string        `json:"name"`
	MetadataName string        `json:"metadata_name"`
	Status       string        `json:"status"`
	Path         string        `json:"path"`
	Stats        FestivalStats `json:"stats"`
}

type FestivalRoadmap struct {
	FestivalPath string         `json:"festival_path"`
	Phases       []RoadmapPhase `json:"phases"`
}

type RoadmapPhase struct {
	Name       string            `json:"name"`
	Path       string            `json:"path"`
	Number     int               `json:"number"`
	Status     string            `json:"status"`
	TotalTasks int               `json:"total_tasks"`
	Sequences  []RoadmapSequence `json:"sequences"`
}

type RoadmapSequence struct {
	Name       string        `json:"name"`
	Path       string        `json:"path"`
	Number     int           `json:"number"`
	Status     string        `json:"status"`
	TotalTasks int           `json:"total_tasks"`
	Steps      []RoadmapStep `json:"steps"`
}

type RoadmapStep struct {
	Number   int           `json:"number"`
	Type     string        `json:"type"`
	Parallel bool          `json:"parallel"`
	Tasks    []RoadmapTask `json:"tasks"`
}

type RoadmapTask struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Number       int      `json:"number"`
	Status       string   `json:"status"`
	IsGate       bool     `json:"is_gate"`
	Dependencies []string `json:"dependencies"`
}

// ExecutionPlan is the normalized plan produced from fest roadmap output.
type ExecutionPlan struct {
	FestivalID string              `json:"festival_id"`
	Sequences  []ExecutionSequence `json:"sequences"`
}

type ExecutionSequence struct {
	ID    string          `json:"id"`
	Tasks []ExecutionTask `json:"tasks"`
}

type ExecutionTask struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	IsGate       bool     `json:"is_gate"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// ProgressSnapshot is the canonical payload source for dashboard festival progress.
type ProgressSnapshot struct {
	Version           string           `json:"version"`
	Source            string           `json:"source"`
	Selector          string           `json:"selector"`
	SnapshotTime      time.Time        `json:"snapshot_time"`
	StaleAfterSeconds int              `json:"stale_after_seconds"`
	FestivalProgress  FestivalProgress `json:"festivalProgress"`
	FallbackReason    string           `json:"fallback_reason,omitempty"`
}

// FestivalProgress is the dashboard-compatible progress payload.
type FestivalProgress struct {
	FestivalID               string          `json:"festivalId"`
	FestivalName             string          `json:"festivalName"`
	Phases                   []FestivalPhase `json:"phases"`
	OverallCompletionPercent int             `json:"overallCompletionPercent"`
}

type FestivalPhase struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Status            string             `json:"status"`
	Sequences         []FestivalSequence `json:"sequences"`
	CompletionPercent int                `json:"completionPercent"`
}

type FestivalSequence struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Status            string         `json:"status"`
	Tasks             []FestivalTask `json:"tasks"`
	CompletionPercent int            `json:"completionPercent"`
}

type FestivalTask struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Autonomy string `json:"autonomy"`
}
