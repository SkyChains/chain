// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build windows
// +build windows

package socket

import (
	"net"

	"github.com/Microsoft/go-winio"

	"github.com/skychains/chain/utils/constants"
)

// listen creates a net.Listen backed by a Windows named pipe
func listen(addr string) (net.Listener, error) {
	return winio.ListenPipe(windowsPipeName(addr), nil)
}

// Dial creates a new *Client connected to a Windows named pipe
func Dial(addr string) (*Client, error) {
	c, err := winio.DialPipe(windowsPipeName(addr), nil)
	if err != nil {
		return nil, err
	}
	return &Client{Conn: c, maxMessageSize: int64(constants.DefaultMaxMessageSize)}, nil
}

// windowsPipeName turns an address into a valid Windows named pipes name
func windowsPipeName(addr string) string {
	return `\\.\pipe\` + addr
}
