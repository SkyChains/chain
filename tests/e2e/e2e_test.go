// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e_test

import (
	"testing"

	// ensure test packages are scanned by ginkgo
	_ "github.com/skychains/chain/tests/e2e/banff"
	_ "github.com/skychains/chain/tests/e2e/c"
	_ "github.com/skychains/chain/tests/e2e/faultinjection"
	_ "github.com/skychains/chain/tests/e2e/p"
	_ "github.com/skychains/chain/tests/e2e/x"
	_ "github.com/skychains/chain/tests/e2e/x/transfer"

	"github.com/skychains/chain/tests/e2e/vms"
	"github.com/skychains/chain/tests/fixture/e2e"
	"github.com/skychains/chain/tests/fixture/tmpnet"

	ginkgo "github.com/onsi/ginkgo/v2"
)

func TestE2E(t *testing.T) {
	ginkgo.RunSpecs(t, "e2e test suites")
}

var flagVars *e2e.FlagVars

func init() {
	flagVars = e2e.RegisterFlags()
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Run only once in the first ginkgo process

	nodes := tmpnet.NewNodesOrPanic(flagVars.NodeCount())
	subnets := vms.XSVMSubnetsOrPanic(nodes...)
	return e2e.NewTestEnvironment(
		flagVars,
		&tmpnet.Network{
			Owner:   "node-e2e",
			Nodes:   nodes,
			Subnets: subnets,
		},
	).Marshal()
}, func(envBytes []byte) {
	// Run in every ginkgo process

	// Initialize the local test environment from the global state
	e2e.InitSharedTestEnvironment(envBytes)
})
