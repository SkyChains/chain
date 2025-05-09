// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"context"

	"github.com/skychains/chain/api/info"
	"github.com/skychains/chain/vms/avm"
	"github.com/skychains/chain/wallet/chain/x/builder"
)

func NewContextFromURI(ctx context.Context, uri string) (*builder.Context, error) {
	infoClient := info.NewClient(uri)
	xChainClient := avm.NewClient(uri, builder.Alias)
	return NewContextFromClients(ctx, infoClient, xChainClient)
}

func NewContextFromClients(
	ctx context.Context,
	infoClient info.Client,
	xChainClient avm.Client,
) (*builder.Context, error) {
	networkID, err := infoClient.GetNetworkID(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := infoClient.GetBlockchainID(ctx, builder.Alias)
	if err != nil {
		return nil, err
	}

	asset, err := xChainClient.GetAssetDescription(ctx, "LUX")
	if err != nil {
		return nil, err
	}

	txFees, err := infoClient.GetTxFee(ctx)
	if err != nil {
		return nil, err
	}

	return &builder.Context{
		NetworkID:        networkID,
		BlockchainID:     chainID,
		LUXAssetID:      asset.AssetID,
		BaseTxFee:        uint64(txFees.TxFee),
		CreateAssetTxFee: uint64(txFees.CreateAssetTxFee),
	}, nil
}
