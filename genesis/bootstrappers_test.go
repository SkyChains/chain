// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/utils/constants"
)

func TestSampleBootstrappers(t *testing.T) {
	require := require.New(t)

	for networkID, networkName := range constants.NetworkIDToNetworkName {
		length := 10
		bootstrappers := SampleBootstrappers(networkID, length)
		t.Logf("%s bootstrappers: %+v", networkName, bootstrappers)

		if networkID == constants.MainnetID || networkID == constants.FujiID {
			require.Len(bootstrappers, length)
		}
	}
}
