// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package subprocess

import (
	"context"
	"os/exec"
	"sync"

	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms/rpcchainvm/runtime"
)

func NewStopper(logger logging.Logger, cmd *exec.Cmd) runtime.Stopper {
	return &stopper{
		cmd:    cmd,
		logger: logger,
	}
}

type stopper struct {
	once   sync.Once
	cmd    *exec.Cmd
	logger logging.Logger
}

func (s *stopper) Stop(ctx context.Context) {
	s.once.Do(func() {
		stop(ctx, s.logger, s.cmd)
	})
}
