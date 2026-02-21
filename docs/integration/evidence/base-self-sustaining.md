# Base Bounty: Self-Sustaining Agent

**Bounty**: $3,000+ -- Self-Sustaining Agent on Base

## Overview

The DeFi agent operates on Base Sepolia (chain ID 84532) as a self-sustaining autonomous trading agent. It registers an on-chain identity, executes trades, tracks profit and loss, and attributes work to its builder -- all using Base-native standards.

## Components

| Component | Responsibility | Code Path |
|-----------|---------------|-----------|
| Identity (ERC-8004) | Registers agent identity on-chain | `internal/base/identity/` |
| Payment (x402) | Machine-to-machine payment via HTTP 402 handshake | `internal/base/payment/` |
| Attribution (ERC-8021) | Appends builder attribution to transaction calldata | `internal/base/attribution/` |
| Trading | Evaluates market state and executes swaps | `internal/base/trading/` |

## ERC-8004: Agent Identity

The agent registers itself on Base Sepolia using the ERC-8004 standard. This provides a verifiable on-chain identity that other agents and protocols can reference.

## x402: Payment Protocol

The agent implements the x402 payment protocol for machine-to-machine payments:

1. Agent sends a request to a paywall-protected resource.
2. Server responds with HTTP 402 and payment requirements.
3. Agent constructs and signs a payment satisfying the requirements.
4. Agent resubmits the request with the payment header.
5. Server verifies payment and grants access.

This enables the agent to autonomously pay for services (data feeds, compute, API access) without human intervention.

## ERC-8021: Builder Attribution

Every transaction the agent submits includes ERC-8021 builder attribution:

- **Format**: 4-byte magic prefix + 20-byte builder address appended to calldata
- **Purpose**: Attributes on-chain activity to the builder/developer for rewards and reputation

## Trading

| Sub-component | Description |
|--------------|-------------|
| MeanReversionStrategy | Evaluates current market state against historical mean; generates buy/sell signals when price deviates beyond threshold |
| TradeExecutor | Submits swap transactions to Uniswap v3 on Base Sepolia |

## P&L Tracking

The `PnLTracker` maintains a thread-safe record of all agent economic activity:

- **Tracks**: Trade outcomes, gas costs, protocol fees, x402 payments
- **Key metric**: `NetPnL` -- sum of all revenue minus all costs
- **Self-sustaining flag**: `IsSelfSustaining` returns true when `NetPnL > 0`

The agent reports its P&L back to the coordinator as part of the three-agent cycle.

## Test Coverage

- **Coverage**: 76.6%
- **Total tests**: 57
- Tests validate: identity registration flow, x402 handshake, attribution encoding, strategy signal generation, trade execution, P&L arithmetic, thread safety.

## Transaction Evidence

| Transaction Type | Hash | Status |
|-----------------|------|--------|
| ERC-8004 identity registration | TBD | Awaiting testnet deployment |
| DEX swap (Uniswap v3) | TBD | Awaiting testnet deployment |
| x402 payment | TBD | Awaiting testnet deployment |

## Evidence Gap

No transactions have been executed on Base Sepolia. To complete evidence:

1. Obtain Base Sepolia ETH from faucet.
2. Deploy agent with Base Sepolia RPC configuration.
3. Execute identity registration, at least one trade, and one x402 payment.
4. Capture transaction hashes and demonstrate `IsSelfSustaining == true`.
