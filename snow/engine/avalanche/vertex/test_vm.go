// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build test

package vertex

import (
	"context"
	"errors"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/consensus/snowstorm"
	"github.com/skychains/chain/snow/engine/snowman/block"
)

var (
	errLinearize = errors.New("unexpectedly called Linearize")

	_ LinearizableVM = (*TestVM)(nil)
)

type TestVM struct {
	block.TestVM

	CantLinearize, CantParse bool

	LinearizeF func(context.Context, ids.ID) error
	ParseTxF   func(context.Context, []byte) (snowstorm.Tx, error)
}

func (vm *TestVM) Default(cant bool) {
	vm.TestVM.Default(cant)

	vm.CantParse = cant
}

func (vm *TestVM) Linearize(ctx context.Context, stopVertexID ids.ID) error {
	if vm.LinearizeF != nil {
		return vm.LinearizeF(ctx, stopVertexID)
	}
	if vm.CantLinearize && vm.T != nil {
		require.FailNow(vm.T, errLinearize.Error())
	}
	return errLinearize
}

func (vm *TestVM) ParseTx(ctx context.Context, b []byte) (snowstorm.Tx, error) {
	if vm.ParseTxF != nil {
		return vm.ParseTxF(ctx, b)
	}
	if vm.CantParse && vm.T != nil {
		require.FailNow(vm.T, errParse.Error())
	}
	return nil, errParse
}
