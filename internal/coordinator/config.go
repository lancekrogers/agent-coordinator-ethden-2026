package coordinator

import (
	"fmt"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

// Config holds configuration for the coordinator engine.
type Config struct {
	// TaskTopicID is the HCS topic for task assignments and results.
	TaskTopicID hiero.TopicID

	// StatusTopicID is the HCS topic for agent status updates.
	StatusTopicID hiero.TopicID

	// PaymentTokenID is the HTS token used for agent payments.
	PaymentTokenID hiero.TokenID

	// TreasuryAccountID holds the payment token supply.
	TreasuryAccountID hiero.AccountID

	// DefaultPaymentAmount is the default token amount paid per task completion.
	DefaultPaymentAmount int64

	// MonitorPollInterval is how often the progress monitor checks for updates.
	MonitorPollInterval time.Duration

	// QualityGateTimeout is the max time to wait for quality gate evaluation.
	QualityGateTimeout time.Duration
}

// DefaultConfig returns sensible defaults for testnet usage.
func DefaultConfig() Config {
	return Config{
		DefaultPaymentAmount: 100,
		MonitorPollInterval:  5 * time.Second,
		QualityGateTimeout:   30 * time.Second,
	}
}

// Validate checks the config for required fields.
func (c Config) Validate() error {
	if c.TaskTopicID.Topic == 0 {
		return fmt.Errorf("coordinator config: task topic ID is required")
	}
	if c.StatusTopicID.Topic == 0 {
		return fmt.Errorf("coordinator config: status topic ID is required")
	}
	if c.PaymentTokenID.Token == 0 {
		return fmt.Errorf("coordinator config: payment token ID is required")
	}
	if c.TreasuryAccountID.Account == 0 {
		return fmt.Errorf("coordinator config: treasury account ID is required")
	}
	if c.DefaultPaymentAmount <= 0 {
		return fmt.Errorf("coordinator config: default payment amount must be positive")
	}
	return nil
}
