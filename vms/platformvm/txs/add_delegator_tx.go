// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"
	"fmt"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/bls"
	"github.com/skychains/chain/utils/math"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/platformvm/fx"
	"github.com/skychains/chain/vms/secp256k1fx"
)

var (
	_ DelegatorTx     = (*AddDelegatorTx)(nil)
	_ ScheduledStaker = (*AddDelegatorTx)(nil)

	errDelegatorWeightMismatch = errors.New("delegator weight is not equal to total stake weight")
	errStakeMustBeLUX         = errors.New("stake must be LUX")
)

// AddDelegatorTx is an unsigned addDelegatorTx
type AddDelegatorTx struct {
	// Metadata, inputs and outputs
	BaseTx `serialize:"true"`
	// Describes the delegatee
	Validator `serialize:"true" json:"validator"`
	// Where to send staked tokens when done validating
	StakeOuts []*lux.TransferableOutput `serialize:"true" json:"stake"`
	// Where to send staking rewards when done validating
	DelegationRewardsOwner fx.Owner `serialize:"true" json:"rewardsOwner"`
}

// InitCtx sets the FxID fields in the inputs and outputs of this
// [UnsignedAddDelegatorTx]. Also sets the [ctx] to the given [vm.ctx] so that
// the addresses can be json marshalled into human readable format
func (tx *AddDelegatorTx) InitCtx(ctx *snow.Context) {
	tx.BaseTx.InitCtx(ctx)
	for _, out := range tx.StakeOuts {
		out.FxID = secp256k1fx.ID
		out.InitCtx(ctx)
	}
	tx.DelegationRewardsOwner.InitCtx(ctx)
}

func (*AddDelegatorTx) SubnetID() ids.ID {
	return constants.PrimaryNetworkID
}

func (tx *AddDelegatorTx) NodeID() ids.NodeID {
	return tx.Validator.NodeID
}

func (*AddDelegatorTx) PublicKey() (*bls.PublicKey, bool, error) {
	return nil, false, nil
}

func (*AddDelegatorTx) PendingPriority() Priority {
	return PrimaryNetworkDelegatorApricotPendingPriority
}

func (*AddDelegatorTx) CurrentPriority() Priority {
	return PrimaryNetworkDelegatorCurrentPriority
}

func (tx *AddDelegatorTx) Stake() []*lux.TransferableOutput {
	return tx.StakeOuts
}

func (tx *AddDelegatorTx) RewardsOwner() fx.Owner {
	return tx.DelegationRewardsOwner
}

// SyntacticVerify returns nil iff [tx] is valid
func (tx *AddDelegatorTx) SyntacticVerify(ctx *snow.Context) error {
	switch {
	case tx == nil:
		return ErrNilTx
	case tx.SyntacticallyVerified: // already passed syntactic verification
		return nil
	}

	if err := tx.BaseTx.SyntacticVerify(ctx); err != nil {
		return err
	}
	if err := verify.All(&tx.Validator, tx.DelegationRewardsOwner); err != nil {
		return fmt.Errorf("failed to verify validator or rewards owner: %w", err)
	}

	totalStakeWeight := uint64(0)
	for _, out := range tx.StakeOuts {
		if err := out.Verify(); err != nil {
			return fmt.Errorf("output verification failed: %w", err)
		}
		newWeight, err := math.Add64(totalStakeWeight, out.Output().Amount())
		if err != nil {
			return err
		}
		totalStakeWeight = newWeight

		assetID := out.AssetID()
		if assetID != ctx.LUXAssetID {
			return fmt.Errorf("%w but is %q", errStakeMustBeLUX, assetID)
		}
	}

	switch {
	case !lux.IsSortedTransferableOutputs(tx.StakeOuts, Codec):
		return errOutputsNotSorted
	case totalStakeWeight != tx.Wght:
		return fmt.Errorf("%w, delegator weight %d total stake weight %d",
			errDelegatorWeightMismatch,
			tx.Wght,
			totalStakeWeight,
		)
	}

	// cache that this is valid
	tx.SyntacticallyVerified = true
	return nil
}

func (tx *AddDelegatorTx) Visit(visitor Visitor) error {
	return visitor.AddDelegatorTx(tx)
}
