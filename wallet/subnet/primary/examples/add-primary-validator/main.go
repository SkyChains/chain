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
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/platformvm/reward"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	startTime := time.Now().Add(time.Minute)
	duration := 3 * 7 * 24 * time.Hour // 3 weeks
	weight := 2_000 * units.Lux
	validatorRewardAddr := key.Address()
	delegatorRewardAddr := key.Address()
	delegationFee := uint32(reward.PercentDenominator / 2) // 50%

	ctx := context.Background()
	infoClient := info.NewClient(uri)

	nodeInfoStartTime := time.Now()
	nodeID, nodePOP, err := infoClient.GetNodeID(ctx)
	if err != nil {
		log.Fatalf("failed to fetch node IDs: %s\n", err)
	}
	log.Printf("fetched node ID %s in %s\n", nodeID, time.Since(nodeInfoStartTime))

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          uri,
		LUXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()
	pBuilder := pWallet.Builder()
	pContext := pBuilder.Context()
	luxAssetID := pContext.LUXAssetID

	addValidatorStartTime := time.Now()
	addValidatorTx, err := pWallet.IssueAddPermissionlessValidatorTx(
		&txs.SubnetValidator{Validator: txs.Validator{
			NodeID: nodeID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(startTime.Add(duration).Unix()),
			Wght:   weight,
		}},
		nodePOP,
		luxAssetID,
		&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{validatorRewardAddr},
		},
		&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{delegatorRewardAddr},
		},
		delegationFee,
	)
	if err != nil {
		log.Fatalf("failed to issue add permissionless validator transaction: %s\n", err)
	}
	log.Printf("added new primary network validator %s with %s in %s\n", nodeID, addValidatorTx.ID(), time.Since(addValidatorStartTime))
}
