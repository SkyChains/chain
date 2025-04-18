// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package c

import (
	"context"

	"github.com/skychains/chain/api/info"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms/avm"
)

const Alias = "C"

type Context struct {
	NetworkID    uint32
	BlockchainID ids.ID
	LUXAssetID  ids.ID
}

func NewContextFromURI(ctx context.Context, uri string) (*Context, error) {
	infoClient := info.NewClient(uri)
	xChainClient := avm.NewClient(uri, "X")
	return NewContextFromClients(ctx, infoClient, xChainClient)
}

func NewContextFromClients(
	ctx context.Context,
	infoClient info.Client,
	xChainClient avm.Client,
) (*Context, error) {
	networkID, err := infoClient.GetNetworkID(ctx)
	if err != nil {
		return nil, err
	}

	blockchainID, err := infoClient.GetBlockchainID(ctx, Alias)
	if err != nil {
		return nil, err
	}

	luxAsset, err := xChainClient.GetAssetDescription(ctx, "LUX")
	if err != nil {
		return nil, err
	}

	return &Context{
		NetworkID:    networkID,
		BlockchainID: blockchainID,
		LUXAssetID:  luxAsset.AssetID,
	}, nil
}

func newSnowContext(c *Context) (*snow.Context, error) {
	lookup := ids.NewAliaser()
	return &snow.Context{
		NetworkID:   c.NetworkID,
		SubnetID:    constants.PrimaryNetworkID,
		ChainID:     c.BlockchainID,
		CChainID:    c.BlockchainID,
		LUXAssetID: c.LUXAssetID,
		Log:         logging.NoLog{},
		BCLookup:    lookup,
	}, lookup.Alias(c.BlockchainID, Alias)
}
