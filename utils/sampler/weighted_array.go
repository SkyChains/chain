// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package sampler

import (
	"cmp"

	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/math"
)

var (
	_ Weighted                             = (*weightedArray)(nil)
	_ utils.Sortable[weightedArrayElement] = weightedArrayElement{}
)

type weightedArrayElement struct {
	cumulativeWeight uint64
	index            int
}

// Note that this sorts in order of decreasing weight.
func (e weightedArrayElement) Compare(other weightedArrayElement) int {
	return cmp.Compare(other.cumulativeWeight, e.cumulativeWeight)
}

// Sampling is performed by executing a modified binary search over the provided
// elements. Rather than cutting the remaining dataset in half, the algorithm
// attempt to just in to where it think the value will be assuming a linear
// distribution of the element weights.
//
// Initialization takes O(n * log(n)) time, where n is the number of elements
// that can be sampled.
// Sampling can take up to O(n) time. If the distribution is linearly
// distributed, then the runtime is constant.
type weightedArray struct {
	arr []weightedArrayElement
}

func (s *weightedArray) Initialize(weights []uint64) error {
	numWeights := len(weights)
	if numWeights <= cap(s.arr) {
		s.arr = s.arr[:numWeights]
	} else {
		s.arr = make([]weightedArrayElement, numWeights)
	}

	for i, weight := range weights {
		s.arr[i] = weightedArrayElement{
			cumulativeWeight: weight,
			index:            i,
		}
	}

	// Optimize so that the array is closer to the uniform distribution
	utils.Sort(s.arr)

	maxIndex := len(s.arr) - 1
	oneIfOdd := 1 & maxIndex
	oneIfEven := 1 - oneIfOdd
	end := maxIndex - oneIfEven
	for i := 1; i < end; i += 2 {
		s.arr[i], s.arr[end] = s.arr[end], s.arr[i]
		end -= 2
	}

	cumulativeWeight := uint64(0)
	for i := 0; i < len(s.arr); i++ {
		newWeight, err := math.Add64(
			cumulativeWeight,
			s.arr[i].cumulativeWeight,
		)
		if err != nil {
			return err
		}
		cumulativeWeight = newWeight
		s.arr[i].cumulativeWeight = cumulativeWeight
	}

	return nil
}

func (s *weightedArray) Sample(value uint64) (int, bool) {
	if len(s.arr) == 0 || s.arr[len(s.arr)-1].cumulativeWeight <= value {
		return 0, false
	}
	minIndex := 0
	maxIndex := len(s.arr) - 1
	maxCumulativeWeight := float64(s.arr[len(s.arr)-1].cumulativeWeight)
	index := int((float64(value) * float64(maxIndex+1)) / maxCumulativeWeight)

	for {
		previousWeight := uint64(0)
		if index > 0 {
			previousWeight = s.arr[index-1].cumulativeWeight
		}
		currentElem := s.arr[index]
		currentWeight := currentElem.cumulativeWeight
		if previousWeight <= value && value < currentWeight {
			return currentElem.index, true
		}

		if value < previousWeight {
			// go to the left
			maxIndex = index - 1
		} else {
			// go to the right
			minIndex = index + 1
		}

		minWeight := uint64(0)
		if minIndex > 0 {
			minWeight = s.arr[minIndex-1].cumulativeWeight
		}
		maxWeight := s.arr[maxIndex].cumulativeWeight

		valueRange := maxWeight - minWeight
		adjustedLookupValue := value - minWeight
		indexRange := maxIndex - minIndex + 1
		lookupMass := float64(adjustedLookupValue) * float64(indexRange)

		index = int(lookupMass/float64(valueRange)) + minIndex
	}
}
