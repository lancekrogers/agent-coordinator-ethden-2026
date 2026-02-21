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

## Live Testnet Results (2026-02-21)

The inference agent connected to 0G Chain RPC (`https://evmrpc-testnet.0g.ai`) successfully. The iNFT minting step is downstream of the compute step, which is currently blocked (no serving contract deployed). Once compute results are available, the minting pipeline is ready to execute.

## Evidence Gap

No minting transactions have been executed on 0G Chain testnet. To complete evidence:

1. Resolve the 0G Compute blocker (deploy serving contract with providers)
2. Once inference results are available, the minting pipeline will execute automatically
3. Capture the iNFT mint transaction hash on 0G Chain (chain ID 16602)
4. Verify encrypted metadata via the explorer
