package innertx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	CosmosCallType = "cosmos"
	CosmosDepth    = 0

	SendCallName       = "cosmos-send"
	DelegateCallName   = "cosmos-delegate"
	MultiCallName      = "cosmos-multi-send"
	UndelegateCallName = "cosmos-undelegate"
	EvmCallName        = "cosmos-call"
	EvmCreateName      = "cosmos-create"

	IsAvailable = true
)

var BIG0 = big.NewInt(0)

type InnerTxKeeper interface {
	InitInnerBlock(hash string)
	UpdateInnerTx(txBytes []byte, blockHeight int64, dept int64, from, to sdk.AccAddress, callType, name string, valueWei sdk.Coins, err error)
}

func CreateInnerTx(dept int64, from, to, callType, name string, valueWei *big.Int, err error) *vm.InnerTx {
	callTx := &vm.InnerTx{
		Dept:     *big.NewInt(dept),
		From:     from,
		IsError:  false,
		To:       to,
		CallType: callType,
		Name:     name,
		ValueWei: valueWei.String(),
	}
	if err != nil {
		callTx.Error = err.Error()
		callTx.IsError = true
	}
	return callTx
}

func AddDefaultInnerTx(evm *vm.EVM, dept int64, from, to, callType, name string, valueWei *big.Int, err error) *vm.InnerTx {
	callTx := CreateInnerTx(dept, from, to, callType, name, valueWei, err)
	evm.InnerTxies = append(evm.InnerTxies, callTx)
	return callTx
}

func UpdateDefaultInnerTx(callTx *vm.InnerTx, to, callType, name string, gasused uint64) {
	callTx.To = to
	callTx.CallType = callType
	callTx.Name = name
	callTx.GasUsed = gasused
}

func ParseInnerTxAndContract(evm *vm.EVM, failed bool) ([]*vm.InnerTx, []*vm.ERC20Contract) {
	if failed {
		for _, errIx := range evm.InnerTxies {
			errIx.IsError = true
		}
		return evm.InnerTxies, evm.Contracts
	} else if len(evm.InnerTxies) > 0 {
		return evm.InnerTxies, evm.Contracts
	}
	return nil, evm.Contracts
}
