// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package c

import (
	"github.com/skychains/coreth/plugin/evm"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/wallet/subnet/primary/common"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

var _ Wallet = (*walletWithOptions)(nil)

func NewWalletWithOptions(
	wallet Wallet,
	options ...common.Option,
) Wallet {
	return &walletWithOptions{
		Wallet:  wallet,
		options: options,
	}
}

type walletWithOptions struct {
	Wallet
	options []common.Option
}

func (w *walletWithOptions) Builder() Builder {
	return NewBuilderWithOptions(
		w.Wallet.Builder(),
		w.options...,
	)
}

func (w *walletWithOptions) IssueImportTx(
	chainID ids.ID,
	to ethcommon.Address,
	options ...common.Option,
) (*evm.Tx, error) {
	return w.Wallet.IssueImportTx(
		chainID,
		to,
		common.UnionOptions(w.options, options)...,
	)
}

func (w *walletWithOptions) IssueExportTx(
	chainID ids.ID,
	outputs []*secp256k1fx.TransferOutput,
	options ...common.Option,
) (*evm.Tx, error) {
	return w.Wallet.IssueExportTx(
		chainID,
		outputs,
		common.UnionOptions(w.options, options)...,
	)
}

func (w *walletWithOptions) IssueUnsignedAtomicTx(
	utx evm.UnsignedAtomicTx,
	options ...common.Option,
) (*evm.Tx, error) {
	return w.Wallet.IssueUnsignedAtomicTx(
		utx,
		common.UnionOptions(w.options, options)...,
	)
}

func (w *walletWithOptions) IssueAtomicTx(
	tx *evm.Tx,
	options ...common.Option,
) error {
	return w.Wallet.IssueAtomicTx(
		tx,
		common.UnionOptions(w.options, options)...,
	)
}
