package types

import (
	"testing"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite
	address1 string
	address2 string
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.address1 = sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes()).String()
	suite.address2 = sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes()).String()
}

func (suite *GenesisTestSuite) TestValidateGenesis() {
	newGen := NewGenesisState(DefaultParams(), []FeeSplit{})
	testCases := []struct {
		name     string
		genState GenesisState
		expPass  bool
	}{
		{
			name:     "valid genesis constructor",
			genState: newGen,
			expPass:  true,
		},
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "valid genesis",
			genState: GenesisState{
				Params:    DefaultParams(),
				FeeSplits: []FeeSplit{},
			},
			expPass: true,
		},
		{
			name: "valid genesis - with fee",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress:   "0xdac17f958d2ee523a2206206994597c13d831ec8",
						DeployerAddress:   suite.address2,
						WithdrawerAddress: suite.address2,
					},
				},
			},
			expPass: true,
		},
		{
			name:     "empty genesis",
			genState: GenesisState{},
			expPass:  false,
		},
		{
			name: "invalid genesis - duplicated fee",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress: suite.address1,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - duplicated fee with different deployer address",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress: suite.address2,
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid genesis - invalid contract address",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress: suite.address1,
						DeployerAddress: suite.address1,
					},
				},
			},
			expPass: false,
		},
		{
			name: "valid genesis - with 0x deployer address",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
					},
				},
			},
			expPass: true,
		},
		{
			name: "valid genesis - with 0x withdraw address",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress:   "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress:   suite.address1,
						WithdrawerAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
					},
				},
			},
			expPass: true,
		},
		{
			name: "invalid genesis - invalid withdrawer address",
			genState: GenesisState{
				Params: DefaultParams(),
				FeeSplits: []FeeSplit{
					{
						ContractAddress:   "0xdac17f958d2ee523a2206206994597c13d831ec7",
						DeployerAddress:   suite.address1,
						WithdrawerAddress: "withdraw",
					},
				},
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
