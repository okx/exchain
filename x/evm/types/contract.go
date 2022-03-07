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

	// ModuleCRC20Contract is the compiled cronos erc20 contract
	ModuleCRC20Contract CompiledContract

	//go:embed contracts/ModuleERC20.json
	moduleERC20Json []byte
)

func init() {
	EVMModuleBechAddr = authtypes.NewModuleAddress(ModuleName)
	EVMModuleETHAddr = common.BytesToAddress(EVMModuleBechAddr.Bytes())

	if err := json.Unmarshal(moduleERC20Json, &ModuleCRC20Contract); err != nil {
		panic(err)
	}
	if len(ModuleCRC20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
