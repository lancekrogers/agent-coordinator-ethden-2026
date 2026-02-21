# 0G Track 3: ERC-7857 iNFT

**Bounty**: $7,000 -- ERC-7857 Intelligent NFT

## Overview

The inference agent mints ERC-7857 intelligent NFTs on 0G Chain after each successful inference run. Each iNFT encapsulates the inference output as encrypted metadata, creating a verifiable on-chain record of AI-generated results.

## Architecture

| Component | Responsibility | Code Path |
|-----------|---------------|-----------|
| Minter | Constructs and submits ERC-7857 mint transactions via JSON-RPC | `internal/zerog/inft/minter.go` |
| Encrypt | Encrypts inference metadata using AES-256-GCM before embedding in iNFT | `internal/zerog/inft/encrypt.go` |

## Target Chain

- **Chain**: 0G Chain
- **Chain ID**: 16602
- **Interaction**: JSON-RPC (eth_sendRawTransaction, eth_getTransactionReceipt)

## Encryption

Inference results are encrypted before being stored as iNFT metadata:

- **Algorithm**: AES-256-GCM
- **Key validation**: 32-byte key length enforced at initialization
- **Nonce**: Random nonce generated per encryption operation (crypto/rand)
- **Output format**: Nonce prepended to ciphertext

This ensures inference outputs remain confidential on-chain while still being verifiable by authorized parties holding the decryption key.

## Minting Flow

1. Inference agent completes a compute job and receives the result.
2. Result metadata (job ID, model, output hash, timestamp) is serialized to JSON.
3. JSON payload is encrypted using AES-256-GCM.
4. Minter constructs an ERC-7857 mint transaction with the encrypted metadata.
5. Transaction is signed and submitted via JSON-RPC to 0G Chain.
6. Minter polls for transaction receipt to confirm successful minting.

## Transaction Evidence

| Transaction Type | Hash | Status |
|-----------------|------|--------|
| ERC-7857 iNFT mint | TBD | Awaiting testnet deployment |

## Evidence Gap

No minting transactions have been executed on 0G Chain testnet. To complete evidence:

1. Obtain 0G Chain testnet tokens.
2. Deploy or reference an ERC-7857 contract on chain ID 16602.
3. Execute a mint transaction with encrypted inference metadata.
4. Capture the transaction hash and explorer link.
