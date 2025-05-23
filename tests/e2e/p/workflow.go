// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build test

package p

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/api/info"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/tests"
	"github.com/skychains/chain/tests/fixture/e2e"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/bls"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/platformvm"
	"github.com/skychains/chain/vms/platformvm/signer"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/secp256k1fx"

	ginkgo "github.com/onsi/ginkgo/v2"
)

// PChainWorkflow is an integration test for normal P-Chain operations
// - Issues an Add Validator and an Add Delegator using the funding address
// - Exports LUX from the P-Chain funding address to the X-Chain created address
// - Exports LUX from the X-Chain created address to the P-Chain created address
// - Checks the expected value of the funding address

var _ = e2e.DescribePChain("[Workflow]", func() {
	require := require.New(ginkgo.GinkgoT())

	ginkgo.It("P-chain main operations",
		func() {
			nodeURI := e2e.Env.GetRandomNodeURI()
			keychain := e2e.Env.NewKeychain(2)
			baseWallet := e2e.NewWallet(keychain, nodeURI)

			pWallet := baseWallet.P()
			pBuilder := pWallet.Builder()
			pContext := pBuilder.Context()
			luxAssetID := pContext.LUXAssetID
			xWallet := baseWallet.X()
			xBuilder := xWallet.Builder()
			xContext := xBuilder.Context()
			pChainClient := platformvm.NewClient(nodeURI.URI)

			tests.Outf("{{blue}} fetching minimal stake amounts {{/}}\n")
			minValStake, minDelStake, err := pChainClient.GetMinStake(e2e.DefaultContext(), constants.PlatformChainID)
			require.NoError(err)
			tests.Outf("{{green}} minimal validator stake: %d {{/}}\n", minValStake)
			tests.Outf("{{green}} minimal delegator stake: %d {{/}}\n", minDelStake)

			tests.Outf("{{blue}} fetching tx fee {{/}}\n")
			infoClient := info.NewClient(nodeURI.URI)
			fees, err := infoClient.GetTxFee(e2e.DefaultContext())
			require.NoError(err)
			txFees := uint64(fees.TxFee)
			tests.Outf("{{green}} txFee: %d {{/}}\n", txFees)

			// amount to transfer from P to X chain
			toTransfer := 1 * units.Lux

			pShortAddr := keychain.Keys[0].Address()
			xTargetAddr := keychain.Keys[1].Address()
			ginkgo.By("check selected keys have sufficient funds", func() {
				pBalances, err := pWallet.Builder().GetBalance()
				pBalance := pBalances[luxAssetID]
				minBalance := minValStake + txFees + minDelStake + txFees + toTransfer + txFees
				require.NoError(err)
				require.GreaterOrEqual(pBalance, minBalance)
			})

			// Use a random node ID to ensure that repeated test runs
			// will succeed against a network that persists across runs.
			validatorID, err := ids.ToNodeID(utils.RandomBytes(ids.NodeIDLen))
			require.NoError(err)

			vdr := &txs.SubnetValidator{
				Validator: txs.Validator{
					NodeID: validatorID,
					End:    uint64(time.Now().Add(72 * time.Hour).Unix()),
					Wght:   minValStake,
				},
				Subnet: constants.PrimaryNetworkID,
			}
			rewardOwner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{pShortAddr},
			}
			shares := uint32(20000) // TODO: retrieve programmatically

			sk, err := bls.NewSecretKey()
			require.NoError(err)
			pop := signer.NewProofOfPossession(sk)

			ginkgo.By("issue add validator tx", func() {
				_, err := pWallet.IssueAddPermissionlessValidatorTx(
					vdr,
					pop,
					luxAssetID,
					rewardOwner,
					rewardOwner,
					shares,
					e2e.WithDefaultContext(),
				)
				require.NoError(err)
			})

			ginkgo.By("issue add delegator tx", func() {
				_, err := pWallet.IssueAddPermissionlessDelegatorTx(
					vdr,
					luxAssetID,
					rewardOwner,
					e2e.WithDefaultContext(),
				)
				require.NoError(err)
			})

			// retrieve initial balances
			pBalances, err := pWallet.Builder().GetBalance()
			require.NoError(err)
			pStartBalance := pBalances[luxAssetID]
			tests.Outf("{{blue}} P-chain balance before P->X export: %d {{/}}\n", pStartBalance)

			xBalances, err := xWallet.Builder().GetFTBalance()
			require.NoError(err)
			xStartBalance := xBalances[luxAssetID]
			tests.Outf("{{blue}} X-chain balance before P->X export: %d {{/}}\n", xStartBalance)

			outputOwner := secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					xTargetAddr,
				},
			}
			output := &secp256k1fx.TransferOutput{
				Amt:          toTransfer,
				OutputOwners: outputOwner,
			}

			ginkgo.By("export lux from P to X chain", func() {
				_, err := pWallet.IssueExportTx(
					xContext.BlockchainID,
					[]*lux.TransferableOutput{
						{
							Asset: lux.Asset{
								ID: luxAssetID,
							},
							Out: output,
						},
					},
					e2e.WithDefaultContext(),
				)
				require.NoError(err)
			})

			// check balances post export
			pBalances, err = pWallet.Builder().GetBalance()
			require.NoError(err)
			pPreImportBalance := pBalances[luxAssetID]
			tests.Outf("{{blue}} P-chain balance after P->X export: %d {{/}}\n", pPreImportBalance)

			xBalances, err = xWallet.Builder().GetFTBalance()
			require.NoError(err)
			xPreImportBalance := xBalances[luxAssetID]
			tests.Outf("{{blue}} X-chain balance after P->X export: %d {{/}}\n", xPreImportBalance)

			require.Equal(xPreImportBalance, xStartBalance) // import not performed yet
			require.Equal(pPreImportBalance, pStartBalance-toTransfer-txFees)

			ginkgo.By("import lux from P into X chain", func() {
				_, err := xWallet.IssueImportTx(
					constants.PlatformChainID,
					&outputOwner,
					e2e.WithDefaultContext(),
				)
				require.NoError(err)
			})

			// check balances post import
			pBalances, err = pWallet.Builder().GetBalance()
			require.NoError(err)
			pFinalBalance := pBalances[luxAssetID]
			tests.Outf("{{blue}} P-chain balance after P->X import: %d {{/}}\n", pFinalBalance)

			xBalances, err = xWallet.Builder().GetFTBalance()
			require.NoError(err)
			xFinalBalance := xBalances[luxAssetID]
			tests.Outf("{{blue}} X-chain balance after P->X import: %d {{/}}\n", xFinalBalance)

			require.Equal(xFinalBalance, xPreImportBalance+toTransfer-txFees) // import not performed yet
			require.Equal(pFinalBalance, pPreImportBalance)
		})
})
