package hts

import (
	"context"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TokenCreator handles HTS token lifecycle operations.
// Used by the coordinator to create the payment token for the agent flow.
type TokenCreator interface {
	// CreateFungibleToken creates a new fungible token on Hedera and returns its ID.
	CreateFungibleToken(ctx context.Context, config TokenConfig) (hiero.TokenID, error)

	// TokenInfo retrieves metadata about an existing token.
	TokenInfo(ctx context.Context, tokenID hiero.TokenID) (*TokenMetadata, error)
}

// TokenTransfer handles HTS token transfer operations between accounts.
// Used by the coordinator to pay agents for completed tasks.
type TokenTransfer interface {
	// Transfer moves tokens from one account to another.
	Transfer(ctx context.Context, req TransferRequest) (*TransferReceipt, error)

	// AssociateToken associates a token with an account so it can receive transfers.
	AssociateToken(ctx context.Context, tokenID hiero.TokenID, accountID hiero.AccountID) error
}
