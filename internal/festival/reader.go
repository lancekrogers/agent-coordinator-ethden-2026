package festival

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"
)

const defaultCommandTimeout = 8 * time.Second

// CommandRunner is an abstraction around command execution for easier testing.
type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) (stdout []byte, stderr []byte, err error)
}

// ExecRunner runs commands on the local machine.
type ExecRunner struct {
	Dir string
}

func (r ExecRunner) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	if _, err := exec.LookPath(name); err != nil {
		return nil, nil, err
	}
	cmd := exec.CommandContext(ctx, name, args...)
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	stdout, err := cmd.Output()
	if err == nil {
		return stdout, nil, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return nil, exitErr.Stderr, err
	}
	return nil, nil, err
}

// Reader fetches and parses festival data from the fest CLI.
type Reader struct {
	runner         CommandRunner
	cfg            ReaderConfig
	logger         *slog.Logger
	commandTimeout time.Duration
}

func NewReader(cfg ReaderConfig, logger *slog.Logger) *Reader {
	if logger == nil {
		logger = slog.Default()
	}
	if cfg.CommandTimeout <= 0 {
		cfg.CommandTimeout = defaultCommandTimeout
	}
	return &Reader{
		runner:         ExecRunner{Dir: cfg.RootDir},
		cfg:            cfg,
		logger:         logger,
		commandTimeout: cfg.CommandTimeout,
	}
}

func (r *Reader) SetRunner(runner CommandRunner) {
	if runner != nil {
		r.runner = runner
	}
}

func (r *Reader) runFest(ctx context.Context, args ...string) ([]byte, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, r.commandTimeout)
	defer cancel()

	start := time.Now()
	stdout, stderr, err := r.runner.Run(cmdCtx, "fest", args...)
	duration := time.Since(start)

	r.logger.Debug("fest command complete",
		"args", strings.Join(args, " "),
		"duration_ms", duration.Milliseconds(),
		"error", err,
		"stderr", string(stderr),
	)

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, ErrFestBinaryMissing
		}
		if strings.Contains(err.Error(), "executable file not found") {
			return nil, ErrFestBinaryMissing
		}
		return nil, err
	}
	return stdout, nil
}

func (r *Reader) ShowAll(ctx context.Context) (ShowAllResponse, error) {
	out, err := r.runFest(ctx, "show", "all", "--json")
	if err != nil {
		if errors.Is(err, ErrFestBinaryMissing) {
			return ShowAllResponse{}, err
		}
		return ShowAllResponse{}, fmt.Errorf("%w: %v", ErrShowAllFailed, err)
	}

	var resp ShowAllResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return ShowAllResponse{}, fmt.Errorf("%w: %v", ErrShowAllParse, err)
	}
	return resp, nil
}

func (r *Reader) ResolveSelector(ctx context.Context) (string, error) {
	if sel := strings.TrimSpace(r.cfg.Selector); sel != "" {
		return sel, nil
	}

	all, err := r.ShowAll(ctx)
	if err != nil {
		return "", err
	}

	selector := firstFestivalName(
		all.Active.Festivals,
		all.Ready.Festivals,
		all.Planning.Festivals,
		all.DungeonSomeday.Festivals,
	)
	if selector == "" && r.cfg.AllowCompleted {
		selector = firstFestivalName(all.DungeonCompleted.Festivals)
	}
	if selector == "" {
		return "", ErrSelectorNotFound
	}
	return selector, nil
}

func firstFestivalName(groups ...[]FestivalSummary) string {
	for _, group := range groups {
		if len(group) == 0 {
			continue
		}
		if strings.TrimSpace(group[0].Name) != "" {
			return group[0].Name
		}
	}
	return ""
}

func (r *Reader) ShowRoadmap(ctx context.Context, selector string) (ShowRoadmapResponse, error) {
	if strings.TrimSpace(selector) == "" {
		return ShowRoadmapResponse{}, ErrSelectorNotFound
	}

	out, err := r.runFest(ctx, "show", "--festival", selector, "--json", "--roadmap")
	if err != nil {
		if errors.Is(err, ErrFestBinaryMissing) {
			return ShowRoadmapResponse{}, err
		}
		return ShowRoadmapResponse{}, fmt.Errorf("%w: %v", ErrRoadmapFailed, err)
	}

	var resp ShowRoadmapResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return ShowRoadmapResponse{}, fmt.Errorf("%w: %v", ErrRoadmapParse, err)
	}
	return resp, nil
}

func (r *Reader) Load(ctx context.Context) (string, ShowRoadmapResponse, error) {
	selector, err := r.ResolveSelector(ctx)
	if err != nil {
		return "", ShowRoadmapResponse{}, err
	}

	roadmap, err := r.ShowRoadmap(ctx, selector)
	if err != nil {
		return "", ShowRoadmapResponse{}, err
	}

	return selector, roadmap, nil
}
