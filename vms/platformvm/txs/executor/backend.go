// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"github.com/SkyChains/chain/snow"
	"github.com/SkyChains/chain/snow/uptime"
	"github.com/SkyChains/chain/utils"
	"github.com/SkyChains/chain/utils/timer/mockable"
	"github.com/SkyChains/chain/vms/platformvm/config"
	"github.com/SkyChains/chain/vms/platformvm/fx"
	"github.com/SkyChains/chain/vms/platformvm/reward"
	"github.com/SkyChains/chain/vms/platformvm/utxo"
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
