// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package gkeystore

import (
	"context"

	"github.com/skychains/chain/api/keystore"
	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/rpcdb"
	"github.com/skychains/chain/vms/rpcchainvm/grpcutils"

	keystorepb "github.com/skychains/chain/proto/pb/keystore"
	rpcdbpb "github.com/skychains/chain/proto/pb/rpcdb"
)

var _ keystorepb.KeystoreServer = (*Server)(nil)

// Server is a snow.Keystore that is managed over RPC.
type Server struct {
	keystorepb.UnsafeKeystoreServer
	ks keystore.BlockchainKeystore
}

// NewServer returns a keystore connected to a remote keystore
func NewServer(ks keystore.BlockchainKeystore) *Server {
	return &Server{
		ks: ks,
	}
}

func (s *Server) GetDatabase(
	_ context.Context,
	req *keystorepb.GetDatabaseRequest,
) (*keystorepb.GetDatabaseResponse, error) {
	db, err := s.ks.GetRawDatabase(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	closer := dbCloser{Database: db}

	serverListener, err := grpcutils.NewListener()
	if err != nil {
		return nil, err
	}

	server := grpcutils.NewServer()
	closer.closer.Add(server)
	rpcdbpb.RegisterDatabaseServer(server, rpcdb.NewServer(&closer))

	// start the db server
	go grpcutils.Serve(serverListener, server)

	return &keystorepb.GetDatabaseResponse{ServerAddr: serverListener.Addr().String()}, nil
}

type dbCloser struct {
	database.Database
	closer grpcutils.ServerCloser
}

func (db *dbCloser) Close() error {
	err := db.Database.Close()
	db.closer.Stop()
	return err
}
