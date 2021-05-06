package types

import (
	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"
)

func IsContractReadonly(code []byte) bool {
	it := asm.NewInstructionIterator(code)
	for it.Next() {
		switch it.Op() {
		case vm.SSTORE, vm.CREATE, vm.CREATE2, vm.DELEGATECALL, vm.CALLCODE, vm.LOG0, vm.LOG1, vm.LOG2, vm.LOG3, vm.LOG4:
			return false
		}
	}
	return true
}
