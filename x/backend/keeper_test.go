package backend

import (
	"fmt"
	"testing"
	"time"

	"github.com/okex/okchain/x/backend/cases"
	"github.com/okex/okchain/x/backend/config"
	"github.com/okex/okchain/x/backend/orm"
	"github.com/okex/okchain/x/dex"
	ordertypes "github.com/okex/okchain/x/order/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	tokentypes "github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestKeeper_AllInOne_Smoke(t *testing.T) {
	app, orders := FireEndBlockerPeriodicMatch(t, true)
	waitInSecond := int(62-time.Now().Second()) % 60
	timer := time.NewTimer(time.Duration(waitInSecond * int(time.Second)))
	select {
	case <-timer.C:
	}

	tickers := app.backendKeeper.GetTickers([]string{}, 100)
	require.True(t, tickers != nil && len(tickers) > 0)

	for _, t := range tickers {
		fmt.Println(t.PrettyString())
	}

	tickers2 := app.backendKeeper.GetTickers([]string{"Not_exist"}, 100)
	require.True(t, tickers2 == nil || len(tickers2) == 0)

	tickers3 := app.backendKeeper.GetTickers([]string{types.TestTokenPair}, 100)
	require.True(t, tickers3 != nil && len(tickers3) == 1)

	candles, _ := app.backendKeeper.GetCandles("not_exists", 60, 100)
	require.True(t, candles == nil || len(candles) == 0)

	candles_no, _ := app.backendKeeper.GetCandles("not_exists", 60, 1001)
	require.True(t, candles_no == nil || len(candles_no) == 0)

	candles1, _ := app.backendKeeper.GetCandles(types.TestTokenPair, 60, 100)
	require.True(t, candles1 != nil || len(candles1) > 0)

	//candles2, _ := app.backendKeeper.GetCandles(types.TestTokenPair, 180, 100)
	//require.True(t, candles2 != nil || len(candles2) > 0)

	ctx := app.NewContext(true, abci.Header{})
	deals, _ := app.backendKeeper.GetDeals(ctx, "nobody", types.TestTokenPair, "", 0, 0, 10, 10)
	require.True(t, deals == nil || len(deals) == 0)

	orders1, cnt1 := app.backendKeeper.GetOrderList(ctx, orders[0].Sender.String(), "", "", true,
		0, 100, 0, 0, false)
	require.Equal(t, cnt1, len(orders1))

	orders2, cnt2 := app.backendKeeper.GetOrderList(ctx, orders[0].Sender.String(), "", "", false,
		0, 100, 0, 0, false)
	require.True(t, orders2 != nil && len(orders2) == cnt2)
	require.True(t, (cnt1+cnt2) == 1)

	_, cnt := app.backendKeeper.GetFeeDetails(ctx, orders[0].Sender.String(), 0, 100)

	//cleanUpInCompletedKline(app.backendKeeper.orm, time.Now().Add(-time.Minute*2).Unix())
	//candles1_, _ := app.backendKeeper.GetCandles(types.TestTokenPair, 60, 100)
	//fmt.Printf("%+v\n", candles1)
	//fmt.Printf("%+v\n", candles1_)
	//require.True(t, candles1_ != nil && len(candles1_) < len(candles1))

	require.True(t, cnt > 0)
}

func TestKeeper_GetCandles(t *testing.T) {
	return

	app, _ := FireEndBlockerPeriodicMatch(t, true)
	time.Sleep(time.Second * 120)

	candles_no, _ := app.backendKeeper.GetCandles("not_exists", 60, 1001)
	require.True(t, candles_no == nil || len(candles_no) == 0)

	candles1, _ := app.backendKeeper.GetCandles(types.TestTokenPair, 60, 100)
	require.True(t, candles1 != nil || len(candles1) >= 2)

}

func TestKeeper_DisableBackend(t *testing.T) {
	app, _ := FireEndBlockerPeriodicMatch(t, false)
	require.Nil(t, app.backendKeeper.Orm)
	require.Nil(t, app.backendKeeper.Cache)
	app.backendKeeper.Stop()
	time.Sleep(time.Second)
}

func TestKeeper_Tx(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 2, true, "")
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{Time: time.Now()}).WithBlockHeight(2)
	feeParams := ordertypes.DefaultParams()
	mapp.orderKeeper.SetParams(ctx, &feeParams)

	tokenPair := dex.GetBuiltInTokenPair()

	mapp.dexKeeper.SaveTokenPair(ctx, tokenPair)

	msgOrderNew := ordertypes.NewMsgNewOrder(addrKeysSlice[0].Address, types.TestTokenPair, types.BuyOrder, "10.0", "1.0")
	msgOrderCancel := ordertypes.NewMsgCancelOrder(addrKeysSlice[0].Address, formatOrderId(2, 1))
	sendCoins, err := sdk.ParseDecCoins("100" + common.TestToken)
	require.Nil(t, err)
	msgSend := tokentypes.NewMsgTokenSend(addrKeysSlice[0].Address, addrKeysSlice[1].Address, sendCoins)

	txs := []auth.StdTx{
		buildTx(mapp, ctx, addrKeysSlice[0], msgOrderNew),
		buildTx(mapp, ctx, addrKeysSlice[0], msgOrderCancel),
		buildTx(mapp, ctx, addrKeysSlice[0], msgSend),
	}

	mockApplyBlock(mapp, ctx, txs)

	ctx = mapp.NewContext(true, abci.Header{})
	getTxs, _ := mapp.backendKeeper.GetTransactionList(ctx, addrKeysSlice[0].Address.String(), 0, 0, 0, 0, 100)
	require.EqualValues(t, 3, len(getTxs))

	getTxs, _ = mapp.backendKeeper.GetTransactionList(ctx, addrKeysSlice[1].Address.String(), 0, 0, 0, 0, 100)
	require.EqualValues(t, 1, len(getTxs))
}

func TestKeeper_CleanUpKlines(t *testing.T) {
	o, _ := orm.MockSqlite3ORM()
	ch := make(chan struct{}, 1)
	conf := config.DefaultConfig()

	cleanUpTime := time.Now().Add(time.Second * 120)
	strClenaUpTime := cleanUpTime.Format("15:04") + ":00"
	conf.CleanUpsTime = strClenaUpTime
	conf.EnableBackend = true
	go CleanUpKlines(ch, o, conf)

	//time.Sleep(121 * time.Second)
}

func sumKlinesVolume(product string, o *orm.ORM, ikline types.IKline) float64 {
	klines, _ := types.NewKlinesFactory(ikline.GetTableName())
	o.GetLatestKlinesByProduct(product, 10000, 0, klines)
	iklines := types.ToIKlinesArray(klines, time.Now().Unix(), false)
	volume := 0.0
	for _, i := range iklines {
		volume += i.GetVolume()
	}

	return volume
}

// TestKeeper_FixJira85 is related to OKDEX-83, OKDEX-85
func TestKeeper_FixJira85(t *testing.T) {
	return
	// FLT Announce : !!! Don't remove the following code!!!!! @wch

	dbDir := cases.GetBackendDBDir()
	mapp, _ := getMockApp(t, 2, true, dbDir)

	timer := time.NewTimer(60 * time.Second)
	select {
	case <-timer.C:
	}

	// 1. TestKline
	product := "btc-235_" + common.NativeToken
	km1Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM1{})
	km3Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM3{})
	km5Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM5{})
	km15Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM15{})
	km360Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM360{})
	km1440Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM1440{})

	require.True(t, km1Sum == km3Sum && km3Sum == km5Sum && km5Sum == km15Sum &&
		km15Sum == km360Sum && km1440Sum == km1Sum && km1Sum == 11.0)

	// 2. TestTicker
	tickers := mapp.backendKeeper.GetTickers(nil, 100)
	require.True(t, tickers != nil && len(tickers) > 1)
	for _, ti := range tickers {
		if ti.Symbol == "btc-235_"+common.NativeToken {
			require.True(t, ti.ChangePercentage == "0.00%")
		}
	}

	// 3. UpdateTickers Again
	ts := time.Now().Unix()
	mapp.backendKeeper.UpdateTickersBuffer(ts-types.SecondsInADay, ts+1, mapp.backendKeeper.Cache.ProductsBuf)
	tickers = mapp.backendKeeper.GetTickers(nil, 100)
	require.True(t, tickers != nil && len(tickers) > 1)
	for _, ti := range tickers {
		if ti.Symbol == "btc-235_"+common.NativeToken {
			require.True(t, ti.ChangePercentage == "0.00%")
		}
	}

}

func TestKeeper_KlineInitialize_RebootTwice(t *testing.T) {
	return
	for i := 0; i < 2; i++ {

		dbDir := cases.GetBackendDBDir()
		mapp, _ := getMockApp(t, 2, true, dbDir)

		timer := time.NewTimer(60 * time.Second)
		select {
		case <-timer.C:
		}

		products := []string{"btc-235_" + common.NativeToken, "atom-564_" + common.NativeToken, "bch-035_" + common.NativeToken}
		expectedSum := []float64{11.0, 10.1, 11.7445}

		for j := 0; j < len(products); j++ {
			product := products[j]
			expSum := expectedSum[j]

			km1Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM1{})
			km3Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM3{})
			km5Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM5{})
			km15Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM15{})
			km30Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM30{})
			km60Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM60{})
			km120Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM120{})
			km240Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM240{})
			km360Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM360{})
			km720Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM720{})
			km1440Sum := sumKlinesVolume(product, mapp.backendKeeper.Orm, &types.KlineM1440{})

			fmt.Println(fmt.Sprintln("Product: ", product, " Expected sum: ", expSum, " Km1Sum: ", km1Sum, "K15Sum", km15Sum, km30Sum, km60Sum, km120Sum, km240Sum, km360Sum, km720Sum, km1440Sum))

			require.True(t, km1Sum == km3Sum && km3Sum == km5Sum && km5Sum == km15Sum && km15Sum == km360Sum &&
				km30Sum == km1Sum && km60Sum == km1Sum && km120Sum == km1Sum && km240Sum == km1Sum && km720Sum == km1Sum &&
				km1440Sum == km1Sum && km1Sum == expSum)

		}
	}
}

func TestKeeper_KlineInitialize_RebootTwice2(t *testing.T) {

	return
	for i := 0; i < 1; i++ {

		dbDir := cases.GetBackendDBDir()
		mapp, _ := getMockApp(t, 2, true, dbDir)

		timer := time.NewTimer(60 * time.Second)
		select {
		case <-timer.C:
		}

		products := []string{"bch-035_" + common.NativeToken, "btc-235_" + common.NativeToken, "atom-564_" + common.NativeToken,
			"dash-150_" + common.NativeToken, "eos-5d4_" + common.NativeToken, "ltc-b72_" + common.NativeToken}
		expectedSum := []float64{12.7445, 11.0, 10.1, 1, 0.45, 2.5099}

		for j := 0; j < len(products); j++ {
			product := products[j]
			expSum := expectedSum[j]

			checkKlinesVolume(t, product, mapp.backendKeeper.Orm, expSum)

		}
	}
}

func checkKlinesVolume(t *testing.T, product string, o *orm.ORM, expSum float64) {

	km1Sum := sumKlinesVolume(product, o, &types.KlineM1{})
	km3Sum := sumKlinesVolume(product, o, &types.KlineM3{})
	km5Sum := sumKlinesVolume(product, o, &types.KlineM5{})
	km15Sum := sumKlinesVolume(product, o, &types.KlineM15{})
	km30Sum := sumKlinesVolume(product, o, &types.KlineM30{})
	km60Sum := sumKlinesVolume(product, o, &types.KlineM60{})
	km120Sum := sumKlinesVolume(product, o, &types.KlineM120{})
	km240Sum := sumKlinesVolume(product, o, &types.KlineM240{})
	km360Sum := sumKlinesVolume(product, o, &types.KlineM360{})
	km720Sum := sumKlinesVolume(product, o, &types.KlineM720{})
	km1440Sum := sumKlinesVolume(product, o, &types.KlineM1440{})

	fmt.Println(fmt.Sprintln("Product: ", product, " Expected sum: ", expSum, " Km1Sum: ", km1Sum, km3Sum, km5Sum, km15Sum, km30Sum, km60Sum, km120Sum, km240Sum, km360Sum, km720Sum, km1440Sum))
	require.True(t, km1Sum == km3Sum && km3Sum == km5Sum && km5Sum == km15Sum && km15Sum == km360Sum &&
		km30Sum == km1Sum && km60Sum == km1Sum && km120Sum == km1Sum && km240Sum == km1Sum && km720Sum == km1Sum &&
		km1440Sum == km1Sum && km1Sum == expSum)

}

func TestKeeper_KlineInitialize_RebootTwice3(t *testing.T) {
	return

	for i := 0; i < 1; i++ {

		dbDir := cases.GetBackendDBDir()
		mapp, _ := getMockApp(t, 2, true, dbDir)

		timer := time.NewTimer(60 * time.Second)
		select {
		case <-timer.C:
		}

		products := []string{"bch-035_" + common.NativeToken, "btc-235_" + common.NativeToken, "atom-564_" + common.NativeToken,
			"dash-150_" + common.NativeToken, "eos-5d4_" + common.NativeToken, "ltc-b72_" + common.NativeToken}
		expectedSum := []float64{11.7445, 11.0, 10.1, 1, 0.45, 2.0099}

		for j := 0; j < len(products); j++ {
			product := products[j]
			expSum := expectedSum[j]
			checkKlinesVolume(t, product, mapp.backendKeeper.Orm, expSum)
		}
	}
}

func TestKeeper_getCandles(t *testing.T) {

	mapp, _ := getMockApp(t, 2, true, "")
	timeMap := GetTimes()
	orm := mapp.backendKeeper.Orm

	k0 := prepareKlineMx("flt_"+common.NativeToken, 60, 0.5, 0.5, 0.5, 0.5, []float64{0.5}, timeMap["-48h"], timeMap["-24h"])
	k1 := prepareKlineMx("flt_"+common.NativeToken, 60, 1, 1, 1, 1, []float64{1}, timeMap["-24h"], timeMap["-15m"])
	k2 := prepareKlineMx("flt_"+common.NativeToken, 60, 2, 2, 2, 2, []float64{2}, timeMap["-15m"], timeMap["-1m"])

	orm.MockCommitKlines(k0, k1, k2)

	endTs := []int64{timeMap["-24h"], timeMap["-15m"], timeMap["-1m"], timeMap["now"] + 120}
	expectedCloses := []string{"0.5000", "1.0000", "2.0000", "2.0000"}
	expectedVolumes := []string{"0.50000000", "1.00000000", "2.00000000", "0.00000000"}
	expectedKlineCount := []int{1, 1000, 1000, 1000}

	for i := 0; i < len(expectedVolumes); i++ {
		restDatas, _ := mapp.backendKeeper.GetCandlesWithTime("flt_"+common.NativeToken, 60, 1000, endTs[i])
		latestKline := restDatas[len(restDatas)-1]
		fmt.Println("[!!!]  ", types.TimeString(endTs[i]), expectedCloses[i], expectedVolumes[i])
		fmt.Printf("[!!!]  %+v\n", latestKline)
		require.True(t, latestKline != nil)
		require.True(t, latestKline[4] == expectedCloses[i])
		require.True(t, latestKline[5] == expectedVolumes[i])
		require.True(t, len(restDatas) == expectedKlineCount[i])
	}
}

func formatOrderId(blockHeight, orderNum int64) string {
	format := "ID%010d-%d"
	if blockHeight > 9999999999 {
		format = "ID%d-%d"
	}
	return fmt.Sprintf(format, blockHeight, orderNum)
}
