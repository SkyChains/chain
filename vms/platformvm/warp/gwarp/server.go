// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package gwarp

import (
	"context"

	"github.com/SkyChains/chain/ids"
	"github.com/SkyChains/chain/vms/platformvm/warp"

	pb "github.com/SkyChains/chain/proto/pb/warp"
)

var _ pb.SignerServer = (*Server)(nil)

type Server struct {
	pb.UnsafeSignerServer
	signer warp.Signer
}

func NewServer(signer warp.Signer) *Server {
	return &Server{signer: signer}
}

func (s *Server) Sign(_ context.Context, unsignedMsg *pb.SignRequest) (*pb.SignResponse, error) {
	sourceChainID, err := ids.ToID(unsignedMsg.SourceChainId)
	if err != nil {
		return nil, err
	}

	msg, err := warp.NewUnsignedMessage(
		unsignedMsg.NetworkId,
		sourceChainID,
		unsignedMsg.Payload,
	)
	if err != nil {
		return nil, err
	}

	sig, err := s.signer.Sign(msg)
	return &pb.SignResponse{
		Signature: sig,
	}, err
}
