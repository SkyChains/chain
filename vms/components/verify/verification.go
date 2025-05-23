// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package verify

import "github.com/skychains/chain/snow"

type Verifiable interface {
	Verify() error
}

type State interface {
	snow.ContextInitializable
	Verifiable
	IsState
}

type IsState interface {
	isState()
}

type IsNotState interface {
	isState() error
}

// All returns nil if all the verifiables were verified with no errors
func All(verifiables ...Verifiable) error {
	for _, verifiable := range verifiables {
		if err := verifiable.Verify(); err != nil {
			return err
		}
	}
	return nil
}
