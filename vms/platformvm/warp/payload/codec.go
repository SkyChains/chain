// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package payload

import (
	"github.com/luxdefi/node/codec"
	"github.com/luxdefi/node/codec/linearcodec"
	"github.com/luxdefi/node/utils"
	"github.com/luxdefi/node/utils/units"
)

const (
	codecVersion = 0

	MaxMessageSize = 24 * units.KiB

	// Note: Modifying this variable can have subtle implications on memory
	// usage when parsing malformed payloads.
	MaxSliceLen = 24 * 1024
)

// Codec does serialization and deserialization for Warp messages.
var c codec.Manager

func init() {
	c = codec.NewManager(MaxMessageSize)
	lc := linearcodec.NewCustomMaxLength(MaxSliceLen)

	err := utils.Err(
		lc.RegisterType(&Hash{}),
		lc.RegisterType(&AddressedCall{}),
		c.RegisterCodec(codecVersion, lc),
	)
	if err != nil {
		panic(err)
	}
}
