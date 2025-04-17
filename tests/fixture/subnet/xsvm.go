// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package subnet

import (
	"math"
	"time"

	"github.com/SkyChains/chain/tests/fixture/tmpnet"
	"github.com/SkyChains/chain/utils/constants"
	"github.com/SkyChains/chain/utils/crypto/secp256k1"
	"github.com/SkyChains/chain/vms/example/xsvm/genesis"
)

func NewXSVMOrPanic(name string, key *secp256k1.PrivateKey, nodes ...*tmpnet.Node) *tmpnet.Subnet {
	if len(nodes) == 0 {
		panic("a subnet must be validated by at least one node")
	}

	genesisBytes, err := genesis.Codec.Marshal(genesis.CodecVersion, &genesis.Genesis{
		Timestamp: time.Now().Unix(),
		Allocations: []genesis.Allocation{
			{
				Address: key.Address(),
				Balance: math.MaxUint64,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	return &tmpnet.Subnet{
		Name: name,
		Chains: []*tmpnet.Chain{
			{
				VMID:         constants.XSVMID,
				Genesis:      genesisBytes,
				PreFundedKey: key,
			},
		},
		ValidatorIDs: tmpnet.NodesToIDs(nodes...),
	}
}
