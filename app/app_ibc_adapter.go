package app

import "github.com/okex/exchain/libs/cosmos-sdk/codec/types"

func MakeIBC()types.InterfaceRegistry{
	interfaceReg:=types.NewInterfaceRegistry()
	ModuleBasics.RegisterInterfaces(interfaceReg)
	return interfaceReg
}
