package types_test

import (
	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
	"github.com/stretchr/testify/suite"
)

var (
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress, _ = sdk.AccAddressFromBech32("ex14r7mrj0nus8k57slulkmrdyeyp7t8xvrdzqmsz")

	// TestPortID defines a reusable port identifier for testing purposes
	TestPortID, _ = types.NewControllerPortID(TestOwnerAddress.String())
)

type TypesTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	chainA ibctesting.TestChainI
	chainB ibctesting.TestChainI
}

func (suite *TypesTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)

	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(0))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(1))
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

// TODO
//func (suite *TypesTestSuite) TestGenerateAddress() {
//	addr := types.GenerateAddress(k.accountKeeper.GetModuleAddress(types.ModuleName),"test-connection-id", "test-port-id")
//	accAddr, err := sdk.AccAddressFromBech32(addr.String())
//
//	suite.Require().NoError(err, "TestGenerateAddress failed")
//	suite.Require().NotEmpty(accAddr)
//}

func (suite *TypesTestSuite) TestValidateAccountAddress() {
	testCases := []struct {
		name    string
		address string
		expPass bool
	}{
		{
			"success",
			TestOwnerAddress.String(),
			true,
		},
		{
			"success with single character",
			"a",
			true,
		},
		{
			"empty string",
			"",
			false,
		},
		{
			"only spaces",
			"     ",
			false,
		},
		{
			"address is too long",
			ibctesting.LongString,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := types.ValidateAccountAddress(tc.address)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

// TODO,修改为amino account
func (suite *TypesTestSuite) TestInterchainAccount() {
}

//
//func (suite *TypesTestSuite) TestGenesisAccountValidate() {
//	pubkey := secp256k1.GenPrivKey().PubKey()
//	addr := sdk.AccAddress(pubkey.Address())
//	baseAcc := authtypes.NewBaseAccountWithAddress(addr)
//	pubkey = secp256k1.GenPrivKey().PubKey()
//	ownerAddr := sdk.AccAddress(pubkey.Address())
//
//	testCases := []struct {
//		name    string
//		acc     authtypes.GenesisAccount
//		expPass bool
//	}{
//		{
//			"success",
//			types.NewInterchainAccount(baseAcc, ownerAddr.String()),
//			true,
//		},
//		{
//			"interchain account with empty AccountOwner field",
//			types.NewInterchainAccount(baseAcc, ""),
//			false,
//		},
//	}
//
//	for _, tc := range testCases {
//		err := tc.acc.Validate()
//
//		if tc.expPass {
//			suite.Require().NoError(err)
//		} else {
//			suite.Require().Error(err)
//		}
//	}
//}
//
//func (suite *TypesTestSuite) TestInterchainAccountMarshalYAML() {
//	addr := suite.chainA.SenderAccount.GetAddress()
//	baseAcc := authtypes.NewBaseAccountWithAddress(addr)
//
//	interchainAcc := types.NewInterchainAccount(baseAcc, suite.chainB.SenderAccount.GetAddress().String())
//	bz, err := interchainAcc.MarshalYAML()
//	suite.Require().NoError(err)
//
//	expected := fmt.Sprintf("address: %s\npublic_key: \"\"\naccount_number: 0\nsequence: 0\naccount_owner: %s\n", suite.chainA.SenderAccount.GetAddress(), suite.chainB.SenderAccount.GetAddress())
//	suite.Require().Equal(expected, string(bz))
//}
//
//func (suite *TypesTestSuite) TestInterchainAccountJSON() {
//	addr := suite.chainA.SenderAccount.GetAddress()
//	ba := authtypes.NewBaseAccountWithAddress(addr)
//
//	interchainAcc := types.NewInterchainAccount(ba, suite.chainB.SenderAccount.GetAddress().String())
//
//	bz, err := json.Marshal(interchainAcc)
//	suite.Require().NoError(err)
//
//	bz1, err := interchainAcc.MarshalJSON()
//	suite.Require().NoError(err)
//	suite.Require().Equal(string(bz), string(bz1))
//
//	var a types.InterchainAccount
//	suite.Require().NoError(json.Unmarshal(bz, &a))
//	suite.Require().Equal(a.String(), interchainAcc.String())
//}
