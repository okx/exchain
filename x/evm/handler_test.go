package evm_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/stretchr/testify/suite"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"

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
	codec   *codec.Codec
}

func (suite *EvmTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
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

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(30000000)))

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
				suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)
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
				suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

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
			suite.ctx = suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			res, err := suite.handler(suite.ctx, tx)

			//nolint
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				var expectedConsumedGas uint64 = 21000
				suite.Require().EqualValues(expectedConsumedGas, suite.ctx.GasMeter().GasConsumed())
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

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(30000000)))

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)
				tx = types.NewMsgEthermint(0, &to, sdk.NewInt(1), 100000, sdk.NewInt(2), []byte{}, from)
				suite.app.EvmKeeper.SetBalance(suite.ctx, ethcmn.BytesToAddress(from.Bytes()), big.NewInt(100))
			},
			true,
		},
		{
			"invalid state transition",
			func() {
				suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)
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
			suite.ctx = suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			res, err := suite.handler(suite.ctx, tx)

			//nolint
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				var expectedConsumedGas uint64 = 21000
				suite.Require().EqualValues(expectedConsumedGas, suite.ctx.GasMeter().GasConsumed())
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

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(3000000000000)))
	suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

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
	err = suite.app.EvmKeeper.SetLogs(suite.ctx, ethcmn.BytesToHash(hash), resultData.Logs)
	suite.Require().NoError(err)

	logs, err := suite.app.EvmKeeper.GetLogs(suite.ctx, ethcmn.BytesToHash(hash))
	suite.Require().NoError(err, "failed to get logs")

	suite.Require().Equal(logs, resultData.Logs)
}

func (suite *EvmTestSuite) TestQueryTxLogs() {

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(1000000000000)))
	suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(1000000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	// send contract deployment transaction with an event in the constructor
	bytecode := common.FromHex("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	err = tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	suite.Require().Equal(len(resultData.Logs), 1)
	suite.Require().Equal(len(resultData.Logs[0].Topics), 2)

	// get logs by tx hash
	hash := resultData.TxHash.Bytes()

	logs, err := suite.app.EvmKeeper.GetLogs(suite.ctx, ethcmn.BytesToHash(hash))
	suite.Require().NoError(err, "failed to get logs")

	suite.Require().Equal(logs, resultData.Logs)

	// query tx logs
	path := []string{"transactionLogs", fmt.Sprintf("0x%x", hash)}
	res, err := suite.querier(suite.ctx, path, abci.RequestQuery{})
	suite.Require().NoError(err, "failed to query txLogs")

	var txLogs types.QueryETHLogs
	suite.codec.MustUnmarshalJSON(res, &txLogs)

	// amino decodes an empty byte array as nil, whereas JSON decodes it as []byte{} causing a discrepancy
	resultData.Logs[0].Data = []byte{}
	suite.Require().Equal(txLogs.Logs[0], resultData.Logs[0])
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

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(1000000000000000000)))
	suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

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

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(3000000000000)))
	suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

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

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
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

	snapshotCommitStateDBJson, err := json.Marshal(suite.app.EvmKeeper.CommitStateDB)
	suite.Require().Nil(err)

	defer func() {
		if r := recover(); r != nil {
			currentCommitStateDBJson, err := json.Marshal(suite.app.EvmKeeper.CommitStateDB)
			suite.Require().Nil(err)
			suite.Require().Equal(snapshotCommitStateDBJson, currentCommitStateDBJson)
		} else {
			suite.Require().Fail("panic did not happen")
		}
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
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(1000000000000000000)))
	suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

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

func (suite *EvmTestSuite) TestErrorWhenDeployContract() {
	gasLimit := uint64(1000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424")

	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())
	suite.Require().NoError(err)

	snapshotCommitStateDBJson, err := json.Marshal(suite.app.EvmKeeper.CommitStateDB)
	suite.Require().Nil(err)

	_, sdkErr := suite.handler(suite.ctx, tx)
	suite.Require().NotNil(sdkErr)

	currentCommitStateDBJson, err := json.Marshal(suite.app.EvmKeeper.CommitStateDB)
	suite.Require().Nil(err)
	suite.Require().Equal(snapshotCommitStateDBJson, currentCommitStateDBJson)
}

func (suite *EvmTestSuite) TestDefaultMsgHandler() {
	tx := sdk.NewTestMsg()
	_, sdkErr := suite.handler(suite.ctx, tx)
	suite.Require().NotNil(sdkErr)
}

func (suite *EvmTestSuite) TestRefundGas() {
	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	sender := ethcmn.HexToAddress(privkey.PubKey().Address().String())
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	feeCollectorAcc.Coins = sdk.NewCoins(sdk.NewCoin(suite.app.EvmKeeper.GetParams(suite.ctx).EvmDenom, sdk.NewInt(1000000000)))

	var tx types.MsgEthereumTx

	gasLimit := uint64(10000000)
	gasPrice := big.NewInt(1)
	userBalance := big.NewInt(100)

	testCases := []struct {
		msg      string
		malleate func()
		consumed *big.Int
	}{
		{
			"Refund on successful transaction",
			func() {
				suite.app.EvmKeeper.SetBalance(suite.ctx, sender, userBalance)
				suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

				tx = types.NewMsgEthereumTx(0, &sender, big.NewInt(100), gasLimit, gasPrice, nil)

				// parse context chain ID to big.Int
				chainID, err := ethermint.ParseChainID(suite.ctx.ChainID())
				suite.Require().NoError(err)

				// sign transaction
				err = tx.Sign(chainID, privkey.ToECDSA())
				suite.Require().NoError(err)
			},
			big.NewInt(0),
		},
		{
			"Refund on failure of transactionDB execution",
			func() {
				suite.app.EvmKeeper.SetBalance(suite.ctx, sender, userBalance)
				suite.app.SupplyKeeper.SetModuleAccount(suite.ctx, feeCollectorAcc)

				bytecode := common.FromHex("0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424")
				tx = types.NewMsgEthereumTx(0, nil, big.NewInt(100), gasLimit, gasPrice, bytecode)

				// parse context chain ID to big.Int
				chainID, err := ethermint.ParseChainID(suite.ctx.ChainID())
				suite.Require().NoError(err)

				err = tx.Sign(chainID, privkey.ToECDSA())
				suite.Require().NoError(err)
			},
			big.NewInt(0),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			suite.SetupTest() // reset
			//nolint
			tc.malleate()

			suite.handler(suite.ctx, tx)

			gasLeft := big.NewInt(1).SetUint64(gasLimit - suite.ctx.GasMeter().GasConsumed())
			gasRefund := big.NewInt(1).Mul(gasPrice, gasLeft)

			balanceAfterHandler := big.NewInt(1).Sub(userBalance, tc.consumed)
			balanceAfterRefund := suite.app.EvmKeeper.GetBalance(suite.ctx, sender)

			suite.Require().Equal(big.NewInt(1).Mul(gasRefund, big.NewInt(1)), big.NewInt(1).Sub(balanceAfterRefund, balanceAfterHandler))
		})
	}
}
