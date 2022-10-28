package keeper_test

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	types2 "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	"github.com/okex/exchain/x/erc20/keeper"
)

const CorrectIbcDenom = "ibc/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

func (suite *KeeperTestSuite) TestSendToIbcHandler() {
	contract := common.BigToAddress(big.NewInt(1))
	sender := common.BigToAddress(big.NewInt(2))
	invalidDenom := "testdenom"
	validDenom := CorrectIbcDenom
	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"non associated coin denom, expect fail",
			func() {
				coin := sdk.NewCoin(invalidDenom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), invalidDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				data = input
			},
			func() {},
			errors.New("contract 0x0000000000000000000000000000000000000001 is not connected to native token"),
		},
		{
			"non IBC denom, expect fail",
			func() {
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, invalidDenom, contract)
				coin := sdk.NewCoin(invalidDenom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), invalidDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				data = input
			},
			func() {},
			errors.New("the native token associated with the contract 0x0000000000000000000000000000000000000001 is not an ibc voucher"),
		},
		{
			"success send to ibc",
			func() {
				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					amount,
				)
				data = input
			},
			func() {},
			nil,
		},
		{
			"denomination trace not found",
			func() {
				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)
				suite.ctx.SetIsCheckTx(false)
				suite.ctx.SetIsTraceTx(false)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				data = input
			},
			func() {},
			errors.New("denomination trace not found: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			handler := keeper.NewSendToIbcEventHandler(suite.app.Erc20Keeper)
			tc.malleate()
			err := handler.Handle(suite.ctx, contract, data)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSendNative20ToIbcHandler() {
	contract := common.BigToAddress(big.NewInt(1))
	sender := common.BigToAddress(big.NewInt(2))
	validDenom := "testdenom"

	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"non associated coin denom, expect fail",
			func() {
				coin := sdk.NewCoin(validDenom, sdk.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				data = input
			},
			func() {},
			errors.New("contract 0x0000000000000000000000000000000000000001 is not connected to native token"),
		},
		{
			"success send to ibc",
			func() {
				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)

				input, err := keeper.SendToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					amount,
				)
				data = input
			},
			func() {},
			nil,
		},
		{
			"portid channel error",
			func() {

				amount := sdk.NewInt(100)
				suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, validDenom, contract)
				coin := sdk.NewCoin(validDenom, amount)
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)
				suite.ctx.SetIsCheckTx(false)
				suite.ctx.SetIsTraceTx(false)
				balance := suite.GetBalance(sdk.AccAddress(contract.Bytes()), validDenom)
				suite.Require().Equal(coin, balance)
				suite.app.TransferKeeper.SetParams(suite.ctx, types2.Params{true, true})
				input, err := keeper.SendNative20ToIbcEvent.Inputs.Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
					"transfer",
					"channel-0",
				)
				data = input
			},
			func() {},
			errors.New("channel not found: port ID (transfer) channel ID (channel-0)"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			handler := keeper.NewSendNative20ToIbcEventHandler(suite.app.Erc20Keeper)
			tc.malleate()
			err := handler.Handle(suite.ctx, contract, data)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}
