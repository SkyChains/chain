// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"context"
	"crypto"
	"net"
	"net/netip"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/message"
	"github.com/skychains/chain/network/throttling"
	"github.com/skychains/chain/snow/networking/router"
	"github.com/skychains/chain/snow/networking/tracker"
	"github.com/skychains/chain/snow/uptime"
	"github.com/skychains/chain/snow/validators"
	"github.com/skychains/chain/staking"
	"github.com/skychains/chain/utils"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/crypto/bls"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/math/meter"
	"github.com/skychains/chain/utils/resource"
	"github.com/skychains/chain/utils/set"
	"github.com/skychains/chain/version"
)

const maxMessageToSend = 1024

// StartTestPeer provides a simple interface to create a peer that has finished
// the p2p handshake.
//
// This function will generate a new TLS key to use when connecting to the peer.
//
// The returned peer will not throttle inbound or outbound messages.
//
//   - [ctx] provides a way of canceling the connection request.
//   - [ip] is the remote that will be dialed to create the connection.
//   - [networkID] will be sent to the peer during the handshake. If the peer is
//     expecting a different [networkID], the handshake will fail and an error
//     will be returned.
//   - [router] will be called with all non-handshake messages received by the
//     peer.
func StartTestPeer(
	ctx context.Context,
	ip netip.AddrPort,
	networkID uint32,
	router router.InboundHandler,
) (Peer, error) {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, constants.NetworkType, ip.String())
	if err != nil {
		return nil, err
	}

	tlsCert, err := staking.NewTLSCert()
	if err != nil {
		return nil, err
	}

	tlsConfg := TLSConfig(*tlsCert, nil)
	clientUpgrader := NewTLSClientUpgrader(
		tlsConfg,
		prometheus.NewCounter(prometheus.CounterOpts{}),
	)

	peerID, conn, cert, err := clientUpgrader.Upgrade(conn)
	if err != nil {
		return nil, err
	}

	mc, err := message.NewCreator(
		logging.NoLog{},
		prometheus.NewRegistry(),
		constants.DefaultNetworkCompressionType,
		10*time.Second,
	)
	if err != nil {
		return nil, err
	}

	metrics, err := NewMetrics(prometheus.NewRegistry())
	if err != nil {
		return nil, err
	}

	resourceTracker, err := tracker.NewResourceTracker(
		prometheus.NewRegistry(),
		resource.NoUsage,
		meter.ContinuousFactory{},
		10*time.Second,
	)
	if err != nil {
		return nil, err
	}

	tlsKey := tlsCert.PrivateKey.(crypto.Signer)
	blsKey, err := bls.NewSecretKey()
	if err != nil {
		return nil, err
	}

	peer := Start(
		&Config{
			Metrics:              metrics,
			MessageCreator:       mc,
			Log:                  logging.NoLog{},
			InboundMsgThrottler:  throttling.NewNoInboundThrottler(),
			Network:              TestNetwork,
			Router:               router,
			VersionCompatibility: version.GetCompatibility(networkID),
			MySubnets:            set.Set[ids.ID]{},
			Beacons:              validators.NewManager(),
			Validators:           validators.NewManager(),
			NetworkID:            networkID,
			PingFrequency:        constants.DefaultPingFrequency,
			PongTimeout:          constants.DefaultPingPongTimeout,
			MaxClockDifference:   time.Minute,
			ResourceTracker:      resourceTracker,
			UptimeCalculator:     uptime.NoOpCalculator,
			IPSigner: NewIPSigner(
				utils.NewAtomic(netip.AddrPortFrom(
					netip.IPv6Loopback(),
					1,
				)),
				tlsKey,
				blsKey,
			),
		},
		conn,
		cert,
		peerID,
		NewBlockingMessageQueue(
			metrics,
			logging.NoLog{},
			maxMessageToSend,
		),
	)
	return peer, peer.AwaitReady(ctx)
}
