// Copyright (C) 2019-2023, Lux Partners Limited All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/pubsub"
	"github.com/luxdefi/node/vms/avm/txs"
	"github.com/luxdefi/node/vms/components/lux"
	"github.com/luxdefi/node/vms/secp256k1fx"
)

type mockFilter struct {
	addr []byte
}

func (f *mockFilter) Check(addr []byte) bool {
	return bytes.Equal(addr, f.addr)
}

func TestFilter(t *testing.T) {
	require := require.New(t)

	addrID := ids.ShortID{1}
	tx := txs.Tx{Unsigned: &txs.BaseTx{BaseTx: lux.BaseTx{
		Outs: []*lux.TransferableOutput{
			{
				Out: &secp256k1fx.TransferOutput{
					OutputOwners: secp256k1fx.OutputOwners{
						Addrs: []ids.ShortID{addrID},
					},
				},
			},
		},
	}}}
	addrBytes := addrID[:]

	fp := pubsub.NewFilterParam()
	require.NoError(fp.Add(addrBytes))

	parser := NewPubSubFilterer(&tx)
	fr, _ := parser.Filter([]pubsub.Filter{&mockFilter{addr: addrBytes}})
	require.Equal([]bool{true}, fr)
}
