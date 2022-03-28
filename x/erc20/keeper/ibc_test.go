package keeper_test

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	erc20Keeper "github.com/okex/exchain/x/erc20/keeper"
	"github.com/okex/exchain/x/erc20/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

func (suite *KeeperTestSuite) TestConvertVouchers() {
	addr1 := common.BigToAddress(big.NewInt(1))
	addr1Bech := sdk.AccAddress(addr1.Bytes())

	amount := int64(123)
	amountDec := sdk.NewDec(amount)

	testCases := []struct {
		msg       string
		from      string
		vouchers  sdk.SysCoins
		malleate  func()
		postcheck func()
		expError  error
	}{
		{
			"Wrong from address",
			"test",
			sdk.NewCoins(sdk.NewCoin(types.IbcDenomDefaultValue, sdk.NewInt(1))),
			func() {},
			func() {},
			errors.New("encoding/hex: invalid byte: U+0074 't'"),
		},
		{
			"Empty address",
			"",
			sdk.NewCoins(sdk.NewCoin(types.IbcDenomDefaultValue, sdk.NewInt(1))),
			func() {},
			func() {},
			errors.New("empty from address string is not allowed"),
		},
		{
			"Correct address with non supported coin denom",
			addr1Bech.String(),
			sdk.NewCoins(sdk.NewCoin("fake", sdk.NewInt(1))),
			func() {},
			func() {},
			errors.New("coin fake is not supported for wrapping"),
		},
		{
			"Correct address with not enough IBC evm token",
			addr1Bech.String(),
			sdk.NewCoins(sdk.NewCoin(types.IbcDenomDefaultValue, sdk.NewInt(123))),
			func() {
				coin := sdk.NewCoin(types.IbcDenomDefaultValue, amountDec.Sub(sdk.NewDec(3)))
				err := suite.MintCoins(addr1Bech, sdk.NewCoins(coin))
				suite.Require().NoError(err)
			},
			func() {},
			fmt.Errorf("insufficient funds: insufficient account funds; 120.000000000000000000%s < 123.000000000000000000%s",
				types.IbcDenomDefaultValue, types.IbcDenomDefaultValue),
		},
		{
			"Correct address with enough IBC-evm token : Should receive evm tokens",
			addr1Bech.String(),
			sdk.NewDecCoinsFromDec(types.IbcDenomDefaultValue, amountDec),
			func() {
				coin := sdk.NewCoin(types.IbcDenomDefaultValue, amountDec)
				err := suite.MintCoins(addr1Bech, sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(addr1Bech, types.IbcDenomDefaultValue)
				suite.Require().Equal(coin, balance)
			},
			func() {
				coin := sdk.NewCoin(sdk.DefaultBondDenom, amountDec)

				balance := suite.GetBalance(addr1Bech, sdk.DefaultBondDenom)
				suite.Require().Equal(coin, balance)
			},
			nil,
		},
		{
			"Correct address with not enough IBC token",
			addr1Bech.String(),
			sdk.NewCoins(sdk.NewCoin(CorrectIbcDenom, amountDec)),
			func() {
				coin := sdk.NewCoin(CorrectIbcDenom, amountDec.Sub(sdk.NewDec(3)))
				err := suite.MintCoins(addr1Bech, sdk.NewCoins(coin))
				suite.Require().NoError(err)

				evmParams := evmtypes.DefaultParams()
				evmParams.EnableCreate = true
				evmParams.EnableCall = true
				suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
			},
			func() {},
			fmt.Errorf("insufficient funds: insufficient account funds; 120.000000000000000000%s < 123.000000000000000000%s",
				CorrectIbcDenom, CorrectIbcDenom),
		},
		{
			"Correct address with IBC token : Should receive ERC20 tokens",
			addr1Bech.String(),
			sdk.NewDecCoinsFromDec(CorrectIbcDenom, amountDec),
			func() {
				coin := sdk.NewCoin(CorrectIbcDenom, amountDec)
				err := suite.MintCoins(addr1Bech, sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.GetBalance(addr1Bech, CorrectIbcDenom)
				suite.Require().Equal(coin, balance)

				evmParams := evmtypes.DefaultParams()
				evmParams.EnableCreate = true
				evmParams.EnableCall = true
				suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
			},
			func() {
				coin := sdk.NewCoin(CorrectIbcDenom, amountDec)

				// 1. Verify balance IBC coin post operation
				balance := suite.GetBalance(addr1Bech, CorrectIbcDenom)
				suite.Require().Equal(sdk.NewDec(0), balance.Amount)
				// 2. Verify ERC20 contract be created
				contract, found := suite.app.Erc20Keeper.GetContractByDenom(suite.ctx, CorrectIbcDenom)
				suite.Require().True(found)
				// 3. Verify balance IBC coin for contract address
				balance = suite.GetBalance(sdk.AccAddress(contract.Bytes()), CorrectIbcDenom)
				suite.Require().Equal(coin, balance)
				// 4. Verify ERC20 balance post operation
				ret, err := suite.app.Erc20Keeper.CallModuleERC20(suite.ctx, contract, "balanceOf", common.BytesToAddress(addr1Bech.Bytes()))
				suite.Require().NoError(err)
				suite.Require().Equal(big.NewInt(amount), big.NewInt(0).SetBytes(ret))
			},
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			tc.malleate()

			err := suite.app.Erc20Keeper.ConvertVouchers(suite.ctx, tc.from, tc.vouchers)
			if tc.expError != nil {
				suite.Require().EqualError(err, tc.expError.Error(), tc.msg)
			} else {
				suite.Require().NoError(err, tc.msg)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestIbcTransferVouchers() {
	addr1 := common.BigToAddress(big.NewInt(1))
	addr1Bech := sdk.AccAddress(addr1.Bytes())

	testCases := []struct {
		name      string
		from      string
		to        string
		coin      sdk.Coins
		malleate  func()
		expError  error
		postCheck func()
	}{
		{
			"Wrong from address",
			"test",
			"to",
			sdk.NewCoins(sdk.NewCoin(CorrectIbcDenom, sdk.NewInt(1))),
			func() {},
			errors.New("encoding/hex: invalid byte: U+0074 't'"),
			func() {},
		},
		{
			"Empty address",
			"",
			"to",
			sdk.NewCoins(sdk.NewCoin(CorrectIbcDenom, sdk.NewInt(1))),
			func() {},
			errors.New("empty from address string is not allowed"),
			func() {},
		},
		{
			"Correct address with non supported coin denom",
			addr1Bech.String(),
			"to",
			sdk.NewCoins(sdk.NewCoin("fake", sdk.NewInt(1))),
			func() {},
			errors.New("coin fake is not supported"),
			func() {},
		},
		//{
		//	"Correct address with too small amount EVM token",
		//	addr1Bech.String(),
		//	"to",
		//	sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(123))),
		//	func() {},
		//	nil,
		//	func() {},
		//},
		//{
		//	"Correct address with not enough EVM token",
		//	addr1Bech.String(),
		//	"to",
		//	sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1230000000000))),
		//	func() {},
		//	errors.New("0aphoton is smaller than 1230000000000aphoton: insufficient funds"),
		//	func() {},
		//},
		//{
		//	"Correct address with enough EVM token : Should receive IBC evm token",
		//	addr1Bech.String(),
		//	"to",
		//	sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1230000000000))),
		//	func() {
		//		// Mint Coin to user and module
		//		suite.MintCoins(addr1Bech, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1230000000000))))
		//		suite.MintCoinsToModule(types.ModuleName, sdk.NewCoins(sdk.NewCoin(types.IbcDenomDefaultValue, sdk.NewInt(123))))
		//		// Verify balance IBC coin pre operation
		//		ibcCoin := suite.GetBalance(addr1Bech, types.IbcDenomDefaultValue)
		//		suite.Require().Equal(sdk.NewInt(0), ibcCoin.Amount)
		//		// Verify balance EVM coin pre operation
		//		evmCoin := suite.GetBalance(addr1Bech, sdk.DefaultBondDenom)
		//		suite.Require().Equal(sdk.NewInt(1230000000000), evmCoin.Amount)
		//	},
		//	nil,
		//	func() {
		//		// Verify balance IBC coin post operation
		//		ibcCoin := suite.GetBalance(addr1Bech, types.IbcDenomDefaultValue)
		//		suite.Require().Equal(sdk.NewInt(123), ibcCoin.Amount)
		//		// Verify balance EVM coin post operation
		//		evmCoin := suite.GetBalance(addr1Bech, sdk.DefaultBondDenom)
		//		suite.Require().Equal(sdk.NewInt(0), evmCoin.Amount)
		//	},
		//},
		{
			"Correct address with non correct IBC token denom",
			addr1Bech.String(),
			"to",
			sdk.NewCoins(sdk.NewCoin("incorrect", sdk.NewInt(123))),
			func() {
				// Add support for the IBC token
				suite.app.Erc20Keeper.SetAutoContractForDenom(suite.ctx, "incorrect", common.HexToAddress("0x11"))
			},
			errors.New("ibc denom is invalid: incorrect is invalid"),
			func() {
			},
		},
		{
			"Correct address with correct IBC token denom",
			addr1Bech.String(),
			"to",
			sdk.NewCoins(sdk.NewCoin(CorrectIbcDenom, sdk.NewInt(123))),
			func() {
				// Mint IBC token for user
				suite.MintCoins(addr1Bech, sdk.NewCoins(sdk.NewCoin(CorrectIbcDenom, sdk.NewInt(123))))
				// Add support for the IBC token
				suite.app.Erc20Keeper.SetAutoContractForDenom(suite.ctx, CorrectIbcDenom, common.HexToAddress("0x11"))
			},
			nil,
			func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			// Create erc20 Keeper with mock transfer keeper
			erc20Keeper := erc20Keeper.NewKeeper(
				suite.app.Codec(),
				suite.app.GetKey(types.StoreKey),
				suite.app.GetSubspace(types.ModuleName),
				suite.app.AccountKeeper,
				suite.app.SupplyKeeper,
				suite.app.BankKeeper,
				suite.app.EvmKeeper,
				IbcKeeperMock{},
			)
			suite.app.Erc20Keeper = erc20Keeper

			tc.malleate()
			err := suite.app.Erc20Keeper.IbcTransferVouchers(suite.ctx, tc.from, tc.to, tc.coin)
			if tc.expError != nil {
				suite.Require().EqualError(err, tc.expError.Error())
			} else {
				suite.Require().NoError(err)
				tc.postCheck()
			}
		})
	}
}
