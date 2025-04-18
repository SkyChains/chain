// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"go.uber.org/zap"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/crypto/bls"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/vms/types"
)

var _ SetCallbackListener = (*logger)(nil)

type logger struct {
	log      logging.Logger
	subnetID ids.ID
	nodeIDs  set.Set[ids.NodeID]
}

// NewLogger returns a callback listener that will log validator set changes for
// the specified validators
func NewLogger(
	log logging.Logger,
	subnetID ids.ID,
	nodeIDs ...ids.NodeID,
) SetCallbackListener {
	nodeIDSet := set.Of(nodeIDs...)
	return &logger{
		log:      log,
		subnetID: subnetID,
		nodeIDs:  nodeIDSet,
	}
}

func (l *logger) OnValidatorAdded(
	nodeID ids.NodeID,
	pk *bls.PublicKey,
	txID ids.ID,
	weight uint64,
) {
	if l.nodeIDs.Contains(nodeID) {
		var pkBytes []byte
		if pk != nil {
			pkBytes = bls.PublicKeyToCompressedBytes(pk)
		}
		l.log.Info("node added to validator set",
			zap.Stringer("subnetID", l.subnetID),
			zap.Stringer("nodeID", nodeID),
			zap.Reflect("publicKey", types.JSONByteSlice(pkBytes)),
			zap.Stringer("txID", txID),
			zap.Uint64("weight", weight),
		)
	}
}

func (l *logger) OnValidatorRemoved(
	nodeID ids.NodeID,
	weight uint64,
) {
	if l.nodeIDs.Contains(nodeID) {
		l.log.Info("node removed from validator set",
			zap.Stringer("subnetID", l.subnetID),
			zap.Stringer("nodeID", nodeID),
			zap.Uint64("weight", weight),
		)
	}
}

func (l *logger) OnValidatorWeightChanged(
	nodeID ids.NodeID,
	oldWeight uint64,
	newWeight uint64,
) {
	if l.nodeIDs.Contains(nodeID) {
		l.log.Info("validator weight changed",
			zap.Stringer("subnetID", l.subnetID),
			zap.Stringer("nodeID", nodeID),
			zap.Uint64("previousWeight ", oldWeight),
			zap.Uint64("newWeight ", newWeight),
		)
	}
}
