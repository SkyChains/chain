// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements ping tests, requires network-runner cluster.
package ping

import (
	"context"
	"time"

	"github.com/luxdefi/luxd/tests/e2e"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[Ping]", func() {
	ginkgo.It("can ping network-runner RPC server",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"require-network-runner",
			"ping",
		),
		func() {
			if e2e.Env.GetRunnerGRPCEndpoint() == "" {
				ginkgo.Skip("no local network-runner, skipping")
			}

			runnerCli := e2e.Env.GetRunnerClient()
			gomega.Expect(runnerCli).ShouldNot(gomega.BeNil())

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			_, err := runnerCli.Ping(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())
		})
})
