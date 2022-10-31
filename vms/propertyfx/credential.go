// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/luxdefi/luxd/vms/secp256k1fx"
)

type Credential struct {
	secp256k1fx.Credential `serialize:"true"`
}
