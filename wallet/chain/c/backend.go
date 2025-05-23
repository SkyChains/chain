// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package c

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/skychains/coreth/plugin/evm"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/utils/math"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/wallet/subnet/primary/common"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	_ Backend = (*backend)(nil)

	errUnknownTxType = errors.New("unknown tx type")
)

// Backend defines the full interface required to support a C-chain wallet.
type Backend interface {
	common.ChainUTXOs
	BuilderBackend
	SignerBackend

	AcceptAtomicTx(ctx context.Context, tx *evm.Tx) error
}

type backend struct {
	common.ChainUTXOs

	accountsLock sync.RWMutex
	accounts     map[ethcommon.Address]*Account
}

type Account struct {
	Balance *big.Int
	Nonce   uint64
}

func NewBackend(
	utxos common.ChainUTXOs,
	accounts map[ethcommon.Address]*Account,
) Backend {
	return &backend{
		ChainUTXOs: utxos,
		accounts:   accounts,
	}
}

func (b *backend) AcceptAtomicTx(ctx context.Context, tx *evm.Tx) error {
	switch tx := tx.UnsignedAtomicTx.(type) {
	case *evm.UnsignedImportTx:
		for _, input := range tx.ImportedInputs {
			utxoID := input.InputID()
			if err := b.RemoveUTXO(ctx, tx.SourceChain, utxoID); err != nil {
				return err
			}
		}

		b.accountsLock.Lock()
		defer b.accountsLock.Unlock()

		for _, output := range tx.Outs {
			account, ok := b.accounts[output.Address]
			if !ok {
				continue
			}

			balance := new(big.Int).SetUint64(output.Amount)
			balance.Mul(balance, luxConversionRate)
			account.Balance.Add(account.Balance, balance)
		}
	case *evm.UnsignedExportTx:
		txID := tx.ID()
		for i, out := range tx.ExportedOutputs {
			err := b.AddUTXO(
				ctx,
				tx.DestinationChain,
				&lux.UTXO{
					UTXOID: lux.UTXOID{
						TxID:        txID,
						OutputIndex: uint32(i),
					},
					Asset: lux.Asset{ID: out.AssetID()},
					Out:   out.Out,
				},
			)
			if err != nil {
				return err
			}
		}

		b.accountsLock.Lock()
		defer b.accountsLock.Unlock()

		for _, input := range tx.Ins {
			account, ok := b.accounts[input.Address]
			if !ok {
				continue
			}

			balance := new(big.Int).SetUint64(input.Amount)
			balance.Mul(balance, luxConversionRate)
			if account.Balance.Cmp(balance) == -1 {
				return errInsufficientFunds
			}
			account.Balance.Sub(account.Balance, balance)

			newNonce, err := math.Add64(input.Nonce, 1)
			if err != nil {
				return err
			}
			account.Nonce = newNonce
		}
	default:
		return fmt.Errorf("%w: %T", errUnknownTxType, tx)
	}
	return nil
}

func (b *backend) Balance(_ context.Context, addr ethcommon.Address) (*big.Int, error) {
	b.accountsLock.RLock()
	defer b.accountsLock.RUnlock()

	account, exists := b.accounts[addr]
	if !exists {
		return nil, database.ErrNotFound
	}
	return account.Balance, nil
}

func (b *backend) Nonce(_ context.Context, addr ethcommon.Address) (uint64, error) {
	b.accountsLock.RLock()
	defer b.accountsLock.RUnlock()

	account, exists := b.accounts[addr]
	if !exists {
		return 0, database.ErrNotFound
	}
	return account.Nonce, nil
}
