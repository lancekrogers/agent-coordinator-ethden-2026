package hts

import (
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TokenConfig holds configuration for creating a new fungible token.
type TokenConfig struct {
	// Name is the human-readable token name (e.g., "Agent Payment Token").
	Name string

	// Symbol is the short token symbol (e.g., "APT").
	Symbol string

	// Decimals is the number of decimal places (0 for whole tokens only).
	Decimals uint32

	// InitialSupply is the number of tokens to mint at creation.
	InitialSupply uint64

	// TreasuryAccountID is the account that holds the initial supply.
	TreasuryAccountID hiero.AccountID

	// AdminKey can modify the token. If nil, the token is immutable.
	AdminKey *hiero.PublicKey

	// SupplyKey can mint/burn tokens. If nil, supply is fixed.
	SupplyKey *hiero.PublicKey
}

// DefaultTokenConfig returns sensible defaults for the agent payment use case.
// The caller must still set TreasuryAccountID.
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		Name:          "Agent Payment Token",
		Symbol:        "APT",
		Decimals:      0,
		InitialSupply: 1000000,
	}
}

// TokenMetadata holds information about an existing HTS token.
type TokenMetadata struct {
	TokenID     hiero.TokenID
	Name        string
	Symbol      string
	Decimals    uint32
	TotalSupply uint64
	TreasuryID  hiero.AccountID
}

// TransferRequest specifies a token transfer between two accounts.
type TransferRequest struct {
	// TokenID is the token to transfer.
	TokenID hiero.TokenID

	// FromAccountID is the sender account.
	FromAccountID hiero.AccountID

	// ToAccountID is the recipient account.
	ToAccountID hiero.AccountID

	// Amount is the number of tokens to transfer (signed: negative = debit).
	Amount int64

	// Memo is an optional memo attached to the transfer transaction.
	Memo string
}

// TransferReceipt holds the result of a completed token transfer.
type TransferReceipt struct {
	// TransactionID is the Hedera transaction ID for this transfer.
	TransactionID hiero.TransactionID

	// TokenID is the token that was transferred.
	TokenID hiero.TokenID

	// FromAccountID is the sender.
	FromAccountID hiero.AccountID

	// ToAccountID is the recipient.
	ToAccountID hiero.AccountID

	// Amount is the number of tokens transferred.
	Amount int64

	// Status is the transaction status from the receipt.
	Status string
}
