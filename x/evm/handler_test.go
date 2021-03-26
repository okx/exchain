package evm_test

import (
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/app"
	"github.com/okex/okexchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm"
	"github.com/okex/okexchain/x/evm/keeper"
	"github.com/okex/okexchain/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type EvmTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	handler sdk.Handler
	querier sdk.Querier
	app     *app.OKExChainApp
	stateDB *types.CommitStateDB
	codec   *codec.Codec
}

func (suite *EvmTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
	suite.stateDB = types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)
	suite.handler = evm.NewHandler(suite.app.EvmKeeper)
	suite.querier = keeper.NewQuerier(*suite.app.EvmKeeper)
	suite.codec = codec.New()

	params := types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	suite.app.EvmKeeper.SetParams(suite.ctx, params)
}

func TestEvmTestSuite(t *testing.T) {
	suite.Run(t, new(EvmTestSuite))
}

func (suite *EvmTestSuite) TestHandleMsgEthereumTx() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	sender := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	var tx types.MsgEthereumTx

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				suite.app.EvmKeeper.SetBalance(suite.ctx, sender, big.NewInt(100))
				tx = types.NewMsgEthereumTx(0, &sender, big.NewInt(100), 3000000, big.NewInt(1), nil)

				// parse context chain ID to big.Int
				chainID, err := ethermint.ParseChainID(suite.ctx.ChainID())
				suite.Require().NoError(err)

				// sign transaction
				err = tx.Sign(chainID, privkey.ToECDSA())
				suite.Require().NoError(err)
			},
			true,
		},
		{
			"insufficient balance",
			func() {
				suite.app.EvmKeeper.SetBalance(suite.ctx, sender, big.NewInt(1))
				tx = types.NewMsgEthereumTxContract(0, big.NewInt(100), 3000000, big.NewInt(1), nil)

				// parse context chain ID to big.Int
				chainID, err := ethermint.ParseChainID(suite.ctx.ChainID())
				suite.Require().NoError(err)

				// sign transaction
				err = tx.Sign(chainID, privkey.ToECDSA())
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"tx encoding failed",
			func() {
				tx = types.NewMsgEthereumTxContract(0, big.NewInt(100), 3000000, big.NewInt(1), nil)
			},
			false,
		},
		{
			"invalid chain ID",
			func() {
				suite.ctx = suite.ctx.WithChainID("chainID")
			},
			false,
		},
		{
			"VerifySig failed",
			func() {
				tx = types.NewMsgEthereumTxContract(0, big.NewInt(100), 3000000, big.NewInt(1), nil)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			suite.SetupTest() // reset
			//nolint
			tc.malleate()

			res, err := suite.handler(suite.ctx, tx)

			//nolint
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *EvmTestSuite) TestMsgEthermint() {
	var (
		tx   types.MsgEthermint
		from = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		to   = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				tx = types.NewMsgEthermint(0, &to, sdk.NewInt(1), 100000, sdk.NewInt(2), []byte("test"), from)
				suite.app.EvmKeeper.SetBalance(suite.ctx, ethcmn.BytesToAddress(from.Bytes()), big.NewInt(100))
			},
			true,
		},
		{
			"invalid state transition",
			func() {
				tx = types.NewMsgEthermint(0, &to, sdk.NewInt(1), 100000, sdk.NewInt(2), []byte("test"), from)
			},
			false,
		},
		{
			"invalid chain ID",
			func() {
				suite.ctx = suite.ctx.WithChainID("chainID")
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run("", func() {
			suite.SetupTest() // reset
			//nolint
			tc.malleate()

			res, err := suite.handler(suite.ctx, tx)

			//nolint
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *EvmTestSuite) TestHandlerLogs() {
	// Test contract:

	// pragma solidity ^0.5.1;

	// contract Test {
	//     event Hello(uint256 indexed world);

	//     constructor() public {
	//         emit Hello(17);
	//     }
	// }

	// {
	// 	"linkReferences": {},
	// 	"object": "6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029",
	// 	"opcodes": "PUSH1 0x80 PUSH1 0x40 MSTORE CALLVALUE DUP1 ISZERO PUSH1 0xF JUMPI PUSH1 0x0 DUP1 REVERT JUMPDEST POP PUSH1 0x11 PUSH32 0x775A94827B8FD9B519D36CD827093C664F93347070A554F65E4A6F56CD738898 PUSH1 0x40 MLOAD PUSH1 0x40 MLOAD DUP1 SWAP2 SUB SWAP1 LOG2 PUSH1 0x35 DUP1 PUSH1 0x4B PUSH1 0x0 CODECOPY PUSH1 0x0 RETURN INVALID PUSH1 0x80 PUSH1 0x40 MSTORE PUSH1 0x0 DUP1 REVERT INVALID LOG1 PUSH6 0x627A7A723058 KECCAK256 PUSH13 0xAB665F0F557620554BB45ADF26 PUSH8 0x8D2BD349B8A4314 0xbd SELFDESTRUCT KECCAK256 0x5e 0xe8 DIFFICULTY 0xe EXTCODECOPY 0x24 STOP 0x29 ",
	// 	"sourceMap": "25:119:0:-;;;90:52;8:9:-1;5:2;;;30:1;27;20:12;5:2;90:52:0;132:2;126:9;;;;;;;;;;25:119;;;;;;"
	// }

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(1000000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	err = tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	suite.Require().Equal(len(resultData.Logs), 1)
	suite.Require().Equal(len(resultData.Logs[0].Topics), 2)

	hash := []byte{1}
	err = suite.stateDB.WithContext(suite.ctx).SetLogs(ethcmn.BytesToHash(hash), resultData.Logs)
	suite.Require().NoError(err)

	logs, err := suite.stateDB.WithContext(suite.ctx).GetLogs(ethcmn.BytesToHash(hash))
	suite.Require().NoError(err, "failed to get logs")

	suite.Require().Equal(logs, resultData.Logs)
}

func (suite *EvmTestSuite) TestDeployAndCallContract() {
	// Test contract:
	//http://remix.ethereum.org/#optimize=false&evmVersion=istanbul&version=soljson-v0.5.15+commit.6a57276f.js
	//2_Owner.sol
	//
	//pragma solidity >=0.4.22 <0.7.0;
	//
	///**
	// * @title Owner
	// * @dev Set & change owner
	// */
	//contract Owner {
	//
	//	address private owner;
	//
	//	// event for EVM logging
	//	event OwnerSet(address indexed oldOwner, address indexed newOwner);
	//
	//	// modifier to check if caller is owner
	//	modifier isOwner() {
	//	// If the first argument of 'require' evaluates to 'false', execution terminates and all
	//	// changes to the state and to Ether balances are reverted.
	//	// This used to consume all gas in old EVM versions, but not anymore.
	//	// It is often a good idea to use 'require' to check if functions are called correctly.
	//	// As a second argument, you can also provide an explanation about what went wrong.
	//	require(msg.sender == owner, "Caller is not owner");
	//	_;
	//}
	//
	//	/**
	//	 * @dev Set contract deployer as owner
	//	 */
	//	constructor() public {
	//	owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
	//	emit OwnerSet(address(0), owner);
	//}
	//
	//	/**
	//	 * @dev Change owner
	//	 * @param newOwner address of new owner
	//	 */
	//	function changeOwner(address newOwner) public isOwner {
	//	emit OwnerSet(owner, newOwner);
	//	owner = newOwner;
	//}
	//
	//	/**
	//	 * @dev Return owner address
	//	 * @return address of owner
	//	 */
	//	function getOwner() external view returns (address) {
	//	return owner;
	//}
	//}

	// Deploy contract - Owner.sol
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	// store - changeOwner
	gasLimit = uint64(100000000000)
	gasPrice = big.NewInt(100)
	receiver := common.HexToAddress(resultData.ContractAddress.String())

	storeAddr := "0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424"
	bytecode = common.FromHex(storeAddr)
	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err = suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err = types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	// query - getOwner
	bytecode = common.FromHex("0x893d20e8")
	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err = suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err = types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	getAddr := strings.ToLower(hexutils.BytesToHex(resultData.Ret))
	suite.Require().Equal(true, strings.HasSuffix(storeAddr, getAddr), "Fail to query the address")
}

func (suite *EvmTestSuite) TestSendTransaction() {

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")
	pub := priv.ToECDSA().Public().(*ecdsa.PublicKey)

	suite.app.EvmKeeper.SetBalance(suite.ctx, ethcrypto.PubkeyToAddress(*pub), big.NewInt(100))

	// send simple value transfer with gasLimit=21000
	tx := types.NewMsgEthereumTx(1, &ethcmn.Address{0x1}, big.NewInt(1), gasLimit, gasPrice, nil)
	err = tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	suite.ctx = suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	var expectedGas uint64 = 1336
	suite.Require().EqualValues(expectedGas, suite.ctx.GasMeter().GasConsumed())
}

func (suite *EvmTestSuite) TestOutOfGasWhenDeployContract() {
	// Test contract:
	//http://remix.ethereum.org/#optimize=false&evmVersion=istanbul&version=soljson-v0.5.15+commit.6a57276f.js
	//2_Owner.sol
	//
	//pragma solidity >=0.4.22 <0.7.0;
	//
	///**
	// * @title Owner
	// * @dev Set & change owner
	// */
	//contract Owner {
	//
	//	address private owner;
	//
	//	// event for EVM logging
	//	event OwnerSet(address indexed oldOwner, address indexed newOwner);
	//
	//	// modifier to check if caller is owner
	//	modifier isOwner() {
	//	// If the first argument of 'require' evaluates to 'false', execution terminates and all
	//	// changes to the state and to Ether balances are reverted.
	//	// This used to consume all gas in old EVM versions, but not anymore.
	//	// It is often a good idea to use 'require' to check if functions are called correctly.
	//	// As a second argument, you can also provide an explanation about what went wrong.
	//	require(msg.sender == owner, "Caller is not owner");
	//	_;
	//}
	//
	//	/**
	//	 * @dev Set contract deployer as owner
	//	 */
	//	constructor() public {
	//	owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
	//	emit OwnerSet(address(0), owner);
	//}
	//
	//	/**
	//	 * @dev Change owner
	//	 * @param newOwner address of new owner
	//	 */
	//	function changeOwner(address newOwner) public isOwner {
	//	emit OwnerSet(owner, newOwner);
	//	owner = newOwner;
	//}
	//
	//	/**
	//	 * @dev Return owner address
	//	 * @return address of owner
	//	 */
	//	function getOwner() external view returns (address) {
	//	return owner;
	//}
	//}

	// Deploy contract - Owner.sol
	gasLimit := uint64(1)
	suite.ctx = suite.ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	defer func() {
		r := recover()
		suite.Require().NotNil(r, "panic for out of gas")
	}()

	suite.handler(suite.ctx, tx)
	suite.Require().Fail("panic did not happen")
}

func (suite *EvmTestSuite) TestRevertErrorWhenCallContract() {
	// Test contract:

	//// SPDX-License-Identifier: GPL-3.0
	//
	//pragma solidity >=0.7.0 <0.8.0;
	//
	///**
	// * @title Storage
	// * @dev Store & retrieve value in a variable
	// */
	//contract Storage {
	//
	//	uint256 number;
	//	event Test(address to);
	//
	//	/**
	//	 * @dev Store value in variable
	//	 * @param num value to store
	//	 */
	//	function store(uint256 num) public {
	//	require(false,"this is my test failed message");
	//	number = num;
	//	emit Test(msg.sender);
	//}
	//
	//	/**
	//	 * @dev Return value
	//	 * @return value of 'number'
	//	 */
	//	function retrieve() public view returns (uint256){
	//	return number;
	//}
	//}

	// Deploy contract - storage.sol
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0x608060405234801561001057600080fd5b50610191806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610087565b6040518082815260200191505060405180910390f35b6100856004803603602081101561006f57600080fd5b8101908080359060200190929190505050610090565b005b60008054905090565b6000610104576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f74686973206973206d792074657374206661696c6564206d657373616765000081525060200191505060405180910390fd5b806000819055507faa9449f2bca09a7b28319d46fd3f3b58a1bb7d94039fc4b69b7bfe5d8535d52733604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a15056fea264697066735822122078908b7dd6de7f67bccf9fa221c027590325c5df3cd7d654ee654e4834ca952b64736f6c63430007060033")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	// store - changeOwner
	gasLimit = uint64(100000000000)
	gasPrice = big.NewInt(100)
	receiver := common.HexToAddress(resultData.ContractAddress.String())

	storeAddr := "0x6057361d0000000000000000000000000000000000000000000000000000000000000001"
	bytecode = common.FromHex(storeAddr)
	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err = suite.handler(suite.ctx, tx)
	suite.Require().Nil(result)
	suite.Require().NotNil(err)
	suite.Require().Equal(err.Error(), "[\"execution reverted\",\"execution reverted:this is my test failed message\",\"HexData\",\"0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001e74686973206973206d792074657374206661696c6564206d6573736167650000\"]")
}

func (suite *EvmTestSuite) TestGasConsume() {
	// Test contract:
	//
	//pragma solidity ^0.8.0;
	//contract Test {
	//	event NotifyUint(string funcName, uint value);
	//	event NotifyBytes32(string funcName, bytes32 value);
	//	event NotifyAddress(string funcName, address value);
	//	event NotifyUint256(string funcName, uint256 value);
	//	event NotifyBytes(string funcName, bytes value);
	//	event NotifyBytes4(string funcName, bytes4 value);
	//
	//	function rand() public payable{
	//	// block releted
	//	emit NotifyUint("block.difficulty", uint(block.difficulty));
	//	emit NotifyUint("block.gaslimit", uint(block.gaslimit));
	//
	//	// not work until solidity v0.8.0
	//	emit NotifyUint("block.chainid", uint(block.chainid));
	//
	//	uint num;
	//	num = uint(block.number);
	//	emit NotifyUint("block.number", num);
	//	emit NotifyBytes32("blockhash", bytes32(blockhash(num)));
	//	emit NotifyBytes32("last blockhash", bytes32(blockhash(num - 1)));
	//	emit NotifyAddress("block.coinbase", address(block.coinbase));
	//
	//	emit NotifyUint("block.timestamp", uint(block.timestamp));
	//	// not work since solidity v0.7.0
	//	//emit NotifyUint("now", uint(now));
	//
	//
	//	// msg releted
	//	emit NotifyBytes("msg.data", bytes(msg.data));
	//	emit NotifyAddress("msg.sender", address(msg.sender));
	//	emit NotifyBytes4("msg.sig", bytes4(msg.sig));
	//	emit NotifyUint("msg.value", uint(msg.value));
	//	emit NotifyUint256("gasleft", uint256(gasleft()));
	//
	//
	//	// tx releted
	//	emit NotifyUint("tx.gasprice", uint(tx.gasprice));
	//	emit NotifyAddress("tx.origin", address(tx.origin));
	//}
	//}

	// Deploy contract - storage.sol
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("608060405234801561001057600080fd5b50610c88806100206000396000f3fe60806040526004361061001e5760003560e01c80633b3dca7614610023575b600080fd5b61002b61002d565b005b7f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f4460405161005c91906107ee565b60405180910390a17f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f4560405161009391906108a6565b60405180910390a17f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f466040516100ca9190610619565b60405180910390a160004390507f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f8160405161010691906106da565b60405180910390a17f66f9aff21094d330be5210a94aa9006bc73c58a78ffbc2fcf7459b469e90a91c814060405161013e919061081c565b60405180910390a17f66f9aff21094d330be5210a94aa9006bc73c58a78ffbc2fcf7459b469e90a91c60018261017491906108f6565b4060405161018291906107c0565b60405180910390a17fd1f988f2f9743dc861b689b05592f3f5a3c320f9ac7c058e7f1aed7c5302c0aa416040516101b99190610878565b60405180910390a17f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f426040516101f09190610647565b60405180910390a17fc5eea5f168bedf3b83d910a0e351deae24dc0af984734c8d3ad930ddb385967c60003660405161022a9291906106a3565b60405180910390a17fd1f988f2f9743dc861b689b05592f3f5a3c320f9ac7c058e7f1aed7c5302c0aa336040516102619190610792565b60405180910390a17f06cfc4b8d6699a3e735d3478223775d7e9ba040403a71b442bca541d0e1e94bf6000357fffffffff00000000000000000000000000000000000000000000000000000000166040516102bc9190610675565b60405180910390a17f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f346040516102f39190610708565b60405180910390a17f428932ecd90e05b8cebd39dcda376284a2ccf362d7081b287d54fbc07b4777265a60405161032a9190610764565b60405180910390a17f01710aa49abe19e9f1e69db68002b08e6fdc26b6603462ec5595143eb16bce0f3a604051610361919061084a565b60405180910390a17fd1f988f2f9743dc861b689b05592f3f5a3c320f9ac7c058e7f1aed7c5302c0aa326040516103989190610736565b60405180910390a150565b6103ac8161092a565b82525050565b6103bb8161093c565b82525050565b6103ca81610946565b82525050565b60006103dc83856108d4565b93506103e983858461099c565b6103f2836109da565b840190509392505050565b600061040a600d836108e5565b9150610415826109eb565b602082019050919050565b600061042d600f836108e5565b915061043882610a14565b602082019050919050565b60006104506007836108e5565b915061045b82610a3d565b602082019050919050565b60006104736008836108e5565b915061047e82610a66565b602082019050919050565b6000610496600c836108e5565b91506104a182610a8f565b602082019050919050565b60006104b96009836108e5565b91506104c482610ab8565b602082019050919050565b60006104dc6009836108e5565b91506104e782610ae1565b602082019050919050565b60006104ff6007836108e5565b915061050a82610b0a565b602082019050919050565b6000610522600a836108e5565b915061052d82610b33565b602082019050919050565b6000610545600e836108e5565b915061055082610b5c565b602082019050919050565b60006105686010836108e5565b915061057382610b85565b602082019050919050565b600061058b6009836108e5565b915061059682610bae565b602082019050919050565b60006105ae600b836108e5565b91506105b982610bd7565b602082019050919050565b60006105d1600e836108e5565b91506105dc82610c00565b602082019050919050565b60006105f4600e836108e5565b91506105ff82610c29565b602082019050919050565b61061381610992565b82525050565b60006040820190508181036000830152610632816103fd565b9050610641602083018461060a565b92915050565b6000604082019050818103600083015261066081610420565b905061066f602083018461060a565b92915050565b6000604082019050818103600083015261068e81610443565b905061069d60208301846103c1565b92915050565b600060408201905081810360008301526106bc81610466565b905081810360208301526106d18184866103d0565b90509392505050565b600060408201905081810360008301526106f381610489565b9050610702602083018461060a565b92915050565b60006040820190508181036000830152610721816104ac565b9050610730602083018461060a565b92915050565b6000604082019050818103600083015261074f816104cf565b905061075e60208301846103a3565b92915050565b6000604082019050818103600083015261077d816104f2565b905061078c602083018461060a565b92915050565b600060408201905081810360008301526107ab81610515565b90506107ba60208301846103a3565b92915050565b600060408201905081810360008301526107d981610538565b90506107e860208301846103b2565b92915050565b600060408201905081810360008301526108078161055b565b9050610816602083018461060a565b92915050565b600060408201905081810360008301526108358161057e565b905061084460208301846103b2565b92915050565b60006040820190508181036000830152610863816105a1565b9050610872602083018461060a565b92915050565b60006040820190508181036000830152610891816105c4565b90506108a060208301846103a3565b92915050565b600060408201905081810360008301526108bf816105e7565b90506108ce602083018461060a565b92915050565b600082825260208201905092915050565b600082825260208201905092915050565b600061090182610992565b915061090c83610992565b92508282101561091f5761091e6109ab565b5b828203905092915050565b600061093582610972565b9050919050565b6000819050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000601f19601f8301169050919050565b7f626c6f636b2e636861696e696400000000000000000000000000000000000000600082015250565b7f626c6f636b2e74696d657374616d700000000000000000000000000000000000600082015250565b7f6d73672e73696700000000000000000000000000000000000000000000000000600082015250565b7f6d73672e64617461000000000000000000000000000000000000000000000000600082015250565b7f626c6f636b2e6e756d6265720000000000000000000000000000000000000000600082015250565b7f6d73672e76616c75650000000000000000000000000000000000000000000000600082015250565b7f74782e6f726967696e0000000000000000000000000000000000000000000000600082015250565b7f6761736c65667400000000000000000000000000000000000000000000000000600082015250565b7f6d73672e73656e64657200000000000000000000000000000000000000000000600082015250565b7f6c61737420626c6f636b68617368000000000000000000000000000000000000600082015250565b7f626c6f636b2e646966666963756c747900000000000000000000000000000000600082015250565b7f626c6f636b686173680000000000000000000000000000000000000000000000600082015250565b7f74782e6761737072696365000000000000000000000000000000000000000000600082015250565b7f626c6f636b2e636f696e62617365000000000000000000000000000000000000600082015250565b7f626c6f636b2e6761736c696d697400000000000000000000000000000000000060008201525056fea264697066735822122002530b5caa31bb9ea0bb1b97d1d7d7225a527bfdedadb22ac0f809fd2d21bed064736f6c63430008010033")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	_, err = suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")
	var expectedConsumedGas sdk.Gas = 660122
	suite.Require().Equal(expectedConsumedGas, suite.ctx.GasMeter().GasConsumed())
}

func (suite *EvmTestSuite) TestErrorWhenDeployContract() {
	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424")

	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	_, sdkErr := suite.handler(suite.ctx, tx)
	suite.Require().NotNil(sdkErr)
}

func (suite *EvmTestSuite) TestDefaultMsgHandler() {
	tx := sdk.NewTestMsg()
	_, sdkErr := suite.handler(suite.ctx, tx)
	suite.Require().NotNil(sdkErr)
}

func (suite *EvmTestSuite) TestSimulateConflict() {

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")
	pub := priv.ToECDSA().Public().(*ecdsa.PublicKey)

	suite.app.EvmKeeper.SetBalance(suite.ctx, ethcrypto.PubkeyToAddress(*pub), big.NewInt(100))
	suite.stateDB.Finalise(false)

	// send simple value transfer with gasLimit=21000
	tx := types.NewMsgEthereumTx(1, &ethcmn.Address{0x1}, big.NewInt(100), gasLimit, gasPrice, nil)
	err = tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	suite.ctx = suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	suite.ctx = suite.ctx.WithIsCheckTx(true)
	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NotNil(result)
	suite.Require().Nil(err)

	suite.ctx = suite.ctx.WithIsCheckTx(false)
	result, err = suite.handler(suite.ctx, tx)
	suite.Require().NotNil(result)
	suite.Require().Nil(err)
	var expectedGas uint64 = 22336
	suite.Require().EqualValues(expectedGas, suite.ctx.GasMeter().GasConsumed())
}
