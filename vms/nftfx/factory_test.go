// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package nftfx

import (
	"testing"
)

func TestFactory(t *testing.T) {
	factory := Factory{}
	if fx, err := factory.New(nil); err != nil {
		t.Fatal(err)
	} else if fx == nil {
		t.Fatalf("Factory.New returned nil")
	}
}
