// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package indexer

import "github.com/skychains/chain/ids"

// Container is something that gets accepted
// (a block, transaction or vertex)
type Container struct {
	// ID of this container
	ID ids.ID `serialize:"true"`
	// Byte representation of this container
	Bytes []byte `serialize:"true"`
	// Unix time, in nanoseconds, at which this container was accepted by this node
	Timestamp int64 `serialize:"true"`
}
