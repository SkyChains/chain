// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"context"
	"fmt"
	"time"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/utils/constants"
)

func (vm *VM) HealthCheck(context.Context) (interface{}, error) {
	localPrimaryValidator, err := vm.state.GetCurrentValidator(
		constants.PrimaryNetworkID,
		vm.ctx.NodeID,
	)
	switch err {
	case nil:
		vm.metrics.SetTimeUntilUnstake(time.Until(localPrimaryValidator.EndTime))
	case database.ErrNotFound:
		vm.metrics.SetTimeUntilUnstake(0)
	default:
		return nil, fmt.Errorf("couldn't get current local validator: %w", err)
	}

	for subnetID := range vm.TrackedSubnets {
		localSubnetValidator, err := vm.state.GetCurrentValidator(
			subnetID,
			vm.ctx.NodeID,
		)
		switch err {
		case nil:
			vm.metrics.SetTimeUntilSubnetUnstake(subnetID, time.Until(localSubnetValidator.EndTime))
		case database.ErrNotFound:
			vm.metrics.SetTimeUntilSubnetUnstake(subnetID, 0)
		default:
			return nil, fmt.Errorf("couldn't get current subnet validator of %q: %w", subnetID, err)
		}
	}
	return nil, nil
}
