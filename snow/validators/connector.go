// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"context"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/version"
)

// Connector represents a handler that is called when a connection is marked as
// connected or disconnected
type Connector interface {
	Connected(
		ctx context.Context,
		nodeID ids.NodeID,
		nodeVersion *version.Application,
	) error
	Disconnected(ctx context.Context, nodeID ids.NodeID) error
}
