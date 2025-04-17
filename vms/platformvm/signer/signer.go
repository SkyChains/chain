// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package signer

import (
	"github.com/skychains/chain/utils/crypto/bls"
	"github.com/skychains/chain/vms/components/verify"
)

type Signer interface {
	verify.Verifiable

	// Key returns the public BLS key if it exists.
	// Note: [nil] will be returned if the key does not exist.
	// Invariant: Only called after [Verify] returns [nil].
	Key() *bls.PublicKey
}
