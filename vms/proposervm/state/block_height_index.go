// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"github.com/skychains/chain/cache"
	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/prefixdb"
	"github.com/skychains/chain/database/versiondb"
	"github.com/skychains/chain/ids"
)

const cacheSize = 8192 // max cache entries

var (
	_ HeightIndex = (*heightIndex)(nil)

	heightPrefix   = []byte("height")
	metadataPrefix = []byte("metadata")

	forkKey = []byte("fork")
)

type HeightIndexGetter interface {
	// GetMinimumHeight return the smallest height of an indexed blockID. If
	// there are no indexed blockIDs, ErrNotFound will be returned.
	GetMinimumHeight() (uint64, error)
	GetBlockIDAtHeight(height uint64) (ids.ID, error)

	// Fork height is stored when the first post-fork block/option is accepted.
	// Before that, fork height won't be found.
	GetForkHeight() (uint64, error)
}

type HeightIndexWriter interface {
	SetForkHeight(height uint64) error
	SetBlockIDAtHeight(height uint64, blkID ids.ID) error
	DeleteBlockIDAtHeight(height uint64) error
}

// HeightIndex contains mapping of blockHeights to accepted proposer block IDs
// along with some metadata (fork height and checkpoint).
type HeightIndex interface {
	HeightIndexWriter
	HeightIndexGetter
}

type heightIndex struct {
	versiondb.Commitable

	// Caches block height -> proposerVMBlockID.
	heightsCache cache.Cacher[uint64, ids.ID]

	heightDB   database.Database
	metadataDB database.Database
}

func NewHeightIndex(db database.Database, commitable versiondb.Commitable) HeightIndex {
	return &heightIndex{
		Commitable: commitable,

		heightsCache: &cache.LRU[uint64, ids.ID]{Size: cacheSize},
		heightDB:     prefixdb.New(heightPrefix, db),
		metadataDB:   prefixdb.New(metadataPrefix, db),
	}
}

func (hi *heightIndex) GetMinimumHeight() (uint64, error) {
	it := hi.heightDB.NewIterator()
	defer it.Release()

	if !it.Next() {
		return 0, database.ErrNotFound
	}

	height, err := database.ParseUInt64(it.Key())
	if err != nil {
		return 0, err
	}
	return height, it.Error()
}

func (hi *heightIndex) GetBlockIDAtHeight(height uint64) (ids.ID, error) {
	if blkID, found := hi.heightsCache.Get(height); found {
		return blkID, nil
	}

	key := database.PackUInt64(height)
	blkID, err := database.GetID(hi.heightDB, key)
	if err != nil {
		return ids.Empty, err
	}
	hi.heightsCache.Put(height, blkID)
	return blkID, err
}

func (hi *heightIndex) SetBlockIDAtHeight(height uint64, blkID ids.ID) error {
	hi.heightsCache.Put(height, blkID)
	key := database.PackUInt64(height)
	return database.PutID(hi.heightDB, key, blkID)
}

func (hi *heightIndex) DeleteBlockIDAtHeight(height uint64) error {
	hi.heightsCache.Evict(height)
	key := database.PackUInt64(height)
	return hi.heightDB.Delete(key)
}

func (hi *heightIndex) GetForkHeight() (uint64, error) {
	return database.GetUInt64(hi.metadataDB, forkKey)
}

func (hi *heightIndex) SetForkHeight(height uint64) error {
	return database.PutUInt64(hi.metadataDB, forkKey, height)
}
