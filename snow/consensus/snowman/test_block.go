// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"context"
	"time"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/choices"
	"github.com/skychains/chain/utils"
)

var (
	_ Block                      = (*TestBlock)(nil)
	_ utils.Sortable[*TestBlock] = (*TestBlock)(nil)
)

// TestBlock is a useful test block
type TestBlock struct {
	choices.TestDecidable

	ParentV    ids.ID
	HeightV    uint64
	TimestampV time.Time
	VerifyV    error
	BytesV     []byte
}

func (b *TestBlock) Parent() ids.ID {
	return b.ParentV
}

func (b *TestBlock) Height() uint64 {
	return b.HeightV
}

func (b *TestBlock) Timestamp() time.Time {
	return b.TimestampV
}

func (b *TestBlock) Verify(context.Context) error {
	return b.VerifyV
}

func (b *TestBlock) Bytes() []byte {
	return b.BytesV
}

func (b *TestBlock) Less(other *TestBlock) bool {
	return b.HeightV < other.HeightV
}
