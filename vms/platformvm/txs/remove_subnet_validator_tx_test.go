// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/components/verify"
	"github.com/skychains/chain/vms/platformvm/stakeable"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/vms/types"
)

var errInvalidSubnetAuth = errors.New("invalid subnet auth")

func TestRemoveSubnetValidatorTxSerialization(t *testing.T) {
	require := require.New(t)

	addr := ids.ShortID{
		0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
		0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
		0x44, 0x55, 0x66, 0x77,
	}

	luxAssetID, err := ids.FromString("FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z")
	require.NoError(err)

	customAssetID := ids.ID{
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
	}

	txID := ids.ID{
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
	}
	nodeID := ids.BuildTestNodeID([]byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x11, 0x22, 0x33, 0x44,
	})
	subnetID := ids.ID{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
	}

	simpleRemoveValidatorTx := &RemoveSubnetValidatorTx{
		BaseTx: BaseTx{
			BaseTx: lux.BaseTx{
				NetworkID:    constants.MainnetID,
				BlockchainID: constants.PlatformChainID,
				Outs:         []*lux.TransferableOutput{},
				Ins: []*lux.TransferableInput{
					{
						UTXOID: lux.UTXOID{
							TxID:        txID,
							OutputIndex: 1,
						},
						Asset: lux.Asset{
							ID: luxAssetID,
						},
						In: &secp256k1fx.TransferInput{
							Amt: units.MilliLux,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{5},
							},
						},
					},
				},
				Memo: types.JSONByteSlice{},
			},
		},
		NodeID: nodeID,
		Subnet: subnetID,
		SubnetAuth: &secp256k1fx.Input{
			SigIndices: []uint32{3},
		},
	}
	require.NoError(simpleRemoveValidatorTx.SyntacticVerify(&snow.Context{
		NetworkID:   1,
		ChainID:     constants.PlatformChainID,
		LUXAssetID: luxAssetID,
	}))

	expectedUnsignedSimpleRemoveValidatorTxBytes := []byte{
		// Codec version
		0x00, 0x00,
		// RemoveSubnetValidatorTx Type ID
		0x00, 0x00, 0x00, 0x17,
		// Mainnet network ID
		0x00, 0x00, 0x00, 0x01,
		// P-chain blockchain ID
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// Number of outputs
		0x00, 0x00, 0x00, 0x00,
		// Number of inputs
		0x00, 0x00, 0x00, 0x01,
		// Inputs[0]
		// TxID
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		// Tx output index
		0x00, 0x00, 0x00, 0x01,
		// Mainnet LUX assetID
		0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
		0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
		0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
		0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
		// secp256k1fx transfer input type ID
		0x00, 0x00, 0x00, 0x05,
		// input amount = 1 MilliLux
		0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x42, 0x40,
		// number of signatures needed in input
		0x00, 0x00, 0x00, 0x01,
		// index of signer
		0x00, 0x00, 0x00, 0x05,
		// length of memo field
		0x00, 0x00, 0x00, 0x00,
		// nodeID to remove
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x11, 0x22, 0x33, 0x44,
		// subnetID to remove from
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		// secp256k1fx authorization type ID
		0x00, 0x00, 0x00, 0x0a,
		// number of signatures needed in authorization
		0x00, 0x00, 0x00, 0x01,
		// index of signer
		0x00, 0x00, 0x00, 0x03,
	}
	var unsignedSimpleRemoveValidatorTx UnsignedTx = simpleRemoveValidatorTx
	unsignedSimpleRemoveValidatorTxBytes, err := Codec.Marshal(CodecVersion, &unsignedSimpleRemoveValidatorTx)
	require.NoError(err)
	require.Equal(expectedUnsignedSimpleRemoveValidatorTxBytes, unsignedSimpleRemoveValidatorTxBytes)

	complexRemoveValidatorTx := &RemoveSubnetValidatorTx{
		BaseTx: BaseTx{
			BaseTx: lux.BaseTx{
				NetworkID:    constants.MainnetID,
				BlockchainID: constants.PlatformChainID,
				Outs: []*lux.TransferableOutput{
					{
						Asset: lux.Asset{
							ID: luxAssetID,
						},
						Out: &stakeable.LockOut{
							Locktime: 87654321,
							TransferableOut: &secp256k1fx.TransferOutput{
								Amt: 1,
								OutputOwners: secp256k1fx.OutputOwners{
									Locktime:  12345678,
									Threshold: 0,
									Addrs:     []ids.ShortID{},
								},
							},
						},
					},
					{
						Asset: lux.Asset{
							ID: customAssetID,
						},
						Out: &stakeable.LockOut{
							Locktime: 876543210,
							TransferableOut: &secp256k1fx.TransferOutput{
								Amt: 0xffffffffffffffff,
								OutputOwners: secp256k1fx.OutputOwners{
									Locktime:  0,
									Threshold: 1,
									Addrs: []ids.ShortID{
										addr,
									},
								},
							},
						},
					},
				},
				Ins: []*lux.TransferableInput{
					{
						UTXOID: lux.UTXOID{
							TxID:        txID,
							OutputIndex: 1,
						},
						Asset: lux.Asset{
							ID: luxAssetID,
						},
						In: &secp256k1fx.TransferInput{
							Amt: units.Lux,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{2, 5},
							},
						},
					},
					{
						UTXOID: lux.UTXOID{
							TxID:        txID,
							OutputIndex: 2,
						},
						Asset: lux.Asset{
							ID: customAssetID,
						},
						In: &stakeable.LockIn{
							Locktime: 876543210,
							TransferableIn: &secp256k1fx.TransferInput{
								Amt: 0xefffffffffffffff,
								Input: secp256k1fx.Input{
									SigIndices: []uint32{0},
								},
							},
						},
					},
					{
						UTXOID: lux.UTXOID{
							TxID:        txID,
							OutputIndex: 3,
						},
						Asset: lux.Asset{
							ID: customAssetID,
						},
						In: &secp256k1fx.TransferInput{
							Amt: 0x1000000000000000,
							Input: secp256k1fx.Input{
								SigIndices: []uint32{},
							},
						},
					},
				},
				Memo: types.JSONByteSlice("😅\nwell that's\x01\x23\x45!"),
			},
		},
		NodeID: nodeID,
		Subnet: subnetID,
		SubnetAuth: &secp256k1fx.Input{
			SigIndices: []uint32{},
		},
	}
	lux.SortTransferableOutputs(complexRemoveValidatorTx.Outs, Codec)
	utils.Sort(complexRemoveValidatorTx.Ins)
	require.NoError(complexRemoveValidatorTx.SyntacticVerify(&snow.Context{
		NetworkID:   1,
		ChainID:     constants.PlatformChainID,
		LUXAssetID: luxAssetID,
	}))

	expectedUnsignedComplexRemoveValidatorTxBytes := []byte{
		// Codec version
		0x00, 0x00,
		// RemoveSubnetValidatorTx Type ID
		0x00, 0x00, 0x00, 0x17,
		// Mainnet network ID
		0x00, 0x00, 0x00, 0x01,
		// P-chain blockchain ID
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// Number of outputs
		0x00, 0x00, 0x00, 0x02,
		// Outputs[0]
		// Mainnet LUX assetID
		0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
		0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
		0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
		0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
		// Stakeable locked output type ID
		0x00, 0x00, 0x00, 0x16,
		// Locktime
		0x00, 0x00, 0x00, 0x00, 0x05, 0x39, 0x7f, 0xb1,
		// secp256k1fx transfer output type ID
		0x00, 0x00, 0x00, 0x07,
		// amount
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		// secp256k1fx output locktime
		0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e,
		// threshold
		0x00, 0x00, 0x00, 0x00,
		// number of addresses
		0x00, 0x00, 0x00, 0x00,
		// Outputs[1]
		// custom asset ID
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		// Stakeable locked output type ID
		0x00, 0x00, 0x00, 0x16,
		// Locktime
		0x00, 0x00, 0x00, 0x00, 0x34, 0x3e, 0xfc, 0xea,
		// secp256k1fx transfer output type ID
		0x00, 0x00, 0x00, 0x07,
		// amount
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		// secp256k1fx output locktime
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// threshold
		0x00, 0x00, 0x00, 0x01,
		// number of addresses
		0x00, 0x00, 0x00, 0x01,
		// address[0]
		0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
		0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
		0x44, 0x55, 0x66, 0x77,
		// number of inputs
		0x00, 0x00, 0x00, 0x03,
		// inputs[0]
		// TxID
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		// Tx output index
		0x00, 0x00, 0x00, 0x01,
		// Mainnet LUX assetID
		0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
		0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
		0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
		0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
		// secp256k1fx transfer input type ID
		0x00, 0x00, 0x00, 0x05,
		// input amount = 1 Lux
		0x00, 0x00, 0x00, 0x00, 0x3b, 0x9a, 0xca, 0x00,
		// number of signatures needed in input
		0x00, 0x00, 0x00, 0x02,
		// index of first signer
		0x00, 0x00, 0x00, 0x02,
		// index of second signer
		0x00, 0x00, 0x00, 0x05,
		// inputs[1]
		// TxID
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		// Tx output index
		0x00, 0x00, 0x00, 0x02,
		// Custom asset ID
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		// Stakeable locked input type ID
		0x00, 0x00, 0x00, 0x15,
		// Locktime
		0x00, 0x00, 0x00, 0x00, 0x34, 0x3e, 0xfc, 0xea,
		// secp256k1fx transfer input type ID
		0x00, 0x00, 0x00, 0x05,
		// input amount
		0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		// number of signatures needed in input
		0x00, 0x00, 0x00, 0x01,
		// index of signer
		0x00, 0x00, 0x00, 0x00,
		// inputs[2]
		// TxID
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		// Tx output index
		0x00, 0x00, 0x00, 0x03,
		// custom asset ID
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		// secp256k1fx transfer input type ID
		0x00, 0x00, 0x00, 0x05,
		// input amount
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// number of signatures needed in input
		0x00, 0x00, 0x00, 0x00,
		// length of memo
		0x00, 0x00, 0x00, 0x14,
		// memo
		0xf0, 0x9f, 0x98, 0x85, 0x0a, 0x77, 0x65, 0x6c,
		0x6c, 0x20, 0x74, 0x68, 0x61, 0x74, 0x27, 0x73,
		0x01, 0x23, 0x45, 0x21,
		// nodeID to remove
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x11, 0x22, 0x33, 0x44,
		// subnetID to remove from
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		// secp256k1fx authorization type ID
		0x00, 0x00, 0x00, 0x0a,
		// number of signatures needed in authorization
		0x00, 0x00, 0x00, 0x00,
	}
	var unsignedComplexRemoveValidatorTx UnsignedTx = complexRemoveValidatorTx
	unsignedComplexRemoveValidatorTxBytes, err := Codec.Marshal(CodecVersion, &unsignedComplexRemoveValidatorTx)
	require.NoError(err)
	require.Equal(expectedUnsignedComplexRemoveValidatorTxBytes, unsignedComplexRemoveValidatorTxBytes)

	aliaser := ids.NewAliaser()
	require.NoError(aliaser.Alias(constants.PlatformChainID, "P"))

	unsignedComplexRemoveValidatorTx.InitCtx(&snow.Context{
		NetworkID:   1,
		ChainID:     constants.PlatformChainID,
		LUXAssetID: luxAssetID,
		BCLookup:    aliaser,
	})

	unsignedComplexRemoveValidatorTxJSONBytes, err := json.MarshalIndent(unsignedComplexRemoveValidatorTx, "", "\t")
	require.NoError(err)
	require.Equal(`{
	"networkID": 1,
	"blockchainID": "11111111111111111111111111111111LpoYY",
	"outputs": [
		{
			"assetID": "FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z",
			"fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
			"output": {
				"locktime": 87654321,
				"output": {
					"addresses": [],
					"amount": 1,
					"locktime": 12345678,
					"threshold": 0
				}
			}
		},
		{
			"assetID": "2Ab62uWwJw1T6VvmKD36ufsiuGZuX1pGykXAvPX1LtjTRHxwcc",
			"fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
			"output": {
				"locktime": 876543210,
				"output": {
					"addresses": [
						"P-lux1g32kvaugnx4tk3z4vemc3xd2hdz92enh972wxr"
					],
					"amount": 18446744073709551615,
					"locktime": 0,
					"threshold": 1
				}
			}
		}
	],
	"inputs": [
		{
			"txID": "2wiU5PnFTjTmoAXGZutHAsPF36qGGyLHYHj9G1Aucfmb3JFFGN",
			"outputIndex": 1,
			"assetID": "FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z",
			"fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
			"input": {
				"amount": 1000000000,
				"signatureIndices": [
					2,
					5
				]
			}
		},
		{
			"txID": "2wiU5PnFTjTmoAXGZutHAsPF36qGGyLHYHj9G1Aucfmb3JFFGN",
			"outputIndex": 2,
			"assetID": "2Ab62uWwJw1T6VvmKD36ufsiuGZuX1pGykXAvPX1LtjTRHxwcc",
			"fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
			"input": {
				"locktime": 876543210,
				"input": {
					"amount": 17293822569102704639,
					"signatureIndices": [
						0
					]
				}
			}
		},
		{
			"txID": "2wiU5PnFTjTmoAXGZutHAsPF36qGGyLHYHj9G1Aucfmb3JFFGN",
			"outputIndex": 3,
			"assetID": "2Ab62uWwJw1T6VvmKD36ufsiuGZuX1pGykXAvPX1LtjTRHxwcc",
			"fxID": "spdxUxVJQbX85MGxMHbKw1sHxMnSqJ3QBzDyDYEP3h6TLuxqQ",
			"input": {
				"amount": 1152921504606846976,
				"signatureIndices": []
			}
		}
	],
	"memo": "0xf09f98850a77656c6c2074686174277301234521",
	"nodeID": "NodeID-2ZbTY9GatRTrfinAoYiYLcf6CvrPAUYgo",
	"subnetID": "SkB92YpWm4UpburLz9tEKZw2i67H3FF6YkjaU4BkFUDTG9Xm",
	"subnetAuthorization": {
		"signatureIndices": []
	}
}`, string(unsignedComplexRemoveValidatorTxJSONBytes))
}

func TestRemoveSubnetValidatorTxSyntacticVerify(t *testing.T) {
	type test struct {
		name        string
		txFunc      func(*gomock.Controller) *RemoveSubnetValidatorTx
		expectedErr error
	}

	var (
		networkID = uint32(1337)
		chainID   = ids.GenerateTestID()
	)

	ctx := &snow.Context{
		ChainID:   chainID,
		NetworkID: networkID,
	}

	// A BaseTx that already passed syntactic verification.
	verifiedBaseTx := BaseTx{
		SyntacticallyVerified: true,
	}
	// Sanity check.
	require.NoError(t, verifiedBaseTx.SyntacticVerify(ctx))

	// A BaseTx that passes syntactic verification.
	validBaseTx := BaseTx{
		BaseTx: lux.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
		},
	}
	// Sanity check.
	require.NoError(t, validBaseTx.SyntacticVerify(ctx))
	// Make sure we're not caching the verification result.
	require.False(t, validBaseTx.SyntacticallyVerified)

	// A BaseTx that fails syntactic verification.
	invalidBaseTx := BaseTx{}

	tests := []test{
		{
			name: "nil tx",
			txFunc: func(*gomock.Controller) *RemoveSubnetValidatorTx {
				return nil
			},
			expectedErr: ErrNilTx,
		},
		{
			name: "already verified",
			txFunc: func(*gomock.Controller) *RemoveSubnetValidatorTx {
				return &RemoveSubnetValidatorTx{BaseTx: verifiedBaseTx}
			},
			expectedErr: nil,
		},
		{
			name: "invalid BaseTx",
			txFunc: func(*gomock.Controller) *RemoveSubnetValidatorTx {
				return &RemoveSubnetValidatorTx{
					// Set subnetID so we don't error on that check.
					Subnet: ids.GenerateTestID(),
					// Set NodeID so we don't error on that check.
					NodeID: ids.GenerateTestNodeID(),
					BaseTx: invalidBaseTx,
				}
			},
			expectedErr: lux.ErrWrongNetworkID,
		},
		{
			name: "invalid subnetID",
			txFunc: func(*gomock.Controller) *RemoveSubnetValidatorTx {
				return &RemoveSubnetValidatorTx{
					BaseTx: validBaseTx,
					// Set NodeID so we don't error on that check.
					NodeID: ids.GenerateTestNodeID(),
					Subnet: constants.PrimaryNetworkID,
				}
			},
			expectedErr: ErrRemovePrimaryNetworkValidator,
		},
		{
			name: "invalid subnetAuth",
			txFunc: func(ctrl *gomock.Controller) *RemoveSubnetValidatorTx {
				// This SubnetAuth fails verification.
				invalidSubnetAuth := verify.NewMockVerifiable(ctrl)
				invalidSubnetAuth.EXPECT().Verify().Return(errInvalidSubnetAuth)
				return &RemoveSubnetValidatorTx{
					// Set subnetID so we don't error on that check.
					Subnet: ids.GenerateTestID(),
					// Set NodeID so we don't error on that check.
					NodeID:     ids.GenerateTestNodeID(),
					BaseTx:     validBaseTx,
					SubnetAuth: invalidSubnetAuth,
				}
			},
			expectedErr: errInvalidSubnetAuth,
		},
		{
			name: "passes verification",
			txFunc: func(ctrl *gomock.Controller) *RemoveSubnetValidatorTx {
				// This SubnetAuth passes verification.
				validSubnetAuth := verify.NewMockVerifiable(ctrl)
				validSubnetAuth.EXPECT().Verify().Return(nil)
				return &RemoveSubnetValidatorTx{
					// Set subnetID so we don't error on that check.
					Subnet: ids.GenerateTestID(),
					// Set NodeID so we don't error on that check.
					NodeID:     ids.GenerateTestNodeID(),
					BaseTx:     validBaseTx,
					SubnetAuth: validSubnetAuth,
				}
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			ctrl := gomock.NewController(t)

			tx := tt.txFunc(ctrl)
			err := tx.SyntacticVerify(ctx)
			require.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}
			require.True(tx.SyntacticallyVerified)
		})
	}
}
