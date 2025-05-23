// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build test

// Implements tests for the banff network upgrade.
package banff

import (
	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/tests"
	"github.com/skychains/chain/tests/fixture/e2e"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/secp256k1fx"

	ginkgo "github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("[Banff]", func() {
	require := require.New(ginkgo.GinkgoT())

	ginkgo.It("can send custom assets X->P and P->X",
		func() {
			keychain := e2e.Env.NewKeychain(1)
			wallet := e2e.NewWallet(keychain, e2e.Env.GetRandomNodeURI())

			// Get the P-chain and the X-chain wallets
			pWallet := wallet.P()
			xWallet := wallet.X()
			xBuilder := xWallet.Builder()
			xContext := xBuilder.Context()

			// Pull out useful constants to use when issuing transactions.
			xChainID := xContext.BlockchainID
			owner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					keychain.Keys[0].Address(),
				},
			}

			var assetID ids.ID
			ginkgo.By("create new X-chain asset", func() {
				assetTx, err := xWallet.IssueCreateAssetTx(
					"RnM",
					"RNM",
					9,
					map[uint32][]verify.State{
						0: {
							&secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
					e2e.WithDefaultContext(),
				)
				require.NoError(err)
				assetID = assetTx.ID()

				tests.Outf("{{green}}created new X-chain asset{{/}}: %s\n", assetID)
			})

			ginkgo.By("export new X-chain asset to P-chain", func() {
				tx, err := xWallet.IssueExportTx(
					constants.PlatformChainID,
					[]*lux.TransferableOutput{
						{
							Asset: lux.Asset{
								ID: assetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
					e2e.WithDefaultContext(),
				)
				require.NoError(err)

				tests.Outf("{{green}}issued X-chain export{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("import new asset from X-chain on the P-chain", func() {
				tx, err := pWallet.IssueImportTx(
					xChainID,
					owner,
					e2e.WithDefaultContext(),
				)
				require.NoError(err)

				tests.Outf("{{green}}issued P-chain import{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("export asset from P-chain to the X-chain", func() {
				tx, err := pWallet.IssueExportTx(
					xChainID,
					[]*lux.TransferableOutput{
						{
							Asset: lux.Asset{
								ID: assetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
					e2e.WithDefaultContext(),
				)
				require.NoError(err)

				tests.Outf("{{green}}issued P-chain export{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("import asset from P-chain on the X-chain", func() {
				tx, err := xWallet.IssueImportTx(
					constants.PlatformChainID,
					owner,
					e2e.WithDefaultContext(),
				)
				require.NoError(err)

				tests.Outf("{{green}}issued X-chain import{{/}}: %s\n", tx.ID())
			})
		})
})
