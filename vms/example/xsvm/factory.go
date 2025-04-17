// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package xsvm

import (
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms"
)

var _ vms.Factory = (*Factory)(nil)

type Factory struct{}

func (*Factory) New(logging.Logger) (interface{}, error) {
	return &VM{}, nil
}
