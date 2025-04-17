// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"github.com/SkyChains/chain/snow"
	"github.com/SkyChains/chain/snow/consensus/snowball"
	"github.com/SkyChains/chain/snow/consensus/snowman"
	"github.com/SkyChains/chain/snow/engine/common"
	"github.com/SkyChains/chain/snow/engine/common/tracker"
	"github.com/SkyChains/chain/snow/engine/snowman/block"
	"github.com/SkyChains/chain/snow/validators"
)

// Config wraps all the parameters needed for a snowman engine
type Config struct {
	common.AllGetsServer

	Ctx                 *snow.ConsensusContext
	VM                  block.ChainVM
	Sender              common.Sender
	Validators          validators.Manager
	ConnectedValidators tracker.Peers
	Params              snowball.Parameters
	Consensus           snowman.Consensus
	PartialSync         bool
}
