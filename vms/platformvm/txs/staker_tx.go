// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"time"

	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/utils/crypto/bls"
	"github.com/SkyChains/chain/vms/components/lux"
	"github.com/SkyChains/chain/vms/platformvm/fx"
)

// ValidatorTx defines the interface for a validator transaction that supports
// delegation.
type ValidatorTx interface {
	UnsignedTx
	PermissionlessStaker

	ValidationRewardsOwner() fx.Owner
	DelegationRewardsOwner() fx.Owner
	Shares() uint32
}

type DelegatorTx interface {
	UnsignedTx
	PermissionlessStaker

	RewardsOwner() fx.Owner
}

type StakerTx interface {
	UnsignedTx
	Staker
}

type PermissionlessStaker interface {
	Staker

	Outputs() []*lux.TransferableOutput
	Stake() []*lux.TransferableOutput
}

type Staker interface {
	SubnetID() ids.ID
	NodeID() ids.NodeID
	// PublicKey returns the BLS public key registered by this transaction. If
	// there was no key registered by this transaction, it will return false.
	PublicKey() (*bls.PublicKey, bool, error)
	EndTime() time.Time
	Weight() uint64
	CurrentPriority() Priority
}

type ScheduledStaker interface {
	Staker
	StartTime() time.Time
	PendingPriority() Priority
}
