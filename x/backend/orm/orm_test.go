package orm

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcommon "github.com/tendermint/tendermint/libs/common"
)

func TestGorm(t *testing.T) {

	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("%+v", r)
			debug.PrintStack()
		}
	}()

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer DeleteDB("test.db")
	defer db.Close()
	//db.LogMode(true)
	p := sdk.NewDecWithPrec(1, 2)

	fp, _ := strconv.ParseFloat(p.String(), 64)

	d1 := types.Deal{
		BlockHeight: 1, OrderID: "order0", Product: "abc_bcd", Price: fp, Quantity: 100,
		Sender: "asdlfkjsd", Side: types.SellOrder, Timestamp: time.Now().Unix()}
	d2 := types.Deal{
		BlockHeight: 2, OrderID: "order1", Product: "abc_bcd", Price: fp, Quantity: 200,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: time.Now().Unix()}

	db.AutoMigrate(&types.Deal{})
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	r1 := tx.Create(&d1).Error
	r2 := tx.Create(&d2).Error
	fmt.Printf("%+v, %+v", r1, r2)
	tx.Commit()

	var queryDeal types.Deal
	db.First(&queryDeal).Limit(1)
	fmt.Printf("%+v", queryDeal)

	var allDeals []types.Deal
	db.Find(&allDeals)
	fmt.Printf("%+v", allDeals)

	_, tsEnd := getTimestampRange()
	sql := fmt.Sprintf("select side, product, sum(Quantity) as quantity, max(Price) as high, min(Price) as low from deals "+
		"where Timestamp >= %d and Timestamp < %d group by side", 0, tsEnd)

	rows, err := db.Raw(sql).Rows()
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	for rows.Next() {

		var side, product string
		var quantity float64
		var high float64
		var low float64

		err = rows.Scan(&side, &product, &quantity, &high, &low)
		require.Nil(t, err)
		fmt.Printf("product: %s, quantity: %f, high: %f, low: %f \n", product, quantity, high, low)
	}

	db.Delete(&types.Deal{})

}

func getTimestampRange() (int64, int64) {
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
	endTime := startTime.Add(time.Minute)
	return startTime.Unix(), endTime.Unix()
}

func TestTimestamp(t *testing.T) {
	now := time.Now()
	unixTimestamp := now.Unix()
	fmt.Println(unixTimestamp)

	str := time.Unix(unixTimestamp, 0).Format("2006-01-02 15:04:05")
	fmt.Println(str)

	gTime := time.Unix(unixTimestamp, 0)
	ts1 := time.Date(gTime.Year(), gTime.Month(), gTime.Day(), gTime.Hour(), gTime.Minute(), 0, 0, time.UTC)
	ts2 := ts1.Add(time.Minute)
	strNewTime1 := ts1.Format("2006-01-02 15:04:05")
	strNewTime2 := ts2.Format("2006-01-02 15:04:05")
	fmt.Println(strNewTime1)
	fmt.Println(strNewTime2)

	return
}

func TestSqlite3_AllInOne(t *testing.T) {
	orm, _ := NewSqlite3ORM(false, "/tmp/", "test.db", nil)
	testORMAllInOne(t, orm)
}

func testORMAllInOne(t *testing.T, orm *ORM) {

	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("%+v\n", e)
			debug.PrintStack()
		}
	}()

	errDeleteDeal := orm.deleteDealBefore(time.Now().Unix() + 1)
	errDeleteMatch := orm.deleteMatchResultBefore(time.Now().Unix() + 1)
	errDeleteKlineM1 := orm.deleteKlineM1Before(time.Now().Unix() + 1)
	errDeleteKlineM3 := orm.DeleteKlineBefore(time.Now().Unix()+1, &types.KlineM3{})
	errDeleteKineM15 := orm.DeleteKlineBefore(time.Now().Unix()+1, &types.KlineM15{})
	err := types.NewErrorsMerged(errDeleteDeal, errDeleteMatch, errDeleteKlineM1, errDeleteKlineM3, errDeleteKineM15)
	require.Nil(t, err)

	p := sdk.NewDecWithPrec(1, 2)
	fp, _ := strconv.ParseFloat(p.String(), 64)
	highPrice, _ := strconv.ParseFloat("100", 64)
	lowPrice, _ := strconv.ParseFloat("0.0001", 64)

	product := "abc_bcd"
	adr1 := "asdlfkjsd"

	ts := time.Now().Unix()
	d1 := types.Deal{
		BlockHeight: 1, OrderID: "order0", Product: product, Price: fp, Quantity: 100,
		Sender: adr1, Side: types.BuyOrder, Timestamp: ts - 60*30}
	d2 := types.Deal{
		BlockHeight: 2, OrderID: "order1", Product: product, Price: fp + 0.1, Quantity: 200,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: ts - 60*15}
	d3 := types.Deal{
		BlockHeight: 3, OrderID: "order1", Product: product, Price: fp, Quantity: 300,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: ts - 60*5}
	d4 := types.Deal{
		BlockHeight: 4, OrderID: "order1", Product: product, Price: fp + 0.2, Quantity: 400,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: ts - 60*3 - 1}

	matches := []*types.MatchResult{
		{BlockHeight: 3, Product: product, Price: fp, Quantity: 300, Timestamp: ts - 60*5},
		{BlockHeight: 4, Product: product, Price: highPrice, Quantity: 200, Timestamp: ts - 60},
		{BlockHeight: 5, Product: product, Price: lowPrice, Quantity: 200, Timestamp: ts - 60},
	}
	addCnt, err := orm.AddMatchResults(matches)
	assert.Equal(t, len(matches), addCnt)
	require.Nil(t, err)

	allDeals := []*types.Deal{&d1, &d2, &d3, &d4}
	addCnt, err = orm.AddDeals(allDeals)
	assert.True(t, addCnt == len(allDeals) && err == nil)

	deals, _ := orm.GetDeals(adr1, "", "", 0, 0, 0, 100)
	assert.True(t, len(deals) == len(allDeals) && deals != nil)

	deals, err = orm.getLatestDeals(product, 100)
	require.Nil(t, err)
	assert.True(t, len(deals) == len(allDeals) && deals != nil)
	var allDealVolume, allKM1Volume, allKM3Volume float64
	for _, d := range deals {
		allDealVolume += d.Quantity
	}

	deals, err = orm.getDealsByTimestampRange(product, 0, time.Now().Unix())
	assert.True(t, err == nil && len(deals) == len(allDeals) && deals != nil)

	openDeal, closeDeal := orm.getOpenCloseDeals(0, time.Now().Unix()+1, product)
	assert.True(t, openDeal != nil && closeDeal != nil)

	minDealTS := orm.getDealsMinTimestamp()
	assert.True(t, minDealTS == (ts-60*30))

	ds := DealDataSource{orm: orm}
	endTS := time.Now().Unix()
	if endTS%60 == 0 {
		endTS += 1
	}
	anchorEndTS, cnt, newKlinesM1, err := orm.CreateKline1M(0, endTS, &ds)
	assert.True(t, err == nil)
	assert.True(t, len(newKlinesM1) == cnt)

	products, _ := orm.getAllUpdatedProducts(0, time.Now().Unix())
	assert.True(t, len(products) > 0)

	_, cnt, newKlinesM1, err = orm.CreateKline1M(anchorEndTS, time.Now().Unix()+1, &ds)
	assert.True(t, err == nil)

	maxTS := orm.getKlineMaxTimestamp(&types.KlineM1{})
	assert.True(t, maxTS < ts)

	r, e := orm.getLatestKlineM1ByProduct(product, 100)
	assert.True(t, r != nil && e == nil)
	for _, v := range *r {
		allKM1Volume += v.Volume
	}

	klineM3, e := types.NewKlineFactory("kline_m3", nil)
	assert.True(t, klineM3 != nil && e == nil)
	klineM15, e := types.NewKlineFactory("kline_m15", nil)
	assert.True(t, klineM15 != nil && e == nil)
	anchorEndTS, _, newKlines, err := orm.MergeKlineM1(0, time.Now().Unix()+1, klineM3.(types.IKline))
	require.Nil(t, err)
	require.True(t, len(newKlines) > 0)
	_, _, newKlines, err = orm.MergeKlineM1(0, time.Now().Unix()+1, klineM15.(types.IKline))
	require.Nil(t, err)
	klineM15List := []types.KlineM15{}
	err = orm.GetLatestKlinesByProduct(product, 100, -1, &klineM15List)
	require.Nil(t, err)

	tickers, err := orm.RefreshTickers(0, time.Now().Unix()+1, nil)
	assert.True(t, err == nil && len(tickers) > 0)
	for _, t := range tickers {
		fmt.Println((*t).PrettyString())
	}

	_, _, newKlines, err = orm.MergeKlineM1(anchorEndTS, time.Now().Unix()+1, klineM3.(types.IKline))
	require.Nil(t, err)
	klineM3List := []types.KlineM3{}
	err = orm.GetLatestKlinesByProduct(product, 100, -1, &klineM3List)
	require.Nil(t, err)
	assert.True(t, len(klineM3List) > 0)

	for _, v := range klineM3List {
		//fmt.Printf("%d, %+v\n", v.GetTimestamp(), v.PrettyTimeString())
		allKM3Volume += v.Volume
	}
	err = orm.GetLatestKlinesByProduct(product, 100, -1, &klineM3List)
	require.Nil(t, err)
	assert.True(t, len(klineM3List) > 0)

	assert.True(t, int64(allDealVolume) == int64(allKM1Volume) && int64(allKM3Volume) == int64(allKM1Volume))

	TestORM_KlineM1ToTicker(t)
}

func TestORM_MergeKlineM1(t *testing.T) {

	orm, err := NewSqlite3ORM(false, "/tmp/", "test.db", nil)
	require.Nil(t, err)
	product := "abc_bcd"

	_, err = orm.getLatestKlineM1ByProduct(product, 100)
	require.Nil(t, err)

	klineM3, e := types.NewKlineFactory("kline_m3", nil)
	assert.True(t, klineM3 != nil && e == nil)

	_, _, _, err = orm.MergeKlineM1(0, time.Now().Unix()+1, klineM3.(types.IKline))
	require.Nil(t, err)

	klineM3List := []types.KlineM3{}
	err = orm.GetLatestKlinesByProduct(product, 100, -1, &klineM3List)
	require.Nil(t, err)
	assert.True(t, len(klineM3List) > 0)

}

func TestORM_KlineM1ToTicker(t *testing.T) {
	orm, _ := NewSqlite3ORM(false, "/tmp/", "test.db", nil)
	tickers1, _ := orm.RefreshTickers(0, time.Now().Unix(), nil)
	assert.True(t, len(tickers1) > 0)

	for _, t := range tickers1 {
		fmt.Printf("%s\n", t.PrettyString())
	}

	orm2, _ := NewSqlite3ORM(false, "/tmp/", "test_nil.db", nil)
	tickers2, _ := orm2.RefreshTickers(0, time.Now().Unix(), nil)
	assert.False(t, len(tickers2) > 0)
}

func TestMap(t *testing.T) {
	m := map[string][]int{}
	m["b"] = []int{100}
	r := m["a"]

	for k, v := range m {
		fmt.Println(fmt.Sprintf("k: %s, v: %+v", k, v))
	}

	for k := range m {
		fmt.Println(fmt.Sprintf("k: %s", k))
	}
	assert.True(t, r == nil)
}

func TestTime(t *testing.T) {

	now := time.Now()
	newNow := time.Unix(now.Unix(), 0).UTC()
	fmt.Printf("now: %+v\n", now.Location())
	fmt.Printf("nowNow: %+v\n", newNow.Location())

	startTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	fmt.Printf("startTime: %+v\n", startTime.Location())
	newStartTime := time.Date(newNow.Year(), newNow.Month(), newNow.Day(), newNow.Hour(), newNow.Minute(), 0, 0, time.UTC)
	endTime := newStartTime.Add(time.Minute)
	fmt.Println(endTime.Location())

	fmt.Printf("now: %+v, start: %+v, newStartTime: %+v, end: %+v, startTS: %d, endTS: %d, diff: %d\n",
		now, startTime, newStartTime, endTime, startTime.Unix(), endTime.Unix(), endTime.Unix()-newStartTime.Unix())

	ts := now.Unix()
	m := (ts / 60) * 60
	fmt.Printf("old: %d, anchor: %d\n", ts, m)
}

func TestORM_Get(t *testing.T) {

	orm, _ := NewSqlite3ORM(false, "/tmp/", "test.db", nil)
	r, e := orm.getLatestKlineM1ByProduct("abc_bcd", 100)
	assert.True(t, r != nil && e == nil)

	fmt.Printf("%+v\n", r)
	fmt.Printf("%+v\n", *r)

	for i, v := range *r {
		fmt.Printf("%+v, %+v\n", i, types.TimeString(v.Timestamp))
	}
}

func TestCandles_NewKlinesFactory(t *testing.T) {

	dbDir, err := os.Getwd()
	require.Nil(t, err)
	orm, _ := NewSqlite3ORM(false, dbDir, "backend.db", nil)
	klines, e := types.NewKlinesFactory("kline_m1")
	assert.True(t, klines != nil && e == nil)

	product := types.TestTokenPair
	err = orm.GetLatestKlinesByProduct(product, 100, 0, klines)
	assert.True(t, err == nil)

	iklines := types.ToIKlinesArray(klines, time.Now().Unix(), true)
	assert.True(t, len(iklines) > 0)
	//for _, k := range iklines {
	//	fmt.Printf("%+v\n", k.PrettyTimeString())
	//}

	result := types.ToRestfulData(&iklines, 100)
	for _, r := range result {
		fmt.Printf("%+v\n", r)
	}

	r, _ := json.Marshal(result)

	err = tcommon.WriteFile("/tmp/k1.txt", r, os.ModePerm)
	require.Nil(t, err)
}

func constructLocalBackendDB(orm *ORM) (err error) {
	m := types.GetAllKlineMap()
	crrTs := time.Now().Unix()
	ds := DealDataSource{orm: orm}
	if _, _, _, err := orm.CreateKline1M(0, crrTs, &ds); err != nil {
		return err
	}

	for freq, tname := range m {
		if freq == 60 {
			continue
		}
		kline, _ := types.NewKlineFactory(tname, nil)
		if _, _, _, err = orm.MergeKlineM1(0, crrTs, kline.(types.IKline)); err != nil {
			return err
		}
	}
	return nil
}

func TestCandles_FromLocalDB(t *testing.T) {
	dbDir, err := os.Getwd()
	require.Nil(t, err)
	orm, err := NewSqlite3ORM(false, dbDir, "backend.db", nil)
	require.Nil(t, err)
	product := types.TestTokenPair
	limit := 10

	maxKlines, err := types.NewKlinesFactory("kline_m1440")
	require.Nil(t, err)
	err = orm.GetLatestKlinesByProduct(product, limit, time.Now().Unix(), maxKlines)
	require.Nil(t, err)
	maxIklines := types.ToIKlinesArray(maxKlines, time.Now().Unix(), true)
	if len(maxIklines) == 0 {
		err := constructLocalBackendDB(orm)
		require.Nil(t, err)
	}

	m := types.GetAllKlineMap()
	for freq, tname := range m {
		if freq > 1440*60 {
			continue
		}

		klines, _ := types.NewKlinesFactory(tname)
		e := orm.GetLatestKlinesByProduct(product, limit, time.Now().Unix(), klines)
		assert.True(t, e == nil)

		iklines := types.ToIKlinesArray(klines, time.Now().Unix(), true)
		assert.True(t, len(iklines) > 0)
		//for _, k := range iklines {
		//	fmt.Printf("%+v\n", k.PrettyTimeString())
		//}

		restDatas := types.ToRestfulData(&iklines, limit)
		assert.True(t, len(restDatas) <= limit)
	}

	maxTS := orm.getDealsMaxTimestamp()
	assert.True(t, maxTS > 0)
}

// Deals
func testORMDeals(t *testing.T, orm *ORM) {

	addDeals := []*types.Deal{
		{Timestamp: 100, BlockHeight: 1, OrderID: "ID1", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
		{Timestamp: 300, BlockHeight: 3, OrderID: "ID2", Sender: "addr1", Product: "btc_" + common.NativeToken, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
		{Timestamp: 200, BlockHeight: 2, OrderID: "ID3", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
		{Timestamp: 400, BlockHeight: 1, OrderID: "ID4", Sender: "addr2", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
	}
	// Test AddDeals
	cnt, err := orm.AddDeals(addDeals)
	require.Nil(t, err)
	require.EqualValues(t, 4, cnt)
	// Test GetDeals
	// filtered by address, sorted by timestamp desc, and paged by offset and limit
	deals, total := orm.GetDeals("addr1", "", "", 0, 0, 1, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 2, len(deals))
	require.EqualValues(t, "ID3", deals[0].OrderID)
	require.EqualValues(t, "ID1", deals[1].OrderID)

	// filtered by address & product & side
	deals, total = orm.GetDeals("addr1", "btc_"+common.NativeToken, types.BuyOrder, 0, 0, 0, 10)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(deals))
	require.EqualValues(t, "ID2", deals[0].OrderID)

	// filtered by address & start end time
	deals, total = orm.GetDeals("addr1", "", "", 200, 300, 0, 10)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(deals))
	require.EqualValues(t, "ID3", deals[0].OrderID)

	// too large offset
	deals, total = orm.GetDeals("addr1", "", "", 0, 0, 3, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(deals))

	// GetDealsV2
	dealsV2 := orm.GetDealsV2("addr1", types.TestTokenPair, types.BuyOrder, "100", "300", 1)
	require.EqualValues(t, 1, len(dealsV2))
	require.EqualValues(t, addDeals[2], &dealsV2[0])

	mrds := MergeResultDataSource{orm}
	oPrice, cPrice := mrds.getOpenClosePrice(0, time.Now().Unix(), types.TestTokenPair)
	require.EqualValues(t, 10, oPrice)
	require.EqualValues(t, 10, cPrice)
}

// Matches
func TestORMMatches(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)

	addMatches := []*types.MatchResult{
		{Timestamp: 100, BlockHeight: 1, Product: types.TestTokenPair, Price: 10.0, Quantity: 1.0},
		{Timestamp: 100, BlockHeight: 1, Product: "btc_" + common.NativeToken, Price: 11.0, Quantity: 2.0},
		{Timestamp: 200, BlockHeight: 2, Product: types.TestTokenPair, Price: 12.0, Quantity: 3.0},
		{Timestamp: 300, BlockHeight: 3, Product: types.TestTokenPair, Price: 13.0, Quantity: 4.0},
	}
	// Test AddMatchResults
	cnt, err := orm.AddMatchResults(addMatches)
	require.Nil(t, err)
	require.EqualValues(t, 4, cnt)

	// Test GetMatchResults
	// filtered by product, sorted by timestamp desc, and paged by offset and limit
	matches, total := orm.GetMatchResults(types.TestTokenPair, 0, 0, 1, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 2, len(matches))
	require.EqualValues(t, 3, matches[0].Quantity)
	require.EqualValues(t, 1, matches[1].Quantity)

	// filtered by address & start end time
	matches, total = orm.GetMatchResults("", 100, 200, 0, 3)
	require.EqualValues(t, 2, total)
	require.EqualValues(t, 2, len(matches))

	// too large offset
	matches, total = orm.GetMatchResults(types.TestTokenPair, 0, 0, 3, 3)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(matches))

	// GetMatchResultsV2
	matchesV2 := orm.GetMatchResultsV2(types.TestTokenPair, "100", "500", 2)
	require.EqualValues(t, 2, len(matchesV2))
	require.EqualValues(t, addMatches[3], &matchesV2[0])
	require.EqualValues(t, addMatches[2], &matchesV2[1])

	//
	matches, err = orm.getLatestMatchResults(types.TestTokenPair, 1)
	require.Nil(t, err)
	require.EqualValues(t, 1, len(matches))
	require.EqualValues(t, 3, matches[0].BlockHeight)

	//
	stamp := orm.getMergeResultMaxTimestamp()
	require.EqualValues(t, 300, stamp)

	//
	mrds := MergeResultDataSource{orm}
	require.EqualValues(t, 100, mrds.getDataSourceMinTimestamp())
	sql := `select product, sum(Quantity) as quantity, max(Price) as high, min(Price) as low, count(price) as cnt from match_results where Timestamp >= 0 and Timestamp < 1574406957 group by product`
	require.EqualValues(t, sql, mrds.getMaxMinSumByGroupSQL(0, 1574406957))

}

func TestSqlite3_ORMDeals(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMDeals(t, orm)
}

// FeeDetail
func testORMFeeDetails(t *testing.T, orm *ORM) {

	feeDetails := []*token.FeeDetail{
		{Address: "addr1", Fee: "0.1" + common.NativeToken, FeeType: types.FeeTypeOrderCancel, Timestamp: 100},
		{Address: "addr1", Fee: "0.5" + common.NativeToken, FeeType: types.FeeTypeOrderNew, Timestamp: 300},
		{Address: "addr1", Fee: "0.2" + common.NativeToken, FeeType: types.FeeTypeOrderDeal, Timestamp: 200},
		{Address: "addr2", Fee: "0.3" + common.NativeToken, FeeType: types.FeeTypeOrderDeal, Timestamp: 100},
	}

	// Test AddFeeDetails
	cnt, err := orm.AddFeeDetails(feeDetails)
	require.EqualValues(t, 4, cnt)
	require.Nil(t, err)

	// Test GetFeeDetails
	fees, total := orm.GetFeeDetails("addr1", 1, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 2, len(fees))
	require.EqualValues(t, 200, fees[0].Timestamp)
	require.EqualValues(t, 100, fees[1].Timestamp)
	// too large offset
	fees, total = orm.GetFeeDetails("addr1", 3, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(fees))

	// GetFeeDetailsV2
	feesV2 := orm.GetFeeDetailsV2("addr1", "100", "300", 1)
	require.EqualValues(t, 1, len(feesV2))
	require.EqualValues(t, feeDetails[2], &feesV2[0])

}

func TestSqlite3_FeeDetails(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMFeeDetails(t, orm)
}

// Order
func testORMOrders(t *testing.T, orm *ORM) {

	orders := []*types.Order{
		{TxHash: "hash1", OrderID: "ID1", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 100},
		{TxHash: "hash2", OrderID: "ID2", Sender: "addr1", Product: "btc_" + common.NativeToken, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 300},
		{TxHash: "hash3", OrderID: "ID3", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 200},
		{TxHash: "hash4", OrderID: "ID4", Sender: "addr2", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 150},
	}
	// Test AddOrders
	cnt, err := orm.AddOrders(orders)
	require.EqualValues(t, 4, cnt)
	require.Nil(t, err)

	// Test GetOrderList
	// filtered by address, sorted by timestamp desc, and paged by offset and limit
	getOrders, total := orm.GetOrderList("addr1", "", "", true, 1, 2, 0, 0, false)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 2, len(getOrders))
	require.EqualValues(t, "ID3", getOrders[0].OrderID)
	require.EqualValues(t, "ID1", getOrders[1].OrderID)

	// filtered by product & side
	getOrders, total = orm.GetOrderList("addr1", "btc_"+common.NativeToken, types.BuyOrder, true, 0, 10, 0, 0, false)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(getOrders))
	require.EqualValues(t, "ID2", getOrders[0].OrderID)

	//// GetOrderListV2 : open order
	openOrdersV2 := orm.GetOrderListV2(types.TestTokenPair, "addr1", types.BuyOrder, true, "10", "300", 1)
	require.Equal(t, 1, len(openOrdersV2))
	require.Equal(t, orders[2], &openOrdersV2[0])

	// TestUpdateOrders
	updateOrders := []*types.Order{
		{TxHash: "hash1", OrderID: "ID1", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 3, FilledAvgPrice: "0", RemainQuantity: "0", Timestamp: 100},
		{TxHash: "hash2", OrderID: "ID2", Sender: "addr1", Product: "btc_" + common.NativeToken, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 2, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 300},
		{TxHash: "hash3", OrderID: "ID3", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 4, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 200},
	}
	cnt, err = orm.UpdateOrders(updateOrders)
	require.Nil(t, err)
	require.EqualValues(t, 3, cnt)

	// filtered closed orders
	getOrders, total = orm.GetOrderList("addr1", "", "", false, 0, 10, 0, 0, false)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 3, len(getOrders))
	require.EqualValues(t, "ID2", getOrders[0].OrderID)
	require.EqualValues(t, "ID3", getOrders[1].OrderID)
	require.EqualValues(t, "ID1", getOrders[2].OrderID)

	// hide no fill orders
	getOrders, total = orm.GetOrderList("addr1", "", "", false, 0, 10, 0, 0, true)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(getOrders))
	require.EqualValues(t, "ID3", getOrders[0].OrderID)

	// too large offset
	getOrders, total = orm.GetOrderList("addr1", "", "", false, 3, 10, 0, 0, false)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(getOrders))

	// GetOrderListV2 : other order: close,filled,cancel,……
	otherOrdersV2 := orm.GetOrderListV2(types.TestTokenPair, "addr1", types.BuyOrder, false, "10", "300", 1)
	require.Equal(t, 1, len(otherOrdersV2))
	require.Equal(t, updateOrders[2], &otherOrdersV2[0])

	// v2 GetOrderByID
	ordersByExistID := orm.GetOrderByID("ID1")
	require.EqualValues(t, updateOrders[0], ordersByExistID)
	ordersByNotExistID := orm.GetOrderByID("not_exist_ID")
	require.Nil(t, ordersByNotExistID)

}

func TestSqlite3_Orders(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMOrders(t, orm)
}

// Transactions
func testORMTransactions(t *testing.T, orm *ORM) {

	txs := []*types.Transaction{
		{TxHash: "hash1", Type: types.TxTypeTransfer, Address: "addr1", Symbol: common.TestToken, Side: types.TxSideFrom, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 100},
		{TxHash: "hash2", Type: types.TxTypeOrderNew, Address: "addr1", Symbol: types.TestTokenPair, Side: types.TxSideBuy, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 300},
		{TxHash: "hash3", Type: types.TxTypeOrderCancel, Address: "addr1", Symbol: types.TestTokenPair, Side: types.TxSideSell, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 200},
		{TxHash: "hash4", Type: types.TxTypeTransfer, Address: "addr2", Symbol: common.TestToken, Side: types.TxSideTo, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 100},
	}
	// Test AddTransactions
	cnt, err := orm.AddTransactions(txs)
	require.Nil(t, err)
	require.EqualValues(t, 4, cnt)

	// Test GetTransactionList
	// filtered by address, sorted by timestamp desc, and paged by offset and limit
	getTxs, total := orm.GetTransactionList("addr1", 0, 0, 0, 1, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 2, len(getTxs))
	require.EqualValues(t, "hash3", getTxs[0].TxHash)
	require.EqualValues(t, "hash1", getTxs[1].TxHash)

	// filtered by address & txType
	getTxs, total = orm.GetTransactionList("addr1", types.TxTypeOrderNew, 0, 0, 0, 10)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(getTxs))
	require.EqualValues(t, "hash2", getTxs[0].TxHash)

	// filtered by address & start end time
	getTxs, total = orm.GetTransactionList("addr1", 0, 200, 300, 0, 10)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(getTxs))
	require.EqualValues(t, "hash3", getTxs[0].TxHash)

	// too large offset
	getTxs, total = orm.GetTransactionList("addr1", 0, 0, 0, 3, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(getTxs))

	// GetTransactionListV2
	getTxsV2 := orm.GetTransactionListV2("addr1", types.TxTypeOrderNew, "10", "400", 1)
	require.EqualValues(t, 1, len(getTxsV2))
	require.EqualValues(t, txs[1], &getTxsV2[0])
}

func TestSqlite3_Transactions(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMTransactions(t, orm)
}

func Test_Time(t *testing.T) {
	now := time.Now()
	time.Sleep(time.Second)

	r := time.Since(now).Nanoseconds() / 1000000
	fmt.Println(r)
}

func TestORM_deleteKlinesAfter(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)

	k, _ := types.NewKlineFactory("kline_m1", nil)
	err := orm.deleteKlinesAfter(0, types.TestTokenPair, k)
	assert.True(t, err == nil)
}

func TestORM_GetDealsMaxTimestamp(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)

	noTS := orm.getDealsMaxTimestamp()
	assert.True(t, noTS == -1)

}
func testORMBatchInsert(t *testing.T, orm *ORM) {
	newOrders := []*types.Order{}

	for i := 0; i < 2000; i++ {
		oid := fmt.Sprintf("FAKEID-%04d", i)
		o := types.Order{TxHash: "hash1", OrderID: oid, Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "1.3", RemainQuantity: "1.5", Timestamp: 100}
		newOrders = append(newOrders, &o)
	}

	updatedOrders := []*types.Order{
		{TxHash: "hash2", OrderID: "FAKEID-0002", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "1.4", RemainQuantity: "1.7", Timestamp: 100},
	}

	txs := []*types.Transaction{
		{TxHash: "FAKEIDHash-1", Type: types.TxTypeTransfer, Address: "addr1", Symbol: common.TestToken, Side: types.TxSideFrom, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 100},
	}

	addDeals := []*types.Deal{
		{Timestamp: 100, BlockHeight: 1, OrderID: "FAKEID-0001", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
		{Timestamp: 300, BlockHeight: 3, OrderID: "FAKEID-0002", Sender: "addr1", Product: "btc_" + common.NativeToken, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
		{Timestamp: 200, BlockHeight: 2, OrderID: "FAKEID-0003", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
		{Timestamp: 400, BlockHeight: 1, OrderID: "FAKEID-0004", Sender: "addr2", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
	}

	mrs := []*types.MatchResult{
		{Timestamp: 100, BlockHeight: 1, Product: types.TestTokenPair, Price: 10.0, Quantity: 1.0},
	}

	feeDetails := []*token.FeeDetail{
		{Address: "addr1", Fee: "0.1" + common.NativeToken, FeeType: types.FeeTypeOrderCancel, Timestamp: 100},
		{Address: "addr1", Fee: "0.5" + common.NativeToken, FeeType: types.FeeTypeOrderNew, Timestamp: 300},
		{Address: "addr1", Fee: "0.2" + common.NativeToken, FeeType: types.FeeTypeOrderDeal, Timestamp: 200},
		{Address: "addr2", Fee: "0.3" + common.NativeToken, FeeType: types.FeeTypeOrderDeal, Timestamp: 100},
	}

	swapInfos := []*types.SwapInfo{
		{Address: "addr1", TokenPairName: types.TestTokenPair, BaseTokenAmount: "10000xxb", QuoteTokenAmount: "10000yyb",
			SellAmount: "10xxb", BuysAmount: "9.8yyb", Price: "1", Timestamp: 100},
	}

	resultMap, e := orm.BatchInsertOrUpdate(newOrders, updatedOrders, addDeals, mrs, feeDetails, txs, swapInfos)
	require.True(t, resultMap != nil && e == nil)

	require.True(t, resultMap != nil && resultMap["newOrders"] == 2000)
	require.True(t, resultMap != nil && resultMap["updatedOrders"] == 1)
	require.True(t, resultMap != nil && resultMap["transactions"] == 1)
	require.True(t, resultMap != nil && resultMap["deals"] == 4)
	require.True(t, resultMap != nil && resultMap["feeDetails"] == 4)
	require.True(t, resultMap != nil && resultMap["swapInfos"] == 1)

	resultMap2, e2 := orm.BatchInsertOrUpdate(newOrders, updatedOrders, addDeals, mrs, feeDetails, txs, swapInfos)
	fmt.Printf("%+v\n", e2)
	require.True(t, resultMap2 != nil, resultMap2)
	require.True(t, e2 != nil, e2)
}

func TestORM_BatchInsert(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMBatchInsert(t, orm)
}

func TestORM_CloseDB(t *testing.T) {
	closeORM, err := NewSqlite3ORM(false, "/tmp/", "test_close.db", nil)
	require.Nil(t, err)
	defer DeleteDB("/tmp/test_close.db")
	err = closeORM.Close()
	require.Nil(t, err)

	// query after close DB
	products, err := closeORM.getAllUpdatedProducts(0, -1)
	require.Error(t, err)
	require.Nil(t, products)
	klines, err := closeORM.getLatestKlineM1ByProduct("abc_bcd", 100)
	require.Nil(t, klines)
	require.Error(t, err)
	matches, err := closeORM.getLatestMatchResults(types.TestTokenPair, 1)
	require.Equal(t, 0, len(matches))
	require.Error(t, err)
	matchResults, err := closeORM.getMatchResultsByTimeRange("", 100, 500)
	require.Equal(t, 0, len(matchResults))
	require.Error(t, err)
	deals, err := closeORM.getLatestDeals("", 100)
	require.Equal(t, 0, len(deals))
	require.Error(t, err)
	deals, err = closeORM.getDealsByTimestampRange("", 0, time.Now().Unix())
	require.Nil(t, deals)
	require.Error(t, err)

	// delete after close DB
	err = closeORM.deleteKlinesBefore(1, &types.KlineM15{})
	require.Error(t, err)
	err = closeORM.deleteKlinesAfter(1, "", &types.KlineM15{})
	require.Error(t, err)
	err = closeORM.deleteDealBefore(time.Now().Unix() + 1)
	require.Error(t, err)
	err = closeORM.deleteMatchResultBefore(1000)
	require.Error(t, err)

	// insert after close DB
	cnt, err := closeORM.AddMatchResults([]*types.MatchResult{
		{Timestamp: 100, BlockHeight: 1, Product: types.TestTokenPair, Price: 10.0, Quantity: 1.0},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddDeals([]*types.Deal{
		{Timestamp: 100, BlockHeight: 1, OrderID: "FAKEID-0001", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: 10.0, Quantity: 1.0, Fee: "0"},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddFeeDetails([]*token.FeeDetail{
		{Address: "addr1", Fee: "0.1" + common.NativeToken, FeeType: types.FeeTypeOrderCancel, Timestamp: 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddOrders([]*types.Order{
		{TxHash: "hash1", OrderID: "ID1", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddTransactions([]*types.Transaction{
		{TxHash: "hash1", Type: types.TxTypeTransfer, Address: "addr1", Symbol: common.TestToken, Side: types.TxSideFrom, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.UpdateOrders([]*types.Order{
		{TxHash: "hash1", OrderID: "ID1", Sender: "addr1", Product: types.TestTokenPair, Side: types.BuyOrder, Price: "10.0", Quantity: "1.1", Status: 0, FilledAvgPrice: "0", RemainQuantity: "1.1", Timestamp: 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)
}

func TestORM_DeferRollbackTx(t *testing.T) {

	orm, err := NewSqlite3ORM(false, "/tmp/", "test_failed.db", nil)
	require.Nil(t, err)
	defer DeleteDB("/tmp/test_failed.db")
	defer orm.deferRollbackTx(orm.db, fmt.Errorf("failed to commit"))
	orm.db.AutoMigrate("rollback")
	panic("orm deferRollbackTx recover will catch the panic")

}
