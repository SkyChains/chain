// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/components/lux"
)

var _ UnsignedTx = (*RewardValidatorTx)(nil)

// RewardValidatorTx is a transaction that represents a proposal to
// remove a validator that is currently validating from the validator set.
//
// If this transaction is accepted and the next block accepted is a Commit
// block, the validator is removed and the address that the validator specified
// receives the staked LUX as well as a validating reward.
//
// If this transaction is accepted and the next block accepted is an Abort
// block, the validator is removed and the address that the validator specified
// receives the staked LUX but no reward.
type RewardValidatorTx struct {
	// ID of the tx that created the delegator/validator being removed/rewarded
	TxID ids.ID `serialize:"true" json:"txID"`

	unsignedBytes []byte // Unsigned byte representation of this data
}

func (tx *RewardValidatorTx) SetBytes(unsignedBytes []byte) {
	tx.unsignedBytes = unsignedBytes
}

func (*RewardValidatorTx) InitCtx(*snow.Context) {}

func (tx *RewardValidatorTx) Bytes() []byte {
	return tx.unsignedBytes
}

func (*RewardValidatorTx) InputIDs() set.Set[ids.ID] {
	return nil
}

func (*RewardValidatorTx) Outputs() []*lux.TransferableOutput {
	return nil
}

func (*RewardValidatorTx) SyntacticVerify(*snow.Context) error {
	return nil
}

func (tx *RewardValidatorTx) Visit(visitor Visitor) error {
	return visitor.RewardValidatorTx(tx)
}
