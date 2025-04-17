// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

// Factory returns new instances of Consensus
type Factory interface {
	New() Consensus
}
