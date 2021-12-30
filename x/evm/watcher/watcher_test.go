package watcher_test

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"
	"math/big"
	"testing"
	"time"
)
type KV struct {
	k []byte
	v []byte
}

func calcHash(kvs []KV) []byte {
	ha := tmhash.New()
	// calc a hash
	for _, kv := range kvs {
		ha.Write(kv.k)
		ha.Write(kv.v)
	}
	return ha.Sum(nil)
}

type WatcherTestSt struct {
	ctx     sdk.Context
	app     *app.OKExChainApp
	handler sdk.Handler
}

func setupTest() *WatcherTestSt {
	w := &WatcherTestSt{}
	checkTx := false
	viper.Set(watcher.FlagFastQuery, true)
	viper.Set(watcher.FlagDBBackend, "memdb")
	w.app = app.Setup(checkTx)
	w.ctx = w.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
	w.handler = evm.NewHandler(w.app.EvmKeeper)

	params := types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	w.app.EvmKeeper.SetParams(w.ctx, params)

	return w
}

func getDBKV(t *testing.T, db dbm.DB) []KV {
	var kvs []KV
	fmt.Println(db.Stats())
	it, _ := db.Iterator(nil, nil)
	for it.Valid() {
		kvs = append(kvs, KV{it.Key(), it.Value()})
		err := db.Delete(it.Key())
		require.Nil(t, err)
		it.Next()
	}
	return kvs
}

func testWatchData(t *testing.T, w *WatcherTestSt) {
	// produce WatchData
	w.app.EvmKeeper.Watcher.Commit()
	time.Sleep(time.Second * 2)
	db := watcher.InstanceOfWatchStore().GetDB()
	pWd := getDBKV(t, db)

	// get WatchData
	wd, err := state.GetWD()
	require.Nil(t, err)
	require.NotEmpty(t, wd)

	db2 := watcher.InstanceOfWatchStore().GetDB()
	fmt.Println(db2.Stats())

	// use WatchData
	state.UseWD(wd)
	time.Sleep(time.Second * 5)
	cWd := getDBKV(t, db2)

	// compare db_kv of producer and consumer
	require.Equal(t, pWd, cWd)
	pHash := calcHash(pWd)
	cHash := calcHash(cWd)
	require.NotEmpty(t, pHash)
	require.NotEmpty(t, cHash)
	require.Equal(t, pHash, cHash)
}

func TestHandleMsgEthereumTx(t *testing.T) {
	w := setupTest()
	privkey, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	sender := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	var tx types.MsgEthereumTx

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				w.app.EvmKeeper.SetBalance(w.ctx, sender, big.NewInt(100))
				tx = types.NewMsgEthereumTx(0, &sender, big.NewInt(100), 3000000, big.NewInt(1), nil)

				// parse context chain ID to big.Int
				chainID, err := ethermint.ParseChainID(w.ctx.ChainID())
				require.NoError(t, err)

				// sign transaction
				err = tx.Sign(chainID, privkey.ToECDSA())
				require.NoError(t, err)
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.msg, func(t *testing.T) {
			w = setupTest() // reset
			//nolint
			tc.malleate()
			w.ctx = w.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			res, err := w.handler(w.ctx, tx)

			//nolint
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, res)
				var expectedConsumedGas uint64 = 21000
				require.EqualValues(t, expectedConsumedGas, w.ctx.GasMeter().GasConsumed())
			} else {
				require.Error(t, err)
				require.Nil(t, res)
			}

			testWatchData(t, w)
		})
	}
}

func TestMsgEthermint(t *testing.T) {
	var (
		tx   types.MsgEthermint
		from = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		to   = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	)
	w := setupTest()
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				tx = types.NewMsgEthermint(0, &to, sdk.NewInt(1), 100000, sdk.NewInt(2), []byte("test"), from)
				w.app.EvmKeeper.SetBalance(w.ctx, ethcmn.BytesToAddress(from.Bytes()), big.NewInt(100))
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w = setupTest() // reset
			//nolint
			tc.malleate()
			w.ctx = w.ctx.WithIsCheckTx(true)
			w.ctx = w.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			res, err := w.handler(w.ctx, tx)

			//nolint
			if tc.expPass {
				require.NoError(t, err)
				require.NotNil(t, res)
				var expectedConsumedGas uint64 = 21064
				require.EqualValues(t, expectedConsumedGas, w.ctx.GasMeter().GasConsumed())
			} else {
				require.Error(t, err)
				require.Nil(t, res)
			}
		})
	}
}
