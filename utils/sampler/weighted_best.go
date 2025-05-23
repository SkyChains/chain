// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package sampler

import (
	"errors"
	"math"
	"time"

	"github.com/skychains/chain/utils/timer/mockable"

	safemath "github.com/skychains/chain/utils/math"
)

var (
	errNoValidWeightedSamplers = errors.New("no valid weighted samplers found")

	_ Weighted = (*weightedBest)(nil)
)

// Sampling is performed by using another implementation of the Weighted
// interface.
//
// Initialization attempts to find the best sampling algorithm given the dataset
// by performing a benchmark of the provided implementations.
type weightedBest struct {
	Weighted
	samplers            []Weighted
	benchmarkIterations int
	clock               mockable.Clock
}

func (s *weightedBest) Initialize(weights []uint64) error {
	totalWeight := uint64(0)
	for _, weight := range weights {
		newWeight, err := safemath.Add64(totalWeight, weight)
		if err != nil {
			return err
		}
		totalWeight = newWeight
	}

	samples := []uint64(nil)
	if totalWeight > 0 {
		samples = make([]uint64, s.benchmarkIterations)
		for i := range samples {
			samples[i] = globalRNG.Uint64Inclusive(totalWeight - 1)
		}
	}

	s.Weighted = nil
	bestDuration := time.Duration(math.MaxInt64)

samplerLoop:
	for _, sampler := range s.samplers {
		if err := sampler.Initialize(weights); err != nil {
			continue
		}

		start := s.clock.Time()
		for _, sample := range samples {
			if _, ok := sampler.Sample(sample); !ok {
				continue samplerLoop
			}
		}
		end := s.clock.Time()
		duration := end.Sub(start)
		if duration < bestDuration {
			bestDuration = duration
			s.Weighted = sampler
		}
	}

	if s.Weighted == nil {
		return errNoValidWeightedSamplers
	}
	return nil
}
