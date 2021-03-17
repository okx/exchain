package evm_test

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/x/evm"
	"github.com/okex/okexchain/x/evm/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
)

func (suite *EvmTestSuite) TestProposalHandler_ManageContractDeploymentWhitelistProposal() {
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
				whitelist := suite.stateDB.GetContractDeploymentWhitelist()
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
				whitelist := suite.stateDB.GetContractDeploymentWhitelist()
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
		suite.Run(tc.msg, func() {
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

func (suite *EvmTestSuite) TestProposalHandler_ManageContractBlockedListProposal() {
	proposal := types.NewManageContractBlockedListProposal(
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
			"add address into blocked list",
			func() {},
			func() {
				blockedList := suite.stateDB.GetContractBlockedList()
				suite.Require().Equal(1, len(blockedList))
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
			"delete address from blocked list",
			func() {
				proposal.IsAdded = false
				govProposal.Content = proposal
			},
			func() {
				blockedList := suite.stateDB.GetContractBlockedList()
				suite.Require().Zero(len(blockedList))
			},
			false,
		},
		{
			"delete an address not in the blocked list",
			func() {},
			func() {},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
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
