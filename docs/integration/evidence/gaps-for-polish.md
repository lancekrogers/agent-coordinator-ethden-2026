# Gaps for Submission-and-Polish

Updated after live testnet integration test (2026-02-21). Items marked with checkmarks are resolved.

## Resolved

- [x] Hedera testnet accounts configured (coordinator, inference, defi)
- [x] HCS topics created (task: 0.0.7999404, status: 0.0.7999405)
- [x] HTS payment token created (0.0.7999406)
- [x] HCS transport implemented in inference and defi agents (Hiero SDK)
- [x] Coordinator main.go fully wired with config, assigner, monitor, payment, result handler
- [x] Three-agent cycle verified on live Hedera testnet
- [x] Coordinator → HCS → Agent → HCS → Coordinator → HTS payment flow confirmed
- [x] Graceful shutdown verified for all three agents

## 1. 0G Serving Contract Deployment (Critical)

- Deploy 0G Serving contract on testnet with at least one registered provider
- Or locate an existing deployment and configure `ZG_SERVING_CONTRACT`
- Without this, inference agent cannot complete the compute → storage → iNFT → DA pipeline
- **Blocks:** 0G Track 2 ($7k) and Track 3 ($7k) evidence

## 2. Base Sepolia Contract Configuration (Critical)

- Configure real DEX router address for Uniswap V3 on Base Sepolia (`DEFI_DEX_ROUTER`)
- Deploy or locate ERC-8004 identity contract (`DEFI_ERC8004_CONTRACT`)
- Register ERC-8021 builder code (`DEFI_BUILDER_CODE`)
- **Blocks:** Base bounty ($3k+) evidence

## 3. Transaction Capture

- Re-run three-agent cycle with all contracts configured
- Record all transaction hashes in `transaction-registry.md`
- Verify each hash on the respective chain explorer
- Capture at least one successful iNFT mint on 0G Chain
- Capture at least one successful DEX swap on Base Sepolia
- Demonstrate `IsSelfSustaining == true` in DeFi agent P&L output

## 4. Architecture Diagram

- Create a visual diagram showing:
  - Three agents and their roles
  - Communication flow via Hedera HCS
  - 0G subsystem interactions (compute, storage, DA, chain)
  - Base Sepolia interactions (identity, trading, payments, attribution)
- Format: PNG or SVG suitable for README and submission page

## 5. Demo Video

- Record screencast of the three-agent cycle running against live testnets
- Show: task assignment, inference execution, iNFT minting, trade execution, P&L report
- Duration target: 2-3 minutes
- Include explorer verification of at least one transaction per chain

## 6. README Updates

- `agent-coordinator/README.md` — architecture overview, setup instructions, bounty references
- `agent-inference/README.md` — 0G integration details, configuration, test instructions
- `agent-defi/README.md` — Base integration details, trading strategy, P&L tracking
- Each README should link to relevant bounty track documentation

## 7. Bounty Submission Forms

- 0G Track 2 (Decentralized GPU Inference) — $7k
- 0G Track 3 (ERC-7857 iNFT) — $7k
- Base (Self-Sustaining Agent) — $3k+
- Each submission requires: project description, repo links, transaction evidence, demo video link
