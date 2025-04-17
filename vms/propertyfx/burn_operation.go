// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/secp256k1fx"
)

type BurnOperation struct {
	secp256k1fx.Input `serialize:"true"`
}

func (*BurnOperation) InitCtx(*snow.Context) {}

func (*BurnOperation) Outs() []verify.State {
	return nil
}
