// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"context"
	"crypto"
	"net"
	"net/netip"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/skychains/chain/ids"
	"github.com/skychains/chain/message"
	"github.com/skychains/chain/network/throttling"
	"github.com/skychains/chain/proto/pb/p2p"
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

type testPeer struct {
	Peer
	inboundMsgChan <-chan message.InboundMessage
}

type rawTestPeer struct {
	config         *Config
	cert           *staking.Certificate
	nodeID         ids.NodeID
	inboundMsgChan <-chan message.InboundMessage
}

func newMessageCreator(t *testing.T) message.Creator {
	t.Helper()

	mc, err := message.NewCreator(
		logging.NoLog{},
		prometheus.NewRegistry(),
		constants.DefaultNetworkCompressionType,
		10*time.Second,
	)
	require.NoError(t, err)

	return mc
}

func newConfig(t *testing.T) Config {
	t.Helper()
	require := require.New(t)

	metrics, err := NewMetrics(prometheus.NewRegistry())
	require.NoError(err)

	resourceTracker, err := tracker.NewResourceTracker(
		prometheus.NewRegistry(),
		resource.NoUsage,
		meter.ContinuousFactory{},
		10*time.Second,
	)
	require.NoError(err)

	return Config{
		ReadBufferSize:       constants.DefaultNetworkPeerReadBufferSize,
		WriteBufferSize:      constants.DefaultNetworkPeerWriteBufferSize,
		Metrics:              metrics,
		MessageCreator:       newMessageCreator(t),
		Log:                  logging.NoLog{},
		InboundMsgThrottler:  throttling.NewNoInboundThrottler(),
		Network:              TestNetwork,
		Router:               nil,
		VersionCompatibility: version.GetCompatibility(constants.LocalID),
		MySubnets:            nil,
		Beacons:              validators.NewManager(),
		Validators:           validators.NewManager(),
		NetworkID:            constants.LocalID,
		PingFrequency:        constants.DefaultPingFrequency,
		PongTimeout:          constants.DefaultPingPongTimeout,
		MaxClockDifference:   time.Minute,
		ResourceTracker:      resourceTracker,
		UptimeCalculator:     uptime.NoOpCalculator,
		IPSigner:             nil,
	}
}

func newRawTestPeer(t *testing.T, config Config) *rawTestPeer {
	t.Helper()
	require := require.New(t)

	tlsCert, err := staking.NewTLSCert()
	require.NoError(err)
	cert, err := staking.ParseCertificate(tlsCert.Leaf.Raw)
	require.NoError(err)
	nodeID := ids.NodeIDFromCert(cert)

	ip := utils.NewAtomic(netip.AddrPortFrom(
		netip.IPv6Loopback(),
		1,
	))
	tls := tlsCert.PrivateKey.(crypto.Signer)
	bls, err := bls.NewSecretKey()
	require.NoError(err)

	config.IPSigner = NewIPSigner(ip, tls, bls)

	inboundMsgChan := make(chan message.InboundMessage)
	config.Router = router.InboundHandlerFunc(func(_ context.Context, msg message.InboundMessage) {
		inboundMsgChan <- msg
	})

	return &rawTestPeer{
		config:         &config,
		cert:           cert,
		nodeID:         nodeID,
		inboundMsgChan: inboundMsgChan,
	}
}

func startTestPeer(self *rawTestPeer, peer *rawTestPeer, conn net.Conn) *testPeer {
	return &testPeer{
		Peer: Start(
			self.config,
			conn,
			peer.cert,
			peer.nodeID,
			NewThrottledMessageQueue(
				self.config.Metrics,
				peer.nodeID,
				logging.NoLog{},
				throttling.NewNoOutboundThrottler(),
			),
		),
		inboundMsgChan: self.inboundMsgChan,
	}
}

func startTestPeers(rawPeer0 *rawTestPeer, rawPeer1 *rawTestPeer) (*testPeer, *testPeer) {
	conn0, conn1 := net.Pipe()
	peer0 := startTestPeer(rawPeer0, rawPeer1, conn0)
	peer1 := startTestPeer(rawPeer1, rawPeer0, conn1)
	return peer0, peer1
}

func awaitReady(t *testing.T, peers ...Peer) {
	t.Helper()
	require := require.New(t)

	for _, peer := range peers {
		require.NoError(peer.AwaitReady(context.Background()))
		require.True(peer.Ready())
	}
}

func TestReady(t *testing.T) {
	require := require.New(t)

	config := newConfig(t)

	rawPeer0 := newRawTestPeer(t, config)
	rawPeer1 := newRawTestPeer(t, config)

	conn0, conn1 := net.Pipe()

	peer0 := startTestPeer(rawPeer0, rawPeer1, conn0)
	require.False(peer0.Ready())

	peer1 := startTestPeer(rawPeer1, rawPeer0, conn1)
	awaitReady(t, peer0, peer1)

	peer0.StartClose()
	require.NoError(peer0.AwaitClosed(context.Background()))
	require.NoError(peer1.AwaitClosed(context.Background()))
}

func TestSend(t *testing.T) {
	require := require.New(t)

	sharedConfig := newConfig(t)

	rawPeer0 := newRawTestPeer(t, sharedConfig)
	rawPeer1 := newRawTestPeer(t, sharedConfig)

	peer0, peer1 := startTestPeers(rawPeer0, rawPeer1)
	awaitReady(t, peer0, peer1)

	outboundGetMsg, err := sharedConfig.MessageCreator.Get(ids.Empty, 1, time.Second, ids.Empty)
	require.NoError(err)

	require.True(peer0.Send(context.Background(), outboundGetMsg))

	inboundGetMsg := <-peer1.inboundMsgChan
	require.Equal(message.GetOp, inboundGetMsg.Op())

	peer1.StartClose()
	require.NoError(peer0.AwaitClosed(context.Background()))
	require.NoError(peer1.AwaitClosed(context.Background()))
}

func TestPingUptimes(t *testing.T) {
	trackedSubnetID := ids.GenerateTestID()
	untrackedSubnetID := ids.GenerateTestID()

	sharedConfig := newConfig(t)
	sharedConfig.MySubnets = set.Of(trackedSubnetID)

	testCases := []struct {
		name        string
		msg         message.OutboundMessage
		shouldClose bool
		assertFn    func(*require.Assertions, *testPeer)
	}{
		{
			name: "primary network only",
			msg: func() message.OutboundMessage {
				pingMsg, err := sharedConfig.MessageCreator.Ping(1, nil)
				require.NoError(t, err)
				return pingMsg
			}(),
			shouldClose: false,
			assertFn: func(require *require.Assertions, peer *testPeer) {
				uptime, ok := peer.ObservedUptime(constants.PrimaryNetworkID)
				require.True(ok)
				require.Equal(uint32(1), uptime)

				uptime, ok = peer.ObservedUptime(trackedSubnetID)
				require.False(ok)
				require.Zero(uptime)
			},
		},
		{
			name: "primary network and subnet",
			msg: func() message.OutboundMessage {
				pingMsg, err := sharedConfig.MessageCreator.Ping(
					1,
					[]*p2p.SubnetUptime{
						{
							SubnetId: trackedSubnetID[:],
							Uptime:   1,
						},
					},
				)
				require.NoError(t, err)
				return pingMsg
			}(),
			shouldClose: false,
			assertFn: func(require *require.Assertions, peer *testPeer) {
				uptime, ok := peer.ObservedUptime(constants.PrimaryNetworkID)
				require.True(ok)
				require.Equal(uint32(1), uptime)

				uptime, ok = peer.ObservedUptime(trackedSubnetID)
				require.True(ok)
				require.Equal(uint32(1), uptime)
			},
		},
		{
			name: "primary network and non tracked subnet",
			msg: func() message.OutboundMessage {
				pingMsg, err := sharedConfig.MessageCreator.Ping(
					1,
					[]*p2p.SubnetUptime{
						{
							// Providing the untrackedSubnetID here should cause
							// the remote peer to disconnect from us.
							SubnetId: untrackedSubnetID[:],
							Uptime:   1,
						},
						{
							SubnetId: trackedSubnetID[:],
							Uptime:   1,
						},
					},
				)
				require.NoError(t, err)
				return pingMsg
			}(),
			shouldClose: true,
			assertFn:    nil,
		},
	}

	// The raw peers are generated outside of the test cases to avoid generating
	// many TLS keys.
	rawPeer0 := newRawTestPeer(t, sharedConfig)
	rawPeer1 := newRawTestPeer(t, sharedConfig)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)

			peer0, peer1 := startTestPeers(rawPeer0, rawPeer1)
			awaitReady(t, peer0, peer1)
			defer func() {
				peer1.StartClose()
				peer0.StartClose()
				require.NoError(peer0.AwaitClosed(context.Background()))
				require.NoError(peer1.AwaitClosed(context.Background()))
			}()

			require.True(peer0.Send(context.Background(), tc.msg))

			if tc.shouldClose {
				require.NoError(peer1.AwaitClosed(context.Background()))
				return
			}

			// we send Get message after ping to ensure Ping is handled by the
			// time Get is handled. This is because Get is routed to the handler
			// whereas Ping is handled by the peer directly. We have no way to
			// know when the peer has handled the Ping message.
			sendAndFlush(t, peer0, peer1)

			tc.assertFn(require, peer1)
		})
	}
}

func TestTrackedSubnets(t *testing.T) {
	sharedConfig := newConfig(t)
	rawPeer0 := newRawTestPeer(t, sharedConfig)
	rawPeer1 := newRawTestPeer(t, sharedConfig)

	makeSubnetIDs := func(numSubnets int) []ids.ID {
		subnetIDs := make([]ids.ID, numSubnets)
		for i := range subnetIDs {
			subnetIDs[i] = ids.GenerateTestID()
		}
		return subnetIDs
	}

	tests := []struct {
		name             string
		trackedSubnets   []ids.ID
		shouldDisconnect bool
	}{
		{
			name:             "primary network only",
			trackedSubnets:   makeSubnetIDs(0),
			shouldDisconnect: false,
		},
		{
			name:             "single subnet",
			trackedSubnets:   makeSubnetIDs(1),
			shouldDisconnect: false,
		},
		{
			name:             "max subnets",
			trackedSubnets:   makeSubnetIDs(maxNumTrackedSubnets),
			shouldDisconnect: false,
		},
		{
			name:             "too many subnets",
			trackedSubnets:   makeSubnetIDs(maxNumTrackedSubnets + 1),
			shouldDisconnect: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			rawPeer0.config.MySubnets = set.Of(test.trackedSubnets...)
			peer0, peer1 := startTestPeers(rawPeer0, rawPeer1)
			if test.shouldDisconnect {
				require.NoError(peer0.AwaitClosed(context.Background()))
				require.NoError(peer1.AwaitClosed(context.Background()))
				return
			}

			defer func() {
				peer1.StartClose()
				peer0.StartClose()
				require.NoError(peer0.AwaitClosed(context.Background()))
				require.NoError(peer1.AwaitClosed(context.Background()))
			}()

			awaitReady(t, peer0, peer1)

			require.Equal(set.Of(constants.PrimaryNetworkID), peer0.TrackedSubnets())

			expectedTrackedSubnets := set.Of(test.trackedSubnets...)
			expectedTrackedSubnets.Add(constants.PrimaryNetworkID)
			require.Equal(expectedTrackedSubnets, peer1.TrackedSubnets())
		})
	}
}

// Test that a peer using the wrong BLS key is disconnected from.
func TestInvalidBLSKeyDisconnects(t *testing.T) {
	require := require.New(t)

	sharedConfig := newConfig(t)

	rawPeer0 := newRawTestPeer(t, sharedConfig)
	rawPeer1 := newRawTestPeer(t, sharedConfig)

	require.NoError(rawPeer0.config.Validators.AddStaker(
		constants.PrimaryNetworkID,
		rawPeer1.nodeID,
		bls.PublicFromSecretKey(rawPeer1.config.IPSigner.blsSigner),
		ids.GenerateTestID(),
		1,
	))

	bogusBLSKey, err := bls.NewSecretKey()
	require.NoError(err)
	require.NoError(rawPeer1.config.Validators.AddStaker(
		constants.PrimaryNetworkID,
		rawPeer0.nodeID,
		bls.PublicFromSecretKey(bogusBLSKey), // This is the wrong BLS key for this peer
		ids.GenerateTestID(),
		1,
	))

	peer0, peer1 := startTestPeers(rawPeer0, rawPeer1)

	// Because peer1 thinks that peer0 is using the wrong BLS key, they should
	// disconnect from each other.
	require.NoError(peer0.AwaitClosed(context.Background()))
	require.NoError(peer1.AwaitClosed(context.Background()))
}

func TestShouldDisconnect(t *testing.T) {
	peerID := ids.GenerateTestNodeID()
	txID := ids.GenerateTestID()
	blsKey, err := bls.NewSecretKey()
	require.NoError(t, err)

	tests := []struct {
		name                     string
		initialPeer              *peer
		expectedPeer             *peer
		expectedShouldDisconnect bool
	}{
		{
			name: "peer is reporting old version",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
				},
				version: &version.Application{
					Name:  version.Client,
					Major: 0,
					Minor: 0,
					Patch: 0,
				},
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
				},
				version: &version.Application{
					Name:  version.Client,
					Major: 0,
					Minor: 0,
					Patch: 0,
				},
			},
			expectedShouldDisconnect: true,
		},
		{
			name: "peer is not a validator",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators:           validators.NewManager(),
				},
				version: version.CurrentApp,
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators:           validators.NewManager(),
				},
				version: version.CurrentApp,
			},
			expectedShouldDisconnect: false,
		},
		{
			name: "peer is a validator without a BLS key",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							nil,
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							nil,
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
			},
			expectedShouldDisconnect: false,
		},
		{
			name: "already verified peer",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:                   peerID,
				version:              version.CurrentApp,
				txIDOfVerifiedBLSKey: txID,
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:                   peerID,
				version:              version.CurrentApp,
				txIDOfVerifiedBLSKey: txID,
			},
			expectedShouldDisconnect: false,
		},
		{
			name: "peer without signature",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
				ip:      &SignedIP{},
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
				ip:      &SignedIP{},
			},
			expectedShouldDisconnect: true,
		},
		{
			name: "peer with invalid signature",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
				ip: &SignedIP{
					BLSSignature: bls.SignProofOfPossession(blsKey, []byte("wrong message")),
				},
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
				ip: &SignedIP{
					BLSSignature: bls.SignProofOfPossession(blsKey, []byte("wrong message")),
				},
			},
			expectedShouldDisconnect: true,
		},
		{
			name: "peer with valid signature",
			initialPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
				ip: &SignedIP{
					BLSSignature: bls.SignProofOfPossession(blsKey, (&UnsignedIP{}).bytes()),
				},
			},
			expectedPeer: &peer{
				Config: &Config{
					Log:                  logging.NoLog{},
					VersionCompatibility: version.GetCompatibility(constants.UnitTestID),
					Validators: func() validators.Manager {
						vdrs := validators.NewManager()
						require.NoError(t, vdrs.AddStaker(
							constants.PrimaryNetworkID,
							peerID,
							bls.PublicFromSecretKey(blsKey),
							txID,
							1,
						))
						return vdrs
					}(),
				},
				id:      peerID,
				version: version.CurrentApp,
				ip: &SignedIP{
					BLSSignature: bls.SignProofOfPossession(blsKey, (&UnsignedIP{}).bytes()),
				},
				txIDOfVerifiedBLSKey: txID,
			},
			expectedShouldDisconnect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			shouldDisconnect := test.initialPeer.shouldDisconnect()
			require.Equal(test.expectedPeer, test.initialPeer)
			require.Equal(test.expectedShouldDisconnect, shouldDisconnect)
		})
	}
}

// Helper to send a message from sender to receiver and assert that the
// receiver receives the message. This can be used to test a prior message
// was handled by the peer.
func sendAndFlush(t *testing.T, sender *testPeer, receiver *testPeer) {
	t.Helper()
	mc := newMessageCreator(t)
	outboundGetMsg, err := mc.Get(ids.Empty, 1, time.Second, ids.Empty)
	require.NoError(t, err)
	require.True(t, sender.Send(context.Background(), outboundGetMsg))
	inboundGetMsg := <-receiver.inboundMsgChan
	require.Equal(t, message.GetOp, inboundGetMsg.Op())
}
