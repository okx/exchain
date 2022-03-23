package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"math/big"
)

// ContractVerifier which verify contract method whether blocked
type ContractVerifier struct {
	params Params
}

// NewContractVerifier return a point of ContractVerifier
func NewContractVerifier(params Params) *ContractVerifier {
	return &ContractVerifier{params: params}
}

// Verify check the contract whether is blocked.
// It never return error,because in okc chain if it blocked, not allow to execute next opCode.In Ethereum call failed,the call err is deal in contract code ether evm.
// If current call/delegatecall/callcode contract method is blocked,it will be panic,then it's deal logic at defer{recover}.
// If contract all method blocked,it will not be panic in Verify. it will be panic in stateDb.GetCode().
func (cv ContractVerifier) Verify(stateDB vm.StateDB, op vm.OpCode, from, to common.Address, input []byte, value *big.Int) error {
	csdb, ok := stateDB.(*CommitStateDB)
	//If stateDB is not okc stateDB ,then return error
	if !ok {
		panic(ErrContractBlockedVerify{"unknown stateDB expected CommitStateDB"})
	}
	//check whether contract has been blocked
	if !cv.params.EnableContractBlockedList {
		return nil
	}
	if op == vm.SELFDESTRUCT {
		//contract not allowed selfdestruct,check from is blocked
		bc := csdb.GetContractMethodBlockedByAddress(from.Bytes())
		if bc != nil && !bc.IsAllMethodBlocked() {
			err := ErrContractBlockedVerify{fmt.Sprintf("Contract %s has been blocked. It's not allow to SELFDESTRUCT", from.String())}
			panic(err)
		}
	} else if op == vm.CALL || op == vm.DELEGATECALL || op == vm.CALLCODE {
		//contract not allowed to call delegateCall and callcode ,check from and input is blocked 。STATICCALL could not check because ，it's readonly.
		if len(input) > 4 {
			input = input[:4]
		}
		method := hexutil.Encode(input)
		if csdb.IsContractMethodBlocked(sdk.AccAddress(to.Bytes()), method) {
			err := ErrContractBlockedVerify{fmt.Sprintf("The method %s of contract %s has been blocked. It's not allow to %s", method, to.String(), op.String())}
			panic(err)
		}
	}
	return nil
}
