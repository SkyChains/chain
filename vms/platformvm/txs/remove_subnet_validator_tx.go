// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"

	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/snow"
	"github.com/luxdefi/luxd/utils/constants"
	"github.com/luxdefi/luxd/vms/components/verify"
)

var (
	_ UnsignedTx = (*RemoveSubnetValidatorTx)(nil)

	errRemovePrimaryNetworkValidator = errors.New("can't remove primary network validator with RemoveSubnetValidatorTx")
)

// Removes a validator from a subnet.
type RemoveSubnetValidatorTx struct {
	BaseTx `serialize:"true"`
	// The node to remove from the subnet.
	NodeID ids.NodeID `serialize:"true" json:"nodeID"`
	// The subnet to remove the node from.
	Subnet ids.ID `serialize:"true" json:"subnetID"`
	// Proves that the issuer has the right to remove the node from the subnet.
	SubnetAuth verify.Verifiable `serialize:"true" json:"subnetAuthorization"`
}

func (tx *RemoveSubnetValidatorTx) SyntacticVerify(ctx *snow.Context) error {
	switch {
	case tx == nil:
		return ErrNilTx
	case tx.SyntacticallyVerified:
		// already passed syntactic verification
		return nil
	case tx.Subnet == constants.PrimaryNetworkID:
		return errRemovePrimaryNetworkValidator
	}

	if err := tx.BaseTx.SyntacticVerify(ctx); err != nil {
		return err
	}
	if err := tx.SubnetAuth.Verify(); err != nil {
		return err
	}

	tx.SyntacticallyVerified = true
	return nil
}

func (tx *RemoveSubnetValidatorTx) Visit(visitor Visitor) error {
	return visitor.RemoveSubnetValidatorTx(tx)
}
