// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package timeout

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/networking/benchlist"
	"github.com/skychains/chain/utils/timer"
)

func TestManagerFire(t *testing.T) {
	benchlist := benchlist.NewNoBenchlist()
	manager, err := NewManager(
		&timer.AdaptiveTimeoutConfig{
			InitialTimeout:     time.Millisecond,
			MinimumTimeout:     time.Millisecond,
			MaximumTimeout:     10 * time.Second,
			TimeoutCoefficient: 1.25,
			TimeoutHalflife:    5 * time.Minute,
		},
		benchlist,
		prometheus.NewRegistry(),
		prometheus.NewRegistry(),
	)
	require.NoError(t, err)
	go manager.Dispatch()
	defer manager.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)

	manager.RegisterRequest(
		ids.EmptyNodeID,
		ids.Empty,
		true,
		ids.RequestID{},
		wg.Done,
	)

	wg.Wait()
}
