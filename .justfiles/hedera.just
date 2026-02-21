# Hedera testnet recipes
# Requires .env with HEDERA_COORDINATOR_*, HEDERA_AGENT1_*, HEDERA_AGENT2_* vars

root := justfile_directory()
env_file := root / ".env"

@default:
    just --list --justfile {{source_file()}}

# Create HCS topics + HTS token, write IDs to .env
[working-directory: '..']
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Creating HCS topics and HTS payment token..."
    output=$(go run ./cmd/setup-testnet)
    echo "$output"
    # Append or update the IDs in .env
    while IFS='=' read -r key value; do
        if grep -q "^${key}=" "{{env_file}}" 2>/dev/null; then
            sed -i '' "s|^${key}=.*|${key}=${value}|" "{{env_file}}"
        else
            echo "${key}=${value}" >> "{{env_file}}"
        fi
    done <<< "$output"
    echo "Updated .env with topic and token IDs"
    grep -E 'HCS_|HTS_' "{{env_file}}"

# Run the full E2E integration test
[working-directory: '..']
e2e:
    go test -tags integration -run TestE2E_FullFestivalCycle -v -timeout 5m ./internal/integration/

# Verify HCS topics exist and are readable
verify-topics:
    #!/usr/bin/env bash
    set -euo pipefail
    source "{{env_file}}"
    echo "Task topic:   ${HCS_TASK_TOPIC_ID:-NOT SET}"
    echo "Status topic: ${HCS_STATUS_TOPIC_ID:-NOT SET}"
    echo "Payment token: ${HTS_PAYMENT_TOKEN_ID:-NOT SET}"
    if [ -n "${HCS_TASK_TOPIC_ID:-}" ]; then
        echo "View on HashScan: https://hashscan.io/testnet/topic/${HCS_TASK_TOPIC_ID}"
        echo "View on HashScan: https://hashscan.io/testnet/topic/${HCS_STATUS_TOPIC_ID}"
        echo "View on HashScan: https://hashscan.io/testnet/token/${HTS_PAYMENT_TOKEN_ID}"
    fi

# Show all Hedera env vars from .env
show-config:
    @grep -E 'HEDERA_|HCS_|HTS_' "{{env_file}}" || echo "No Hedera config found in .env"
