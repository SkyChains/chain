// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"

	"github.com/luxdefi/node/chains/atomic"
	"github.com/luxdefi/node/database"
	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/snow"
	"github.com/luxdefi/node/utils"
	"github.com/luxdefi/node/utils/logging"
	"github.com/luxdefi/node/utils/timer/mockable"
	"github.com/luxdefi/node/vms/components/verify"
	"github.com/luxdefi/node/vms/platformvm/block"
	"github.com/luxdefi/node/vms/platformvm/metrics"
	"github.com/luxdefi/node/vms/platformvm/state"
	"github.com/luxdefi/node/vms/platformvm/txs"
	"github.com/luxdefi/node/vms/platformvm/validators"
	"github.com/luxdefi/node/vms/secp256k1fx"
)

func TestAcceptorVisitProposalBlock(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)

	lastAcceptedID := ids.GenerateTestID()

	blk, err := block.NewApricotProposalBlock(
		lastAcceptedID,
		1,
		&txs.Tx{
			Unsigned: &txs.AddDelegatorTx{
				// Without the line below, this function will error.
				DelegationRewardsOwner: &secp256k1fx.OutputOwners{},
			},
			Creds: []verify.Verifiable{},
		},
	)
	require.NoError(err)

	blkID := blk.ID()

	s := state.NewMockState(ctrl)
	s.EXPECT().Checksum().Return(ids.Empty).Times(1)

	acceptor := &acceptor{
		backend: &backend{
			ctx: &snow.Context{
				Log: logging.NoLog{},
			},
			blkIDToState: map[ids.ID]*blockState{
				blkID: {},
			},
			state: s,
		},
		metrics:    metrics.Noop,
		validators: validators.TestManager,
	}

	require.NoError(acceptor.ApricotProposalBlock(blk))

	require.Equal(blkID, acceptor.backend.lastAccepted)

	_, exists := acceptor.GetState(blkID)
	require.False(exists)

	s.EXPECT().GetLastAccepted().Return(lastAcceptedID).Times(1)

	_, exists = acceptor.GetState(lastAcceptedID)
	require.True(exists)
}

func TestAcceptorVisitAtomicBlock(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)

	s := state.NewMockState(ctrl)
	sharedMemory := atomic.NewMockSharedMemory(ctrl)

	parentID := ids.GenerateTestID()
	acceptor := &acceptor{
		backend: &backend{
			lastAccepted: parentID,
			blkIDToState: make(map[ids.ID]*blockState),
			state:        s,
			ctx: &snow.Context{
				Log:          logging.NoLog{},
				SharedMemory: sharedMemory,
			},
		},
		metrics:    metrics.Noop,
		validators: validators.TestManager,
	}

	blk, err := block.NewApricotAtomicBlock(
		parentID,
		1,
		&txs.Tx{
			Unsigned: &txs.AddDelegatorTx{
				// Without the line below, this function will error.
				DelegationRewardsOwner: &secp256k1fx.OutputOwners{},
			},
			Creds: []verify.Verifiable{},
		},
	)
	require.NoError(err)

	// Set expected calls on the state.
	// We should error after [commonAccept] is called.
	s.EXPECT().SetLastAccepted(blk.ID()).Times(1)
	s.EXPECT().SetHeight(blk.Height()).Times(1)
	s.EXPECT().AddStatelessBlock(blk).Times(1)

	err = acceptor.ApricotAtomicBlock(blk)
	require.ErrorIs(err, errMissingBlockState)

	// Set [blk]'s state in the map as though it had been verified.
	onAcceptState := state.NewMockDiff(ctrl)
	childID := ids.GenerateTestID()
	atomicRequests := map[ids.ID]*atomic.Requests{ids.GenerateTestID(): nil}
	acceptor.backend.blkIDToState[blk.ID()] = &blockState{
		onAcceptState:  onAcceptState,
		atomicRequests: atomicRequests,
	}
	// Give [blk] a child.
	childOnAcceptState := state.NewMockDiff(ctrl)
	childOnAbortState := state.NewMockDiff(ctrl)
	childOnCommitState := state.NewMockDiff(ctrl)
	childState := &blockState{
		onAcceptState: childOnAcceptState,
		proposalBlockState: proposalBlockState{
			onAbortState:  childOnAbortState,
			onCommitState: childOnCommitState,
		},
	}
	acceptor.backend.blkIDToState[childID] = childState

	// Set expected calls on dependencies.
	s.EXPECT().SetLastAccepted(blk.ID()).Times(1)
	s.EXPECT().SetHeight(blk.Height()).Times(1)
	s.EXPECT().AddStatelessBlock(blk).Times(1)
	batch := database.NewMockBatch(ctrl)
	s.EXPECT().CommitBatch().Return(batch, nil).Times(1)
	s.EXPECT().Abort().Times(1)
	onAcceptState.EXPECT().Apply(s).Times(1)
	sharedMemory.EXPECT().Apply(atomicRequests, batch).Return(nil).Times(1)
	s.EXPECT().Checksum().Return(ids.Empty).Times(1)

	require.NoError(acceptor.ApricotAtomicBlock(blk))
}

func TestAcceptorVisitStandardBlock(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)

	s := state.NewMockState(ctrl)
	sharedMemory := atomic.NewMockSharedMemory(ctrl)

	parentID := ids.GenerateTestID()
	clk := &mockable.Clock{}
	acceptor := &acceptor{
		backend: &backend{
			lastAccepted: parentID,
			blkIDToState: make(map[ids.ID]*blockState),
			state:        s,
			ctx: &snow.Context{
				Log:          logging.NoLog{},
				SharedMemory: sharedMemory,
			},
		},
		metrics:    metrics.Noop,
		validators: validators.TestManager,
	}

	blk, err := block.NewBanffStandardBlock(
		clk.Time(),
		parentID,
		1,
		[]*txs.Tx{
			{
				Unsigned: &txs.AddDelegatorTx{
					// Without the line below, this function will error.
					DelegationRewardsOwner: &secp256k1fx.OutputOwners{},
				},
				Creds: []verify.Verifiable{},
			},
		},
	)
	require.NoError(err)

	// Set expected calls on the state.
	// We should error after [commonAccept] is called.
	s.EXPECT().SetLastAccepted(blk.ID()).Times(1)
	s.EXPECT().SetHeight(blk.Height()).Times(1)
	s.EXPECT().AddStatelessBlock(blk).Times(1)

	err = acceptor.BanffStandardBlock(blk)
	require.ErrorIs(err, errMissingBlockState)

	// Set [blk]'s state in the map as though it had been verified.
	onAcceptState := state.NewMockDiff(ctrl)
	childID := ids.GenerateTestID()
	atomicRequests := map[ids.ID]*atomic.Requests{ids.GenerateTestID(): nil}
	calledOnAcceptFunc := false
	acceptor.backend.blkIDToState[blk.ID()] = &blockState{
		onAcceptState:  onAcceptState,
		atomicRequests: atomicRequests,
		standardBlockState: standardBlockState{
			onAcceptFunc: func() {
				calledOnAcceptFunc = true
			},
		},
	}
	// Give [blk] a child.
	childOnAcceptState := state.NewMockDiff(ctrl)
	childOnAbortState := state.NewMockDiff(ctrl)
	childOnCommitState := state.NewMockDiff(ctrl)
	childState := &blockState{
		onAcceptState: childOnAcceptState,
		proposalBlockState: proposalBlockState{
			onAbortState:  childOnAbortState,
			onCommitState: childOnCommitState,
		},
	}
	acceptor.backend.blkIDToState[childID] = childState

	// Set expected calls on dependencies.
	s.EXPECT().SetLastAccepted(blk.ID()).Times(1)
	s.EXPECT().SetHeight(blk.Height()).Times(1)
	s.EXPECT().AddStatelessBlock(blk).Times(1)
	batch := database.NewMockBatch(ctrl)
	s.EXPECT().CommitBatch().Return(batch, nil).Times(1)
	s.EXPECT().Abort().Times(1)
	onAcceptState.EXPECT().Apply(s).Times(1)
	sharedMemory.EXPECT().Apply(atomicRequests, batch).Return(nil).Times(1)
	s.EXPECT().Checksum().Return(ids.Empty).Times(1)

	require.NoError(acceptor.BanffStandardBlock(blk))
	require.True(calledOnAcceptFunc)
	require.Equal(blk.ID(), acceptor.backend.lastAccepted)
}

func TestAcceptorVisitCommitBlock(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)

	s := state.NewMockState(ctrl)
	sharedMemory := atomic.NewMockSharedMemory(ctrl)

	parentID := ids.GenerateTestID()
	acceptor := &acceptor{
		backend: &backend{
			lastAccepted: parentID,
			blkIDToState: make(map[ids.ID]*blockState),
			state:        s,
			ctx: &snow.Context{
				Log:          logging.NoLog{},
				SharedMemory: sharedMemory,
			},
		},
		metrics:      metrics.Noop,
		validators:   validators.TestManager,
		bootstrapped: &utils.Atomic[bool]{},
	}

	blk, err := block.NewApricotCommitBlock(parentID, 1 /*height*/)
	require.NoError(err)

	err = acceptor.ApricotCommitBlock(blk)
	require.ErrorIs(err, state.ErrMissingParentState)

	// Set [blk]'s parent in the state map.
	parentOnAcceptState := state.NewMockDiff(ctrl)
	parentOnAbortState := state.NewMockDiff(ctrl)
	parentOnCommitState := state.NewMockDiff(ctrl)
	parentStatelessBlk := block.NewMockBlock(ctrl)
	parentState := &blockState{
		statelessBlock: parentStatelessBlk,
		onAcceptState:  parentOnAcceptState,
		proposalBlockState: proposalBlockState{
			onAbortState:  parentOnAbortState,
			onCommitState: parentOnCommitState,
		},
	}
	acceptor.backend.blkIDToState[parentID] = parentState

	blkID := blk.ID()
	// Set expected calls on dependencies.
	// Make sure the parent is accepted first.
	gomock.InOrder(
		parentStatelessBlk.EXPECT().ID().Return(parentID).Times(2),
		s.EXPECT().SetLastAccepted(parentID).Times(1),
		parentStatelessBlk.EXPECT().Height().Return(blk.Height()-1).Times(1),
		s.EXPECT().SetHeight(blk.Height()-1).Times(1),
		s.EXPECT().AddStatelessBlock(parentState.statelessBlock).Times(1),

		s.EXPECT().SetLastAccepted(blkID).Times(1),
		s.EXPECT().SetHeight(blk.Height()).Times(1),
		s.EXPECT().AddStatelessBlock(blk).Times(1),
	)

	err = acceptor.ApricotCommitBlock(blk)
	require.ErrorIs(err, errMissingBlockState)

	// Set [blk]'s state in the map as though it had been verified.
	acceptor.backend.blkIDToState[parentID] = parentState
	onAcceptState := state.NewMockDiff(ctrl)
	acceptor.backend.blkIDToState[blkID] = &blockState{
		onAcceptState: onAcceptState,
	}

	// Set expected calls on dependencies.
	// Make sure the parent is accepted first.
	gomock.InOrder(
		parentStatelessBlk.EXPECT().ID().Return(parentID).Times(2),
		s.EXPECT().SetLastAccepted(parentID).Times(1),
		parentStatelessBlk.EXPECT().Height().Return(blk.Height()-1).Times(1),
		s.EXPECT().SetHeight(blk.Height()-1).Times(1),
		s.EXPECT().AddStatelessBlock(parentState.statelessBlock).Times(1),

		s.EXPECT().SetLastAccepted(blkID).Times(1),
		s.EXPECT().SetHeight(blk.Height()).Times(1),
		s.EXPECT().AddStatelessBlock(blk).Times(1),

		onAcceptState.EXPECT().Apply(s).Times(1),
		s.EXPECT().Commit().Return(nil).Times(1),
		s.EXPECT().Checksum().Return(ids.Empty).Times(1),
	)

	require.NoError(acceptor.ApricotCommitBlock(blk))
	require.Equal(blk.ID(), acceptor.backend.lastAccepted)
}

func TestAcceptorVisitAbortBlock(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)

	s := state.NewMockState(ctrl)
	sharedMemory := atomic.NewMockSharedMemory(ctrl)

	parentID := ids.GenerateTestID()
	acceptor := &acceptor{
		backend: &backend{
			lastAccepted: parentID,
			blkIDToState: make(map[ids.ID]*blockState),
			state:        s,
			ctx: &snow.Context{
				Log:          logging.NoLog{},
				SharedMemory: sharedMemory,
			},
		},
		metrics:      metrics.Noop,
		validators:   validators.TestManager,
		bootstrapped: &utils.Atomic[bool]{},
	}

	blk, err := block.NewApricotAbortBlock(parentID, 1 /*height*/)
	require.NoError(err)

	err = acceptor.ApricotAbortBlock(blk)
	require.ErrorIs(err, state.ErrMissingParentState)

	// Set [blk]'s parent in the state map.
	parentOnAcceptState := state.NewMockDiff(ctrl)
	parentOnAbortState := state.NewMockDiff(ctrl)
	parentOnCommitState := state.NewMockDiff(ctrl)
	parentStatelessBlk := block.NewMockBlock(ctrl)
	parentState := &blockState{
		statelessBlock: parentStatelessBlk,
		onAcceptState:  parentOnAcceptState,
		proposalBlockState: proposalBlockState{
			onAbortState:  parentOnAbortState,
			onCommitState: parentOnCommitState,
		},
	}
	acceptor.backend.blkIDToState[parentID] = parentState

	blkID := blk.ID()
	// Set expected calls on dependencies.
	// Make sure the parent is accepted first.
	gomock.InOrder(
		parentStatelessBlk.EXPECT().ID().Return(parentID).Times(2),
		s.EXPECT().SetLastAccepted(parentID).Times(1),
		parentStatelessBlk.EXPECT().Height().Return(blk.Height()-1).Times(1),
		s.EXPECT().SetHeight(blk.Height()-1).Times(1),
		s.EXPECT().AddStatelessBlock(parentState.statelessBlock).Times(1),

		s.EXPECT().SetLastAccepted(blkID).Times(1),
		s.EXPECT().SetHeight(blk.Height()).Times(1),
		s.EXPECT().AddStatelessBlock(blk).Times(1),
	)

	err = acceptor.ApricotAbortBlock(blk)
	require.ErrorIs(err, errMissingBlockState)

	// Set [blk]'s state in the map as though it had been verified.
	acceptor.backend.blkIDToState[parentID] = parentState

	onAcceptState := state.NewMockDiff(ctrl)
	acceptor.backend.blkIDToState[blkID] = &blockState{
		onAcceptState: onAcceptState,
	}

	// Set expected calls on dependencies.
	// Make sure the parent is accepted first.
	gomock.InOrder(
		parentStatelessBlk.EXPECT().ID().Return(parentID).Times(2),
		s.EXPECT().SetLastAccepted(parentID).Times(1),
		parentStatelessBlk.EXPECT().Height().Return(blk.Height()-1).Times(1),
		s.EXPECT().SetHeight(blk.Height()-1).Times(1),
		s.EXPECT().AddStatelessBlock(parentState.statelessBlock).Times(1),

		s.EXPECT().SetLastAccepted(blkID).Times(1),
		s.EXPECT().SetHeight(blk.Height()).Times(1),
		s.EXPECT().AddStatelessBlock(blk).Times(1),

		onAcceptState.EXPECT().Apply(s).Times(1),
		s.EXPECT().Commit().Return(nil).Times(1),
		s.EXPECT().Checksum().Return(ids.Empty).Times(1),
	)

	require.NoError(acceptor.ApricotAbortBlock(blk))
	require.Equal(blk.ID(), acceptor.backend.lastAccepted)
}
