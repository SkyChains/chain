// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package execute

import (
	"github.com/skychains/chain/database"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/vms/example/xsvm/block"
	"github.com/skychains/chain/vms/example/xsvm/genesis"
	"github.com/skychains/chain/vms/example/xsvm/state"
)

func Genesis(db database.KeyValueReaderWriterDeleter, chainID ids.ID, g *genesis.Genesis) error {
	isInitialized, err := state.IsInitialized(db)
	if err != nil {
		return err
	}
	if isInitialized {
		return nil
	}

	blk, err := genesis.Block(g)
	if err != nil {
		return err
	}

	for _, allocation := range g.Allocations {
		if err := state.SetBalance(db, allocation.Address, chainID, allocation.Balance); err != nil {
			return err
		}
	}

	blkID, err := blk.ID()
	if err != nil {
		return err
	}

	blkBytes, err := block.Codec.Marshal(block.CodecVersion, blk)
	if err != nil {
		return err
	}

	if err := state.AddBlock(db, blk.Height, blkID, blkBytes); err != nil {
		return err
	}
	if err := state.SetLastAccepted(db, blkID); err != nil {
		return err
	}
	return state.SetInitialized(db)
}
