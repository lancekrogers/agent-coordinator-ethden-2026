# Transaction Registry

All on-chain transactions from the live testnet integration test (2026-02-21).

## Hedera Testnet (Coordinator)

| Transaction Type | Description | Status | Notes |
|-----------------|-------------|--------|-------|
| HCS Task Assignment (x2) | Coordinator publishes `task-inference-01` and `task-defi-01` to topic 0.0.7999404 | Confirmed | Messages delivered to both agents |
| HCS Result Report (inference) | inference-001 reports task_result (failed) to topic 0.0.7999405 | Confirmed | Coordinator received at 05:04:15 |
| HCS Result Report (defi) | defi-001 reports task_result (completed) to topic 0.0.7999405 | Confirmed | Coordinator received at 05:04:34 |
| HTS Payment | 100 tokens transferred to defi-001 (account 0.0.7985425) via token 0.0.7999406 | Confirmed | Settled at 05:04:38 |

**Hedera Explorer:** All HCS messages and HTS transfers can be verified at [hashscan.io](https://hashscan.io/testnet/) using topic IDs `0.0.7999404` and `0.0.7999405`.

## 0G Chain (Inference Agent)

| Transaction Type | Description | Status | Notes |
|-----------------|-------------|--------|-------|
| Compute Job Submission | Inference job to 0G Compute | Blocked | No serving contract deployed |
| Storage Upload | Result uploaded to 0G Storage | Blocked | Depends on compute result |
| DA Audit Publish | Audit record to 0G DA | Blocked | Depends on compute result |
| iNFT Mint (ERC-7857) | Encrypted metadata on 0G Chain | Blocked | Depends on storage content |

**0G Chain RPC:** `https://evmrpc-testnet.0g.ai` (chain ID 16602) — connection verified.

## Base Sepolia (DeFi Agent)

| Transaction Type | Description | Status | Notes |
|-----------------|-------------|--------|-------|
| ERC-8004 Registration | Agent identity registered on-chain | Executed | Stub contract (no real contract address) |
| DEX Swap | Trade on Uniswap v3 | Executed | Stub router (zero address) |

**Base Sepolia RPC:** `https://sepolia.base.org` (chain ID 84532) — connection verified.

## Summary

- **4 Hedera testnet transactions confirmed** (2 HCS publishes, 2 HCS receives, 1 HTS payment)
- **0 0G Chain transactions** (blocked by missing serving contract)
- **2 Base Sepolia transactions** (executed against stub contracts)
