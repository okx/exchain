package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mock"
	"github.com/stretchr/testify/suite"
)

type FeeSplitTestSuite struct {
	suite.Suite
	address1 sdk.AccAddress
	address2 sdk.AccAddress
}

func TestFeeSplitSuite(t *testing.T) {
	suite.Run(t, new(FeeSplitTestSuite))
}

func (suite *FeeSplitTestSuite) SetupTest() {
	_, testAccounts := mock.GeneratePrivKeyAddressPairs(2)
	suite.address1 = testAccounts[0]
	suite.address2 = testAccounts[1]
}

func (suite *FeeSplitTestSuite) TestFeeNew() {
	testCases := []struct {
		name       string
		contract   common.Address
		deployer   sdk.AccAddress
		withdraw   sdk.AccAddress
		expectPass bool
	}{
		{
			"Create fee split- pass",
			ethsecp256k1.GenerateAddress(),
			suite.address1,
			suite.address2,
			true,
		},
		{
			"Create fee, omit withdraw - pass",
			ethsecp256k1.GenerateAddress(),
			suite.address1,
			nil,
			true,
		},
		{
			"Create fee split- invalid contract address",
			common.Address{},
			suite.address1,
			suite.address2,
			false,
		},
		{
			"Create fee split- invalid deployer address",
			ethsecp256k1.GenerateAddress(),
			sdk.AccAddress{},
			suite.address2,
			false,
		},
	}

	for _, tc := range testCases {
		i := NewFeeSplit(tc.contract, tc.deployer, tc.withdraw)
		err := i.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

func (suite *FeeSplitTestSuite) TestFee() {
	testCases := []struct {
		msg        string
		feeSplit   FeeSplit
		expectPass bool
	}{
		{
			"Create fee split- pass",
			FeeSplit{
				ethsecp256k1.GenerateAddress(),
				suite.address1,
				suite.address2,
			},
			true,
		},
		{
			"Create fee split- invalid contract address (Empty)",
			FeeSplit{
				common.Address{},
				suite.address1,
				suite.address2,
			},
			false,
		},
		{
			"Create fee split- invalid deployer address(Empty 1)",
			FeeSplit{
				ethsecp256k1.GenerateAddress(),
				nil,
				suite.address2,
			},
			false,
		},
		{
			"Create fee split- invalid deployer address(Empty 2)",
			FeeSplit{
				ethsecp256k1.GenerateAddress(),
				sdk.AccAddress{},
				suite.address2,
			},
			false,
		},
	}

	for _, tc := range testCases {
		err := tc.feeSplit.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.msg)
		} else {
			suite.Require().Error(err, tc.msg)
		}
	}
}

func (suite *FeeSplitTestSuite) TestFeeSplitGetters() {
	contract := ethsecp256k1.GenerateAddress()
	fs := FeeSplit{
		contract,
		suite.address1,
		suite.address2,
	}
	suite.Equal(fs.ContractAddress, contract)
	suite.Equal(fs.DeployerAddress, suite.address1)
	suite.Equal(fs.WithdrawerAddress, suite.address2)

	fs = FeeSplit{
		ContractAddress: contract,
		DeployerAddress: suite.address1,
	}
	suite.Equal(fs.ContractAddress, contract)
	suite.Equal(fs.DeployerAddress, suite.address1)
	suite.Equal(len(fs.WithdrawerAddress), 0)
}
