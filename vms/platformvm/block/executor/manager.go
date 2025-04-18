// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"errors"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/consensus/snowman"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/platformvm/block"
	"github.com/skychains/chain/vms/platformvm/metrics"
	"github.com/skychains/chain/vms/platformvm/state"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/platformvm/txs/executor"
	"github.com/skychains/chain/vms/platformvm/txs/mempool"
	"github.com/skychains/chain/vms/platformvm/validators"
)

var (
	_ Manager = (*manager)(nil)

	ErrChainNotSynced = errors.New("chain not synced")
)

type Manager interface {
	state.Versions

	// Returns the ID of the most recently accepted block.
	LastAccepted() ids.ID

	SetPreference(blkID ids.ID) (updated bool)
	Preferred() ids.ID

	GetBlock(blkID ids.ID) (snowman.Block, error)
	GetStatelessBlock(blkID ids.ID) (block.Block, error)
	NewBlock(block.Block) snowman.Block

	// VerifyTx verifies that the transaction can be issued based on the currently
	// preferred state. This should *not* be used to verify transactions in a block.
	VerifyTx(tx *txs.Tx) error

	// VerifyUniqueInputs verifies that the inputs are not duplicated in the
	// provided blk or any of its ancestors pinned in memory.
	VerifyUniqueInputs(blkID ids.ID, inputs set.Set[ids.ID]) error
}

func NewManager(
	mempool mempool.Mempool,
	metrics metrics.Metrics,
	s state.State,
	txExecutorBackend *executor.Backend,
	validatorManager validators.Manager,
) Manager {
	lastAccepted := s.GetLastAccepted()
	backend := &backend{
		Mempool:      mempool,
		lastAccepted: lastAccepted,
		state:        s,
		ctx:          txExecutorBackend.Ctx,
		blkIDToState: map[ids.ID]*blockState{},
	}

	return &manager{
		backend: backend,
		verifier: &verifier{
			backend:           backend,
			txExecutorBackend: txExecutorBackend,
		},
		acceptor: &acceptor{
			backend:      backend,
			metrics:      metrics,
			validators:   validatorManager,
			bootstrapped: txExecutorBackend.Bootstrapped,
		},
		rejector: &rejector{
			backend:         backend,
			addTxsToMempool: !txExecutorBackend.Config.PartialSyncPrimaryNetwork,
		},
		preferred:         lastAccepted,
		txExecutorBackend: txExecutorBackend,
	}
}

type manager struct {
	*backend
	verifier block.Visitor
	acceptor block.Visitor
	rejector block.Visitor

	preferred         ids.ID
	txExecutorBackend *executor.Backend
}

func (m *manager) GetBlock(blkID ids.ID) (snowman.Block, error) {
	blk, err := m.backend.GetBlock(blkID)
	if err != nil {
		return nil, err
	}
	return m.NewBlock(blk), nil
}

func (m *manager) GetStatelessBlock(blkID ids.ID) (block.Block, error) {
	return m.backend.GetBlock(blkID)
}

func (m *manager) NewBlock(blk block.Block) snowman.Block {
	return &Block{
		manager: m,
		Block:   blk,
	}
}

func (m *manager) SetPreference(blkID ids.ID) bool {
	updated := m.preferred != blkID
	m.preferred = blkID
	return updated
}

func (m *manager) Preferred() ids.ID {
	return m.preferred
}

func (m *manager) VerifyTx(tx *txs.Tx) error {
	if !m.txExecutorBackend.Bootstrapped.Get() {
		return ErrChainNotSynced
	}

	stateDiff, err := state.NewDiff(m.preferred, m)
	if err != nil {
		return err
	}

	nextBlkTime, _, err := state.NextBlockTime(stateDiff, m.txExecutorBackend.Clk)
	if err != nil {
		return err
	}

	_, err = executor.AdvanceTimeTo(m.txExecutorBackend, stateDiff, nextBlkTime)
	if err != nil {
		return err
	}

	return tx.Unsigned.Visit(&executor.StandardTxExecutor{
		Backend: m.txExecutorBackend,
		State:   stateDiff,
		Tx:      tx,
	})
}

func (m *manager) VerifyUniqueInputs(blkID ids.ID, inputs set.Set[ids.ID]) error {
	return m.backend.verifyUniqueInputs(blkID, inputs)
}
