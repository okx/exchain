package types_test

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	lru "github.com/hashicorp/golang-lru"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

func (suite *StateDBTestSuite) TestStateObject_State() {
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

		value := suite.stateObject.GetState(nil, tc.key)
		suite.Require().Equal(tc.expValue, value, tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateObject_AddBalance() {
	testCase := []struct {
		name       string
		amount     *big.Int
		expBalance *big.Int
	}{
		{"zero amount", big.NewInt(0), big.NewInt(0)},
		{"positive amount", big.NewInt(10), big.NewInt(10)},
		{"negative amount", big.NewInt(-1), big.NewInt(9)},
	}

	for _, tc := range testCase {
		suite.stateObject.AddBalance(tc.amount)
		suite.Require().Equal(tc.expBalance, suite.stateObject.Balance(), tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateObject_SubBalance() {
	testCase := []struct {
		name       string
		amount     *big.Int
		expBalance *big.Int
	}{
		{"zero amount", big.NewInt(0), big.NewInt(0)},
		{"negative amount", big.NewInt(-10), big.NewInt(10)},
		{"positive amount", big.NewInt(1), big.NewInt(9)},
	}

	for _, tc := range testCase {
		suite.stateObject.SubBalance(tc.amount)
		suite.Require().Equal(tc.expBalance, suite.stateObject.Balance(), tc.name)
	}
}

func (suite *StateDBTestSuite) TestStateObject_Code() {
	testCase := []struct {
		name     string
		expCode  []byte
		malleate func()
	}{
		{
			"cached code",
			[]byte("code"),
			func() {
				suite.stateObject.SetCode(ethcmn.BytesToHash([]byte("code_hash")), []byte("code"))
			},
		},
		{
			"empty code hash",
			nil,
			func() {
				suite.stateObject.SetCode(ethcmn.Hash{}, nil)
			},
		},
		{
			"empty code",
			nil,
			func() {
				suite.stateObject.SetCode(ethcmn.BytesToHash([]byte("code_hash")), nil)
			},
		},
	}

	for _, tc := range testCase {
		tc.malleate()

		code := suite.stateObject.Code(nil)
		suite.Require().Equal(tc.expCode, code, tc.name)
	}
}

func TestDefaultGenesisState(t *testing.T) {
	key := []byte("test1")
	value := ethcrypto.Keccak256Hash(key)
	hash1 := getHash([]byte("test2"))
	tmtypes.HashPool.Put(&hash1)
	hash := getHash(key)
	require.Equal(t, value, hash)

	hash1 = getHash([]byte("test2-more"))
	tmtypes.HashPool.Put(&hash1)
	hash = getHash(key)
	require.Equal(t, value, hash)

	hash1 = getHash([]byte("test2-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	tmtypes.HashPool.Put(&hash1)
	hash = getHash(key)
	require.Equal(t, value, hash)

	hash1 = getHash([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	tmtypes.HashPool.Put(&hash1)
	hash = getHash(key)
	require.Equal(t, value, hash)

	result := string(key)
	result1 := tmtypes.ByteSliceToStr(key)
	require.Equal(t, result1, result)

	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	ethAddr := ethcrypto.PubkeyToAddress(priv.PublicKey)
	prefix := ethAddr.Bytes()
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)

	compositeKeyGC := getKey(prefix, key)
	require.Equal(t, compositeKey, compositeKeyGC)

	testHash := []byte("defer tmtypes.BytesPool.Put(&compositeKey)defer tmtypes.BytesPool.Put(&compositeKey)")
	tmtypes.BytesPool.Put(&testHash)
	compositeKeyGC = getKey(prefix, key)
	require.Equal(t, compositeKey, compositeKeyGC)

}

func getKey(prefix, key []byte) (hash []byte) {
	bp := tmtypes.BytesPool.Get().(*[]byte)
	compositeKey := *bp
	compositeKey = compositeKey[:0]
	compositeKey = append(compositeKey, prefix...)
	compositeKey = append(compositeKey, key...)
	return compositeKey
}

func getHash(key []byte) ethcmn.Hash {
	crytoState := tmtypes.EthCryptoState.Get().(ethcrypto.KeccakState)
	defer tmtypes.EthCryptoState.Put(crytoState)
	crytoState.Reset()
	crytoState.Write(key)

	hashP := tmtypes.HashPool.Get().(*ethcmn.Hash)
	hash := *hashP
	crytoState.Read(hash[:])
	return hash
}

func (suite *StateDBTestSuite) TestStateObject_GetSate() {
	suite.stateObject.SetState(nil, ethcmn.BytesToHash([]byte("key")), ethcmn.Hash{})

}

func TestKeccak256HashWithCache(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	ethAddr := ethcrypto.PubkeyToAddress(priv.PublicKey)
	key, err := hexutil.Decode("0x2B2641734D81a6B93C9aE1Ee6290258FB6666921")
	require.NoError(t, err)

	compositeKey := make([]byte, 0)
	compositeKey = append(compositeKey, ethAddr.Bytes()...)
	compositeKey = append(compositeKey, key...)
	hash1 := types.Keccak256HashWithCache(compositeKey)
	hash2 := types.Keccak256HashWithCacheNew(compositeKey)
	require.Equal(t, hash2, hash1)

	cache, _ := lru.NewARC(1)
	lifei := []byte("aaa-test")
	test := tmtypes.ByteSliceToStr(lifei)
	keyValue := []byte("bbb-value")
	temp := &keyValue
	cache.Add(test, keyValue)
	value, ok := cache.Get("aaa-test")
	t.Log("result1", string(value.([]byte)), ok)
	//t.Log("startStr", startStr)
	////value, ok := cache.Get(startStr)
	////t.Log("result0", value, ok)
	////temp := &lifei
	////	(*temp)[0] = 105
	////	(*temp)[1] = 105
	////	(*temp)[2] = 105
	//lifei[0] = 105
	//lifei[1] = 105
	//lifei[2] = 105
	////t.Log("test", strings.Compare(startStr, *test))
	//t.Log("testStr", *test, "startStr", startStr)
	//test1 := tmtypes.UnsafeToString(lifei)
	////t.Log("test", strings.Compare(startStr, *test1))
	//t.Log("testStr", *test1, "startStr", startStr)
	//
	//value, ok := cache.Get("aaa-lifei")
	//t.Log("result1", value, ok)
	//value, ok = cache.Get("iii-lifei")
	//t.Log("result2", value, ok)
	(*temp)[0] = 105
	(*temp)[1] = 105
	(*temp)[2] = 105
	value, ok = cache.Get("aaa-test")
	t.Log("result1", string(value.([]byte)), ok)
}
