# agent-coordinator

Coordinator agent — reads festival plans, assigns tasks via HCS, monitors progress, enforces quality gates, manages HTS payments. Uses daemon Execute RPC to run fest commands within the campaign sandbox. fest detects campaign root and festivals/ automatically.

## Build

```bash
just build   # Build binary to bin/
just run     # Run the agent
just test    # Run tests
```

## Structure

- `cmd/` — Entry point
- `internal/` — Private packages
- `justfile` — Build recipes

## Development

- Follow Go conventions from root CLAUDE.md
- Always pass context.Context as first parameter for I/O
- Use the project's error framework, not fmt.Errorf
- Keep files under 500 lines, functions under 50 lines
