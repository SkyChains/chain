// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package nftfx

import (
	"errors"

	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/secp256k1fx"
)

var errNilTransferOperation = errors.New("nil transfer operation")

type TransferOperation struct {
	Input  secp256k1fx.Input `serialize:"true" json:"input"`
	Output TransferOutput    `serialize:"true" json:"output"`
}

func (op *TransferOperation) InitCtx(ctx *snow.Context) {
	op.Output.OutputOwners.InitCtx(ctx)
}

func (op *TransferOperation) Cost() (uint64, error) {
	return op.Input.Cost()
}

func (op *TransferOperation) Outs() []verify.State {
	return []verify.State{&op.Output}
}

func (op *TransferOperation) Verify() error {
	switch {
	case op == nil:
		return errNilTransferOperation
	default:
		return verify.All(&op.Input, &op.Output)
	}
}
