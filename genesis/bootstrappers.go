// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"encoding/json"
	"fmt"

	_ "embed"

	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/utils/constants"
	"github.com/luxdefi/node/utils/ips"
	"github.com/luxdefi/node/utils/math"
	"github.com/luxdefi/node/utils/sampler"
)

var (
	//go:embed bootstrappers.json
	bootstrappersPerNetworkJSON []byte

	bootstrappersPerNetwork map[string][]Bootstrapper
)

func init() {
	if err := json.Unmarshal(bootstrappersPerNetworkJSON, &bootstrappersPerNetwork); err != nil {
		panic(fmt.Sprintf("failed to decode bootstrappers.json %v", err))
	}
}

// Represents the relationship between the nodeID and the nodeIP.
// The bootstrapper is sometimes called "anchor" or "beacon" node.
type Bootstrapper struct {
	ID ids.NodeID `json:"id"`
	IP ips.IPDesc `json:"ip"`
}

// SampleBootstrappers returns the some beacons this node should connect to
func SampleBootstrappers(networkID uint32, count int) []Bootstrapper {
	networkName := constants.NetworkIDToNetworkName[networkID]
	bootstrappers := bootstrappersPerNetwork[networkName]
	count = math.Min(count, len(bootstrappers))

	s := sampler.NewUniform()
	s.Initialize(uint64(len(bootstrappers)))
	indices, _ := s.Sample(count)

	sampled := make([]Bootstrapper, 0, len(indices))
	for _, index := range indices {
		sampled = append(sampled, bootstrappers[int(index)])
	}
	return sampled
}
