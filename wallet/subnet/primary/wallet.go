// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package primary

import (
	"context"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/keychain"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/wallet/chain/c"
	"github.com/skychains/chain/wallet/chain/p"
	"github.com/skychains/chain/wallet/chain/x"
	"github.com/skychains/chain/wallet/subnet/primary/common"

	pbuilder "github.com/skychains/chain/wallet/chain/p/builder"
	psigner "github.com/skychains/chain/wallet/chain/p/signer"
	xbuilder "github.com/skychains/chain/wallet/chain/x/builder"
	xsigner "github.com/skychains/chain/wallet/chain/x/signer"
)

var _ Wallet = (*wallet)(nil)

// Wallet provides chain wallets for the primary network.
type Wallet interface {
	P() p.Wallet
	X() x.Wallet
	C() c.Wallet
}

type wallet struct {
	p p.Wallet
	x x.Wallet
	c c.Wallet
}

func (w *wallet) P() p.Wallet {
	return w.p
}

func (w *wallet) X() x.Wallet {
	return w.x
}

func (w *wallet) C() c.Wallet {
	return w.c
}

// Creates a new default wallet
func NewWallet(p p.Wallet, x x.Wallet, c c.Wallet) Wallet {
	return &wallet{
		p: p,
		x: x,
		c: c,
	}
}

// Creates a Wallet with the given set of options
func NewWalletWithOptions(w Wallet, options ...common.Option) Wallet {
	return NewWallet(
		p.NewWalletWithOptions(w.P(), options...),
		x.NewWalletWithOptions(w.X(), options...),
		c.NewWalletWithOptions(w.C(), options...),
	)
}

type WalletConfig struct {
	// Base URI to use for all node requests.
	URI string // required
	// Keys to use for signing all transactions.
	LUXKeychain keychain.Keychain // required
	EthKeychain  c.EthKeychain     // required
	// Set of P-chain transactions that the wallet should know about to be able
	// to generate transactions.
	PChainTxs map[ids.ID]*txs.Tx // optional
	// Set of P-chain transactions that the wallet should fetch to be able to
	// generate transactions.
	PChainTxsToFetch set.Set[ids.ID] // optional
}

// MakeWallet returns a wallet that supports issuing transactions to the chains
// living in the primary network.
//
// On creation, the wallet attaches to the provided uri and fetches all UTXOs
// that reference any of the provided keys. If the UTXOs are modified through an
// external issuance process, such as another instance of the wallet, the UTXOs
// may become out of sync. The wallet will also fetch all requested P-chain
// transactions.
//
// The wallet manages all state locally, and performs all tx signing locally.
func MakeWallet(ctx context.Context, config *WalletConfig) (Wallet, error) {
	luxAddrs := config.LUXKeychain.Addresses()
	luxState, err := FetchState(ctx, config.URI, luxAddrs)
	if err != nil {
		return nil, err
	}

	ethAddrs := config.EthKeychain.EthAddresses()
	ethState, err := FetchEthState(ctx, config.URI, ethAddrs)
	if err != nil {
		return nil, err
	}

	pChainTxs := config.PChainTxs
	if pChainTxs == nil {
		pChainTxs = make(map[ids.ID]*txs.Tx)
	}

	for txID := range config.PChainTxsToFetch {
		txBytes, err := luxState.PClient.GetTx(ctx, txID)
		if err != nil {
			return nil, err
		}
		tx, err := txs.Parse(txs.Codec, txBytes)
		if err != nil {
			return nil, err
		}
		pChainTxs[txID] = tx
	}

	pUTXOs := common.NewChainUTXOs(constants.PlatformChainID, luxState.UTXOs)
	pBackend := p.NewBackend(luxState.PCTX, pUTXOs, pChainTxs)
	pBuilder := pbuilder.New(luxAddrs, luxState.PCTX, pBackend)
	pSigner := psigner.New(config.LUXKeychain, pBackend)

	xChainID := luxState.XCTX.BlockchainID
	xUTXOs := common.NewChainUTXOs(xChainID, luxState.UTXOs)
	xBackend := x.NewBackend(luxState.XCTX, xUTXOs)
	xBuilder := xbuilder.New(luxAddrs, luxState.XCTX, xBackend)
	xSigner := xsigner.New(config.LUXKeychain, xBackend)

	cChainID := luxState.CCTX.BlockchainID
	cUTXOs := common.NewChainUTXOs(cChainID, luxState.UTXOs)
	cBackend := c.NewBackend(cUTXOs, ethState.Accounts)
	cBuilder := c.NewBuilder(luxAddrs, ethAddrs, luxState.CCTX, cBackend)
	cSigner := c.NewSigner(config.LUXKeychain, config.EthKeychain, cBackend)

	return NewWallet(
		p.NewWallet(pBuilder, pSigner, luxState.PClient, pBackend),
		x.NewWallet(xBuilder, xSigner, luxState.XClient, xBackend),
		c.NewWallet(cBuilder, cSigner, luxState.CClient, ethState.Client, cBackend),
	), nil
}
