// Copyright (C) 2019-2023, Lux Partners Limited All rights reserved.
// See the file LICENSE for licensing terms.

package p

import (
	"sync"

	stdcontext "context"

	"github.com/luxdefi/node/database"
	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/utils/constants"
	"github.com/luxdefi/node/utils/set"
	"github.com/luxdefi/node/vms/components/lux"
	"github.com/luxdefi/node/vms/platformvm/txs"
	"github.com/luxdefi/node/wallet/subnet/primary/common"
)

var _ Backend = (*backend)(nil)

// Backend defines the full interface required to support a P-chain wallet.
type Backend interface {
	common.ChainUTXOs
	BuilderBackend
	SignerBackend

	AcceptTx(ctx stdcontext.Context, tx *txs.Tx) error
}

type backend struct {
	Context
	common.ChainUTXOs

	txsLock sync.RWMutex
	// txID -> tx
	txs map[ids.ID]*txs.Tx
}

func NewBackend(ctx Context, utxos common.ChainUTXOs, txs map[ids.ID]*txs.Tx) Backend {
	return &backend{
		Context:    ctx,
		ChainUTXOs: utxos,
		txs:        txs,
	}
}

func (b *backend) AcceptTx(ctx stdcontext.Context, tx *txs.Tx) error {
	txID := tx.ID()
	err := tx.Unsigned.Visit(&backendVisitor{
		b:    b,
		ctx:  ctx,
		txID: txID,
	})
	if err != nil {
		return err
	}

	producedUTXOSlice := tx.UTXOs()
	err = b.addUTXOs(ctx, constants.PlatformChainID, producedUTXOSlice)
	if err != nil {
		return err
	}

	b.txsLock.Lock()
	defer b.txsLock.Unlock()

	b.txs[txID] = tx
	return nil
}

func (b *backend) addUTXOs(ctx stdcontext.Context, destinationChainID ids.ID, utxos []*lux.UTXO) error {
	for _, utxo := range utxos {
		if err := b.AddUTXO(ctx, destinationChainID, utxo); err != nil {
			return err
		}
	}
	return nil
}

func (b *backend) removeUTXOs(ctx stdcontext.Context, sourceChain ids.ID, utxoIDs set.Set[ids.ID]) error {
	for utxoID := range utxoIDs {
		if err := b.RemoveUTXO(ctx, sourceChain, utxoID); err != nil {
			return err
		}
	}
	return nil
}

func (b *backend) GetTx(_ stdcontext.Context, txID ids.ID) (*txs.Tx, error) {
	b.txsLock.RLock()
	defer b.txsLock.RUnlock()

	tx, exists := b.txs[txID]
	if !exists {
		return nil, database.ErrNotFound
	}
	return tx, nil
}
