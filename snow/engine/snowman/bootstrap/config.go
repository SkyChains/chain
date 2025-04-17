// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package bootstrap

import (
	"github.com/skychains/chain/database"
	"github.com/skychains/chain/network/p2p"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/engine/common/tracker"
	"github.com/skychains/chain/snow/engine/snowman/block"
	"github.com/skychains/chain/snow/validators"
)

type Config struct {
	common.AllGetsServer

	Ctx     *snow.ConsensusContext
	Beacons validators.Manager

	SampleK          int
	StartupTracker   tracker.Startup
	Sender           common.Sender
	BootstrapTracker common.BootstrapTracker
	Timer            common.Timer

	// PeerTracker manages the set of nodes that we fetch the next block from.
	PeerTracker *p2p.PeerTracker

	// This node will only consider the first [AncestorsMaxContainersReceived]
	// containers in an ancestors message it receives.
	AncestorsMaxContainersReceived int

	// Database used to track the fetched, but not yet executed, blocks during
	// bootstrapping.
	DB database.Database

	VM block.ChainVM

	Bootstrapped func()
}
