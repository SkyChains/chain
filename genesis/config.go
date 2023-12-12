// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/luxdefi/node/ids"
	"github.com/luxdefi/node/utils"
	"github.com/luxdefi/node/utils/constants"
	"github.com/luxdefi/node/utils/formatting/address"
	"github.com/luxdefi/node/utils/math"
	"github.com/luxdefi/node/vms/platformvm/signer"
)

var (
	_ utils.Sortable[Allocation] = Allocation{}

	errInvalidGenesisJSON = errors.New("could not unmarshal genesis JSON")
)

type LockedAmount struct {
	Amount   uint64 `json:"amount"`
	Locktime uint64 `json:"locktime"`
}

type Allocation struct {
	ETHAddr        ids.ShortID    `json:"ethAddr"`
	LUXAddr       ids.ShortID    `json:"luxAddr"`
	InitialAmount  uint64         `json:"initialAmount"`
	UnlockSchedule []LockedAmount `json:"unlockSchedule"`
}

func (a Allocation) Unparse(networkID uint32) (UnparsedAllocation, error) {
	ua := UnparsedAllocation{
		InitialAmount:  a.InitialAmount,
		UnlockSchedule: a.UnlockSchedule,
		ETHAddr:        "0x" + hex.EncodeToString(a.ETHAddr.Bytes()),
	}
	luxAddr, err := address.Format(
		"X",
		constants.GetHRP(networkID),
		a.LUXAddr.Bytes(),
	)
	ua.LUXAddr = luxAddr
	return ua, err
}

func (a Allocation) Less(other Allocation) bool {
	return a.InitialAmount < other.InitialAmount ||
		(a.InitialAmount == other.InitialAmount && a.LUXAddr.Less(other.LUXAddr))
}

type Staker struct {
	NodeID        ids.NodeID                `json:"nodeID"`
	RewardAddress ids.ShortID               `json:"rewardAddress"`
	DelegationFee uint32                    `json:"delegationFee"`
	Signer        *signer.ProofOfPossession `json:"signer,omitempty"`
}

func (s Staker) Unparse(networkID uint32) (UnparsedStaker, error) {
	luxAddr, err := address.Format(
		"X",
		constants.GetHRP(networkID),
		s.RewardAddress.Bytes(),
	)
	return UnparsedStaker{
		NodeID:        s.NodeID,
		RewardAddress: luxAddr,
		DelegationFee: s.DelegationFee,
		Signer:        s.Signer,
	}, err
}

// Config contains the genesis addresses used to construct a genesis
type Config struct {
	NetworkID uint32 `json:"networkID"`

	Allocations []Allocation `json:"allocations"`

	StartTime                  uint64        `json:"startTime"`
	InitialStakeDuration       uint64        `json:"initialStakeDuration"`
	InitialStakeDurationOffset uint64        `json:"initialStakeDurationOffset"`
	InitialStakedFunds         []ids.ShortID `json:"initialStakedFunds"`
	InitialStakers             []Staker      `json:"initialStakers"`

	CChainGenesis string `json:"cChainGenesis"`

	Message string `json:"message"`
}

func (c Config) Unparse() (UnparsedConfig, error) {
	uc := UnparsedConfig{
		NetworkID:                  c.NetworkID,
		Allocations:                make([]UnparsedAllocation, len(c.Allocations)),
		StartTime:                  c.StartTime,
		InitialStakeDuration:       c.InitialStakeDuration,
		InitialStakeDurationOffset: c.InitialStakeDurationOffset,
		InitialStakedFunds:         make([]string, len(c.InitialStakedFunds)),
		InitialStakers:             make([]UnparsedStaker, len(c.InitialStakers)),
		CChainGenesis:              c.CChainGenesis,
		Message:                    c.Message,
	}
	for i, a := range c.Allocations {
		ua, err := a.Unparse(uc.NetworkID)
		if err != nil {
			return uc, err
		}
		uc.Allocations[i] = ua
	}
	for i, isa := range c.InitialStakedFunds {
		luxAddr, err := address.Format(
			"X",
			constants.GetHRP(uc.NetworkID),
			isa.Bytes(),
		)
		if err != nil {
			return uc, err
		}
		uc.InitialStakedFunds[i] = luxAddr
	}
	for i, is := range c.InitialStakers {
		uis, err := is.Unparse(c.NetworkID)
		if err != nil {
			return uc, err
		}
		uc.InitialStakers[i] = uis
	}

	return uc, nil
}

func (c *Config) InitialSupply() (uint64, error) {
	initialSupply := uint64(0)
	for _, allocation := range c.Allocations {
		newInitialSupply, err := math.Add64(initialSupply, allocation.InitialAmount)
		if err != nil {
			return 0, err
		}
		for _, unlock := range allocation.UnlockSchedule {
			newInitialSupply, err = math.Add64(newInitialSupply, unlock.Amount)
			if err != nil {
				return 0, err
			}
		}
		initialSupply = newInitialSupply
	}
	return initialSupply, nil
}

var (
	// MainnetConfig is the config that should be used to generate the mainnet
	// genesis.
	MainnetConfig Config

	// TestNetConfig is the config that should be used to generate the testnet
	// genesis.
	TestNetConfig Config

	// LocalConfig is the config that should be used to generate a local
	// genesis.
	LocalConfig Config
)

func init() {
	unparsedMainnetConfig := UnparsedConfig{}
	unparsedTestNetConfig := UnparsedConfig{}
	unparsedLocalConfig := UnparsedConfig{}

	err := utils.Err(
		json.Unmarshal(mainnetGenesisConfigJSON, &unparsedMainnetConfig),
		json.Unmarshal(testnetGenesisConfigJSON, &unparsedTestNetConfig),
		json.Unmarshal(localGenesisConfigJSON, &unparsedLocalConfig),
	)
	if err != nil {
		panic(err)
	}

	MainnetConfig, err = unparsedMainnetConfig.Parse()
	if err != nil {
		panic(err)
	}

	TestNetConfig, err = unparsedTestNetConfig.Parse()
	if err != nil {
		panic(err)
	}

	LocalConfig, err = unparsedLocalConfig.Parse()
	if err != nil {
		panic(err)
	}
}

func GetConfig(networkID uint32) *Config {
	switch networkID {
	case constants.MainnetID:
		return &MainnetConfig
	case constants.TestNetID:
		return &TestNetConfig
	case constants.LocalID:
		return &LocalConfig
	default:
		tempConfig := LocalConfig
		tempConfig.NetworkID = networkID
		return &tempConfig
	}
}

// GetConfigFile loads a *Config from a provided filepath.
func GetConfigFile(fp string) (*Config, error) {
	bytes, err := os.ReadFile(filepath.Clean(fp))
	if err != nil {
		return nil, fmt.Errorf("unable to load file %s: %w", fp, err)
	}
	return parseGenesisJSONBytesToConfig(bytes)
}

// GetConfigContent loads a *Config from a provided environment variable
func GetConfigContent(genesisContent string) (*Config, error) {
	bytes, err := base64.StdEncoding.DecodeString(genesisContent)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64 content: %w", err)
	}
	return parseGenesisJSONBytesToConfig(bytes)
}

func parseGenesisJSONBytesToConfig(bytes []byte) (*Config, error) {
	var unparsedConfig UnparsedConfig
	if err := json.Unmarshal(bytes, &unparsedConfig); err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidGenesisJSON, err)
	}

	config, err := unparsedConfig.Parse()
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}
	return &config, nil
}
