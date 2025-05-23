// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"time"

	"github.com/skychains/coreth/plugin/evm"

	"github.com/skychains/chain/genesis"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/components/lux"
	"github.com/skychains/chain/vms/secp256k1fx"
	"github.com/skychains/chain/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	luxAddr := key.Address()
	ethAddr := evm.PublicKeyToEthAddress(key.PublicKey())

	ctx := context.Background()

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:         uri,
		LUXKeychain: kc,
		EthKeychain: kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()
	cWallet := wallet.C()

	// Pull out useful constants to use when issuing transactions.
	cContext := cWallet.Builder().Context()
	cChainID := cContext.BlockchainID
	luxAssetID := cContext.LUXAssetID
	owner := secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			luxAddr,
		},
	}

	exportStartTime := time.Now()
	exportTx, err := pWallet.IssueExportTx(cChainID, []*lux.TransferableOutput{{
		Asset: lux.Asset{ID: luxAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt:          units.Lux,
			OutputOwners: owner,
		},
	}})
	if err != nil {
		log.Fatalf("failed to issue export transaction: %s\n", err)
	}
	log.Printf("issued export %s in %s\n", exportTx.ID(), time.Since(exportStartTime))

	importStartTime := time.Now()
	importTx, err := cWallet.IssueImportTx(constants.PlatformChainID, ethAddr)
	if err != nil {
		log.Fatalf("failed to issue import transaction: %s\n", err)
	}
	log.Printf("issued import %s to %s in %s\n", importTx.ID(), ethAddr.Hex(), time.Since(importStartTime))
}
