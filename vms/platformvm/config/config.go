// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"time"

	"github.com/skychains/chain/chains"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/uptime"
	"github.com/skychains/chain/snow/validators"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/platformvm/reward"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/platformvm/txs/fee"
	"github.com/skychains/chain/vms/platformvm/upgrade"
)

// Struct collecting all foundational parameters of PlatformVM
type Config struct {
	// The node's chain manager
	Chains chains.Manager

	// Node's validator set maps subnetID -> validators of the subnet
	//
	// Invariant: The primary network's validator set should have been added to
	//            the manager before calling VM.Initialize.
	// Invariant: The primary network's validator set should be empty before
	//            calling VM.Initialize.
	Validators validators.Manager

	// All static fees config active before E-upgrade
	StaticFeeConfig fee.StaticConfig

	// Provides access to the uptime manager as a thread safe data structure
	UptimeLockedCalculator uptime.LockedCalculator

	// True if the node is being run with staking enabled
	SybilProtectionEnabled bool

	// If true, only the P-chain will be instantiated on the primary network.
	PartialSyncPrimaryNetwork bool

	// Set of subnets that this node is validating
	TrackedSubnets set.Set[ids.ID]

	// The minimum amount of tokens one must bond to be a validator
	MinValidatorStake uint64

	// The maximum amount of tokens that can be bonded on a validator
	MaxValidatorStake uint64

	// Minimum stake, in nLUX, that can be delegated on the primary network
	MinDelegatorStake uint64

	// Minimum fee that can be charged for delegation
	MinDelegationFee uint32

	// UptimePercentage is the minimum uptime required to be rewarded for staking
	UptimePercentage float64

	// Minimum amount of time to allow a staker to stake
	MinStakeDuration time.Duration

	// Maximum amount of time to allow a staker to stake
	MaxStakeDuration time.Duration

	// Config for the minting function
	RewardConfig reward.Config

	// All network upgrade timestamps
	UpgradeConfig upgrade.Config

	// UseCurrentHeight forces [GetMinimumHeight] to return the current height
	// of the P-Chain instead of the oldest block in the [recentlyAccepted]
	// window.
	//
	// This config is particularly useful for triggering proposervm activation
	// on recently created subnets (without this, users need to wait for
	// [recentlyAcceptedWindowTTL] to pass for activation to occur).
	UseCurrentHeight bool
}

// Create the blockchain described in [tx], but only if this node is a member of
// the subnet that validates the chain
func (c *Config) CreateChain(chainID ids.ID, tx *txs.CreateChainTx) {
	if c.SybilProtectionEnabled && // Sybil protection is enabled, so nodes might not validate all chains
		constants.PrimaryNetworkID != tx.SubnetID && // All nodes must validate the primary network
		!c.TrackedSubnets.Contains(tx.SubnetID) { // This node doesn't validate this blockchain
		return
	}

	chainParams := chains.ChainParameters{
		ID:          chainID,
		SubnetID:    tx.SubnetID,
		GenesisData: tx.GenesisData,
		VMID:        tx.VMID,
		FxIDs:       tx.FxIDs,
	}

	c.Chains.QueueChainCreation(chainParams)
}
