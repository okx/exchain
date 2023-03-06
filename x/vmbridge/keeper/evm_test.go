package keeper_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	keeper2 "github.com/okx/okbchain/x/vmbridge/keeper"
	"github.com/okx/okbchain/x/vmbridge/types"
	wasmtypes "github.com/okx/okbchain/x/wasm/types"

	"math/big"
)

func (suite *KeeperTestSuite) TestKeeper_SendToEvm() {

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
			types.ErrIsNotETHAddr,
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
			types.ErrIsNotETHAddr,
			true,
		},
		{
			"caller(ex wasm),contract(0x),recipient(0x),amount(0)",
			func() {
				amount = sdk.NewInt(0)
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
			"caller(ex wasm),contract(0x),recipient(0x),amount(-1)",
			func() {
				amount = sdk.NewInt(-1)
			},
			func() {

			},
			errors.New("[\"execution reverted\",\"0x4e487b710000000000000000000000000000000000000000000000000000000000000011\",\"HexData\",\"0x4e487b710000000000000000000000000000000000000000000000000000000000000011\"]"),
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
			nil,
			//errors.New("[\"execution reverted\",\"execution reverted:ERC20: mint to the zero address\",\"HexData\",\"0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001f45524332303a206d696e7420746f20746865207a65726f206164647265737300\"]"),
			true,
		},
		{
			"caller(ex wasm),contract(0x wasm),recipient(0x),amount(1)",
			func() {
				buffer := make([]byte, 32)
				buffer[31] = 0x1
				contract = "0x" + hex.EncodeToString(buffer)
			},
			func() {
			},
			errors.New("abi: attempting to unmarshall an empty string while arguments are expected"),
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
			errors.New("execution reverted"),
			true,
		},
		{
			"caller(ex wasm is no exist in erc20 contrat),contract(ex),recipient(0x),amount(1)",
			func() {
				buffer := make([]byte, 32)
				buffer[19] = 0x1
				caller = sdk.AccAddress(buffer).String()
			},
			func() {
			},
			errors.New("execution reverted"),
			true,
		},
		{
			"caller(0x wasm),contract(0x),recipient(ex),amount(1)",
			func() {
				wasmAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				temp := hex.EncodeToString(wasmAddr.Bytes())
				suite.T().Log(temp)
				caller = "0x" + temp
			},
			func() {
			},
			errors.New("execution reverted"),
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			tc.malleate()

			success, err := suite.app.VMBridgeKeeper.SendToEvm(suite.ctx, caller, contract, recipient, amount)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.success, success)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSendToWasmEventHandler_Handle() {
	contractAccAddr, err := sdk.AccAddressFromBech32("ex1fnkz39vpxmukf6mp78essh8g0hrzp3gylyd2u8")
	suite.Require().NoError(err)
	contract := common.BytesToAddress(contractAccAddr.Bytes())
	//addr := sdk.AccAddress{0x1}
	ethAddr := common.BigToAddress(big.NewInt(1))
	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"normal topic,recipient is 0x",
			func() {
				wasmAddrStr := suite.wasmContract.String()
				input, err := getSendToWasmEventData(wasmAddrStr, ethAddr.String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func() {
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.wasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			types.ErrIsNotOKCAddr,
		},
		{
			"normal topic,recipient is ex",
			func() {
				wasmAddrStr := suite.wasmContract.String()
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				input, err := getSendToWasmEventData(wasmAddrStr, queryAddr.String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
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
			"normal topic,amount is zero",
			func() {
				wasmAddrStr := suite.wasmContract.String()
				queryAddr := sdk.AccAddress(ethAddr.Bytes())
				input, err := getSendToWasmEventData(wasmAddrStr, queryAddr.String(), big.NewInt(0))
				suite.Require().NoError(err)
				data = input
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
			"error input",
			func() {
				data = []byte("ddddddd")
			},
			func() {
			},
			nil,
		},
		{
			"wasmAddStr is not wasm",
			func() {
				wasmAddrStr := sdk.AccAddress(make([]byte, 20)).String()
				input, err := getSendToWasmEventData(wasmAddrStr, sdk.AccAddress(ethAddr.Bytes()).String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func() {
			},
			types.ErrIsNotWasmAddr,
		},
		{
			"wasmAddStr is not exist",
			func() {
				wasmAddrStr := sdk.AccAddress(make([]byte, 32)).String()
				input, err := getSendToWasmEventData(wasmAddrStr, sdk.AccAddress(ethAddr.Bytes()).String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func() {
			},
			sdkerrors.Wrap(wasmtypes.ErrNotFound, "contract"),
		},
		{
			"recipient  is a error addr",
			func() {
				wasmAddrStr := suite.wasmContract.String()
				input, err := getSendToWasmEventData(wasmAddrStr, "ex111", big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func() {
			},
			errors.New("decoding bech32 failed: invalid bech32 string length 5"),
		},
		{
			"caller is not expect",
			func() {
				contract = common.BigToAddress(big.NewInt(1000))
				wasmAddrStr := suite.wasmContract.String()
				input, err := getSendToWasmEventData(wasmAddrStr, sdk.AccAddress(ethAddr.Bytes()).String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func() {
			},
			errors.New("execute wasm contract failed: The Contract addr is not expect)"),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			handler := keeper2.NewSendToWasmEventHandler(*suite.keeper)
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

func (suite *KeeperTestSuite) TestSendToWasmEvent_Unpack() {
	ethAddr := common.BigToAddress(big.NewInt(1))
	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func(wasmAddr string, recipient string, amount sdk.Int, err error)
		error     error
	}{
		{
			"normal topic",
			func() {
				wasmAddrStr := suite.wasmContract.String()
				input, err := getSendToWasmEventData(wasmAddrStr, ethAddr.String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, recipient string, amount sdk.Int, err error) {
				suite.Require().NoError(err)
				suite.Require().Equal(suite.wasmContract.String(), wasmAddr)
				suite.Require().Equal(ethAddr.String(), recipient)
				suite.Require().Equal(big.NewInt(1), amount.BigInt())
			},
			nil,
		},
		{
			"recipient is bytes",
			func() {
				testABIJson := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"recipient\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"__OKCSendToWasm\",\"type\":\"event\"}]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)

				ethAddrAcc, err := sdk.AccAddressFromBech32(ethAddr.String())
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.SendToWasmEventName].Inputs.Pack(suite.wasmContract.String(), []byte(ethAddrAcc.String()), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, recipient string, amount sdk.Int, err error) {
				suite.Require().NoError(err)
				suite.Require().NotEqual(ethAddr.String(), recipient)
				suite.Require().Equal(big.NewInt(1), amount.BigInt())
			},
			nil,
		},
		{
			"wasmAddr is bytes",
			func() {
				testABIJson := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"wasmAddr\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"__OKCSendToWasm\",\"type\":\"event\"}]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.SendToWasmEventName].Inputs.Pack([]byte(suite.wasmContract.String()), ethAddr.String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, recipient string, amount sdk.Int, err error) {
				suite.Require().NoError(err)
				suite.Require().Equal(suite.wasmContract.String(), wasmAddr)
				suite.Require().Equal(ethAddr.String(), recipient)
				suite.Require().Equal(big.NewInt(1), amount.BigInt())
			},
			nil,
		},
		{
			"event __OKCSendToWasm(string,uint256) ",
			func() {
				testABIJson := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"__OKCSendToWasm\",\"type\":\"event\"}]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.SendToWasmEventName].Inputs.Pack(ethAddr.String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, recipient string, amount sdk.Int, err error) {
			},
			errors.New("abi: cannot marshal in to go type: length insufficient 160 require 16417"),
		},
		{
			"event __OKCSendToWasm(string,string,string,uint256) ",
			func() {
				testABIJson := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient2\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"__OKCSendToWasm\",\"type\":\"event\"}]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.SendToWasmEventName].Inputs.Pack("1", "2", "3", big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, recipient string, amount sdk.Int, err error) {
				suite.Require().Equal("1", wasmAddr)
				suite.Require().Equal("2", recipient)
				suite.Require().NotEqual(big.NewInt(1), amount.BigInt())
			},
			nil,
			//errors.New("argument count mismatch: got 2 for 4"),
		},
		{
			"amount is negative",
			func() {
				testABIJson := "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"int8\",\"name\":\"amount\",\"type\":\"int8\"}],\"name\":\"__OKCSendToWasm\",\"type\":\"event\"}]\n"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.SendToWasmEventName].Inputs.Pack(suite.wasmContract.String(), ethAddr.String(), int8(-1))
				suite.T().Log(testABIEvent.Events[types.SendToWasmEventName].ID, types.SendToWasmEvent.ID)
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, recipient string, amount sdk.Int, err error) {
				suite.Require().Equal(errors.New("recover err: NewIntFromBigInt() out of bound"), err)
				suite.Require().Equal(suite.wasmContract.String(), wasmAddr)
				suite.Require().Equal(ethAddr.String(), recipient)
			},
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()
			unpacked, err := types.SendToWasmEvent.Inputs.Unpack(data)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				feild1, field2, feild3, err := getUnpack(unpacked)
				tc.postcheck(feild1, field2, feild3, err)
			}
		})
	}
}

func getUnpack(unpacked []interface{}) (wasmAddr string, recipient string, amount sdk.Int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover err: %v", r)
		}
	}()
	wasmAddr, ok := unpacked[0].(string)
	if !ok {
		return wasmAddr, recipient, amount, errors.New("the 1 feild is not string")
	}

	recipient, ok = unpacked[1].(string)
	if !ok {
		return wasmAddr, recipient, amount, errors.New("the 2 feild is not string")
	}

	temp, ok := unpacked[2].(*big.Int)
	if !ok {
		return wasmAddr, recipient, amount, errors.New("the 3 feild is not *big.Int")
	}
	amount = sdk.NewIntFromBigInt(temp)
	return
}
func getSendToWasmEventData(wasmAddr, recipient string, amount *big.Int) ([]byte, error) {
	return types.SendToWasmEvent.Inputs.Pack(wasmAddr, recipient, amount)
}
