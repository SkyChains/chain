// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package syncer

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/engine/common/tracker"
	"github.com/skychains/chain/snow/engine/snowman/block"
	"github.com/skychains/chain/snow/engine/snowman/getter"
	"github.com/skychains/chain/snow/validators"
	"github.com/skychains/chain/utils/hashing"
)

const (
	key         uint64 = 2022
	minorityKey uint64 = 2000
)

var (
	_ block.ChainVM         = fullVM{}
	_ block.StateSyncableVM = fullVM{}

	unknownSummaryID = ids.ID{'g', 'a', 'r', 'b', 'a', 'g', 'e'}

	summaryBytes = []byte{'s', 'u', 'm', 'm', 'a', 'r', 'y'}
	summaryID    ids.ID

	minoritySummaryBytes = []byte{'m', 'i', 'n', 'o', 'r', 'i', 't', 'y'}
	minoritySummaryID    ids.ID
)

func init() {
	var err error
	summaryID, err = ids.ToID(hashing.ComputeHash256(summaryBytes))
	if err != nil {
		panic(err)
	}

	minoritySummaryID, err = ids.ToID(hashing.ComputeHash256(minoritySummaryBytes))
	if err != nil {
		panic(err)
	}
}

type fullVM struct {
	*block.TestVM
	*block.TestStateSyncableVM
}

func buildTestPeers(t *testing.T, subnetID ids.ID) validators.Manager {
	// We consider more than maxOutstandingBroadcastRequests peers to test
	// capping the number of requests sent out.
	vdrs := validators.NewManager()
	for idx := 0; idx < 2*maxOutstandingBroadcastRequests; idx++ {
		beaconID := ids.GenerateTestNodeID()
		require.NoError(t, vdrs.AddStaker(subnetID, beaconID, nil, ids.Empty, 1))
	}
	return vdrs
}

func buildTestsObjects(
	t *testing.T,
	ctx *snow.ConsensusContext,
	startupTracker tracker.Startup,
	beacons validators.Manager,
	alpha uint64,
) (
	*stateSyncer,
	*fullVM,
	*common.SenderTest,
) {
	require := require.New(t)

	fullVM := &fullVM{
		TestVM: &block.TestVM{
			TestVM: common.TestVM{T: t},
		},
		TestStateSyncableVM: &block.TestStateSyncableVM{
			T: t,
		},
	}
	sender := &common.SenderTest{T: t}
	dummyGetter, err := getter.New(
		fullVM,
		sender,
		ctx.Log,
		time.Second,
		2000,
		ctx.Registerer,
	)
	require.NoError(err)

	cfg, err := NewConfig(
		dummyGetter,
		ctx,
		startupTracker,
		sender,
		beacons,
		beacons.Count(ctx.SubnetID),
		alpha,
		nil,
		fullVM,
	)
	require.NoError(err)
	commonSyncer := New(cfg, func(context.Context, uint32) error {
		return nil
	})
	require.IsType(&stateSyncer{}, commonSyncer)
	syncer := commonSyncer.(*stateSyncer)
	require.NotNil(syncer.stateSyncVM)

	fullVM.GetOngoingSyncStateSummaryF = func(context.Context) (block.StateSummary, error) {
		return nil, database.ErrNotFound
	}

	return syncer, fullVM, sender
}

func pickRandomFrom(nodes map[ids.NodeID]uint32) ids.NodeID {
	for node := range nodes {
		return node
	}
	return ids.EmptyNodeID
}
