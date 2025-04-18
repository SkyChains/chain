// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements static handlers tests for avm and platformvm
package statichandlers

import (
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/tests/fixture/e2e"
	"github.com/skychains/chain/utils/cb58"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/secp256k1"
	"github.com/skychains/chain/utils/formatting"
	"github.com/skychains/chain/utils/formatting/address"
	"github.com/skychains/chain/utils/json"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/avm"
	"github.com/skychains/chain/vms/platformvm/api"
	"github.com/skychains/chain/vms/platformvm/reward"
)

var _ = ginkgo.Describe("[StaticHandlers]", func() {
	require := require.New(ginkgo.GinkgoT())

	ginkgo.It("can make calls to avm static api",
		func() {
			addrMap := map[string]string{}
			for _, addrStr := range []string{
				"A9bTQjfYGBFK3JPRJqF2eh3JYL7cHocvy",
				"6mxBGnjGDCKgkVe7yfrmvMA7xE7qCv3vv",
				"6ncQ19Q2U4MamkCYzshhD8XFjfwAWFzTa",
				"Jz9ayEDt7dx9hDx45aXALujWmL9ZUuqe7",
			} {
				addr, err := ids.ShortFromString(addrStr)
				require.NoError(err)
				addrMap[addrStr], err = address.FormatBech32(constants.NetworkIDToHRP[constants.LocalID], addr[:])
				require.NoError(err)
			}
			avmArgs := avm.BuildGenesisArgs{
				Encoding: formatting.Hex,
				GenesisData: map[string]avm.AssetDefinition{
					"asset1": {
						Name:         "myFixedCapAsset",
						Symbol:       "MFCA",
						Denomination: 8,
						InitialState: map[string][]interface{}{
							"fixedCap": {
								avm.Holder{
									Amount:  100000,
									Address: addrMap["A9bTQjfYGBFK3JPRJqF2eh3JYL7cHocvy"],
								},
								avm.Holder{
									Amount:  100000,
									Address: addrMap["6mxBGnjGDCKgkVe7yfrmvMA7xE7qCv3vv"],
								},
								avm.Holder{
									Amount:  json.Uint64(50000),
									Address: addrMap["6ncQ19Q2U4MamkCYzshhD8XFjfwAWFzTa"],
								},
								avm.Holder{
									Amount:  json.Uint64(50000),
									Address: addrMap["Jz9ayEDt7dx9hDx45aXALujWmL9ZUuqe7"],
								},
							},
						},
					},
					"asset2": {
						Name:   "myVarCapAsset",
						Symbol: "MVCA",
						InitialState: map[string][]interface{}{
							"variableCap": {
								avm.Owners{
									Threshold: 1,
									Minters: []string{
										addrMap["A9bTQjfYGBFK3JPRJqF2eh3JYL7cHocvy"],
										addrMap["6mxBGnjGDCKgkVe7yfrmvMA7xE7qCv3vv"],
									},
								},
								avm.Owners{
									Threshold: 2,
									Minters: []string{
										addrMap["6ncQ19Q2U4MamkCYzshhD8XFjfwAWFzTa"],
										addrMap["Jz9ayEDt7dx9hDx45aXALujWmL9ZUuqe7"],
									},
								},
							},
						},
					},
					"asset3": {
						Name: "myOtherVarCapAsset",
						InitialState: map[string][]interface{}{
							"variableCap": {
								avm.Owners{
									Threshold: 1,
									Minters: []string{
										addrMap["A9bTQjfYGBFK3JPRJqF2eh3JYL7cHocvy"],
									},
								},
							},
						},
					},
				},
			}
			staticClient := avm.NewStaticClient(e2e.Env.GetRandomNodeURI().URI)
			resp, err := staticClient.BuildGenesis(e2e.DefaultContext(), &avmArgs)
			require.NoError(err)
			require.Equal(resp.Bytes, "0x0000000000030006617373657431000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f6d794669786564436170417373657400044d4643410800000001000000000000000400000007000000000000c350000000000000000000000001000000013f78e510df62bc48b0829ec06d6a6b98062d695300000007000000000000c35000000000000000000000000100000001c54903de5177a16f7811771ef2f4659d9e8646710000000700000000000186a0000000000000000000000001000000013f58fda2e9ea8d9e4b181832a07b26dae286f2cb0000000700000000000186a000000000000000000000000100000001645938bb7ae2193270e6ffef009e3664d11e07c10006617373657432000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d6d79566172436170417373657400044d5643410000000001000000000000000200000006000000000000000000000001000000023f58fda2e9ea8d9e4b181832a07b26dae286f2cb645938bb7ae2193270e6ffef009e3664d11e07c100000006000000000000000000000001000000023f78e510df62bc48b0829ec06d6a6b98062d6953c54903de5177a16f7811771ef2f4659d9e864671000661737365743300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000126d794f7468657256617243617041737365740000000000000100000000000000010000000600000000000000000000000100000001645938bb7ae2193270e6ffef009e3664d11e07c1279fa028")
		})

	ginkgo.It("can make calls to platformvm static api", func() {
		keys := []*secp256k1.PrivateKey{}
		for _, key := range []string{
			"24jUJ9vZexUM6expyMcT48LBx27k1m7xpraoV62oSQAHdziao5",
			"2MMvUMsxx6zsHSNXJdFD8yc5XkancvwyKPwpw4xUK3TCGDuNBY",
			"cxb7KpGWhDMALTjNNSJ7UQkkomPesyWAPUaWRGdyeBNzR6f35",
			"ewoqjP7PxY4yr3iLTpLisriqt94hdyDFNgchSxGGztUrTXtNN",
			"2RWLv6YVEXDiWLpaCbXhhqxtLbnFaKQsWPSSMSPhpWo47uJAeV",
		} {
			privKeyBytes, err := cb58.Decode(key)
			require.NoError(err)
			pk, err := secp256k1.ToPrivateKey(privKeyBytes)
			require.NoError(err)
			keys = append(keys, pk)
		}

		genesisUTXOs := make([]api.UTXO, len(keys))
		hrp := constants.NetworkIDToHRP[constants.UnitTestID]
		for i, key := range keys {
			id := key.PublicKey().Address()
			addr, err := address.FormatBech32(hrp, id.Bytes())
			require.NoError(err)
			genesisUTXOs[i] = api.UTXO{
				Amount:  json.Uint64(50000 * units.MilliLux),
				Address: addr,
			}
		}

		genesisValidators := make([]api.GenesisPermissionlessValidator, len(keys))
		for i, key := range keys {
			id := key.PublicKey().Address()
			addr, err := address.FormatBech32(hrp, id.Bytes())
			require.NoError(err)
			genesisValidators[i] = api.GenesisPermissionlessValidator{
				GenesisValidator: api.GenesisValidator{
					StartTime: json.Uint64(time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC).Unix()),
					EndTime:   json.Uint64(time.Date(1997, 1, 30, 0, 0, 0, 0, time.UTC).Unix()),
					NodeID:    ids.BuildTestNodeID(id[:]),
				},
				RewardOwner: &api.Owner{
					Threshold: 1,
					Addresses: []string{addr},
				},
				Staked: []api.UTXO{{
					Amount:  json.Uint64(10000),
					Address: addr,
				}},
				DelegationFee: reward.PercentDenominator,
			}
		}

		buildGenesisArgs := api.BuildGenesisArgs{
			NetworkID:     json.Uint32(constants.UnitTestID),
			LuxAssetID:   ids.ID{'a', 'v', 'a', 'x'},
			UTXOs:         genesisUTXOs,
			Validators:    genesisValidators,
			Chains:        nil,
			Time:          json.Uint64(time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC).Unix()),
			InitialSupply: json.Uint64(360 * units.MegaLux),
			Encoding:      formatting.Hex,
		}

		staticClient := api.NewStaticClient(e2e.Env.GetRandomNodeURI().URI)
		resp, err := staticClient.BuildGenesis(e2e.DefaultContext(), &buildGenesisArgs)
		require.NoError(err)
		require.Equal(resp.Bytes, "0x0000000000050000000000000000000000000000000000000000000000000000000000000000000000006176617800000000000000000000000000000000000000000000000000000000000000070000000ba43b740000000000000000000000000100000001fceda8f90fcb5d30614b99d79fc4baa293077626000000000000000000000000000000000000000000000000000000000000000000000000000000016176617800000000000000000000000000000000000000000000000000000000000000070000000ba43b7400000000000000000000000001000000016ead693c17abb1be422bb50b30b9711ff98d667e000000000000000000000000000000000000000000000000000000000000000000000000000000026176617800000000000000000000000000000000000000000000000000000000000000070000000ba43b740000000000000000000000000100000001f2420846876e69f473dda256172967e992f0ee31000000000000000000000000000000000000000000000000000000000000000000000000000000036176617800000000000000000000000000000000000000000000000000000000000000070000000ba43b7400000000000000000000000001000000013cb7d3842e8cee6a0ebd09f1fe884f6861e1b29c000000000000000000000000000000000000000000000000000000000000000000000000000000046176617800000000000000000000000000000000000000000000000000000000000000070000000ba43b74000000000000000000000000010000000187c4ec0736fdad03fd9ec8c3ba609de958601a7b00000000000000050000000c0000000a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fceda8f90fcb5d30614b99d79fc4baa2930776260000000032c9a9000000000032efe480000000000000271000000001617661780000000000000000000000000000000000000000000000000000000000000007000000000000271000000000000000000000000100000001fceda8f90fcb5d30614b99d79fc4baa2930776260000000b00000000000000000000000100000001fceda8f90fcb5d30614b99d79fc4baa29307762600000000000000000000000c0000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000006ead693c17abb1be422bb50b30b9711ff98d667e0000000032c9a9000000000032efe4800000000000002710000000016176617800000000000000000000000000000000000000000000000000000000000000070000000000002710000000000000000000000001000000016ead693c17abb1be422bb50b30b9711ff98d667e0000000b000000000000000000000001000000016ead693c17abb1be422bb50b30b9711ff98d667e00000000000000000000000c0000000a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f2420846876e69f473dda256172967e992f0ee310000000032c9a9000000000032efe480000000000000271000000001617661780000000000000000000000000000000000000000000000000000000000000007000000000000271000000000000000000000000100000001f2420846876e69f473dda256172967e992f0ee310000000b00000000000000000000000100000001f2420846876e69f473dda256172967e992f0ee3100000000000000000000000c0000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000003cb7d3842e8cee6a0ebd09f1fe884f6861e1b29c0000000032c9a9000000000032efe4800000000000002710000000016176617800000000000000000000000000000000000000000000000000000000000000070000000000002710000000000000000000000001000000013cb7d3842e8cee6a0ebd09f1fe884f6861e1b29c0000000b000000000000000000000001000000013cb7d3842e8cee6a0ebd09f1fe884f6861e1b29c00000000000000000000000c0000000a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000087c4ec0736fdad03fd9ec8c3ba609de958601a7b0000000032c9a9000000000032efe48000000000000027100000000161766178000000000000000000000000000000000000000000000000000000000000000700000000000027100000000000000000000000010000000187c4ec0736fdad03fd9ec8c3ba609de958601a7b0000000b0000000000000000000000010000000187c4ec0736fdad03fd9ec8c3ba609de958601a7b0000000000000000000000000000000032c9a90004fefa17b724000000008e96cbef")
	})
})
