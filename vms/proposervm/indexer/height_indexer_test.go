// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package indexer

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/memdb"
	"github.com/skychains/chain/database/versiondb"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/choices"
	"github.com/skychains/chain/snow/consensus/snowman"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms/proposervm/block"
	"github.com/skychains/chain/vms/proposervm/state"
)

func TestHeightBlockIndexPostFork(t *testing.T) {
	require := require.New(t)

	db := memdb.New()
	vdb := versiondb.New(db)
	storedState := state.New(vdb)

	// Build a chain of post fork blocks
	var (
		blkNumber = uint64(10)
		lastBlkID = ids.Empty.Prefix(0) // initially set to a dummyGenesisID
		proBlks   = make(map[ids.ID]snowman.Block)
	)

	for blkHeight := uint64(1); blkHeight <= blkNumber; blkHeight++ {
		blockBytes := ids.Empty.Prefix(blkHeight + blkNumber + 1)
		dummyTS := time.Time{}
		dummyPCH := uint64(2022)

		// store postForkStatelessBlk in State ...
		postForkStatelessBlk, err := block.BuildUnsigned(
			lastBlkID,
			dummyTS,
			dummyPCH,
			blockBytes[:],
		)
		require.NoError(err)
		require.NoError(storedState.PutBlock(postForkStatelessBlk, choices.Accepted))

		// ... and create a corresponding test block just for block server
		postForkBlk := &snowman.TestBlock{
			TestDecidable: choices.TestDecidable{
				IDV:     postForkStatelessBlk.ID(),
				StatusV: choices.Accepted,
			},
			HeightV: blkHeight,
		}
		proBlks[postForkBlk.ID()] = postForkBlk

		lastBlkID = postForkStatelessBlk.ID()
	}

	blkSrv := &TestBlockServer{
		CantGetFullPostForkBlock: true,
		CantCommit:               true,

		GetFullPostForkBlockF: func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
			blk, found := proBlks[blkID]
			if !found {
				return nil, database.ErrNotFound
			}
			return blk, nil
		},
		CommitF: func() error {
			return nil
		},
	}

	hIndex := newHeightIndexer(blkSrv,
		logging.NoLog{},
		storedState,
	)
	hIndex.commitFrequency = 0 // commit each block

	// checkpoint last accepted block and show the whole chain in reindexed
	require.NoError(hIndex.state.SetCheckpoint(lastBlkID))
	require.NoError(hIndex.RepairHeightIndex(context.Background()))
	require.True(hIndex.IsRepaired())

	// check that height index is fully built
	loadedForkHeight, err := storedState.GetForkHeight()
	require.NoError(err)
	require.Equal(uint64(1), loadedForkHeight)
	for height := uint64(1); height <= blkNumber; height++ {
		_, err := storedState.GetBlockIDAtHeight(height)
		require.NoError(err)
	}
}

func TestHeightBlockIndexAcrossFork(t *testing.T) {
	require := require.New(t)

	db := memdb.New()
	vdb := versiondb.New(db)
	storedState := state.New(vdb)

	// Build a chain of post fork blocks
	var (
		blkNumber  = uint64(10)
		forkHeight = blkNumber / 2
		lastBlkID  = ids.Empty.Prefix(0) // initially set to a last pre fork blk
		proBlks    = make(map[ids.ID]snowman.Block)
	)

	for blkHeight := forkHeight; blkHeight <= blkNumber; blkHeight++ {
		blockBytes := ids.Empty.Prefix(blkHeight + blkNumber + 1)
		dummyTS := time.Time{}
		dummyPCH := uint64(2022)

		// store postForkStatelessBlk in State ...
		postForkStatelessBlk, err := block.BuildUnsigned(
			lastBlkID,
			dummyTS,
			dummyPCH,
			blockBytes[:],
		)
		require.NoError(err)
		require.NoError(storedState.PutBlock(postForkStatelessBlk, choices.Accepted))

		// ... and create a corresponding test block just for block server
		postForkBlk := &snowman.TestBlock{
			TestDecidable: choices.TestDecidable{
				IDV:     postForkStatelessBlk.ID(),
				StatusV: choices.Accepted,
			},
			HeightV: blkHeight,
		}
		proBlks[postForkBlk.ID()] = postForkBlk

		lastBlkID = postForkStatelessBlk.ID()
	}

	blkSrv := &TestBlockServer{
		CantGetFullPostForkBlock: true,
		CantCommit:               true,

		GetFullPostForkBlockF: func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
			blk, found := proBlks[blkID]
			if !found {
				return nil, database.ErrNotFound
			}
			return blk, nil
		},
		CommitF: func() error {
			return nil
		},
	}

	hIndex := newHeightIndexer(blkSrv,
		logging.NoLog{},
		storedState,
	)
	hIndex.commitFrequency = 0 // commit each block

	// checkpoint last accepted block and show the whole chain in reindexed
	require.NoError(hIndex.state.SetCheckpoint(lastBlkID))
	require.NoError(hIndex.RepairHeightIndex(context.Background()))
	require.True(hIndex.IsRepaired())

	// check that height index is fully built
	loadedForkHeight, err := storedState.GetForkHeight()
	require.NoError(err)
	require.Equal(forkHeight, loadedForkHeight)
	for height := uint64(0); height < forkHeight; height++ {
		_, err := storedState.GetBlockIDAtHeight(height)
		require.ErrorIs(err, database.ErrNotFound)
	}
	for height := forkHeight; height <= blkNumber; height++ {
		_, err := storedState.GetBlockIDAtHeight(height)
		require.NoError(err)
	}
}

func TestHeightBlockIndexResumeFromCheckPoint(t *testing.T) {
	require := require.New(t)

	db := memdb.New()
	vdb := versiondb.New(db)
	storedState := state.New(vdb)

	// Build a chain of post fork blocks
	var (
		blkNumber  = uint64(10)
		forkHeight = blkNumber / 2
		lastBlkID  = ids.Empty.Prefix(0) // initially set to a last pre fork blk
		proBlks    = make(map[ids.ID]snowman.Block)
	)

	for blkHeight := forkHeight; blkHeight <= blkNumber; blkHeight++ {
		blockBytes := ids.Empty.Prefix(blkHeight + blkNumber + 1)
		dummyTS := time.Time{}
		dummyPCH := uint64(2022)

		// store postForkStatelessBlk in State ...
		postForkStatelessBlk, err := block.BuildUnsigned(
			lastBlkID,
			dummyTS,
			dummyPCH,
			blockBytes[:],
		)
		require.NoError(err)
		require.NoError(storedState.PutBlock(postForkStatelessBlk, choices.Accepted))

		// ... and create a corresponding test block just for block server
		postForkBlk := &snowman.TestBlock{
			TestDecidable: choices.TestDecidable{
				IDV:     postForkStatelessBlk.ID(),
				StatusV: choices.Accepted,
			},
			HeightV: blkHeight,
		}
		proBlks[postForkBlk.ID()] = postForkBlk

		lastBlkID = postForkStatelessBlk.ID()
	}

	blkSrv := &TestBlockServer{
		CantGetFullPostForkBlock: true,
		CantCommit:               true,

		GetFullPostForkBlockF: func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
			blk, found := proBlks[blkID]
			if !found {
				return nil, database.ErrNotFound
			}
			return blk, nil
		},
		CommitF: func() error {
			return nil
		},
	}

	hIndex := newHeightIndexer(blkSrv,
		logging.NoLog{},
		storedState,
	)
	hIndex.commitFrequency = 0 // commit each block

	// pick a random block in the chain and checkpoint it;...
	rndPostForkHeight := rand.Intn(int(blkNumber-forkHeight)) + int(forkHeight) // #nosec G404
	var checkpointBlk snowman.Block
	for _, blk := range proBlks {
		if blk.Height() != uint64(rndPostForkHeight) {
			continue // not the blk we are looking for
		}

		checkpointBlk = blk
		require.NoError(hIndex.state.SetCheckpoint(checkpointBlk.ID()))
		break
	}

	// perform repair and show index is built
	require.NoError(hIndex.RepairHeightIndex(context.Background()))
	require.True(hIndex.IsRepaired())

	// check that height index is fully built
	loadedForkHeight, err := storedState.GetForkHeight()
	require.NoError(err)
	require.Equal(forkHeight, loadedForkHeight)
	for height := forkHeight; height <= checkpointBlk.Height(); height++ {
		_, err := storedState.GetBlockIDAtHeight(height)
		require.NoError(err)
	}
}
