package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/x/evm/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestProposal_ManageContractDeploymentWhitelistProposal() {
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()
	proposal := types.NewManageContractDeploymentWhitelistProposal(
		"default title",
		"default description",
		types.AddressList{addr1, addr2},
		true,
	)

	minDeposit := suite.app.EvmKeeper.GetMinDeposit(suite.ctx, proposal)
	require.Equal(suite.T(), sdk.SysCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}, minDeposit)

	maxDepositPeriod := suite.app.EvmKeeper.GetMaxDepositPeriod(suite.ctx, proposal)
	require.Equal(suite.T(), time.Hour*24, maxDepositPeriod)

	votingPeriod := suite.app.EvmKeeper.GetVotingPeriod(suite.ctx, proposal)
	require.Equal(suite.T(), time.Hour*72, votingPeriod)

	testCases := []struct {
		msg     string
		prepare func()
	}{
		{
			"pass check",
			func() {},
		},
		{
			"pass check when trying to add addresses already exists in whitelist",
			func() {
				suite.stateDB.SetContractDeploymentWhitelist(types.AddressList{addr1, addr2})
			},
		},
		{
			"pass check when trying to delete addresses from whitelist",
			func() {
				proposal.IsAdded = false
			},
		},
		{
			"pass check when trying to delete addresses from whitelist which contains none of them",
			func() {
				// clear whitelist in the store
				suite.stateDB.DeleteContractDeploymentWhitelist(suite.stateDB.GetContractDeploymentWhitelist())
				suite.Require().Zero(len(suite.stateDB.GetContractDeploymentWhitelist()))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.prepare()

			msg := govtypes.NewMsgSubmitProposal(proposal, minDeposit, addr1)
			err := suite.app.EvmKeeper.CheckMsgSubmitProposal(suite.ctx, msg)
			suite.Require().NoError(err)
		})
	}
}

func (suite *KeeperTestSuite) TestProposal_ManageContractBlockedListProposal() {
	addr := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addrUnqualified := ethcmn.BytesToAddress([]byte{0x1}).Bytes()
	proposal := types.NewManageContractBlockedListProposal(
		"default title",
		"default description",
		addr,
		true,
	)

	minDeposit := suite.app.EvmKeeper.GetMinDeposit(suite.ctx, proposal)
	require.Equal(suite.T(), sdk.SysCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}, minDeposit)

	maxDepositPeriod := suite.app.EvmKeeper.GetMaxDepositPeriod(suite.ctx, proposal)
	require.Equal(suite.T(), time.Hour*24, maxDepositPeriod)

	votingPeriod := suite.app.EvmKeeper.GetVotingPeriod(suite.ctx, proposal)
	require.Equal(suite.T(), time.Hour*72, votingPeriod)

	// check submit proposal
	msg := govtypes.NewMsgSubmitProposal(proposal, minDeposit, addr)
	require.NoError(suite.T(), suite.app.EvmKeeper.CheckMsgSubmitProposal(suite.ctx, msg))

	testCases := []struct {
		msg           string
		malleate      func()
		expectedError bool
	}{
		{
			"pass check",
			func() {},
			false,
		},
		{
			"try to add an address already exists in blocked list",
			func() {
				suite.stateDB.SetContractBlockedListMember(addr)
			},
			true,
		},
		{
			"try to delete an address not in the blocked list",
			func() {
				proposal.IsAdded = false
				proposal.ContractAddr = addrUnqualified
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			tc.malleate()

			err := suite.app.EvmKeeper.CheckMsgManageContractBlockedListProposal(suite.ctx, proposal)
			if tc.expectedError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
