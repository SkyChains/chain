// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"time"

	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/vms/avm/txs"
)

// Block defines the common stateless interface for all blocks
type Block interface {
	snow.ContextInitializable

	ID() ids.ID
	Parent() ids.ID
	Height() uint64
	// Timestamp that this block was created at
	Timestamp() time.Time
	MerkleRoot() ids.ID
	Bytes() []byte

	// Txs returns the transactions contained in the block
	Txs() []*txs.Tx

	// note: initialize does not assume that the transactions are initialized,
	// and initializes them itself.
	initialize(bytes []byte, cm codec.Manager) error
}
