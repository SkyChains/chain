// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/bls"
)

func TestMessage(t *testing.T) {
	require := require.New(t)

	unsignedMsg, err := NewUnsignedMessage(
		constants.UnitTestID,
		ids.GenerateTestID(),
		[]byte("payload"),
	)
	require.NoError(err)

	msg, err := NewMessage(
		unsignedMsg,
		&BitSetSignature{
			Signers:   []byte{1, 2, 3},
			Signature: [bls.SignatureLen]byte{4, 5, 6},
		},
	)
	require.NoError(err)

	msgBytes := msg.Bytes()
	msg2, err := ParseMessage(msgBytes)
	require.NoError(err)
	require.Equal(msg, msg2)
}

func TestParseMessageJunk(t *testing.T) {
	require := require.New(t)

	bytes := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	_, err := ParseMessage(bytes)
	require.ErrorIs(err, codec.ErrUnknownVersion)
}
