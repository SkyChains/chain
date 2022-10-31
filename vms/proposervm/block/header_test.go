// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import "github.com/stretchr/testify/require"

func equalHeader(require *require.Assertions, want, have Header) {
	require.Equal(want.ChainID(), have.ChainID())
	require.Equal(want.ParentID(), have.ParentID())
	require.Equal(want.BodyID(), have.BodyID())
}
