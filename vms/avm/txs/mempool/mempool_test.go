// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package mempool

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/vms/avm/txs"
	"github.com/skychains/chain/vms/components/lux"
)

func newMempool(toEngine chan<- common.Message) (Mempool, error) {
	return New("mempool", prometheus.NewRegistry(), toEngine)
}

func TestRequestBuildBlock(t *testing.T) {
	require := require.New(t)

	toEngine := make(chan common.Message, 1)
	mempool, err := newMempool(toEngine)
	require.NoError(err)

	mempool.RequestBuildBlock()
	select {
	case <-toEngine:
		require.FailNow("should not have sent message to engine")
	default:
	}

	tx := newTx(0, 32)
	require.NoError(mempool.Add(tx))

	mempool.RequestBuildBlock()
	mempool.RequestBuildBlock() // Must not deadlock
	select {
	case <-toEngine:
	default:
		require.FailNow("should have sent message to engine")
	}
	select {
	case <-toEngine:
		require.FailNow("should have only sent one message to engine")
	default:
	}
}

func newTx(index uint32, size int) *txs.Tx {
	tx := &txs.Tx{Unsigned: &txs.BaseTx{BaseTx: lux.BaseTx{
		Ins: []*lux.TransferableInput{{
			UTXOID: lux.UTXOID{
				TxID:        ids.ID{'t', 'x', 'I', 'D'},
				OutputIndex: index,
			},
		}},
	}}}
	tx.SetBytes(utils.RandomBytes(size), utils.RandomBytes(size))
	return tx
}
