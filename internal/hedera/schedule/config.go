package schedule

import (
	"fmt"
	"time"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

const (
	defaultHeartbeatInterval = 30 * time.Second
	defaultHeartbeatMemo     = "agent-heartbeat"
	minHeartbeatInterval     = 5 * time.Second
)

// HeartbeatConfig holds configuration for the heartbeat runner.
type HeartbeatConfig struct {
	// Interval between heartbeat submissions.
	Interval time.Duration

	// Memo attached to each heartbeat scheduled transaction.
	Memo string

	// AgentID identifies which agent this heartbeat belongs to.
	AgentID string

	// AccountID is the Hedera account submitting the heartbeat.
	AccountID hiero.AccountID

	// TopicID is an optional HCS topic to publish heartbeat notifications to.
	TopicID *hiero.TopicID
}

// DefaultHeartbeatConfig returns sensible defaults for testnet usage.
// The caller must set AgentID and AccountID.
func DefaultHeartbeatConfig() HeartbeatConfig {
	return HeartbeatConfig{
		Interval: defaultHeartbeatInterval,
		Memo:     defaultHeartbeatMemo,
	}
}

// Validate checks the heartbeat config for required fields and valid values.
func (c HeartbeatConfig) Validate() error {
	if c.Interval < minHeartbeatInterval {
		return fmt.Errorf("heartbeat interval %v is below minimum %v", c.Interval, minHeartbeatInterval)
	}
	if c.AgentID == "" {
		return fmt.Errorf("heartbeat agent ID is required")
	}
	if c.AccountID.Account == 0 {
		return fmt.Errorf("heartbeat account ID is required")
	}
	return nil
}
