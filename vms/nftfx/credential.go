// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package nftfx

import "github.com/skychains/chain/vms/secp256k1fx"

type Credential struct {
	secp256k1fx.Credential `serialize:"true"`
}
