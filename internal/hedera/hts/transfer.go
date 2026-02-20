package hts

import (
	"context"
	"fmt"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TransferService implements the TokenTransfer interface using the Hiero (Hedera) SDK.
type TransferService struct {
	client *hiero.Client
}

// NewTransferService creates a new TransferService with the given Hiero client.
func NewTransferService(client *hiero.Client) *TransferService {
	return &TransferService{client: client}
}

// Transfer moves tokens from one account to another atomically.
func (s *TransferService) Transfer(ctx context.Context, req TransferRequest) (*TransferReceipt, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("transfer token %s from %s to %s: %w",
			req.TokenID, req.FromAccountID, req.ToAccountID, err)
	}

	if req.Amount <= 0 {
		return nil, fmt.Errorf("transfer token %s: amount must be positive, got %d",
			req.TokenID, req.Amount)
	}

	tx := hiero.NewTransferTransaction().
		AddTokenTransfer(req.TokenID, req.FromAccountID, -req.Amount).
		AddTokenTransfer(req.TokenID, req.ToAccountID, req.Amount)

	if req.Memo != "" {
		tx.SetTransactionMemo(req.Memo)
	}

	frozen, err := tx.FreezeWith(s.client)
	if err != nil {
		return nil, fmt.Errorf("transfer %d of token %s from %s to %s: freeze: %w",
			req.Amount, req.TokenID, req.FromAccountID, req.ToAccountID, err)
	}

	resp, err := frozen.Execute(s.client)
	if err != nil {
		return nil, fmt.Errorf("transfer %d of token %s from %s to %s: execute: %w",
			req.Amount, req.TokenID, req.FromAccountID, req.ToAccountID, err)
	}

	receipt, err := resp.GetReceipt(s.client)
	if err != nil {
		return nil, fmt.Errorf("transfer %d of token %s from %s to %s: receipt: %w",
			req.Amount, req.TokenID, req.FromAccountID, req.ToAccountID, err)
	}

	return &TransferReceipt{
		TransactionID: resp.TransactionID,
		TokenID:       req.TokenID,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Status:        receipt.Status.String(),
	}, nil
}

// AssociateToken associates a token with an account so it can receive transfers.
func (s *TransferService) AssociateToken(ctx context.Context, tokenID hiero.TokenID, accountID hiero.AccountID) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("associate token %s with account %s: %w", tokenID, accountID, err)
	}

	tx, err := hiero.NewTokenAssociateTransaction().
		SetAccountID(accountID).
		SetTokenIDs(tokenID).
		FreezeWith(s.client)
	if err != nil {
		return fmt.Errorf("associate token %s with account %s: freeze: %w", tokenID, accountID, err)
	}

	resp, err := tx.Execute(s.client)
	if err != nil {
		return fmt.Errorf("associate token %s with account %s: execute: %w", tokenID, accountID, err)
	}

	_, err = resp.GetReceipt(s.client)
	if err != nil {
		return fmt.Errorf("associate token %s with account %s: receipt: %w", tokenID, accountID, err)
	}

	return nil
}

// Compile-time interface compliance check.
var _ TokenTransfer = (*TransferService)(nil)
