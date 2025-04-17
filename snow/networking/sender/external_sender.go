// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package sender

import (
	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/message"
	"github.com/SkyChains/chain/snow/engine/common"
	"github.com/SkyChains/chain/subnets"
	"github.com/SkyChains/chain/utils/set"
)

// ExternalSender sends consensus messages to other validators
// Right now this is implemented in the networking package
type ExternalSender interface {
	Send(
		msg message.OutboundMessage,
		config common.SendConfig,
		subnetID ids.ID,
		allower subnets.Allower,
	) set.Set[ids.NodeID]
}
