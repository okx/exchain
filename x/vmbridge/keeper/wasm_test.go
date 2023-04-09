package keeper_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/vmbridge/keeper"
	"github.com/okex/exchain/x/vmbridge/types"
	wasmtypes "github.com/okex/exchain/x/wasm/types"
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
				queryAddr := sdk.WasmAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
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
				queryAddr := sdk.WasmAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
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
				queryAddr := sdk.WasmAddress(ethAddr.Bytes())
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			types.ErrIsNotOKCAddr,
		},
		{
			"recipient is wasmaddr",
			func() {
				recipient = sdk.AccAddress(make([]byte, 20)).String()
			},
			func() {
				queryAddr := sdk.WasmAddress(make([]byte, 32))
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
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
				queryAddr := sdk.WasmAddress(make([]byte, 32))
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, sdk.AccToAWasmddress(suite.wasmContract), []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"1\"}", string(result))
			},
			types.ErrIsNotOKCAddr,
		},
		{
			"normal topic,amount is zero",
			func() {
				amount = sdk.NewInt(0)
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
			"wasmAddStr is not exist",
			func() {
				wasmContractAddr = sdk.AccAddress(make([]byte, 20)).String()
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

func (suite *KeeperTestSuite) TestKeeper_CallToWasm() {
	//addr := sdk.AccAddress{0x1}
	tempAddr, err := sdk.AccAddressFromBech32("ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq")
	suite.Require().NoError(err)

	caller := sdk.AccAddress(suite.freeCallEvmContract.Bytes())
	wasmContractAddr := suite.freeCallWasmContract.String()
	calldata := "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"
	value := sdk.NewInt(0)
	reset := func() {
		caller = sdk.AccAddress(suite.freeCallEvmContract.Bytes())
		wasmContractAddr = suite.freeCallWasmContract.String()
		calldata = "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"0xCf164e001d86639231d92Ab1D71DB8353E43C295\"}}"
		value = sdk.NewInt(0)
	}
	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
	}{
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(normal)",
			func() {
			},
			func() {
				queryAddr := sdk.AccToAWasmddress(caller)
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = sdk.AccToAWasmddress(tempAddr)
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))
			},
			nil,
		},
		{
			"caller(20),wasmContract(ex 20),value(0),calldata(normal)",
			func() {
				wasmContractAddr = sdk.WasmToAccAddress(suite.freeCallWasmContract).String()
			},
			func() {
				queryAddr := sdk.AccToAWasmddress(caller)
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = sdk.AccToAWasmddress(tempAddr)
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))
			},
			nil,
		},
		{
			"caller(not engouh cw20 token),wasmContract(32),value(0),calldata(normal)",
			func() {
				buffer := make([]byte, 20)
				buffer[19] = 0x1
				caller = sdk.AccAddress(buffer)
			},
			func() {

			},
			errors.New("execute wasm contract failed: Insufficient funds (balance 0, required=100)"),
		},
		{
			"caller(20),wasmContract(0x 32),value(0),calldata(normal)",
			func() {
				wasmContractAddr = hex.EncodeToString(make([]byte, 32))
			},
			func() {
			},
			errors.New("incorrect address length"),
		},
		{
			"caller(20),wasmContract(no exist),value(0),calldata(normal)",
			func() {
				wasmContractAddr = sdk.WasmAddress(make([]byte, 20)).String()
			},
			func() {
			},
			errors.New("not found: contract"),
		},
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(empty)",
			func() {
				calldata = ""
			},
			func() {

			},
			errors.New("unexpected end of JSON input"),
		},
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(ex wasm addr)",
			func() {
				calldata = "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
			},
			func() {

			},
			errors.New("execute wasm contract failed: Generic error: addr_validate errored: Address is not normalized"),
		},
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(ex wasm addr)",
			func() {
				calldata = "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"0xCf164e001d86639231d92Ab1D71DB8353E43c295\"}}"
			},
			func() {

			},
			errors.New("execute wasm contract failed: Generic error: addr_validate errored: Address is not normalized"),
		},
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(not standard schema)",
			func() {
				calldata = "{\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"},\"transfer\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
			},
			func() {

			},
			errors.New("execute wasm contract failed: Error parsing into type cw_erc20::msg::ExecuteMsg: Expected this character to start a JSON value."),
		},
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(not method schema)",
			func() {
				calldata = "{\"transfer1\":{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
			},
			func() {

			},
			errors.New("execute wasm contract failed: Error parsing into type cw_erc20::msg::ExecuteMsg: unknown variant `transfer1`, expected one of `approve`, `transfer`, `transfer_from`, `burn`, `mint_c_w20`, `call_to_evm`"),
		},
		{
			"caller(20),wasmContract(0x 20),value(0),calldata(not method schema)",
			func() {
				calldata = "{\"transfer1:{\"amount\":\"100\",\"recipient\":\"ex1eutyuqqase3eyvwe92caw8dcx5ly8s544q3hmq\"}}"
			},
			func() {

			},
			errors.New("invalid character 'a' after object key"),
		},
		{
			"caller(20),wasmContract(ex 20),value(1),calldata(normal)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(1).BigInt())
				suite.SetAccountCoins(caller, sdk.NewInt(1))
			},
			func() {
				queryAddr := sdk.AccToAWasmddress(caller)
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = sdk.AccToAWasmddress(tempAddr)
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))

				balance := suite.queryCoins(caller)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())

				balance = suite.queryCoins(sdk.WasmToAccAddress(suite.freeCallWasmContract))
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}.String(), balance.String())
			},
			nil,
		},
		{
			"caller(20),wasmContract(ex 20),value(2 insufficient),calldata(normal)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(2).BigInt())
				suite.SetAccountCoins(caller, sdk.NewInt(1))
			},
			func() {
				queryAddr := caller
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = tempAddr
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))

				balance := suite.queryCoins(caller)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())

				balance = suite.queryCoins(sdk.WasmToAccAddress(suite.freeCallWasmContract))
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}.String(), balance.String())
			},
			errors.New("insufficient funds: insufficient account funds; 1.000000000000000000okt < 2.000000000000000000okt"),
		},
		{
			"caller(20),wasmContract(ex 20),value(-1 negative),calldata(normal)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(-1).BigInt())
				suite.SetAccountCoins(caller, sdk.NewInt(1))
			},
			func() {
				queryAddr := caller
				result, err := suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"99999900\"}", string(result))

				queryAddr = tempAddr
				result, err = suite.app.WasmKeeper.QuerySmart(suite.ctx, suite.freeCallWasmContract, []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", queryAddr.String())))
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100\"}", string(result))

				balance := suite.queryCoins(caller)
				suite.Require().Equal(sdk.Coins{}.String(), balance.String())

				balance = suite.queryCoins(sdk.WasmToAccAddress(suite.freeCallWasmContract))
				suite.Require().Equal(sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)}.String(), balance.String())
			},
			types.ErrAmountNegative,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			tc.malleate()
			err := suite.keeper.CallToWasm(suite.ctx, caller, wasmContractAddr, value, calldata)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				tc.postcheck()
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServer_CallToEvmEvent() {
	caller := suite.freeCallWasmContract.String()
	contract := suite.freeCallEvmContract.String()
	contractEx := sdk.AccAddress(suite.freeCallEvmContract.Bytes()).String()
	callDataFormat := "{\"call_to_evm\":{\"value\":\"0\",\"evmaddr\":\"%s\",\"calldata\":\"%s\"}}"
	callData := fmt.Sprintf(callDataFormat, contract, "init-to-call-evm")
	value := sdk.NewInt(0)
	evmReturnPrefix := "callByWasm return: %s ---data: "
	evmInput, err := getCallByWasmInput(suite.evmABI, caller, callData)
	suite.Require().NoError(err)

	reset := func() {
		caller = suite.freeCallWasmContract.String()
		contract = suite.freeCallEvmContract.String()
		callData = fmt.Sprintf(callDataFormat, contract, "init-to-call-evm")
		value = sdk.NewInt(0)
		evmInput, err = getCallByWasmInput(suite.evmABI, caller, callData)
		suite.Require().NoError(err)
	}
	testCases := []struct {
		msg       string
		malleate  func()
		postcheck func()
		error     error
		expect    string
	}{
		{
			"caller(0x),contract(0x),calldata(normal),amount(0)",
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
			"caller(ex 20),contract(0x 20),calldata(normal),amount(0)",
			func() {
				buffer := make([]byte, 20)
				buffer[19] = 0x1
				caller = sdk.AccAddress(buffer).String()
				evmInput, err = getCallByWasmInput(suite.evmABI, caller, callData)
				suite.Require().NoError(err)
			},
			func() {
			},
			nil,
			fmt.Sprintf(evmReturnPrefix, "ex1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpxuz0nc") + fmt.Sprintf(callDataFormat, contract, "init-to-call-evm"),
		},
		{
			"caller(0x 20),contract(ex 20),calldata(normal),amount(0)",
			func() {
				contract = contractEx
				evmInput, err = getCallByWasmInput(suite.evmABI, caller, callData)
				suite.Require().NoError(err)
			},
			func() {
			},
			errors.New("the evm execute: the address prefix must be 0x"),
			"",
		},
		{
			"caller(error),contract(0x),calldata(normal),amount(0)",
			func() {
				caller = "ex1231bdjasd1"
				evmInput, err = getCallByWasmInput(suite.evmABI, caller, callData)
				suite.Require().NoError(err)
			},
			func() {
			},
			errors.New("the evm execute: decoding bech32 failed: invalid index of 1"),
			"",
		},
		{
			"caller(ex),contract(0x ),calldata(emppty),amount(0)",
			func() {
				callData = fmt.Sprintf(callDataFormat, contract, "")
				evmInput, err = getCallByWasmInput(suite.evmABI, caller, callData)
				suite.Require().NoError(err)
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
			"caller(ex),contract(0x),calldata(normal),amount(1)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(1).BigInt())
				suite.SetAccountCoins(sdk.WasmToAccAddress(suite.freeCallWasmContract), sdk.NewInt(1))
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
			"caller(ex),contract(0x),calldata(normal),amount(2)",
			func() {
				value = sdk.NewIntFromBigInt(sdk.NewDec(2).BigInt())
				suite.SetAccountCoins(sdk.WasmToAccAddress(suite.freeCallWasmContract), sdk.NewInt(1))
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
			errors.New("the evm execute: insufficient balance for transfer"),
			"",
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			reset()
			tc.malleate()
			msgServer := keeper.NewMsgServerImpl(*suite.app.VMBridgeKeeper)

			msg := types.MsgCallToEvm{Sender: caller, Evmaddr: contract, Calldata: string(evmInput), Value: value}
			result, err := msgServer.CallToEvmEvent(sdk.WrapSDKContext(suite.ctx), &msg)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
				response, err := getCallByWasmOutput(suite.evmABI, []byte(result.Response))
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expect, response)
				tc.postcheck()
			}
		})
	}

}
