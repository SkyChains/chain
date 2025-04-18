// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build test

package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var _ Timer = (*TimerTest)(nil)

// TimerTest is a test timer
type TimerTest struct {
	T *testing.T

	CantRegisterTimout bool

	RegisterTimeoutF func(time.Duration)
}

// Default set the default callable value to [cant]
func (t *TimerTest) Default(cant bool) {
	t.CantRegisterTimout = cant
}

func (t *TimerTest) RegisterTimeout(delay time.Duration) {
	if t.RegisterTimeoutF != nil {
		t.RegisterTimeoutF(delay)
	} else if t.CantRegisterTimout && t.T != nil {
		require.FailNow(t.T, "Unexpectedly called RegisterTimeout")
	}
}
