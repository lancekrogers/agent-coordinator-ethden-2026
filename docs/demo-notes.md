# Demo Notes: Hedera Track 3

Demo script for the three-agent autonomous economy. Target audience: Hedera Track 3 judges evaluating native (non-EVM) Hedera usage.

## Key Message

All agent coordination and payment settlement uses native Hedera services (HCS + HTS) through the Go SDK. Zero Solidity. Zero EVM contracts. The entire multi-agent economy runs on Hedera's native layer.

## Demo Structure

### 1. Opening (10s)

**Say:** "This is a three-agent autonomous economy where AI agents coordinate, execute tasks, and get paid -- all on Hedera's native services."

**Show:** Dashboard overview with the three agents visible (coordinator, inference, DeFi).

### 2. HCS Messaging (20s)

**Say:**
- "Every message between agents flows through Hedera Consensus Service"
- "We have two HCS topics: one for task assignments, one for status updates"
- "Every message is immutable, timestamped, and publicly verifiable on the Hedera network"

**Show:**
- Dashboard HCS feed panel showing live messages
- Point out the topic IDs: Task Topic `0.0.7999404`, Status Topic `0.0.7999405`
- Click a message to show the JSON payload structure

**Visual cue:** Highlight the `task_assignment` message type flowing from coordinator to agents.

### 3. Task Lifecycle (20s)

**Say:**
- "Watch a task go from assignment to completion entirely through HCS"
- "The coordinator publishes a task assignment -- the inference agent picks it up, runs inference on 0G Compute, and publishes the result back"
- "Every state transition is a new HCS message with a sequence number, so we have a full audit trail"

**Show:**
- Trigger a task assignment (or show a recent one)
- Follow the message chain: `task_assignment` -> `task_result` -> `pnl_report`
- Point out the sequence numbers incrementing
- Show the task state: `assigned` -> `in_progress` -> `complete`

**Visual cue:** Dashboard agent activity panel showing status changes in real-time.

### 4. HTS Payment (15s)

**Say:**
- "When a task completes, the coordinator automatically pays the agent using Hedera Token Service"
- "This is a native HTS transfer -- not a smart contract call"
- "The payment confirmation is then published back to HCS, creating a complete audit loop"

**Show:**
- Payment settlement message in the HCS feed
- Open HashScan to show the HTS transfer transaction
- Point out the AGNT token transfer from coordinator treasury to agent account

**Visual cue:** HashScan transaction page showing the token transfer.

### 5. Key Differentiator (10s)

**Say:**
- "The entire system -- messaging, coordination, and payments -- runs on Hedera's native layer"
- "No Solidity, no EVM contracts. Just the Go SDK talking directly to HCS and HTS"
- "This proves that complex multi-agent workflows don't need smart contracts when you have native Hedera services"

**Show:** The project structure showing `hiero-sdk-go` imports and the absence of any `.sol` files.

## Hedera-Specific Highlights

When judges ask questions, emphasize:

- **HCS as a message bus**: Replaces centralized message queues. Every agent message is a permanent, ordered record.
- **HTS for settlements**: Native token transfers that finalize in seconds. No gas wars, no contract deployment, no approval patterns.
- **Topic architecture**: Two-topic design (command vs. status) separates concerns. Coordinator publishes to task topic, agents publish to status topic. Clean routing.
- **Go SDK**: `hiero-sdk-go v2.75.0` for all Hedera interactions. Production-grade SDK with full type safety.
- **Mirror Node for observability**: Dashboard reads from the Hedera Mirror Node API, showing the same data that any third party can verify.

## HashScan Links

Keep these tabs open during the demo:

- Task Topic: `https://hashscan.io/testnet/topic/0.0.7999404`
- Status Topic: `https://hashscan.io/testnet/topic/0.0.7999405`
- Coordinator Account: `https://hashscan.io/testnet/account/{COORDINATOR_ACCOUNT_ID}`

## Fallback Plan

If the live testnet is slow or unresponsive during the demo:

1. Show the integration test running locally: `go test -v -run TestLocalThreeAgentCycle ./internal/integration/`
2. Show pre-recorded HCS messages on HashScan
3. Walk through the code: HCS publisher, HTS transfer, result handler
