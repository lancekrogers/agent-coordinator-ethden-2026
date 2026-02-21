# Gaps for Submission-and-Polish

Ordered list of items the submission-and-polish festival phase must address before final hackathon submission.

## 1. Testnet Credential Setup

- Hedera testnet account (portal.hedera.com) -- operator ID and private key
- 0G testnet tokens from faucet -- for compute, storage, DA, and chain transactions
- Base Sepolia ETH from faucet -- for identity registration, trading, and gas

## 2. Live Agent Deployment

- Deploy coordinator against Hedera testnet with HCS topic creation
- Deploy inference agent against 0G testnet endpoints (compute, storage, DA, chain)
- Deploy DeFi agent against Base Sepolia RPC
- Verify all three agents can communicate end-to-end over live networks

## 3. Transaction Capture

- Execute the full three-agent cycle on testnets
- Record all transaction hashes in `transaction-registry.md`
- Verify each hash resolves on the respective chain explorer
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

- `agent-coordinator/README.md` -- architecture overview, setup instructions, bounty references
- `agent-inference/README.md` -- 0G integration details, configuration, test instructions
- `agent-defi/README.md` -- Base integration details, trading strategy, P&L tracking
- Each README should link to relevant bounty track documentation

## 7. Bounty Submission Forms

- 0G Track 2 (Decentralized GPU Inference) -- $7k
- 0G Track 3 (ERC-7857 iNFT) -- $7k
- Base (Self-Sustaining Agent) -- $3k+
- Each submission requires: project description, repo links, transaction evidence, demo video link
