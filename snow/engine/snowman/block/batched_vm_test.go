// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/consensus/snowman"
	"github.com/skychains/chain/utils/logging"
)

var errTest = errors.New("non-nil error")

func TestGetAncestorsDatabaseNotFound(t *testing.T) {
	require := require.New(t)

	vm := &TestVM{}
	someID := ids.GenerateTestID()
	vm.GetBlockF = func(_ context.Context, id ids.ID) (snowman.Block, error) {
		require.Equal(someID, id)
		return nil, database.ErrNotFound
	}
	containers, err := GetAncestors(context.Background(), logging.NoLog{}, vm, someID, 10, 10, 1*time.Second)
	require.NoError(err)
	require.Empty(containers)
}

// TestGetAncestorsPropagatesErrors checks errors other than
// database.ErrNotFound propagate to caller.
func TestGetAncestorsPropagatesErrors(t *testing.T) {
	require := require.New(t)

	vm := &TestVM{}
	someID := ids.GenerateTestID()
	vm.GetBlockF = func(_ context.Context, id ids.ID) (snowman.Block, error) {
		require.Equal(someID, id)
		return nil, errTest
	}
	containers, err := GetAncestors(context.Background(), logging.NoLog{}, vm, someID, 10, 10, 1*time.Second)
	require.Nil(containers)
	require.ErrorIs(err, errTest)
}
