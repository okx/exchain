package types_test

import (
	govtypes "github.com/okex/exchain/libs/cosmos-sdk/x/gov/types"
	// "github.com/cosmos/cosmos-sdk/codec"
	// codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	// govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	ibctesting "github.com/okex/exchain/libs/ibc-go/testing"
)

func (suite *TypesTestSuite) testValidateBasic() {
	subjectPath := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(subjectPath)
	subject := subjectPath.EndpointA.ClientID

	substitutePath := ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupClients(substitutePath)
	substitute := substitutePath.EndpointA.ClientID

	testCases := []struct {
		name     string
		proposal govtypes.Content
		expPass  bool
	}{
		{
			"success",
			types.NewClientUpdateProposal(ibctesting.Title, ibctesting.Description, subject, substitute),
			true,
		},
		{
			"fails validate abstract - empty title",
			types.NewClientUpdateProposal("", ibctesting.Description, subject, substitute),
			false,
		},
		{
			"subject and substitute use the same identifier",
			types.NewClientUpdateProposal(ibctesting.Title, ibctesting.Description, subject, subject),
			false,
		},
		{
			"invalid subject clientID",
			types.NewClientUpdateProposal(ibctesting.Title, ibctesting.Description, ibctesting.InvalidID, substitute),
			false,
		},
		{
			"invalid substitute clientID",
			types.NewClientUpdateProposal(ibctesting.Title, ibctesting.Description, subject, ibctesting.InvalidID),
			false,
		},
	}

	for _, tc := range testCases {

		err := tc.proposal.ValidateBasic()

		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
