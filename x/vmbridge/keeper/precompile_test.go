package keeper_test

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/vmbridge/types"
	"math/big"
	"strconv"
	"strings"
)

var (
	testPrecompileCodeA    = "608060405273db327e55ca2c68b23f83a0fbe29b592702e1d4d76000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006457600080fd5b50610b76806100746000396000f3fe60806040526004361061004a5760003560e01c80635b3082c21461004f57806363de1b5d1461007f5780636bbb9b13146100af5780638381f58a146100df578063be2b0ac21461010a575b600080fd5b610069600480360381019061006491906106cc565b610147565b60405161007691906108ba565b60405180910390f35b61009960048036038101906100949190610670565b610161565b6040516100a69190610898565b60405180910390f35b6100c960048036038101906100c49190610744565b610314565b6040516100d69190610898565b60405180910390f35b3480156100eb57600080fd5b506100f46104ca565b6040516101019190610913565b60405180910390f35b34801561011657600080fd5b50610131600480360381019061012c91906105de565b6104d0565b60405161013e91906108ba565b60405180910390f35b606060405180602001604052806000815250905092915050565b60606001805461017191906109c7565b60018190555060008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1634866040516024016101c391906108ba565b6040516020818303038152906040527fbe2b0ac2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161024d9190610881565b60006040518083038185875af1925050503d806000811461028a576040519150601f19603f3d011682016040523d82523d6000602084013e61028f565b606091505b509150915083156102f557816102a457600080fd5b6000818060200190518101906102ba9190610627565b90507fe390e3d6b4766bc311796e6b5ce75dd6d51f0cb55cea58be963a5e7972ade65c816040516102eb91906108ba565b60405180910390a1505b6001805461030391906109c7565b600181905550809250505092915050565b60606001805461032491906109c7565b60018190555060008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163487876040516024016103789291906108dc565b6040516020818303038152906040527f5b3082c2000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516104029190610881565b60006040518083038185875af1925050503d806000811461043f576040519150601f19603f3d011682016040523d82523d6000602084013e610444565b606091505b509150915083156104aa578161045957600080fd5b60008180602001905181019061046f9190610627565b90507fe390e3d6b4766bc311796e6b5ce75dd6d51f0cb55cea58be963a5e7972ade65c816040516104a091906108ba565b60405180910390a1505b600180546104b891906109c7565b60018190555080925050509392505050565b60015481565b6060604051806020016040528060008152509050919050565b60006104fc6104f784610953565b61092e565b90508281526020810184848401111561051857610517610b09565b5b610523848285610a33565b509392505050565b600061053e61053984610953565b61092e565b90508281526020810184848401111561055a57610559610b09565b5b610565848285610a42565b509392505050565b60008135905061057c81610b29565b92915050565b600082601f83011261059757610596610b04565b5b81356105a78482602086016104e9565b91505092915050565b600082601f8301126105c5576105c4610b04565b5b81516105d584826020860161052b565b91505092915050565b6000602082840312156105f4576105f3610b13565b5b600082013567ffffffffffffffff81111561061257610611610b0e565b5b61061e84828501610582565b91505092915050565b60006020828403121561063d5761063c610b13565b5b600082015167ffffffffffffffff81111561065b5761065a610b0e565b5b610667848285016105b0565b91505092915050565b6000806040838503121561068757610686610b13565b5b600083013567ffffffffffffffff8111156106a5576106a4610b0e565b5b6106b185828601610582565b92505060206106c28582860161056d565b9150509250929050565b600080604083850312156106e3576106e2610b13565b5b600083013567ffffffffffffffff81111561070157610700610b0e565b5b61070d85828601610582565b925050602083013567ffffffffffffffff81111561072e5761072d610b0e565b5b61073a85828601610582565b9150509250929050565b60008060006060848603121561075d5761075c610b13565b5b600084013567ffffffffffffffff81111561077b5761077a610b0e565b5b61078786828701610582565b935050602084013567ffffffffffffffff8111156107a8576107a7610b0e565b5b6107b486828701610582565b92505060406107c58682870161056d565b9150509250925092565b60006107da82610984565b6107e4818561099a565b93506107f4818560208601610a42565b6107fd81610b18565b840191505092915050565b600061081382610984565b61081d81856109ab565b935061082d818560208601610a42565b80840191505092915050565b60006108448261098f565b61084e81856109b6565b935061085e818560208601610a42565b61086781610b18565b840191505092915050565b61087b81610a29565b82525050565b600061088d8284610808565b915081905092915050565b600060208201905081810360008301526108b281846107cf565b905092915050565b600060208201905081810360008301526108d48184610839565b905092915050565b600060408201905081810360008301526108f68185610839565b9050818103602083015261090a8184610839565b90509392505050565b60006020820190506109286000830184610872565b92915050565b6000610938610949565b90506109448282610a75565b919050565b6000604051905090565b600067ffffffffffffffff82111561096e5761096d610ad5565b5b61097782610b18565b9050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b60006109d282610a29565b91506109dd83610a29565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610a1257610a11610aa6565b5b828201905092915050565b60008115159050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015610a60578082015181840152602081019050610a45565b83811115610a6f576000848401525b50505050565b610a7e82610b18565b810181811067ffffffffffffffff82111715610a9d57610a9c610ad5565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b610b3281610a1d565b8114610b3d57600080fd5b5056fea2646970667358221220a0661117352f3c64e9a2b9fc3db330faf807a6e0345cd74e8381ad531347486464736f6c63430008070033"
	testPrecompileABIAJson = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"callToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"msgData\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"requireASuccess\",\"type\":\"bool\"}],\"name\":\"callWasm\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"pushLog\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"msgData\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"requireASuccess\",\"type\":\"bool\"}],\"name\":\"queryWasm\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"number\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"queryToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	testPrecompileCodeB    = "608060405234801561001057600080fd5b5061094a806100206000396000f3fe6080604052600436106100345760003560e01c80638381f58a14610039578063988950c714610064578063e3cb5bf114610094575b600080fd5b34801561004557600080fd5b5061004e6100c4565b60405161005b919061069e565b60405180910390f35b61007e6004803603810190610079919061047c565b6100ca565b60405161008b9190610607565b60405180910390f35b6100ae60048036038101906100a991906103f9565b610216565b6040516100bb9190610607565b60405180910390f35b60005481565b606060016000546100db9190610752565b6000819055506000808773ffffffffffffffffffffffffffffffffffffffff163488888860405160240161011193929190610659565b6040516020818303038152906040527f6bbb9b13000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161019b91906105f0565b60006040518083038185875af1925050503d80600081146101d8576040519150601f19603f3d011682016040523d82523d6000602084013e6101dd565b606091505b509150915060016000546101f19190610752565b6000819055508315610208578161020757600080fd5b5b809250505095945050505050565b606060016000546102279190610752565b6000819055506000808673ffffffffffffffffffffffffffffffffffffffff1634878760405160240161025b929190610629565b6040516020818303038152906040527f63de1b5d000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516102e591906105f0565b60006040518083038185875af1925050503d8060008114610322576040519150601f19603f3d011682016040523d82523d6000602084013e610327565b606091505b5091509150600160005461033b9190610752565b6000819055508315610352578161035157600080fd5b5b8092505050949350505050565b600061037261036d846106de565b6106b9565b90508281526020810184848401111561038e5761038d6108c6565b5b6103998482856107f0565b509392505050565b6000813590506103b0816108e6565b92915050565b6000813590506103c5816108fd565b92915050565b600082601f8301126103e0576103df6108c1565b5b81356103f084826020860161035f565b91505092915050565b60008060008060808587031215610413576104126108d0565b5b6000610421878288016103a1565b945050602085013567ffffffffffffffff811115610442576104416108cb565b5b61044e878288016103cb565b935050604061045f878288016103b6565b9250506060610470878288016103b6565b91505092959194509250565b600080600080600060a08688031215610498576104976108d0565b5b60006104a6888289016103a1565b955050602086013567ffffffffffffffff8111156104c7576104c66108cb565b5b6104d3888289016103cb565b945050604086013567ffffffffffffffff8111156104f4576104f36108cb565b5b610500888289016103cb565b9350506060610511888289016103b6565b9250506080610522888289016103b6565b9150509295509295909350565b610538816107ba565b82525050565b60006105498261070f565b6105538185610725565b93506105638185602086016107ff565b61056c816108d5565b840191505092915050565b60006105828261070f565b61058c8185610736565b935061059c8185602086016107ff565b80840191505092915050565b60006105b38261071a565b6105bd8185610741565b93506105cd8185602086016107ff565b6105d6816108d5565b840191505092915050565b6105ea816107e6565b82525050565b60006105fc8284610577565b915081905092915050565b60006020820190508181036000830152610621818461053e565b905092915050565b6000604082019050818103600083015261064381856105a8565b9050610652602083018461052f565b9392505050565b6000606082019050818103600083015261067381866105a8565b9050818103602083015261068781856105a8565b9050610696604083018461052f565b949350505050565b60006020820190506106b360008301846105e1565b92915050565b60006106c36106d4565b90506106cf8282610832565b919050565b6000604051905090565b600067ffffffffffffffff8211156106f9576106f8610892565b5b610702826108d5565b9050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b600082825260208201905092915050565b600061075d826107e6565b9150610768836107e6565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561079d5761079c610863565b5b828201905092915050565b60006107b3826107c6565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b8381101561081d578082015181840152602081019050610802565b8381111561082c576000848401525b50505050565b61083b826108d5565b810181811067ffffffffffffffff8211171561085a57610859610892565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b6108ef816107a8565b81146108fa57600080fd5b50565b610906816107ba565b811461091157600080fd5b5056fea26469706673582212207229166639d952ed2d2e9407fd0d2b5cedfdf8cfc2ed0f9e47b488ea5806f71e64736f6c63430008070033"
	testPrecompileABIBJson = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractA\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"msgData\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"requireASuccess\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"requireBSuccess\",\"type\":\"bool\"}],\"name\":\"callWasm\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractA\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"msgData\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"requireASuccess\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"requireBSuccess\",\"type\":\"bool\"}],\"name\":\"queryWasm\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"response\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"number\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]\n"
	testPrecompileABIA     abi.ABI
	testPrecompileABIB     abi.ABI
	callWasmMethod         = "callWasm"
	queryWasmMethod        = "queryWasm"
)

func init() {
	var err error
	testPrecompileABIA, err = abi.JSON(strings.NewReader(testPrecompileABIAJson))
	if err != nil {
		panic(err)
	}
	testPrecompileABIB, err = abi.JSON(strings.NewReader(testPrecompileABIBJson))
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) precompile_setup() (contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
	contractA = suite.deployEvmContract(testPrecompileCodeA)
	contractB = suite.deployEvmContract(testPrecompileCodeB)
	wasmContract = suite.deployWasmContract("precompile.wasm", []byte(fmt.Sprintf("{\"decimals\":10,\"initial_balances\":[{\"address\":\"%s\",\"amount\":\"100000000\"}],\"name\":\"my test token\", \"symbol\":\"MTT\"}", common.BytesToAddress(suite.addr).String())))
	return
}

func (suite *KeeperTestSuite) TestPrecompileHooks() {
	//contractA, contractB := suite.precompile_setup()
	cmBridgePrecompileAddress := common.HexToAddress("0xDb327e55CA2C68b23f83a0fbe29b592702e1d4d7")

	testAddr := common.HexToAddress("0x09084cc9c3e579Fd4aa383D3fD6C543f7FFC36c7")
	callWasmMsgFormat := "{\"transfer\":{\"amount\":\"%d\",\"recipient\":\"%s\"}}"
	caller := common.BytesToAddress(suite.addr)
	contract := cmBridgePrecompileAddress
	amount := big.NewInt(0)
	evmCalldata := make([]byte, 0)
	reset := func() {
		caller = common.BytesToAddress(suite.addr)
		contract = cmBridgePrecompileAddress
		amount = big.NewInt(0)
		evmCalldata = make([]byte, 0)
	}
	var err error
	testCases := []struct {
		msg       string
		malleate  func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress)
		postcheck func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress)
		error     error
		success   bool
	}{
		{
			"call to wasm at tx",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(10, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(99999990, result)

				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("the result wasm contract data", resultStr)
			},
			nil,
			true,
		},
		{
			"call to wasm at tx with okt",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
				amount = sdk.NewDec(1).BigInt()
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(10, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(99999990, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(9999), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(sdk.NewDec(1), coin[0].Amount)

				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("the result wasm contract data", resultStr)
			},
			nil,
			true,
		},
		{
			"call to wasm at tx with okt insuffent",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
				amount = sdk.NewDec(10001).BigInt()
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(100000000, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
			},
			errors.New("insufficient balance for transfer"),
			false,
		},
		{
			"call to wasm at tx with wasm failed",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, common.Address{0x1}.String(), hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(100000000, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
			},
			errors.New("not found: contract"),
			false,
		},
		{
			"call to wasm at tx with err wasmContractAddr ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, "common.Address{0x1}.String()", hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(100000000, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
			},
			errors.New("encoding/hex: invalid byte: U+006F 'o'"),
			false,
		},
		{
			"call to wasm at tx with err wasmContractAddr ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, "0x1234", hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(100000000, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
			},
			errors.New("incorrect address length"),
			false,
		},
		{
			"call to wasm at tx with err wasm calldata ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := "fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())"
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(100000000, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
			},
			errors.New("invalid character 'm' in literal false (expecting 'a')"),
			false,
		},
		{
			"call to wasm at tx with err hex wasm calldata ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := "fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())"
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileCallToWasm, wasmContract.String(), wasmCallData)
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), common.BytesToAddress(suite.addr).String())
				suite.Require().NoError(err)
				suite.Require().Equal(100000000, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
			},
			errors.New("encoding/hex: invalid byte: U+006D 'm'"),
			false,
		},
		{
			"call to wasm at contractA",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = testPrecompileABIA.Pack(callWasmMethod, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)), true)
				suite.Require().NoError(err)
				contract = contractA
				suite.transferWasmBalance(*ctx, wasmContract, sdk.WasmAddress(contractA.Bytes()), 10)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(10, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), contractA.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))

				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("the result wasm contract data", resultStr)

				numberMethod, err := testPrecompileABIA.Pack("number")
				suite.Require().NoError(err)
				numberResult, err := suite.queryEvmContract(*ctx, contractA, numberMethod)
				suite.Require().NoError(err)
				pack, err := testPrecompileABIA.Methods["number"].Outputs.Unpack(numberResult)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pack))
				number := pack[0].(*big.Int)
				suite.Require().Equal(int64(2), number.Int64())
			},
			nil,
			false,
		},
		{
			"call to wasm at contractA with okt",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = testPrecompileABIA.Pack(callWasmMethod, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)), true)
				suite.Require().NoError(err)
				contract = contractA
				suite.transferWasmBalance(*ctx, wasmContract, sdk.WasmAddress(contractA.Bytes()), 10)
				amount = sdk.NewDec(1).BigInt()
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(10, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), contractA.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(9999), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(sdk.NewDec(1), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, contractA.Bytes())
				suite.Require().Equal(0, len(coin))

				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("the result wasm contract data", resultStr)

				numberMethod, err := testPrecompileABIA.Pack("number")
				suite.Require().NoError(err)
				numberResult, err := suite.queryEvmContract(*ctx, contractA, numberMethod)
				suite.Require().NoError(err)
				pack, err := testPrecompileABIA.Methods["number"].Outputs.Unpack(numberResult)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pack))
				number := pack[0].(*big.Int)
				suite.Require().Equal(int64(2), number.Int64())
			},
			nil,
			false,
		},
		{
			"call to wasm at contractA with okt but tx failed",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 10, testAddr.String())
				evmCalldata, err = testPrecompileABIA.Pack(callWasmMethod, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)), true)
				suite.Require().NoError(err)
				contract = contractA
				//suite.transferWasmBalance(*ctx, wasmContract, sdk.WasmAddress(contractA.Bytes()), 10)
				amount = sdk.NewDec(1).BigInt()
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), contractA.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(10000), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
				coin = suite.queryPrecompileCoins(*ctx, contractA.Bytes())
				suite.Require().Equal(0, len(coin))

				suite.Require().Nil(data)

				numberMethod, err := testPrecompileABIA.Pack("number")
				suite.Require().NoError(err)
				numberResult, err := suite.queryEvmContract(*ctx, contractA, numberMethod)
				suite.Require().NoError(err)
				pack, err := testPrecompileABIA.Methods["number"].Outputs.Unpack(numberResult)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pack))
				number := pack[0].(*big.Int)
				suite.Require().Equal(int64(0), number.Int64())
			},
			errors.New("execution reverted"),
			false,
		},
		{
			"call to wasm at contractA with okt but tx not failed",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmCallData := fmt.Sprintf(callWasmMsgFormat, 100, testAddr.String())
				evmCalldata, err = testPrecompileABIA.Pack(callWasmMethod, wasmContract.String(), hex.EncodeToString([]byte(wasmCallData)), false)
				suite.Require().NoError(err)
				contract = contractA
				suite.transferWasmBalance(*ctx, wasmContract, sdk.WasmAddress(contractA.Bytes()), 10)
				amount = sdk.NewDec(1).BigInt()
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				result, err := suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), testAddr.String())
				suite.Require().NoError(err)
				suite.Require().Equal(0, result)

				result, err = suite.queryPrecompileWasmBalance(*ctx, caller.String(), wasmContract.String(), contractA.String())
				suite.Require().NoError(err)
				suite.Require().Equal(10, result)

				coin := suite.queryPrecompileCoins(*ctx, sdk.AccAddress(caller.Bytes()))
				suite.Require().Equal(sdk.NewDec(9999), coin[0].Amount)
				coin = suite.queryPrecompileCoins(*ctx, wasmContract.Bytes())
				suite.Require().Equal(0, len(coin))
				coin = suite.queryPrecompileCoins(*ctx, contractA.Bytes())
				suite.Require().Equal(sdk.NewDec(1), coin[0].Amount)

				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("", resultStr)

				numberMethod, err := testPrecompileABIA.Pack("number")
				suite.Require().NoError(err)
				numberResult, err := suite.queryEvmContract(*ctx, contractA, numberMethod)
				suite.Require().NoError(err)
				pack, err := testPrecompileABIA.Methods["number"].Outputs.Unpack(numberResult)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pack))
				number := pack[0].(*big.Int)
				suite.Require().Equal(int64(2), number.Int64())
			},
			nil,
			false,
		},
		{
			"smart query to wasm  at tx ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				testQueryMsg := fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", common.BytesToAddress(suite.addr).String())
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: wasmContract.String(), Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("{\"balance\":\"100000000\"}", resultStr)
			},
			nil,
			true,
		},
		{
			"smart query to wasm  at tx contract not found ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				testQueryMsg := fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", common.BytesToAddress(suite.addr).String())
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: common.Address{0x1}.String(), Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				suite.Require().Nil(data)
			},
			errors.New("codespace: wasm, code: 8"),
			true,
		},
		{
			"smart query to wasm  at tx contractaddr is error ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				testQueryMsg := fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", common.BytesToAddress(suite.addr).String())
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: "common.Address{0x1}.String()", Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				suite.Require().Nil(data)
			},
			errors.New("codespace: sdk, code: 7"),
			true,
		},
		{
			"smart query to wasm  at tx smart json is error1 ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				testQueryMsg := fmt.Sprintf("{\"error\":{\"address\":\"%s\"}}", common.BytesToAddress(suite.addr).String())
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: wasmContract.String(), Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				suite.Require().Nil(data)
			},
			errors.New("codespace: wasm, code: 9"),
			true,
		},
		{
			"smart query to wasm  at tx smart json is error2 ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				testQueryMsg := fmt.Sprintf("{\"balance\":{\"error\":\"%s\"}}", common.BytesToAddress(suite.addr).String())
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: wasmContract.String(), Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				suite.Require().Nil(data)
			},
			errors.New("codespace: wasm, code: 9"),
			true,
		},
		{
			"smart query to wasm  at tx smart json is error3 ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				testQueryMsg := fmt.Sprintf("\"\":{\"error\":\"%s\"}}", common.BytesToAddress(suite.addr).String())
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: wasmContract.String(), Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				suite.Require().Nil(data)
			},
			errors.New("codespace: wasm, code: 14"),
			true,
		},
		{
			"raw query to wasm  at tx ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				raw, err := hex.DecodeString("0006636F6E666967636F6E7374616E7473")
				suite.Require().NoError(err)
				wasmsmartRequest := wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: wasmContract.String(), Key: raw}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("{\"name\":\"my test token\",\"symbol\":\"MTT\",\"decimals\":10}", resultStr)
			},
			nil,
			true,
		},
		{
			"raw query to wasm  at tx contract no found",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				raw, err := hex.DecodeString("0006636F6E666967636F6E7374616E7473")
				suite.Require().NoError(err)
				wasmsmartRequest := wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: common.Address{0x1}.String(), Key: raw}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("", resultStr)
			},
			nil,
			true,
		},
		{
			"raw query to wasm  at tx key is not exist",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				raw, err := hex.DecodeString("0006636F6E666967636F6E7374616E7472")
				suite.Require().NoError(err)
				wasmsmartRequest := wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: wasmContract.String(), Key: raw}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("", resultStr)
			},
			nil,
			true,
		},
		{
			"raw query to wasm  at tx key is empty",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				raw, err := hex.DecodeString("")
				suite.Require().NoError(err)
				wasmsmartRequest := wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: wasmContract.String(), Key: raw}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				resultStr, err := decodeCallWasmOutput(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal("", resultStr)
			},
			nil,
			true,
		},
		{
			"contractinfo query to wasm  at tx ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmsmartRequest := wasmvmtypes.WasmQuery{ContractInfo: &wasmvmtypes.ContractInfoQuery{ContractAddr: wasmContract.String()}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				pack, err := testPrecompileABIA.Methods[callWasmMethod].Outputs.Unpack(data.Ret)
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(pack))
				suite.Require().Equal("{\"code_id\":3,\"creator\":\"0x0102030405060708091011121314151617181920\",\"admin\":\"0x0102030405060708091011121314151617181920\",\"pinned\":false}", string(pack[0].([]byte)))
			},
			nil,
			true,
		},
		{
			"contractinfo query to  at tx contract no found  ",
			func(ctx *sdk.Context, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				wasmsmartRequest := wasmvmtypes.WasmQuery{ContractInfo: &wasmvmtypes.ContractInfoQuery{ContractAddr: common.Address{0x1}.String()}}
				buff, err := json.Marshal(wasmsmartRequest)
				suite.Require().NoError(err)
				evmCalldata, err = types.PreCompileABI.Pack(types.PrecompileQueryToWasm, hex.EncodeToString(buff))
				suite.Require().NoError(err)
			},
			func(ctx *sdk.Context, data *evmtypes.ResultData, contractA, contractB common.Address, wasmContract sdk.WasmAddress) {
				suite.Require().Nil(data)
			},
			errors.New("no such contract: 0x0100000000000000000000000000000000000000"),
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()
			contractA, contractB, wasmContract := suite.precompile_setup()
			reset()
			subCtx, _ := suite.ctx.CacheContext()
			tc.malleate(&subCtx, contractA, contractB, wasmContract)
			_, evmResult, err := suite.app.VMBridgeKeeper.CallEvm(subCtx, caller, &contract, amount, evmCalldata)
			if tc.error != nil {
				suite.Require().EqualError(err, tc.error.Error())
			} else {
				suite.Require().NoError(err)
			}
			tc.postcheck(&subCtx, evmResult, contractA, contractB, wasmContract)
		})
	}
}

func (suite *KeeperTestSuite) queryPrecompileWasmBalance(ctx sdk.Context, caller, wasmContract, to string) (int, error) {
	testQueryMsg := fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", to)
	wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: wasmContract, Msg: []byte(testQueryMsg)}}
	buff, err := json.Marshal(wasmsmartRequest)
	if err != nil {
		return 0, err
	}
	ret, err := suite.app.VMBridgeKeeper.QueryToWasm(ctx, caller, buff)
	if err != nil {
		return 0, err
	}
	var response struct {
		Balance string `json:"balance"`
	}
	if err := json.Unmarshal(ret, &response); err != nil {
		return 0, err
	}
	return strconv.Atoi(response.Balance)
}

func (suite *KeeperTestSuite) queryPrecompileCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := suite.app.AccountKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0))
	}
	return acc.GetCoins()
}

func (suite *KeeperTestSuite) transferWasmBalance(ctx sdk.Context, wasmcontract, to sdk.WasmAddress, amount int) {
	msg := []byte(fmt.Sprintf("{\"transfer\":{\"amount\":\"%d\",\"recipient\":\"%s\"}}", amount, to.String()))
	suite.executeWasmContract(ctx, suite.addr.Bytes(), wasmcontract, msg, sdk.Coins{})
}

func decodeCallWasmOutput(input []byte) (string, error) {
	pack, err := testPrecompileABIA.Methods[callWasmMethod].Outputs.Unpack(input)
	if err != nil {
		return "", err
	}
	if len(pack) != 1 {
		return "", errors.New("decodeCallToWasmOutput failed: got multi result")
	}
	buff := pack[0].([]byte)
	if len(buff) >= 96 {
		pack, err = testPrecompileABIA.Methods[callWasmMethod].Outputs.Unpack(buff)
		if err != nil {
			return "", err
		}
		if len(pack) != 1 {
			return "", errors.New("decodeCallToWasmOutput failed: got multi result")
		}
	}
	return string(pack[0].([]byte)), nil
}
