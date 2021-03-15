package evm_test

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/x/evm"
	"github.com/okex/okexchain/x/evm/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
)

func (suite *EvmTestSuite) TestProposalHandler() {
	proposal := types.NewManageContractDeploymentWhitelistProposal(
		"default title",
		"default description",
		ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
		true,
	)

	suite.govHandler = evm.NewManageContractDeploymentWhitelistProposalHandler(suite.app.EvmKeeper)
	govProposal := govtypes.Proposal{
		Content: proposal,
	}

	testCases := []struct {
		msg           string
		malleate      func()
		statusCheck   func()
		expectedError bool
	}{
		{
			"add address into whitelist",
			func() {},
			func() {
				whitelist := suite.app.EvmKeeper.GetContractDeploymentWhitelist(suite.ctx)
				suite.Require().Equal(1, len(whitelist))
			},
			false,
		},
		{
			"add address repeatedly",
			func() {},
			func() {},
			true,
		},
		{
			"delete address from whitelist",
			func() {
				proposal.IsAdded = false
				govProposal.Content = proposal
			},
			func() {
				whitelist := suite.app.EvmKeeper.GetContractDeploymentWhitelist(suite.ctx)
				suite.Require().Zero(len(whitelist))
			},
			false,
		},
		{
			"delete an address not in the whitelist",
			func() {},
			func() {},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run("", func() {
			tc.malleate()

			err := suite.govHandler(suite.ctx, &govProposal)

			tc.statusCheck()

			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
