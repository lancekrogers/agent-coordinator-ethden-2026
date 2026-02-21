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

## Evidence Gap

Testnet deployment has not been performed. To complete evidence for this track:

1. Obtain 0G testnet tokens from faucet.
2. Deploy inference agent against 0G testnet endpoints.
3. Execute a full inference cycle and capture transaction hashes.
4. Record explorer links for each transaction type.
