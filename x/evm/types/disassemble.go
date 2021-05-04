package types

import (
	"github.com/ethereum/go-ethereum/core/asm"
	"github.com/ethereum/go-ethereum/core/vm"
)

func IsContractReadonly(code []byte) bool {
	it := asm.NewInstructionIterator(code)
	for it.Next() {
		num := it.PC()
		if num == 0 {
			num = 1
		}
		switch it.Op() {
		case vm.SSTORE, vm.CREATE, vm.CREATE2:
			return false
		}
	}
	return true
}
