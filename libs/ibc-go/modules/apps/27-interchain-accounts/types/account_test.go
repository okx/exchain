package types_test

import (
	"testing"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
	"github.com/stretchr/testify/suite"
)

var (
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "ex14r7mrj0nus8k57slulkmrdyeyp7t8xvrdzqmsz"

	// TestPortID defines a reusable port identifier for testing purposes
	TestPortID, _ = types.NewControllerPortID(TestOwnerAddress)
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
			TestOwnerAddress,
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
