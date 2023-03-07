package mpt

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/okx/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okx/exchain/libs/cosmos-sdk/store/types"
	abci "github.com/okx/exchain/libs/tendermint/abci/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

var (
	commonKeys   = []string{"key1", "key2", "key3", "key4", "key5"}
	commonValues = []string{"value1", "value2", "value3", "value4", "value5"}

	randKeyNum = 1000
)

func randBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, _ = rand.Read(b)
	return b
}

type StoreTestSuite struct {
	suite.Suite

	mptStore *MptStore
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (suite *StoreTestSuite) SetupTest() {
	// set exchaind path
	serverDir, err := ioutil.TempDir("", ".exchaind")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(serverDir)
	viper.Set(flags.FlagHome, serverDir)

	mptStore, err := mockMptStore(nil, types.CommitID{})
	if err != nil {
		panic(err)
	}
	for _, key := range commonKeys {
		mptStore.Set([]byte(key), []byte(commonValues[0]))
	}
	for i := 0; i < randKeyNum; i++ {
		key := randBytes(12)
		value := randBytes(32)
		mptStore.Set(key, value)
	}
	mptStore.CommitterCommit(nil)

	suite.mptStore = mptStore
}

func (suite *StoreTestSuite) TestLoadStore() {
	store := suite.mptStore
	key := []byte(commonKeys[0])

	// Create non-pruned height H
	valueH := randBytes(32)
	store.Set(key, valueH)
	cIDH, _ := store.CommitterCommit(nil)

	// Create pruned height Hp
	valueHp := randBytes(32)
	store.Set(key, valueHp)
	cIDHp, _ := store.CommitterCommit(nil)

	// Create current height Hc
	valueHc := randBytes(32)
	store.Set(key, valueHc)
	cIDHc, _ := store.CommitterCommit(nil)

	// Querying an existing store at some previous non-pruned height H
	hStore, err := store.GetImmutable(cIDH.Version)
	suite.Require().NoError(err)
	suite.Require().Equal(hStore.Get(key), valueH)

	// Querying an existing store at some previous pruned height Hp
	hpStore, err := store.GetImmutable(cIDHp.Version)
	suite.Require().NoError(err)
	suite.Require().Equal(hpStore.Get(key), valueHp)

	// Querying an existing store at current height Hc
	hcStore, err := store.GetImmutable(cIDHc.Version)
	suite.Require().NoError(err)
	suite.Require().Equal(hcStore.Get(key), valueHc)
}

func (suite *StoreTestSuite) TestMPTStoreGetSetHasDelete() {
	store := suite.mptStore
	key, originValue := commonKeys[0], commonValues[0]

	exists := store.Has([]byte(key))
	suite.Require().True(exists)

	value := store.Get([]byte(key))
	suite.Require().EqualValues(value, originValue)

	value2 := "notgoodbye"
	store.Set([]byte(key), []byte(value2))

	value = store.Get([]byte(key))
	suite.Require().EqualValues(value, value2)

	store.Delete([]byte(key))
	exists = store.Has([]byte(key))
	suite.Require().False(exists)
}

func (suite *StoreTestSuite) TestMPTStoreNoNilSet() {
	store := suite.mptStore
	suite.Require().Panics(func() { store.Set([]byte("key"), nil) }, "setting a nil value should panic")
}

func (suite *StoreTestSuite) TestGetImmutable() {
	store := suite.mptStore
	key := []byte(commonKeys[0])
	oldValue := store.Get(key)

	newValue := randBytes(32)
	store.Set(key, newValue)
	cID, _ := store.CommitterCommit(nil)

	_, err := store.GetImmutable(cID.Version + 1)
	suite.Require().NoError(err)

	oldStore, err := store.GetImmutable(cID.Version - 1)
	suite.Require().NoError(err)
	suite.Require().Equal(oldStore.Get(key), oldValue)

	curStore, err := store.GetImmutable(cID.Version)
	suite.Require().NoError(err)
	suite.Require().Equal(curStore.Get(key), newValue)

	res := curStore.Query(abci.RequestQuery{Data: key, Height: cID.Version, Path: "/key", Prove: true})
	suite.Require().Equal(res.Value, newValue)
	suite.Require().NotNil(res.Proof)

	suite.Require().Panics(func() { curStore.Set(nil, nil) })
	suite.Require().NotPanics(func() { curStore.Delete(nil) })
	suite.Require().Panics(func() { curStore.CommitterCommit(nil) })
}

func (suite *StoreTestSuite) TestTestIterator() {
	store := suite.mptStore
	iter := store.Iterator(nil, nil)
	i := 0
	for ; iter.Valid(); iter.Next() {
		suite.Require().NotNil(iter.Key())
		suite.Require().NotNil(iter.Value())
		i++
	}

	suite.Require().Equal(len(commonKeys)+randKeyNum, i)
}

func nextVersion(iStore *MptStore) {
	key := []byte(fmt.Sprintf("Key for tree: %d", iStore.LastCommitID().Version))
	value := []byte(fmt.Sprintf("Value for tree: %d", iStore.LastCommitID().Version))
	iStore.Set(key, value)
	iStore.CommitterCommit(nil)
}

func (suite *StoreTestSuite) TestMPTNoPrune() {
	store := suite.mptStore
	nextVersion(store)

	for i := 1; i < 100; i++ {
		for j := 1; j <= i; j++ {
			rootHash := store.GetMptRootHash(uint64(j))
			suite.Require().NotEqual(NilHash, rootHash)
			tire, err := store.db.OpenTrie(rootHash)
			suite.Require().NoError(err)
			suite.Require().NotEqual(EmptyCodeHash, tire.Hash())
		}

		nextVersion(store)
	}
}

func (suite *StoreTestSuite) TestMPTStoreQuery() {
	store := suite.mptStore

	k1, v1 := []byte(commonKeys[0]), []byte(commonValues[0])
	k2, v2 := []byte(commonKeys[1]), []byte(commonValues[1])
	v3 := []byte(commonValues[2])

	cid, _ := store.CommitterCommit(nil)
	ver := cid.Version
	query := abci.RequestQuery{Path: "/key", Data: k1, Height: ver}
	querySub := abci.RequestQuery{Path: "/subspace", Data: []byte("key"), Height: ver}

	// query subspace before anything set
	qres := store.Query(querySub)
	suite.Require().NotEqual(uint32(0), qres.Code)

	// set data
	store.Set(k1, v1)
	store.Set(k2, v2)

	// query data without commit
	qres = store.Query(query)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v1, qres.Value)

	// commit it, but still don't see on old version
	cid, _ = store.CommitterCommit(nil)
	qres = store.Query(query)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v1, qres.Value)

	// but yes on the new version
	query.Height = cid.Version
	qres = store.Query(query)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v1, qres.Value)

	// modify
	store.Set(k1, v3)
	cid, _ = store.CommitterCommit(nil)

	// query will return old values, as height is fixed
	qres = store.Query(query)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v1, qres.Value)

	// update to latest in the query and we are happy
	query.Height = cid.Version
	qres = store.Query(query)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v3, qres.Value)
	query2 := abci.RequestQuery{Path: "/key", Data: k2, Height: cid.Version}

	qres = store.Query(query2)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v2, qres.Value)

	// default (height 0) will show latest
	query0 := abci.RequestQuery{Path: "/key", Data: k1}
	qres = store.Query(query0)
	suite.Require().Equal(uint32(0), qres.Code)
	suite.Require().Equal(v3, qres.Value)
}
