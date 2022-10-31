// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luxdefi/luxd/message"
	"github.com/luxdefi/luxd/utils/logging"
)

func TestBlockingMessageQueue(t *testing.T) {
	require := require.New(t)

	q := NewBlockingMessageQueue(
		SendFailedFunc(func(msg message.OutboundMessage) {
			t.Fail()
		}),
		logging.NoLog{},
		0,
	)

	mc := newMessageCreator(t)
	msg, err := mc.Ping()
	require.NoError(err)

	numToSend := 10
	go func() {
		for i := 0; i < numToSend; i++ {
			q.Push(context.Background(), msg)
		}
	}()

	for i := 0; i < numToSend; i++ {
		_, ok := q.Pop()
		require.True(ok)
	}

	_, ok := q.PopNow()
	require.False(ok)

	q.Close()

	_, ok = q.Pop()
	require.False(ok)
}
