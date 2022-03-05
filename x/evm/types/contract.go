package types

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

var (
	EVMModuleETHAddr  common.Address
	EVMModuleBechAddr sdk.AccAddress
)

func init() {
	EVMModuleBechAddr = authtypes.NewModuleAddress(ModuleName)
	EVMModuleETHAddr = common.BytesToAddress(EVMModuleBechAddr.Bytes())
}
