// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/luxdefi/luxd/database"
	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/snow/choices"
	"github.com/luxdefi/luxd/snow/consensus/snowman"
	"github.com/luxdefi/luxd/utils/constants"
	"github.com/luxdefi/luxd/utils/crypto"
	"github.com/luxdefi/luxd/vms/components/lux"
	"github.com/luxdefi/luxd/vms/platformvm/blocks"
	"github.com/luxdefi/luxd/vms/platformvm/reward"
	"github.com/luxdefi/luxd/vms/platformvm/state"
	"github.com/luxdefi/luxd/vms/platformvm/status"
	"github.com/luxdefi/luxd/vms/platformvm/txs"
	"github.com/luxdefi/luxd/vms/platformvm/txs/executor"
	"github.com/luxdefi/luxd/vms/platformvm/validator"
	"github.com/luxdefi/luxd/vms/secp256k1fx"
)

func TestApricotProposalBlockTimeVerification(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	env := newEnvironment(t, ctrl)
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	// create apricotParentBlk. It's a standard one for simplicity
	parentHeight := uint64(2022)

	apricotParentBlk, err := blocks.NewApricotStandardBlock(
		ids.Empty, // does not matter
		parentHeight,
		nil, // txs do not matter in this test
	)
	require.NoError(err)
	parentID := apricotParentBlk.ID()

	// store parent block, with relevant quantities
	onParentAccept := state.NewMockDiff(ctrl)
	env.blkManager.(*manager).blkIDToState[parentID] = &blockState{
		statelessBlock: apricotParentBlk,
		onAcceptState:  onParentAccept,
	}
	env.blkManager.(*manager).lastAccepted = parentID
	chainTime := env.clk.Time().Truncate(time.Second)
	env.mockedState.EXPECT().GetTimestamp().Return(chainTime).AnyTimes()
	env.mockedState.EXPECT().GetLastAccepted().Return(parentID).AnyTimes()

	// create a proposal transaction to be included into proposal block
	utx := &txs.AddValidatorTx{
		BaseTx:    txs.BaseTx{},
		Validator: validator.Validator{End: uint64(chainTime.Unix())},
		StakeOuts: []*lux.TransferableOutput{
			{
				Asset: lux.Asset{
					ID: env.ctx.LUXAssetID,
				},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			},
		},
		RewardsOwner:     &secp256k1fx.OutputOwners{},
		DelegationShares: uint32(defaultTxFee),
	}
	addValTx := &txs.Tx{Unsigned: utx}
	require.NoError(addValTx.Sign(txs.Codec, nil))
	blkTx := &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: addValTx.ID(),
		},
	}

	// setup state to validate proposal block transaction
	onParentAccept.EXPECT().GetTimestamp().Return(chainTime).AnyTimes()

	currentStakersIt := state.NewMockStakerIterator(ctrl)
	currentStakersIt.EXPECT().Next().Return(true)
	currentStakersIt.EXPECT().Value().Return(&state.Staker{
		TxID:      addValTx.ID(),
		NodeID:    utx.NodeID(),
		SubnetID:  utx.SubnetID(),
		StartTime: utx.StartTime(),
		EndTime:   chainTime,
	}).Times(2)
	currentStakersIt.EXPECT().Release()
	onParentAccept.EXPECT().GetCurrentStakerIterator().Return(currentStakersIt, nil)
	onParentAccept.EXPECT().GetCurrentValidator(utx.SubnetID(), utx.NodeID()).Return(&state.Staker{
		TxID:      addValTx.ID(),
		NodeID:    utx.NodeID(),
		SubnetID:  utx.SubnetID(),
		StartTime: utx.StartTime(),
		EndTime:   chainTime,
	}, nil)
	onParentAccept.EXPECT().GetTx(addValTx.ID()).Return(addValTx, status.Committed, nil)
	onParentAccept.EXPECT().GetCurrentSupply(constants.PrimaryNetworkID).Return(uint64(1000), nil).AnyTimes()

	env.mockedState.EXPECT().GetUptime(gomock.Any()).Return(
		time.Duration(1000), /*upDuration*/
		time.Time{},         /*lastUpdated*/
		nil,                 /*err*/
	).AnyTimes()

	// wrong height
	statelessProposalBlock, err := blocks.NewApricotProposalBlock(
		parentID,
		parentHeight,
		blkTx,
	)
	require.NoError(err)

	block := env.blkManager.NewBlock(statelessProposalBlock)
	require.Error(block.Verify())

	// valid
	statelessProposalBlock, err = blocks.NewApricotProposalBlock(
		parentID,
		parentHeight+1,
		blkTx,
	)
	require.NoError(err)

	block = env.blkManager.NewBlock(statelessProposalBlock)
	require.NoError(block.Verify())
}

func TestBanffProposalBlockTimeVerification(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	env := newEnvironment(t, ctrl)
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()
	env.clk.Set(defaultGenesisTime)
	env.config.BanffTime = time.Time{} // activate Banff

	// create parentBlock. It's a standard one for simplicity
	parentTime := defaultGenesisTime
	parentHeight := uint64(2022)

	banffParentBlk, err := blocks.NewApricotStandardBlock(
		genesisBlkID, // does not matter
		parentHeight,
		nil, // txs do not matter in this test
	)
	require.NoError(err)
	parentID := banffParentBlk.ID()

	// store parent block, with relevant quantities
	chainTime := parentTime
	env.mockedState.EXPECT().GetTimestamp().Return(chainTime).AnyTimes()

	onParentAccept := state.NewMockDiff(ctrl)
	onParentAccept.EXPECT().GetTimestamp().Return(parentTime).AnyTimes()
	onParentAccept.EXPECT().GetCurrentSupply(constants.PrimaryNetworkID).Return(uint64(1000), nil).AnyTimes()

	env.blkManager.(*manager).blkIDToState[parentID] = &blockState{
		statelessBlock: banffParentBlk,
		onAcceptState:  onParentAccept,
		timestamp:      parentTime,
	}
	env.blkManager.(*manager).lastAccepted = parentID
	env.mockedState.EXPECT().GetLastAccepted().Return(parentID).AnyTimes()
	env.mockedState.EXPECT().GetStatelessBlock(gomock.Any()).DoAndReturn(
		func(blockID ids.ID) (blocks.Block, choices.Status, error) {
			if blockID == parentID {
				return banffParentBlk, choices.Accepted, nil
			}
			return nil, choices.Rejected, database.ErrNotFound
		}).AnyTimes()

	// setup state to validate proposal block transaction
	nextStakerTime := chainTime.Add(executor.SyncBound).Add(-1 * time.Second)
	unsignedNextStakerTx := &txs.AddValidatorTx{
		BaseTx:    txs.BaseTx{},
		Validator: validator.Validator{End: uint64(nextStakerTime.Unix())},
		StakeOuts: []*lux.TransferableOutput{
			{
				Asset: lux.Asset{
					ID: env.ctx.LUXAssetID,
				},
				Out: &secp256k1fx.TransferOutput{
					Amt: 1,
				},
			},
		},
		RewardsOwner:     &secp256k1fx.OutputOwners{},
		DelegationShares: uint32(defaultTxFee),
	}
	nextStakerTx := &txs.Tx{Unsigned: unsignedNextStakerTx}
	require.NoError(nextStakerTx.Sign(txs.Codec, nil))

	nextStakerTxID := nextStakerTx.ID()
	onParentAccept.EXPECT().GetCurrentValidator(unsignedNextStakerTx.SubnetID(), unsignedNextStakerTx.NodeID()).Return(&state.Staker{
		TxID:      nextStakerTxID,
		NodeID:    unsignedNextStakerTx.NodeID(),
		SubnetID:  unsignedNextStakerTx.SubnetID(),
		StartTime: unsignedNextStakerTx.StartTime(),
		EndTime:   chainTime,
	}, nil)
	onParentAccept.EXPECT().GetTx(nextStakerTxID).Return(nextStakerTx, status.Processing, nil)

	currentStakersIt := state.NewMockStakerIterator(ctrl)
	currentStakersIt.EXPECT().Next().Return(true).AnyTimes()
	currentStakersIt.EXPECT().Value().Return(&state.Staker{
		TxID:     nextStakerTxID,
		EndTime:  nextStakerTime,
		NextTime: nextStakerTime,
		Priority: txs.PrimaryNetworkValidatorCurrentPriority,
	}).AnyTimes()
	currentStakersIt.EXPECT().Release().AnyTimes()
	onParentAccept.EXPECT().GetCurrentStakerIterator().Return(currentStakersIt, nil).AnyTimes()

	pendingStakersIt := state.NewMockStakerIterator(ctrl)
	pendingStakersIt.EXPECT().Next().Return(false).AnyTimes() // no pending stakers
	pendingStakersIt.EXPECT().Release().AnyTimes()
	onParentAccept.EXPECT().GetPendingStakerIterator().Return(pendingStakersIt, nil).AnyTimes()

	env.mockedState.EXPECT().GetUptime(gomock.Any()).Return(
		time.Duration(1000), /*upDuration*/
		time.Time{},         /*lastUpdated*/
		nil,                 /*err*/
	).AnyTimes()

	// create proposal tx to be included in the proposal block
	blkTx := &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: nextStakerTxID,
		},
	}
	require.NoError(blkTx.Sign(txs.Codec, nil))

	{
		// wrong height
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			parentTime.Add(time.Second),
			parentID,
			banffParentBlk.Height(),
			blkTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.Error(block.Verify())
	}

	{
		// wrong version
		statelessProposalBlock, err := blocks.NewApricotProposalBlock(
			parentID,
			banffParentBlk.Height()+1,
			blkTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.Error(block.Verify())
	}

	{
		// wrong timestamp, earlier than parent
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			parentTime.Add(-1*time.Second),
			parentID,
			banffParentBlk.Height()+1,
			blkTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.Error(block.Verify())
	}

	{
		// wrong timestamp, violated synchrony bound
		beyondSyncBoundTimeStamp := env.clk.Time().Add(executor.SyncBound).Add(time.Second)
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			beyondSyncBoundTimeStamp,
			parentID,
			banffParentBlk.Height()+1,
			blkTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.Error(block.Verify())
	}

	{
		// wrong timestamp, skipped staker set change event
		skippedStakerEventTimeStamp := nextStakerTime.Add(time.Second)
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			skippedStakerEventTimeStamp,
			parentID,
			banffParentBlk.Height()+1,
			blkTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.Error(block.Verify())
	}

	{
		// wrong tx content (no advance time txs)
		invalidTx := &txs.Tx{
			Unsigned: &txs.AdvanceTimeTx{
				Time: uint64(nextStakerTime.Unix()),
			},
		}
		require.NoError(invalidTx.Sign(txs.Codec, nil))
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			parentTime.Add(time.Second),
			parentID,
			banffParentBlk.Height()+1,
			invalidTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.Error(block.Verify())
	}

	{
		// include too many transactions
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			nextStakerTime,
			parentID,
			banffParentBlk.Height()+1,
			blkTx,
		)
		require.NoError(err)

		statelessProposalBlock.Transactions = []*txs.Tx{blkTx}
		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.ErrorIs(block.Verify(), errBanffProposalBlockWithMultipleTransactions)
	}

	{
		// valid
		statelessProposalBlock, err := blocks.NewBanffProposalBlock(
			nextStakerTime,
			parentID,
			banffParentBlk.Height()+1,
			blkTx,
		)
		require.NoError(err)

		block := env.blkManager.NewBlock(statelessProposalBlock)
		require.NoError(block.Verify())
	}
}

func TestBanffProposalBlockUpdateStakers(t *testing.T) {
	// Chronological order (not in scale):
	// Staker0:    |--- ??? // Staker0 end time depends on the test
	// Staker1:        |------------------------------------------------------|
	// Staker2:            |------------------------|
	// Staker3:                |------------------------|
	// Staker3sub:                 |----------------|
	// Staker4:                |------------------------|
	// Staker5:                                     |--------------------|

	// Staker0 it's here just to allow to issue a proposal block with the chosen endTime.
	staker0RewardAddress := ids.GenerateTestShortID()
	staker0 := staker{
		nodeID:        ids.NodeID(staker0RewardAddress),
		rewardAddress: staker0RewardAddress,
		startTime:     defaultGenesisTime,
		endTime:       time.Time{}, // actual endTime depends on specific test
	}

	staker1RewardAddress := ids.GenerateTestShortID()
	staker1 := staker{
		nodeID:        ids.NodeID(staker1RewardAddress),
		rewardAddress: staker1RewardAddress,
		startTime:     defaultGenesisTime.Add(1 * time.Minute),
		endTime:       defaultGenesisTime.Add(10 * defaultMinStakingDuration).Add(1 * time.Minute),
	}

	staker2RewardAddress := ids.GenerateTestShortID()
	staker2 := staker{
		nodeID:        ids.NodeID(staker2RewardAddress),
		rewardAddress: staker2RewardAddress,
		startTime:     staker1.startTime.Add(1 * time.Minute),
		endTime:       staker1.startTime.Add(1 * time.Minute).Add(defaultMinStakingDuration),
	}

	staker3RewardAddress := ids.GenerateTestShortID()
	staker3 := staker{
		nodeID:        ids.NodeID(staker3RewardAddress),
		rewardAddress: staker3RewardAddress,
		startTime:     staker2.startTime.Add(1 * time.Minute),
		endTime:       staker2.endTime.Add(1 * time.Minute),
	}

	staker3Sub := staker{
		nodeID:        staker3.nodeID,
		rewardAddress: staker3.rewardAddress,
		startTime:     staker3.startTime.Add(1 * time.Minute),
		endTime:       staker3.endTime.Add(-1 * time.Minute),
	}

	staker4RewardAddress := ids.GenerateTestShortID()
	staker4 := staker{
		nodeID:        ids.NodeID(staker4RewardAddress),
		rewardAddress: staker4RewardAddress,
		startTime:     staker3.startTime,
		endTime:       staker3.endTime,
	}

	staker5RewardAddress := ids.GenerateTestShortID()
	staker5 := staker{
		nodeID:        ids.NodeID(staker5RewardAddress),
		rewardAddress: staker5RewardAddress,
		startTime:     staker2.endTime,
		endTime:       staker2.endTime.Add(defaultMinStakingDuration),
	}

	tests := []test{
		{
			description:   "advance time to before staker1 start with subnet",
			stakers:       []staker{staker1, staker2, staker3, staker4, staker5},
			subnetStakers: []staker{staker1, staker2, staker3, staker4, staker5},
			advanceTimeTo: []time.Time{staker1.startTime.Add(-1 * time.Second)},
			expectedStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: pending,
				staker2.nodeID: pending,
				staker3.nodeID: pending,
				staker4.nodeID: pending,
				staker5.nodeID: pending,
			},
			expectedSubnetStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: pending,
				staker2.nodeID: pending,
				staker3.nodeID: pending,
				staker4.nodeID: pending,
				staker5.nodeID: pending,
			},
		},
		{
			description:   "advance time to staker 1 start with subnet",
			stakers:       []staker{staker1, staker2, staker3, staker4, staker5},
			subnetStakers: []staker{staker1},
			advanceTimeTo: []time.Time{staker1.startTime},
			expectedStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,
				staker2.nodeID: pending,
				staker3.nodeID: pending,
				staker4.nodeID: pending,
				staker5.nodeID: pending,
			},
			expectedSubnetStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,
				staker2.nodeID: pending,
				staker3.nodeID: pending,
				staker4.nodeID: pending,
				staker5.nodeID: pending,
			},
		},
		{
			description:   "advance time to the staker2 start",
			stakers:       []staker{staker1, staker2, staker3, staker4, staker5},
			advanceTimeTo: []time.Time{staker1.startTime, staker2.startTime},
			expectedStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,
				staker2.nodeID: current,
				staker3.nodeID: pending,
				staker4.nodeID: pending,
				staker5.nodeID: pending,
			},
		},
		{
			description:   "staker3 should validate only primary network",
			stakers:       []staker{staker1, staker2, staker3, staker4, staker5},
			subnetStakers: []staker{staker1, staker2, staker3Sub, staker4, staker5},
			advanceTimeTo: []time.Time{staker1.startTime, staker2.startTime, staker3.startTime},
			expectedStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,
				staker2.nodeID: current,
				staker3.nodeID: current,
				staker4.nodeID: current,
				staker5.nodeID: pending,
			},
			expectedSubnetStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID:    current,
				staker2.nodeID:    current,
				staker3Sub.nodeID: pending,
				staker4.nodeID:    current,
				staker5.nodeID:    pending,
			},
		},
		{
			description:   "advance time to staker3 start with subnet",
			stakers:       []staker{staker1, staker2, staker3, staker4, staker5},
			subnetStakers: []staker{staker1, staker2, staker3Sub, staker4, staker5},
			advanceTimeTo: []time.Time{staker1.startTime, staker2.startTime, staker3.startTime, staker3Sub.startTime},
			expectedStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,
				staker2.nodeID: current,
				staker3.nodeID: current,
				staker4.nodeID: current,
				staker5.nodeID: pending,
			},
			expectedSubnetStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,
				staker2.nodeID: current,
				staker3.nodeID: current,
				staker4.nodeID: current,
				staker5.nodeID: pending,
			},
		},
		{
			description:   "advance time to staker5 end",
			stakers:       []staker{staker1, staker2, staker3, staker4, staker5},
			advanceTimeTo: []time.Time{staker1.startTime, staker2.startTime, staker3.startTime, staker5.startTime},
			expectedStakers: map[ids.NodeID]stakerStatus{
				staker1.nodeID: current,

				// given its txID, staker2 will be
				// rewarded and moved out of current stakers set
				// staker2.nodeID: current,
				staker3.nodeID: current,
				staker4.nodeID: current,
				staker5.nodeID: current,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(ts *testing.T) {
			require := require.New(ts)
			env := newEnvironment(t, nil)
			defer func() {
				if err := shutdownEnvironment(env); err != nil {
					t.Fatal(err)
				}
			}()

			env.config.BanffTime = time.Time{} // activate Banff
			env.config.WhitelistedSubnets.Add(testSubnet1.ID())

			for _, staker := range test.stakers {
				tx, err := env.txBuilder.NewAddValidatorTx(
					env.config.MinValidatorStake,
					uint64(staker.startTime.Unix()),
					uint64(staker.endTime.Unix()),
					staker.nodeID,
					staker.rewardAddress,
					reward.PercentDenominator,
					[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
					ids.ShortEmpty,
				)
				require.NoError(err)

				staker := state.NewPendingStaker(
					tx.ID(),
					tx.Unsigned.(*txs.AddValidatorTx),
				)

				env.state.PutPendingValidator(staker)
				env.state.AddTx(tx, status.Committed)
				require.NoError(env.state.Commit())
			}

			for _, subStaker := range test.subnetStakers {
				tx, err := env.txBuilder.NewAddSubnetValidatorTx(
					10, // Weight
					uint64(subStaker.startTime.Unix()),
					uint64(subStaker.endTime.Unix()),
					subStaker.nodeID, // validator ID
					testSubnet1.ID(), // Subnet ID
					[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
					ids.ShortEmpty,
				)
				require.NoError(err)

				subnetStaker := state.NewPendingStaker(
					tx.ID(),
					tx.Unsigned.(*txs.AddSubnetValidatorTx),
				)

				env.state.PutPendingValidator(subnetStaker)
				env.state.AddTx(tx, status.Committed)
				require.NoError(env.state.Commit())
			}

			for _, newTime := range test.advanceTimeTo {
				env.clk.Set(newTime)

				// add Staker0 (with the right end time) to state
				// so to allow proposalBlk issuance
				staker0.endTime = newTime
				addStaker0, err := env.txBuilder.NewAddValidatorTx(
					10,
					uint64(staker0.startTime.Unix()),
					uint64(staker0.endTime.Unix()),
					staker0.nodeID,
					staker0.rewardAddress,
					reward.PercentDenominator,
					[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
					ids.ShortEmpty,
				)
				require.NoError(err)

				// store Staker0 to state
				staker0 := state.NewCurrentStaker(
					addStaker0.ID(),
					addStaker0.Unsigned.(*txs.AddValidatorTx),
					0,
				)
				env.state.PutCurrentValidator(staker0)
				env.state.AddTx(addStaker0, status.Committed)
				require.NoError(env.state.Commit())

				s0RewardTx := &txs.Tx{
					Unsigned: &txs.RewardValidatorTx{
						TxID: staker0.TxID,
					},
				}
				require.NoError(s0RewardTx.Sign(txs.Codec, nil))

				// build proposal block moving ahead chain time
				// as well as rewarding staker0
				preferredID := env.state.GetLastAccepted()
				parentBlk, _, err := env.state.GetStatelessBlock(preferredID)
				require.NoError(err)
				statelessProposalBlock, err := blocks.NewBanffProposalBlock(
					newTime,
					parentBlk.ID(),
					parentBlk.Height()+1,
					s0RewardTx,
				)
				require.NoError(err)

				// verify and accept the block
				block := env.blkManager.NewBlock(statelessProposalBlock)
				require.NoError(block.Verify())
				options, err := block.(snowman.OracleBlock).Options()
				require.NoError(err)

				require.NoError(options[0].Verify())

				require.NoError(block.Accept())
				require.NoError(options[0].Accept())
			}
			require.NoError(env.state.Commit())

			for stakerNodeID, status := range test.expectedStakers {
				switch status {
				case pending:
					_, err := env.state.GetPendingValidator(constants.PrimaryNetworkID, stakerNodeID)
					require.NoError(err)
					require.False(env.config.Validators.Contains(constants.PrimaryNetworkID, stakerNodeID))
				case current:
					_, err := env.state.GetCurrentValidator(constants.PrimaryNetworkID, stakerNodeID)
					require.NoError(err)
					require.True(env.config.Validators.Contains(constants.PrimaryNetworkID, stakerNodeID))
				}
			}

			for stakerNodeID, status := range test.expectedSubnetStakers {
				switch status {
				case pending:
					require.False(env.config.Validators.Contains(testSubnet1.ID(), stakerNodeID))
				case current:
					require.True(env.config.Validators.Contains(testSubnet1.ID(), stakerNodeID))
				}
			}
		})
	}
}

func TestBanffProposalBlockRemoveSubnetValidator(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, nil)
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()
	env.config.BanffTime = time.Time{} // activate Banff
	env.config.WhitelistedSubnets.Add(testSubnet1.ID())

	// Add a subnet validator to the staker set
	subnetValidatorNodeID := ids.NodeID(preFundedKeys[0].PublicKey().Address())
	// Starts after the corre
	subnetVdr1StartTime := defaultValidateStartTime
	subnetVdr1EndTime := defaultValidateStartTime.Add(defaultMinStakingDuration)
	tx, err := env.txBuilder.NewAddSubnetValidatorTx(
		1,                                  // Weight
		uint64(subnetVdr1StartTime.Unix()), // Start time
		uint64(subnetVdr1EndTime.Unix()),   // end time
		subnetValidatorNodeID,              // Node ID
		testSubnet1.ID(),                   // Subnet ID
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	staker := state.NewCurrentStaker(
		tx.ID(),
		tx.Unsigned.(*txs.AddSubnetValidatorTx),
		0,
	)

	env.state.PutCurrentValidator(staker)
	env.state.AddTx(tx, status.Committed)
	require.NoError(env.state.Commit())

	// The above validator is now part of the staking set

	// Queue a staker that joins the staker set after the above validator leaves
	subnetVdr2NodeID := ids.NodeID(preFundedKeys[1].PublicKey().Address())
	tx, err = env.txBuilder.NewAddSubnetValidatorTx(
		1, // Weight
		uint64(subnetVdr1EndTime.Add(time.Second).Unix()),                                // Start time
		uint64(subnetVdr1EndTime.Add(time.Second).Add(defaultMinStakingDuration).Unix()), // end time
		subnetVdr2NodeID, // Node ID
		testSubnet1.ID(), // Subnet ID
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	staker = state.NewPendingStaker(
		tx.ID(),
		tx.Unsigned.(*txs.AddSubnetValidatorTx),
	)

	env.state.PutPendingValidator(staker)
	env.state.AddTx(tx, status.Committed)
	require.NoError(env.state.Commit())

	// The above validator is now in the pending staker set

	// Advance time to the first staker's end time.
	env.clk.Set(subnetVdr1EndTime)

	// add Staker0 (with the right end time) to state
	// so to allow proposalBlk issuance
	staker0StartTime := defaultValidateStartTime
	staker0EndTime := subnetVdr1EndTime
	addStaker0, err := env.txBuilder.NewAddValidatorTx(
		10,
		uint64(staker0StartTime.Unix()),
		uint64(staker0EndTime.Unix()),
		ids.GenerateTestNodeID(),
		ids.GenerateTestShortID(),
		reward.PercentDenominator,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	// store Staker0 to state
	staker = state.NewCurrentStaker(
		addStaker0.ID(),
		addStaker0.Unsigned.(*txs.AddValidatorTx),
		0,
	)
	env.state.PutCurrentValidator(staker)
	env.state.AddTx(addStaker0, status.Committed)
	require.NoError(env.state.Commit())

	// create rewardTx for staker0
	s0RewardTx := &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: addStaker0.ID(),
		},
	}
	require.NoError(s0RewardTx.Sign(txs.Codec, nil))

	// build proposal block moving ahead chain time
	preferredID := env.state.GetLastAccepted()
	parentBlk, _, err := env.state.GetStatelessBlock(preferredID)
	require.NoError(err)
	statelessProposalBlock, err := blocks.NewBanffProposalBlock(
		subnetVdr1EndTime,
		parentBlk.ID(),
		parentBlk.Height()+1,
		s0RewardTx,
	)
	require.NoError(err)
	propBlk := env.blkManager.NewBlock(statelessProposalBlock)
	require.NoError(propBlk.Verify()) // verify and update staker set

	options, err := propBlk.(snowman.OracleBlock).Options()
	require.NoError(err)
	commitBlk := options[0]
	require.NoError(commitBlk.Verify())

	blkStateMap := env.blkManager.(*manager).blkIDToState
	updatedState := blkStateMap[commitBlk.ID()].onAcceptState
	_, err = updatedState.GetCurrentValidator(testSubnet1.ID(), subnetValidatorNodeID)
	require.ErrorIs(err, database.ErrNotFound)

	// Check VM Validators are removed successfully
	require.NoError(propBlk.Accept())
	require.NoError(commitBlk.Accept())
	require.False(env.config.Validators.Contains(testSubnet1.ID(), subnetVdr2NodeID))
	require.False(env.config.Validators.Contains(testSubnet1.ID(), subnetValidatorNodeID))
}

func TestBanffProposalBlockWhitelistedSubnet(t *testing.T) {
	require := require.New(t)

	for _, whitelist := range []bool{true, false} {
		t.Run(fmt.Sprintf("whitelisted %t", whitelist), func(ts *testing.T) {
			env := newEnvironment(t, nil)
			defer func() {
				if err := shutdownEnvironment(env); err != nil {
					t.Fatal(err)
				}
			}()
			env.config.BanffTime = time.Time{} // activate Banff
			if whitelist {
				env.config.WhitelistedSubnets.Add(testSubnet1.ID())
			}

			// Add a subnet validator to the staker set
			subnetValidatorNodeID := ids.NodeID(preFundedKeys[0].PublicKey().Address())

			subnetVdr1StartTime := defaultGenesisTime.Add(1 * time.Minute)
			subnetVdr1EndTime := defaultGenesisTime.Add(10 * defaultMinStakingDuration).Add(1 * time.Minute)
			tx, err := env.txBuilder.NewAddSubnetValidatorTx(
				1,                                  // Weight
				uint64(subnetVdr1StartTime.Unix()), // Start time
				uint64(subnetVdr1EndTime.Unix()),   // end time
				subnetValidatorNodeID,              // Node ID
				testSubnet1.ID(),                   // Subnet ID
				[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
				ids.ShortEmpty,
			)
			require.NoError(err)

			staker := state.NewPendingStaker(
				tx.ID(),
				tx.Unsigned.(*txs.AddSubnetValidatorTx),
			)

			env.state.PutPendingValidator(staker)
			env.state.AddTx(tx, status.Committed)
			require.NoError(env.state.Commit())

			// Advance time to the staker's start time.
			env.clk.Set(subnetVdr1StartTime)

			// add Staker0 (with the right end time) to state
			// so to allow proposalBlk issuance
			staker0StartTime := defaultGenesisTime
			staker0EndTime := subnetVdr1StartTime
			addStaker0, err := env.txBuilder.NewAddValidatorTx(
				10,
				uint64(staker0StartTime.Unix()),
				uint64(staker0EndTime.Unix()),
				ids.GenerateTestNodeID(),
				ids.GenerateTestShortID(),
				reward.PercentDenominator,
				[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
				ids.ShortEmpty,
			)
			require.NoError(err)

			// store Staker0 to state
			staker = state.NewCurrentStaker(
				addStaker0.ID(),
				addStaker0.Unsigned.(*txs.AddValidatorTx),
				0,
			)
			env.state.PutCurrentValidator(staker)
			env.state.AddTx(addStaker0, status.Committed)
			require.NoError(env.state.Commit())

			// create rewardTx for staker0
			s0RewardTx := &txs.Tx{
				Unsigned: &txs.RewardValidatorTx{
					TxID: addStaker0.ID(),
				},
			}
			require.NoError(s0RewardTx.Sign(txs.Codec, nil))

			// build proposal block moving ahead chain time
			preferredID := env.state.GetLastAccepted()
			parentBlk, _, err := env.state.GetStatelessBlock(preferredID)
			require.NoError(err)
			statelessProposalBlock, err := blocks.NewBanffProposalBlock(
				subnetVdr1StartTime,
				parentBlk.ID(),
				parentBlk.Height()+1,
				s0RewardTx,
			)
			require.NoError(err)
			propBlk := env.blkManager.NewBlock(statelessProposalBlock)
			require.NoError(propBlk.Verify()) // verify update staker set
			options, err := propBlk.(snowman.OracleBlock).Options()
			require.NoError(err)
			commitBlk := options[0]
			require.NoError(commitBlk.Verify())

			require.NoError(propBlk.Accept())
			require.NoError(commitBlk.Accept())
			require.Equal(whitelist, env.config.Validators.Contains(testSubnet1.ID(), subnetValidatorNodeID))
		})
	}
}

func TestBanffProposalBlockDelegatorStakerWeight(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, nil)
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()
	env.config.BanffTime = time.Time{} // activate Banff

	// Case: Timestamp is after next validator start time
	// Add a pending validator
	pendingValidatorStartTime := defaultGenesisTime.Add(1 * time.Second)
	pendingValidatorEndTime := pendingValidatorStartTime.Add(defaultMaxStakingDuration)
	nodeID := ids.GenerateTestNodeID()
	rewardAddress := ids.GenerateTestShortID()
	_, err := addPendingValidator(
		env,
		pendingValidatorStartTime,
		pendingValidatorEndTime,
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
	)
	require.NoError(err)

	// add Staker0 (with the right end time) to state
	// just to allow proposalBlk issuance (with a reward Tx)
	staker0StartTime := defaultGenesisTime
	staker0EndTime := pendingValidatorStartTime
	addStaker0, err := env.txBuilder.NewAddValidatorTx(
		10,
		uint64(staker0StartTime.Unix()),
		uint64(staker0EndTime.Unix()),
		ids.GenerateTestNodeID(),
		ids.GenerateTestShortID(),
		reward.PercentDenominator,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	// store Staker0 to state
	staker := state.NewCurrentStaker(
		addStaker0.ID(),
		addStaker0.Unsigned.(*txs.AddValidatorTx),
		0,
	)
	env.state.PutCurrentValidator(staker)
	env.state.AddTx(addStaker0, status.Committed)
	require.NoError(env.state.Commit())

	// create rewardTx for staker0
	s0RewardTx := &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: addStaker0.ID(),
		},
	}
	require.NoError(s0RewardTx.Sign(txs.Codec, nil))

	// build proposal block moving ahead chain time
	preferredID := env.state.GetLastAccepted()
	parentBlk, _, err := env.state.GetStatelessBlock(preferredID)
	require.NoError(err)
	statelessProposalBlock, err := blocks.NewBanffProposalBlock(
		pendingValidatorStartTime,
		parentBlk.ID(),
		parentBlk.Height()+1,
		s0RewardTx,
	)
	require.NoError(err)
	propBlk := env.blkManager.NewBlock(statelessProposalBlock)
	require.NoError(propBlk.Verify())

	options, err := propBlk.(snowman.OracleBlock).Options()
	require.NoError(err)
	commitBlk := options[0]
	require.NoError(commitBlk.Verify())

	require.NoError(propBlk.Accept())
	require.NoError(commitBlk.Accept())

	// Test validator weight before delegation
	primarySet, ok := env.config.Validators.GetValidators(constants.PrimaryNetworkID)
	require.True(ok)
	vdrWeight, _ := primarySet.GetWeight(nodeID)
	require.Equal(env.config.MinValidatorStake, vdrWeight)

	// Add delegator
	pendingDelegatorStartTime := pendingValidatorStartTime.Add(1 * time.Second)
	pendingDelegatorEndTime := pendingDelegatorStartTime.Add(1 * time.Second)

	addDelegatorTx, err := env.txBuilder.NewAddDelegatorTx(
		env.config.MinDelegatorStake,
		uint64(pendingDelegatorStartTime.Unix()),
		uint64(pendingDelegatorEndTime.Unix()),
		nodeID,
		preFundedKeys[0].PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{
			preFundedKeys[0],
			preFundedKeys[1],
			preFundedKeys[4],
		},
		ids.ShortEmpty,
	)
	require.NoError(err)

	staker = state.NewPendingStaker(
		addDelegatorTx.ID(),
		addDelegatorTx.Unsigned.(*txs.AddDelegatorTx),
	)

	env.state.PutPendingDelegator(staker)
	env.state.AddTx(addDelegatorTx, status.Committed)
	env.state.SetHeight( /*dummyHeight*/ uint64(1))
	require.NoError(env.state.Commit())

	// add Staker0 (with the right end time) to state
	// so to allow proposalBlk issuance
	staker0EndTime = pendingDelegatorStartTime
	addStaker0, err = env.txBuilder.NewAddValidatorTx(
		10,
		uint64(staker0StartTime.Unix()),
		uint64(staker0EndTime.Unix()),
		ids.GenerateTestNodeID(),
		ids.GenerateTestShortID(),
		reward.PercentDenominator,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	// store Staker0 to state
	staker = state.NewCurrentStaker(
		addStaker0.ID(),
		addStaker0.Unsigned.(*txs.AddValidatorTx),
		0,
	)
	env.state.PutCurrentValidator(staker)
	env.state.AddTx(addStaker0, status.Committed)
	require.NoError(env.state.Commit())

	// create rewardTx for staker0
	s0RewardTx = &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: addStaker0.ID(),
		},
	}
	require.NoError(s0RewardTx.Sign(txs.Codec, nil))

	// Advance Time
	preferredID = env.state.GetLastAccepted()
	parentBlk, _, err = env.state.GetStatelessBlock(preferredID)
	require.NoError(err)
	statelessProposalBlock, err = blocks.NewBanffProposalBlock(
		pendingDelegatorStartTime,
		parentBlk.ID(),
		parentBlk.Height()+1,
		s0RewardTx,
	)
	require.NoError(err)

	propBlk = env.blkManager.NewBlock(statelessProposalBlock)
	require.NoError(propBlk.Verify())

	options, err = propBlk.(snowman.OracleBlock).Options()
	require.NoError(err)
	commitBlk = options[0]
	require.NoError(commitBlk.Verify())

	require.NoError(propBlk.Accept())
	require.NoError(commitBlk.Accept())

	// Test validator weight after delegation
	vdrWeight, _ = primarySet.GetWeight(nodeID)
	require.Equal(env.config.MinDelegatorStake+env.config.MinValidatorStake, vdrWeight)
}

func TestBanffProposalBlockDelegatorStakers(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, nil)
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()
	env.config.BanffTime = time.Time{} // activate Banff

	// Case: Timestamp is after next validator start time
	// Add a pending validator
	pendingValidatorStartTime := defaultGenesisTime.Add(1 * time.Second)
	pendingValidatorEndTime := pendingValidatorStartTime.Add(defaultMinStakingDuration)
	factory := crypto.FactorySECP256K1R{}
	nodeIDKey, _ := factory.NewPrivateKey()
	rewardAddress := nodeIDKey.PublicKey().Address()
	nodeID := ids.NodeID(rewardAddress)

	_, err := addPendingValidator(
		env,
		pendingValidatorStartTime,
		pendingValidatorEndTime,
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
	)
	require.NoError(err)

	// add Staker0 (with the right end time) to state
	// so to allow proposalBlk issuance
	staker0StartTime := defaultGenesisTime
	staker0EndTime := pendingValidatorStartTime
	addStaker0, err := env.txBuilder.NewAddValidatorTx(
		10,
		uint64(staker0StartTime.Unix()),
		uint64(staker0EndTime.Unix()),
		ids.GenerateTestNodeID(),
		ids.GenerateTestShortID(),
		reward.PercentDenominator,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	// store Staker0 to state
	staker := state.NewCurrentStaker(
		addStaker0.ID(),
		addStaker0.Unsigned.(*txs.AddValidatorTx),
		0,
	)
	env.state.PutCurrentValidator(staker)
	env.state.AddTx(addStaker0, status.Committed)
	require.NoError(env.state.Commit())

	// create rewardTx for staker0
	s0RewardTx := &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: addStaker0.ID(),
		},
	}
	require.NoError(s0RewardTx.Sign(txs.Codec, nil))

	// build proposal block moving ahead chain time
	preferredID := env.state.GetLastAccepted()
	parentBlk, _, err := env.state.GetStatelessBlock(preferredID)
	require.NoError(err)
	statelessProposalBlock, err := blocks.NewBanffProposalBlock(
		pendingValidatorStartTime,
		parentBlk.ID(),
		parentBlk.Height()+1,
		s0RewardTx,
	)
	require.NoError(err)
	propBlk := env.blkManager.NewBlock(statelessProposalBlock)
	require.NoError(propBlk.Verify())

	options, err := propBlk.(snowman.OracleBlock).Options()
	require.NoError(err)
	commitBlk := options[0]
	require.NoError(commitBlk.Verify())

	require.NoError(propBlk.Accept())
	require.NoError(commitBlk.Accept())

	// Test validator weight before delegation
	primarySet, ok := env.config.Validators.GetValidators(constants.PrimaryNetworkID)
	require.True(ok)
	vdrWeight, _ := primarySet.GetWeight(nodeID)
	require.Equal(env.config.MinValidatorStake, vdrWeight)

	// Add delegator
	pendingDelegatorStartTime := pendingValidatorStartTime.Add(1 * time.Second)
	pendingDelegatorEndTime := pendingDelegatorStartTime.Add(defaultMinStakingDuration)
	addDelegatorTx, err := env.txBuilder.NewAddDelegatorTx(
		env.config.MinDelegatorStake,
		uint64(pendingDelegatorStartTime.Unix()),
		uint64(pendingDelegatorEndTime.Unix()),
		nodeID,
		preFundedKeys[0].PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{
			preFundedKeys[0],
			preFundedKeys[1],
			preFundedKeys[4],
		},
		ids.ShortEmpty,
	)
	require.NoError(err)

	staker = state.NewPendingStaker(
		addDelegatorTx.ID(),
		addDelegatorTx.Unsigned.(*txs.AddDelegatorTx),
	)

	env.state.PutPendingDelegator(staker)
	env.state.AddTx(addDelegatorTx, status.Committed)
	env.state.SetHeight( /*dummyHeight*/ uint64(1))
	require.NoError(env.state.Commit())

	// add Staker0 (with the right end time) to state
	// so to allow proposalBlk issuance
	staker0EndTime = pendingDelegatorStartTime
	addStaker0, err = env.txBuilder.NewAddValidatorTx(
		10,
		uint64(staker0StartTime.Unix()),
		uint64(staker0EndTime.Unix()),
		ids.GenerateTestNodeID(),
		ids.GenerateTestShortID(),
		reward.PercentDenominator,
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0], preFundedKeys[1]},
		ids.ShortEmpty,
	)
	require.NoError(err)

	// store Staker0 to state
	staker = state.NewCurrentStaker(
		addStaker0.ID(),
		addStaker0.Unsigned.(*txs.AddValidatorTx),
		0,
	)
	env.state.PutCurrentValidator(staker)
	env.state.AddTx(addStaker0, status.Committed)
	require.NoError(env.state.Commit())

	// create rewardTx for staker0
	s0RewardTx = &txs.Tx{
		Unsigned: &txs.RewardValidatorTx{
			TxID: addStaker0.ID(),
		},
	}
	require.NoError(s0RewardTx.Sign(txs.Codec, nil))

	// Advance Time
	preferredID = env.state.GetLastAccepted()
	parentBlk, _, err = env.state.GetStatelessBlock(preferredID)
	require.NoError(err)
	statelessProposalBlock, err = blocks.NewBanffProposalBlock(
		pendingDelegatorStartTime,
		parentBlk.ID(),
		parentBlk.Height()+1,
		s0RewardTx,
	)
	require.NoError(err)
	propBlk = env.blkManager.NewBlock(statelessProposalBlock)
	require.NoError(propBlk.Verify())

	options, err = propBlk.(snowman.OracleBlock).Options()
	require.NoError(err)
	commitBlk = options[0]
	require.NoError(commitBlk.Verify())

	require.NoError(propBlk.Accept())
	require.NoError(commitBlk.Accept())

	// Test validator weight after delegation
	vdrWeight, _ = primarySet.GetWeight(nodeID)
	require.Equal(env.config.MinDelegatorStake+env.config.MinValidatorStake, vdrWeight)
}
