// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/skychains/chain/utils"
)

type gossipTrackerMetrics struct {
	trackedPeersSize prometheus.Gauge
	validatorsSize   prometheus.Gauge
}

func newGossipTrackerMetrics(registerer prometheus.Registerer, namespace string) (gossipTrackerMetrics, error) {
	m := gossipTrackerMetrics{
		trackedPeersSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "tracked_peers_size",
				Help:      "amount of peers that are being tracked",
			},
		),
		validatorsSize: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "validators_size",
				Help:      "number of validators this node is tracking",
			},
		),
	}

	err := utils.Err(
		registerer.Register(m.trackedPeersSize),
		registerer.Register(m.validatorsSize),
	)
	return m, err
}
