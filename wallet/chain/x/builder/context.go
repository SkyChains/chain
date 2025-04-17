// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package builder

import (
	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/snow"
	"github.com/SkyChains/chain/utils/constants"
	"github.com/SkyChains/chain/utils/logging"
)

const Alias = "X"

type Context struct {
	NetworkID        uint32
	BlockchainID     ids.ID
	LUXAssetID      ids.ID
	BaseTxFee        uint64
	CreateAssetTxFee uint64
}

func NewSnowContext(
	networkID uint32,
	blockchainID ids.ID,
	luxAssetID ids.ID,
) (*snow.Context, error) {
	lookup := ids.NewAliaser()
	return &snow.Context{
		NetworkID:   networkID,
		SubnetID:    constants.PrimaryNetworkID,
		ChainID:     blockchainID,
		XChainID:    blockchainID,
		LUXAssetID: luxAssetID,
		Log:         logging.NoLog{},
		BCLookup:    lookup,
	}, lookup.Alias(blockchainID, Alias)
}
