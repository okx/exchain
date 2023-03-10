package keeper_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/okx/okbchain/app"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	minttypes "github.com/okx/okbchain/libs/cosmos-sdk/x/mint"
	transfertypes "github.com/okx/okbchain/libs/ibc-go/modules/apps/transfer/types"
	clienttypes "github.com/okx/okbchain/libs/ibc-go/modules/core/02-client/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	tmbytes "github.com/okx/okbchain/libs/tendermint/libs/bytes"
	"github.com/okx/okbchain/x/erc20/keeper"
	"github.com/okx/okbchain/x/erc20/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

var (
	Uint256, _ = abi.NewType("uint256", "", nil)
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.OKBChainApp

	querier sdk.Querier
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.NewContext(checkTx, abci.Header{
		Height:  1,
		ChainID: "ethermint-3",
		Time:    time.Now().UTC(),
	})
	suite.querier = keeper.NewQuerier(suite.app.Erc20Keeper)
	suite.app.Erc20Keeper.SetParams(suite.ctx, types.DefaultParams())
}

func (suite *KeeperTestSuite) MintCoins(address sdk.AccAddress, coins sdk.Coins) error {
	err := suite.app.SupplyKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	if err != nil {
		return err
	}
	err = suite.app.SupplyKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
	if err != nil {
		return err
	}
	return nil
}

func (suite *KeeperTestSuite) GetBalance(address sdk.AccAddress, denom string) sdk.Coin {
	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, address)
	return sdk.Coin{denom, acc.GetCoins().AmountOf(denom)}
}

type IbcKeeperMock struct{}

func (i IbcKeeperMock) SendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.CoinAdapter,
	sender sdk.AccAddress, receiver string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) error {
	return nil
}

func (i IbcKeeperMock) DenomPathFromHash(ctx sdk.Context, denom string) (string, error) { //nolint
	if denom == "ibc/ddcd907790b8aa2bf9b2b3b614718fa66bfc7540e832ce3e3696ea717dceff49" {
		return "transfer/channel-0", nil
	}
	if denom == "ibc/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" {
		return "transfer/channel-0", nil
	}
	return "", errors.New("not fount")
}

func (i IbcKeeperMock) GetDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) (transfertypes.DenomTrace, bool) {
	return transfertypes.DenomTrace{}, false
}

func (suite *KeeperTestSuite) TestDenomContractMap() {
	denom1 := "testdenom1"
	denom2 := "testdenom2"

	autoContract := common.BigToAddress(big.NewInt(1))
	externalContract := common.BigToAddress(big.NewInt(2))

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"success, happy path",
			func() {
				keeper := suite.app.Erc20Keeper

				contract, found := keeper.GetContractByDenom(suite.ctx, denom1)
				suite.Require().False(found)

				keeper.SetContractForDenom(suite.ctx, denom1, autoContract)

				contract, found = keeper.GetContractByDenom(suite.ctx, denom1)
				suite.Require().True(found)
				suite.Require().Equal(autoContract, contract)

				denom, found := keeper.GetDenomByContract(suite.ctx, contract)
				suite.Require().True(found)
				suite.Require().Equal(denom1, denom)

				keeper.SetContractForDenom(suite.ctx, denom1, externalContract)

				contract, found = keeper.GetContractByDenom(suite.ctx, denom1)
				suite.Require().True(found)
				suite.Require().Equal(externalContract, contract)
			},
		},
		{
			"failure, multiple denoms map to same contract",
			func() {
				keeper := suite.app.Erc20Keeper
				keeper.SetContractForDenom(suite.ctx, denom1, autoContract)
				err := keeper.SetContractForDenom(suite.ctx, denom2, autoContract)
				suite.Require().Error(err)
			},
		},
		{
			"failure, multiple denoms map to same external contract",
			func() {
				keeper := suite.app.Erc20Keeper
				err := keeper.SetContractForDenom(suite.ctx, denom1, externalContract)
				suite.Require().NoError(err)
				err = keeper.SetContractForDenom(suite.ctx, denom2, externalContract)
				suite.Require().Error(err)
			},
		},
		{
			"success, delete contract",
			func() {
				keeper := suite.app.Erc20Keeper
				r := keeper.DeleteContractForDenom(suite.ctx, denom1)
				suite.Require().Equal(r, false)
				err := keeper.SetContractForDenom(suite.ctx, denom1, externalContract)
				suite.Require().NoError(err)
				r = keeper.DeleteContractForDenom(suite.ctx, denom1)
				suite.Require().Equal(r, true)
			},
		},
		{
			"success, multiple denoms map to different contracts",
			func() {
				keeper := suite.app.Erc20Keeper
				err := keeper.SetContractForDenom(suite.ctx, denom1, autoContract)
				suite.Require().NoError(err)
				err = keeper.SetContractForDenom(suite.ctx, denom2, externalContract)
				suite.Require().NoError(err)
				out := keeper.GetContracts(suite.ctx)
				suite.Require().Equal(out[0].Contract, autoContract.String())
				suite.Require().Equal(out[1].Contract, externalContract.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}

func (suite *KeeperTestSuite) TestProxyContractRedirect() {
	denom := "testdenom1"
	addr1 := common.BigToAddress(big.NewInt(2))

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"success, proxy contract redirect owner",
			func() {
				suite.app.Erc20Keeper.InitInternalTemplateContract(suite.ctx)
				evmParams := evmtypes.DefaultParams()
				evmParams.EnableCreate = true
				evmParams.EnableCall = true
				suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, denom, addr1)
				err := suite.app.Erc20Keeper.ProxyContractRedirect(suite.ctx, denom, types.RedirectOwner, addr1)
				suite.Require().NoError(err)
			},
		},
		{
			"success, proxy contract redirect contract",
			func() {
				suite.app.Erc20Keeper.InitInternalTemplateContract(suite.ctx)
				evmParams := evmtypes.DefaultParams()
				evmParams.EnableCreate = true
				evmParams.EnableCall = true
				suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, denom, addr1)
				err := suite.app.Erc20Keeper.ProxyContractRedirect(suite.ctx, denom, types.RedirectImplementation, addr1)
				suite.Require().NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}

func (suite *KeeperTestSuite) TestSetGetTemplateContract() {
	f := func(bin string) string {
		json := `[{	"inputs": [{"internalType": "uint256","name": "a","type": "uint256"	},{	"internalType": "uint256","name": "b","type": "uint256"}],"stateMutability": "nonpayable","type": "constructor"}]`
		str := fmt.Sprintf(`{"abi":%s,"bin":"%s"}`, json, bin)
		return str
	}

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"success, default data ",
			func() {
				_, found := suite.app.Erc20Keeper.GetImplementTemplateContract(suite.ctx)
				suite.Require().Equal(true, found)
			},
		},
		{
			"success,set contract first",
			func() {
				c1 := f("c1")
				err := suite.app.Erc20Keeper.SetTemplateContract(suite.ctx, types.ProposalTypeContextTemplateImpl, c1)
				suite.Require().NoError(err)
				c11, found := suite.app.Erc20Keeper.GetImplementTemplateContract(suite.ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(found, true)
				suite.Require().NotEqual(c11, types.ModuleERC20Contract)
				suite.Require().Equal(c11.Bin, "c1")
			},
		},
		{
			"success, set contract twice",
			func() {
				c1 := f("c1")
				c2 := f("c2")
				err := suite.app.Erc20Keeper.SetTemplateContract(suite.ctx, types.ProposalTypeContextTemplateImpl, c1)
				suite.Require().NoError(err)
				err = suite.app.Erc20Keeper.SetTemplateContract(suite.ctx, types.ProposalTypeContextTemplateImpl, c2)
				suite.Require().NoError(err)
				c11, found := suite.app.Erc20Keeper.GetImplementTemplateContract(suite.ctx)
				suite.Require().Equal(found, true)
				suite.Require().NoError(err)
				suite.Require().NotEqual(c11, types.ModuleERC20Contract)
				suite.Require().NotEqual(c11.Bin, "c1")
				suite.Require().Equal(c11.Bin, "c2")
			},
		},
		{
			"success ,proxy contract",
			func() {
				_, found := suite.app.Erc20Keeper.GetProxyTemplateContract(suite.ctx)
				suite.Require().Equal(true, found)
			},
		},
		{
			"success ,set proxy contract",
			func() {
				_, found := suite.app.Erc20Keeper.GetProxyTemplateContract(suite.ctx)
				suite.Require().Equal(true, found)
				proxy := f("proxy")
				suite.app.Erc20Keeper.SetTemplateContract(suite.ctx, types.ProposalTypeContextTemplateProxy, proxy)
				cc, found := suite.app.Erc20Keeper.GetProxyTemplateContract(suite.ctx)
				suite.Require().Equal(true, found)
				suite.Require().Equal(cc.Bin, "proxy")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}
