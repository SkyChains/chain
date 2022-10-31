// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package builder

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luxdefi/luxd/chains/atomic"
	"github.com/luxdefi/luxd/database/prefixdb"
	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/utils/crypto"
	"github.com/luxdefi/luxd/vms/components/lux"
	"github.com/luxdefi/luxd/vms/platformvm/status"
	"github.com/luxdefi/luxd/vms/platformvm/txs"
	"github.com/luxdefi/luxd/vms/secp256k1fx"
)

func TestAtomicTxImports(t *testing.T) {
	require := require.New(t)

	env := newEnvironment(t)
	env.ctx.Lock.Lock()
	defer func() {
		if err := shutdownEnvironment(env); err != nil {
			t.Fatal(err)
		}
	}()

	utxoID := lux.UTXOID{
		TxID:        ids.Empty.Prefix(1),
		OutputIndex: 1,
	}
	amount := uint64(70000)
	recipientKey := preFundedKeys[1]

	m := atomic.NewMemory(prefixdb.New([]byte{5}, env.baseDB))

	env.msm.SharedMemory = m.NewSharedMemory(env.ctx.ChainID)
	peerSharedMemory := m.NewSharedMemory(env.ctx.XChainID)
	utxo := &lux.UTXO{
		UTXOID: utxoID,
		Asset:  lux.Asset{ID: luxAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: amount,
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{recipientKey.PublicKey().Address()},
			},
		},
	}
	utxoBytes, err := txs.Codec.Marshal(txs.Version, utxo)
	require.NoError(err)

	inputID := utxo.InputID()
	err = peerSharedMemory.Apply(map[ids.ID]*atomic.Requests{
		env.ctx.ChainID: {PutRequests: []*atomic.Element{{
			Key:   inputID[:],
			Value: utxoBytes,
			Traits: [][]byte{
				recipientKey.PublicKey().Address().Bytes(),
			},
		}}},
	})
	require.NoError(err)

	tx, err := env.txBuilder.NewImportTx(
		env.ctx.XChainID,
		recipientKey.PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{recipientKey},
		ids.ShortEmpty, // change addr
	)
	require.NoError(err)

	require.NoError(env.Builder.Add(tx))
	b, err := env.Builder.BuildBlock()
	require.NoError(err)
	// Test multiple verify calls work
	require.NoError(b.Verify())
	require.NoError(b.Accept())
	_, txStatus, err := env.state.GetTx(tx.ID())
	require.NoError(err)
	// Ensure transaction is in the committed state
	require.Equal(txStatus, status.Committed)
}
