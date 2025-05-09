// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package local

import (
	"time"

	"github.com/skychains/chain/config"
	"github.com/skychains/chain/tests/fixture/tmpnet"
)

const (
	// Constants defining the names of shell variables whose value can
	// configure local network orchestration.
	LuxdPathEnvName   = "LUXD_PATH"
	NetworkDirEnvName = "TMPNET_NETWORK_DIR"
	RootDirEnvName    = "TMPNET_ROOT_DIR"

	DefaultNetworkStartTimeout = 2 * time.Minute
	DefaultNodeInitTimeout     = 10 * time.Second
	DefaultNodeStopTimeout     = 5 * time.Second
)

// A set of flags appropriate for local testing.
func LocalFlags() tmpnet.FlagsMap {
	// Supply only non-default configuration to ensure that default values will be used.
	return tmpnet.FlagsMap{
		config.NetworkPeerListGossipFreqKey: "250ms",
		config.NetworkMaxReconnectDelayKey:  "1s",
		config.PublicIPKey:                  "127.0.0.1",
		config.HTTPHostKey:                  "127.0.0.1",
		config.StakingHostKey:               "127.0.0.1",
		config.HealthCheckFreqKey:           "2s",
		config.AdminAPIEnabledKey:           true,
		config.IpcAPIEnabledKey:             true,
		config.IndexEnabledKey:              true,
		config.LogDisplayLevelKey:           "INFO",
		config.LogLevelKey:                  "DEBUG",
		config.MinStakeDurationKey:          tmpnet.DefaultMinStakeDuration.String(),
	}
}

// C-Chain config for local testing.
func LocalCChainConfig() tmpnet.FlagsMap {
	// Supply only non-default configuration to ensure that default
	// values will be used. Available C-Chain configuration options are
	// defined in the `github.com/skychains/coreth/evm` package.
	return tmpnet.FlagsMap{
		"log-level": "trace",
	}
}
