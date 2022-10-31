// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package beacon

import (
	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/utils/ips"
)

var _ Beacon = (*beacon)(nil)

type Beacon interface {
	ID() ids.NodeID
	IP() ips.IPPort
}

type beacon struct {
	id ids.NodeID
	ip ips.IPPort
}

func New(id ids.NodeID, ip ips.IPPort) Beacon {
	return &beacon{
		id: id,
		ip: ip,
	}
}

func (b *beacon) ID() ids.NodeID { return b.id }
func (b *beacon) IP() ips.IPPort { return b.ip }
