package types

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"

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

	// ModuleERC20Contract is the compiled okc erc20 contract
	ModuleERC20Contract CompiledContract

	//go:embed contracts/ModuleERC20.json
	moduleERC20Json []byte
)

const (
	IbcEvmModuleName = "ibc-evm"

	ContractMintMethod = "mint_by_okc_module"
)

func init() {
	EVMModuleBechAddr = authtypes.NewModuleAddress(IbcEvmModuleName)
	EVMModuleETHAddr = common.BytesToAddress(EVMModuleBechAddr.Bytes())

	if err := json.Unmarshal(moduleERC20Json, &ModuleERC20Contract); err != nil {
		panic(err)
	}
	if len(ModuleERC20Contract.Bin) == 0 {
		panic("load contract failed")
	}
}

func (c CompiledContract) ValidBasic() error {
	if len(c.Bin) == 0 {
		return errors.New("empty bin data")
	}
	_, err := hex.DecodeString(c.Bin)
	if nil != err {
		return err
	}
	return nil
}

func MustMarshalCompileContract(data CompiledContract) []byte {
	ret, err := MarshalCompileContract(data)
	if nil != err {
		panic(err)
	}
	return ret
}

func MarshalCompileContract(data CompiledContract) ([]byte, error) {
	return json.Marshal(data)
}

func MustUnmarshalCompileContract(data []byte) CompiledContract {
	ret, err := UnmarshalCompileContract(data)
	if nil != err {
		panic(err)
	}
	return ret
}

func UnmarshalCompileContract(data []byte) (CompiledContract, error) {
	var ret CompiledContract
	err := json.Unmarshal(data, &ret)
	if nil != err {
		return CompiledContract{}, err
	}
	return ret, nil
}

func GetInternalTemplateContract() []byte {
	return moduleERC20Json
}
