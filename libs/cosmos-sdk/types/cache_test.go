package types

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type mockAccount struct {
	data int
}

func newMockAccount(data int) *mockAccount {
	return &mockAccount{
		data: data,
	}
}

func (m *mockAccount) Copy() interface{} {
	return m
}

func (m *mockAccount) GetAddress() AccAddress {
	return nil
}

func (m *mockAccount) SetAddress(AccAddress) error {
	return nil
}

func (m *mockAccount) GetPubKey() crypto.PubKey {
	return nil
}

func (m *mockAccount) SetPubKey(crypto.PubKey) error {
	return nil
}

func (m *mockAccount) GetAccountNumber() uint64 {
	return uint64(m.data)
}

func (m *mockAccount) SetAccountNumber(uint64) error {
	return nil
}
func (m *mockAccount) GetSequence() uint64 {
	return 0
}
func (m *mockAccount) SetSequence(uint64) error {
	return nil
}
func (m *mockAccount) GetCoins() Coins {
	return nil
}
func (m *mockAccount) SetCoins(Coins) error {
	return nil
}
func (m *mockAccount) SpendableCoins(blockTime time.Time) Coins {
	return nil
}
func (m *mockAccount) String() string {
	return ""
}

func newCache(parent *Cache) *Cache {
	return NewCache(parent, true)
}
func bz(s string) ethcmn.Address { return ethcmn.BytesToAddress([]byte(s)) }

func keyFmt(i int) ethcmn.Address   { return bz(fmt.Sprintf("key%0.8d", i)) }
func accountValueFmt(i int) account { return newMockAccount(i) }

func TestCache(t *testing.T) {
	parent := newCache(nil)
	st := newCache(parent)

	key1Value, _, _ := st.GetAccount(keyFmt(1))
	require.Empty(t, key1Value, "should 'key1' to be empty")

	// put something in mem and in cache
	parent.UpdateAccount(keyFmt(1).Bytes(), accountValueFmt(1), 0, true)
	st.UpdateAccount(keyFmt(1).Bytes(), accountValueFmt(1), 0, true)
	key1Value, _, _ = st.GetAccount(keyFmt(1))
	require.Equal(t, key1Value.GetAccountNumber(), uint64(1))

	// update it in cache, shoudn't change mem
	st.UpdateAccount(keyFmt(1).Bytes(), accountValueFmt(2), 0, true)
	key1ValueInParent, _, _ := parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ := st.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInParent.GetAccountNumber(), uint64(1))
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(2))

	// write it . should change mem
	st.Write(true)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInParent.GetAccountNumber(), uint64(2))
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(2))

	// more writes and checks
	st.Write(true)
	st.Write(true)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInParent.GetAccountNumber(), uint64(2))
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(2))

	// make a new one, check it
	st = newCache(parent)
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(2))

	// make a new one and delete - should not be removed from mem
	st = newCache(parent)
	st.UpdateAccount(keyFmt(1).Bytes(), nil, 0, true)
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	require.Empty(t, key1ValueInSt)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInParent.GetAccountNumber(), uint64(2))

	// Write. should now be removed from both
	st.Write(true)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	require.Empty(t, key1ValueInParent)
	require.Empty(t, key1ValueInSt)
}

func TestCacheNested(t *testing.T) {
	parent := newCache(nil)
	st := newCache(parent)

	// set. check its there on st and not on mem.
	st.UpdateAccount(keyFmt(1).Bytes(), accountValueFmt(1), 0, true)
	key1ValueInParent, _, _ := parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ := st.GetAccount(keyFmt(1))
	require.Empty(t, key1ValueInParent)
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(1))

	// make a new from st and check
	st2 := newCache(st)
	key1ValueInSt, _, _ = st2.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(1))

	// update the value on st2, check it only effects st2
	st2.UpdateAccount(keyFmt(1).Bytes(), accountValueFmt(3), 0, true)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	key1ValueInSt2, _, _ := st2.GetAccount(keyFmt(1))
	require.Empty(t, key1ValueInParent)
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(1))
	require.Equal(t, key1ValueInSt2.GetAccountNumber(), uint64(3))

	// st2 write to its parent, st. doesnt effect parent
	st2.Write(true)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	key1ValueInSt, _, _ = st.GetAccount(keyFmt(1))
	require.Empty(t, key1ValueInParent)
	require.Equal(t, key1ValueInSt.GetAccountNumber(), uint64(3))

	// updates parent
	st.Write(true)
	key1ValueInParent, _, _ = parent.GetAccount(keyFmt(1))
	require.Equal(t, key1ValueInParent.GetAccountNumber(), uint64(3))
}

func BenchmarkCacheKVStoreGetNoKeyFound(b *testing.B) {
	st := newCache(nil)
	b.ResetTimer()
	// assumes b.N < 2**24
	for i := 0; i < b.N; i++ {
		st.GetAccount(ethcmn.BytesToAddress([]byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}))
	}
}

func BenchmarkCacheKVStoreGetKeyFound(b *testing.B) {
	st := newCache(nil)
	for i := 0; i < b.N; i++ {
		arr := []byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}
		st.UpdateAccount(arr, nil, 0, true)
	}
	b.ResetTimer()
	// assumes b.N < 2**24
	for i := 0; i < b.N; i++ {
		st.GetAccount(ethcmn.BytesToAddress([]byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}))
	}
}
