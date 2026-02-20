package integration

import (
	"fmt"
	"os"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// TestnetConfig holds all testnet account configuration.
type TestnetConfig struct {
	Network            string
	CoordinatorAccount AccountConfig
	Agent1Account      AccountConfig
	Agent2Account      AccountConfig
}

// AccountConfig holds a single Hedera account's configuration.
type AccountConfig struct {
	AccountID  hiero.AccountID
	PrivateKey hiero.PrivateKey
}

// LoadTestnetConfig loads configuration from environment variables.
func LoadTestnetConfig() (*TestnetConfig, error) {
	coordAcct, err := loadAccountConfig("HEDERA_COORDINATOR")
	if err != nil {
		return nil, fmt.Errorf("load coordinator config: %w", err)
	}

	agent1Acct, err := loadAccountConfig("HEDERA_AGENT1")
	if err != nil {
		return nil, fmt.Errorf("load agent1 config: %w", err)
	}

	agent2Acct, err := loadAccountConfig("HEDERA_AGENT2")
	if err != nil {
		return nil, fmt.Errorf("load agent2 config: %w", err)
	}

	return &TestnetConfig{
		Network:            envOrDefault("HEDERA_NETWORK", "testnet"),
		CoordinatorAccount: *coordAcct,
		Agent1Account:      *agent1Acct,
		Agent2Account:      *agent2Acct,
	}, nil
}

func loadAccountConfig(prefix string) (*AccountConfig, error) {
	accountIDStr := os.Getenv(prefix + "_ACCOUNT_ID")
	if accountIDStr == "" {
		return nil, fmt.Errorf("%s_ACCOUNT_ID not set", prefix)
	}

	accountID, err := hiero.AccountIDFromString(accountIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse %s_ACCOUNT_ID %q: %w", prefix, accountIDStr, err)
	}

	privateKeyStr := os.Getenv(prefix + "_PRIVATE_KEY")
	if privateKeyStr == "" {
		return nil, fmt.Errorf("%s_PRIVATE_KEY not set", prefix)
	}

	privateKey, err := hiero.PrivateKeyFromString(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("parse %s_PRIVATE_KEY: %w", prefix, err)
	}

	return &AccountConfig{
		AccountID:  accountID,
		PrivateKey: privateKey,
	}, nil
}

func envOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// NewClientForAccount creates a Hedera client configured for a specific account on testnet.
func NewClientForAccount(acct AccountConfig) (*hiero.Client, error) {
	client := hiero.ClientForTestnet()
	client.SetOperator(acct.AccountID, acct.PrivateKey)
	return client, nil
}
