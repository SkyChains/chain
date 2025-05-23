// Copyright (C) 2019-2024, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package antithesis

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/compose-spec/compose-go/types"
	"gopkg.in/yaml.v3"

	"github.com/skychains/chain/config"
	"github.com/skychains/chain/tests/fixture/tmpnet"
	"github.com/skychains/chain/utils/constants"
	"github.com/skychains/chain/utils/logging"
	"github.com/skychains/chain/utils/perms"
)

const bootstrapIndex = 0

// Initialize the given path with the docker-compose configuration (compose file and
// volumes) needed for an Antithesis test setup.
func GenerateComposeConfig(
	network *tmpnet.Network,
	nodeImageName string,
	workloadImageName string,
	targetPath string,
) error {
	// Generate a compose project for the specified network
	project, err := newComposeProject(network, nodeImageName, workloadImageName)
	if err != nil {
		return fmt.Errorf("failed to create compose project: %w", err)
	}

	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to convert target path to absolute path: %w", err)
	}

	if err := os.MkdirAll(absPath, perms.ReadWriteExecute); err != nil {
		return fmt.Errorf("failed to create target path %q: %w", absPath, err)
	}

	// Write the compose file
	bytes, err := yaml.Marshal(&project)
	if err != nil {
		return fmt.Errorf("failed to marshal compose project: %w", err)
	}
	composePath := filepath.Join(targetPath, "docker-compose.yml")
	if err := os.WriteFile(composePath, bytes, perms.ReadWrite); err != nil {
		return fmt.Errorf("failed to write genesis: %w", err)
	}

	// Create the volume paths
	for _, service := range project.Services {
		for _, volume := range service.Volumes {
			volumePath := filepath.Join(absPath, volume.Source)
			if err := os.MkdirAll(volumePath, perms.ReadWriteExecute); err != nil {
				return fmt.Errorf("failed to create volume path %q: %w", volumePath, err)
			}
		}
	}
	return nil
}

// Create a new docker compose project for an antithesis test setup
// for the provided network configuration.
func newComposeProject(network *tmpnet.Network, nodeImageName string, workloadImageName string) (*types.Project, error) {
	networkName := "lux-testnet"
	baseNetworkAddress := "10.0.20"

	services := make(types.Services, len(network.Nodes)+1)
	uris := make([]string, len(network.Nodes))
	var (
		bootstrapIP  string
		bootstrapIDs string
	)
	for i, node := range network.Nodes {
		address := fmt.Sprintf("%s.%d", baseNetworkAddress, 3+i)

		tlsKey, err := node.Flags.GetStringVal(config.StakingTLSKeyContentKey)
		if err != nil {
			return nil, err
		}
		tlsCert, err := node.Flags.GetStringVal(config.StakingCertContentKey)
		if err != nil {
			return nil, err
		}
		signerKey, err := node.Flags.GetStringVal(config.StakingSignerKeyContentKey)
		if err != nil {
			return nil, err
		}

		env := types.Mapping{
			config.NetworkNameKey:             constants.LocalName,
			config.LogLevelKey:                logging.Debug.String(),
			config.LogDisplayLevelKey:         logging.Trace.String(),
			config.HTTPHostKey:                "0.0.0.0",
			config.PublicIPKey:                address,
			config.StakingTLSKeyContentKey:    tlsKey,
			config.StakingCertContentKey:      tlsCert,
			config.StakingSignerKeyContentKey: signerKey,
		}

		// Apply configuration appropriate to a test network
		for k, v := range tmpnet.DefaultTestFlags() {
			switch value := v.(type) {
			case string:
				env[k] = value
			case bool:
				env[k] = strconv.FormatBool(value)
			default:
				return nil, fmt.Errorf("unable to convert unsupported type %T to string", v)
			}
		}

		serviceName := getServiceName(i)

		volumes := []types.ServiceVolumeConfig{
			{
				Type:   types.VolumeTypeBind,
				Source: fmt.Sprintf("./volumes/%s/logs", serviceName),
				Target: "/root/.node/logs",
			},
		}

		trackSubnets, err := node.Flags.GetStringVal(config.TrackSubnetsKey)
		if err != nil {
			return nil, err
		}
		if len(trackSubnets) > 0 {
			env[config.TrackSubnetsKey] = trackSubnets
			if i == bootstrapIndex {
				// DB volume for bootstrap node will need to initialized with the subnet
				volumes = append(volumes, types.ServiceVolumeConfig{
					Type:   types.VolumeTypeBind,
					Source: fmt.Sprintf("./volumes/%s/db", serviceName),
					Target: "/root/.node/db",
				})
			}
		}

		if i == 0 {
			bootstrapIP = address + ":9651"
			bootstrapIDs = node.NodeID.String()
		} else {
			env[config.BootstrapIPsKey] = bootstrapIP
			env[config.BootstrapIDsKey] = bootstrapIDs
		}

		// The env is defined with the keys and then converted to env
		// vars because only the keys are available as constants.
		env = keyMapToEnvVarMap(env)

		services[i+1] = types.ServiceConfig{
			Name:          serviceName,
			ContainerName: serviceName,
			Hostname:      serviceName,
			Image:         nodeImageName,
			Volumes:       volumes,
			Environment:   env.ToMappingWithEquals(),
			Networks: map[string]*types.ServiceNetworkConfig{
				networkName: {
					Ipv4Address: address,
				},
			},
		}

		// Collect URIs for the workload container
		uris[i] = fmt.Sprintf("http://%s:9650", address)
	}

	workloadEnv := types.Mapping{
		"AVAWL_URIS": strings.Join(uris, " "),
	}
	chainIDs := []string{}
	for _, subnet := range network.Subnets {
		for _, chain := range subnet.Chains {
			chainIDs = append(chainIDs, chain.ChainID.String())
		}
	}
	if len(chainIDs) > 0 {
		workloadEnv["AVAWL_CHAIN_IDS"] = strings.Join(chainIDs, " ")
	}

	workloadName := "workload"
	services[0] = types.ServiceConfig{
		Name:          workloadName,
		ContainerName: workloadName,
		Hostname:      workloadName,
		Image:         workloadImageName,
		Environment:   workloadEnv.ToMappingWithEquals(),
		Networks: map[string]*types.ServiceNetworkConfig{
			networkName: {
				Ipv4Address: baseNetworkAddress + ".129",
			},
		},
	}

	return &types.Project{
		Networks: types.Networks{
			networkName: types.NetworkConfig{
				Driver: "bridge",
				Ipam: types.IPAMConfig{
					Config: []*types.IPAMPool{
						{
							Subnet: baseNetworkAddress + ".0/24",
						},
					},
				},
			},
		},
		Services: services,
	}, nil
}

// Convert a mapping of lux config keys to a mapping of env vars
func keyMapToEnvVarMap(keyMap types.Mapping) types.Mapping {
	envVarMap := make(types.Mapping, len(keyMap))
	for key, val := range keyMap {
		// e.g. network-id -> AVAGO_NETWORK_ID
		envVar := strings.ToUpper(config.EnvPrefix + "_" + config.DashesToUnderscores.Replace(key))
		envVarMap[envVar] = val
	}
	return envVarMap
}

// Retrieve the service name for a node at the given index. Common to
// GenerateComposeConfig and InitDBVolumes to ensure consistency
// between db volumes configuration and volume paths.
func getServiceName(index int) string {
	baseName := "lux"
	if index == 0 {
		return baseName + "-bootstrap-node"
	}
	return fmt.Sprintf("%s-node-%d", baseName, index)
}
