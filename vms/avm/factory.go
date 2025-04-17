// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms"
	"github.com/skychains/chain/vms/avm/config"
)

var _ vms.Factory = (*Factory)(nil)

type Factory struct {
	config.Config
}

func (f *Factory) New(logging.Logger) (interface{}, error) {
	return &VM{Config: f.Config}, nil
}
