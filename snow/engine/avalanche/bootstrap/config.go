// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package bootstrap

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/network/p2p"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/engine/lux/bootstrap/queue"
	"github.com/skychains/chain/snow/engine/lux/vertex"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/engine/common/tracker"
)

type Config struct {
	common.AllGetsServer

	Ctx *snow.ConsensusContext

	StartupTracker tracker.Startup
	Sender         common.Sender

	// PeerTracker manages the set of nodes that we fetch the next block from.
	PeerTracker *p2p.PeerTracker

	// This node will only consider the first [AncestorsMaxContainersReceived]
	// containers in an ancestors message it receives.
	AncestorsMaxContainersReceived int

	// VtxBlocked tracks operations that are blocked on vertices
	VtxBlocked *queue.JobsWithMissing
	// TxBlocked tracks operations that are blocked on transactions
	TxBlocked *queue.Jobs

	Manager vertex.Manager
	VM      vertex.LinearizableVM

	// If StopVertexID is empty, the engine will generate the stop vertex based
	// on the current state.
	StopVertexID ids.ID
}
