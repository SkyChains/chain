// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package rpcchainvm

import (
	"github.com/skychains/chain/database"
	"github.com/skychains/chain/snow/engine/snowman/block"

	vmpb "github.com/skychains/chain/proto/pb/vm"
)

var (
	errEnumToError = map[vmpb.Error]error{
		vmpb.Error_ERROR_CLOSED:                     database.ErrClosed,
		vmpb.Error_ERROR_NOT_FOUND:                  database.ErrNotFound,
		vmpb.Error_ERROR_STATE_SYNC_NOT_IMPLEMENTED: block.ErrStateSyncableVMNotImplemented,
	}
	errorToErrEnum = map[error]vmpb.Error{
		database.ErrClosed:                     vmpb.Error_ERROR_CLOSED,
		database.ErrNotFound:                   vmpb.Error_ERROR_NOT_FOUND,
		block.ErrStateSyncableVMNotImplemented: vmpb.Error_ERROR_STATE_SYNC_NOT_IMPLEMENTED,
	}
)

func errorToRPCError(err error) error {
	if _, ok := errorToErrEnum[err]; ok {
		return nil
	}
	return err
}
