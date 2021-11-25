package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/x/evm/types"
)

func (suite *KeeperTestSuite) TestBloomFilter() {
	// Prepare db for logs
	tHash := ethcmn.BytesToHash([]byte{0x1})
	bHash := ethcmn.BytesToHash([]byte{0x1})
	suite.stateDB.WithContext(suite.ctx).Prepare(tHash, bHash, 0)
	contractAddress := ethcmn.BigToAddress(big.NewInt(1))
	log := ethtypes.Log{Address: contractAddress}
	logs := []*ethtypes.Log{&log}

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
				suite.stateDB.WithContext(suite.ctx).SetLogs(tHash, logs)
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
		logs, err := suite.stateDB.WithContext(suite.ctx).GetLogs(tHash)
		if !tc.isBloom {
			suite.Require().NoError(err, tc.name)
			suite.Require().Len(logs, tc.numLogs, tc.name)
			if len(logs) != 0 {
				suite.Require().Equal(log, *logs[0], tc.name)
			}
		} else {
			// get logs bloom from the log
			bloomInt := ethtypes.LogsBloom(logs)
			bloomFilter := ethtypes.BytesToBloom(bloomInt)
			suite.Require().True(ethtypes.BloomLookup(bloomFilter, contractAddress), tc.name)
			suite.Require().False(ethtypes.BloomLookup(bloomFilter, ethcmn.BigToAddress(big.NewInt(2))), tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestStateDB_Balance() {
	testCase := []struct {
		name     string
		malleate func()
		balance  *big.Int
	}{
		{
			"set balance",
			func() {
				suite.stateDB.WithContext(suite.ctx).SetBalance(suite.address, big.NewInt(100))
			},
			big.NewInt(100),
		},
		{
			"sub balance",
			func() {
				suite.stateDB.WithContext(suite.ctx).SubBalance(suite.address, big.NewInt(100))
			},
			big.NewInt(0),
		},
		{
			"add balance",
			func() {
				suite.stateDB.WithContext(suite.ctx).AddBalance(suite.address, big.NewInt(200))
			},
			big.NewInt(200),
		},
	}

	for _, tc := range testCase {
		tc.malleate()
		suite.Require().Equal(tc.balance, suite.stateDB.WithContext(suite.ctx).GetBalance(suite.address), tc.name)
	}
}

func (suite *KeeperTestSuite) TestStateDBNonce() {
	nonce := uint64(123)
	suite.stateDB.WithContext(suite.ctx).SetNonce(suite.address, nonce)
	suite.Require().Equal(nonce, suite.stateDB.WithContext(suite.ctx).GetNonce(suite.address))
}

func (suite *KeeperTestSuite) TestStateDB_Error() {
	nonce := suite.stateDB.WithContext(suite.ctx).GetNonce(ethcmn.Address{})
	suite.Require().Equal(0, int(nonce))
	suite.Require().Error(suite.stateDB.WithContext(suite.ctx).Error())
}

func (suite *KeeperTestSuite) TestStateDB_Database() {
	suite.Require().Nil(suite.stateDB.WithContext(suite.ctx).Database())
}

func (suite *KeeperTestSuite) TestStateDB_State() {
	key := ethcmn.BytesToHash([]byte("foo"))
	val := ethcmn.BytesToHash([]byte("bar"))
	suite.stateDB.WithContext(suite.ctx).SetState(suite.address, key, val)

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
		value := suite.stateDB.WithContext(suite.ctx).GetState(tc.address, tc.key)
		suite.Require().Equal(tc.value, value, tc.name)
	}
}

func (suite *KeeperTestSuite) TestStateDB_Code() {
	testCase := []struct {
		name     string
		address  ethcmn.Address
		code     []byte
		codeHash ethcmn.Hash
		malleate func()
	}{
		{
			"no stored code for state object",
			suite.address,
			nil,
			ethcmn.HexToHash("0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"),
			func() {},
		},
		{
			"existing address",
			suite.address,
			[]byte("code"),
			ethcmn.HexToHash("0x2dc081a8d6d4714c79b5abd2e9b08c3a33b4ef1dcf946ef8b8cf6c495014f47b"),
			func() {
				suite.stateDB.WithContext(suite.ctx).SetCode(suite.address, []byte("code"))
			},
		},
		{
			"state object not found",
			ethcmn.Address{},
			nil,
			ethcmn.HexToHash("0"),
			func() {},
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		suite.Require().Equal(tc.code, suite.stateDB.WithContext(suite.ctx).GetCode(tc.address), tc.name)
		suite.Require().Equal(len(tc.code), suite.stateDB.WithContext(suite.ctx).GetCodeSize(tc.address), tc.name)
		suite.Require().Equal(tc.codeHash, suite.stateDB.WithContext(suite.ctx).GetCodeHash(tc.address), tc.name)
	}
}

func (suite *KeeperTestSuite) TestStateDB_Logs() {
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

		err := suite.stateDB.WithContext(suite.ctx).SetLogs(hash, logs)
		suite.Require().NoError(err, tc.name)
		dbLogs, err := suite.stateDB.WithContext(suite.ctx).GetLogs(hash)
		suite.Require().NoError(err, tc.name)
		suite.Require().Equal(logs, dbLogs, tc.name)
	}
}

func (suite *KeeperTestSuite) TestStateDB_Preimage() {
	hash := ethcmn.BytesToHash([]byte("hash"))
	preimage := []byte("preimage")

	suite.stateDB.WithContext(suite.ctx).AddPreimage(hash, preimage)
	suite.Require().Equal(preimage, suite.stateDB.WithContext(suite.ctx).Preimages()[hash])
}

func (suite *KeeperTestSuite) TestStateDB_Refund() {
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

			suite.stateDB.WithContext(suite.ctx).AddRefund(tc.addAmount)
			suite.Require().Equal(tc.addAmount, suite.stateDB.WithContext(suite.ctx).GetRefund())

			if tc.expPanic {
				suite.Panics(func() {
					suite.stateDB.WithContext(suite.ctx).SubRefund(tc.subAmount)
				})
			} else {
				suite.stateDB.WithContext(suite.ctx).SubRefund(tc.subAmount)
				suite.Require().Equal(tc.expRefund, suite.stateDB.WithContext(suite.ctx).GetRefund())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestStateDB_CreateAccount() {
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
				suite.stateDB.WithContext(suite.ctx).AddBalance(suite.address, prevBalance)
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
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.malleate()

			suite.stateDB.WithContext(suite.ctx).CreateAccount(tc.address)
			suite.Require().True(suite.stateDB.WithContext(suite.ctx).Exist(tc.address))
			suite.Require().Equal(prevBalance, suite.stateDB.WithContext(suite.ctx).GetBalance(tc.address))
		})
	}
}

func (suite *KeeperTestSuite) TestStateDB_ClearStateObj() {
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	suite.stateDB.WithContext(suite.ctx).CreateAccount(addr)
	suite.Require().True(suite.stateDB.WithContext(suite.ctx).Exist(addr))

	suite.stateDB.WithContext(suite.ctx).ClearStateObjects()
	suite.Require().False(suite.stateDB.WithContext(suite.ctx).Exist(addr))
}

func (suite *KeeperTestSuite) TestStateDB_Reset() {
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	suite.stateDB.WithContext(suite.ctx).CreateAccount(addr)
	suite.Require().True(suite.stateDB.WithContext(suite.ctx).Exist(addr))

	err = suite.stateDB.WithContext(suite.ctx).Reset(ethcmn.BytesToHash(nil))
	suite.Require().NoError(err)
	suite.Require().False(suite.stateDB.WithContext(suite.ctx).Exist(addr))
}

func (suite *KeeperTestSuite) TestSuiteDB_Prepare() {
	thash := ethcmn.BytesToHash([]byte("thash"))
	bhash := ethcmn.BytesToHash([]byte("bhash"))
	txi := 1

	suite.stateDB.WithContext(suite.ctx).Prepare(thash, bhash, txi)

	suite.Require().Equal(txi, suite.stateDB.WithContext(suite.ctx).TxIndex())
	suite.Require().Equal(bhash, suite.stateDB.WithContext(suite.ctx).BlockHash())
}

func (suite *KeeperTestSuite) TestSuiteDB_Empty() {
	suite.Require().True(suite.stateDB.WithContext(suite.ctx).Empty(suite.address))

	suite.stateDB.WithContext(suite.ctx).SetBalance(suite.address, big.NewInt(100))
	suite.Require().False(suite.stateDB.WithContext(suite.ctx).Empty(suite.address))
}

func (suite *KeeperTestSuite) TestSuiteDB_Suicide() {
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
			_, err := suite.stateDB.WithContext(suite.ctx).Commit(tc.delete)
			suite.Require().NoError(err, tc.name)
			suite.Require().False(suite.stateDB.WithContext(suite.ctx).Exist(suite.address), tc.name)
			continue
		}

		if tc.expPass {
			suite.stateDB.WithContext(suite.ctx).SetBalance(suite.address, tc.amount)
			suicide := suite.stateDB.WithContext(suite.ctx).Suicide(suite.address)
			suite.Require().True(suicide, tc.name)
			suite.Require().True(suite.stateDB.WithContext(suite.ctx).HasSuicided(suite.address), tc.name)
		} else {
			//Suicide only works for an account with non-zero balance/nonce
			priv, err := ethsecp256k1.GenerateKey()
			suite.Require().NoError(err)

			addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
			suicide := suite.stateDB.WithContext(suite.ctx).Suicide(addr)
			suite.Require().False(suicide, tc.name)
			suite.Require().False(suite.stateDB.WithContext(suite.ctx).HasSuicided(addr), tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestCommitStateDB_Commit() {
	testCase := []struct {
		name       string
		malleate   func()
		deleteObjs bool
		expPass    bool
	}{
		{
			"commit suicided",
			func() {
				ok := suite.stateDB.WithContext(suite.ctx).Suicide(suite.address)
				suite.Require().True(ok)
			},
			true, true,
		},
		{
			"commit with dirty value",
			func() {
				suite.stateDB.WithContext(suite.ctx).SetCode(suite.address, []byte("code"))
			},
			false, true,
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		_, err := suite.stateDB.WithContext(suite.ctx).Commit(tc.deleteObjs)

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

func (suite *KeeperTestSuite) TestCommitStateDB_Finalize() {
	testCase := []struct {
		name       string
		malleate   func()
		deleteObjs bool
		expPass    bool
	}{
		{
			"finalize suicided",
			func() {
				ok := suite.stateDB.WithContext(suite.ctx).Suicide(suite.address)
				suite.Require().True(ok)
			},
			true, true,
		},
		{
			"finalize, not suicided",
			func() {
				suite.stateDB.WithContext(suite.ctx).AddBalance(suite.address, big.NewInt(5))
			},
			false, true,
		},
		{
			"finalize, dirty storage",
			func() {
				suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte("key")), ethcmn.BytesToHash([]byte("value")))
			},
			false, true,
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		suite.stateDB.WithContext(suite.ctx).Finalise(tc.deleteObjs)

		if !tc.expPass {
			hash := suite.stateDB.WithContext(suite.ctx).GetCommittedState(suite.address, ethcmn.BytesToHash([]byte("key")))
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

	_, err := suite.stateDB.WithContext(suite.ctx).Commit(false)
	suite.Require().Nil(err, "successful get the root hash of the state")
}

func (suite *KeeperTestSuite) TestCommitStateDB_GetCommittedState() {
	hash := suite.stateDB.WithContext(suite.ctx).GetCommittedState(ethcmn.Address{}, ethcmn.BytesToHash([]byte("key")))
	suite.Require().Equal(ethcmn.Hash{}, hash)
}

func (suite *KeeperTestSuite) TestCommitStateDB_Snapshot() {
	id := suite.stateDB.WithContext(suite.ctx).Snapshot()
	suite.Require().NotPanics(func() {
		suite.stateDB.WithContext(suite.ctx).RevertToSnapshot(id)
	})

	suite.Require().Panics(func() {
		suite.stateDB.WithContext(suite.ctx).RevertToSnapshot(-1)
	}, "invalid revision should panic")
}

func (suite *KeeperTestSuite) TestCommitStateDB_ForEachStorage() {
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
					suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte(fmt.Sprintf("key%d", i))), ethcmn.BytesToHash([]byte(fmt.Sprintf("value%d", i))))
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
				suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte("key")), ethcmn.BytesToHash([]byte("value")))
				suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte("filterkey")), ethcmn.BytesToHash([]byte("filtervalue")))
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
			suite.stateDB.WithContext(suite.ctx).Finalise(false)

			err := suite.stateDB.WithContext(suite.ctx).ForEachStorage(suite.address, tc.callback)
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

func (suite *KeeperTestSuite) TestStorageTrie() {
	for i := 0; i < 5; i++ {
		suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte(fmt.Sprintf("key%d", i))), ethcmn.BytesToHash([]byte(fmt.Sprintf("value%d", i))))
	}

	trie := suite.stateDB.WithContext(suite.ctx).StorageTrie(suite.address)
	suite.Require().Equal(nil, trie, "Ethermint does not use a direct storage trie.")
}
