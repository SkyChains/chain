// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"github.com/skychains/chain/api"
	"github.com/skychains/chain/pubsub"
	"github.com/skychains/chain/vms/avm/txs"
	"github.com/skychains/chain/vms/components/lux"
)

var _ pubsub.Filterer = (*connector)(nil)

type connector struct {
	tx *txs.Tx
}

func NewPubSubFilterer(tx *txs.Tx) pubsub.Filterer {
	return &connector{tx: tx}
}

// Apply the filter on the addresses.
func (f *connector) Filter(filters []pubsub.Filter) ([]bool, interface{}) {
	resp := make([]bool, len(filters))
	for _, utxo := range f.tx.UTXOs() {
		addressable, ok := utxo.Out.(lux.Addressable)
		if !ok {
			continue
		}

		for _, address := range addressable.Addresses() {
			for i, c := range filters {
				if resp[i] {
					continue
				}
				resp[i] = c.Check(address)
			}
		}
	}
	return resp, api.JSONTxID{
		TxID: f.tx.ID(),
	}
}
