// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package sender

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/message"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/subnets"
	"github.com/skychains/chain/utils/set"
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
