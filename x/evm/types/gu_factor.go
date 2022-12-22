package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"math/big"
)

var (
	DefaultGuFactor = sdk.OneDec()
)

type GuFactor struct {
	Factor sdk.Dec `json:"gu_factor" yaml:"gu_factor"`
}

func (factor GuFactor) ValidateBasic() error {
	// proportion must (0,1]. if proportion <= 0 or  > 1
	if factor.Factor.LTE(sdk.ZeroDec()) {
		return ErrGUFactor
	}
	return nil
}

func UnmarshalGuFactor(data string) (GuFactor, error) {
	var factor GuFactor
	err := json.Unmarshal([]byte(data), &factor)
	return factor, err
}

type GuFactorHook struct {
}

func (hook GuFactorHook) UpdateGuFactor(csdb *CommitStateDB, op vm.OpCode, from, to common.Address, input []byte, value *big.Int) {
	if op == vm.CALL || op == vm.DELEGATECALL || op == vm.CALLCODE {
		//contract not allowed to call delegateCall and callcode ,check from and input is blocked 。STATICCALL could not check because ，it's readonly.
		if len(input) > 4 {
			input = input[:4]
		}
		method := hexutil.Encode(input)

		if bc := csdb.GetContractMethodBlockedByAddress(to.Bytes()); bc != nil {
			if contractMethod := bc.BlockMethods.GetMethod(method); contractMethod != nil {
				if factor := contractMethod.GetGuFactor(); factor != nil {
					if factor.Factor.GT(csdb.GuFactor) {
						csdb.GuFactor = factor.Factor
					}
				}
			}
		}
	}
}
