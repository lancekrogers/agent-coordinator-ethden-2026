package hts

import (
	"context"
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TokenService implements the TokenCreator interface using the Hiero (Hedera) SDK.
type TokenService struct {
	client *hiero.Client
}

// NewTokenService creates a new TokenService with the given Hiero client.
func NewTokenService(client *hiero.Client) *TokenService {
	return &TokenService{client: client}
}

// CreateFungibleToken creates a new fungible token on Hedera and returns its ID.
func (s *TokenService) CreateFungibleToken(ctx context.Context, config TokenConfig) (hiero.TokenID, error) {
	if err := ctx.Err(); err != nil {
		return hiero.TokenID{}, fmt.Errorf("create token %q: %w", config.Name, err)
	}

	tx := hiero.NewTokenCreateTransaction().
		SetTokenName(config.Name).
		SetTokenSymbol(config.Symbol).
		SetDecimals(uint(config.Decimals)).
		SetInitialSupply(config.InitialSupply).
		SetTreasuryAccountID(config.TreasuryAccountID).
		SetTokenType(hiero.TokenTypeFungibleCommon)

	if config.AdminKey != nil {
		tx = tx.SetAdminKey(*config.AdminKey)
	}
	if config.SupplyKey != nil {
		tx = tx.SetSupplyKey(*config.SupplyKey)
	}

	frozen, err := tx.FreezeWith(s.client)
	if err != nil {
		return hiero.TokenID{}, fmt.Errorf("create token %q: freeze: %w", config.Name, err)
	}

	resp, err := frozen.Execute(s.client)
	if err != nil {
		return hiero.TokenID{}, fmt.Errorf("create token %q: execute: %w", config.Name, err)
	}

	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return hiero.TokenID{}, fmt.Errorf("create token %q: receipt: %w", config.Name, err)
	}

	if receipt.TokenID == nil {
		return hiero.TokenID{}, fmt.Errorf("create token %q: receipt contained nil token ID", config.Name)
	}

	return *receipt.TokenID, nil
}

// TokenInfo retrieves metadata about an existing token.
func (s *TokenService) TokenInfo(ctx context.Context, tokenID hiero.TokenID) (*TokenMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("token info %s: %w", tokenID, err)
	}

	info, err := hiero.NewTokenInfoQuery().
		SetTokenID(tokenID).
		Execute(s.client)
	if err != nil {
		return nil, fmt.Errorf("token info %s: execute query: %w", tokenID, err)
	}

	return &TokenMetadata{
		TokenID:     tokenID,
		Name:        info.Name,
		Symbol:      info.Symbol,
		Decimals:    info.Decimals,
		TotalSupply: info.TotalSupply,
		TreasuryID:  info.Treasury,
	}, nil
}

// Compile-time interface compliance check.
var _ TokenCreator = (*TokenService)(nil)
