// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avax

import (
	"bytes"
	"errors"
	"fmt"
<<<<<<< HEAD
=======
	"sort"
>>>>>>> 7b681477c (Add `avax.UTXOIDFromString` helper (#2138))
	"strconv"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/vms/components/verify"
)

var (
	errNilUTXOID                 = errors.New("nil utxo ID is not valid")
	errMalformedUTXOIDString     = errors.New("unexpected number of tokens in string")
	errFailedDecodingUTXOIDTxID  = errors.New("failed decoding UTXOID TxID")
	errFailedDecodingUTXOIDIndex = errors.New("failed decoding UTXOID index")

	_ verify.Verifiable       = (*UTXOID)(nil)
	_ utils.Sortable[*UTXOID] = (*UTXOID)(nil)
)

type UTXOID struct {
	// Serialized:
	TxID        ids.ID `serialize:"true" json:"txID"`
	OutputIndex uint32 `serialize:"true" json:"outputIndex"`

	// Symbol is false if the UTXO should be part of the DB
	Symbol bool `json:"-"`
	// id is the unique ID of a UTXO, it is calculated from TxID and OutputIndex
	id ids.ID
}

// InputSource returns the source of the UTXO that this input is spending
func (utxo *UTXOID) InputSource() (ids.ID, uint32) {
	return utxo.TxID, utxo.OutputIndex
}

// InputID returns a unique ID of the UTXO that this input is spending
func (utxo *UTXOID) InputID() ids.ID {
	if utxo.id == ids.Empty {
		utxo.id = utxo.TxID.Prefix(uint64(utxo.OutputIndex))
	}
	return utxo.id
}

// Symbolic returns if this is the ID of a UTXO in the DB, or if it is a
// symbolic input
func (utxo *UTXOID) Symbolic() bool {
	return utxo.Symbol
}

func (utxo *UTXOID) String() string {
	return fmt.Sprintf("%s:%d", utxo.TxID, utxo.OutputIndex)
}

// UTXOIDFromString attempts to parse a string into a UTXOID
func UTXOIDFromString(s string) (*UTXOID, error) {
	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return nil, errMalformedUTXOIDString
	}

	txID, err := ids.FromString(ss[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errFailedDecodingUTXOIDTxID, err)
	}

	idx, err := strconv.ParseUint(ss[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errFailedDecodingUTXOIDIndex, err)
	}

	return &UTXOID{
		TxID:        txID,
		OutputIndex: uint32(idx),
	}, nil
}

func (utxo *UTXOID) Verify() error {
	switch {
	case utxo == nil:
		return errNilUTXOID
	default:
		return nil
	}
}

func (utxo *UTXOID) Less(other *UTXOID) bool {
	utxoID, utxoIndex := utxo.InputSource()
	otherID, otherIndex := other.InputSource()

	switch bytes.Compare(utxoID[:], otherID[:]) {
	case -1:
		return true
	case 0:
		return utxoIndex < otherIndex
	default:
		return false
	}
}
<<<<<<< HEAD
=======

func (utxos innerSortUTXOIDs) Len() int {
	return len(utxos)
}

func (utxos innerSortUTXOIDs) Swap(i, j int) {
	utxos[j], utxos[i] = utxos[i], utxos[j]
}

func SortUTXOIDs(utxos []*UTXOID) {
	sort.Sort(innerSortUTXOIDs(utxos))
}

func IsSortedAndUniqueUTXOIDs(utxos []*UTXOID) bool {
	return utils.IsSortedAndUnique(innerSortUTXOIDs(utxos))
}
>>>>>>> 55bd9343c (Add EmptyLines linter (#2233))
