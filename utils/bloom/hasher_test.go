// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package bloom

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/utils/units"
)

func TestCollisionResistance(t *testing.T) {
	require := require.New(t)

	f, err := New(8, 16*units.KiB)
	require.NoError(err)

	Add(f, []byte("hello world?"), []byte("so salty"))
	collision := Contains(f, []byte("hello world!"), []byte("so salty"))
	require.False(collision)
}

func BenchmarkHash(b *testing.B) {
	key := ids.GenerateTestID()
	salt := ids.GenerateTestID()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Hash(key[:], salt[:])
	}
}
