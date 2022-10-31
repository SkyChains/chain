// Copyright (C) 2022, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// e2e implements the e2e tests.
package e2e

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/onsi/gomega"

	runner_sdk "github.com/luxdefi/lux-network-runner-sdk"

	"github.com/luxdefi/luxd/ids"
	"github.com/luxdefi/luxd/tests"
	"github.com/luxdefi/luxd/utils/crypto"
	"github.com/luxdefi/luxd/vms/secp256k1fx"
)

type ClusterType byte

const (
	Unknown ClusterType = iota
	StandAlone
	PreExisting

	// Enough for primary.NewWallet to fetch initial UTXOs.
	DefaultWalletCreationTimeout = 5 * time.Second

	// Defines default tx confirmation timeout.
	// Enough for test/custom networks.
	DefaultConfirmTxTimeout = 20 * time.Second

	DefaultShutdownTimeout = 2 * time.Minute
)

// Env is the global struct containing all we need to test
var Env = &TestEnvinronment{
	testEnvironmentConfig: &testEnvironmentConfig{
		clusterType: Unknown,
	},
}

type testEnvironmentConfig struct {
	clusterType               ClusterType
	logLevel                  string
	networkRunnerGRPCEndpoint string
	luxGoExecPath       string
	luxGoLogLevel       string
	testKeysFile              string

	// we snapshot initial state, right after starting cluster
	// to be able to reset state if needed and isolate tests
	snapshotName string
}

type TestEnvinronment struct {
	*testEnvironmentConfig

	runnerMu     sync.RWMutex
	runnerCli    runner_sdk.Client
	runnerGRPCEp string

	urisMu sync.RWMutex
	uris   []string

	testKeysMu sync.RWMutex
	testKeys   []*crypto.PrivateKeySECP256K1R

	snapMu  sync.RWMutex
	snapped bool
}

// should be called only once
// must be called before StartCluster
// Note that either networkRunnerGRPCEp or uris must be specified
func (te *TestEnvinronment) ConfigCluster(
	logLevel string,
	networkRunnerGRPCEp string,
	luxGoExecPath string,
	luxGoLogLevel string,
	uris string,
	testKeysFile string,
) error {
	if luxGoExecPath != "" {
		if _, err := os.Stat(luxGoExecPath); err != nil {
			return fmt.Errorf("could not find luxd binary: %w", err)
		}
	}

	te.testKeysFile = testKeysFile
	te.snapshotName = "ginkgo" + time.Now().String()
	switch {
	case networkRunnerGRPCEp != "" && len(uris) == 0:
		te.clusterType = StandAlone
		te.logLevel = logLevel
		te.networkRunnerGRPCEndpoint = networkRunnerGRPCEp
		te.luxGoExecPath = luxGoExecPath
		te.luxGoLogLevel = luxGoLogLevel

		err := te.setRunnerClient(te.logLevel, te.networkRunnerGRPCEndpoint)
		if err != nil {
			return fmt.Errorf("could not setup network-runner client: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		presp, err := te.GetRunnerClient().Ping(ctx)
		cancel()
		if err != nil {
			return fmt.Errorf("could not ping network-runner client: %w", err)
		}
		tests.Outf("{{green}}network-runner running in PID %d{{/}}\n", presp.Pid)

		// URIs will be set upon cluster start here
		return nil

	case networkRunnerGRPCEp == "" && len(uris) != 0:
		te.clusterType = PreExisting
		uriSlice := strings.Split(uris, ",")
		te.setURIs(uriSlice)
		tests.Outf("{{green}}URIs:{{/}} %q\n", uriSlice)
		return nil

	default:
		return fmt.Errorf("either network-runner-grpc-endpoint or uris should be specified, not both")
	}
}

func (te *TestEnvinronment) LoadKeys() error {
	// load test keys
	if len(te.testKeysFile) == 0 {
		return fmt.Errorf("test keys file not provided")
	}
	testKeys, err := tests.LoadHexTestKeys(te.testKeysFile)
	if err != nil {
		return fmt.Errorf("failed loading test keys: %w", err)
	}
	te.setTestKeys(testKeys)
	return nil
}

func (te *TestEnvinronment) StartCluster() error {
	switch te.clusterType {
	case StandAlone:
		tests.Outf("{{magenta}}starting network-runner with %q{{/}}\n", te.luxGoExecPath)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		resp, err := te.GetRunnerClient().Start(ctx, te.luxGoExecPath,
			runner_sdk.WithNumNodes(5),
			runner_sdk.WithGlobalNodeConfig(fmt.Sprintf(`{"log-level":"%s"}`, te.luxGoLogLevel)),
		)
		cancel()
		if err != nil {
			return fmt.Errorf("could not start network-runner: %w", err)
		}
		tests.Outf("{{green}}successfully started network-runner: {{/}} %+v\n", resp.ClusterInfo.NodeNames)

		// start is async, so wait some time for cluster health
		time.Sleep(time.Minute)

		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
		_, err = te.GetRunnerClient().Health(ctx)
		cancel()
		if err != nil {
			return fmt.Errorf("could not check health network-runner: %w", err)
		}

		return te.refreshURIs()

	case PreExisting:
		return nil // nothing to do, really

	default:
		return fmt.Errorf("unhandled cluster type")
	}
}

func (te *TestEnvinronment) refreshURIs() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	uriSlice, err := te.GetRunnerClient().URIs(ctx)
	cancel()
	if err != nil {
		return fmt.Errorf("could not retrieve URIs: %w", err)
	}
	te.setURIs(uriSlice)
	tests.Outf("{{green}}URIs:{{/}} %q\n", uriSlice)
	return nil
}

func (te *TestEnvinronment) setRunnerClient(logLevel string, gRPCEp string) error {
	te.runnerMu.Lock()
	defer te.runnerMu.Unlock()

	cli, err := runner_sdk.New(runner_sdk.Config{
		LogLevel:    logLevel,
		Endpoint:    gRPCEp,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	if te.runnerCli != nil {
		te.runnerCli.Close()
	}
	te.runnerCli = cli
	te.runnerGRPCEp = gRPCEp
	return err
}

func (te *TestEnvinronment) GetRunnerClient() (cli runner_sdk.Client) {
	te.runnerMu.RLock()
	cli = te.runnerCli
	te.runnerMu.RUnlock()
	return cli
}

func (te *TestEnvinronment) closeRunnerClient() (err error) {
	te.runnerMu.Lock()
	err = te.runnerCli.Close()
	te.runnerMu.Unlock()
	return err
}

func (te *TestEnvinronment) GetRunnerGRPCEndpoint() (ep string) {
	te.runnerMu.RLock()
	ep = te.runnerGRPCEp
	te.runnerMu.RUnlock()
	return ep
}

func (te *TestEnvinronment) setURIs(us []string) {
	te.urisMu.Lock()
	te.uris = us
	te.urisMu.Unlock()
}

func (te *TestEnvinronment) GetURIs() []string {
	te.urisMu.RLock()
	us := te.uris
	te.urisMu.RUnlock()
	return us
}

func (te *TestEnvinronment) setTestKeys(ks []*crypto.PrivateKeySECP256K1R) {
	te.testKeysMu.Lock()
	te.testKeys = ks
	te.testKeysMu.Unlock()
}

func (te *TestEnvinronment) GetTestKeys() ([]*crypto.PrivateKeySECP256K1R, []ids.ShortID, *secp256k1fx.Keychain) {
	te.testKeysMu.RLock()
	testKeys := te.testKeys
	te.testKeysMu.RUnlock()
	testKeyAddrs := make([]ids.ShortID, len(testKeys))
	for i := range testKeyAddrs {
		testKeyAddrs[i] = testKeys[i].PublicKey().Address()
	}
	keyChain := secp256k1fx.NewKeychain(testKeys...)
	return testKeys, testKeyAddrs, keyChain
}

func (te *TestEnvinronment) ShutdownCluster() error {
	if te.GetRunnerGRPCEndpoint() == "" {
		// we connected directly to existing cluster
		// nothing to shutdown
		return nil
	}

	runnerCli := te.GetRunnerClient()
	if runnerCli == nil {
		return fmt.Errorf("not network-runner cli")
	}

	tests.Outf("{{red}}shutting down network-runner cluster{{/}}\n")
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	_, err := runnerCli.Stop(ctx)
	cancel()
	if err != nil {
		return err
	}

	tests.Outf("{{red}}shutting down network-runner client{{/}}\n")
	return te.closeRunnerClient()
}

func (te *TestEnvinronment) SnapInitialState() error {
	te.snapMu.RLock()
	defer te.snapMu.RUnlock()

	if te.snapped {
		return nil // initial state snapshot already captured
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	_, err := te.runnerCli.SaveSnapshot(ctx, te.snapshotName)
	cancel()
	if err != nil {
		return err
	}
	te.snapped = true
	return nil
}

func (te *TestEnvinronment) RestoreInitialState(switchOffNetworkFirst bool) error {
	te.snapMu.Lock()
	defer te.snapMu.Unlock()

	if switchOffNetworkFirst {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		_, err := te.GetRunnerClient().Stop(ctx)
		cancel()
		gomega.Expect(err).Should(gomega.BeNil())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	_, err := te.GetRunnerClient().LoadSnapshot(ctx, te.snapshotName)
	cancel()
	if err != nil {
		return err
	}

	// make sure cluster goes back to health before moving on
	ctx, cancel = context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	_, err = te.GetRunnerClient().Health(ctx)
	cancel()
	if err != nil {
		return fmt.Errorf("could not check health network-runner: %w", err)
	}

	return te.refreshURIs()
}
