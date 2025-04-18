// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	dto "github.com/prometheus/client_model/go"
)

var _ OptionalGatherer = (*optionalGatherer)(nil)

// OptionalGatherer extends the Gatherer interface by allowing the optional
// registration of a single gatherer. If no gatherer is registered, Gather will
// return no metrics and no error. If a gatherer is registered, Gather will
// return the results of calling Gather on the provided gatherer.
type OptionalGatherer interface {
	prometheus.Gatherer

	// Register the provided gatherer. If a gatherer was previously registered,
	// an error will be returned.
	Register(gatherer prometheus.Gatherer) error
}

type optionalGatherer struct {
	lock     sync.RWMutex
	gatherer prometheus.Gatherer
}

func NewOptionalGatherer() OptionalGatherer {
	return &optionalGatherer{}
}

func (g *optionalGatherer) Gather() ([]*dto.MetricFamily, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	if g.gatherer == nil {
		return nil, nil
	}
	return g.gatherer.Gather()
}

func (g *optionalGatherer) Register(gatherer prometheus.Gatherer) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.gatherer != nil {
		return fmt.Errorf("%w; existing: %#v; new: %#v",
			errReregisterGatherer,
			g.gatherer,
			gatherer,
		)
	}
	g.gatherer = gatherer
	return nil
}
