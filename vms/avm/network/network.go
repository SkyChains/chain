// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/network/p2p"
	"github.com/skychains/chain/network/p2p/gossip"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/validators"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/vms/avm/txs"
	"github.com/skychains/chain/vms/avm/txs/mempool"
)

const txGossipHandlerID = 0

var (
	_ common.AppHandler    = (*Network)(nil)
	_ validators.Connector = (*Network)(nil)
)

type Network struct {
	*p2p.Network

	log       logging.Logger
	parser    txs.Parser
	mempool   *gossipMempool
	appSender common.AppSender

	txPushGossiper        *gossip.PushGossiper[*txs.Tx]
	txPushGossipFrequency time.Duration
	txPullGossiper        gossip.Gossiper
	txPullGossipFrequency time.Duration
}

func New(
	log logging.Logger,
	nodeID ids.NodeID,
	subnetID ids.ID,
	vdrs validators.State,
	parser txs.Parser,
	txVerifier TxVerifier,
	mempool mempool.Mempool,
	appSender common.AppSender,
	registerer prometheus.Registerer,
	config Config,
) (*Network, error) {
	p2pNetwork, err := p2p.NewNetwork(log, appSender, registerer, "p2p")
	if err != nil {
		return nil, err
	}

	marshaller := &txParser{
		parser: parser,
	}
	validators := p2p.NewValidators(
		p2pNetwork.Peers,
		log,
		subnetID,
		vdrs,
		config.MaxValidatorSetStaleness,
	)
	txGossipClient := p2pNetwork.NewClient(
		txGossipHandlerID,
		p2p.WithValidatorSampling(validators),
	)
	txGossipMetrics, err := gossip.NewMetrics(registerer, "tx")
	if err != nil {
		return nil, err
	}

	gossipMempool, err := newGossipMempool(
		mempool,
		registerer,
		log,
		txVerifier,
		parser,
		config.ExpectedBloomFilterElements,
		config.ExpectedBloomFilterFalsePositiveProbability,
		config.MaxBloomFilterFalsePositiveProbability,
	)
	if err != nil {
		return nil, err
	}

	txPushGossiper, err := gossip.NewPushGossiper[*txs.Tx](
		marshaller,
		gossipMempool,
		validators,
		txGossipClient,
		txGossipMetrics,
		gossip.BranchingFactor{
			StakePercentage: config.PushGossipPercentStake,
			Validators:      config.PushGossipNumValidators,
			Peers:           config.PushGossipNumPeers,
		},
		gossip.BranchingFactor{
			Validators: config.PushRegossipNumValidators,
			Peers:      config.PushRegossipNumPeers,
		},
		config.PushGossipDiscardedCacheSize,
		config.TargetGossipSize,
		config.PushGossipMaxRegossipFrequency,
	)
	if err != nil {
		return nil, err
	}

	var txPullGossiper gossip.Gossiper = gossip.NewPullGossiper[*txs.Tx](
		log,
		marshaller,
		gossipMempool,
		txGossipClient,
		txGossipMetrics,
		config.PullGossipPollSize,
	)

	// Gossip requests are only served if a node is a validator
	txPullGossiper = gossip.ValidatorGossiper{
		Gossiper:   txPullGossiper,
		NodeID:     nodeID,
		Validators: validators,
	}

	handler := gossip.NewHandler[*txs.Tx](
		log,
		marshaller,
		gossipMempool,
		txGossipMetrics,
		config.TargetGossipSize,
	)

	validatorHandler := p2p.NewValidatorHandler(
		p2p.NewThrottlerHandler(
			handler,
			p2p.NewSlidingWindowThrottler(
				config.PullGossipThrottlingPeriod,
				config.PullGossipThrottlingLimit,
			),
			log,
		),
		validators,
		log,
	)

	// We allow pushing txs between all peers, but only serve gossip requests
	// from validators
	txGossipHandler := txGossipHandler{
		appGossipHandler:  handler,
		appRequestHandler: validatorHandler,
	}

	if err := p2pNetwork.AddHandler(txGossipHandlerID, txGossipHandler); err != nil {
		return nil, err
	}

	return &Network{
		Network:               p2pNetwork,
		log:                   log,
		parser:                parser,
		mempool:               gossipMempool,
		appSender:             appSender,
		txPushGossiper:        txPushGossiper,
		txPushGossipFrequency: config.PushGossipFrequency,
		txPullGossiper:        txPullGossiper,
		txPullGossipFrequency: config.PullGossipFrequency,
	}, nil
}

func (n *Network) PushGossip(ctx context.Context) {
	gossip.Every(ctx, n.log, n.txPushGossiper, n.txPushGossipFrequency)
}

func (n *Network) PullGossip(ctx context.Context) {
	gossip.Every(ctx, n.log, n.txPullGossiper, n.txPullGossipFrequency)
}

// IssueTxFromRPC attempts to add a tx to the mempool, after verifying it. If
// the tx is added to the mempool, it will attempt to push gossip the tx to
// random peers in the network.
//
// If the tx is already in the mempool, mempool.ErrDuplicateTx will be
// returned.
// If the tx is not added to the mempool, an error will be returned.
func (n *Network) IssueTxFromRPC(tx *txs.Tx) error {
	if err := n.mempool.Add(tx); err != nil {
		return err
	}
	n.txPushGossiper.Add(tx)
	return nil
}

// IssueTxFromRPCWithoutVerification attempts to add a tx to the mempool,
// without first verifying it. If the tx is added to the mempool, it will
// attempt to push gossip the tx to random peers in the network.
//
// If the tx is already in the mempool, mempool.ErrDuplicateTx will be
// returned.
// If the tx is not added to the mempool, an error will be returned.
func (n *Network) IssueTxFromRPCWithoutVerification(tx *txs.Tx) error {
	if err := n.mempool.AddWithoutVerification(tx); err != nil {
		return err
	}
	n.txPushGossiper.Add(tx)
	return nil
}
