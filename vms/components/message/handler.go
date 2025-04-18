// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"go.uber.org/zap"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/logging"
)

var _ Handler = NoopHandler{}

type Handler interface {
	HandleTx(nodeID ids.NodeID, requestID uint32, msg *Tx) error
}

type NoopHandler struct {
	Log logging.Logger
}

func (h NoopHandler) HandleTx(nodeID ids.NodeID, requestID uint32, _ *Tx) error {
	h.Log.Debug("dropping unexpected Tx message",
		zap.Stringer("nodeID", nodeID),
		zap.Uint32("requestID", requestID),
	)
	return nil
}
