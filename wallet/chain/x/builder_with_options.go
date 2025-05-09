// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/vms/avm/txs"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/wallet/subnet/primary/common"
)

var _ Builder = (*builderWithOptions)(nil)

type builderWithOptions struct {
	Builder
	options []common.Option
}

// NewBuilderWithOptions returns a new transaction builder that will use the
// given options by default.
//
//   - [builder] is the builder that will be called to perform the underlying
//     operations.
//   - [options] will be provided to the builder in addition to the options
//     provided in the method calls.
func NewBuilderWithOptions(builder Builder, options ...common.Option) Builder {
	return &builderWithOptions{
		Builder: builder,
		options: options,
	}
}

func (b *builderWithOptions) GetFTBalance(
	options ...common.Option,
) (map[ids.ID]uint64, error) {
	return b.Builder.GetFTBalance(
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) GetImportableBalance(
	chainID ids.ID,
	options ...common.Option,
) (map[ids.ID]uint64, error) {
	return b.Builder.GetImportableBalance(
		chainID,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewBaseTx(
	outputs []*lux.TransferableOutput,
	options ...common.Option,
) (*txs.BaseTx, error) {
	return b.Builder.NewBaseTx(
		outputs,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewCreateAssetTx(
	name string,
	symbol string,
	denomination byte,
	initialState map[uint32][]verify.State,
	options ...common.Option,
) (*txs.CreateAssetTx, error) {
	return b.Builder.NewCreateAssetTx(
		name,
		symbol,
		denomination,
		initialState,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTx(
	operations []*txs.Operation,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTx(
		operations,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxMintFT(
	outputs map[ids.ID]*secp256k1fx.TransferOutput,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxMintFT(
		outputs,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxMintNFT(
	assetID ids.ID,
	payload []byte,
	owners []*secp256k1fx.OutputOwners,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxMintNFT(
		assetID,
		payload,
		owners,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxMintProperty(
	assetID ids.ID,
	owner *secp256k1fx.OutputOwners,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxMintProperty(
		assetID,
		owner,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxBurnProperty(
	assetID ids.ID,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxBurnProperty(
		assetID,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewImportTx(
	chainID ids.ID,
	to *secp256k1fx.OutputOwners,
	options ...common.Option,
) (*txs.ImportTx, error) {
	return b.Builder.NewImportTx(
		chainID,
		to,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewExportTx(
	chainID ids.ID,
	outputs []*lux.TransferableOutput,
	options ...common.Option,
) (*txs.ExportTx, error) {
	return b.Builder.NewExportTx(
		chainID,
		outputs,
		common.UnionOptions(b.options, options)...,
	)
}
