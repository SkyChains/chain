// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"
	"math"

	"github.com/skychains/chain/codec"
	"github.com/skychains/chain/codec/linearcodec"
	"github.com/skychains/chain/utils/wrappers"
	"github.com/skychains/chain/vms/platformvm/signer"
	"github.com/skychains/chain/vms/platformvm/stakeable"
	"github.com/skychains/chain/vms/secp256k1fx"
)

const CodecVersion = 0

var (
	Codec codec.Manager

	// GenesisCodec allows txs of larger than usual size to be parsed.
	// While this gives flexibility in accommodating large genesis txs
	// it must not be used to parse new, unverified txs which instead
	// must be processed by Codec
	GenesisCodec codec.Manager
)

func init() {
	c := linearcodec.NewDefault()
	gc := linearcodec.NewDefault()

	errs := wrappers.Errs{}
	for _, c := range []linearcodec.Codec{c, gc} {
		// Order in which type are registered affect the byte representation
		// generated by marshalling ops. To maintain codec type ordering,
		// we skip positions for the blocks.
		c.SkipRegistrations(5)

		errs.Add(RegisterUnsignedTxsTypes(c))

		c.SkipRegistrations(4)

		errs.Add(RegisterDUnsignedTxsTypes(c))
	}

	Codec = codec.NewDefaultManager()
	GenesisCodec = codec.NewManager(math.MaxInt32)
	errs.Add(
		Codec.RegisterCodec(CodecVersion, c),
		GenesisCodec.RegisterCodec(CodecVersion, gc),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}

// RegisterUnsignedTxsTypes allows registering relevant type of unsigned package
// in the right sequence. Following repackaging of platformvm package, a few
// subpackage-level codecs were introduced, each handling serialization of
// specific types.
//
// RegisterUnsignedTxsTypes is made exportable so to guarantee that other codecs
// are coherent with components one.
func RegisterUnsignedTxsTypes(targetCodec linearcodec.Codec) error {
	errs := wrappers.Errs{}

	// The secp256k1fx is registered here because this is the same place it is
	// registered in the AVM. This ensures that the typeIDs match up for utxos
	// in shared memory.
	errs.Add(targetCodec.RegisterType(&secp256k1fx.TransferInput{}))
	targetCodec.SkipRegistrations(1)
	errs.Add(targetCodec.RegisterType(&secp256k1fx.TransferOutput{}))
	targetCodec.SkipRegistrations(1)
	errs.Add(
		targetCodec.RegisterType(&secp256k1fx.Credential{}),
		targetCodec.RegisterType(&secp256k1fx.Input{}),
		targetCodec.RegisterType(&secp256k1fx.OutputOwners{}),

		targetCodec.RegisterType(&AddValidatorTx{}),
		targetCodec.RegisterType(&AddSubnetValidatorTx{}),
		targetCodec.RegisterType(&AddDelegatorTx{}),
		targetCodec.RegisterType(&CreateChainTx{}),
		targetCodec.RegisterType(&CreateSubnetTx{}),
		targetCodec.RegisterType(&ImportTx{}),
		targetCodec.RegisterType(&ExportTx{}),
		targetCodec.RegisterType(&AdvanceTimeTx{}),
		targetCodec.RegisterType(&RewardValidatorTx{}),

		targetCodec.RegisterType(&stakeable.LockIn{}),
		targetCodec.RegisterType(&stakeable.LockOut{}),

		// Banff additions:
		targetCodec.RegisterType(&RemoveSubnetValidatorTx{}),
		targetCodec.RegisterType(&TransformSubnetTx{}),
		targetCodec.RegisterType(&AddPermissionlessValidatorTx{}),
		targetCodec.RegisterType(&AddPermissionlessDelegatorTx{}),

		targetCodec.RegisterType(&signer.Empty{}),
		targetCodec.RegisterType(&signer.ProofOfPossession{}),
	)
	return errs.Err
}

func RegisterDUnsignedTxsTypes(targetCodec linearcodec.Codec) error {
	return errors.Join(
		targetCodec.RegisterType(&TransferSubnetOwnershipTx{}),
		targetCodec.RegisterType(&BaseTx{}),
	)
}
