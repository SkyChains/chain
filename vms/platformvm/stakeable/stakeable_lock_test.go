// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package stakeable

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/skychains/chain/vms/components/lux"
)

var errTest = errors.New("hi mom")

func TestLockOutVerify(t *testing.T) {
	tests := []struct {
		name             string
		locktime         uint64
		transferableOutF func(*gomock.Controller) lux.TransferableOut
		expectedErr      error
	}{
		{
			name:     "happy path",
			locktime: 1,
			transferableOutF: func(ctrl *gomock.Controller) lux.TransferableOut {
				o := lux.NewMockTransferableOut(ctrl)
				o.EXPECT().Verify().Return(nil)
				return o
			},
			expectedErr: nil,
		},
		{
			name:     "invalid locktime",
			locktime: 0,
			transferableOutF: func(*gomock.Controller) lux.TransferableOut {
				return nil
			},
			expectedErr: errInvalidLocktime,
		},
		{
			name:     "nested",
			locktime: 1,
			transferableOutF: func(*gomock.Controller) lux.TransferableOut {
				return &LockOut{}
			},
			expectedErr: errNestedStakeableLocks,
		},
		{
			name:     "inner output fails verification",
			locktime: 1,
			transferableOutF: func(ctrl *gomock.Controller) lux.TransferableOut {
				o := lux.NewMockTransferableOut(ctrl)
				o.EXPECT().Verify().Return(errTest)
				return o
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			lockOut := &LockOut{
				Locktime:        tt.locktime,
				TransferableOut: tt.transferableOutF(ctrl),
			}
			require.Equal(t, tt.expectedErr, lockOut.Verify())
		})
	}
}

func TestLockInVerify(t *testing.T) {
	tests := []struct {
		name            string
		locktime        uint64
		transferableInF func(*gomock.Controller) lux.TransferableIn
		expectedErr     error
	}{
		{
			name:     "happy path",
			locktime: 1,
			transferableInF: func(ctrl *gomock.Controller) lux.TransferableIn {
				o := lux.NewMockTransferableIn(ctrl)
				o.EXPECT().Verify().Return(nil)
				return o
			},
			expectedErr: nil,
		},
		{
			name:     "invalid locktime",
			locktime: 0,
			transferableInF: func(*gomock.Controller) lux.TransferableIn {
				return nil
			},
			expectedErr: errInvalidLocktime,
		},
		{
			name:     "nested",
			locktime: 1,
			transferableInF: func(*gomock.Controller) lux.TransferableIn {
				return &LockIn{}
			},
			expectedErr: errNestedStakeableLocks,
		},
		{
			name:     "inner input fails verification",
			locktime: 1,
			transferableInF: func(ctrl *gomock.Controller) lux.TransferableIn {
				o := lux.NewMockTransferableIn(ctrl)
				o.EXPECT().Verify().Return(errTest)
				return o
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			lockOut := &LockIn{
				Locktime:       tt.locktime,
				TransferableIn: tt.transferableInF(ctrl),
			}
			require.Equal(t, tt.expectedErr, lockOut.Verify())
		})
	}
}
