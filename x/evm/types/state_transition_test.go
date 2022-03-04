package types_test

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
)

var (
	callAddr                  = "0x2B2641734D81a6B93C9aE1Ee6290258FB6666921"
	callCode                  = "0x608060405260043610610083576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630c55699c1461008857806350cd4df2146100b35780637811c6c1146100de578063a6516bda14610121578063a7126c2d14610164578063a9421619146101a7578063d3ab86a1146101ea575b600080fd5b34801561009457600080fd5b5061009d610241565b6040518082815260200191505060405180910390f35b3480156100bf57600080fd5b506100c8610247565b6040518082815260200191505060405180910390f35b3480156100ea57600080fd5b5061011f600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061024d565b005b34801561012d57600080fd5b50610162600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610304565b005b34801561017057600080fd5b506101a5600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506103bb565b005b3480156101b357600080fd5b506101e8600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610470565b005b3480156101f657600080fd5b506101ff610527565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60005481565b60015481565b600060405180807f696e6328290000000000000000000000000000000000000000000000000000008152506005019050604051809103902090508173ffffffffffffffffffffffffffffffffffffffff16817c010000000000000000000000000000000000000000000000000000000090046040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016000604051808303816000875af292505050505050565b600060405180807f6f6e6328290000000000000000000000000000000000000000000000000000008152506005019050604051809103902090508173ffffffffffffffffffffffffffffffffffffffff16817c010000000000000000000000000000000000000000000000000000000090046040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016000604051808303816000875af192505050505050565b600060405180807f696e6328290000000000000000000000000000000000000000000000000000008152506005019050604051809103902090508173ffffffffffffffffffffffffffffffffffffffff16817c010000000000000000000000000000000000000000000000000000000090046040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381865af492505050505050565b600060405180807f696e6328290000000000000000000000000000000000000000000000000000008152506005019050604051809103902090508173ffffffffffffffffffffffffffffffffffffffff16817c010000000000000000000000000000000000000000000000000000000090046040518163ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016000604051808303816000875af192505050505050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a7230582003530ba5d655e02d210fb630e4067ad896add11d3c99c6c69165d11ce4855ca90029"
	callAcc, _                = sdk.AccAddressFromBech32(callAddr)
	callEthAcc                = common.BytesToAddress(callAcc.Bytes())
	callBuffer                = hexutil.MustDecode(callCode)
	blockedAddr               = "0xf297Ab486Be410A2649901849B0477D519E99960"
	blockedCode               = "0x60806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630c55699c14610072578063371303c01461009d57806350cd4df2146100b4578063579be378146100df578063d3ab86a1146100f6575b600080fd5b34801561007e57600080fd5b5061008761014d565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b2610153565b005b3480156100c057600080fd5b506100c96101ba565b6040518082815260200191505060405180910390f35b3480156100eb57600080fd5b506100f46101c0565b005b34801561010257600080fd5b5061010b6101d9565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60005481565b33600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160008154809291906001900391905055506000808154809291906001019190505550565b60015481565b3373ffffffffffffffffffffffffffffffffffffffff16ff5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820b537b2bbcf121c2be169c4f990888d02d3bbab4fd6a806c3d4a0f3643cebd4590029"
	blockedAcc, _             = sdk.AccAddressFromBech32(blockedAddr)
	blockedBuffer             = hexutil.MustDecode(blockedCode)
	blockedEthAcc             = common.BytesToAddress(blockedAcc.Bytes())
	callMethodBlocked         = hexutil.MustDecode("0xa9421619000000000000000000000000f297ab486be410a2649901849b0477d519e99960")
	selfdestructMethodBlocked = hexutil.MustDecode("0xa6516bda000000000000000000000000f297ab486be410a2649901849b0477d519e99960")
	callcodeMethodBlocked     = hexutil.MustDecode("0x7811c6c1000000000000000000000000f297ab486be410a2649901849b0477d519e99960")
	delegatecallMethodBlocked = hexutil.MustDecode("0xa7126c2d000000000000000000000000f297ab486be410a2649901849b0477d519e99960")
	blockedMethods            = types.ContractMethods{
		types.ContractMethod{
			Sign:  "0x371303c0",
			Extra: "inc()",
		},
	}
	blockedContract = types.BlockedContract{
		Address:      blockedAcc,
		BlockMethods: blockedMethods,
	}
)

//Call Code ABI
//[
//  {
//    "constant": false,
//    "inputs": [],
//    "name": "inc",
//    "outputs": [],
//    "payable": false,
//    "stateMutability": "nonpayable",
//    "type": "function"
//  },
//  {
//    "constant": false,
//    "inputs": [],
//    "name": "onc",
//    "outputs": [],
//    "payable": false,
//    "stateMutability": "nonpayable",
//    "type": "function"
//  },
//]
//Blocked Code ABI
//[
//  {
//    "constant": false,
//    "inputs": [
//      {
//        "name": "contractAddress",
//        "type": "address"
//      }
//    ],
//    "name": "inc_call",
//    "outputs": [],
//    "payable": false,
//    "stateMutability": "nonpayable",
//    "type": "function"
//  },
//  {
//    "constant": false,
//    "inputs": [
//      {
//        "name": "contractAddress",
//        "type": "address"
//      }
//    ],
//    "name": "inc_call_selfdestruct",
//    "outputs": [],
//    "payable": false,
//    "stateMutability": "nonpayable",
//    "type": "function"
//  },
//  {
//    "constant": false,
//    "inputs": [
//      {
//        "name": "contractAddress",
//        "type": "address"
//      }
//    ],
//    "name": "inc_callcode",
//    "outputs": [],
//    "payable": false,
//    "stateMutability": "nonpayable",
//    "type": "function"
//  },
//  {
//    "constant": false,
//    "inputs": [
//      {
//        "name": "contractAddress",
//        "type": "address"
//      }
//    ],
//    "name": "inc_delegatecall",
//    "outputs": [],
//    "payable": false,
//    "stateMutability": "nonpayable",
//    "type": "function"
//  },
//]
func (suite *StateDBTestSuite) TestGetHashFn() {
	testCase := []struct {
		name         string
		height       uint64
		malleate     func()
		expEmptyHash bool
	}{
		{
			"valid hash, case 1",
			1,
			func() {
				suite.ctx = suite.ctx.WithBlockHeader(
					abci.Header{
						ChainID:        "ethermint-1",
						Height:         1,
						ValidatorsHash: []byte("val_hash"),
					},
				)
				hash := ethcmn.BytesToHash([]byte("test hash"))
				suite.stateDB.SetBlockHash(hash)
			},
			false,
		},
		{
			"case 1, nil tendermint hash",
			1,
			func() {},
			true,
		},
		{
			"valid hash, case 2",
			1,
			func() {
				suite.ctx = suite.ctx.WithBlockHeader(
					abci.Header{
						ChainID:        "ethermint-1",
						Height:         100,
						ValidatorsHash: []byte("val_hash"),
					},
				)
				hash := ethcmn.BytesToHash([]byte("test hash"))
				suite.stateDB.WithContext(suite.ctx).SetHeightHash(1, hash)
			},
			false,
		},
		{
			"height not found, case 2",
			1,
			func() {
				suite.ctx = suite.ctx.WithBlockHeader(
					abci.Header{
						ChainID:        "ethermint-1",
						Height:         100,
						ValidatorsHash: []byte("val_hash"),
					},
				)
			},
			true,
		},
		{
			"empty hash, case 3",
			1000,
			func() {
				suite.ctx = suite.ctx.WithBlockHeader(
					abci.Header{
						ChainID:        "ethermint-1",
						Height:         100,
						ValidatorsHash: []byte("val_hash"),
					},
				)
			},
			true,
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			tc.malleate()

			hash := types.GetHashFn(suite.ctx, suite.stateDB)(tc.height)
			if tc.expEmptyHash {
				suite.Require().Equal(common.Hash{}.String(), hash.String())
			} else {
				suite.Require().NotEqual(common.Hash{}.String(), hash.String())
			}
		})
	}
}

func (suite *StateDBTestSuite) TestTransitionDb() {
	suite.stateDB.SetNonce(suite.address, 123)

	addr := sdk.AccAddress(suite.address.Bytes())
	balance := ethermint.NewPhotonCoin(sdk.NewInt(5000))
	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, addr)
	_ = acc.SetCoins(sdk.NewCoins(balance))
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	suite.stateDB = types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	recipient := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	testCase := []struct {
		name     string
		malleate func()
		state    types.StateTransition
		expPass  bool
	}{
		{
			"passing state transition",
			func() {},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     11,
				Recipient:    &recipient,
				Amount:       sdk.NewDec(50).BigInt(),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			true,
		},
		{
			"contract creation",
			func() {},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     11,
				Recipient:    nil,
				Amount:       sdk.NewDec(10).BigInt(),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			true,
		},
		{
			"fail by sending more than balance",
			func() {},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     11,
				Recipient:    &recipient,
				Amount:       sdk.NewDec(500000).BigInt(),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"call disabled",
			func() {
				params := types.NewParams(true, false, false, false, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        big.NewInt(10),
				GasLimit:     11,
				Recipient:    &recipient,
				Amount:       big.NewInt(50),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"create disabled",
			func() {
				params := types.NewParams(false, true, false, false, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        big.NewInt(10),
				GasLimit:     11,
				Recipient:    nil,
				Amount:       big.NewInt(50),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"nil gas price",
			func() {
				suite.stateDB.SetParams(types.DefaultParams())
				invalidGas := sdk.DecCoins{
					{Denom: ethermint.NativeToken},
				}
				suite.ctx = suite.ctx.WithMinGasPrices(invalidGas)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     11,
				Recipient:    &recipient,
				Amount:       sdk.NewDec(10).BigInt(),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"state transition simulation",
			func() {
				params := types.NewParams(false, true, false, false, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     11,
				Recipient:    &recipient,
				Amount:       sdk.NewDec(10).BigInt(),
				Payload:      []byte("data"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     true,
			},
			true,
		},
		{
			"contract failed call addr is blocked",
			func() {
				params := types.NewParams(false, true, false, true, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)

				suite.stateDB.SetCode(common.BytesToAddress(callAcc.Bytes()), callBuffer)
				suite.stateDB.SetCode(common.BytesToAddress(blockedAcc.Bytes()), blockedBuffer)
				blockedList := types.AddressList{blockedAcc}
				suite.stateDB.SetContractBlockedList(blockedList)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     100000000,
				Recipient:    &blockedEthAcc,
				Amount:       sdk.NewDec(0).BigInt(),
				Payload:      hexutil.MustDecode("0x371303c0"),
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"contract failed call contract method blocked",
			func() {
				params := types.NewParams(false, true, false, true, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)
				suite.stateDB.SetCode(common.BytesToAddress(callAcc.Bytes()), callBuffer)
				suite.stateDB.SetCode(common.BytesToAddress(blockedAcc.Bytes()), blockedBuffer)
				suite.stateDB.SetContractMethodBlocked(blockedContract)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     100000000,
				Recipient:    &callEthAcc,
				Amount:       sdk.NewDec(0).BigInt(),
				Payload:      callMethodBlocked,
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"contract failed callcode contract method blocked",
			func() {
				params := types.NewParams(false, true, false, true, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)
				suite.stateDB.SetCode(common.BytesToAddress(callAcc.Bytes()), callBuffer)
				suite.stateDB.SetCode(common.BytesToAddress(blockedAcc.Bytes()), blockedBuffer)
				suite.stateDB.SetContractMethodBlocked(blockedContract)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     100000000,
				Recipient:    &callEthAcc,
				Amount:       sdk.NewDec(0).BigInt(),
				Payload:      callcodeMethodBlocked,
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"contract failed delegate call contract method blocked",
			func() {
				params := types.NewParams(false, true, false, true, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)
				suite.stateDB.SetCode(common.BytesToAddress(callAcc.Bytes()), callBuffer)
				suite.stateDB.SetCode(common.BytesToAddress(blockedAcc.Bytes()), blockedBuffer)
				suite.stateDB.SetContractMethodBlocked(blockedContract)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     100000000,
				Recipient:    &callEthAcc,
				Amount:       sdk.NewDec(0).BigInt(),
				Payload:      delegatecallMethodBlocked,
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
		{
			"contract failed selfdestruct contract method blocked",
			func() {
				params := types.NewParams(false, true, false, true, types.DefaultMaxGasLimitPerTx)
				suite.stateDB.SetParams(params)

				suite.stateDB.CreateAccount(callEthAcc)
				suite.stateDB.CreateAccount(blockedEthAcc)
				suite.stateDB.SetCode(callEthAcc, callBuffer)
				suite.stateDB.SetCode(blockedEthAcc, blockedBuffer)

				suite.stateDB.SetContractMethodBlocked(blockedContract)
			},
			types.StateTransition{
				AccountNonce: 123,
				Price:        sdk.NewDec(10).BigInt(),
				GasLimit:     100000000,
				Recipient:    &callEthAcc,
				Amount:       sdk.NewDec(0).BigInt(),
				Payload:      selfdestructMethodBlocked,
				ChainID:      big.NewInt(1),
				Csdb:         suite.stateDB,
				TxHash:       &ethcmn.Hash{},
				Sender:       suite.address,
				Simulate:     suite.ctx.IsCheckTx(),
			},
			false,
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		cc := types.DefaultChainConfig()
		_, _, err, _, _ = tc.state.TransitionDb(suite.ctx, &cc)

		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
	fromBalance := suite.app.EvmKeeper.GetBalance(suite.ctx, suite.address)
	toBalance := suite.app.EvmKeeper.GetBalance(suite.ctx, recipient)
	suite.Require().Equal(fromBalance, sdk.NewDec(4940).BigInt())
	suite.Require().Equal(toBalance, sdk.NewDec(50).BigInt())
}
