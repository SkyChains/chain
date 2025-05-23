// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package mempool

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/vms/avm/txs"

	txmempool "github.com/skychains/chain/vms/txs/mempool"
)

var _ Mempool = (*mempool)(nil)

// Mempool contains transactions that have not yet been put into a block.
type Mempool interface {
	txmempool.Mempool[*txs.Tx]

	// RequestBuildBlock notifies the consensus engine that a block should be
	// built if there is at least one transaction in the mempool.
	RequestBuildBlock()
}

type mempool struct {
	txmempool.Mempool[*txs.Tx]

	toEngine chan<- common.Message
}

func New(
	namespace string,
	registerer prometheus.Registerer,
	toEngine chan<- common.Message,
) (Mempool, error) {
	metrics, err := txmempool.NewMetrics(namespace, registerer)
	if err != nil {
		return nil, err
	}
	pool := txmempool.New[*txs.Tx](
		metrics,
	)
	return &mempool{
		Mempool:  pool,
		toEngine: toEngine,
	}, nil
}

func (m *mempool) RequestBuildBlock() {
	if m.Len() == 0 {
		return
	}

	select {
	case m.toEngine <- common.PendingTxs:
	default:
	}
}
