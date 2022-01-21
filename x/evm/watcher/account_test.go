package watcher

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

type mockDuplicateAccount struct {
	*auth.BaseAccount
	Addr byte
	Seq  byte
}

func (a *mockDuplicateAccount) GetAddress() sdk.AccAddress {
	return []byte{a.Addr}
}

func newMockAccount(byteAddr, seq byte) *mockDuplicateAccount {
	ret := &mockDuplicateAccount{Addr: byteAddr, Seq: seq}
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := auth.NewBaseAccount(addr, nil, pubkey, 0, 0)
	ret.BaseAccount = baseAcc
	return ret
}

func localTest(t *testing.T) {
	viper.Set(FlagFastQuery, true)
	rootDir, err := ioutil.TempDir("", fmt.Sprintf("%s-%s_", "test", "test"))
	require.NoError(t, err)
	viper.Set("home", rootDir)
	t.Cleanup(func() {
		os.RemoveAll(rootDir)
	})
}

func TestDuplicateAddAccount(t *testing.T) {
	localTest(t)
	w := NewWatcher(nil)
	a := newMockAccount(1, 1)
	w.SaveAccount(a, false)
	require.Equal(t, 1, len(w.staleBatch))
	w.SaveAccount(a, false)
	require.Equal(t, 1, len(w.staleBatch))
	ind:=w.regionKeySet[regionAccountIndirectly]
	require.NotNil(t, ind)
	require.Equal(t, 1, len(ind))
}

func TestAccountDirectly(t *testing.T) {
	localTest(t)
	w := NewWatcher(nil)
	saveTwoAccounts(t, w)
}

func saveTwoAccounts(t *testing.T, w *Watcher) {
	a := newMockAccount(1, 1)
	w.SaveAccount(a, false)
	indirectly := w.regionKeySet[regionAccountIndirectly]
	require.NotNil(t, indirectly)
	require.Nil(t, w.regionKeySet[regionAccountDirectly])
	require.Equal(t, 1, len(w.staleBatch))
	require.Equal(t, 1, len(indirectly))

	w.SaveAccount(a, true)
	directly := w.regionKeySet[regionAccountDirectly]
	require.NotNil(t, directly)
	require.Equal(t, 1, len(directly))
	require.Equal(t, 1, len(w.batch))
	require.Equal(t, 2, len(w.regionKeySet))
}
func TestAppendAfterNewHeight(t *testing.T) {
	localTest(t)
	w := NewWatcher(nil)
	saveTwoAccounts(t, w)

	w.NewHeight(123, common.BytesToHash([]byte("123")), types.Header{})
	w.Reset()
	require.Equal(t, 0, len(w.regionKeySet))
	require.Equal(t, 0, len(w.batch))

	saveTwoAccounts(t, w)
}

func TestReplaceMoreTimes(t *testing.T) {
	localTest(t)
	w := NewWatcher(nil)

	limit := 100
	for i := 0; i < limit; i++ {
		a1 := newMockAccount(1, byte(i))
		w.SaveAccount(a1, true)
		require.Equal(t, 1, len(w.regionKeySet))
		require.Equal(t, 1, len(w.batch))
	}
	var m map[string]interface{}
	jsonV := w.batch[0].GetValue()
	err := json.Unmarshal([]byte(jsonV), &m)
	require.NoError(t, err)
	require.Equal(t, 99, int(m["Seq"].(float64)))
}
