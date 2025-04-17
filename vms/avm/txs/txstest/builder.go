// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txstest

import (
	"context"
	"fmt"

	"github.com/SkyChains/chain/codec"
	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/snow"
	"github.com/SkyChains/chain/vms/avm/config"
	"github.com/SkyChains/chain/vms/avm/state"
	"github.com/SkyChains/chain/vms/avm/txs"
	"github.com/SkyChains/chain/vms/components/lux"
	"github.com/SkyChains/chain/vms/components/verify"
	"github.com/SkyChains/chain/vms/secp256k1fx"
	"github.com/SkyChains/chain/wallet/chain/x/builder"
	"github.com/SkyChains/chain/wallet/chain/x/signer"
	"github.com/SkyChains/chain/wallet/subnet/primary/common"
)

type Builder struct {
	utxos *utxos
	ctx   *builder.Context
}

func New(
	codec codec.Manager,
	ctx *snow.Context,
	cfg *config.Config,
	feeAssetID ids.ID,
	state state.State,
) *Builder {
	utxos := newUTXOs(ctx, state, ctx.SharedMemory, codec)
	return &Builder{
		utxos: utxos,
		ctx:   newContext(ctx, cfg, feeAssetID),
	}
}

func (b *Builder) CreateAssetTx(
	name, symbol string,
	denomination byte,
	initialStates map[uint32][]verify.State,
	kc *secp256k1fx.Keychain,
	changeAddr ids.ShortID,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	utx, err := xBuilder.NewCreateAssetTx(
		name,
		symbol,
		denomination,
		initialStates,
		common.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{changeAddr},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed building base tx: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) BaseTx(
	outs []*lux.TransferableOutput,
	memo []byte,
	kc *secp256k1fx.Keychain,
	changeAddr ids.ShortID,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	utx, err := xBuilder.NewBaseTx(
		outs,
		common.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{changeAddr},
		}),
		common.WithMemo(memo),
	)
	if err != nil {
		return nil, fmt.Errorf("failed building base tx: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) MintNFT(
	assetID ids.ID,
	payload []byte,
	owners []*secp256k1fx.OutputOwners,
	kc *secp256k1fx.Keychain,
	changeAddr ids.ShortID,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	utx, err := xBuilder.NewOperationTxMintNFT(
		assetID,
		payload,
		owners,
		common.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{changeAddr},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed minting NFTs: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) MintFTs(
	outputs map[ids.ID]*secp256k1fx.TransferOutput,
	kc *secp256k1fx.Keychain,
	changeAddr ids.ShortID,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	utx, err := xBuilder.NewOperationTxMintFT(
		outputs,
		common.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{changeAddr},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed minting FTs: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) Operation(
	ops []*txs.Operation,
	kc *secp256k1fx.Keychain,
	changeAddr ids.ShortID,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	utx, err := xBuilder.NewOperationTx(
		ops,
		common.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{changeAddr},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed building operation tx: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) ImportTx(
	sourceChain ids.ID,
	to ids.ShortID,
	kc *secp256k1fx.Keychain,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	outOwner := &secp256k1fx.OutputOwners{
		Locktime:  0,
		Threshold: 1,
		Addrs:     []ids.ShortID{to},
	}

	utx, err := xBuilder.NewImportTx(
		sourceChain,
		outOwner,
	)
	if err != nil {
		return nil, fmt.Errorf("failed building import tx: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) ExportTx(
	destinationChain ids.ID,
	to ids.ShortID,
	exportedAssetID ids.ID,
	exportedAmt uint64,
	kc *secp256k1fx.Keychain,
	changeAddr ids.ShortID,
) (*txs.Tx, error) {
	xBuilder, xSigner := b.builders(kc)

	outputs := []*lux.TransferableOutput{{
		Asset: lux.Asset{ID: exportedAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: exportedAmt,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  0,
				Threshold: 1,
				Addrs:     []ids.ShortID{to},
			},
		},
	}}

	utx, err := xBuilder.NewExportTx(
		destinationChain,
		outputs,
		common.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{changeAddr},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed building export tx: %w", err)
	}

	return signer.SignUnsigned(context.Background(), xSigner, utx)
}

func (b *Builder) builders(kc *secp256k1fx.Keychain) (builder.Builder, signer.Signer) {
	var (
		addrs = kc.Addresses()
		wa    = &walletUTXOsAdapter{
			utxos: b.utxos,
			addrs: addrs,
		}
		builder = builder.New(addrs, b.ctx, wa)
		signer  = signer.New(kc, wa)
	)
	return builder, signer
}
