package keeper_test

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestWhitelist() {
	whitelist := suite.app.EvmKeeper.GetContractDeploymentWhitelist(suite.ctx)
	require.Zero(suite.T(), len(whitelist))

	// create addresses for test
	ethAddr := ethcmn.BytesToAddress([]byte{0x0})
	addr := ethAddr.Bytes()

	addrUnqualified := ethcmn.BytesToAddress([]byte{0x1}).Bytes()

	// setter
	suite.app.EvmKeeper.SetContractDeploymentWhitelistMember(suite.ctx, addr)
	whitelist = suite.app.EvmKeeper.GetContractDeploymentWhitelist(suite.ctx)
	require.Equal(suite.T(), 1, len(whitelist))

	// check for whitelist
	require.True(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addr, nil))
	require.True(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addr, &ethAddr))
	require.False(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addrUnqualified, nil))
	require.True(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addrUnqualified, &ethAddr))

	// delete
	suite.app.EvmKeeper.DeleteContractDeploymentWhitelistMember(suite.ctx, addr)

	// check for whitelist
	whitelist = suite.app.EvmKeeper.GetContractDeploymentWhitelist(suite.ctx)
	require.Zero(suite.T(), len(whitelist))

	require.False(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addr, nil))
	require.True(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addr, &ethAddr))
	require.False(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addrUnqualified, nil))
	require.True(suite.T(), suite.app.EvmKeeper.IsContractDeployerQualified(suite.ctx, addrUnqualified, &ethAddr))
}
