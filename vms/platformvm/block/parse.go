// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import "github.com/SkyChains/chain/codec"

func Parse(c codec.Manager, b []byte) (Block, error) {
	var blk Block
	if _, err := c.Unmarshal(b, &blk); err != nil {
		return nil, err
	}
	return blk, blk.initialize(b)
}
