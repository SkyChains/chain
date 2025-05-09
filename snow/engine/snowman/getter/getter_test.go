// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package getter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/consensus/snowman"
	"github.com/skychains/chain/snow/consensus/snowman/snowmantest"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/engine/snowman/block"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/set"
)

var errUnknownBlock = errors.New("unknown block")

type StateSyncEnabledMock struct {
	*block.TestVM
	*block.MockStateSyncableVM
}

func newTest(t *testing.T) (common.AllGetsServer, StateSyncEnabledMock, *common.SenderTest) {
	ctrl := gomock.NewController(t)

	vm := StateSyncEnabledMock{
		TestVM:              &block.TestVM{},
		MockStateSyncableVM: block.NewMockStateSyncableVM(ctrl),
	}

	sender := &common.SenderTest{
		T: t,
	}
	sender.Default(true)

	bs, err := New(
		vm,
		sender,
		logging.NoLog{},
		time.Second,
		2000,
		prometheus.NewRegistry(),
	)
	require.NoError(t, err)

	return bs, vm, sender
}

func TestAcceptedFrontier(t *testing.T) {
	require := require.New(t)
	bs, vm, sender := newTest(t)

	blkID := ids.GenerateTestID()
	vm.LastAcceptedF = func(context.Context) (ids.ID, error) {
		return blkID, nil
	}

	var accepted ids.ID
	sender.SendAcceptedFrontierF = func(_ context.Context, _ ids.NodeID, _ uint32, containerID ids.ID) {
		accepted = containerID
	}

	require.NoError(bs.GetAcceptedFrontier(context.Background(), ids.EmptyNodeID, 0))
	require.Equal(blkID, accepted)
}

func TestFilterAccepted(t *testing.T) {
	require := require.New(t)
	bs, vm, sender := newTest(t)

	acceptedBlk := snowmantest.BuildChild(snowmantest.Genesis)
	require.NoError(acceptedBlk.Accept(context.Background()))

	unknownBlkID := ids.GenerateTestID()

	vm.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		switch blkID {
		case snowmantest.GenesisID:
			return snowmantest.Genesis, nil
		case acceptedBlk.ID():
			return acceptedBlk, nil
		case unknownBlkID:
			return nil, errUnknownBlock
		default:
			require.FailNow(errUnknownBlock.Error())
			return nil, errUnknownBlock
		}
	}

	var accepted []ids.ID
	sender.SendAcceptedF = func(_ context.Context, _ ids.NodeID, _ uint32, frontier []ids.ID) {
		accepted = frontier
	}

	blkIDs := set.Of(snowmantest.GenesisID, acceptedBlk.ID(), unknownBlkID)
	require.NoError(bs.GetAccepted(context.Background(), ids.EmptyNodeID, 0, blkIDs))

	require.Len(accepted, 2)
	require.Contains(accepted, snowmantest.GenesisID)
	require.Contains(accepted, acceptedBlk.ID())
	require.NotContains(accepted, unknownBlkID)
}
