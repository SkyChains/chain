// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package handler

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/set"
)

func getIDs(idsBytes [][]byte) (set.Set[ids.ID], error) {
	var res set.Set[ids.ID]
	for _, bytes := range idsBytes {
		id, err := ids.ToID(bytes)
		if err != nil {
			return nil, err
		}
		res.Add(id)
	}
	return res, nil
}
