# Transaction Registry

All on-chain transactions required for bounty submission evidence. Hashes and explorer links will be populated after testnet deployment.

## Hedera (Coordinator)

| Chain | Transaction Type | Description | Hash | Explorer Link |
|-------|-----------------|-------------|------|---------------|
| Hedera Testnet | HCS Topic Creation | Create topic for agent communication | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| Hedera Testnet | HCS Task Assignment | Coordinator publishes task to inference agent | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| Hedera Testnet | HCS Result Report | Inference agent reports result to coordinator | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| Hedera Testnet | HTS Payment | Token transfer for agent compensation | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |

## 0G (Inference Agent)

| Chain | Transaction Type | Description | Hash | Explorer Link |
|-------|-----------------|-------------|------|---------------|
| 0G Testnet | Compute Job Submission | Inference job submitted to 0G Compute | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| 0G Testnet | Storage Upload | Inference result uploaded to 0G Storage | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| 0G Testnet | DA Audit Publish | Audit record published to 0G DA | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| 0G Chain (16602) | iNFT Mint | ERC-7857 iNFT minted with encrypted metadata | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |

## Base Sepolia (DeFi Agent)

| Chain | Transaction Type | Description | Hash | Explorer Link |
|-------|-----------------|-------------|------|---------------|
| Base Sepolia (84532) | ERC-8004 Registration | Agent identity registered on-chain | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |
| Base Sepolia (84532) | DEX Swap | Trade executed on Uniswap v3 | TBD - awaiting testnet deployment | TBD - awaiting testnet deployment |

## Notes

- All hashes will be populated during the submission-and-polish phase.
- Each hash should be verified on the respective chain explorer before final submission.
- Multiple transactions of the same type may be recorded; this table tracks the first confirmed instance of each.
