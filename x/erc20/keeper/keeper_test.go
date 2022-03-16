package keeper_test

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	minttypes "github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/erc20/keeper"
	"github.com/stretchr/testify/suite"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.OKExChainApp

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

type IbcKeeperMock struct {
}

func (i IbcKeeperMock) SendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.CoinAdapter, sender sdk.AccAddress, receiver string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) error {
	return nil
}

func (i IbcKeeperMock) DenomPathFromHash(ctx sdk.Context, denom string) (string, error) { //nolint
	if denom == "ibc/DDCD907790B8AA2BF9B2B3B614718FA66BFC7540E832CE3E3696EA717DCEFF49" {
		return "transfer/channel-0", nil
	}
	if denom == "ibc/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" {
		return "transfer/channel-0", nil
	}
	return "", errors.New("not fount")
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

				keeper.SetAutoContractForDenom(suite.ctx, denom1, autoContract)

				contract, found = keeper.GetContractByDenom(suite.ctx, denom1)
				suite.Require().True(found)
				suite.Require().Equal(autoContract, contract)

				denom, found := keeper.GetDenomByContract(suite.ctx, contract)
				suite.Require().True(found)
				suite.Require().Equal(denom1, denom)

				keeper.SetExternalContractForDenom(suite.ctx, denom1, externalContract)

				contract, found = keeper.GetContractByDenom(suite.ctx, denom1)
				suite.Require().True(found)
				suite.Require().Equal(externalContract, contract)
			},
		},
		{
			"failure, multiple denoms map to same contract",
			func() {
				keeper := suite.app.Erc20Keeper
				keeper.SetAutoContractForDenom(suite.ctx, denom1, autoContract)
				err := keeper.SetExternalContractForDenom(suite.ctx, denom2, autoContract)
				suite.Require().Error(err)
			},
		},
		{
			"failure, multiple denoms map to same external contract",
			func() {
				keeper := suite.app.Erc20Keeper
				err := keeper.SetExternalContractForDenom(suite.ctx, denom1, externalContract)
				suite.Require().NoError(err)
				err = keeper.SetExternalContractForDenom(suite.ctx, denom2, externalContract)
				suite.Require().Error(err)
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
