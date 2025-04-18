// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package signer

import (
	"context"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/crypto/keychain"
	"github.com/skychains/chain/vms/avm/txs"
	"github.com/skychains/chain/vms/components/lux"
)

var _ Signer = (*signer)(nil)

type Signer interface {
	// Sign adds as many missing signatures as possible to the provided
	// transaction.
	//
	// If there are already some signatures on the transaction, those signatures
	// will not be removed.
	//
	// If the signer doesn't have the ability to provide a required signature,
	// the signature slot will be skipped without reporting an error.
	Sign(ctx context.Context, tx *txs.Tx) error
}

type Backend interface {
	GetUTXO(ctx context.Context, chainID, utxoID ids.ID) (*lux.UTXO, error)
}

type signer struct {
	kc      keychain.Keychain
	backend Backend
}

func New(kc keychain.Keychain, backend Backend) Signer {
	return &signer{
		kc:      kc,
		backend: backend,
	}
}

func (s *signer) Sign(ctx context.Context, tx *txs.Tx) error {
	return tx.Unsigned.Visit(&visitor{
		kc:      s.kc,
		backend: s.backend,
		ctx:     ctx,
		tx:      tx,
	})
}

func SignUnsigned(
	ctx context.Context,
	signer Signer,
	utx txs.UnsignedTx,
) (*txs.Tx, error) {
	tx := &txs.Tx{Unsigned: utx}
	return tx, signer.Sign(ctx, tx)
}
