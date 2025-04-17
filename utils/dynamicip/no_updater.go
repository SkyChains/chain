// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package dynamicip

import "github.com/SkyChains/chain/utils/logging"

var _ Updater = noUpdater{}

func NewNoUpdater() Updater {
	return noUpdater{}
}

type noUpdater struct{}

func (noUpdater) Dispatch(logging.Logger) {}

func (noUpdater) Stop() {}
