package daemon

import (
	"fmt"
	"time"
)

const (
	defaultDaemonAddr  = "localhost:50051"
	defaultDialTimeout = 10 * time.Second
	defaultCallTimeout = 30 * time.Second
)

// Config holds configuration for connecting to the obey daemon.
type Config struct {
	// Address is the daemon gRPC endpoint (host:port).
	Address string

	// DialTimeout is the maximum time to wait for initial connection.
	DialTimeout time.Duration

	// CallTimeout is the default timeout for individual RPC calls.
	CallTimeout time.Duration

	// TLSEnabled enables TLS for the gRPC connection.
	TLSEnabled bool

	// TLSCertPath is the path to the TLS certificate (required if TLSEnabled).
	TLSCertPath string
}

// DefaultConfig returns sensible defaults for local development.
func DefaultConfig() Config {
	return Config{
		Address:     defaultDaemonAddr,
		DialTimeout: defaultDialTimeout,
		CallTimeout: defaultCallTimeout,
		TLSEnabled:  false,
	}
}

// Validate checks the config for required fields.
func (c Config) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("daemon config: address is required")
	}
	if c.DialTimeout <= 0 {
		return fmt.Errorf("daemon config: dial timeout must be positive")
	}
	if c.CallTimeout <= 0 {
		return fmt.Errorf("daemon config: call timeout must be positive")
	}
	if c.TLSEnabled && c.TLSCertPath == "" {
		return fmt.Errorf("daemon config: TLS cert path required when TLS is enabled")
	}
	return nil
}
