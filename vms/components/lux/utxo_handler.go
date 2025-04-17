// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package lux

import "github.com/skychains/chain/ids"

// Removes the UTXOs consumed by [ins] from the UTXO set
func Consume(utxoDB UTXODeleter, ins []*TransferableInput) {
	for _, input := range ins {
		utxoDB.DeleteUTXO(input.InputID())
	}
}

// Adds the UTXOs created by [outs] to the UTXO set.
// [txID] is the ID of the tx that created [outs].
func Produce(
	utxoDB UTXOAdder,
	txID ids.ID,
	outs []*TransferableOutput,
) {
	for index, out := range outs {
		utxoDB.AddUTXO(&UTXO{
			UTXOID: UTXOID{
				TxID:        txID,
				OutputIndex: uint32(index),
			},
			Asset: out.Asset,
			Out:   out.Output(),
		})
	}
}
