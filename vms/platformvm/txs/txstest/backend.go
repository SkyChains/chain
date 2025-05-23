// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txstest

import (
	"context"
	"math"

	"github.com/skychains/chain/chains/atomic"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/platformvm/fx"
	"github.com/skychains/chain/vms/platformvm/state"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/wallet/chain/p/builder"
	"github.com/skychains/chain/wallet/chain/p/signer"
)

var (
	_ builder.Backend = (*Backend)(nil)
	_ signer.Backend  = (*Backend)(nil)
)

func newBackend(
	addrs set.Set[ids.ShortID],
	state state.State,
	sharedMemory atomic.SharedMemory,
) *Backend {
	return &Backend{
		addrs:        addrs,
		state:        state,
		sharedMemory: sharedMemory,
	}
}

type Backend struct {
	addrs        set.Set[ids.ShortID]
	state        state.State
	sharedMemory atomic.SharedMemory
}

func (b *Backend) UTXOs(_ context.Context, sourceChainID ids.ID) ([]*lux.UTXO, error) {
	if sourceChainID == constants.PlatformChainID {
		return lux.GetAllUTXOs(b.state, b.addrs)
	}

	utxos, _, _, err := lux.GetAtomicUTXOs(
		b.sharedMemory,
		txs.Codec,
		sourceChainID,
		b.addrs,
		ids.ShortEmpty,
		ids.Empty,
		math.MaxInt,
	)
	return utxos, err
}

func (b *Backend) GetUTXO(_ context.Context, chainID, utxoID ids.ID) (*lux.UTXO, error) {
	if chainID == constants.PlatformChainID {
		return b.state.GetUTXO(utxoID)
	}

	utxoBytes, err := b.sharedMemory.Get(chainID, [][]byte{utxoID[:]})
	if err != nil {
		return nil, err
	}

	utxo := lux.UTXO{}
	if _, err := txs.Codec.Unmarshal(utxoBytes[0], &utxo); err != nil {
		return nil, err
	}
	return &utxo, nil
}

func (b *Backend) GetSubnetOwner(_ context.Context, subnetID ids.ID) (fx.Owner, error) {
	return b.state.GetSubnetOwner(subnetID)
}
