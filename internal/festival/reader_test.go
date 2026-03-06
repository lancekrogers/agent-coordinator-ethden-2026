package festival

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"log/slog"
)

type commandResult struct {
	stdout []byte
	stderr []byte
	err    error
}

type mockRunner struct {
	responses map[string]commandResult
}

func (m mockRunner) Run(_ context.Context, _ string, args ...string) ([]byte, []byte, error) {
	key := strings.Join(args, " ")
	res, ok := m.responses[key]
	if !ok {
		return nil, nil, errors.New("unexpected command: " + key)
	}
	return res.stdout, res.stderr, res.err
}

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", "fest", name)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}

func TestResolveSelector_UsesConfiguredSelector(t *testing.T) {
	r := NewReader(ReaderConfig{Selector: "explicit-selector"}, slog.Default())
	r.SetRunner(mockRunner{})

	sel, err := r.ResolveSelector(context.Background())
	if err != nil {
		t.Fatalf("ResolveSelector error: %v", err)
	}
	if sel != "explicit-selector" {
		t.Fatalf("selector = %q, want explicit-selector", sel)
	}
}

func TestResolveSelector_PriorityFallback(t *testing.T) {
	r := NewReader(ReaderConfig{}, slog.Default())
	r.SetRunner(mockRunner{responses: map[string]commandResult{
		"show all --json": {stdout: fixture(t, "show_all_no_active.json")},
	}})

	sel, err := r.ResolveSelector(context.Background())
	if err != nil {
		t.Fatalf("ResolveSelector error: %v", err)
	}
	if sel != "fest-ready" {
		t.Fatalf("selector = %q, want fest-ready", sel)
	}
}

func TestResolveSelector_AllowCompleted(t *testing.T) {
	r := NewReader(ReaderConfig{AllowCompleted: true}, slog.Default())
	r.SetRunner(mockRunner{responses: map[string]commandResult{
		"show all --json": {
			stdout: []byte(`{"active":{"count":0,"festivals":[]},"ready":{"count":0,"festivals":[]},"planning":{"count":0,"festivals":[]},"dungeon/someday":{"count":0,"festivals":[]},"dungeon/completed":{"count":1,"festivals":[{"name":"completed-fest"}]}}`),
		},
	}})

	sel, err := r.ResolveSelector(context.Background())
	if err != nil {
		t.Fatalf("ResolveSelector error: %v", err)
	}
	if sel != "completed-fest" {
		t.Fatalf("selector = %q, want completed-fest", sel)
	}
}

func TestResolveSelector_NotFound(t *testing.T) {
	r := NewReader(ReaderConfig{}, slog.Default())
	r.SetRunner(mockRunner{responses: map[string]commandResult{
		"show all --json": {stdout: []byte(`{"active":{"count":0,"festivals":[]},"ready":{"count":0,"festivals":[]},"planning":{"count":0,"festivals":[]},"dungeon/someday":{"count":0,"festivals":[]},"dungeon/completed":{"count":0,"festivals":[]}}`)},
	}})

	_, err := r.ResolveSelector(context.Background())
	if !errors.Is(err, ErrSelectorNotFound) {
		t.Fatalf("ResolveSelector err = %v, want %v", err, ErrSelectorNotFound)
	}
}

func TestShowAll_BinaryMissing(t *testing.T) {
	r := NewReader(ReaderConfig{}, slog.Default())
	r.SetRunner(mockRunner{responses: map[string]commandResult{
		"show all --json": {err: execErrNotFound{}},
	}})

	_, err := r.ShowAll(context.Background())
	if !errors.Is(err, ErrFestBinaryMissing) {
		t.Fatalf("ShowAll err = %v, want %v", err, ErrFestBinaryMissing)
	}
}

type execErrNotFound struct{}

func (execErrNotFound) Error() string { return "executable file not found in $PATH" }

func TestShowRoadmap_ParseFailure(t *testing.T) {
	r := NewReader(ReaderConfig{}, slog.Default())
	r.SetRunner(mockRunner{responses: map[string]commandResult{
		"show --festival fest-ready --json --roadmap": {stdout: fixture(t, "show_roadmap_malformed.json")},
	}})

	_, err := r.ShowRoadmap(context.Background(), "fest-ready")
	if !errors.Is(err, ErrRoadmapParse) {
		t.Fatalf("ShowRoadmap err = %v, want %v", err, ErrRoadmapParse)
	}
}

func TestLoad_ReturnsSelectorAndRoadmap(t *testing.T) {
	r := NewReader(ReaderConfig{}, slog.Default())
	r.SetRunner(mockRunner{responses: map[string]commandResult{
		"show all --json": {stdout: fixture(t, "show_all_no_active.json")},
		"show --festival fest-ready --json --roadmap": {stdout: fixture(t, "show_roadmap_valid.json")},
	}})

	sel, resp, err := r.Load(context.Background())
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if sel != "fest-ready" {
		t.Fatalf("selector = %q, want fest-ready", sel)
	}
	if resp.Festival.MetadataID != "FR0001" {
		t.Fatalf("metadata_id = %q, want FR0001", resp.Festival.MetadataID)
	}
}
