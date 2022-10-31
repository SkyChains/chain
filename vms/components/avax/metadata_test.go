// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package lux

import (
	"testing"
)

func TestMetaDataVerifyNil(t *testing.T) {
	md := (*Metadata)(nil)
	if err := md.Verify(); err == nil {
		t.Fatalf("Should have errored due to nil metadata")
	}
}

func TestMetaDataVerifyUninitialized(t *testing.T) {
	md := &Metadata{}
	if err := md.Verify(); err == nil {
		t.Fatalf("Should have errored due to uninitialized metadata")
	}
}
