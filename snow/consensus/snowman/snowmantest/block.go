// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowmantest

import (
	"cmp"
	"context"
	"time"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/choices"
	"github.com/skychains/chain/utils"
)

const (
	GenesisHeight        uint64 = 0
	GenesisUnixTimestamp int64  = 1
)

var (
	_ utils.Sortable[*Block] = (*Block)(nil)

	GenesisID        = ids.GenerateTestID()
	GenesisTimestamp = time.Unix(GenesisUnixTimestamp, 0)
	GenesisBytes     = GenesisID[:]
	Genesis          = BuildChain(1)[0]
)

func BuildChild(parent *Block) *Block {
	blkID := ids.GenerateTestID()
	return &Block{
		TestDecidable: choices.TestDecidable{
			IDV:     blkID,
			StatusV: choices.Processing,
		},
		ParentV:    parent.ID(),
		HeightV:    parent.Height() + 1,
		TimestampV: parent.Timestamp(),
		BytesV:     blkID[:],
	}
}

func BuildChain(length int) []*Block {
	if length == 0 {
		return nil
	}

	genesis := &Block{
		TestDecidable: choices.TestDecidable{
			IDV:     GenesisID,
			StatusV: choices.Accepted,
		},
		HeightV:    GenesisHeight,
		TimestampV: GenesisTimestamp,
		BytesV:     GenesisBytes,
	}
	return append([]*Block{genesis}, BuildDescendants(genesis, length-1)...)
}

func BuildDescendants(parent *Block, length int) []*Block {
	chain := make([]*Block, length)
	for i := range chain {
		parent = BuildChild(parent)
		chain[i] = parent
	}
	return chain
}

type Block struct {
	choices.TestDecidable

	ParentV    ids.ID
	HeightV    uint64
	TimestampV time.Time
	VerifyV    error
	BytesV     []byte
}

func (b *Block) Parent() ids.ID {
	return b.ParentV
}

func (b *Block) Height() uint64 {
	return b.HeightV
}

func (b *Block) Timestamp() time.Time {
	return b.TimestampV
}

func (b *Block) Verify(context.Context) error {
	return b.VerifyV
}

func (b *Block) Bytes() []byte {
	return b.BytesV
}

func (b *Block) Compare(other *Block) int {
	return cmp.Compare(b.HeightV, other.HeightV)
}
