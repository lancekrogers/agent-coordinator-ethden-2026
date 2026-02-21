# agent-coordinator

Task orchestration agent for the multi-agent economy on Hedera.

Part of the [ETHDenver 2026 Agent Economy](../README.md) submission.

## Overview

Reads festival plans, assigns tasks to agents via Hedera Consensus Service (HCS), monitors progress, enforces quality gates, and manages HTS token payments upon task completion. Includes a testnet setup utility that provisions HCS topics and the HTS payment token.

## Built with Obedience Corp

This project is part of an [Obedience Corp](https://obediencecorp.com) campaign — built and planned using **camp** (campaign management) and **fest** (festival methodology). This repository, its git history, and the planning artifacts in `festivals/` are a live example of these tools in action.

The coordinator connects to the **obey daemon** via gRPC for sandboxed command execution. Festival plans define the task graph (phases -> sequences -> tasks) that the coordinator publishes to agents via HCS. Quality gates enforce `fest_commit` checkpoints before allowing task completion.

## System Context

```
                    ┌─────────────┐
           tasks    │ Coordinator │    tasks        <-- you are here
          ┌────────>│  (Hedera)   │<────────┐
          │         └─────────────┘         │
          │               │                 │
          │          assignments             │
          │               │                 │
    ┌─────┴─────┐         │         ┌───────┴──────┐
    │ Inference │         │         │  DeFi Agent  │
    │   (0G)    │         └────────>│   (Base)     │
    └───────────┘                   └──────────────┘
```

## Quick Start

```bash
cp .env.example .env   # fill in Hedera accounts + topic/token IDs
just build
just run
```

To provision testnet infrastructure from scratch:

```bash
just hedera setup          # creates HCS topics + HTS token
just hedera show-config    # prints the generated .env values
```

## Prerequisites

- Go 1.24+
- 3 Hedera testnet accounts -- coordinator (treasury) + 2 agents ([portal.hedera.com](https://portal.hedera.com))
- HCS topics and HTS token (created via `just hedera setup`)

## Configuration

| Variable | Description |
|----------|-------------|
| `HEDERA_NETWORK` | Network name (`testnet`) |
| `HEDERA_COORDINATOR_ACCOUNT_ID` | Coordinator / treasury account |
| `HEDERA_COORDINATOR_PRIVATE_KEY` | Coordinator private key |
| `HEDERA_AGENT1_ACCOUNT_ID` | Inference agent account |
| `HEDERA_AGENT1_PRIVATE_KEY` | Inference agent key |
| `HEDERA_AGENT2_ACCOUNT_ID` | DeFi agent account |
| `HEDERA_AGENT2_PRIVATE_KEY` | DeFi agent key |
| `HCS_TASK_TOPIC_ID` | HCS topic for task assignments |
| `HCS_STATUS_TOPIC_ID` | HCS topic for status updates |
| `HTS_PAYMENT_TOKEN_ID` | HTS fungible token for payments |
| `DAEMON_ADDRESS` | Daemon gRPC address (default: localhost:50051) |
| `DAEMON_TLS_ENABLED` | Enable TLS for daemon connection |

## HCS Message Protocol

| Type | Direction | Description |
|------|-----------|-------------|
| `task_assignment` | Coordinator -> Agent | Assigns a task with parameters |
| `status_update` | Agent -> Coordinator | Reports task progress |
| `task_result` | Agent -> Coordinator | Delivers task output |
| `pnl_report` | DeFi Agent -> Coordinator | P&L metrics |
| `heartbeat` | Agent -> Coordinator | Liveness signal |
| `quality_gate` | Coordinator -> Agent | Quality check enforcement |
| `payment_settled` | Coordinator -> Agent | HTS payment confirmation |

Task state machine: `pending` -> `assigned` -> `in_progress` -> `review` -> `complete` -> `paid`

## Project Structure

```
cmd/
  coordinator/             Coordinator entry point
  setup-testnet/           Provisions HCS topics + HTS token
internal/
  config/                  Config loading and validation
  coordinator/             Assigner, monitor, payment, result handler, quality gates
  daemon/                  Daemon RPC client
  festival/                Festival plan reader
  hedera/
    hcs/                   HCS publisher, subscriber, topic lifecycle
    hts/                   HTS token creation and transfer
  integration/             E2E integration test helpers
pkg/daemon/                Shared daemon proto bindings
proto/                     Protobuf definitions
docs/integration/          Integration test evidence and logs
```

## Development

```bash
just build                 # Build binary to bin/
just run                   # Run the coordinator
just test                  # Run tests
just lint                  # golangci-lint
just hedera setup          # Provision HCS topics + HTS token
just hedera e2e            # Full E2E integration test
just hedera verify-topics  # Check topic/token existence
just hedera show-config    # Display Hedera env vars
```

## Architecture

`main.go` initializes the Hedera client, config, and all coordinator services via dependency injection. The Assigner publishes task assignments to agents via HCS. The Monitor and ResultHandler run as background goroutines, listening on the status topic for agent updates. Quality gates validate task completion before the Payment service executes HTS transfers and publishes settlement confirmations. All state management is thread-safe with `sync.RWMutex`.

## License

Apache-2.0
