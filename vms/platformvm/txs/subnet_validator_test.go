// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
)

func TestSubnetValidatorVerifySubnetID(t *testing.T) {
	require := require.New(t)

	// Error path
	{
		vdr := &SubnetValidator{
			Subnet: constants.PrimaryNetworkID,
		}

		err := vdr.Verify()
		require.ErrorIs(err, errBadSubnetID)
	}

	// Happy path
	{
		vdr := &SubnetValidator{
			Subnet: ids.GenerateTestID(),
			Validator: Validator{
				Wght: 1,
			},
		}

		require.NoError(vdr.Verify())
	}
}
