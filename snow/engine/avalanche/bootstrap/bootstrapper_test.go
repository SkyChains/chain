// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package bootstrap

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database/memdb"
	"github.com/skychains/chain/database/prefixdb"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/network/p2p"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/choices"
	"github.com/skychains/chain/snow/consensus/lux"
	"github.com/skychains/chain/snow/consensus/snowstorm"
	"github.com/skychains/chain/snow/engine/lux/bootstrap/queue"
	"github.com/skychains/chain/snow/engine/lux/getter"
	"github.com/skychains/chain/snow/engine/lux/vertex"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/engine/common/tracker"
	"github.com/skychains/chain/snow/snowtest"
	"github.com/skychains/chain/snow/validators"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/version"

	p2ppb "github.com/skychains/chain/proto/pb/p2p"
)

var (
	errUnknownVertex       = errors.New("unknown vertex")
	errParsedUnknownVertex = errors.New("parsed unknown vertex")
	errUnknownTx           = errors.New("unknown tx")
)

type testTx struct {
	snowstorm.Tx

	tx *snowstorm.TestTx
}

func (t *testTx) Accept(ctx context.Context) error {
	if err := t.Tx.Accept(ctx); err != nil {
		return err
	}
	t.tx.DependenciesV = nil
	return nil
}

func newConfig(t *testing.T) (Config, ids.NodeID, *common.SenderTest, *vertex.TestManager, *vertex.TestVM) {
	require := require.New(t)

	snowCtx := snowtest.Context(t, snowtest.CChainID)
	ctx := snowtest.ConsensusContext(snowCtx)

	vdrs := validators.NewManager()
	db := memdb.New()
	sender := &common.SenderTest{T: t}
	manager := vertex.NewTestManager(t)
	vm := &vertex.TestVM{}
	vm.T = t

	sender.Default(true)
	manager.Default(true)
	vm.Default(true)

	peer := ids.GenerateTestNodeID()
	require.NoError(vdrs.AddStaker(constants.PrimaryNetworkID, peer, nil, ids.Empty, 1))

	vtxBlocker, err := queue.NewWithMissing(prefixdb.New([]byte("vtx"), db), "vtx", prometheus.NewRegistry())
	require.NoError(err)

	txBlocker, err := queue.New(prefixdb.New([]byte("tx"), db), "tx", prometheus.NewRegistry())
	require.NoError(err)

	peerTracker := tracker.NewPeers()
	totalWeight, err := vdrs.TotalWeight(constants.PrimaryNetworkID)
	require.NoError(err)
	startupTracker := tracker.NewStartup(peerTracker, totalWeight/2+1)
	vdrs.RegisterSetCallbackListener(constants.PrimaryNetworkID, startupTracker)

	avaGetHandler, err := getter.New(manager, sender, ctx.Log, time.Second, 2000, prometheus.NewRegistry())
	require.NoError(err)

	p2pTracker, err := p2p.NewPeerTracker(
		logging.NoLog{},
		"",
		prometheus.NewRegistry(),
		nil,
		version.CurrentApp,
	)
	require.NoError(err)

	p2pTracker.Connected(peer, version.CurrentApp)

	return Config{
		AllGetsServer:                  avaGetHandler,
		Ctx:                            ctx,
		StartupTracker:                 startupTracker,
		Sender:                         sender,
		PeerTracker:                    p2pTracker,
		AncestorsMaxContainersReceived: 2000,
		VtxBlocked:                     vtxBlocker,
		TxBlocked:                      txBlocker,
		Manager:                        manager,
		VM:                             vm,
	}, peer, sender, manager, vm
}

// Three vertices in the accepted frontier. None have parents. No need to fetch
// anything
func TestBootstrapperSingleFrontier(t *testing.T) {
	require := require.New(t)

	config, _, _, manager, vm := newConfig(t)

	vtxID0 := ids.Empty.Prefix(0)
	vtxID1 := ids.Empty.Prefix(1)
	vtxID2 := ids.Empty.Prefix(2)

	vtxBytes0 := []byte{0}
	vtxBytes1 := []byte{1}
	vtxBytes2 := []byte{2}

	vtx0 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID0,
			StatusV: choices.Processing,
		},
		HeightV: 0,
		BytesV:  vtxBytes0,
	}
	vtx1 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID1,
			StatusV: choices.Processing,
		},
		ParentsV: []lux.Vertex{
			vtx0,
		},
		HeightV: 1,
		BytesV:  vtxBytes1,
	}
	vtx2 := &lux.TestVertex{ // vtx2 is the stop vertex
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID2,
			StatusV: choices.Processing,
		},
		ParentsV: []lux.Vertex{
			vtx1,
		},
		HeightV: 2,
		BytesV:  vtxBytes2,
	}

	config.StopVertexID = vtxID2
	bs, err := New(
		config,
		func(context.Context, uint32) error {
			config.Ctx.State.Set(snow.EngineState{
				Type:  p2ppb.EngineType_ENGINE_TYPE_LUX,
				State: snow.NormalOp,
			})
			return nil
		},
		prometheus.NewRegistry(),
	)
	require.NoError(err)

	manager.GetVtxF = func(_ context.Context, vtxID ids.ID) (lux.Vertex, error) {
		switch vtxID {
		case vtxID0:
			return vtx0, nil
		case vtxID1:
			return vtx1, nil
		case vtxID2:
			return vtx2, nil
		default:
			require.FailNow(errUnknownVertex.Error())
			return nil, errUnknownVertex
		}
	}

	manager.ParseVtxF = func(_ context.Context, vtxBytes []byte) (lux.Vertex, error) {
		switch {
		case bytes.Equal(vtxBytes, vtxBytes0):
			return vtx0, nil
		case bytes.Equal(vtxBytes, vtxBytes1):
			return vtx1, nil
		case bytes.Equal(vtxBytes, vtxBytes2):
			return vtx2, nil
		default:
			require.FailNow(errParsedUnknownVertex.Error())
			return nil, errParsedUnknownVertex
		}
	}

	manager.StopVertexAcceptedF = func(context.Context) (bool, error) {
		return vtx2.Status() == choices.Accepted, nil
	}

	manager.EdgeF = func(context.Context) []ids.ID {
		require.Equal(choices.Accepted, vtx2.Status())
		return []ids.ID{vtxID2}
	}

	vm.LinearizeF = func(_ context.Context, stopVertexID ids.ID) error {
		require.Equal(vtxID2, stopVertexID)
		return nil
	}

	vm.CantSetState = false
	require.NoError(bs.Start(context.Background(), 0))
	require.Equal(snow.NormalOp, config.Ctx.State.Get().State)
	require.Equal(choices.Accepted, vtx0.Status())
	require.Equal(choices.Accepted, vtx1.Status())
	require.Equal(choices.Accepted, vtx2.Status())
}

// Accepted frontier has one vertex, which has one vertex as a dependency.
// Requests again and gets an unexpected vertex. Requests again and gets the
// expected vertex and an additional vertex that should not be accepted.
func TestBootstrapperByzantineResponses(t *testing.T) {
	require := require.New(t)

	config, peerID, sender, manager, vm := newConfig(t)

	vtxID0 := ids.Empty.Prefix(0)
	vtxID1 := ids.Empty.Prefix(1)
	vtxID2 := ids.Empty.Prefix(2)

	vtxBytes0 := []byte{0}
	vtxBytes1 := []byte{1}
	vtxBytes2 := []byte{2}

	vtx0 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID0,
			StatusV: choices.Unknown,
		},
		HeightV: 0,
		BytesV:  vtxBytes0,
	}
	vtx1 := &lux.TestVertex{ // vtx1 is the stop vertex
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID1,
			StatusV: choices.Processing,
		},
		ParentsV: []lux.Vertex{vtx0},
		HeightV:  1,
		BytesV:   vtxBytes1,
	}
	// Should not receive transitive votes from [vtx1]
	vtx2 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID2,
			StatusV: choices.Unknown,
		},
		HeightV: 0,
		BytesV:  vtxBytes2,
	}

	config.StopVertexID = vtxID1
	bs, err := New(
		config,
		func(context.Context, uint32) error {
			config.Ctx.State.Set(snow.EngineState{
				Type:  p2ppb.EngineType_ENGINE_TYPE_LUX,
				State: snow.NormalOp,
			})
			return nil
		},
		prometheus.NewRegistry(),
	)
	require.NoError(err)

	manager.GetVtxF = func(_ context.Context, vtxID ids.ID) (lux.Vertex, error) {
		switch vtxID {
		case vtxID1:
			return vtx1, nil
		case vtxID0:
			return nil, errUnknownVertex
		default:
			require.FailNow(errUnknownVertex.Error())
			return nil, errUnknownVertex
		}
	}

	requestID := new(uint32)
	reqVtxID := ids.Empty
	sender.SendGetAncestorsF = func(_ context.Context, vdr ids.NodeID, reqID uint32, vtxID ids.ID) {
		require.Equal(peerID, vdr)
		require.Equal(vtxID0, vtxID)

		*requestID = reqID
		reqVtxID = vtxID
	}

	manager.ParseVtxF = func(_ context.Context, vtxBytes []byte) (lux.Vertex, error) {
		switch {
		case bytes.Equal(vtxBytes, vtxBytes0):
			vtx0.StatusV = choices.Processing
			return vtx0, nil
		case bytes.Equal(vtxBytes, vtxBytes1):
			vtx1.StatusV = choices.Processing
			return vtx1, nil
		case bytes.Equal(vtxBytes, vtxBytes2):
			vtx2.StatusV = choices.Processing
			return vtx2, nil
		default:
			require.FailNow(errParsedUnknownVertex.Error())
			return nil, errParsedUnknownVertex
		}
	}

	vm.CantSetState = false
	require.NoError(bs.Start(context.Background(), 0)) // should request vtx0
	require.Equal(vtxID0, reqVtxID)

	oldReqID := *requestID
	require.NoError(bs.Ancestors(context.Background(), peerID, *requestID, [][]byte{vtxBytes2})) // send unexpected vertex
	require.NotEqual(oldReqID, *requestID)                                                       // should have sent a new request

	oldReqID = *requestID
	manager.GetVtxF = func(_ context.Context, vtxID ids.ID) (lux.Vertex, error) {
		switch vtxID {
		case vtxID1:
			return vtx1, nil
		case vtxID0:
			return vtx0, nil
		default:
			require.FailNow(errUnknownVertex.Error())
			return nil, errUnknownVertex
		}
	}

	manager.StopVertexAcceptedF = func(context.Context) (bool, error) {
		return vtx1.Status() == choices.Accepted, nil
	}

	manager.EdgeF = func(context.Context) []ids.ID {
		require.Equal(choices.Accepted, vtx1.Status())
		return []ids.ID{vtxID1}
	}

	vm.LinearizeF = func(_ context.Context, stopVertexID ids.ID) error {
		require.Equal(vtxID1, stopVertexID)
		return nil
	}

	require.NoError(bs.Ancestors(context.Background(), peerID, *requestID, [][]byte{vtxBytes0, vtxBytes2})) // send expected vertex and vertex that should not be accepted
	require.Equal(oldReqID, *requestID)                                                                     // shouldn't have sent a new request
	require.Equal(snow.NormalOp, config.Ctx.State.Get().State)
	require.Equal(choices.Accepted, vtx0.Status())
	require.Equal(choices.Accepted, vtx1.Status())
	require.Equal(choices.Processing, vtx2.Status())
}

// Vertex has a dependency and tx has a dependency
func TestBootstrapperTxDependencies(t *testing.T) {
	require := require.New(t)

	config, peerID, sender, manager, vm := newConfig(t)

	txID0 := ids.GenerateTestID()
	txID1 := ids.GenerateTestID()

	txBytes0 := []byte{0}
	txBytes1 := []byte{1}

	innerTx0 := &snowstorm.TestTx{
		TestDecidable: choices.TestDecidable{
			IDV:     txID0,
			StatusV: choices.Processing,
		},
		BytesV: txBytes0,
	}

	// Depends on tx0
	tx1 := &snowstorm.TestTx{
		TestDecidable: choices.TestDecidable{
			IDV:     txID1,
			StatusV: choices.Processing,
		},
		DependenciesV: set.Of(innerTx0.IDV),
		BytesV:        txBytes1,
	}

	tx0 := &testTx{
		Tx: innerTx0,
		tx: tx1,
	}

	vtxID0 := ids.GenerateTestID()
	vtxID1 := ids.GenerateTestID()

	vtxBytes0 := []byte{2}
	vtxBytes1 := []byte{3}
	vm.ParseTxF = func(_ context.Context, b []byte) (snowstorm.Tx, error) {
		switch {
		case bytes.Equal(b, txBytes0):
			return tx0, nil
		case bytes.Equal(b, txBytes1):
			return tx1, nil
		default:
			return nil, errUnknownTx
		}
	}

	vtx0 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID0,
			StatusV: choices.Unknown,
		},
		HeightV: 0,
		TxsV:    []snowstorm.Tx{tx1},
		BytesV:  vtxBytes0,
	}
	vtx1 := &lux.TestVertex{ // vtx1 is the stop vertex
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID1,
			StatusV: choices.Processing,
		},
		ParentsV: []lux.Vertex{vtx0}, // Depends on vtx0
		HeightV:  1,
		TxsV:     []snowstorm.Tx{tx0},
		BytesV:   vtxBytes1,
	}

	config.StopVertexID = vtxID1
	bs, err := New(
		config,
		func(context.Context, uint32) error {
			config.Ctx.State.Set(snow.EngineState{
				Type:  p2ppb.EngineType_ENGINE_TYPE_LUX,
				State: snow.NormalOp,
			})
			return nil
		},
		prometheus.NewRegistry(),
	)
	require.NoError(err)

	manager.ParseVtxF = func(_ context.Context, vtxBytes []byte) (lux.Vertex, error) {
		switch {
		case bytes.Equal(vtxBytes, vtxBytes1):
			return vtx1, nil
		case bytes.Equal(vtxBytes, vtxBytes0):
			return vtx0, nil
		default:
			require.FailNow(errParsedUnknownVertex.Error())
			return nil, errParsedUnknownVertex
		}
	}
	manager.GetVtxF = func(_ context.Context, vtxID ids.ID) (lux.Vertex, error) {
		switch vtxID {
		case vtxID1:
			return vtx1, nil
		case vtxID0:
			return nil, errUnknownVertex
		default:
			require.FailNow(errUnknownVertex.Error())
			return nil, errUnknownVertex
		}
	}

	reqIDPtr := new(uint32)
	sender.SendGetAncestorsF = func(_ context.Context, vdr ids.NodeID, reqID uint32, vtxID ids.ID) {
		require.Equal(peerID, vdr)
		require.Equal(vtxID0, vtxID)

		*reqIDPtr = reqID
	}

	vm.CantSetState = false
	require.NoError(bs.Start(context.Background(), 0))

	manager.ParseVtxF = func(_ context.Context, vtxBytes []byte) (lux.Vertex, error) {
		switch {
		case bytes.Equal(vtxBytes, vtxBytes1):
			return vtx1, nil
		case bytes.Equal(vtxBytes, vtxBytes0):
			vtx0.StatusV = choices.Processing
			return vtx0, nil
		default:
			require.FailNow(errParsedUnknownVertex.Error())
			return nil, errParsedUnknownVertex
		}
	}

	manager.StopVertexAcceptedF = func(context.Context) (bool, error) {
		return vtx1.Status() == choices.Accepted, nil
	}

	manager.EdgeF = func(context.Context) []ids.ID {
		require.Equal(choices.Accepted, vtx1.Status())
		return []ids.ID{vtxID1}
	}

	vm.LinearizeF = func(_ context.Context, stopVertexID ids.ID) error {
		require.Equal(vtxID1, stopVertexID)
		return nil
	}

	require.NoError(bs.Ancestors(context.Background(), peerID, *reqIDPtr, [][]byte{vtxBytes0}))
	require.Equal(snow.NormalOp, config.Ctx.State.Get().State)
	require.Equal(choices.Accepted, tx0.Status())
	require.Equal(choices.Accepted, tx1.Status())
	require.Equal(choices.Accepted, vtx0.Status())
	require.Equal(choices.Accepted, vtx1.Status())
}

// Ancestors only contains 1 of the two needed vertices; have to issue another GetAncestors
func TestBootstrapperIncompleteAncestors(t *testing.T) {
	require := require.New(t)

	config, peerID, sender, manager, vm := newConfig(t)

	vtxID0 := ids.Empty.Prefix(0)
	vtxID1 := ids.Empty.Prefix(1)
	vtxID2 := ids.Empty.Prefix(2)

	vtxBytes0 := []byte{0}
	vtxBytes1 := []byte{1}
	vtxBytes2 := []byte{2}

	vtx0 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID0,
			StatusV: choices.Unknown,
		},
		HeightV: 0,
		BytesV:  vtxBytes0,
	}
	vtx1 := &lux.TestVertex{
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID1,
			StatusV: choices.Unknown,
		},
		ParentsV: []lux.Vertex{vtx0},
		HeightV:  1,
		BytesV:   vtxBytes1,
	}
	vtx2 := &lux.TestVertex{ // vtx2 is the stop vertex
		TestDecidable: choices.TestDecidable{
			IDV:     vtxID2,
			StatusV: choices.Processing,
		},
		ParentsV: []lux.Vertex{vtx1},
		HeightV:  2,
		BytesV:   vtxBytes2,
	}

	config.StopVertexID = vtxID2
	bs, err := New(
		config,
		func(context.Context, uint32) error {
			config.Ctx.State.Set(snow.EngineState{
				Type:  p2ppb.EngineType_ENGINE_TYPE_LUX,
				State: snow.NormalOp,
			})
			return nil
		},
		prometheus.NewRegistry(),
	)
	require.NoError(err)

	manager.GetVtxF = func(_ context.Context, vtxID ids.ID) (lux.Vertex, error) {
		switch {
		case vtxID == vtxID0:
			return nil, errUnknownVertex
		case vtxID == vtxID1:
			return nil, errUnknownVertex
		case vtxID == vtxID2:
			return vtx2, nil
		default:
			require.FailNow(errUnknownVertex.Error())
			return nil, errUnknownVertex
		}
	}
	manager.ParseVtxF = func(_ context.Context, vtxBytes []byte) (lux.Vertex, error) {
		switch {
		case bytes.Equal(vtxBytes, vtxBytes0):
			vtx0.StatusV = choices.Processing
			return vtx0, nil
		case bytes.Equal(vtxBytes, vtxBytes1):
			vtx1.StatusV = choices.Processing
			return vtx1, nil
		case bytes.Equal(vtxBytes, vtxBytes2):
			return vtx2, nil
		default:
			require.FailNow(errParsedUnknownVertex.Error())
			return nil, errParsedUnknownVertex
		}
	}
	reqIDPtr := new(uint32)
	requested := ids.Empty
	sender.SendGetAncestorsF = func(_ context.Context, vdr ids.NodeID, reqID uint32, vtxID ids.ID) {
		require.Equal(peerID, vdr)
		require.Contains([]ids.ID{vtxID1, vtxID0}, vtxID)

		*reqIDPtr = reqID
		requested = vtxID
	}

	vm.CantSetState = false
	require.NoError(bs.Start(context.Background(), 0)) // should request vtx1
	require.Equal(vtxID1, requested)

	require.NoError(bs.Ancestors(context.Background(), peerID, *reqIDPtr, [][]byte{vtxBytes1})) // Provide vtx1; should request vtx0
	require.Equal(snow.Bootstrapping, bs.Context().State.Get().State)
	require.Equal(vtxID0, requested)

	manager.StopVertexAcceptedF = func(context.Context) (bool, error) {
		return vtx2.Status() == choices.Accepted, nil
	}

	manager.EdgeF = func(context.Context) []ids.ID {
		require.Equal(choices.Accepted, vtx2.Status())
		return []ids.ID{vtxID2}
	}

	vm.LinearizeF = func(_ context.Context, stopVertexID ids.ID) error {
		require.Equal(vtxID2, stopVertexID)
		return nil
	}

	require.NoError(bs.Ancestors(context.Background(), peerID, *reqIDPtr, [][]byte{vtxBytes0})) // Provide vtx0; can finish now
	require.Equal(snow.NormalOp, bs.Context().State.Get().State)
	require.Equal(choices.Accepted, vtx0.Status())
	require.Equal(choices.Accepted, vtx1.Status())
	require.Equal(choices.Accepted, vtx2.Status())
}
