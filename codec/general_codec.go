// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package codec

// GeneralCodec marshals and unmarshals structs including interfaces
type GeneralCodec interface {
	Codec
	Registry
}
