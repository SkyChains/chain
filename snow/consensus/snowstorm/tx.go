// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowstorm

import (
	"context"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/choices"
	"github.com/skychains/chain/utils/set"
)

// Tx consumes state.
type Tx interface {
	choices.Decidable

	// MissingDependencies returns the set of transactions that must be accepted
	// before this transaction is accepted.
	MissingDependencies() (set.Set[ids.ID], error)

	// Verify that the state transition this transaction would make if it were
	// accepted is valid. If the state transition is invalid, a non-nil error
	// should be returned.
	//
	// It is guaranteed that when Verify is called, all the dependencies of
	// this transaction have already been successfully verified.
	Verify(context.Context) error

	// Bytes returns the binary representation of this transaction.
	//
	// This is used for sending transactions to peers. Another node should be
	// able to parse these bytes to the same transaction.
	Bytes() []byte
}
