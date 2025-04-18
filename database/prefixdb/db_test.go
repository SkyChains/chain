// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package prefixdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/database"
	"github.com/skychains/chain/database/memdb"
)

func TestInterface(t *testing.T) {
	for name, test := range database.Tests {
		t.Run(name, func(t *testing.T) {
			db := memdb.New()
			test(t, New([]byte("hello"), db))
			test(t, New([]byte("world"), db))
			test(t, New([]byte("wor"), New([]byte("ld"), db)))
			test(t, New([]byte("ld"), New([]byte("wor"), db)))
			test(t, NewNested([]byte("wor"), New([]byte("ld"), db)))
			test(t, NewNested([]byte("ld"), New([]byte("wor"), db)))
		})
	}
}

func TestPrefixLimit(t *testing.T) {
	testString := []string{"hello", "world", "a\xff", "\x01\xff\xff\xff\xff"}
	expected := []string{"hellp", "worle", "b\x00", "\x02\x00\x00\x00\x00"}
	for i, str := range testString {
		db := newDB([]byte(str), nil)
		require.Equal(t, db.dbLimit, []byte(expected[i]))
	}
}

func FuzzKeyValue(f *testing.F) {
	database.FuzzKeyValue(f, New([]byte(""), memdb.New()))
}

func FuzzNewIteratorWithPrefix(f *testing.F) {
	database.FuzzNewIteratorWithPrefix(f, New([]byte(""), memdb.New()))
}

func FuzzNewIteratorWithStartAndPrefix(f *testing.F) {
	database.FuzzNewIteratorWithStartAndPrefix(f, New([]byte(""), memdb.New()))
}

func BenchmarkInterface(b *testing.B) {
	for _, size := range database.BenchmarkSizes {
		keys, values := database.SetupBenchmark(b, size[0], size[1], size[2])
		for name, bench := range database.Benchmarks {
			b.Run(fmt.Sprintf("prefixdb_%d_pairs_%d_keys_%d_values_%s", size[0], size[1], size[2], name), func(b *testing.B) {
				db := New([]byte("hello"), memdb.New())
				bench(b, db, keys, values)
			})
		}
	}
}
