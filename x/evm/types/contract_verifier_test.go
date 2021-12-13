package types_test

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

func (suite *StateDBTestSuite) TestContractVerifier_Verify() {
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	//addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()
	sign := hexutil.Encode([]byte("hehe"))
	bcMethodOne1 := types.BlockedContract{
		Address: addr1,
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  sign,
				Extra: "aaaa()",
			},
		},
	}

	testCase := []struct {
		name    string
		execute func() sdk.Error
		expPass bool
	}{
		{
			"stateDB is ethereum",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(&ethstate.StateDB{}, vm.CALL, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), nil, nil)
			},
			false,
		},
		{
			"blocked list is disable",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = false
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.CALL, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), nil, nil)
			},
			true,
		},
		{
			"contract address is all method blocked (CALL)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractBlockedList(types.AddressList{addr1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.CALL, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), nil, nil)
			},
			true,
		},
		{
			"contract address is all method blocked (SELFDESTRUCT)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractBlockedList(types.AddressList{addr1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.SELFDESTRUCT, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), nil, nil)
			},
			true,
		},
		{
			"contract method is blocked (SELFDESTRUCT)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.SELFDESTRUCT, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), nil, nil)
			},
			false,
		},
		{
			"contract method is blocked (CALL)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.CALL, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), hexutil.MustDecode(sign), nil)
			},
			false,
		},
		{
			"contract method is blocked (DELEGATECALL)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.DELEGATECALL, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), hexutil.MustDecode(sign), nil)
			},
			false,
		},
		{
			"contract method is blocked (DELEGATECALL)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.CALLCODE, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), hexutil.MustDecode(sign), nil)
			},
			false,
		},
		{
			"contract method is not blocked (CALL)",
			func() (err sdk.Error) {
				defer func() {
					if e := recover(); e != nil {
						switch rType := e.(type) {
						case types.ErrContractBlockedVerify:
							err = types.ErrCallBlockedContract(rType.Descriptor)
						default:
							panic(e)
						}
					}
				}()

				params := suite.stateDB.GetParams()
				params.EnableContractBlockedList = true
				suite.stateDB.SetContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
				verifier := types.NewContractVerifier(params)

				return verifier.Verify(suite.stateDB, vm.CALLCODE, ethcmn.BytesToAddress(addr1), ethcmn.BytesToAddress(addr1), hexutil.MustDecode("0x1111"), nil)
			},
			true,
		},
	}
	for _, tc := range testCase {
		err := tc.execute()

		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
