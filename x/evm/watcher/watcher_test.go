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
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

const addrHex = "0x756F45E3FA69347A9A973A725E3C98bC4db0b4c1"

type WatcherTestw struct {
	ctx     sdk.Context
	querier sdk.Querier
	app     *app.OKExChainApp
	stateDB *types.CommitStateDB
	address ethcmn.Address
	handler    sdk.Handler
}


func SetupTest() *WatcherTestw {
	w := &WatcherTestw{}
	checkTx := false
	viper.Set(watcher.FlagFastQuery, true)
	viper.Set(watcher.FlagDBBackend, "memdb")
	w.app = app.Setup(checkTx)
	w.ctx = w.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
//	w.stateDB = types.CreateEmptyCommitStateDB(w.app.EvmKeeper.GenerateCSDBParams(), w.ctx)
	w.handler = evm.NewHandler(w.app.EvmKeeper)
//	w.querier = keeper.NewQuerier(*w.app.EvmKeeper)
//	w.address = ethcmn.HexToAddress(addrHex)

	params := types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	w.app.EvmKeeper.SetParams(w.ctx, params)

	return w
}



func TestHandleMsgEthereumTx(t *testing.T) {
	w := SetupTest()
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
			w = SetupTest() // reset
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

			// produce WatchData
			w.app.EvmKeeper.Watcher.Commit()
			time.Sleep(time.Second * 2)
			db := watcher.InstanceOfWatchStore().GetDB()
			fmt.Println(db.Stats())
			it, _ := db.Iterator(nil, nil)
			for ;it.Valid(); {
				fmt.Println(string(it.Key()), string(it.Value()))
				db.Delete(it.Key())
				it.Next()
			}

			// get WatchData
			wd, err := state.GetWD()
			require.Nil(t, err)
			fmt.Println(string(wd))

			// use WatchData
			//db2 := dbm.NewMemDB()
			//w.app.EvmKeeper.Watcher.SetStore(db2)
			w = SetupTest()
			db2 := watcher.InstanceOfWatchStore().GetDB()
			fmt.Println(db2.Stats())
			state.UseWD(wd)
			time.Sleep(time.Second * 2)
			fmt.Println(db2.Stats())
			it, _ = db2.Iterator(nil, nil)
			for ;it.Valid(); {
				fmt.Println(string(it.Key()), string(it.Value()))
				it.Next()
			}

		})
	}
}

func TestMsgEthermint(t *testing.T) {
	var (
		tx   types.MsgEthermint
		from = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		to   = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	)
	w := SetupTest()
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
			w = SetupTest() // reset
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

			// produce WatchData
			w.app.EvmKeeper.Watcher.Commit()
			time.Sleep(time.Second * 2)
			db := watcher.InstanceOfWatchStore().GetDB()
			fmt.Println(db.Stats())
			it, _ := db.Iterator(nil, nil)
			for ;it.Valid(); {
				fmt.Println(string(it.Key()), string(it.Value()))
				db.Delete(it.Key())
				it.Next()
			}

			// get WatchData
			wd, err := state.GetWD()
			require.Nil(t, err)
			fmt.Println(string(wd))

			// use WatchData
			//db2 := dbm.NewMemDB()
			//w.app.EvmKeeper.Watcher.SetStore(db2)
			w = SetupTest()
			db2 := watcher.InstanceOfWatchStore().GetDB()
			fmt.Println(db2.Stats())
			state.UseWD(wd)
			time.Sleep(time.Second * 2)
			fmt.Println(db2.Stats())
			it, _ = db2.Iterator(nil, nil)
			for ;it.Valid(); {
				fmt.Println(string(it.Key()), string(it.Value()))
				it.Next()
			}
		})
	}
}