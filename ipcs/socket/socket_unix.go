// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

//go:build !windows && !plan9 && !js
// +build !windows,!plan9,!js

package socket

import (
	"net"
	"os"
	"syscall"
	"time"

	"github.com/skychains/chain/utils/constants"
)

var staleSocketTimeout = 100 * time.Millisecond

func listen(addr string) (net.Listener, error) {
	uAddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}

	// Try to listen on the socket.
	l, err := net.ListenUnix("unix", uAddr)
	if err == nil {
		return l, nil
	}

	// Check to see if the socket is stale and remove it if it is.
	if err := removeIfStaleUnixSocket(addr); err != nil {
		return nil, err
	}

	// Try listening again now that it shouldn't be stale.
	return net.ListenUnix("unix", uAddr)
}

// Dial creates a new *Client connected to the given address over a Unix socket
func Dial(addr string) (*Client, error) {
	unixAddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		return nil, err
	}

	c, err := net.DialUnix("unix", nil, unixAddr)
	if err != nil {
		if isTimeoutError(err) {
			return nil, errReadTimeout{c.RemoteAddr()}
		}
		return nil, err
	}

	return &Client{Conn: c, maxMessageSize: int64(constants.DefaultMaxMessageSize)}, nil
}

// removeIfStaleUnixSocket takes in a path and removes it iff it is a socket
// that is refusing connections
func removeIfStaleUnixSocket(socketPath string) error {
	// Ensure it's a socket; if not return without an error
	st, err := os.Stat(socketPath)
	if err != nil {
		return nil
	}
	if st.Mode()&os.ModeType != os.ModeSocket {
		return nil
	}

	// Try to connect
	conn, err := net.DialTimeout("unix", socketPath, staleSocketTimeout)
	switch {
	// The connection was refused so this socket is stale; remove it
	case isSyscallError(err, syscall.ECONNREFUSED):
		return os.Remove(socketPath)

	// The socket is alive so close this connection and leave the socket alone
	case err == nil:
		return conn.Close()

	default:
		return nil
	}
}
