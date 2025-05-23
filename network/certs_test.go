// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"crypto/tls"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/network/peer"
	"github.com/skychains/chain/staking"
)

var (
	certLock   sync.Mutex
	tlsCerts   []*tls.Certificate
	tlsConfigs []*tls.Config
)

func getTLS(t *testing.T, index int) (ids.NodeID, *tls.Certificate, *tls.Config) {
	certLock.Lock()
	defer certLock.Unlock()

	for len(tlsCerts) <= index {
		cert, err := staking.NewTLSCert()
		require.NoError(t, err)
		tlsConfig := peer.TLSConfig(*cert, nil)

		tlsCerts = append(tlsCerts, cert)
		tlsConfigs = append(tlsConfigs, tlsConfig)
	}

	tlsCert := tlsCerts[index]
	cert := staking.CertificateFromX509(tlsCert.Leaf)
	nodeID := ids.NodeIDFromCert(cert)
	return nodeID, tlsCert, tlsConfigs[index]
}
