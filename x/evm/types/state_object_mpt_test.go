package types_test

import ethcmn "github.com/ethereum/go-ethereum/common"

func (suite *StateDBMptTestSuite) TestStateObject_State() {
	testCase := []struct {
		name     string
		key      ethcmn.Hash
		expValue ethcmn.Hash
		malleate func()
	}{
		{
			"no set value, load from KVStore",
			ethcmn.BytesToHash([]byte("key")),
			ethcmn.Hash{},
			func() {},
		},
		{
			"no-op SetState",
			ethcmn.BytesToHash([]byte("key")),
			ethcmn.Hash{},
			func() {
				suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key")), ethcmn.Hash{})
			},
		},
		{
			"cached value",
			ethcmn.BytesToHash([]byte("key1")),
			ethcmn.BytesToHash([]byte("value1")),
			func() {
				suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key1")), ethcmn.BytesToHash([]byte("value1")))
			},
		},
		{
			"update value",
			ethcmn.BytesToHash([]byte("key1")),
			ethcmn.BytesToHash([]byte("value2")),
			func() {
				suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key1")), ethcmn.BytesToHash([]byte("value2")))
			},
		},
		{
			"update various keys",
			ethcmn.BytesToHash([]byte("key1")),
			ethcmn.BytesToHash([]byte("value1")),
			func() {
				suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key1")), ethcmn.BytesToHash([]byte("value1")))
				suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key2")), ethcmn.BytesToHash([]byte("value2")))
				suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key3")), ethcmn.BytesToHash([]byte("value3")))
			},
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		value := suite.stateObject.GetState(suite.stateDB.Database(), tc.key)
		suite.Require().Equal(tc.expValue, value, tc.name)
	}
}
