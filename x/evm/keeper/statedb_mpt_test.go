package keeper_test

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/evm/types"

	ethermint "github.com/okx/okbchain/app/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

func (suite *KeeperMptTestSuite) TestCommitStateDB_CommitMpt() {
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

		hash, err := suite.stateDB.WithContext(suite.ctx).Commit(tc.deleteObjs)
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

func (suite *KeeperMptTestSuite) TestCommitStateDB_ForEachStorageMpt() {
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
			suite.stateDB.WithContext(suite.ctx).Commit(false)
			suite.app.Commit(abci.RequestCommit{})
			types.ResetCommitStateDB(suite.stateDB, suite.app.EvmKeeper.GenerateCSDBParams(), &suite.ctx)
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

func (suite *KeeperMptTestSuite) TestCommitStateDB_GetCommittedStateMpt() {
	suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte("key")), ethcmn.BytesToHash([]byte("value")))
	suite.stateDB.Commit(false)

	hash := suite.stateDB.WithContext(suite.ctx).GetCommittedState(suite.address, ethcmn.BytesToHash([]byte("key")))
	suite.Require().Equal(ethcmn.BytesToHash([]byte("value")), hash)
}

func (suite *KeeperMptTestSuite) TestStateDB_CodeMpt() {
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
				suite.stateDB.Commit(false)
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
		suite.Require().Equal(tc.code, suite.stateDB.WithContext(suite.ctx).GetCodeByHashInRawDB(tc.codeHash), tc.name)
	}
}

func (suite *KeeperMptTestSuite) TestStateDB_StateMpt() {
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

func (suite *KeeperMptTestSuite) TestStorageTrieMpt() {
	for i := 0; i < 5; i++ {
		suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte(fmt.Sprintf("key%d", i))), ethcmn.BytesToHash([]byte(fmt.Sprintf("value%d", i))))
	}

	trie := suite.stateDB.WithContext(suite.ctx).StorageTrie(suite.address)
	suite.Require().NotNil(trie, "Ethermint now use a direct storage trie.")
}
