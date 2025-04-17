// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"time"

	"github.com/skychains/chain/chains/atomic"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/platformvm/block"
	"github.com/skychains/chain/vms/platformvm/state"
)

type proposalBlockState struct {
	onDecisionState state.Diff
	onCommitState   state.Diff
	onAbortState    state.Diff
}

// The state of a block.
// Note that not all fields will be set for a given block.
type blockState struct {
	proposalBlockState
	statelessBlock block.Block

	onAcceptState state.Diff
	onAcceptFunc  func()

	inputs         set.Set[ids.ID]
	timestamp      time.Time
	atomicRequests map[ids.ID]*atomic.Requests
}
