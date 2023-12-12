// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package p2p

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/utils/logging"
	"github.com/luxdefi/node/utils/set"
)

var _ ValidatorSet = (*testValidatorSet)(nil)

type testValidatorSet struct {
	validators set.Set[ids.NodeID]
}

func (t testValidatorSet) Has(_ context.Context, nodeID ids.NodeID) bool {
	return t.validators.Contains(nodeID)
}

func TestValidatorHandlerAppGossip(t *testing.T) {
	nodeID := ids.GenerateTestNodeID()
	validatorSet := set.Of(nodeID)

	tests := []struct {
		name         string
		validatorSet ValidatorSet
		nodeID       ids.NodeID
		expected     bool
	}{
		{
			name:         "message dropped",
			validatorSet: testValidatorSet{},
			nodeID:       nodeID,
		},
		{
			name: "message handled",
			validatorSet: testValidatorSet{
				validators: validatorSet,
			},
			nodeID:   nodeID,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			called := false
			handler := ValidatorHandler{
				Handler: testHandler{
					appGossipF: func(context.Context, ids.NodeID, []byte) {
						called = true
					},
				},
				ValidatorSet: tt.validatorSet,
				Log:          logging.NoLog{},
			}

			handler.AppGossip(context.Background(), tt.nodeID, []byte("foobar"))
			require.Equal(tt.expected, called)
		})
	}
}

func TestValidatorHandlerAppRequest(t *testing.T) {
	nodeID := ids.GenerateTestNodeID()
	validatorSet := set.Of(nodeID)

	tests := []struct {
		name         string
		validatorSet ValidatorSet
		nodeID       ids.NodeID
		expected     error
	}{
		{
			name:         "message dropped",
			validatorSet: testValidatorSet{},
			nodeID:       nodeID,
			expected:     ErrNotValidator,
		},
		{
			name: "message handled",
			validatorSet: testValidatorSet{
				validators: validatorSet,
			},
			nodeID: nodeID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			handler := ValidatorHandler{
				Handler:      NoOpHandler{},
				ValidatorSet: tt.validatorSet,
				Log:          logging.NoLog{},
			}

			_, err := handler.AppRequest(context.Background(), tt.nodeID, time.Time{}, []byte("foobar"))
			require.ErrorIs(err, tt.expected)
		})
	}
}
