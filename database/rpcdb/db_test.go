// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package rpcdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/corruptabledb"
	"github.com/skychains/chain/database/memdb"
	"github.com/skychains/chain/vms/rpcchainvm/grpcutils"

	rpcdbpb "github.com/skychains/chain/proto/pb/rpcdb"
)

type testDatabase struct {
	client *DatabaseClient
	server *memdb.Database
}

func setupDB(t testing.TB) *testDatabase {
	require := require.New(t)

	db := &testDatabase{
		server: memdb.New(),
	}

	listener, err := grpcutils.NewListener()
	require.NoError(err)
	serverCloser := grpcutils.ServerCloser{}

	server := grpcutils.NewServer()
	rpcdbpb.RegisterDatabaseServer(server, NewServer(db.server))
	serverCloser.Add(server)

	go grpcutils.Serve(listener, server)

	conn, err := grpcutils.Dial(listener.Addr().String())
	require.NoError(err)

	db.client = NewClient(rpcdbpb.NewDatabaseClient(conn))

	t.Cleanup(func() {
		serverCloser.Stop()
		_ = conn.Close()
		_ = listener.Close()
	})

	return db
}

func TestInterface(t *testing.T) {
	for name, test := range database.Tests {
		t.Run(name, func(t *testing.T) {
			db := setupDB(t)
			test(t, db.client)
		})
	}
}

func FuzzKeyValue(f *testing.F) {
	db := setupDB(f)
	database.FuzzKeyValue(f, db.client)
}

func FuzzNewIteratorWithPrefix(f *testing.F) {
	db := setupDB(f)
	database.FuzzNewIteratorWithPrefix(f, db.client)
}

func FuzzNewIteratorWithStartAndPrefix(f *testing.F) {
	db := setupDB(f)
	database.FuzzNewIteratorWithStartAndPrefix(f, db.client)
}

func BenchmarkInterface(b *testing.B) {
	for _, size := range database.BenchmarkSizes {
		keys, values := database.SetupBenchmark(b, size[0], size[1], size[2])
		for name, bench := range database.Benchmarks {
			b.Run(fmt.Sprintf("rpcdb_%d_pairs_%d_keys_%d_values_%s", size[0], size[1], size[2], name), func(b *testing.B) {
				db := setupDB(b)
				bench(b, db.client, keys, values)
			})
		}
	}
}

func TestHealthCheck(t *testing.T) {
	scenarios := []struct {
		name         string
		testDatabase *testDatabase
		testFn       func(db *corruptabledb.Database) error
		wantErr      bool
		wantErrMsg   string
	}{
		{
			name:         "healthcheck success",
			testDatabase: setupDB(t),
			testFn: func(_ *corruptabledb.Database) error {
				return nil
			},
		},
		{
			name:         "healthcheck failed db closed",
			testDatabase: setupDB(t),
			testFn: func(db *corruptabledb.Database) error {
				return db.Close()
			},
			wantErr:    true,
			wantErrMsg: "closed",
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			require := require.New(t)

			baseDB := setupDB(t)
			db := corruptabledb.New(baseDB.server)
			defer db.Close()
			require.NoError(scenario.testFn(db))

			// check db HealthCheck
			_, err := db.HealthCheck(context.Background())
			if scenario.wantErr {
				require.Error(err) //nolint:forbidigo
				require.Contains(err.Error(), scenario.wantErrMsg)
				return
			}
			require.NoError(err)

			// check rpc HealthCheck
			_, err = baseDB.client.HealthCheck(context.Background())
			require.NoError(err)
		})
	}
}
