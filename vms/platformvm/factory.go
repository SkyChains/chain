// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms"
	"github.com/skychains/chain/vms/platformvm/config"
)

var _ vms.Factory = (*Factory)(nil)

// Factory can create new instances of the Platform Chain
type Factory struct {
	config.Config
}

// New returns a new instance of the Platform Chain
func (f *Factory) New(logging.Logger) (interface{}, error) {
	return &VM{Config: f.Config}, nil
}
