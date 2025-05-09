// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/skychains/chain/api/info"
	"github.com/skychains/chain/genesis"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	subnetIDStr := "29uVeLPJB1eQJkzRemU8g8wZDw5uJRqpab5U2mX9euieVwiEbL"
	startTime := time.Now().Add(time.Minute)
	duration := 2 * 7 * 24 * time.Hour // 2 weeks
	weight := units.Schmeckle

	subnetID, err := ids.FromString(subnetIDStr)
	if err != nil {
		log.Fatalf("failed to parse subnet ID: %s\n", err)
	}

	ctx := context.Background()
	infoClient := info.NewClient(uri)

	nodeInfoStartTime := time.Now()
	nodeID, _, err := infoClient.GetNodeID(ctx)
	if err != nil {
		log.Fatalf("failed to fetch node IDs: %s\n", err)
	}
	log.Printf("fetched node ID %s in %s\n", nodeID, time.Since(nodeInfoStartTime))

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

	addValidatorStartTime := time.Now()
	addValidatorTx, err := pWallet.IssueAddSubnetValidatorTx(&txs.SubnetValidator{
		Validator: txs.Validator{
			NodeID: nodeID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(startTime.Add(duration).Unix()),
			Wght:   weight,
		},
		Subnet: subnetID,
	})
	if err != nil {
		log.Fatalf("failed to issue add subnet validator transaction: %s\n", err)
	}
	log.Printf("added new subnet validator %s to %s with %s in %s\n", nodeID, subnetID, addValidatorTx.ID(), time.Since(addValidatorStartTime))
}
