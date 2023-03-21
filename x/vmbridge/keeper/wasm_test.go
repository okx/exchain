package keeper_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/x/vmbridge/keeper"
	"github.com/okx/okbchain/x/vmbridge/types"
	wasmtypes "github.com/okx/okbchain/x/wasm/types"
	"math/big"
)

func (suite *KeeperTestSuite) TestKeeper_SendToWasm() {
	contractAccAddr, err := sdk.AccAddressFromBech32("ex1fnkz39vpxmukf6mp78essh8g0hrzp3gylyd2u8")
	suite.Require().NoError(err)
	contract := common.BytesToAddress(contractAccAddr.Bytes())
	//addr := sdk.AccAddress{0x1}
	ethAddr := common.BigToAddress(big.NewInt(1))

	caller := sdk.AccAddress(contract.Bytes())
	wasmContractAddr := suite.wasmContract.String()
	recipient := ethAddr.String()
	amount := sdk.NewInt(1)
	reset := func() {
		caller = sdk.AccAddress(contract.Bytes())
		wasmContractAddr = suite.wasmContract.String()
		recipient = sdk.AccAddress(ethAddr.Bytes()).String()
		amount = sdk.NewInt(1)
	}
	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"normal",
			func() {
			},
			func() {
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			nil,
		},
		{
			"recipient is ex",
			func() {
				recipient = sdk.AccAddress(ethAddr.Bytes()).String()
			},
			func() {
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			nil,
		},
		{
			"recipient is 0x",
			func() {
				recipient = ethAddr.String()
			},
			func() {
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			types.ErrIsNotOKBCAddr,
		},
		{
			"recipient is wasmaddr",
			func() {
				recipient = sdk.AccAddress(make([]byte, 32)).String()
			},
			func() {
				queryAddr := sdk.AccAddress(make([]byte, 32))
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			nil,
		},
		{
			"recipient is wasmaddr 0x",
			func() {
				recipient = "0x" + hex.EncodeToString(make([]byte, 32))
			},
			func() {
				queryAddr := sdk.AccAddress(make([]byte, 32))
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			types.ErrIsNotOKBCAddr,
		},
		{
			"normal topic,amount is zero",
			func() {
				amount = sdk.NewInt(0)
			},
			func() {
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"0\"}", string(result))
			},
			nil,
		},
		{
			"wasmAddStr is not wasm",
			func() {
				wasmContractAddr = sdk.AccAddress(make([]byte, 20)).String()

			},
			func() {
			},
			types.ErrIsNotWasmAddr,
		},
		{
			"wasmAddStr is not exist",
			func() {
				wasmContractAddr = sdk.AccAddress(make([]byte, 32)).String()
			},
			func() {
			},
			sdkerrors.Wrap(wasmtypes.ErrNotFound, "contract"),
		},
		{
			"recipient  is a error addr",
			func() {
				recipient = "ex111"
			},
			func() {
			},
			errors.New("decoding bech32 failed: invalid bech32 string length 5"),
		},
		{
			"caller is not expect",
			func() {
				caller = sdk.AccAddress(common.BigToAddress(big.NewInt(1000)).Bytes())
			},
			func() {
			},
			errors.New("execute wasm contract failed: The Contract addr is not expect)"),
		},
		{
			"amount is negative",
			func() {
				amount = sdk.NewInt(-1)
			},
			func() {
			},
			types.ErrAmountNegative,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			tc.malleate()
			err := suite.keeper.SendToWasm(suite.ctx, caller, wasmContractAddr, recipient, amount)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_SendToEvmEvent() {
	caller := suite.wasmContract.String()
	contract := suite.evmContract.String()
	recipient := sdk.AccAddress(common.BigToAddress(big.NewInt(1)).Bytes()).String()
	amount := sdk.NewInt(1)

	reset := func() {
		caller = suite.wasmContract.String()
		contract = suite.evmContract.String()
		recipient = common.BigToAddress(big.NewInt(1)).String()
		amount = sdk.NewInt(1)
	}
	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
		success   bool
	}{
		{
			"caller(ex wasm),contract(0x),recipient(0x),amount(1)",
			func() {

			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				balance := suite.queryBalance(common.BytesToAddress(aimAddr.Bytes()))
				suite.Require().Equal(amount.Int64(), balance.Int64())
			},
			nil,
			true,
		},
		{
			"caller(ex wasm),contract(ex),recipient(0x),amount(1)",
			func() {
				temp, err := sdk.AccAddressFromBech32(contract)
				suite.Require().NoError(err)
				contract = temp.String()
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				balance := suite.queryBalance(common.BytesToAddress(aimAddr.Bytes()))
				suite.Require().Equal(amount.Int64(), balance.Int64())
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, types.ErrIsNotETHAddr.Error()),
			true,
		},
		{
			"caller(ex wasm),contract(0x),recipient(ex),amount(1)",
			func() {
				temp, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				recipient = temp.String()
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				balance := suite.queryBalance(common.BytesToAddress(aimAddr.Bytes()))
				suite.Require().Equal(amount.Int64(), balance.Int64())
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, types.ErrIsNotETHAddr.Error()),
			true,
		},
		{
			"caller(ex wasm),contract(0x),recipient(ex),amount(0)",
			func() {
				amount = sdk.NewInt(0)
				temp, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				recipient = temp.String()
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				balance := suite.queryBalance(common.BytesToAddress(aimAddr.Bytes()))
				suite.Require().Equal(amount.Int64(), balance.Int64())
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, types.ErrIsNotETHAddr.Error()),
			true,
		},
		{
			"caller(ex wasm),contract(0x),recipient(0x),amount(-1)",
			func() {
				amount = sdk.NewInt(-1)
			},
			func() {

			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, "[\"execution reverted\",\"0x4e487b710000000000000000000000000000000000000000000000000000000000000011\",\"HexData\",\"0x4e487b710000000000000000000000000000000000000000000000000000000000000011\"]"),
			true,
		},
		{
			"caller(ex wasm),contract(0x),recipient(0x wasm),amount(1)", // recipent is not wasm addr but is check in SendToEvmEvent Check.
			func() {
				buffer := make([]byte, 32)
				buffer[31] = 0x1

				recipient = "0x" + hex.EncodeToString(buffer)
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(recipient)
				suite.Require().NoError(err)
				balance := suite.queryBalance(common.BytesToAddress(aimAddr.Bytes()))
				suite.Require().Equal(amount.Int64(), balance.Int64())
			},
			nil, // This case is checkout in msg.validateBasic. so this case pass
			//errors.New("[\"execution reverted\",\"execution reverted:ERC20: mint to the zero address\",\"HexData\",\"0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001f45524332303a206d696e7420746f20746865207a65726f206164647265737300\"]"),
			true,
		},
		{
			"caller(ex wasm),contract(ex wasm),recipient(0x),amount(1)",
			func() {
				buffer := make([]byte, 32)
				buffer[31] = 0x1
				contract = "0x" + hex.EncodeToString(buffer)
			},
			func() {
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, "abi: attempting to unmarshall an empty string while arguments are expected"),
			true,
		},
		{
			"caller(ex nowasm),contract(0x),recipient(0x),amount(1)",
			func() {
				buffer := make([]byte, 20)
				buffer[19] = 0x1
				caller = sdk.AccAddress(buffer).String()
			},
			func() {
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, "execution reverted"),
			true,
		},
		{
			"caller(ex wasm is no exist in erc20 contrat),contract(0x),recipient(0x),amount(1)",
			func() {
				buffer := make([]byte, 32)
				buffer[19] = 0x1
				caller = sdk.AccAddress(buffer).String()
			},
			func() {
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, "execution reverted"),
			true,
		},
		{
			"caller(0x wasm),contract(0x),recipient(0x),amount(1)",
			func() {
				wasmAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				temp := hex.EncodeToString(wasmAddr.Bytes())
				suite.T().Log(temp)
				caller = "0x" + temp
			},
			func() {
			},
			sdkerrors.Wrap(types.ErrEvmExecuteFailed, "execution reverted"),
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			tc.malleate()
			msgServer := keeper.NewMsgServerImpl(*suite.app.VMBridgeKeeper)

			msg := types.MsgSendToEvm{Sender: caller, Contract: contract, Recipient: recipient, Amount: amount}
			success, err := msgServer.SendToEvmEvent(sdk.WrapSDKContext(suite.ctx), &msg)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.success, success.Success)
				tc.postcheck()
			}
		})
	}

}
