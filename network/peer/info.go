// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"net/netip"
	"time"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/json"
	"github.com/skychains/chain/utils/set"
)

type Info struct {
	IP                    netip.AddrPort         `json:"ip"`
	PublicIP              netip.AddrPort         `json:"publicIP,omitempty"`
	ID                    ids.NodeID             `json:"nodeID"`
	Version               string                 `json:"version"`
	LastSent              time.Time              `json:"lastSent"`
	LastReceived          time.Time              `json:"lastReceived"`
	ObservedUptime        json.Uint32            `json:"observedUptime"`
	ObservedSubnetUptimes map[ids.ID]json.Uint32 `json:"observedSubnetUptimes"`
	TrackedSubnets        set.Set[ids.ID]        `json:"trackedSubnets"`
	SupportedACPs         set.Set[uint32]        `json:"supportedACPs"`
	ObjectedACPs          set.Set[uint32]        `json:"objectedACPs"`
}
