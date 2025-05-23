// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"reflect"

	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/timer/mockable"
	"github.com/skychains/chain/utils/wrappers"
	"github.com/skychains/chain/vms/secp256k1fx"
)

var (
	_ codec.Registry = (*codecRegistry)(nil)
	_ secp256k1fx.VM = (*fxVM)(nil)
)

type codecRegistry struct {
	codecs      []codec.Registry
	index       int
	typeToIndex map[reflect.Type]int
}

func (cr *codecRegistry) RegisterType(val interface{}) error {
	valType := reflect.TypeOf(val)
	cr.typeToIndex[valType] = cr.index

	errs := wrappers.Errs{}
	for _, c := range cr.codecs {
		errs.Add(c.RegisterType(val))
	}
	return errs.Err
}

type fxVM struct {
	typeToFxIndex map[reflect.Type]int

	clock         *mockable.Clock
	log           logging.Logger
	codecRegistry codec.Registry
}

func (vm *fxVM) Clock() *mockable.Clock {
	return vm.clock
}

func (vm *fxVM) CodecRegistry() codec.Registry {
	return vm.codecRegistry
}

func (vm *fxVM) Logger() logging.Logger {
	return vm.log
}
