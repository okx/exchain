package types

import (
	_ "embed"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// CompiledContract contains compiled bytecode and abi
type CompiledContract struct {
	ABI abi.ABI
	Bin string
}

var (
	EVMModuleETHAddr  common.Address
	EVMModuleBechAddr sdk.AccAddress

	// ModuleERC20Contract is the compiled cronos erc20 contract
	ModuleERC20Contract CompiledContract
	// TODO cronos ---> oec

	//go:embed contracts/ModuleERC20.json
	moduleERC20Json []byte
)

func init() {
	EVMModuleBechAddr = authtypes.NewModuleAddress(ModuleName)
	EVMModuleETHAddr = common.BytesToAddress(EVMModuleBechAddr.Bytes())
	// 0x603871c2ddd41c26Ee77495E2E31e6De7f9957e0

	if err := json.Unmarshal(moduleERC20Json, &ModuleERC20Contract); err != nil {
		panic(err)
	}
	if len(ModuleERC20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
