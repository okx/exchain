package dydx

import "github.com/ethereum/go-ethereum/accounts/abi"

var (
	SolTyBytes32 abi.Type
	SolTyUint256 abi.Type
	SolTyAddress abi.Type
	SolTyBool    abi.Type
)

func init() {
	var err error
	SolTyBytes32, err = abi.NewType("bytes32", "", nil)
	if err != nil {
		panic(err)
	}
	SolTyUint256, err = abi.NewType("uint256", "", nil)
	if err != nil {
		panic(err)
	}
	SolTyAddress, err = abi.NewType("address", "", nil)
	if err != nil {
		panic(err)
	}
	SolTyBool, err = abi.NewType("bool", "", nil)
	if err != nil {
		panic(err)
	}
}
