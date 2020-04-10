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
	"github.com/okex/okchain/x/backend/cases"
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token"
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
	db.LogMode(true)
	p := sdk.NewDecWithPrec(1, 2)
	p.String()

	fp, _ := strconv.ParseFloat(p.String(), 64)

	d1 := types.Deal{
		BlockHeight: 1, OrderId: "order0", Product: "abc_bcd", Price: fp, Quantity: 100,
		Sender: "asdlfkjsd", Side: types.SellOrder, Timestamp: time.Now().Unix()}
	d2 := types.Deal{
		BlockHeight: 2, OrderId: "order1", Product: "abc_bcd", Price: fp, Quantity: 200,
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

		rows.Scan(&side, &product, &quantity, &high, &low)
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

//func TestMysql_AllInOne(t *testing.T) {
//	orm, _ := NewMysqlORM()
//	testORMAllInOne(t, orm)
//}

func testORMAllInOne(t *testing.T, orm *ORM) {

	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("%+v\n", e)
			debug.PrintStack()
		}
	}()

	err := orm.DeleteDealBefore(time.Now().Unix() + 1)
	orm.DeleteMatchResultBefore(time.Now().Unix() + 1)
	orm.DeleteKlineM1Before(time.Now().Unix() + 1)
	orm.DeleteKlineBefore(time.Now().Unix()+1, &types.KlineM3{})
	orm.DeleteKlineBefore(time.Now().Unix()+1, &types.KlineM15{})
	assert.True(t, err == nil)

	p := sdk.NewDecWithPrec(1, 2)
	fp, _ := strconv.ParseFloat(p.String(), 64)
	highPrice, _ := strconv.ParseFloat("100", 64)
	lowPrice, _ := strconv.ParseFloat("0.0001", 64)

	product := "abc_bcd"
	adr1 := "asdlfkjsd"

	ts := time.Now().Unix()
	d1 := types.Deal{
		BlockHeight: 1, OrderId: "order0", Product: product, Price: fp, Quantity: 100,
		Sender: adr1, Side: types.BuyOrder, Timestamp: ts - 60*30}
	d2 := types.Deal{
		BlockHeight: 2, OrderId: "order1", Product: product, Price: fp + 0.1, Quantity: 200,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: ts - 60*15}
	d3 := types.Deal{
		BlockHeight: 3, OrderId: "order1", Product: product, Price: fp, Quantity: 300,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: ts - 60*5}
	d4 := types.Deal{
		BlockHeight: 4, OrderId: "order1", Product: product, Price: fp + 0.2, Quantity: 400,
		Sender: "asdlfkjsd", Side: types.BuyOrder, Timestamp: ts - 60*3 - 1}

	matches := []*types.MatchResult{
		{BlockHeight: 3, Product: product, Price: fp, Quantity: 300, Timestamp: ts - 60*5},
		{BlockHeight: 4, Product: product, Price: highPrice, Quantity: 200, Timestamp: ts - 60},
		{BlockHeight: 5, Product: product, Price: lowPrice, Quantity: 200, Timestamp: ts - 60},
	}
	addCnt, err := orm.AddMatchResults(matches)
	assert.Equal(t, len(matches), addCnt)
	require.Nil(t, err)

	all_deals := []*types.Deal{&d1, &d2, &d3, &d4}
	addCnt, err = orm.AddDeals(all_deals)
	assert.True(t, addCnt == len(all_deals) && err == nil)

	deals, _ := orm.GetDeals(adr1, "", "", 0, 0, 0, 100)
	assert.True(t, len(deals) == len(all_deals) && deals != nil)

	deals, err = orm.GetLatestDeals(product, 100)
	assert.True(t, len(deals) == len(all_deals) && deals != nil)
	var allDealVolume, allKM1Volume, allKM3Volume float64
	for _, d := range deals {
		fmt.Printf("%+v\n", d)
		allDealVolume += d.Quantity
	}

	deals, err = orm.GetDealsByTimestampRange(product, 0, time.Now().Unix())
	assert.True(t, len(deals) == len(all_deals) && deals != nil)

	openDeal, closeDeal := orm.getOpenCloseDeals(0, time.Now().Unix()+1, product)
	assert.True(t, openDeal != nil && closeDeal != nil)

	minDealTS := orm.GetDealsMinTimestamp()
	assert.True(t, minDealTS == (ts-60*30))

	ds := DealDataSource{orm: orm}
	anchorEndTS, cnt, err := orm.CreateKline1min(0, time.Now().Unix()+1, &ds)
	fmt.Printf("CreateKline1min ERROR: %+v", err)
	assert.True(t, err == nil, cnt == 3)

	products, _ := orm.GetAllUpdatedProducts(0, time.Now().Unix())
	assert.True(t, products != nil && len(products) > 0)
	fmt.Printf("%+v \n", products)

	anchorEndTS, cnt, err = orm.CreateKline1min(anchorEndTS, time.Now().Unix()+1, &ds)
	fmt.Printf("CreateKline1min ERROR: %+v", err)
	assert.True(t, err == nil, cnt == 1)

	maxTS := orm.GetKlineMaxTimestamp(&types.KlineM1{})
	assert.True(t, maxTS < ts)

	r, e := orm.GetLatestKlineM1ByProduct(product, 100)
	assert.True(t, r != nil && e == nil)
	fmt.Printf("NOW : %s\n", types.TimeString(ts))
	for _, v := range *r {
		//fmt.Printf("%d, %+v\n", v.GetTimestamp(), v.PrettyTimeString())
		allKM1Volume += v.Volume
	}

	kM3, e := types.NewKlineFactory("kline_m3", nil)
	kM15, e := types.NewKlineFactory("kline_m15", nil)
	assert.True(t, kM3 != nil && e == nil)
	anchorEndTS, cnt, err = orm.MergeKlineM1(0, time.Now().Unix()+1, kM3.(types.IKline))

	orm.MergeKlineM1(0, time.Now().Unix()+1, kM15.(types.IKline))
	klineM15List := []types.KlineM15{}
	orm.GetLatestKlinesByProduct(product, 100, -1, &klineM15List)

	tickers, err := orm.RefreshTickers(0, time.Now().Unix()+1, nil)
	assert.True(t, tickers != nil && len(tickers) > 0)
	for _, t := range tickers {
		fmt.Println((*t).PrettyString())
	}

	anchorEndTS, cnt, err = orm.MergeKlineM1(anchorEndTS, time.Now().Unix()+1, kM3.(types.IKline))

	kM3List := []types.KlineM3{}
	orm.GetLatestKlinesByProduct(product, 100, -1, &kM3List)
	assert.True(t, kM3List != nil && len(kM3List) > 0)

	for _, v := range kM3List {
		//fmt.Printf("%d, %+v\n", v.GetTimestamp(), v.PrettyTimeString())
		allKM3Volume += v.Volume
	}
	orm.GetLatestKlinesByProduct(product, 100, -1, &kM3List)
	assert.True(t, kM3List != nil && len(kM3List) > 0)

	assert.True(t, int64(allDealVolume) == int64(allKM1Volume) && int64(allKM3Volume) == int64(allKM1Volume))

	TestORM_KlineM1ToTicker(t)
}

func TestORM_MergeKlineM1(t *testing.T) {

	orm, _ := NewSqlite3ORM(false, "/tmp/", "test.db", nil)
	product := "abc_bcd"

	orm.GetLatestKlineM1ByProduct(product, 100)
	//for _, v := range *klinesM1 {
	//	fmt.Printf("%d, %+v\n", v.GetTimestamp(), v.PrettyTimeString())
	//}

	kM3, e := types.NewKlineFactory("kline_m3", nil)
	assert.True(t, kM3 != nil && e == nil)

	orm.MergeKlineM1(0, time.Now().Unix()+1, kM3.(types.IKline))

	kM3List := []types.KlineM3{}
	orm.GetLatestKlinesByProduct(product, 100, -1, &kM3List)
	assert.True(t, kM3List != nil && len(kM3List) > 0)

	//for _, v := range kM3List {
	//	fmt.Printf("%d, %+v\n", v.GetTimestamp(), v.PrettyTimeString())
	//}
}

func TestORM_KlineM1ToTicker(t *testing.T) {
	orm, _ := NewSqlite3ORM(false, "/tmp/", "test.db", nil)
	tickers1, _ := orm.RefreshTickers(0, time.Now().Unix(), nil)
	assert.True(t, tickers1 != nil && len(tickers1) > 0)

	for _, t := range tickers1 {
		fmt.Printf("%s\n", t.PrettyString())
	}

	orm2, _ := NewSqlite3ORM(false, "/tmp/", "test_nil.db", nil)
	tickers2, _ := orm2.RefreshTickers(0, time.Now().Unix(), nil)
	assert.False(t, tickers2 != nil && len(tickers2) > 0)
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
	r, e := orm.GetLatestKlineM1ByProduct("abc_bcd", 100)
	assert.True(t, r != nil && e == nil)

	fmt.Printf("%+v\n", r)
	fmt.Printf("%+v\n", *r)

	for i, v := range *r {
		fmt.Printf("%+v, %+v\n", i, types.TimeString(v.Timestamp))
	}
}

func TestCandles_NewKlinesFactory(t *testing.T) {

	dbDir := cases.GetBackendDBDir()
	orm, _ := NewSqlite3ORM(false, dbDir, "backend.db", nil)
	klines, e := types.NewKlinesFactory("kline_m1")
	assert.True(t, klines != nil && e == nil)

	product := types.TestTokenPair
	err := orm.GetLatestKlinesByProduct(product, 100, 0, klines)
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

	tcommon.WriteFile("/tmp/k1.txt", r, os.ModePerm)
}

func constructLocalBackendDB(orm *ORM) {
	m := types.GetAllKlineMap()
	crrTs := time.Now().Unix()
	ds := DealDataSource{orm: orm}
	orm.CreateKline1min(0, crrTs, &ds)
	for freq, tname := range m {
		if freq == 60 {
			continue
		}
		kline, _ := types.NewKlineFactory(tname, nil)
		orm.MergeKlineM1(0, crrTs, kline.(types.IKline))
	}
}

func TestCandles_FromLocalDB(t *testing.T) {
	orm, _ := NewSqlite3ORM(false, cases.GetBackendDBDir(), "backend.db", nil)
	product := types.TestTokenPair
	limit := 10

	maxKlines, _ := types.NewKlinesFactory("kline_m1440")
	orm.GetLatestKlinesByProduct(product, limit, time.Now().Unix(), maxKlines)
	maxIklines := types.ToIKlinesArray(maxKlines, time.Now().Unix(), true)
	if len(maxIklines) == 0 {
		constructLocalBackendDB(orm)
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

	maxTS := orm.GetDealsMaxTimestamp()
	assert.True(t, maxTS > 0)
}

// Deals
func testORMDeals(t *testing.T, orm *ORM) {

	addDeals := []*types.Deal{
		{100, 1, "ID1", "addr1", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
		{300, 3, "ID2", "addr1", "btc_" + common.NativeToken, types.BuyOrder, 10.0, 1.0, "0"},
		{200, 2, "ID3", "addr1", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
		{400, 1, "ID4", "addr2", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
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
	require.EqualValues(t, "ID3", deals[0].OrderId)
	require.EqualValues(t, "ID1", deals[1].OrderId)

	// filtered by address & product & side
	deals, total = orm.GetDeals("addr1", "btc_"+common.NativeToken, types.BuyOrder, 0, 0, 0, 10)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(deals))
	require.EqualValues(t, "ID2", deals[0].OrderId)

	// filtered by address & start end time
	deals, total = orm.GetDeals("addr1", "", "", 200, 300, 0, 10)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(deals))
	require.EqualValues(t, "ID3", deals[0].OrderId)

	// too large offset
	deals, total = orm.GetDeals("addr1", "", "", 0, 0, 3, 2)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(deals))

	// GetDealsV2
	dealsV2 := orm.GetDealsV2("addr1", types.TestTokenPair, types.BuyOrder, "100", "300", 1)
	require.EqualValues(t, 1, len(dealsV2))
	require.EqualValues(t, addDeals[2], &dealsV2[0])

	mrds := MergeResultDataSource{orm}
	oPrice, cPrice := mrds.GetOpenClosePrice(0, time.Now().Unix(), types.TestTokenPair)
	require.EqualValues(t, 10, oPrice)
	require.EqualValues(t, 10, cPrice)
}

// Matches
func TestORMMatches(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)

	addMatches := []*types.MatchResult{
		{100, 1, types.TestTokenPair, 10.0, 1.0},
		{100, 1, "btc_" + common.NativeToken, 11.0, 2.0},
		{200, 2, types.TestTokenPair, 12.0, 3.0},
		{300, 3, types.TestTokenPair, 13.0, 4.0},
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
	matches, err = orm.GetLatestMatchResults(types.TestTokenPair, 1)
	require.EqualValues(t, 1, len(matches))
	require.EqualValues(t, 3, matches[0].BlockHeight)

	//
	stamp := orm.GetMergeResultMaxTimestamp()
	require.EqualValues(t, 300, stamp)

	//
	mrds := MergeResultDataSource{orm}
	require.EqualValues(t, 100, mrds.GetDataSourceMinTimestamp())
	sql := `select product, sum(Quantity) as quantity, max(Price) as high, min(Price) as low, count(price) as cnt from match_results where Timestamp >= 0 and Timestamp < 1574406957 group by product`
	require.EqualValues(t, sql, mrds.GetMaxMinSumByGroupSQL(0, 1574406957))

}

func TestSqlite3_ORMDeals(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMDeals(t, orm)
}

// FeeDetail
func testORMFeeDetails(t *testing.T, orm *ORM) {

	feeDetails := []*token.FeeDetail{
		{"addr1", "0.1" + common.NativeToken, types.FeeTypeOrderCancel, 100},
		{"addr1", "0.5" + common.NativeToken, types.FeeTypeOrderNew, 300},
		{"addr1", "0.2" + common.NativeToken, types.FeeTypeOrderDeal, 200},
		{"addr2", "0.3" + common.NativeToken, types.FeeTypeOrderDeal, 100},
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
		{"hash1", "ID1", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "0", "1.1", 100},
		{"hash2", "ID2", "addr1", "btc_" + common.NativeToken, types.BuyOrder, "10.0", "1.1", 0, "0", "1.1", 300},
		{"hash3", "ID3", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "0", "1.1", 200},
		{"hash4", "ID4", "addr2", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "0", "1.1", 150},
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
	require.EqualValues(t, "ID3", getOrders[0].OrderId)
	require.EqualValues(t, "ID1", getOrders[1].OrderId)

	// filtered by product & side
	getOrders, total = orm.GetOrderList("addr1", "btc_"+common.NativeToken, types.BuyOrder, true, 0, 10, 0, 0, false)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(getOrders))
	require.EqualValues(t, "ID2", getOrders[0].OrderId)

	//// GetOrderListV2 : open order
	openOrdersV2 := orm.GetOrderListV2(types.TestTokenPair, "addr1", types.BuyOrder, true, "10", "300", 1)
	require.Equal(t, 1, len(openOrdersV2))
	require.Equal(t, orders[2], &openOrdersV2[0])

	// TestUpdateOrders
	updateOrders := []*types.Order{
		{"hash1", "ID1", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 3, "0", "0", 100},
		{"hash2", "ID2", "addr1", "btc_" + common.NativeToken, types.BuyOrder, "10.0", "1.1", 2, "0", "1.1", 300},
		{"hash3", "ID3", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 4, "0", "1.1", 200},
	}
	cnt, err = orm.UpdateOrders(updateOrders)
	require.Nil(t, err)
	require.EqualValues(t, 3, cnt)

	// filtered closed orders
	getOrders, total = orm.GetOrderList("addr1", "", "", false, 0, 10, 0, 0, false)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 3, len(getOrders))
	require.EqualValues(t, "ID2", getOrders[0].OrderId)
	require.EqualValues(t, "ID3", getOrders[1].OrderId)
	require.EqualValues(t, "ID1", getOrders[2].OrderId)

	// hide no fill orders
	getOrders, total = orm.GetOrderList("addr1", "", "", false, 0, 10, 0, 0, true)
	require.EqualValues(t, 1, total)
	require.EqualValues(t, 1, len(getOrders))
	require.EqualValues(t, "ID3", getOrders[0].OrderId)

	// too large offset
	getOrders, total = orm.GetOrderList("addr1", "", "", false, 3, 10, 0, 0, false)
	require.EqualValues(t, 3, total)
	require.EqualValues(t, 0, len(getOrders))

	// GetOrderListV2 : other order: close,filled,cancel,……
	otherOrdersV2 := orm.GetOrderListV2(types.TestTokenPair, "addr1", types.BuyOrder, false, "10", "300", 1)
	require.Equal(t, 1, len(otherOrdersV2))
	require.Equal(t, updateOrders[2], &otherOrdersV2[0])

	// v2 GetOrderById
	ordersByExistId := orm.GetOrderById("ID1")
	require.EqualValues(t, updateOrders[0], ordersByExistId)
	ordersByNotExistId := orm.GetOrderById("not_exist_ID")
	require.Nil(t, ordersByNotExistId)

}

func TestSqlite3_Orders(t *testing.T) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)
	testORMOrders(t, orm)
}

// Transactions
func testORMTransactions(t *testing.T, orm *ORM) {
	orm, dbPath := MockSqlite3ORM()
	defer DeleteDB(dbPath)

	txs := []*types.Transaction{
		{"hash1", types.TxTypeTransfer, "addr1", common.TestToken, types.TxSideFrom, "10.0", "0.1" + common.NativeToken, 100},
		{"hash2", types.TxTypeOrderNew, "addr1", types.TestTokenPair, types.TxSideBuy, "10.0", "0.1" + common.NativeToken, 300},
		{"hash3", types.TxTypeOrderCancel, "addr1", types.TestTokenPair, types.TxSideSell, "10.0", "0.1" + common.NativeToken, 200},
		{"hash4", types.TxTypeTransfer, "addr2", common.TestToken, types.TxSideTo, "10.0", "0.1" + common.NativeToken, 100},
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

	noTS := orm.GetDealsMaxTimestamp()
	assert.True(t, noTS == -1)

}
func testORMBatchInsert(t *testing.T, orm *ORM) {
	newOrders := []*types.Order{}

	for i := 0; i < 2000; i++ {
		oid := fmt.Sprintf("FAKEID-%04d", i)
		o := types.Order{"hash1", oid, "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "1.3", "1.5", 100}
		newOrders = append(newOrders, &o)
	}

	updatedOrders := []*types.Order{
		{"hash2", "FAKEID-0002", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "1.4", "1.7", 100},
	}

	txs := []*types.Transaction{
		{"FAKEIDHash-1", types.TxTypeTransfer, "addr1", common.TestToken, types.TxSideFrom, "10.0", "0.1" + common.NativeToken, 100},
	}

	addDeals := []*types.Deal{
		{100, 1, "FAKEID-0001", "addr1", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
		{300, 3, "FAKEID-0002", "addr1", "btc_" + common.NativeToken, types.BuyOrder, 10.0, 1.0, "0"},
		{200, 2, "FAKEID-0003", "addr1", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
		{400, 1, "FAKEID-0004", "addr2", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
	}

	mrs := []*types.MatchResult{
		{100, 1, types.TestTokenPair, 10.0, 1.0},
	}

	feeDetails := []*token.FeeDetail{
		{"addr1", "0.1" + common.NativeToken, types.FeeTypeOrderCancel, 100},
		{"addr1", "0.5" + common.NativeToken, types.FeeTypeOrderNew, 300},
		{"addr1", "0.2" + common.NativeToken, types.FeeTypeOrderDeal, 200},
		{"addr2", "0.3" + common.NativeToken, types.FeeTypeOrderDeal, 100},
	}

	resultMap, e := orm.BatchInsertOrUpdate(newOrders, updatedOrders, addDeals, mrs, feeDetails, txs)
	require.True(t, resultMap != nil && e == nil)

	require.True(t, resultMap != nil && resultMap["newOrders"] == 2000)
	require.True(t, resultMap != nil && resultMap["updatedOrders"] == 1)
	require.True(t, resultMap != nil && resultMap["transactions"] == 1)
	require.True(t, resultMap != nil && resultMap["deals"] == 4)
	require.True(t, resultMap != nil && resultMap["feeDetails"] == 4)

	resultMap2, e2 := orm.BatchInsertOrUpdate(newOrders, updatedOrders, addDeals, mrs, feeDetails, txs)
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
	products, err := closeORM.GetAllUpdatedProducts(0, -1)
	require.Error(t, err)
	require.Nil(t, products)
	klines, err := closeORM.GetLatestKlineM1ByProduct("abc_bcd", 100)
	require.Nil(t, klines)
	require.Error(t, err)
	matches, err := closeORM.GetLatestMatchResults(types.TestTokenPair, 1)
	require.Equal(t, 0, len(matches))
	require.Error(t, err)
	matchResults, err := closeORM.GetMatchResultsByTimeRange("", 100, 500)
	require.Equal(t, 0, len(matchResults))
	require.Error(t, err)
	deals, err := closeORM.GetLatestDeals("", 100)
	require.Equal(t, 0, len(deals))
	require.Error(t, err)
	deals, err = closeORM.GetDealsByTimestampRange("", 0, time.Now().Unix())
	require.Nil(t, deals)
	require.Error(t, err)

	// delete after close DB
	err = closeORM.deleteKlinesBefore(1, &types.KlineM15{})
	require.Error(t, err)
	err = closeORM.deleteKlinesAfter(1, "", &types.KlineM15{})
	require.Error(t, err)
	err = closeORM.DeleteDealBefore(time.Now().Unix() + 1)
	require.Error(t, err)
	err = closeORM.DeleteMatchResultBefore(1000)
	require.Error(t, err)

	// insert after close DB
	cnt, err := closeORM.AddMatchResults([]*types.MatchResult{
		{100, 1, types.TestTokenPair, 10.0, 1.0},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddDeals([]*types.Deal{
		{100, 1, "FAKEID-0001", "addr1", types.TestTokenPair, types.BuyOrder, 10.0, 1.0, "0"},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddFeeDetails([]*token.FeeDetail{
		{"addr1", "0.1" + common.NativeToken, types.FeeTypeOrderCancel, 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddOrders([]*types.Order{
		{"hash1", "ID1", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "0", "1.1", 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.AddTransactions([]*types.Transaction{
		{"hash1", types.TxTypeTransfer, "addr1", common.TestToken, types.TxSideFrom, "10.0", "0.1" + common.NativeToken, 100},
	})
	require.Error(t, err)
	require.Equal(t, 0, cnt)

	cnt, err = closeORM.UpdateOrders([]*types.Order{
		{"hash1", "ID1", "addr1", types.TestTokenPair, types.BuyOrder, "10.0", "1.1", 0, "0", "1.1", 100},
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
