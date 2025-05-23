// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/skychains/chain/genesis"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/wallet/subnet/primary"

	xsgenesis "github.com/skychains/chain/vms/example/xsvm/genesis"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	subnetIDStr := "29uVeLPJB1eQJkzRemU8g8wZDw5uJRqpab5U2mX9euieVwiEbL"
	genesis := &xsgenesis.Genesis{
		Timestamp: time.Now().Unix(),
		Allocations: []xsgenesis.Allocation{
			{
				Address: genesis.EWOQKey.Address(),
				Balance: math.MaxUint64,
			},
		},
	}
	vmID := constants.XSVMID
	name := "let there"

	subnetID, err := ids.FromString(subnetIDStr)
	if err != nil {
		log.Fatalf("failed to parse subnet ID: %s\n", err)
	}

	genesisBytes, err := xsgenesis.Codec.Marshal(xsgenesis.CodecVersion, genesis)
	if err != nil {
		log.Fatalf("failed to create genesis bytes: %s\n", err)
	}

	ctx := context.Background()

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting and registers [subnetID].
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:              uri,
		LUXKeychain:     kc,
		EthKeychain:      kc,
		PChainTxsToFetch: set.Of(subnetID),
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()

	createChainStartTime := time.Now()
	createChainTx, err := pWallet.IssueCreateChainTx(
		subnetID,
		genesisBytes,
		vmID,
		nil,
		name,
	)
	if err != nil {
		log.Fatalf("failed to issue create chain transaction: %s\n", err)
	}
	log.Printf("created new chain %s in %s\n", createChainTx.ID(), time.Since(createChainStartTime))
}
