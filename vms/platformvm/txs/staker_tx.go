// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"time"

	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/vms/components/lux"
	"github.com/luxdefi/luxd/vms/platformvm/fx"
)

// ValidatorTx defines the interface for a validator transaction that supports
// delegation.
type ValidatorTx interface {
	UnsignedTx
	Validator
}

type DelegatorTx interface {
	UnsignedTx
	Delegator
}

type StakerTx interface {
	UnsignedTx
	Staker
}

type Validator interface {
	PermissionlessStaker

	ValidationRewardsOwner() fx.Owner
	DelegationRewardsOwner() fx.Owner
	Shares() uint32
}

type Delegator interface {
	PermissionlessStaker

	RewardsOwner() fx.Owner
}

type PermissionlessStaker interface {
	Staker

	Outputs() []*lux.TransferableOutput
	Stake() []*lux.TransferableOutput
}

type Staker interface {
	SubnetID() ids.ID
	NodeID() ids.NodeID
	StartTime() time.Time
	EndTime() time.Time
	Weight() uint64
	PendingPriority() Priority
	CurrentPriority() Priority
}
