// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/networking/router"
	"github.com/skychains/chain/version"
)

var _ router.ExternalHandler = (*testHandler)(nil)

type testHandler struct {
	router.InboundHandler
	ConnectedF    func(nodeID ids.NodeID, nodeVersion *version.Application, subnetID ids.ID)
	DisconnectedF func(nodeID ids.NodeID)
}

func (h *testHandler) Connected(id ids.NodeID, nodeVersion *version.Application, subnetID ids.ID) {
	if h.ConnectedF != nil {
		h.ConnectedF(id, nodeVersion, subnetID)
	}
}

func (h *testHandler) Disconnected(id ids.NodeID) {
	if h.DisconnectedF != nil {
		h.DisconnectedF(id)
	}
}
