// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowstorm

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/utils/set"
)

var _ Tx = (*TestTx)(nil)

// TestTx is a useful test tx
type TestTx struct {
	choices.TestDecidable

	DependenciesV    []Tx
	DependenciesErrV error
	InputIDsV        []ids.ID
	HasWhitelistV    bool
	WhitelistV       set.Set[ids.ID]
	WhitelistErrV    error
	VerifyV          error
	BytesV           []byte
}

func (t *TestTx) Dependencies() ([]Tx, error) {
	return t.DependenciesV, t.DependenciesErrV
}

func (t *TestTx) InputIDs() []ids.ID {
	return t.InputIDsV
}

func (t *TestTx) HasWhitelist() bool {
	return t.HasWhitelistV
}

<<<<<<< HEAD
func (t *TestTx) Whitelist(context.Context) (set.Set[ids.ID], error) {
	return t.WhitelistV, t.WhitelistErrV
}

func (t *TestTx) Verify(context.Context) error {
=======
func (t *TestTx) Whitelist() (ids.Set, error) {
	return t.WhitelistV, t.WhitelistErrV
}

func (t *TestTx) Verify() error {
>>>>>>> 55bd9343c (Add EmptyLines linter (#2233))
	return t.VerifyV
}

func (t *TestTx) Bytes() []byte {
	return t.BytesV
}
