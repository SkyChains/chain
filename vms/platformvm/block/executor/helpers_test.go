// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/skychains/chain/chains"
	"github.com/skychains/chain/chains/atomic"
	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/codec/linearcodec"
	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/memdb"
	"github.com/skychains/chain/database/prefixdb"
	"github.com/skychains/chain/database/versiondb"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow"
	"github.com/skychains/chain/snow/engine/common"
	"github.com/skychains/chain/snow/snowtest"
	"github.com/skychains/chain/snow/uptime"
	"github.com/skychains/chain/snow/validators"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/secp256k1"
	"github.com/skychains/chain/utils/formatting"
	"github.com/skychains/chain/utils/formatting/address"
	"github.com/skychains/chain/utils/json"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/timer/mockable"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/platformvm/api"
	"github.com/skychains/chain/vms/platformvm/config"
	"github.com/skychains/chain/vms/platformvm/fx"
	"github.com/skychains/chain/vms/platformvm/metrics"
	"github.com/skychains/chain/vms/platformvm/reward"
	"github.com/skychains/chain/vms/platformvm/state"
	"github.com/skychains/chain/vms/platformvm/status"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/platformvm/txs/executor"
	"github.com/skychains/chain/vms/platformvm/txs/fee"
	"github.com/skychains/chain/vms/platformvm/txs/mempool"
	"github.com/skychains/chain/vms/platformvm/txs/txstest"
	"github.com/skychains/chain/vms/platformvm/upgrade"
	"github.com/skychains/chain/vms/platformvm/utxo"
	"github.com/skychains/chain/vms/secp256k1fx"

	pvalidators "github.com/skychains/chain/vms/platformvm/validators"
	walletsigner "github.com/skychains/chain/wallet/chain/p/signer"
	walletcommon "github.com/skychains/chain/wallet/subnet/primary/common"
)

const (
	pending stakerStatus = iota
	current

	defaultWeight = 10000
	trackChecksum = false

	apricotPhase3 fork = iota
	apricotPhase5
	banff
	cortina
	durango
	eUpgrade
)

var (
	defaultMinStakingDuration = 24 * time.Hour
	defaultMaxStakingDuration = 365 * 24 * time.Hour
	defaultGenesisTime        = time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultValidateStartTime  = defaultGenesisTime
	defaultValidateEndTime    = defaultValidateStartTime.Add(10 * defaultMinStakingDuration)
	defaultMinValidatorStake  = 5 * units.MilliLux
	defaultBalance            = 100 * defaultMinValidatorStake
	preFundedKeys             = secp256k1.TestKeys()
	luxAssetID               = ids.ID{'y', 'e', 'e', 't'}
	defaultTxFee              = uint64(100)

	genesisBlkID ids.ID
	testSubnet1  *txs.Tx

	// Node IDs of genesis validators. Initialized in init function
	genesisNodeIDs []ids.NodeID
)

func init() {
	genesisNodeIDs = make([]ids.NodeID, len(preFundedKeys))
	for i := range preFundedKeys {
		genesisNodeIDs[i] = ids.GenerateTestNodeID()
	}
}

type stakerStatus uint

type fork uint8

type staker struct {
	nodeID             ids.NodeID
	rewardAddress      ids.ShortID
	startTime, endTime time.Time
}

type test struct {
	description           string
	stakers               []staker
	subnetStakers         []staker
	advanceTimeTo         []time.Time
	expectedStakers       map[ids.NodeID]stakerStatus
	expectedSubnetStakers map[ids.NodeID]stakerStatus
}

type environment struct {
	blkManager Manager
	mempool    mempool.Mempool
	sender     *common.SenderTest

	isBootstrapped *utils.Atomic[bool]
	config         *config.Config
	clk            *mockable.Clock
	baseDB         *versiondb.Database
	ctx            *snow.Context
	fx             fx.Fx
	state          state.State
	mockedState    *state.MockState
	uptimes        uptime.Manager
	utxosVerifier  utxo.Verifier
	factory        *txstest.WalletFactory
	backend        *executor.Backend
}

func newEnvironment(t *testing.T, ctrl *gomock.Controller, f fork) *environment {
	res := &environment{
		isBootstrapped: &utils.Atomic[bool]{},
		config:         defaultConfig(t, f),
		clk:            defaultClock(),
	}
	res.isBootstrapped.Set(true)

	res.baseDB = versiondb.New(memdb.New())
	atomicDB := prefixdb.New([]byte{1}, res.baseDB)
	m := atomic.NewMemory(atomicDB)

	res.ctx = snowtest.Context(t, snowtest.PChainID)
	res.ctx.LUXAssetID = luxAssetID
	res.ctx.SharedMemory = m.NewSharedMemory(res.ctx.ChainID)

	res.fx = defaultFx(res.clk, res.ctx.Log, res.isBootstrapped.Get())

	rewardsCalc := reward.NewCalculator(res.config.RewardConfig)

	if ctrl == nil {
		res.state = defaultState(res.config, res.ctx, res.baseDB, rewardsCalc)
		res.uptimes = uptime.NewManager(res.state, res.clk)
		res.utxosVerifier = utxo.NewVerifier(res.ctx, res.clk, res.fx)
		res.factory = txstest.NewWalletFactory(
			res.ctx,
			res.config,
			res.state,
		)
	} else {
		genesisBlkID = ids.GenerateTestID()
		res.mockedState = state.NewMockState(ctrl)
		res.uptimes = uptime.NewManager(res.mockedState, res.clk)
		res.utxosVerifier = utxo.NewVerifier(res.ctx, res.clk, res.fx)
		res.factory = txstest.NewWalletFactory(
			res.ctx,
			res.config,
			res.mockedState,
		)

		// setup expectations strictly needed for environment creation
		res.mockedState.EXPECT().GetLastAccepted().Return(genesisBlkID).Times(1)
	}

	res.backend = &executor.Backend{
		Config:       res.config,
		Ctx:          res.ctx,
		Clk:          res.clk,
		Bootstrapped: res.isBootstrapped,
		Fx:           res.fx,
		FlowChecker:  res.utxosVerifier,
		Uptimes:      res.uptimes,
		Rewards:      rewardsCalc,
	}

	registerer := prometheus.NewRegistry()
	res.sender = &common.SenderTest{T: t}

	metrics := metrics.Noop

	var err error
	res.mempool, err = mempool.New("mempool", registerer, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create mempool: %w", err))
	}

	if ctrl == nil {
		res.blkManager = NewManager(
			res.mempool,
			metrics,
			res.state,
			res.backend,
			pvalidators.TestManager,
		)
		addSubnet(res)
	} else {
		res.blkManager = NewManager(
			res.mempool,
			metrics,
			res.mockedState,
			res.backend,
			pvalidators.TestManager,
		)
		// we do not add any subnet to state, since we can mock
		// whatever we need
	}

	t.Cleanup(func() {
		res.ctx.Lock.Lock()
		defer res.ctx.Lock.Unlock()

		if res.mockedState != nil {
			// state is mocked, nothing to do here
			return
		}

		require := require.New(t)

		if res.isBootstrapped.Get() {
			validatorIDs := res.config.Validators.GetValidatorIDs(constants.PrimaryNetworkID)

			require.NoError(res.uptimes.StopTracking(validatorIDs, constants.PrimaryNetworkID))
			require.NoError(res.state.Commit())
		}

		if res.state != nil {
			require.NoError(res.state.Close())
		}

		require.NoError(res.baseDB.Close())
	})

	return res
}

func addSubnet(env *environment) {
	builder, signer := env.factory.NewWallet(preFundedKeys[0])
	utx, err := builder.NewCreateSubnetTx(
		&secp256k1fx.OutputOwners{
			Threshold: 2,
			Addrs: []ids.ShortID{
				preFundedKeys[0].PublicKey().Address(),
				preFundedKeys[1].PublicKey().Address(),
				preFundedKeys[2].PublicKey().Address(),
			},
		},
		walletcommon.WithChangeOwner(&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
		}),
	)
	if err != nil {
		panic(err)
	}
	testSubnet1, err = walletsigner.SignUnsigned(context.Background(), signer, utx)
	if err != nil {
		panic(err)
	}

	genesisID := env.state.GetLastAccepted()
	stateDiff, err := state.NewDiff(genesisID, env.blkManager)
	if err != nil {
		panic(err)
	}

	executor := executor.StandardTxExecutor{
		Backend: env.backend,
		State:   stateDiff,
		Tx:      testSubnet1,
	}
	err = testSubnet1.Unsigned.Visit(&executor)
	if err != nil {
		panic(err)
	}

	stateDiff.AddTx(testSubnet1, status.Committed)
	if err := stateDiff.Apply(env.state); err != nil {
		panic(err)
	}
}

func defaultState(
	cfg *config.Config,
	ctx *snow.Context,
	db database.Database,
	rewards reward.Calculator,
) state.State {
	genesisBytes := buildGenesisTest(ctx)
	execCfg, _ := config.GetExecutionConfig([]byte(`{}`))
	state, err := state.New(
		db,
		genesisBytes,
		prometheus.NewRegistry(),
		cfg,
		execCfg,
		ctx,
		metrics.Noop,
		rewards,
	)
	if err != nil {
		panic(err)
	}

	// persist and reload to init a bunch of in-memory stuff
	state.SetHeight(0)
	if err := state.Commit(); err != nil {
		panic(err)
	}
	genesisBlkID = state.GetLastAccepted()
	return state
}

func defaultConfig(t *testing.T, f fork) *config.Config {
	c := &config.Config{
		Chains:                 chains.TestManager,
		UptimeLockedCalculator: uptime.NewLockedCalculator(),
		Validators:             validators.NewManager(),
		StaticFeeConfig: fee.StaticConfig{
			TxFee:                 defaultTxFee,
			CreateSubnetTxFee:     100 * defaultTxFee,
			CreateBlockchainTxFee: 100 * defaultTxFee,
		},
		MinValidatorStake: 5 * units.MilliLux,
		MaxValidatorStake: 500 * units.MilliLux,
		MinDelegatorStake: 1 * units.MilliLux,
		MinStakeDuration:  defaultMinStakingDuration,
		MaxStakeDuration:  defaultMaxStakingDuration,
		RewardConfig: reward.Config{
			MaxConsumptionRate: .12 * reward.PercentDenominator,
			MinConsumptionRate: .10 * reward.PercentDenominator,
			MintingPeriod:      365 * 24 * time.Hour,
			SupplyCap:          720 * units.MegaLux,
		},
		UpgradeConfig: upgrade.Config{
			ApricotPhase3Time: mockable.MaxTime,
			ApricotPhase5Time: mockable.MaxTime,
			BanffTime:         mockable.MaxTime,
			CortinaTime:       mockable.MaxTime,
			DurangoTime:       mockable.MaxTime,
			EUpgradeTime:      mockable.MaxTime,
		},
	}

	switch f {
	case eUpgrade:
		c.UpgradeConfig.EUpgradeTime = time.Time{} // neglecting fork ordering this for package tests
		fallthrough
	case durango:
		c.UpgradeConfig.DurangoTime = time.Time{} // neglecting fork ordering for this package's tests
		fallthrough
	case cortina:
		c.UpgradeConfig.CortinaTime = time.Time{} // neglecting fork ordering for this package's tests
		fallthrough
	case banff:
		c.UpgradeConfig.BanffTime = time.Time{} // neglecting fork ordering for this package's tests
		fallthrough
	case apricotPhase5:
		c.UpgradeConfig.ApricotPhase5Time = defaultValidateEndTime
		fallthrough
	case apricotPhase3:
		c.UpgradeConfig.ApricotPhase3Time = defaultValidateEndTime
	default:
		require.FailNow(t, "unhandled fork", f)
	}

	return c
}

func defaultClock() *mockable.Clock {
	clk := &mockable.Clock{}
	clk.Set(defaultGenesisTime)
	return clk
}

type fxVMInt struct {
	registry codec.Registry
	clk      *mockable.Clock
	log      logging.Logger
}

func (fvi *fxVMInt) CodecRegistry() codec.Registry {
	return fvi.registry
}

func (fvi *fxVMInt) Clock() *mockable.Clock {
	return fvi.clk
}

func (fvi *fxVMInt) Logger() logging.Logger {
	return fvi.log
}

func defaultFx(clk *mockable.Clock, log logging.Logger, isBootstrapped bool) fx.Fx {
	fxVMInt := &fxVMInt{
		registry: linearcodec.NewDefault(),
		clk:      clk,
		log:      log,
	}
	res := &secp256k1fx.Fx{}
	if err := res.Initialize(fxVMInt); err != nil {
		panic(err)
	}
	if isBootstrapped {
		if err := res.Bootstrapped(); err != nil {
			panic(err)
		}
	}
	return res
}

func buildGenesisTest(ctx *snow.Context) []byte {
	genesisUTXOs := make([]api.UTXO, len(preFundedKeys))
	for i, key := range preFundedKeys {
		id := key.PublicKey().Address()
		addr, err := address.FormatBech32(constants.UnitTestHRP, id.Bytes())
		if err != nil {
			panic(err)
		}
		genesisUTXOs[i] = api.UTXO{
			Amount:  json.Uint64(defaultBalance),
			Address: addr,
		}
	}

	genesisValidators := make([]api.GenesisPermissionlessValidator, len(genesisNodeIDs))
	for i, nodeID := range genesisNodeIDs {
		addr, err := address.FormatBech32(constants.UnitTestHRP, nodeID.Bytes())
		if err != nil {
			panic(err)
		}
		genesisValidators[i] = api.GenesisPermissionlessValidator{
			GenesisValidator: api.GenesisValidator{
				StartTime: json.Uint64(defaultValidateStartTime.Unix()),
				EndTime:   json.Uint64(defaultValidateEndTime.Unix()),
				NodeID:    nodeID,
			},
			RewardOwner: &api.Owner{
				Threshold: 1,
				Addresses: []string{addr},
			},
			Staked: []api.UTXO{{
				Amount:  json.Uint64(defaultWeight),
				Address: addr,
			}},
			DelegationFee: reward.PercentDenominator,
		}
	}

	buildGenesisArgs := api.BuildGenesisArgs{
		NetworkID:     json.Uint32(constants.UnitTestID),
		LuxAssetID:   ctx.LUXAssetID,
		UTXOs:         genesisUTXOs,
		Validators:    genesisValidators,
		Chains:        nil,
		Time:          json.Uint64(defaultGenesisTime.Unix()),
		InitialSupply: json.Uint64(360 * units.MegaLux),
		Encoding:      formatting.Hex,
	}

	buildGenesisResponse := api.BuildGenesisReply{}
	platformvmSS := api.StaticService{}
	if err := platformvmSS.BuildGenesis(nil, &buildGenesisArgs, &buildGenesisResponse); err != nil {
		panic(fmt.Errorf("problem while building platform chain's genesis state: %w", err))
	}

	genesisBytes, err := formatting.Decode(buildGenesisResponse.Encoding, buildGenesisResponse.Bytes)
	if err != nil {
		panic(err)
	}

	return genesisBytes
}

func addPendingValidator(
	env *environment,
	startTime time.Time,
	endTime time.Time,
	nodeID ids.NodeID,
	rewardAddress ids.ShortID,
	keys []*secp256k1.PrivateKey,
) (*txs.Tx, error) {
	builder, signer := env.factory.NewWallet(keys...)
	utx, err := builder.NewAddValidatorTx(
		&txs.Validator{
			NodeID: nodeID,
			Start:  uint64(startTime.Unix()),
			End:    uint64(endTime.Unix()),
			Wght:   env.config.MinValidatorStake,
		},
		&secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{rewardAddress},
		},
		reward.PercentDenominator,
	)
	if err != nil {
		return nil, err
	}
	addPendingValidatorTx, err := walletsigner.SignUnsigned(context.Background(), signer, utx)
	if err != nil {
		return nil, err
	}

	staker, err := state.NewPendingStaker(
		addPendingValidatorTx.ID(),
		addPendingValidatorTx.Unsigned.(*txs.AddValidatorTx),
	)
	if err != nil {
		return nil, err
	}

	env.state.PutPendingValidator(staker)
	env.state.AddTx(addPendingValidatorTx, status.Committed)
	dummyHeight := uint64(1)
	env.state.SetHeight(dummyHeight)
	if err := env.state.Commit(); err != nil {
		return nil, err
	}
	return addPendingValidatorTx, nil
}
