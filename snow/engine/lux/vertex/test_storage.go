// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/snow/consensus/lux"
)

var (
	errGet                = errors.New("unexpectedly called Get")
	errEdge               = errors.New("unexpectedly called Edge")
	errStopVertexAccepted = errors.New("unexpectedly called StopVertexAccepted")

	_ Storage = (*TestStorage)(nil)
)

type TestStorage struct {
	T                                            *testing.T
	CantGetVtx, CantEdge, CantStopVertexAccepted bool
	GetVtxF                                      func(context.Context, ids.ID) (lux.Vertex, error)
	EdgeF                                        func(context.Context) []ids.ID
	StopVertexAcceptedF                          func(context.Context) (bool, error)
}

func (s *TestStorage) Default(cant bool) {
	s.CantGetVtx = cant
	s.CantEdge = cant
}

func (s *TestStorage) GetVtx(ctx context.Context, vtxID ids.ID) (lux.Vertex, error) {
	if s.GetVtxF != nil {
		return s.GetVtxF(ctx, vtxID)
	}
	if s.CantGetVtx && s.T != nil {
		require.FailNow(s.T, errGet.Error())
	}
	return nil, errGet
}

func (s *TestStorage) Edge(ctx context.Context) []ids.ID {
	if s.EdgeF != nil {
		return s.EdgeF(ctx)
	}
	if s.CantEdge && s.T != nil {
		require.FailNow(s.T, errEdge.Error())
	}
	return nil
}

func (s *TestStorage) StopVertexAccepted(ctx context.Context) (bool, error) {
	if s.StopVertexAcceptedF != nil {
		return s.StopVertexAcceptedF(ctx)
	}
	if s.CantStopVertexAccepted && s.T != nil {
		require.FailNow(s.T, errStopVertexAccepted.Error())
	}
	return false, nil
}
