// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"errors"
	"time"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/proto/pb/p2p"
)

var (
	errMissingField = errors.New("message missing field")

	_ chainIDGetter = (*p2p.GetStateSummaryFrontier)(nil)
	_ chainIDGetter = (*p2p.StateSummaryFrontier)(nil)
	_ chainIDGetter = (*p2p.GetAcceptedStateSummary)(nil)
	_ chainIDGetter = (*p2p.AcceptedStateSummary)(nil)
	_ chainIDGetter = (*p2p.GetAcceptedFrontier)(nil)
	_ chainIDGetter = (*p2p.AcceptedFrontier)(nil)
	_ chainIDGetter = (*p2p.GetAccepted)(nil)
	_ chainIDGetter = (*p2p.Accepted)(nil)
	_ chainIDGetter = (*p2p.GetAncestors)(nil)
	_ chainIDGetter = (*p2p.Ancestors)(nil)
	_ chainIDGetter = (*p2p.Get)(nil)
	_ chainIDGetter = (*p2p.Put)(nil)
	_ chainIDGetter = (*p2p.PushQuery)(nil)
	_ chainIDGetter = (*p2p.PullQuery)(nil)
	_ chainIDGetter = (*p2p.Chits)(nil)
	_ chainIDGetter = (*p2p.AppRequest)(nil)
	_ chainIDGetter = (*p2p.AppResponse)(nil)
	_ chainIDGetter = (*p2p.AppGossip)(nil)

	_ requestIDGetter = (*p2p.GetStateSummaryFrontier)(nil)
	_ requestIDGetter = (*p2p.StateSummaryFrontier)(nil)
	_ requestIDGetter = (*p2p.GetAcceptedStateSummary)(nil)
	_ requestIDGetter = (*p2p.AcceptedStateSummary)(nil)
	_ requestIDGetter = (*p2p.GetAcceptedFrontier)(nil)
	_ requestIDGetter = (*p2p.AcceptedFrontier)(nil)
	_ requestIDGetter = (*p2p.GetAccepted)(nil)
	_ requestIDGetter = (*p2p.Accepted)(nil)
	_ requestIDGetter = (*p2p.GetAncestors)(nil)
	_ requestIDGetter = (*p2p.Ancestors)(nil)
	_ requestIDGetter = (*p2p.Get)(nil)
	_ requestIDGetter = (*p2p.Put)(nil)
	_ requestIDGetter = (*p2p.PushQuery)(nil)
	_ requestIDGetter = (*p2p.PullQuery)(nil)
	_ requestIDGetter = (*p2p.Chits)(nil)
	_ requestIDGetter = (*p2p.AppRequest)(nil)
	_ requestIDGetter = (*p2p.AppResponse)(nil)

	_ engineTypeGetter = (*p2p.GetAncestors)(nil)

	_ deadlineGetter = (*p2p.GetStateSummaryFrontier)(nil)
	_ deadlineGetter = (*p2p.GetAcceptedStateSummary)(nil)
	_ deadlineGetter = (*p2p.GetAcceptedFrontier)(nil)
	_ deadlineGetter = (*p2p.GetAccepted)(nil)
	_ deadlineGetter = (*p2p.GetAncestors)(nil)
	_ deadlineGetter = (*p2p.Get)(nil)
	_ deadlineGetter = (*p2p.PushQuery)(nil)
	_ deadlineGetter = (*p2p.PullQuery)(nil)
	_ deadlineGetter = (*p2p.AppRequest)(nil)
)

type chainIDGetter interface {
	GetChainId() []byte
}

func GetChainID(m any) (ids.ID, error) {
	msg, ok := m.(chainIDGetter)
	if !ok {
		return ids.Empty, errMissingField
	}
	chainIDBytes := msg.GetChainId()
	return ids.ToID(chainIDBytes)
}

type sourceChainIDGetter interface {
	GetSourceChainID() ids.ID
}

func GetSourceChainID(m any) (ids.ID, error) {
	msg, ok := m.(sourceChainIDGetter)
	if !ok {
		return GetChainID(m)
	}
	return msg.GetSourceChainID(), nil
}

type requestIDGetter interface {
	GetRequestId() uint32
}

func GetRequestID(m any) (uint32, bool) {
	if msg, ok := m.(requestIDGetter); ok {
		requestID := msg.GetRequestId()
		return requestID, true
	}

	// AppGossip is the only inbound message not containing a requestID. For
	// ease of handling, imagine that it does have a requestID.
	_, ok := m.(*p2p.AppGossip)
	return 0, ok
}

type engineTypeGetter interface {
	GetEngineType() p2p.EngineType
}

func GetEngineType(m any) (p2p.EngineType, bool) {
	msg, ok := m.(engineTypeGetter)
	if !ok {
		return p2p.EngineType_ENGINE_TYPE_UNSPECIFIED, false
	}
	return msg.GetEngineType(), true
}

type deadlineGetter interface {
	GetDeadline() uint64
}

func GetDeadline(m any) (time.Duration, bool) {
	msg, ok := m.(deadlineGetter)
	if !ok {
		return 0, false
	}
	deadline := msg.GetDeadline()
	return time.Duration(deadline), true
}
