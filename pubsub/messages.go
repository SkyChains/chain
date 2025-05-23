// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package pubsub

import (
	"github.com/skychains/chain/api"
	"github.com/skychains/chain/utils/formatting/address"
	"github.com/skychains/chain/utils/json"
)

// NewBloom command for a new bloom filter
//
// Deprecated: The pubsub server is deprecated.
type NewBloom struct {
	// MaxElements size of bloom filter
	MaxElements json.Uint64 `json:"maxElements"`
	// CollisionProb expected error rate of filter
	CollisionProb json.Float64 `json:"collisionProb"`
}

// NewSet command for a new map set
//
// Deprecated: The pubsub server is deprecated.
type NewSet struct{}

// AddAddresses command to add addresses
//
// Deprecated: The pubsub server is deprecated.
type AddAddresses struct {
	api.JSONAddresses

	// addressIds array of addresses, kept as a [][]byte for use in the bloom filter
	addressIds [][]byte
}

// Command execution command
//
// Deprecated: The pubsub server is deprecated.
type Command struct {
	NewBloom     *NewBloom     `json:"newBloom,omitempty"`
	NewSet       *NewSet       `json:"newSet,omitempty"`
	AddAddresses *AddAddresses `json:"addAddresses,omitempty"`
}

func (c *Command) String() string {
	switch {
	case c.NewBloom != nil:
		return "newBloom"
	case c.NewSet != nil:
		return "newSet"
	case c.AddAddresses != nil:
		return "addAddresses"
	default:
		return "unknown"
	}
}

func (c *NewBloom) IsParamsValid() bool {
	p := float64(c.CollisionProb)
	return c.MaxElements > 0 && 0 < p && p <= 1
}

// parseAddresses converts the bech32 addresses to their byte format.
func (c *AddAddresses) parseAddresses() error {
	if c.addressIds == nil {
		c.addressIds = make([][]byte, len(c.Addresses))
	}
	for i, addrStr := range c.Addresses {
		_, _, addrBytes, err := address.Parse(addrStr)
		if err != nil {
			return err
		}
		c.addressIds[i] = addrBytes
	}
	return nil
}
