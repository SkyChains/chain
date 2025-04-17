// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/bls"
)

func TestSigner(t *testing.T) {
	for name, test := range SignerTests {
		t.Run(name, func(t *testing.T) {
			sk, err := bls.NewSecretKey()
			require.NoError(t, err)

			chainID := ids.GenerateTestID()
			s := NewSigner(sk, constants.UnitTestID, chainID)

			test(t, s, sk, constants.UnitTestID, chainID)
		})
	}
}
