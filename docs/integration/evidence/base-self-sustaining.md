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

## Live Testnet Results (2026-02-21)

The DeFi agent was deployed against live Hedera testnet and Base Sepolia:

- **HCS transport:** Working — agent received `task-defi-01` from coordinator via HCS topic `0.0.7999404`
- **Base Sepolia RPC:** Connected successfully to `https://sepolia.base.org`
- **Identity:** Registration completed (ERC-8004 on Base Sepolia)
- **Trading:** 2 trades executed during 117s test run (60s interval)
- **Strategy signals:** Mean reversion sell signal at 71.4% confidence
- **Result reported:** Completed status sent to coordinator via HCS (duration=346ms)
- **Payment received:** Coordinator triggered 100 HTS token payment to defi-001

**Partial blockers:**
- DEX Router is set to zero address — trades execute against stub
- Builder code not configured — ERC-8021 attribution disabled
- x402 payment not demonstrated — no paywall-protected resource configured

## Evidence Gap

To complete evidence for this track:

1. Deploy Uniswap V3 router on Base Sepolia or configure existing deployment address
2. Register an ERC-8021 builder code
3. Configure a paywall-protected resource for x402 demonstration
4. Re-run agent and capture real DEX swap transaction hashes
5. Demonstrate `IsSelfSustaining == true` with real P&L data
