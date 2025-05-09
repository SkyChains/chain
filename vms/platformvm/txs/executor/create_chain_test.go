// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/secp256k1"
	"github.com/skychains/chain/utils/hashing"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/utils/units"
	"github.com/skychains/chain/vms/platformvm/state"
	"github.com/skychains/chain/vms/platformvm/txs"
	"github.com/skychains/chain/vms/platformvm/txs/txstest"
	"github.com/skychains/chain/vms/platformvm/utxo"
	"github.com/skychains/chain/vms/secp256k1fx"

	walletsigner "github.com/skychains/chain/wallet/chain/p/signer"
)

// Ensure Execute fails when there are not enough control sigs
func TestCreateChainTxInsufficientControlSigs(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, banff)
	env.ctx.Lock.Lock()
	defer env.ctx.Lock.Unlock()

	builder, signer := env.factory.NewWallet(preFundedKeys[0], preFundedKeys[1])
	utx, err := builder.NewCreateChainTx(
		testSubnet1.ID(),
		nil,
		constants.AVMID,
		nil,
		"chain name",
	)
	require.NoError(err)
	tx, err := walletsigner.SignUnsigned(context.Background(), signer, utx)
	require.NoError(err)

	// Remove a signature
	tx.Creds[0].(*secp256k1fx.Credential).Sigs = tx.Creds[0].(*secp256k1fx.Credential).Sigs[1:]

	stateDiff, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	executor := StandardTxExecutor{
		Backend: &env.backend,
		State:   stateDiff,
		Tx:      tx,
	}
	err = tx.Unsigned.Visit(&executor)
	require.ErrorIs(err, errUnauthorizedSubnetModification)
}

// Ensure Execute fails when an incorrect control signature is given
func TestCreateChainTxWrongControlSig(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, banff)
	env.ctx.Lock.Lock()
	defer env.ctx.Lock.Unlock()

	builder, signer := env.factory.NewWallet(testSubnet1ControlKeys[0], testSubnet1ControlKeys[1])
	utx, err := builder.NewCreateChainTx(
		testSubnet1.ID(),
		nil,
		constants.AVMID,
		nil,
		"chain name",
	)
	require.NoError(err)
	tx, err := walletsigner.SignUnsigned(context.Background(), signer, utx)
	require.NoError(err)

	// Generate new, random key to sign tx with
	key, err := secp256k1.NewPrivateKey()
	require.NoError(err)

	// Replace a valid signature with one from another key
	sig, err := key.SignHash(hashing.ComputeHash256(tx.Unsigned.Bytes()))
	require.NoError(err)
	copy(tx.Creds[0].(*secp256k1fx.Credential).Sigs[0][:], sig)

	stateDiff, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	executor := StandardTxExecutor{
		Backend: &env.backend,
		State:   stateDiff,
		Tx:      tx,
	}
	err = tx.Unsigned.Visit(&executor)
	require.ErrorIs(err, errUnauthorizedSubnetModification)
}

// Ensure Execute fails when the Subnet the blockchain specifies as
// its validator set doesn't exist
func TestCreateChainTxNoSuchSubnet(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, banff)
	env.ctx.Lock.Lock()
	defer env.ctx.Lock.Unlock()

	builder, signer := env.factory.NewWallet(testSubnet1ControlKeys[0], testSubnet1ControlKeys[1])
	utx, err := builder.NewCreateChainTx(
		testSubnet1.ID(),
		nil,
		constants.AVMID,
		nil,
		"chain name",
	)
	require.NoError(err)
	tx, err := walletsigner.SignUnsigned(context.Background(), signer, utx)
	require.NoError(err)

	tx.Unsigned.(*txs.CreateChainTx).SubnetID = ids.GenerateTestID()

	stateDiff, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	executor := StandardTxExecutor{
		Backend: &env.backend,
		State:   stateDiff,
		Tx:      tx,
	}
	err = tx.Unsigned.Visit(&executor)
	require.ErrorIs(err, database.ErrNotFound)
}

// Ensure valid tx passes semanticVerify
func TestCreateChainTxValid(t *testing.T) {
	require := require.New(t)
	env := newEnvironment(t, banff)
	env.ctx.Lock.Lock()
	defer env.ctx.Lock.Unlock()

	builder, signer := env.factory.NewWallet(testSubnet1ControlKeys[0], testSubnet1ControlKeys[1])
	utx, err := builder.NewCreateChainTx(
		testSubnet1.ID(),
		nil,
		constants.AVMID,
		nil,
		"chain name",
	)
	require.NoError(err)
	tx, err := walletsigner.SignUnsigned(context.Background(), signer, utx)
	require.NoError(err)

	stateDiff, err := state.NewDiff(lastAcceptedID, env)
	require.NoError(err)

	executor := StandardTxExecutor{
		Backend: &env.backend,
		State:   stateDiff,
		Tx:      tx,
	}
	require.NoError(tx.Unsigned.Visit(&executor))
}

func TestCreateChainTxAP3FeeChange(t *testing.T) {
	ap3Time := defaultGenesisTime.Add(time.Hour)
	tests := []struct {
		name          string
		time          time.Time
		fee           uint64
		expectedError error
	}{
		{
			name:          "pre-fork - correctly priced",
			time:          defaultGenesisTime,
			fee:           0,
			expectedError: nil,
		},
		{
			name:          "post-fork - incorrectly priced",
			time:          ap3Time,
			fee:           100*defaultTxFee - 1*units.NanoLux,
			expectedError: utxo.ErrInsufficientUnlockedFunds,
		},
		{
			name:          "post-fork - correctly priced",
			time:          ap3Time,
			fee:           100 * defaultTxFee,
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			env := newEnvironment(t, banff)
			env.config.UpgradeConfig.ApricotPhase3Time = ap3Time

			addrs := set.NewSet[ids.ShortID](len(preFundedKeys))
			for _, key := range preFundedKeys {
				addrs.Add(key.Address())
			}

			env.state.SetTimestamp(test.time) // to duly set fee

			cfg := *env.config

			cfg.StaticFeeConfig.CreateBlockchainTxFee = test.fee
			factory := txstest.NewWalletFactory(env.ctx, &cfg, env.state)
			builder, signer := factory.NewWallet(preFundedKeys...)
			utx, err := builder.NewCreateChainTx(
				testSubnet1.ID(),
				nil,
				ids.GenerateTestID(),
				nil,
				"",
			)
			require.NoError(err)
			tx, err := walletsigner.SignUnsigned(context.Background(), signer, utx)
			require.NoError(err)

			stateDiff, err := state.NewDiff(lastAcceptedID, env)
			require.NoError(err)

			stateDiff.SetTimestamp(test.time)

			executor := StandardTxExecutor{
				Backend: &env.backend,
				State:   stateDiff,
				Tx:      tx,
			}
			err = tx.Unsigned.Visit(&executor)
			require.ErrorIs(err, test.expectedError)
		})
	}
}
