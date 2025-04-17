// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"cmp"

	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/vms/avm/txs"
)

var _ utils.Sortable[*GenesisAsset] = (*GenesisAsset)(nil)

type Genesis struct {
	Txs []*GenesisAsset `serialize:"true"`
}

type GenesisAsset struct {
	Alias             string `serialize:"true"`
	txs.CreateAssetTx `serialize:"true"`
}

func (g *GenesisAsset) Compare(other *GenesisAsset) int {
	return cmp.Compare(g.Alias, other.Alias)
}
