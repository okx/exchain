package types

import (
	"github.com/ethereum/go-ethereum/common"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"testing"
)

func TestGetEVMABIConfig(t *testing.T) {
	addr := authtypes.NewModuleAddress(ModuleName)
	ethAddr := common.BytesToAddress(addr.Bytes())
	t.Log(addr.String(), ethAddr.String())
}
