// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package xsvm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"go.uber.org/zap"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/versiondb"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/consensus/snowman"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/json"
	"github.com/skychains/chain/version"
	"github.com/skychains/chain/vms/example/xsvm/api"
	"github.com/skychains/chain/vms/example/xsvm/builder"
	"github.com/skychains/chain/vms/example/xsvm/chain"
	"github.com/skychains/chain/vms/example/xsvm/execute"
	"github.com/skychains/chain/vms/example/xsvm/genesis"
	"github.com/skychains/chain/vms/example/xsvm/state"

	smblock "github.com/skychains/chain/snow/engine/snowman/block"
	xsblock "github.com/skychains/chain/vms/example/xsvm/block"
)

var (
	_ smblock.ChainVM                      = (*VM)(nil)
	_ smblock.BuildBlockWithContextChainVM = (*VM)(nil)
)

type VM struct {
	common.AppHandler

	chainContext *snow.Context
	db           database.Database
	genesis      *genesis.Genesis
	engineChan   chan<- common.Message

	chain   chain.Chain
	builder builder.Builder
}

func (vm *VM) Initialize(
	_ context.Context,
	chainContext *snow.Context,
	db database.Database,
	genesisBytes []byte,
	_ []byte,
	_ []byte,
	engineChan chan<- common.Message,
	_ []*common.Fx,
	_ common.AppSender,
) error {
	vm.AppHandler = common.NewNoOpAppHandler(chainContext.Log)

	chainContext.Log.Info("initializing xsvm",
		zap.Stringer("version", Version),
	)

	vm.chainContext = chainContext
	vm.db = db
	g, err := genesis.Parse(genesisBytes)
	if err != nil {
		return fmt.Errorf("failed to parse genesis bytes: %w", err)
	}

	vdb := versiondb.New(vm.db)
	if err := execute.Genesis(vdb, chainContext.ChainID, g); err != nil {
		return fmt.Errorf("failed to initialize genesis state: %w", err)
	}
	if err := vdb.Commit(); err != nil {
		return err
	}

	vm.genesis = g
	vm.engineChan = engineChan

	vm.chain, err = chain.New(chainContext, vm.db)
	if err != nil {
		return fmt.Errorf("failed to initialize chain manager: %w", err)
	}

	vm.builder = builder.New(chainContext, engineChan, vm.chain)

	chainContext.Log.Info("initialized xsvm",
		zap.Stringer("lastAcceptedID", vm.chain.LastAccepted()),
	)
	return nil
}

func (vm *VM) SetState(_ context.Context, state snow.State) error {
	vm.chain.SetChainState(state)
	return nil
}

func (vm *VM) Shutdown(context.Context) error {
	if vm.chainContext == nil {
		return nil
	}
	return vm.db.Close()
}

func (*VM) Version(context.Context) (string, error) {
	return Version.String(), nil
}

func (vm *VM) CreateHandlers(context.Context) (map[string]http.Handler, error) {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	api := api.NewServer(
		vm.chainContext,
		vm.genesis,
		vm.db,
		vm.chain,
		vm.builder,
	)
	return map[string]http.Handler{
		"": server,
	}, server.RegisterService(api, constants.XSVMName)
}

func (*VM) HealthCheck(context.Context) (interface{}, error) {
	return http.StatusOK, nil
}

func (*VM) Connected(context.Context, ids.NodeID, *version.Application) error {
	return nil
}

func (*VM) Disconnected(context.Context, ids.NodeID) error {
	return nil
}

func (vm *VM) GetBlock(_ context.Context, blkID ids.ID) (snowman.Block, error) {
	return vm.chain.GetBlock(blkID)
}

func (vm *VM) ParseBlock(_ context.Context, blkBytes []byte) (snowman.Block, error) {
	blk, err := xsblock.Parse(blkBytes)
	if err != nil {
		return nil, err
	}
	return vm.chain.NewBlock(blk)
}

func (vm *VM) BuildBlock(ctx context.Context) (snowman.Block, error) {
	return vm.builder.BuildBlock(ctx, nil)
}

func (vm *VM) SetPreference(_ context.Context, preferred ids.ID) error {
	vm.builder.SetPreference(preferred)
	return nil
}

func (vm *VM) LastAccepted(context.Context) (ids.ID, error) {
	return vm.chain.LastAccepted(), nil
}

func (vm *VM) BuildBlockWithContext(ctx context.Context, blockContext *smblock.Context) (snowman.Block, error) {
	return vm.builder.BuildBlock(ctx, blockContext)
}

func (vm *VM) GetBlockIDAtHeight(_ context.Context, height uint64) (ids.ID, error) {
	return state.GetBlockIDByHeight(vm.db, height)
}
