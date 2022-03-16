package keeper_test

import (
	"math/big"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/erc20/types"
	govtypes "github.com/okex/exchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestProposal_TokenMappingProposal() {
	denom1 := "testdenom1"
	externalContract := ethcmn.BigToAddress(big.NewInt(2))
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()

	proposal := types.NewTokenMappingProposal(
		"default title",
		"default description",
		denom1,
		&externalContract,
	)

	minDeposit := suite.app.Erc20Keeper.GetMinDeposit(suite.ctx, proposal)
	require.Equal(suite.T(), sdk.SysCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}, minDeposit)

	maxDepositPeriod := suite.app.Erc20Keeper.GetMaxDepositPeriod(suite.ctx, proposal)
	require.Equal(suite.T(), time.Hour*24, maxDepositPeriod)

	votingPeriod := suite.app.Erc20Keeper.GetVotingPeriod(suite.ctx, proposal)
	require.Equal(suite.T(), time.Hour*72, votingPeriod)

	testCases := []struct {
		msg     string
		prepare func()
	}{
		{
			"pass check",
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			//suite.SetupTest()
			tc.prepare()

			msg := govtypes.NewMsgSubmitProposal(proposal, minDeposit, addr1)
			err := suite.app.Erc20Keeper.CheckMsgSubmitProposal(suite.ctx, msg)
			suite.Require().NoError(err)
		})
	}
}
