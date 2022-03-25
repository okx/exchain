package types_test

import (
	"fmt"
	"math/big"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/suite"
)

type StateDBTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *app.OKExChainApp
	stateDB     *types.CommitStateDB
	address     ethcmn.Address
	stateObject types.StateObject
}

func TestStateDBTestSuite(t *testing.T) {
	suite.Run(t, new(StateDBTestSuite))
}

func (suite *StateDBTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-1"})
	suite.ctx.SetDeliver()
	suite.stateDB = types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)

	privkey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	suite.address = ethcmn.BytesToAddress(privkey.PubKey().Address().Bytes())

	balance := sdk.NewCoins(ethermint.NewPhotonCoin(sdk.ZeroInt()))
	acc := &ethermint.EthAccount{
		BaseAccount: auth.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), balance, nil, 0, 0),
		CodeHash:    ethcrypto.Keccak256(nil),
		StateRoot: ethtypes.EmptyRootHash,
	}

	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
	suite.stateObject = suite.stateDB.GetOrNewStateObject(suite.address)
	params := types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	suite.stateDB.SetParams(params)
}

func (suite *StateDBTestSuite) TestParams() {
	defaultParams := types.DefaultParams()
	defaultParams.EnableCreate = true
	defaultParams.EnableCall = true
	params := suite.stateDB.GetParams()
	suite.Require().Equal(defaultParams, params)
	suite.stateDB.SetParams(params)
	newParams := suite.stateDB.GetParams()
	suite.Require().Equal(newParams, params)
}

func (suite *StateDBTestSuite) TestGetHeightHash() {
	hash := suite.stateDB.GetHeightHash(0)
	suite.Require().Equal(ethcmn.Hash{}.String(), hash.String())

	expHash := ethcmn.BytesToHash([]byte("hash"))
	suite.stateDB.SetHeightHash(10, expHash)

	hash = suite.stateDB.GetHeightHash(10)
	suite.Require().Equal(expHash.String(), hash.String())
}

func (suite *StateDBTestSuite) TestBloomFilter() {
	// Prepare db for logs
	tHash := ethcmn.BytesToHash([]byte{0x1})
	bhash := ethcmn.BytesToHash([]byte{0x1})
	suite.stateDB.Prepare(tHash, bhash, 0)
	contractAddress := ethcmn.BigToAddress(big.NewInt(1))
	log := ethtypes.Log{Address: contractAddress}

	testCase := []struct {
		name     string
		malleate func()
		numLogs  int
		isBloom  bool
	}{
		{
			"no logs",
			func() {},
			0,
			false,
		},
		{
			"add log",
			func() {
				suite.stateDB.AddLog(&log)
			},
			1,
			false,
		},
		{
			"bloom",
			func() {},
			0,
			true,
		},
	}

	for _, tc := range testCase {
		tc.malleate()
		logs := suite.stateDB.GetLogs()
		if !tc.isBloom {
			suite.Require().Len(logs, tc.numLogs, tc.name)
			if len(logs) != 0 {
				suite.Require().Equal(log, *logs[0], tc.name)
			}
		} else {
			// get logs bloom from the log
			bloomBytes := ethtypes.LogsBloom(logs)
			bloomFilter := ethtypes.BytesToBloom(bloomBytes)
			suite.Require().True(ethtypes.BloomLookup(bloomFilter, contractAddress), tc.name)
			suite.Require().False(ethtypes.BloomLookup(bloomFilter, ethcmn.BigToAddress(big.NewInt(2))), tc.name)
		}
	}
}

func (suite *StateDBTestSuite) TestStateDB_Balance() {
	testCase := []struct {
		name     string
		malleate func()
		balance  *big.Int
	}{
		{
			"set balance",
			func() {
				suite.stateDB.SetBalance(suite.address, big.NewInt(100))
			},
			big.NewInt(100),
		},
		{
			"sub balance",
			func() {
				suite.stateDB.SubBalance(suite.address, big.NewInt(100))
			},
			big.NewInt(0),
		},
		{
			"add balance",
			func() {
				suite.stateDB.AddBalance(suite.address, big.NewInt(200))
			},
			big.NewInt(200),
		},
	}

	for _, tc := range testCase {
		tc.malleate()
		suite.Require().Equal(tc.balance, suite.stateDB.GetBalance(suite.address), tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateDBNonce() {
	nonce := uint64(123)
	suite.stateDB.SetNonce(suite.address, nonce)
	suite.Require().Equal(nonce, suite.stateDB.GetNonce(suite.address))
}

func (suite *StateDBTestSuite) TestStateDB_Error() {
	nonce := suite.stateDB.GetNonce(ethcmn.Address{})
	suite.Require().Equal(0, int(nonce))
	suite.Require().Error(suite.stateDB.Error())
}

func (suite *StateDBTestSuite) TestStateDB_Database() {
	suite.Require().NotNil(suite.stateDB.Database())
}

func (suite *StateDBTestSuite) TestStateDB_State() {
	key := ethcmn.BytesToHash([]byte("foo"))
	val := ethcmn.BytesToHash([]byte("bar"))
	suite.stateDB.SetState(suite.address, key, val)

	testCase := []struct {
		name    string
		address ethcmn.Address
		key     ethcmn.Hash
		value   ethcmn.Hash
	}{
		{
			"found state",
			suite.address,
			ethcmn.BytesToHash([]byte("foo")),
			ethcmn.BytesToHash([]byte("bar")),
		},
		{
			"state not found",
			suite.address,
			ethcmn.BytesToHash([]byte("key")),
			ethcmn.Hash{},
		},
		{
			"object not found",
			ethcmn.Address{},
			ethcmn.BytesToHash([]byte("foo")),
			ethcmn.Hash{},
		},
	}
	for _, tc := range testCase {
		value := suite.stateDB.GetState(tc.address, tc.key)
		suite.Require().Equal(tc.value, value, tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateDB_Code() {
	testCase := []struct {
		name     string
		address  ethcmn.Address
		code     []byte
		malleate func()
	}{
		{
			"no stored code for state object",
			suite.address,
			nil,
			func() {},
		},
		{
			"existing address",
			suite.address,
			[]byte("code"),
			func() {
				suite.stateDB.SetCode(suite.address, []byte("code"))
			},
		},
		{
			"state object not found",
			ethcmn.Address{},
			nil,
			func() {},
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		suite.Require().Equal(tc.code, suite.stateDB.GetCode(tc.address), tc.name)
		suite.Require().Equal(len(tc.code), suite.stateDB.GetCodeSize(tc.address), tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateDB_Logs() {
	testCase := []struct {
		name string
		log  ethtypes.Log
	}{
		{
			"state db log",
			ethtypes.Log{
				Address:     suite.address,
				Topics:      []ethcmn.Hash{ethcmn.BytesToHash([]byte("topic"))},
				Data:        []byte("data"),
				BlockNumber: 1,
				TxHash:      ethcmn.Hash{},
				TxIndex:     1,
				BlockHash:   ethcmn.Hash{},
				Index:       1,
				Removed:     false,
			},
		},
	}

	for _, tc := range testCase {
		hash := ethcmn.BytesToHash([]byte("hash"))
		logs := []*ethtypes.Log{&tc.log}

		suite.stateDB.SetLogs(logs)
		dbLogs := suite.stateDB.GetLogs()
		suite.Require().Equal(logs, dbLogs, tc.name)

		suite.stateDB.DeleteLogs()
		dbLogs = suite.stateDB.GetLogs()
		suite.Require().Empty(dbLogs, tc.name)

		suite.stateDB.Prepare(hash, ethcmn.BytesToHash([]byte("bhash")), 1)
		suite.stateDB.AddLog(&tc.log)
		newLogs := suite.stateDB.GetLogs()
		suite.Require().Equal(logs, newLogs, tc.name)

		//resets state but checking to see if storekey still persists.
		suite.stateDB.Reset(hash)
		newLogs = suite.stateDB.GetLogs()
		suite.Require().Equal(logs, newLogs, tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateDB_Preimage() {
	hash := ethcmn.BytesToHash([]byte("hash"))
	preimage := []byte("preimage")

	suite.stateDB.AddPreimage(hash, preimage)
	suite.Require().Equal(preimage, suite.stateDB.Preimages()[hash])
}

func (suite *StateDBTestSuite) TestStateDB_Refund() {
	testCase := []struct {
		name      string
		addAmount uint64
		subAmount uint64
		expRefund uint64
		expPanic  bool
	}{
		{
			"refund 0",
			0, 0, 0,
			false,
		},
		{
			"refund positive amount",
			100, 0, 100,
			false,
		},
		{
			"refund panic",
			100, 200, 100,
			true,
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			suite.stateDB.AddRefund(tc.addAmount)
			suite.Require().Equal(tc.addAmount, suite.stateDB.GetRefund())

			if tc.expPanic {
				suite.Panics(func() {
					suite.stateDB.SubRefund(tc.subAmount)
				})
			} else {
				suite.stateDB.SubRefund(tc.subAmount)
				suite.Require().Equal(tc.expRefund, suite.stateDB.GetRefund())
			}
		})
	}
}

func (suite *StateDBTestSuite) TestStateDB_CreateAccount() {
	prevBalance := big.NewInt(12)

	testCase := []struct {
		name     string
		address  ethcmn.Address
		malleate func()
	}{
		{
			"existing account",
			suite.address,
			func() {
				suite.stateDB.AddBalance(suite.address, prevBalance)
			},
		},
		{
			"new account",
			ethcmn.HexToAddress("0x756F45E3FA69347A9A973A725E3C98bC4db0b4c1"),
			func() {
				prevBalance = big.NewInt(0)
			},
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		suite.stateDB.CreateAccount(tc.address)
		suite.Require().True(suite.stateDB.Exist(tc.address), tc.name)
		suite.Require().Equal(prevBalance, suite.stateDB.GetBalance(tc.address), tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateDB_ClearStateObj() {
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	suite.stateDB.CreateAccount(addr)
	suite.Require().True(suite.stateDB.Exist(addr))

	suite.stateDB.ClearStateObjects()
	suite.Require().False(suite.stateDB.Exist(addr))
}

func (suite *StateDBTestSuite) TestStateDB_Reset() {
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	suite.stateDB.CreateAccount(addr)
	suite.Require().True(suite.stateDB.Exist(addr))

	err = suite.stateDB.Reset(ethcmn.BytesToHash(nil))
	suite.Require().NoError(err)
	suite.Require().False(suite.stateDB.Exist(addr))
}

func (suite *StateDBTestSuite) TestSuiteDB_Prepare() {
	thash := ethcmn.BytesToHash([]byte("thash"))
	bhash := ethcmn.BytesToHash([]byte("bhash"))
	txi := 1

	suite.stateDB.Prepare(thash, bhash, txi)
	suite.stateDB.SetBlockHash(bhash)

	suite.Require().Equal(txi, suite.stateDB.TxIndex())
	suite.Require().Equal(bhash, suite.stateDB.BlockHash())
}

func (suite *StateDBTestSuite) TestSuiteDB_Empty() {
	suite.Require().True(suite.stateDB.Empty(suite.address))

	suite.stateDB.SetBalance(suite.address, big.NewInt(100))
	suite.Require().False(suite.stateDB.Empty(suite.address))
}

func (suite *StateDBTestSuite) TestSuiteDB_Suicide() {
	testCase := []struct {
		name    string
		amount  *big.Int
		expPass bool
		delete  bool
	}{
		{
			"suicide zero balance",
			big.NewInt(0),
			false, false,
		},
		{
			"suicide with balance",
			big.NewInt(100),
			true, false,
		},
		{
			"delete",
			big.NewInt(0),
			true, true,
		},
	}

	for _, tc := range testCase {
		if tc.delete {
			_, err := suite.stateDB.Commit(tc.delete)
			suite.Require().NoError(err, tc.name)
			suite.Require().False(suite.stateDB.Exist(suite.address), tc.name)
			continue
		}

		if tc.expPass {
			suite.stateDB.SetBalance(suite.address, tc.amount)
			suicide := suite.stateDB.Suicide(suite.address)
			suite.Require().True(suicide, tc.name)
			suite.Require().True(suite.stateDB.HasSuicided(suite.address), tc.name)
		} else {
			//Suicide only works for an account with non-zero balance/nonce
			priv, err := ethsecp256k1.GenerateKey()
			suite.Require().NoError(err)

			addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
			suicide := suite.stateDB.Suicide(addr)
			suite.Require().False(suicide, tc.name)
			suite.Require().False(suite.stateDB.HasSuicided(addr), tc.name)
		}
	}
}

func (suite *StateDBTestSuite) TestCommitStateDB_Commit() {
	testCase := []struct {
		name       string
		malleate   func()
		deleteObjs bool
		expPass    bool
	}{
		{
			"commit suicided",
			func() {
				ok := suite.stateDB.Suicide(suite.address)
				suite.Require().True(ok)
			},
			true, true,
		},
		{
			"commit with dirty value",
			func() {
				suite.stateDB.SetCode(suite.address, []byte("code"))
			},
			false, true,
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		hash, err := suite.stateDB.Commit(tc.deleteObjs)
		suite.Require().Equal(ethcmn.Hash{}, hash)

		if !tc.expPass {
			suite.Require().Error(err, tc.name)
			continue
		}

		suite.Require().NoError(err, tc.name)
		acc := suite.app.AccountKeeper.GetAccount(suite.ctx, sdk.AccAddress(suite.address.Bytes()))

		if tc.deleteObjs {
			suite.Require().Nil(acc, tc.name)
			continue
		}

		suite.Require().NotNil(acc, tc.name)
		ethAcc, ok := acc.(*ethermint.EthAccount)
		suite.Require().True(ok)
		suite.Require().Equal(ethcrypto.Keccak256([]byte("code")), ethAcc.CodeHash)
	}
}

func (suite *StateDBTestSuite) TestCommitStateDB_Finalize() {
	testCase := []struct {
		name       string
		malleate   func()
		deleteObjs bool
		expPass    bool
	}{
		{
			"finalize suicided",
			func() {
				ok := suite.stateDB.Suicide(suite.address)
				suite.Require().True(ok)
			},
			true, true,
		},
		{
			"finalize, not suicided",
			func() {
				suite.stateDB.AddBalance(suite.address, big.NewInt(5))
			},
			false, true,
		},
		{
			"finalize, dirty storage",
			func() {
				suite.stateDB.SetState(suite.address, ethcmn.BytesToHash([]byte("key")), ethcmn.BytesToHash([]byte("value")))
			},
			false, true,
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		suite.stateDB.IntermediateRoot(tc.deleteObjs)

		if !tc.expPass {
			hash := suite.stateDB.GetCommittedState(suite.address, ethcmn.BytesToHash([]byte("key")))
			suite.Require().NotEqual(ethcmn.Hash{}, hash, tc.name)
			continue
		}

		acc := suite.app.AccountKeeper.GetAccount(suite.ctx, sdk.AccAddress(suite.address.Bytes()))

		if tc.deleteObjs {
			suite.Require().Nil(acc, tc.name)
			continue
		}

		suite.Require().NotNil(acc, tc.name)
	}
}
func (suite *StateDBTestSuite) TestCommitStateDB_GetCommittedState() {
	hash := suite.stateDB.GetCommittedState(ethcmn.Address{}, ethcmn.BytesToHash([]byte("key")))
	suite.Require().Equal(ethcmn.Hash{}, hash)
}

func (suite *StateDBTestSuite) TestCommitStateDB_Snapshot() {
	id := suite.stateDB.Snapshot()
	suite.Require().NotPanics(func() {
		suite.stateDB.RevertToSnapshot(id)
	})

	suite.Require().Panics(func() {
		suite.stateDB.RevertToSnapshot(-1)
	}, "invalid revision should panic")
}

func (suite *StateDBTestSuite) TestCommitStateDB_ForEachStorage() {
	var storage types.Storage

	testCase := []struct {
		name      string
		malleate  func()
		callback  func(key, value ethcmn.Hash) (stop bool)
		expValues []ethcmn.Hash
	}{
		{
			"aggregate state",
			func() {
				for i := 0; i < 5; i++ {
					suite.stateDB.SetState(suite.address, ethcmn.BytesToHash([]byte(fmt.Sprintf("key%d", i))), ethcmn.BytesToHash([]byte(fmt.Sprintf("value%d", i))))
				}
			},
			func(key, value ethcmn.Hash) bool {
				storage = append(storage, types.NewState(key, value))
				return false
			},
			[]ethcmn.Hash{
				ethcmn.BytesToHash([]byte("value0")),
				ethcmn.BytesToHash([]byte("value1")),
				ethcmn.BytesToHash([]byte("value2")),
				ethcmn.BytesToHash([]byte("value3")),
				ethcmn.BytesToHash([]byte("value4")),
			},
		},
		{
			"filter state",
			func() {
				suite.stateDB.SetState(suite.address, ethcmn.BytesToHash([]byte("key")), ethcmn.BytesToHash([]byte("value")))
				suite.stateDB.SetState(suite.address, ethcmn.BytesToHash([]byte("filterkey")), ethcmn.BytesToHash([]byte("filtervalue")))
			},
			func(key, value ethcmn.Hash) bool {
				if value == ethcmn.BytesToHash([]byte("filtervalue")) {
					storage = append(storage, types.NewState(key, value))
					return true
				}
				return false
			},
			[]ethcmn.Hash{
				ethcmn.BytesToHash([]byte("filtervalue")),
			},
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.malleate()
			suite.stateDB.Commit(false)

			err := suite.stateDB.ForEachStorage(suite.address, tc.callback)
			suite.Require().NoError(err)
			suite.Require().Equal(len(tc.expValues), len(storage), fmt.Sprintf("Expected values:\n%v\nStorage Values\n%v", tc.expValues, storage))

			vals := make([]ethcmn.Hash, len(storage))
			for i := range storage {
				vals[i] = storage[i].Value
			}

			suite.Require().ElementsMatch(tc.expValues, vals)
		})
		storage = types.Storage{}
	}
}

func (suite *StateDBTestSuite) TestCommitStateDB_AccessList() {
	addr := ethcmn.Address([20]byte{77})
	hash := ethcmn.Hash([32]byte{99})

	suite.Require().False(suite.stateDB.AddressInAccessList(addr))

	suite.stateDB.AddAddressToAccessList(addr)
	suite.Require().True(suite.stateDB.AddressInAccessList(addr))
	addrIn, slotIn := suite.stateDB.SlotInAccessList(addr, hash)
	suite.Require().True(addrIn)
	suite.Require().False(slotIn)

	suite.stateDB.AddSlotToAccessList(addr, hash)
	addrIn, slotIn = suite.stateDB.SlotInAccessList(addr, hash)
	suite.Require().True(addrIn)
	suite.Require().True(slotIn)
}

func (suite *StateDBTestSuite) TestCommitStateDB_ContractDeploymentWhitelist() {
	// create addresses for test
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()

	testCase := []struct {
		name           string
		targetAddrList types.AddressList
		// true -> add, false -> delete
		isAdded     bool
		expectedLen int
	}{
		{
			"add empty list into whitelist",
			types.AddressList{},
			true,
			0,
		},
		{
			"add list with one member into whitelist",
			types.AddressList{addr1},
			true,
			1,
		},
		{
			"add list with two members into the whitelist that has contained one member already",
			types.AddressList{addr1, addr2},
			true,
			2,
		},
		{
			"delete empty from whitelist",
			types.AddressList{},
			false,
			2,
		},
		{
			"delete list with one member from whitelist",
			types.AddressList{addr1},
			false,
			1,
		},
		{
			"delete list with two members from the whitelist that has contained one member only",
			types.AddressList{addr1, addr2},
			false,
			0,
		},
		{
			"delete list with two members from the empty whitelist",
			types.AddressList{addr1, addr2},
			false,
			0,
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			if tc.isAdded {
				suite.stateDB.SetContractDeploymentWhitelist(tc.targetAddrList)
			} else {
				suite.stateDB.DeleteContractDeploymentWhitelist(tc.targetAddrList)
			}

			whitelist := suite.stateDB.GetContractDeploymentWhitelist()
			suite.Require().Equal(tc.expectedLen, len(whitelist))
		})
	}
}

func (suite *StateDBTestSuite) TestCommitStateDB_ContractBlockedList() {
	// create addresses for test
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()

	testCase := []struct {
		name           string
		targetAddrList types.AddressList
		// true -> add, false -> delete
		isAdded     bool
		expectedLen int
	}{
		{
			"add empty list into blocked list",
			types.AddressList{},
			true,
			0,
		},
		{
			"add list with one member into blocked list",
			types.AddressList{addr1},
			true,
			1,
		},
		{
			"add list with two members into the blocked list that has contained one member already",
			types.AddressList{addr1, addr2},
			true,
			2,
		},
		{
			"delete empty from blocked list",
			types.AddressList{},
			false,
			2,
		},
		{
			"delete list with one member from blocked list",
			types.AddressList{addr1},
			false,
			1,
		},
		{
			"delete list with two members from the blocked list that has contained one member only",
			types.AddressList{addr1, addr2},
			false,
			0,
		},
		{
			"delete list with two members from the empty blocked list",
			types.AddressList{addr1, addr2},
			false,
			0,
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			if tc.isAdded {
				suite.stateDB.SetContractBlockedList(tc.targetAddrList)
			} else {
				suite.stateDB.DeleteContractBlockedList(tc.targetAddrList)
			}

			blockedList := suite.stateDB.GetContractBlockedList()
			suite.Require().Equal(tc.expectedLen, len(blockedList))
		})
	}
}

func (suite *StateDBTestSuite) TestCommitStateDB_ContractMethodBlockedList() {
	// create addresses for test
	bcMethodOne1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "aaaa",
				Extra: "aaaa()",
			},
		},
	}
	bcMethodTwo1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "aaaa",
				Extra: "aaaa()",
			},
		},
	}

	bcMethodOne3 := types.BlockedContract{
		Address: bcMethodOne1.Address,
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "bbbb",
				Extra: "bbbb()",
			},
		},
	}
	methods := types.ContractMethods{}
	methods = append(methods, bcMethodOne1.BlockMethods...)
	methods = append(methods, bcMethodOne3.BlockMethods...)
	expectBcMethodOne3 := types.NewBlockContract(bcMethodOne1.Address, methods)

	bcMethodOne4 := types.BlockedContract{
		Address: bcMethodOne1.Address,
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "bbbb",
				Extra: "bbbb()",
			},
			types.ContractMethod{
				Sign:  "cccc",
				Extra: "cccc()",
			},
		},
	}
	methods = types.ContractMethods{}
	methods = append(methods, bcMethodOne1.BlockMethods...)
	methods = append(methods, bcMethodOne4.BlockMethods...)
	expectBcMethodOne4 := types.NewBlockContract(bcMethodOne1.Address, methods)

	bcMethodOne5 := types.BlockedContract{
		Address: bcMethodOne1.Address,
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "bbbb",
				Extra: "bbbb()",
			},
			types.ContractMethod{
				Sign:  "cccc",
				Extra: "cccc()",
			},
			types.ContractMethod{
				Sign:  "dddd",
				Extra: "dddd()",
			},
		},
	}

	testCase := []struct {
		name           string
		targetAddrList types.BlockedContractList
		// true -> add, false -> delete
		isAdded              bool
		expectedLen          int
		expectedContractList types.BlockedContractList
		success              bool
	}{
		{
			"add empty blocked contract list",
			types.BlockedContractList{},
			true,
			0,
			nil,
			true,
		},
		{
			"add list with one member into blocked list",
			types.BlockedContractList{bcMethodOne1},
			true,
			1,
			types.BlockedContractList{bcMethodOne1},
			true,
		},
		{
			"add list with two members into the blocked list that has contained one member already",
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			true,
			2,
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			true,
		},
		{
			"add list with one members(method empty) into the blocked list that has contained one member already",
			types.BlockedContractList{bcMethodOne1},
			true,
			2,
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			true,
		},
		{
			"delete empty from blocked list",
			types.BlockedContractList{},
			false,
			2,
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			true,
		},
		{
			"delete list with one members from the blocked list that has contained one member only",
			types.BlockedContractList{bcMethodTwo1},
			false,
			1,
			types.BlockedContractList{bcMethodOne1},
			true,
		},
		{
			"delete list with two members from the empty blocked list",
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			false,
			0,
			nil,
			false,
		},
		{
			"reset contract method blocked list",
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			true,
			2,
			types.BlockedContractList{bcMethodOne1, bcMethodTwo1},
			true,
		},
		{
			"add new method into contract method blocked list",
			types.BlockedContractList{bcMethodOne3},
			true,
			2,
			types.BlockedContractList{*expectBcMethodOne3, bcMethodTwo1},
			true,
		},
		{
			"add new methods which is one method exist into contract method blocked list",
			types.BlockedContractList{bcMethodOne4},
			true,
			2,
			types.BlockedContractList{*expectBcMethodOne4, bcMethodTwo1},
			true,
		},
		{
			"delete methods which is not exist from contract method blocked list",
			types.BlockedContractList{bcMethodOne5},
			false,
			2,
			types.BlockedContractList{*expectBcMethodOne4, bcMethodTwo1},
			false,
		},
		{
			"delete all methods from contract method blocked list",
			types.BlockedContractList{*expectBcMethodOne4},
			false,
			1,
			types.BlockedContractList{bcMethodTwo1},
			true,
		},
	}

	for _, tc := range testCase {
		suite.Run(tc.name, func() {
			var err sdk.Error
			if tc.isAdded {
				err = suite.stateDB.InsertContractMethodBlockedList(tc.targetAddrList)
			} else {
				err = suite.stateDB.DeleteContractMethodBlockedList(tc.targetAddrList)
			}
			if tc.success {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

			blockedList := suite.stateDB.GetContractMethodBlockedList()
			suite.Require().Equal(tc.expectedLen, len(blockedList))
			if tc.expectedLen != 0 {
				ok := types.BlockedContractListIsEqual(suite.T(), tc.expectedContractList, blockedList)
				suite.Require().True(ok)
			}
		})
	}
}

func (suite *StateDBTestSuite) TestCommitStateDB_ContractMethodBlockedList_BlockedList() {
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	bcMethodOne1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "aaaa",
				Extra: "aaaa()",
			},
		},
	}
	bcMethodTwo1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "aaaa",
				Extra: "aaaa()",
			},
		},
	}

	// set contract method blocked list with same blocked list
	suite.stateDB.SetContractBlockedList(types.AddressList{addr1})
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
	bcl := suite.stateDB.GetContractMethodBlockedList()
	ok := types.BlockedContractListIsEqual(suite.T(), types.BlockedContractList{bcMethodOne1}, bcl)
	suite.Require().True(ok)
	// set contract method blocked list with not same blocked list
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodTwo1})
	bcl = suite.stateDB.GetContractMethodBlockedList()
	ok = types.BlockedContractListIsEqual(suite.T(), types.BlockedContractList{bcMethodOne1, bcMethodTwo1}, bcl)
	suite.Require().True(ok)

	// set blocked list with  same method blocked list
	suite.stateDB.SetContractBlockedList(types.AddressList{addr1})
	bcl = suite.stateDB.GetContractMethodBlockedList()
	expect := types.NewBlockContract(bcMethodOne1.Address, nil)
	ok = types.BlockedContractListIsEqual(suite.T(), types.BlockedContractList{*expect, bcMethodTwo1}, bcl)
	suite.Require().True(ok)

}

func (suite *StateDBTestSuite) TestCommitStateDB_GetContractMethodBlockedByAddress() {
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()
	bcMethodOne1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "aaaa",
				Extra: "aaaa()",
			},
		},
	}
	bcMethodTwo1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "bbbb",
				Extra: "bbbb()",
			},
		},
	}

	// get blocked list is not exist
	bc := suite.stateDB.GetContractMethodBlockedByAddress(addr1)
	suite.Require().Nil(bc)

	// get blocked list
	suite.stateDB.SetContractBlockedList(types.AddressList{addr1, addr2})
	bc = suite.stateDB.GetContractMethodBlockedByAddress(addr1)
	suite.Require().NotNil(bc)
	suite.Require().Equal(0, len(bc.BlockMethods))
	suite.Require().Equal(addr1, bc.Address.Bytes())

	// get blocked list
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
	bc = suite.stateDB.GetContractMethodBlockedByAddress(addr1)
	suite.Require().NotNil(bc)
	ok := types.BlockedContractListIsEqual(suite.T(), types.BlockedContractList{bcMethodOne1}, types.BlockedContractList{*bc})
	suite.Require().True(ok)

	// get blocked list from cache
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodTwo1})
	bc = suite.stateDB.GetContractMethodBlockedByAddress(addr2)
	suite.Require().NotNil(bc)
	ok = types.BlockedContractListIsEqual(suite.T(), types.BlockedContractList{bcMethodTwo1}, types.BlockedContractList{*bc})
	suite.Require().True(ok)

	bc = suite.stateDB.GetContractMethodBlockedByAddress(addr2)
	suite.Require().NotNil(bc)
	ok = types.BlockedContractListIsEqual(suite.T(), types.BlockedContractList{bcMethodTwo1}, types.BlockedContractList{*bc})
	suite.Require().True(ok)
}
func (suite *StateDBTestSuite) TestCacheSet() {
	addr := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	method1 := types.ContractMethod{
		Sign:  "aaaa",
		Extra: "aaaa()",
	}
	method2 := types.ContractMethod{
		Sign:  "bbbb",
		Extra: "bbbb()",
	}
	sourceBc := types.BlockedContract{
		Address: addr,
		BlockMethods: types.ContractMethods{
			method1, method2,
		},
	}
	sourceBcl := types.BlockedContractList{sourceBc}
	suite.stateDB.InsertContractMethodBlockedList(sourceBcl)
	bc := suite.stateDB.GetContractMethodBlockedByAddress(addr)
	methods := &bc.BlockMethods
	*methods = (*methods)[0:0]
	*methods = append(*methods, types.ContractMethod{
		Sign:  "dddd",
		Extra: "dddd()",
	})
	bc = suite.stateDB.GetContractMethodBlockedByAddress(addr)
	ok := types.BlockedContractListIsEqual(suite.T(), sourceBcl, types.BlockedContractList{*bc})
	suite.Require().True(ok)
}
func (suite *StateDBTestSuite) TestCommitStateDB_IsContractMethodBlocked() {
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()
	bcMethodOne1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x0}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "aaaa",
				Extra: "aaaa()",
			},
		},
	}
	bcMethodTwo1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "bbbb",
				Extra: "bbbb()",
			},
		},
	}

	// contract method is not exist
	ok := suite.stateDB.IsContractMethodBlocked(addr1, "")
	suite.Require().False(ok)
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "aaaa")
	suite.Require().False(ok)

	// contract all method is blocked
	suite.stateDB.SetContractBlockedList(types.AddressList{addr1, addr2})
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "")
	suite.Require().False(ok)
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "aaaa")
	suite.Require().False(ok)

	// contract aaaa method is blocked
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodOne1})
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "")
	suite.Require().False(ok)
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "aaaa")
	suite.Require().True(ok)

	// contract aaaa method is blocked
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodTwo1})
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "")
	suite.Require().False(ok)
	ok = suite.stateDB.IsContractMethodBlocked(addr1, "bbbb")
	suite.Require().False(ok)
}

func (suite *StateDBTestSuite) TestCommitStateDB_IsContractInBlockedList() {
	addr1 := ethcmn.BytesToAddress([]byte{0x0}).Bytes()
	addr2 := ethcmn.BytesToAddress([]byte{0x1}).Bytes()

	bcMethodTwo1 := types.BlockedContract{
		Address: ethcmn.BytesToAddress([]byte{0x1}).Bytes(),
		BlockMethods: types.ContractMethods{
			types.ContractMethod{
				Sign:  "bbbb",
				Extra: "bbbb()",
			},
		},
	}
	// contract is not exist
	ok := suite.stateDB.IsContractInBlockedList(addr1)
	suite.Require().False(ok)
	ok = suite.stateDB.IsContractInBlockedList(addr2)
	suite.Require().False(ok)
	// contract is exist
	suite.stateDB.SetContractBlockedList(types.AddressList{addr1})
	ok = suite.stateDB.IsContractInBlockedList(addr1)
	suite.Require().True(ok)
	ok = suite.stateDB.IsContractInBlockedList(addr2)
	suite.Require().False(ok)

	// contract method is blocked
	suite.stateDB.InsertContractMethodBlockedList(types.BlockedContractList{bcMethodTwo1})
	ok = suite.stateDB.IsContractInBlockedList(addr2)
	suite.Require().False(ok)
	ok = suite.stateDB.IsContractInBlockedList(addr1)
	suite.Require().True(ok)
}
