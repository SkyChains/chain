// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package bloom

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/utils/units"
)

func TestNew(t *testing.T) {
	var (
		require  = require.New(t)
		maxN     = 10000
		p        = 0.1
		maxBytes = 1 * units.MiB // 1 MiB
	)
	f, err := New(maxN, p, maxBytes)
	require.NoError(err)
	require.NotNil(f)

	f.Add([]byte("hello"))

	checked := f.Check([]byte("hello"))
	require.True(checked, "should have contained the key")

	checked = f.Check([]byte("bye"))
	require.False(checked, "shouldn't have contained the key")
}
