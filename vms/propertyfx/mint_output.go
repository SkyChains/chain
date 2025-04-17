// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/secp256k1fx"
)

var _ verify.State = (*MintOutput)(nil)

type MintOutput struct {
	verify.IsState `json:"-"`

	secp256k1fx.OutputOwners `serialize:"true"`
}
