// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"context"
	"time"

	"github.com/skychains/chain/api/health"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/consensus/snowball"
	"github.com/skychains/chain/utils/bag"
)

// Consensus represents a general snowman instance that can be used directly to
// process a series of dependent operations.
type Consensus interface {
	health.Checker

	// Takes in the context, snowball parameters, and the last accepted block.
	Initialize(
		ctx *snow.ConsensusContext,
		params snowball.Parameters,
		lastAcceptedID ids.ID,
		lastAcceptedHeight uint64,
		lastAcceptedTime time.Time,
	) error

	// Returns the number of blocks processing
	NumProcessing() int

	// Add a new block.
	//
	// Add should not be called multiple times with the same block.
	// The parent block should either be the last accepted block or processing.
	//
	// Returns if a critical error has occurred.
	Add(Block) error

	// Processing returns true if the block ID is currently processing.
	Processing(ids.ID) bool

	// IsPreferred returns true if the block ID is preferred. Only the last
	// accepted block and processing blocks are considered preferred.
	IsPreferred(ids.ID) bool

	// Returns the ID and height of the last accepted decision.
	LastAccepted() (ids.ID, uint64)

	// Returns the ID of the tail of the strongly preferred sequence of
	// decisions.
	Preference() ids.ID

	// Returns the ID of the strongly preferred decision with the provided
	// height. Only the last accepted decision and processing decisions are
	// tracked.
	PreferenceAtHeight(height uint64) (ids.ID, bool)

	// RecordPoll collects the results of a network poll. Assumes all decisions
	// have been previously added. Returns if a critical error has occurred.
	RecordPoll(context.Context, bag.Bag[ids.ID]) error
}
