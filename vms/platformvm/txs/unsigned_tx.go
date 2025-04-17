// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/secp256k1fx"
)

// UnsignedTx is an unsigned transaction
type UnsignedTx interface {
	// TODO: Remove this initialization pattern from both the platformvm and the
	// avm.
	snow.ContextInitializable
	secp256k1fx.UnsignedTx
	SetBytes(unsignedBytes []byte)

	// InputIDs returns the set of inputs this transaction consumes
	InputIDs() set.Set[ids.ID]

	Outputs() []*lux.TransferableOutput

	// Attempts to verify this transaction without any provided state.
	SyntacticVerify(ctx *snow.Context) error

	// Visit calls [visitor] with this transaction's concrete type
	Visit(visitor Visitor) error
}
