// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/codec/linearcodec"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/units"
)

const (
	codecVersion   = 0
	maxMessageSize = 512 * units.KiB
	maxSliceLen    = maxMessageSize
)

// Codec does serialization and deserialization
var c codec.Manager

func init() {
	c = codec.NewManager(maxMessageSize)
	lc := linearcodec.NewCustomMaxLength(maxSliceLen)

	err := utils.Err(
		lc.RegisterType(&Tx{}),
		c.RegisterCodec(codecVersion, lc),
	)
	if err != nil {
		panic(err)
	}
}
