package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/x/evm/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestProposal() {
	addr := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addrUnqualified := ethcmn.BytesToAddress([]byte{0x1}).Bytes()
	proposal := types.NewManageContractDeploymentWhitelistProposal(
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

	// check ManageContractDeploymentWhitelistProposal in details
	require.NoError(suite.T(), suite.app.EvmKeeper.CheckMsgManageContractDeploymentWhitelistProposal(suite.ctx, proposal))
	// error check
	// to add a address already exists in whitelist
	suite.app.EvmKeeper.SetContractDeploymentWhitelistMember(suite.ctx, addr)
	require.Error(suite.T(), suite.app.EvmKeeper.CheckMsgManageContractDeploymentWhitelistProposal(suite.ctx, proposal))
	// to delete a address not in the whitelist
	proposal.DeployerAddr = addrUnqualified
	proposal.IsAdded = false
	require.Error(suite.T(), suite.app.EvmKeeper.CheckMsgManageContractDeploymentWhitelistProposal(suite.ctx, proposal))
}
