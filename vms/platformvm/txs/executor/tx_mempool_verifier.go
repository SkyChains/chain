// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"errors"
	"fmt"
	"time"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/vms/platformvm/state"
	"github.com/skychains/chain/vms/platformvm/txs"
)

var _ txs.Visitor = (*MempoolTxVerifier)(nil)

type MempoolTxVerifier struct {
	*Backend
	ParentID      ids.ID
	StateVersions state.Versions
	Tx            *txs.Tx
}

func (*MempoolTxVerifier) AdvanceTimeTx(*txs.AdvanceTimeTx) error {
	return ErrWrongTxType
}

func (*MempoolTxVerifier) RewardValidatorTx(*txs.RewardValidatorTx) error {
	return ErrWrongTxType
}

func (v *MempoolTxVerifier) AddValidatorTx(tx *txs.AddValidatorTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) AddSubnetValidatorTx(tx *txs.AddSubnetValidatorTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) AddDelegatorTx(tx *txs.AddDelegatorTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) CreateChainTx(tx *txs.CreateChainTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) CreateSubnetTx(tx *txs.CreateSubnetTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) ImportTx(tx *txs.ImportTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) ExportTx(tx *txs.ExportTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) RemoveSubnetValidatorTx(tx *txs.RemoveSubnetValidatorTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) TransformSubnetTx(tx *txs.TransformSubnetTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) AddPermissionlessValidatorTx(tx *txs.AddPermissionlessValidatorTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) AddPermissionlessDelegatorTx(tx *txs.AddPermissionlessDelegatorTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) TransferSubnetOwnershipTx(tx *txs.TransferSubnetOwnershipTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) BaseTx(tx *txs.BaseTx) error {
	return v.standardTx(tx)
}

func (v *MempoolTxVerifier) standardTx(tx txs.UnsignedTx) error {
	baseState, err := v.standardBaseState()
	if err != nil {
		return err
	}

	executor := StandardTxExecutor{
		Backend: v.Backend,
		State:   baseState,
		Tx:      v.Tx,
	}
	err = tx.Visit(&executor)
	// We ignore [errFutureStakeTime] here because the time will be advanced
	// when this transaction is issued.
	if errors.Is(err, ErrFutureStakeTime) {
		return nil
	}
	return err
}

func (v *MempoolTxVerifier) standardBaseState() (state.Diff, error) {
	state, err := state.NewDiff(v.ParentID, v.StateVersions)
	if err != nil {
		return nil, err
	}

	nextBlkTime, err := v.nextBlockTime(state)
	if err != nil {
		return nil, err
	}

	changes, err := AdvanceTimeTo(v.Backend, state, nextBlkTime)
	if err != nil {
		return nil, err
	}
	changes.Apply(state)
	state.SetTimestamp(nextBlkTime)

	return state, nil
}

func (v *MempoolTxVerifier) nextBlockTime(state state.Diff) (time.Time, error) {
	var (
		parentTime  = state.GetTimestamp()
		nextBlkTime = v.Clk.Time()
	)
	if parentTime.After(nextBlkTime) {
		nextBlkTime = parentTime
	}
	nextStakerChangeTime, err := GetNextStakerChangeTime(state)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not calculate next staker change time: %w", err)
	}
	if !nextBlkTime.Before(nextStakerChangeTime) {
		nextBlkTime = nextStakerChangeTime
	}
	return nextBlkTime, nil
}
