// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package formatting

import "fmt"

// PrefixedStringer extends a stringer that adds a prefix
type PrefixedStringer interface {
	fmt.Stringer

	PrefixedString(prefix string) string
}
