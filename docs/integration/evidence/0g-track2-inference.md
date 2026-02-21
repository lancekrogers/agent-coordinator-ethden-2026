# 0G Track 2: Decentralized GPU Inference

**Bounty**: $7,000 -- Decentralized GPU Inference

## Inference Agent Architecture

The inference agent communicates with three 0G subsystems through dedicated clients:

| Component | Responsibility | Code Path |
|-----------|---------------|-----------|
| ComputeBroker | Submits inference jobs via OpenAI-compatible REST API to 0G Compute | `internal/zerog/compute/broker.go` |
| StorageClient | Uploads inference results in chunks to 0G Storage | `internal/zerog/storage/client.go` |
| AuditPublisher | Publishes audit records to 0G DA with retry and exponential backoff | `internal/zerog/da/publisher.go` |

## Pipeline

The inference agent operates within the three-agent cycle as follows:

1. Coordinator assigns a task via Hedera HCS message.
2. Inference agent receives the task and submits an inference job to **0G Compute** (OpenAI-compatible REST endpoint).
3. Agent polls 0G Compute for job completion.
4. Completed result is uploaded to **0G Storage** using chunked upload protocol.
5. An audit record (job ID, storage reference, timestamp) is published to **0G DA** for data availability verification.
6. Agent reports the result back to the coordinator via HCS.

## Implementation Details

- **ComputeBroker**: Wraps the 0G Compute Network's OpenAI-compatible API. Handles job submission, status polling, and result retrieval. Configurable endpoint, model selection, and timeout.
- **StorageClient**: Implements chunked file upload to 0G Storage. Handles large inference outputs by splitting into segments and reassembling references.
- **AuditPublisher**: Publishes structured audit data to 0G's Data Availability layer. Includes retry logic with exponential backoff to handle transient network errors.

## Test Coverage

- **Coverage**: 75.2%
- **Total tests**: 63
- Tests validate: job submission formatting, polling state machine, chunk boundary handling, DA publish retries, error propagation, context cancellation.

## Transaction Evidence

| Transaction Type | Hash | Status |
|-----------------|------|--------|
| 0G Compute job submission | TBD | Awaiting testnet deployment |
| 0G Storage upload | TBD | Awaiting testnet deployment |
| 0G DA audit publish | TBD | Awaiting testnet deployment |

## Live Testnet Results (2026-02-21)

The inference agent was deployed against live Hedera testnet and 0G Chain testnet:

- **HCS transport:** Working — agent received `task-inference-01` from coordinator via HCS topic `0.0.7999404`
- **0G Chain RPC:** Connected successfully to `https://evmrpc-testnet.0g.ai`
- **Task processing:** Agent correctly unpacked the task payload (model=test-model, input="Analyze market sentiment for ETH")
- **Failure reported:** Agent reported failure status back to coordinator via HCS topic `0.0.7999405`

**Blocker:** 0G Compute execution failed because `ZG_SERVING_CONTRACT` is empty — no serving contract is deployed on the 0G testnet. The agent's on-chain query returns "no contract code at given address".

## Evidence Gap

To complete evidence for this track:

1. Deploy or locate an existing 0G Serving contract on testnet with registered providers
2. Re-run the three-agent cycle with the serving contract configured
3. Capture compute job submission, storage upload, DA publish, and iNFT mint transaction hashes
4. Record explorer links for each transaction type
