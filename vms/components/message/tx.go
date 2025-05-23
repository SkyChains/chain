// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import "github.com/skychains/chain/ids"

var _ Message = (*Tx)(nil)

type Tx struct {
	message

	Tx []byte `serialize:"true"`
}

func (msg *Tx) Handle(handler Handler, nodeID ids.NodeID, requestID uint32) error {
	return handler.HandleTx(nodeID, requestID, msg)
}
