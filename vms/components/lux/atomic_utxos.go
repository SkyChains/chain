// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package lux

import (
	"fmt"

	"github.com/skychains/chain/chains/atomic"
	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/set"
)

var _ AtomicUTXOManager = (*atomicUTXOManager)(nil)

type AtomicUTXOManager interface {
	// GetAtomicUTXOs returns exported UTXOs such that at least one of the
	// addresses in [addrs] is referenced.
	//
	// Returns at most [limit] UTXOs.
	//
	// Returns:
	// * The fetched UTXOs
	// * The address associated with the last UTXO fetched
	// * The ID of the last UTXO fetched
	// * Any error that may have occurred upstream.
	GetAtomicUTXOs(
		chainID ids.ID,
		addrs set.Set[ids.ShortID],
		startAddr ids.ShortID,
		startUTXOID ids.ID,
		limit int,
	) ([]*UTXO, ids.ShortID, ids.ID, error)
}

type atomicUTXOManager struct {
	sm    atomic.SharedMemory
	codec codec.Manager
}

func NewAtomicUTXOManager(sm atomic.SharedMemory, codec codec.Manager) AtomicUTXOManager {
	return &atomicUTXOManager{
		sm:    sm,
		codec: codec,
	}
}

func (a *atomicUTXOManager) GetAtomicUTXOs(
	chainID ids.ID,
	addrs set.Set[ids.ShortID],
	startAddr ids.ShortID,
	startUTXOID ids.ID,
	limit int,
) ([]*UTXO, ids.ShortID, ids.ID, error) {
	addrsList := make([][]byte, addrs.Len())
	i := 0
	for addr := range addrs {
		copied := addr
		addrsList[i] = copied[:]
		i++
	}

	allUTXOBytes, lastAddr, lastUTXO, err := a.sm.Indexed(
		chainID,
		addrsList,
		startAddr.Bytes(),
		startUTXOID[:],
		limit,
	)
	if err != nil {
		return nil, ids.ShortID{}, ids.ID{}, fmt.Errorf("error fetching atomic UTXOs: %w", err)
	}

	lastAddrID, err := ids.ToShortID(lastAddr)
	if err != nil {
		lastAddrID = ids.ShortEmpty
	}
	lastUTXOID, err := ids.ToID(lastUTXO)
	if err != nil {
		lastUTXOID = ids.Empty
	}

	utxos := make([]*UTXO, len(allUTXOBytes))
	for i, utxoBytes := range allUTXOBytes {
		utxo := &UTXO{}
		if _, err := a.codec.Unmarshal(utxoBytes, utxo); err != nil {
			return nil, ids.ShortID{}, ids.ID{}, fmt.Errorf("error parsing UTXO: %w", err)
		}
		utxos[i] = utxo
	}
	return utxos, lastAddrID, lastUTXOID, nil
}
