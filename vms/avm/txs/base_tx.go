// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/snow"
	"github.com/SkyChains/chain/utils/set"
	"github.com/SkyChains/chain/vms/components/lux"
	"github.com/SkyChains/chain/vms/secp256k1fx"
)

var (
	_ UnsignedTx             = (*BaseTx)(nil)
	_ secp256k1fx.UnsignedTx = (*BaseTx)(nil)
)

// BaseTx is the basis of all transactions.
type BaseTx struct {
	lux.BaseTx `serialize:"true"`

	bytes []byte
}

func (t *BaseTx) InitCtx(ctx *snow.Context) {
	for _, out := range t.Outs {
		out.InitCtx(ctx)
	}
}

func (t *BaseTx) SetBytes(bytes []byte) {
	t.bytes = bytes
}

func (t *BaseTx) Bytes() []byte {
	return t.bytes
}

func (t *BaseTx) InputIDs() set.Set[ids.ID] {
	inputIDs := set.NewSet[ids.ID](len(t.Ins))
	for _, in := range t.Ins {
		inputIDs.Add(in.InputID())
	}
	return inputIDs
}

func (t *BaseTx) Visit(v Visitor) error {
	return v.BaseTx(t)
}
