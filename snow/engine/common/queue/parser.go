// Copyright (C) 2019-2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package queue

// Parser allows parsing a job from bytes.
type Parser interface {
	Parse([]byte) (Job, error)
}
