package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite
	address1  sdk.AccAddress
	address2  sdk.AccAddress
	contract1 common.Address
	contract2 common.Address
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.address1 = sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	suite.address2 = sdk.AccAddress(ethsecp256k1.GenerateAddress().Bytes())
	suite.contract1 = ethsecp256k1.GenerateAddress()
	suite.contract2 = ethsecp256k1.GenerateAddress()
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
						ContractAddress: suite.contract1,
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress:   suite.contract2,
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
						ContractAddress: suite.contract1,
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress: suite.contract1,
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
						ContractAddress: suite.contract1,
						DeployerAddress: suite.address1,
					},
					{
						ContractAddress: suite.contract1,
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
						ContractAddress: common.Address{},
						DeployerAddress: suite.address1,
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
