// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/uptime"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/timer/mockable"
	"github.com/skychains/chain/vms/platformvm/config"
	"github.com/skychains/chain/vms/platformvm/fx"
	"github.com/skychains/chain/vms/platformvm/reward"
	"github.com/skychains/chain/vms/platformvm/utxo"
)

type Backend struct {
	Config       *config.Config
	Ctx          *snow.Context
	Clk          *mockable.Clock
	Fx           fx.Fx
	FlowChecker  utxo.Verifier
	Uptimes      uptime.Calculator
	Rewards      reward.Calculator
	Bootstrapped *utils.Atomic[bool]
}
