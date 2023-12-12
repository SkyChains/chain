// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package tx

import (
	"github.com/luxdefi/node/cache"
	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/utils/crypto/secp256k1"
	"github.com/luxdefi/node/utils/hashing"
)

var secpCache = secp256k1.RecoverCache{
	LRU: cache.LRU[ids.ID, *secp256k1.PublicKey]{
		Size: 2048,
	},
}

type Tx struct {
	Unsigned  `serialize:"true" json:"unsigned"`
	Signature [secp256k1.SignatureLen]byte `serialize:"true" json:"signature"`
}

func Parse(bytes []byte) (*Tx, error) {
	tx := &Tx{}
	_, err := Codec.Unmarshal(bytes, tx)
	return tx, err
}

func Sign(utx Unsigned, key *secp256k1.PrivateKey) (*Tx, error) {
	unsignedBytes, err := Codec.Marshal(Version, &utx)
	if err != nil {
		return nil, err
	}

	sig, err := key.Sign(unsignedBytes)
	if err != nil {
		return nil, err
	}

	tx := &Tx{
		Unsigned: utx,
	}
	copy(tx.Signature[:], sig)
	return tx, nil
}

func (tx *Tx) ID() (ids.ID, error) {
	bytes, err := Codec.Marshal(Version, tx)
	return hashing.ComputeHash256Array(bytes), err
}

func (tx *Tx) SenderID() (ids.ShortID, error) {
	unsignedBytes, err := Codec.Marshal(Version, &tx.Unsigned)
	if err != nil {
		return ids.ShortEmpty, err
	}

	pk, err := secpCache.RecoverPublicKey(unsignedBytes, tx.Signature[:])
	if err != nil {
		return ids.ShortEmpty, err
	}
	return pk.Address(), nil
}
