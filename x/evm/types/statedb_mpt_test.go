package types_test

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/okx/okbchain/app/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/tendermint/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StateDBMptTestSuite struct {
	StateDBTestSuite
}

func (suite *StateDBMptTestSuite) SetupTest() {
	types.UnittestOnlySetMilestoneMarsHeight(1)

	suite.StateDBTestSuite.SetupTest()
}

func TestStateDBMptTestSuite(t *testing.T) {
	suite.Run(t, new(StateDBMptTestSuite))
}

func (suite *StateDBMptTestSuite) TestGetHeightHashMpt() {
	hash := suite.stateDB.GetHeightHash(0)
	suite.Require().Equal(ethcmn.Hash{}.String(), hash.String())

	expHash := ethcmn.BytesToHash([]byte("hash"))
	suite.stateDB.SetHeightHash(10, expHash)

	hash = suite.stateDB.GetHeightHash(10)
	suite.Require().Equal(expHash.String(), hash.String())
}

func (suite *StateDBMptTestSuite) TestStateDB_StateMpt() {
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

func (suite *StateDBMptTestSuite) TestCommitStateDB_CommitMpt() {
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
