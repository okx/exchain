package backend

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/backend/cache"
	"github.com/okex/okexchain/x/backend/orm"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/common"
	"github.com/stretchr/testify/assert"
)

func prepareKlineMx(product string, refreshInterval int, open, close, low, high float64, volumes []float64, startTS, endTS int64) []interface{} {

	destKName := types.GetKlineTableNameByFreq(refreshInterval)
	destK := types.MustNewKlineFactory(destKName, nil)
	destIKline := destK.(types.IKline)

	klines := []interface{}{}
	for i := 0; i < len(volumes); i++ {
		ts := destIKline.GetAnchorTimeTS(endTS - int64(destIKline.GetFreqInSecond()*i))

		b := types.BaseKline{
			Product:   product,
			High:      high,
			Low:       low,
			Volume:    volumes[i],
			Timestamp: ts,
			Open:      open,
			Close:     close,
		}

		newDestK, _ := types.NewKlineFactory(destIKline.GetTableName(), &b)
		klines = append(klines, newDestK)
	}

	return klines
}

func prepareMatches(product string, prices []float64, quantities []float64, endTS int64) []*types.MatchResult {

	matchResults := make([]*types.MatchResult, 0, len(prices))
	for i := 0; i < len(prices); i++ {
		match := &types.MatchResult{
			Timestamp:   endTS - int64(len(prices)) + int64(i),
			BlockHeight: endTS + int64(i),
			Product:     product,
			Price:       prices[i],
			Quantity:    quantities[i],
		}

		matchResults = append(matchResults, match)
	}

	return matchResults
}

func aTicker(product string, open, close, high, low, price, volume float64) *types.Ticker {
	t := types.Ticker{
		Timestamp: time.Now().Unix(),
		Product:   product,
		Open:      open,
		Close:     close,
		High:      high,
		Low:       low,
		Price:     price,
		Volume:    volume,
		Symbol:    product,
	}
	return &t
}

func GetTimes() map[string]int64 {

	timeMap := map[string]int64{}
	strNow := "2019-06-10 15:40:59"
	tm, _ := time.Parse("2006-01-02 15:04:05", strNow)
	nowTS := tm.Unix()
	timeMap["now"] = nowTS
	timeMap["-21d"] = nowTS - (types.SecondsInADay * 21)
	timeMap["-14d"] = nowTS - (types.SecondsInADay * 14)
	timeMap["-48h"] = nowTS - (types.SecondsInADay * 2)
	timeMap["-24h"] = nowTS - (types.SecondsInADay * 1)
	timeMap["-60m"] = nowTS - (60 * 60)
	timeMap["-30m"] = nowTS - (60 * 30)
	timeMap["-15m"] = nowTS - (60 * 15)
	timeMap["-5m"] = nowTS - (60 * 5)
	timeMap["-2m"] = nowTS - (60 * 2)
	timeMap["-1m"] = nowTS - (60 * 1)

	return timeMap
}

func baseCaseRunner(t *testing.T, product string, productBuffer []string, startTS, endTS int64,
	kline15s []interface{}, kline1s []interface{}, matches []*types.MatchResult, expectedTicker *types.Ticker, fakeLatestTickers *map[string]*types.Ticker, orm *orm.ORM, doCreate bool) error {
	// 1. Prepare Datas
	//defer func() {
	//	r := recover()
	//	if r != nil {
	//		assert.True(t, false)
	//	}
	//}()

	if doCreate {
		if len(matches) > 0 {
			_, err := orm.AddMatchResults(matches)
			require.Nil(t, err)
		}
		orm.CommitKlines(kline15s, kline1s)
	}

	// 2. UpdateTickerBuffer
	keeper := Keeper{Orm: orm, Cache: cache.NewCache()}
	keeper.Cache.LatestTicker = *fakeLatestTickers
	keeper.UpdateTickersBuffer(startTS, endTS, productBuffer)

	// 3. CheckResults
	gotTicker := keeper.Cache.LatestTicker[product]
	if gotTicker != nil {
		fmt.Println(fmt.Sprintf("   Got: %s", gotTicker.PrettyString()))
	}

	if expectedTicker != nil {
		fmt.Println(fmt.Sprintf("Expect: %s", expectedTicker.PrettyString()))
	}

	if expectedTicker == nil {
		assert.True(t, gotTicker == nil && nil == expectedTicker)
	} else {
		assert.Equal(t, gotTicker.Price, expectedTicker.Price, gotTicker.PrettyString(), expectedTicker.PrettyString())
		assert.Equal(t, gotTicker.Product, expectedTicker.Product, gotTicker.PrettyString(), expectedTicker.PrettyString())
		assert.Equal(t, gotTicker.Open, expectedTicker.Open, gotTicker.PrettyString(), expectedTicker.PrettyString())
		assert.Equal(t, gotTicker.Close, expectedTicker.Close, gotTicker.PrettyString(), expectedTicker.PrettyString())
		assert.Equal(t, gotTicker.High, expectedTicker.High, gotTicker.PrettyString(), expectedTicker.PrettyString())
		assert.Equal(t, gotTicker.Low, expectedTicker.Low, gotTicker.PrettyString(), expectedTicker.PrettyString())
		assert.Equal(t, gotTicker.Volume, expectedTicker.Volume, gotTicker.PrettyString(), expectedTicker.PrettyString())
	}

	return nil
}

func simpleCaseRunner(t *testing.T, product string, productBuffer []string, startTS, endTS int64,
	kline15s []interface{}, kline1s []interface{}, matches []*types.MatchResult, expectedTicker *types.Ticker, fakeLatestTickers *map[string]*types.Ticker) (err error) {
	o, dbPath := orm.MockSqlite3ORM()
	defer orm.DeleteDB(dbPath)

	return baseCaseRunner(t, product, productBuffer, startTS, endTS, kline15s, kline1s, matches, expectedTicker, fakeLatestTickers, o, true)
}

func TestTicker_S1(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := []interface{}{}
	kline1s := []interface{}{}
	matches := []*types.MatchResult{}
	fakeLatestTickers := &map[string]*types.Ticker{}

	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"], kline15s, kline1s, matches, nil, fakeLatestTickers)

	assert.True(t, err == nil)
}

func TestTicker_S2(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := []interface{}{}
	kline1s := []interface{}{}
	matches := prepareMatches(product, []float64{100.0}, []float64{2.0}, timeMap["-48h"])
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, d := range matches {
		fmt.Println(d)
	}
	err := simpleCaseRunner(t, product, []string{product}, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 100.0, 100.0, 100.0, 100.0, 0), fakeLatestTickers)

	assert.True(t, err == nil)
}

func TestTicker_S3(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := []interface{}{}
	kline1s := []interface{}{}
	matches := prepareMatches(product, []float64{100.0}, []float64{2.0}, timeMap["now"])
	fakeLatestTickers := &map[string]*types.Ticker{}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 100.0, 100.0, 100.0, 100.0, 2.0), fakeLatestTickers)

	assert.True(t, err == nil)
}

func TestTicker_S4(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := []interface{}{}
	kline1s := []interface{}{}
	matches := prepareMatches(product, []float64{100.0, 101.0}, []float64{2.0, 3.0}, timeMap["now"])

	fakeLatestTickers := &map[string]*types.Ticker{}

	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 101.0, 101.0, 100.0, 101.0, 5.0), fakeLatestTickers)

	assert.True(t, err == nil)
}

func TestTicker_S5(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-24h"], timeMap["-30m"])
	kline1s := []interface{}{}
	matches := []*types.MatchResult{}
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline15s {
		fmt.Println(k)
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 100.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S6(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := []interface{}{}
	kline1s := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-15m"], timeMap["-2m"])
	matches := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["-2m"])
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline1s {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 100.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S7(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-24h"], timeMap["-30m"])
	kline1s := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-30m"], timeMap["-2m"])
	matches := []*types.MatchResult{}
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline1s {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 200.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S8(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-24h"], timeMap["-30m"])
	kline1s := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-30m"], timeMap["-2m"])
	matches := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["now"])
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline1s {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 300.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S9(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-24h"], timeMap["-15m"])

	kline1sB := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-30m"], timeMap["-15m"])
	kline1sA := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-15m"], timeMap["-2m"])
	kline1sA = append(kline1sA, kline1sB...)

	matches := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["now"])
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline1sA {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1sA, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 300.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S10(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-24h"], timeMap["-15m"])

	kline1sB := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-30m"], timeMap["-15m"])
	kline1sA := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-15m"], timeMap["-2m"])
	kline1sA = append(kline1sA, kline1sB...)

	matchesA := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["now"])
	matchesB := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["-2m"])
	matchesA = append(matchesA, matchesB...)

	fakeLatestTickers := &map[string]*types.Ticker{}
	for _, k := range kline1sA {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1sA, matchesA,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 300.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S11(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15sA := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-48h"], timeMap["-24m"])
	kline15sB := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-24h"], timeMap["-15m"])
	kline15sA = append(kline15sA, kline15sB...)

	kline1sB := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-30m"], timeMap["-15m"])
	kline1sA := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-15m"], timeMap["-2m"])
	kline1sA = append(kline1sA, kline1sB...)

	matchesA := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["now"])
	matchesB := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["-2m"])
	matchesA = append(matchesA, matchesB...)

	fakeLatestTickers := &map[string]*types.Ticker{}
	for _, k := range kline1sA {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15sA, kline1sA, matchesA,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 300.0), fakeLatestTickers)

	assert.True(t, err == nil)

}

func TestTicker_S12(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15sA := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-48h"], timeMap["-24h"])
	kline15sB := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100, 200}, timeMap["-24h"], timeMap["-15m"])
	kline15sA = append(kline15sA, kline15sB...)

	kline1sB := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-30m"], timeMap["-15m"])
	kline1sA := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100, 200}, timeMap["-15m"], timeMap["-2m"])
	kline1sA = append(kline1sA, kline1sB...)

	matchesA := prepareMatches(product, []float64{100.0, 99.0}, []float64{25, 25}, timeMap["now"])
	matchesB := prepareMatches(product, []float64{100.0}, []float64{25}, timeMap["-2m"])
	matchesA = append(matchesA, matchesB...)
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline15sA {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}

	for _, k := range kline1sA {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}

	for _, d := range matchesA {
		fmt.Println(d)
	}

	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15sA, kline1sA, matchesA,
		aTicker(product, 100.0, 99.0, 210.0, 99.0, 99.0, 750.0), fakeLatestTickers)
	assert.True(t, err == nil)

}

func TestTicker_S13(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 220.0, 99.0, 220.0, []float64{100}, timeMap["-48h"], timeMap["-24h"])
	kline1s := []interface{}{}
	matches := prepareMatches(product, []float64{100.0, 99.0, 210.0, 220.0}, []float64{25, 25, 25, 25}, timeMap["-24h"])
	fakeLatestTickers := &map[string]*types.Ticker{}

	for _, k := range kline1s {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 220.0, 220.0, 99.0, 220.0, 100.0), fakeLatestTickers)
	assert.True(t, err == nil)

	err = simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+15*60, kline15s, kline1s, matches,
		aTicker(product, 220.0, 220.0, 220.0, 220.0, 220.0, 0), fakeLatestTickers)
	assert.True(t, err == nil)
}

func TestTicker_S14(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	// open
	kline15s1 := prepareKlineMx(product, 15*60, 100.0, 100.0, 100.0, 100.0, []float64{100}, timeMap["-24h"], timeMap["-30m"]-60*15*5)
	kline15s2 := prepareKlineMx(product, 15*60, 99.0, 99.0, 99.0, 99.0, []float64{100}, timeMap["-24h"], timeMap["-30m"]-60*15*4)
	kline15s3 := prepareKlineMx(product, 15*60, 220.0, 220.0, 220.0, 220.0, []float64{100}, timeMap["-24h"], timeMap["-30m"]-60*15*3)
	kline15s4 := prepareKlineMx(product, 15*60, 99.0, 99.0, 99.0, 99.0, []float64{100}, timeMap["-24h"], timeMap["-30m"]-60*15*2)

	// close
	kline15s5 := prepareKlineMx(product, 15*60, 98.0, 98.0, 98.0, 98.0, []float64{100}, timeMap["-24h"], timeMap["-30m"]-60*15*1)
	kline1s := []interface{}{}
	matches := []*types.MatchResult{}
	fakeLatestTickers := &map[string]*types.Ticker{}

	klines15s := []interface{}{}
	klines15s = append(klines15s, kline15s1...)
	klines15s = append(klines15s, kline15s2...)
	klines15s = append(klines15s, kline15s3...)
	klines15s = append(klines15s, kline15s4...)
	klines15s = append(klines15s, kline15s5...)

	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, klines15s, kline1s, matches,
		aTicker(product, 100.0, 98.0, 220.0, 98.0, 98.0, 500.0), fakeLatestTickers)
	assert.True(t, err == nil)

}

func TestTicker_C1(t *testing.T) {

	latestTickers := map[string]*types.Ticker{}
	latestTickers["not_exist"] = aTicker("not_exist", 100.0, 230.0, 220.0, 99.0, 230.0, 100.0)

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 220.0, 99.0, 220.0, []float64{100}, timeMap["-48h"], timeMap["-24h"])
	kline1s := []interface{}{}
	matches := prepareMatches(product, []float64{100.0, 99.0, 210.0, 220.0}, []float64{25, 25, 25, 25}, timeMap["-24h"])

	for _, k := range kline1s {
		fmt.Println(k.(types.IKline).PrettyTimeString())
	}
	err := simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 220.0, 220.0, 99.0, 220.0, 100.0), &latestTickers)
	assert.True(t, err == nil)

	err = simpleCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+15*60, kline15s, kline1s, matches,
		aTicker(product, 220.0, 220.0, 220.0, 220.0, 220.0, 0), &latestTickers)
	assert.True(t, err == nil)

	oldTicker := latestTickers["not_exist"]
	assert.True(t, oldTicker.Open == 230.0)
	assert.True(t, oldTicker.Close == 230.0)
	assert.True(t, oldTicker.High == 230.0)
	assert.True(t, oldTicker.Low == 230.0)
	assert.True(t, oldTicker.Price == 230.0)
	assert.True(t, oldTicker.Volume == 0)
}

func TestTicker_C3(t *testing.T) {

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-14d"]-types.SecondsInADay, timeMap["-14d"])
	kline1s := prepareKlineMx(product, 60, 100.0, 200.0, 99.0, 210.0, []float64{100}, timeMap["-14d"]-types.SecondsInADay, timeMap["-14d"])
	matches := prepareMatches(product, []float64{100.0, 99.0, 210.0, 200.0}, []float64{25, 25, 25, 25}, timeMap["-14d"])
	fakeLatestTickers := &map[string]*types.Ticker{}

	o, dbPath := orm.MockSqlite3ORM()
	defer orm.DeleteDB(dbPath)

	err := baseCaseRunner(t, product, nil, timeMap["-21d"], timeMap["-21d"]+1, kline15s, kline1s, matches,
		nil, fakeLatestTickers, o, true)
	assert.True(t, err == nil)

	err = baseCaseRunner(t, product, nil, timeMap["-14d"], timeMap["-14d"]+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 100.0), fakeLatestTickers, o, false)
	assert.True(t, err == nil)

	var SecondInAMinute int64 = 60
	err = baseCaseRunner(t, product, nil, timeMap["-14d"]+SecondInAMinute*2, timeMap["-14d"]+SecondInAMinute*2+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 100.0), fakeLatestTickers, o, false)
	assert.True(t, err == nil)

	err = baseCaseRunner(t, product, nil, timeMap["-14d"]+SecondInAMinute*15, timeMap["-14d"]+SecondInAMinute*15+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 100.0), fakeLatestTickers, o, false)
	assert.True(t, err == nil)

	err = baseCaseRunner(t, product, nil, timeMap["-14d"]+types.SecondsInADay, timeMap["-14d"]+types.SecondsInADay+1, kline15s, kline1s, matches,
		aTicker(product, 100.0, 200.0, 210.0, 99.0, 200.0, 100.0), fakeLatestTickers, o, false)
	assert.True(t, err == nil)

	err = baseCaseRunner(t, product, nil, timeMap["-14d"]+types.SecondsInADay+15*SecondInAMinute,
		timeMap["-14d"]+types.SecondsInADay+15*SecondInAMinute+1, kline15s, kline1s, matches,
		aTicker(product, 200.0, 200.0, 200.0, 200.0, 200.0, 0.0), fakeLatestTickers, o, false)
	assert.True(t, err == nil)

	err = baseCaseRunner(t, product, nil, timeMap["now"], timeMap["now"]+1, kline15s, kline1s, matches,
		aTicker(product, 200.0, 200.0, 200.0, 200.0, 200.0, 0.0), fakeLatestTickers, o, false)
	assert.True(t, err == nil)

}

func TestTicker_C4(t *testing.T) {
	//return

	product := "btc_" + common.NativeToken
	timeMap := GetTimes()
	kline15s := prepareKlineMx(product, 15*60, 100.0, 100.0, 100.0, 100.0, []float64{100}, timeMap["-24h"], timeMap["-60m"])
	kline1s1 := prepareKlineMx(product, 60, 40.0, 40.0, 40.0, 40.0, []float64{40, 60}, timeMap["-15m"], timeMap["-5m"])
	kline1s2 := prepareKlineMx(product, 60, 98.0, 99.0, 98.0, 99.0, []float64{100}, timeMap["-5m"], timeMap["now"])
	matches := prepareMatches(product, []float64{98.0, 99.0}, []float64{98, 2}, timeMap["now"])
	var kline1s []interface{}
	kline1s = append(kline1s, kline1s1...)
	kline1s = append(kline1s, kline1s2...)

	fakeLatestTickers := &map[string]*types.Ticker{}

	orm, dbPath := orm.MockSqlite3ORM()
	//defer DeleteDB(dbPath)

	err := baseCaseRunner(t, product, nil, timeMap["-21d"], timeMap["-21d"]+1, kline15s, kline1s, matches,
		nil, fakeLatestTickers, orm, true)
	assert.True(t, err == nil)

	tickerInNext1M := aTicker(product, 100.0, 99.0, 100.0, 40.0, 99.0, 300.0)
	for i := 1; i <= 60; i++ {
		err = baseCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+int64(i), kline15s, kline1s, matches, tickerInNext1M, fakeLatestTickers, orm, false)
		require.Nil(t, err)
	}

	tickerInNext1M2M := tickerInNext1M
	for i := 61; i <= 120; i++ {
		err = baseCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+int64(i), kline15s, kline1s, matches, tickerInNext1M2M, fakeLatestTickers, orm, false)
		require.Nil(t, err)

	}

	tickerInNext2M5M := tickerInNext1M
	for i := 121; i <= 300; i++ {
		err = baseCaseRunner(t, product, nil, timeMap["-24h"], timeMap["now"]+int64(i), kline15s, kline1s, matches, tickerInNext2M5M, fakeLatestTickers, orm, false)
		require.Nil(t, err)
	}

	tickerInNext5M15M := tickerInNext1M
	for i := 301; i <= 900; i++ {
		expectTS := timeMap["now"] + int64(i)
		tickerInNext5M15M.Timestamp = expectTS
		err = baseCaseRunner(t, product, nil, timeMap["-24h"], expectTS, kline15s, kline1s, matches, tickerInNext5M15M, fakeLatestTickers, orm, false)
		require.Nil(t, err)
	}

	fmt.Println(dbPath)
}
