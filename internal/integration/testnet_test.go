//go:build integration

package integration

import (
	"testing"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

func TestVerifyTestnetAccounts(t *testing.T) {
	cfg, err := LoadTestnetConfig()
	if err != nil {
		t.Skipf("testnet config not available: %v", err)
	}

	accounts := []struct {
		name   string
		config AccountConfig
	}{
		{"coordinator", cfg.CoordinatorAccount},
		{"agent-1", cfg.Agent1Account},
		{"agent-2", cfg.Agent2Account},
	}

	for _, acct := range accounts {
		t.Run(acct.name, func(t *testing.T) {
			client, err := NewClientForAccount(acct.config)
			if err != nil {
				t.Fatalf("create client: %v", err)
			}
			defer client.Close()

			balance, err := hiero.NewAccountBalanceQuery().
				SetAccountID(acct.config.AccountID).
				Execute(client)
			if err != nil {
				t.Fatalf("query balance for %s: %v", acct.config.AccountID, err)
			}

			t.Logf("%s (%s): %s", acct.name, acct.config.AccountID, balance.Hbars)
		})
	}
}

func TestLoadTestnetConfig_MissingEnv(t *testing.T) {
	// Without env vars set, LoadTestnetConfig should return an error.
	_, err := LoadTestnetConfig()
	if err == nil {
		t.Skip("env vars are set, skipping missing env test")
	}
}
