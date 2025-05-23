// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/platformvm/txs"
)

func TestNewApricotAtomicBlock(t *testing.T) {
	require := require.New(t)

	parentID := ids.GenerateTestID()
	height := uint64(1337)
	tx := &txs.Tx{
		Unsigned: &txs.ImportTx{
			BaseTx: txs.BaseTx{
				BaseTx: lux.BaseTx{
					Ins:  []*lux.TransferableInput{},
					Outs: []*lux.TransferableOutput{},
				},
			},
			ImportedInputs: []*lux.TransferableInput{},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Initialize(txs.Codec))

	blk, err := NewApricotAtomicBlock(
		parentID,
		height,
		tx,
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotEmpty(blk.Bytes())
	require.NotEmpty(blk.Tx.Bytes())
	require.NotEqual(ids.Empty, blk.Tx.ID())
	require.Equal(tx.Bytes(), blk.Tx.Bytes())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}
