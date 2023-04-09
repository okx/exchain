package keeper_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	keeper2 "github.com/okex/exchain/x/vmbridge/keeper"
	"github.com/okex/exchain/x/vmbridge/types"
	wasmtypes "github.com/okex/exchain/x/wasm/types"
	"github.com/stretchr/testify/require"
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
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
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
				queryAddr := sdk.WasmAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
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
				queryAddr := sdk.WasmAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
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
			"wasmAddStr is not exist",
			func() {
				wasmAddrStr := sdk.AccAddress(make([]byte, 20)).String()
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

func (suite *KeeperTestSuite) TestKeeper_CallToEvm() {

	caller := suite.freeCallWasmContract.String()
	callerhex := common.BytesToAddress(suite.freeCallWasmContract.Bytes()).String()
	contract := suite.freeCallEvmContract.String()
	contractEx := sdk.AccAddress(suite.freeCallEvmContract.Bytes()).String()
	callDataFormat := "{\"call_to_evm\":{\"value\":\"0\",\"evmaddr\":\"%s\",\"calldata\":\"%s\"}}"
	callData := fmt.Sprintf(callDataFormat, contract, "init-to-call-evm")
	value := sdk.NewInt(0)
	evmReturnPrefix := "callByWasm return: %s ---data: "
	reset := func() {
		caller = suite.freeCallWasmContract.String()
		contract = suite.freeCallEvmContract.String()
		callData = fmt.Sprintf(callDataFormat, contract, "init-to-call-evm")
		value = sdk.NewInt(0)
	}
	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
		expect    string
	}{
		{
			"caller(ex 32),contract(0x 20),calldata(normal),amount(0)",
			func() {

			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, caller) + fmt.Sprintf(callDataFormat, contract, "init-to-call-evm"),
		},
		{
			"caller(0x 32),contract(0x 20),calldata(normal),amount(0)",
			func() {
				caller = callerhex
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, callerhex) + fmt.Sprintf(callDataFormat, contract, "init-to-call-evm"),
		},
		{
			"caller(ex 20),contract(0x 20),calldata(normal),amount(0)",
			func() {
				buffer := make([]byte, 20)
				buffer[19] = 0x1
				caller = sdk.AccAddress(buffer).String()
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, "ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpxuz0nc") + fmt.Sprintf(callDataFormat, contract, "init-to-call-evm"),
		},
		{
			"caller(0x 20),contract(0x 20),calldata(normal),amount(0)",
			func() {
				buffer := make([]byte, 20)
				buffer[19] = 0x1
				caller = common.BytesToAddress(buffer).String()
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, "0x0000000000000000000000000000000000000001") + fmt.Sprintf(callDataFormat, contract, "init-to-call-evm"),
		},
		{
			"caller(ex 32),contract(0x 32),calldata(normal),amount(0)",
			func() {
				buffer := make([]byte, 32)
				buffer[19] = 0x1
				contract = common.BytesToAddress(buffer).String()
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			errors.New("abi: attempting to unmarshall an empty string while arguments are expected"),
			"",
		},
		{
			"caller(ex 32),contract(ex),calldata(normal),amount(0)",
			func() {
				contract = contractEx
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			types.ErrIsNotETHAddr,
			"",
		},
		{
			"caller(ex 32),contract(0x 20),calldata(emppty),amount(0)",
			func() {
				callData = fmt.Sprintf(callDataFormat, contract, "")
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, caller) + fmt.Sprintf(callDataFormat, contract, ""),
		},
		{
			"caller(ex 32),contract(0x 20),calldata(normal),amount(1)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(1).BigInt())
				suite.SetAccountCoins(suite.freeCallWasmContract, sdk.NewInt(1))
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())

				aimAddr, err = sdk.AccAddressFromBech32(contract)
				suite.Require().NoError(err)
				balance = suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}.String(), balance.String())
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, caller) + fmt.Sprintf(callDataFormat, contract, "init-to-call-evm"),
		},
		{
			"caller(ex 32),contract(0x 20),calldata(normal),amount(2)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(2).BigInt())
				suite.SetAccountCoins(suite.freeCallWasmContract, sdk.NewInt(1))
			},
			func() {
				aimAddr, err := sdk.AccAddressFromBech32(caller)
				suite.Require().NoError(err)
				balance := suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())

				aimAddr, err = sdk.AccAddressFromBech32(contract)
				suite.Require().NoError(err)
				balance = suite.queryCoins(aimAddr)
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}.String(), balance.String())
			},
			errors.New("insufficient funds: insufficient account funds; 1.000000000000000000okt < 2.000000000000000000okt"),
			"",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			tc.malleate()

			response, err := suite.app.VMBridgeKeeper.CallToEvm(suite.ctx, caller, contract, callData, value)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expect, response)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCallToWasmEventHandler_Handle() {
	tempAddr, err := sdk.AccAddressFromBech32("ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq")
	suite.Require().NoError(err)

	caller := suite.freeCallEvmContract

	wasmContractAddr := suite.freeCallWasmContract.String()
	calldata := "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
	value := sdk.NewInt(0)
	data, err := getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
	require.NoError(suite.T(), err)

	reset := func() {
		caller = suite.freeCallEvmContract
		wasmContractAddr = suite.freeCallWasmContract.String()
		calldata = "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
		value = sdk.NewInt(0)
		data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
		require.NoError(suite.T(), err)
	}
	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"caller(exist),wasmContract(ex 32),value(0),data(normal)",
			func() {
			},
			func() {
				queryAddr := sdk.AccAddress(caller.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = tempAddr
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))
			},
			nil,
		},
		{
			"caller(no exist),wasmContract(ex 32),value(0),data(normal)",
			func() {
				caller = common.BytesToAddress(make([]byte, 20))
			},
			func() {
			},
			errors.New("execute wasm contract failed: Insufficient funds (balance 0, required=100)"),
		},
		{
			"caller(exist),wasmContract(0x 32),value(0),data(normal)",
			func() {
				data, err = getCallToWasmEventData(hex.EncodeToString(suite.freeCallWasmContract), value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
				queryAddr := sdk.AccAddress(caller.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = tempAddr
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))
			},
			nil,
		},
		{
			"caller(exist),wasmContract(ex not found),value(0),data(normal)",
			func() {
				data, err = getCallToWasmEventData(sdk.AccAddress(make([]byte, 32)).String(), value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {

			},
			errors.New("not found: contract"),
		},
		{
			"caller(exist),wasmContract(ex 20),value(0),data(normal)",
			func() {

				data, err = getCallToWasmEventData(sdk.AccAddress(make([]byte, 20)).String(), value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {

			},
			types.ErrIsNotWasmAddr,
		},
		{
			"caller(exist),wasmContract(0x 20),value(0),data(normal)",
			func() {

				data, err = getCallToWasmEventData(hex.EncodeToString(make([]byte, 20)), value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {

			},
			types.ErrIsNotWasmAddr,
		},
		{
			"caller(exist),wasmContract(ex 32),value(1),data(normal)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(1).BigInt())
				suite.SetAccountCoins(sdk.AccAddress(caller.Bytes()), sdk.NewInt(1))
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
				queryAddr := sdk.AccAddress(caller.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = tempAddr
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))

				balance := suite.queryCoins(caller.Bytes())
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())

				balance = suite.queryCoins(suite.freeCallWasmContract)
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}.String(), balance.String())
			},
			nil,
		},
		{
			"caller(exist),wasmContract(ex 32),value(2),data(normal)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(2).BigInt())
				suite.SetAccountCoins(sdk.AccAddress(caller.Bytes()), sdk.NewInt(1))
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
			},
			errors.New("insufficient funds: insufficient account funds; 1.000000000000000000okt < 2.000000000000000000okt"),
		},
		{
			"caller(exist),wasmContract(ex 32),value(-1),data(normal)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(-1).BigInt())
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
			},
			nil, //because it has been recover check err
		},
		{
			"caller(exist),wasmContract(ex 32),value(1),data(error msg)",
			func() {
				calldata := "11111111122222222"
				value := sdk.NewInt(0)
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
			},
			errors.New("json: cannot unmarshal number into Go value of type map[string]interface {}"),
		},
		{
			"caller(exist),wasmContract(ex 32),value(-1),data(empty msg)",
			func() {
				calldata := ""
				value := sdk.NewInt(0)
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)

			},
			func() {
			},
			errors.New("unexpected end of JSON input"),
		},
		{
			"caller(exist),wasmContract(ex 32),value(-1),data(nofound method msg)",
			func() {
				calldata := "{\"transfer1\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
				value := sdk.NewInt(0)
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
			},
			errors.New("execute wasm contract failed: Error parsing into type cw_erc20::msg::ExecuteMsg: unknown variant `transfer1`, expected one of `approve`, `transfer`, `transfer_from`, `burn`, `mint_c_w20`, `call_to_evm`"),
		},
		{
			"caller(exist),wasmContract(ex 32),value(-1),data(multi method msg)",
			func() {
				calldata := "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"},\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
				value := sdk.NewInt(0)
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
			},
			errors.New("execute wasm contract failed: Error parsing into type cw_erc20::msg::ExecuteMsg: Expected this character to start a JSON value."),
		},
		{
			"caller(exist),wasmContract(ex 32),value(-1),data(other method msg)",
			func() {
				calldata := "{\"transfer_from\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
				value := sdk.NewInt(0)
				data, err = getCallToWasmEventData(wasmContractAddr, value.BigInt(), hex.EncodeToString([]byte(calldata)))
				require.NoError(suite.T(), err)
			},
			func() {
			},
			errors.New("execute wasm contract failed: Error parsing into type cw_erc20::msg::ExecuteMsg: missing field `owner`"),
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			handler := keeper2.NewCallToWasmEventHandler(*suite.keeper)
			tc.malleate()

			if tc.msg == "caller(exist),wasmContract(ex 32),value(-1),data(normal)" {
				defer func() {
					r := recover()
					suite.Require().NotNil(r)
					suite.Require().Equal(r.(string), "NewIntFromBigInt() out of bound")
				}()
			}
			err := handler.Handle(suite.ctx, caller, data)

			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCallToWasmEvent_Unpack() {
	normalCallData := "test calldata"
	var data []byte

	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func(wasmAddr string, value sdk.Int, calldata string, err error)
		error     error
	}{
		{
			"normal topic",
			func() {
				wasmAddrStr := suite.freeCallWasmContract.String()
				input, err := getCallToWasmEventData(wasmAddrStr, big.NewInt(1), hex.EncodeToString([]byte(normalCallData)))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, value sdk.Int, calldata string, err error) {
				suite.Require().NoError(err)
				suite.Require().Equal(suite.freeCallWasmContract.String(), wasmAddr)
				suite.Require().Equal(normalCallData, calldata)
				suite.Require().Equal(big.NewInt(1), value.BigInt())
			},
			nil,
		},
		{
			"calldata is bytes",
			func() {
				testABIJson := "[{\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"value\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes\",\n        \"name\": \"calldata\",\n        \"type\": \"bytes\"\n      }\n    ],\n    \"name\": \"__OKCCallToWasm\",\n    \"type\": \"event\"\n  }]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)

				input, err := testABIEvent.Events[types.CallToWasmEventName].Inputs.Pack(suite.freeCallWasmContract.String(), big.NewInt(1), []byte(hex.EncodeToString([]byte(normalCallData))))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, value sdk.Int, calldata string, err error) {
				suite.Require().NoError(err)
				suite.Require().Equal(normalCallData, calldata)
				suite.Require().Equal(big.NewInt(1), value.BigInt())
			},
			nil,
		},
		{
			"wasmAddr is bytes",
			func() {
				testABIJson := "[{\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"bytes\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"value\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"calldata\",\n        \"type\": \"string\"\n      }\n    ],\n    \"name\": \"__OKCCallToWasm\",\n    \"type\": \"event\"\n  }]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.CallToWasmEventName].Inputs.Pack([]byte(suite.freeCallWasmContract.String()), big.NewInt(1), hex.EncodeToString([]byte(normalCallData)))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, value sdk.Int, calldata string, err error) {
				suite.Require().NoError(err)
				suite.Require().Equal(suite.freeCallWasmContract.String(), wasmAddr)
				suite.Require().Equal(normalCallData, calldata)
				suite.Require().Equal(big.NewInt(1), value.BigInt())
			},
			nil,
		},
		{
			"event __OKCCallToWasm(string,uint256) ",
			func() {
				testABIJson := "[{\n    \"anonymous\":false,\n    \"inputs\":[\n        {\n            \"indexed\":false,\n            \"internalType\":\"string\",\n            \"name\":\"wasmAddr\",\n            \"type\":\"string\"\n        },\n        {\n            \"indexed\":false,\n            \"internalType\":\"uint256\",\n            \"name\":\"value\",\n            \"type\":\"uint256\"\n        }\n    ],\n    \"name\":\"__OKCCallToWasm\",\n    \"type\":\"event\"\n}]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.CallToWasmEventName].Inputs.Pack(suite.freeCallWasmContract.String(), big.NewInt(1))
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, value sdk.Int, calldata string, err error) {
			},
			errors.New("abi: length larger than int64: 6901746346790563787434755862277025452451108972170386555162524223799389"),
		},
		{
			"event __OKCCallToWasm(string,uint256,string,string) ",
			func() {
				testABIJson := "[{\n    \"anonymous\":false,\n    \"inputs\":[\n        {\n            \"indexed\":false,\n            \"internalType\":\"string\",\n            \"name\":\"wasmAddr\",\n            \"type\":\"string\"\n        },\n        {\n            \"indexed\":false,\n            \"internalType\":\"uint256\",\n            \"name\":\"value\",\n            \"type\":\"uint256\"\n        },\n        {\n            \"indexed\":false,\n            \"internalType\":\"string\",\n            \"name\":\"calldata\",\n            \"type\":\"string\"\n        },\n        {\n            \"indexed\":false,\n            \"internalType\":\"string\",\n            \"name\":\"calldata1\",\n            \"type\":\"string\"\n        }\n    ],\n    \"name\":\"__OKCCallToWasm\",\n    \"type\":\"event\"\n}]"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.CallToWasmEventName].Inputs.Pack("1", big.NewInt(1), hex.EncodeToString([]byte(normalCallData)), "3")
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, value sdk.Int, calldata string, err error) {
				suite.Require().Equal("1", wasmAddr)
				suite.Require().Equal(normalCallData, calldata)
				suite.Require().Equal(big.NewInt(1), value.BigInt())
			},
			nil,
			//errors.New("argument count mismatch: got 2 for 4"),
		},
		{
			"value is negative",
			func() {
				testABIJson := "[{\n    \"anonymous\":false,\n    \"inputs\":[\n        {\n            \"indexed\":false,\n            \"internalType\":\"string\",\n            \"name\":\"wasmAddr\",\n            \"type\":\"string\"\n        },\n        {\n            \"indexed\":false,\n            \"internalType\":\"int8\",\n            \"name\":\"value\",\n            \"type\":\"int8\"\n        },\n        {\n            \"indexed\":false,\n            \"internalType\":\"string\",\n            \"name\":\"calldata\",\n            \"type\":\"string\"\n        }\n    ],\n    \"name\":\"__OKCCallToWasm\",\n    \"type\":\"event\"\n}]\n"

				testABIEvent, err := abi.JSON(bytes.NewReader([]byte(testABIJson)))
				suite.Require().NoError(err)
				input, err := testABIEvent.Events[types.CallToWasmEventName].Inputs.Pack(suite.wasmContract.String(), int8(-1), hex.EncodeToString([]byte(normalCallData)))
				suite.T().Log(testABIEvent.Events[types.CallToWasmEventName].ID, types.CallToWasmEvent.ID)
				suite.Require().NoError(err)
				data = input
			},
			func(wasmAddr string, value sdk.Int, calldata string, err error) {
				suite.Require().Equal(errors.New("recover err: NewIntFromBigInt() out of bound"), err)
				suite.Require().Equal(suite.wasmContract.String(), wasmAddr)
				suite.Require().Equal("", calldata)
			},
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()
			unpacked, err := types.CallToWasmEvent.Inputs.Unpack(data)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				feild1, field2, feild3, err := getCallToWasmUnpack(unpacked)
				tc.postcheck(feild1, field2, feild3, err)
			}
		})
	}
}

func getCallToWasmEventData(wasmAddr string, value *big.Int, calldata string) ([]byte, error) {
	return types.CallToWasmEvent.Inputs.Pack(wasmAddr, value, calldata)
}

func getCallToWasmUnpack(unpacked []interface{}) (wasmAddr string, value sdk.Int, calldata string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover err: %v", r)
		}
	}()
	wasmAddr, ok := unpacked[0].(string)
	if !ok {
		return wasmAddr, value, calldata, errors.New("the 1 feild is not string")
	}

	temp, ok := unpacked[1].(*big.Int)
	if !ok {
		return wasmAddr, value, calldata, errors.New("the 2 feild is not *big.Int")
	}
	value = sdk.NewIntFromBigInt(temp)

	temp1, ok := unpacked[2].(string)
	if !ok {
		return wasmAddr, value, calldata, errors.New("the 3 feild is not string")
	}

	buff, err := hex.DecodeString(temp1)
	if err != nil {
		return wasmAddr, value, calldata, errors.New("the 3 feild must be hex")
	}
	calldata = string(buff)
	return
}
