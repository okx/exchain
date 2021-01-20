package orm

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	okexchaincfg "github.com/cosmos/cosmos-sdk/server/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/token"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
)

// nolint
const (
	EngineTypeSqlite = okexchaincfg.BackendOrmEngineTypeSqlite
	EngineTypeMysql  = okexchaincfg.BackendOrmEngineTypeMysql
)

// nolint
type OrmEngineInfo = okexchaincfg.BackendOrmEngineInfo

// ORM is designed for deal with database by orm
// http://gorm.io/docs/query.html
type ORM struct {
	db                     *gorm.DB
	logger                 *log.Logger
	bufferLock             sync.Locker
	singleEntryLock        sync.Locker
	lastK15Timestamp       int64
	klineM15sBuffer        map[string][]types.KlineM15
	lastK1Timestamp        int64
	klineM1sBuffer         map[string][]types.KlineM1
	maxBlockTimestampMutex *sync.RWMutex
	maxBlockTimestamp      int64
}

func (o *ORM) SetMaxBlockTimestamp(maxBlockTimestamp int64) {
	o.maxBlockTimestampMutex.Lock()
	defer o.maxBlockTimestampMutex.Unlock()
	o.maxBlockTimestamp = maxBlockTimestamp
}

func (o *ORM) GetMaxBlockTimestamp() int64 {
	o.maxBlockTimestampMutex.RLock()
	defer o.maxBlockTimestampMutex.RUnlock()
	return o.maxBlockTimestamp
}

// New return pointer to ORM to deal with database，called at NewKeeper
func New(enableLog bool, engineInfo *OrmEngineInfo, logger *log.Logger) (m *ORM, err error) {
	orm := ORM{}
	var db *gorm.DB

	switch engineInfo.EngineType {
	case EngineTypeSqlite:
		_, e := os.Stat(engineInfo.ConnectStr)
		if e != nil && !os.IsExist(e) {
			dbDir := filepath.Dir(engineInfo.ConnectStr)
			if _, err := os.Stat(dbDir); err != nil {
				if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
					panic(err)
				}
				orm.Debug(fmt.Sprintf("%s created", dbDir))
			}
		}
	case EngineTypeMysql:
	default:

	}

	if db, err = gorm.Open(engineInfo.EngineType, engineInfo.ConnectStr); err != nil {
		e := fmt.Errorf(fmt.Sprintf("ConnectStr: %s, error: %+v", engineInfo.ConnectStr, err))
		panic(e)
	}

	orm.logger = logger
	orm.db = db
	orm.lastK1Timestamp = -1
	orm.lastK15Timestamp = -1
	orm.bufferLock = new(sync.Mutex)
	orm.singleEntryLock = new(sync.Mutex)
	orm.maxBlockTimestampMutex = new(sync.RWMutex)
	orm.db.LogMode(enableLog)
	orm.db.AutoMigrate(&types.MatchResult{})
	orm.db.AutoMigrate(&types.Deal{})
	orm.db.AutoMigrate(&token.FeeDetail{})
	orm.db.AutoMigrate(&types.Order{})
	orm.db.AutoMigrate(&types.Transaction{})
	orm.db.AutoMigrate(&types.SwapInfo{})
	orm.db.AutoMigrate(&types.SwapWhitelist{})
	orm.db.AutoMigrate(&types.ClaimInfo{})

	allKlinesMap := types.GetAllKlineMap()
	for _, v := range allKlinesMap {
		k := types.MustNewKlineFactory(v, nil)
		orm.db.AutoMigrate(k)
	}
	return &orm, nil
}

// Debug log  debug info when use orm
func (orm *ORM) Debug(msg string) {
	if orm.logger != nil {
		(*orm.logger).Debug(msg)
	}
}

// Error log occurred error when use orm
func (orm *ORM) Error(msg string) {
	if orm.logger != nil {
		(*orm.logger).Error(msg)
	}
}

func (orm *ORM) deferRollbackTx(trx *gorm.DB, returnErr error) {
	e := recover()
	if e != nil {
		orm.Error(fmt.Sprintf("ORM Panic : %+v", e))
		debug.PrintStack()
	}
	if e != nil || returnErr != nil {
		trx.Rollback()
	}
}

// Close close the database by orm
func (orm *ORM) Close() error {
	return orm.db.Close()
}

// nolint
func (orm *ORM) AddMatchResults(results []*types.MatchResult) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	cnt := 0
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	for _, result := range results {
		if result != nil {
			ret := tx.Create(result)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

// nolint
func (orm *ORM) GetMatchResults(product string, startTime, endTime int64,
	offset, limit int) ([]types.MatchResult, int) {
	var matchResults []types.MatchResult
	query := orm.db.Model(types.MatchResult{})

	if startTime == 0 && endTime == 0 {
		endTime = time.Now().Unix()
	}

	if product != "" {
		query = query.Where("product = ?", product)
	}

	if startTime > 0 {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("timestamp < ?", endTime)
	}
	var total int
	query.Count(&total)
	if offset >= total {
		return matchResults, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&matchResults)
	return matchResults, total
}

func (orm *ORM) deleteMatchResultBefore(timestamp int64) (err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(&types.MatchResult{}, " Timestamp < ? ", timestamp)
	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

func (orm *ORM) getMatchResultsByTimeRange(product string, startTime, endTime int64) ([]types.MatchResult, error) {
	var matchResults []types.MatchResult
	r := orm.db.Model(types.MatchResult{}).Where("Product = ? and Timestamp >= ? and Timestamp < ?",
		product, startTime, endTime).Order("Timestamp desc").Find(&matchResults)
	if r.Error == nil {
		return matchResults, nil
	}
	return matchResults, r.Error
}

func (orm *ORM) getLatestMatchResults(product string, limit int) ([]types.MatchResult, error) {
	var matchResults []types.MatchResult
	r := orm.db.Where("Product = ?", product).Order("Timestamp desc").Limit(limit).Find(&matchResults)
	if r.Error != nil {
		return matchResults, r.Error
	}

	return matchResults, r.Error
}

// nolint
func (orm *ORM) AddDeals(deals []*types.Deal) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	cnt := 0
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	for _, deal := range deals {
		if deal != nil {
			ret := tx.Create(deal)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

func (orm *ORM) deleteDealBefore(timestamp int64) (err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(&types.Deal{}, " Timestamp < ? ", timestamp)
	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

func (orm *ORM) getLatestDeals(product string, limit int) ([]types.Deal, error) {
	var deals []types.Deal
	r := orm.db.Where("Product = ?", product).Order("Timestamp desc").Limit(limit).Find(&deals)
	if r.Error != nil {
		return nil, r.Error
	}

	return deals, r.Error
}

// nolint
func (orm *ORM) GetDeals(address, product, side string, startTime, endTime int64, offset, limit int) ([]types.Deal, int) {
	var deals []types.Deal
	query := orm.db.Model(types.Deal{})

	if startTime == 0 && endTime == 0 {
		endTime = time.Now().Unix()
	}

	if address != "" {
		query = query.Where("sender = ?", address)
	}
	if product != "" {
		query = query.Where("product = ?", product)
	}

	if side != "" {
		query = query.Where("side = ?", side)
	}

	if startTime > 0 {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("timestamp < ?", endTime)
	}
	var total int
	query.Count(&total)
	if offset >= total {
		return deals, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&deals)
	return deals, total
}

// nolint
func (orm *ORM) GetDexFees(dexHandlingAddr, product string, offset, limit int) ([]types.DexFees, int) {
	var deals []types.Deal
	query := orm.db.Model(types.Deal{})

	if dexHandlingAddr != "" {
		query = query.Where("fee_receiver = ?", dexHandlingAddr)
	}
	if product != "" {
		query = query.Where("product = ?", product)
	}

	var total int
	query.Count(&total)
	if offset >= total {
		return nil, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&deals)
	if len(deals) == 0 {
		return nil, 0
	}
	var dexFees []types.DexFees
	for _, deal := range deals {
		dexFees = append(dexFees, types.DexFees{
			Timestamp:       deal.Timestamp,
			OrderID:         deal.OrderID,
			Product:         deal.Product,
			Fee:             deal.Fee,
			HandlingFeeAddr: deal.FeeReceiver,
		})
	}

	return dexFees, total
}

func (orm *ORM) getDealsByTimestampRange(product string, startTS, endTS int64) ([]types.Deal, error) {
	var deals []types.Deal
	r := orm.db.Model(types.Deal{}).Where(
		"Product = ? and Timestamp >= ? and Timestamp < ?", product, startTS, endTS).Order("Timestamp desc").Find(&deals)
	if r.Error == nil {
		return deals, nil
	}
	return nil, r.Error
}

func (orm *ORM) getOpenCloseDeals(startTS, endTS int64, product string) (open *types.Deal, close *types.Deal) {
	var openDeal, closeDeal types.Deal
	orm.db.Model(types.Deal{}).Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp desc").Limit(1).First(&closeDeal)
	orm.db.Model(types.Deal{}).Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp asc").Limit(1).First(&openDeal)

	if startTS <= openDeal.Timestamp && openDeal.Timestamp < endTS {
		return &openDeal, &closeDeal
	}

	return nil, nil
}

func (orm *ORM) getOpenCloseKline(startTS, endTS int64, product string, firstK interface{}, lastK interface{}) error {
	defer types.PrintStackIfPanic()

	orm.db.Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp desc").Limit(1).First(lastK)
	orm.db.Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).Order("Timestamp asc").Limit(1).First(firstK)

	return nil
}

func (orm *ORM) getMinTimestamp(tbName string) int64 {
	sql := fmt.Sprintf("select min(Timestamp) as ts from %s", tbName)
	ts := int64(-1)
	count := 0

	raw := orm.db.Raw(sql)
	raw.Count(&count)
	if count == 0 {
		return ts
	}

	if err := raw.Row().Scan(&ts); err != nil {
		orm.Error("failed to execute scan result, error:" + err.Error())
	}

	return ts
}

func (orm *ORM) getMaxTimestamp(tbName string) int64 {
	sql := fmt.Sprintf("select max(Timestamp) as ts from %s", tbName)
	ts := int64(-1)
	count := 0

	raw := orm.db.Raw(sql)
	raw.Count(&count)
	if count == 0 {
		return ts
	}

	if err := raw.Row().Scan(&ts); err != nil {
		orm.Error("failed to execute scan result, error:" + err.Error())
	}

	return ts
}

func (orm *ORM) getMergingKlineTimestamp(tbName string, timestamp int64) int64 {
	sql := fmt.Sprintf("select max(Timestamp) as ts from %s where Timestamp <=%d", tbName, timestamp)
	ts := int64(-1)
	count := 0

	raw := orm.db.Raw(sql)
	raw.Count(&count)
	if count == 0 {
		return ts
	}

	if err := raw.Row().Scan(&ts); err != nil {
		orm.Error("failed to execute scan result, error:" + err.Error())
	}

	return ts
}

func (orm *ORM) getDealsMinTimestamp() int64 {
	return orm.getMinTimestamp("deals")
}

func (orm *ORM) getDealsMaxTimestamp() int64 {
	return orm.getMaxTimestamp("deals")
}

func (orm *ORM) getKlineMaxTimestamp(k types.IKline) int64 {
	return orm.getMaxTimestamp(k.GetTableName())
}

func (orm *ORM) getKlineMinTimestamp(k types.IKline) int64 {
	return orm.getMinTimestamp(k.GetTableName())
}

func (orm *ORM) getMergeResultMinTimestamp() int64 {
	return orm.getMinTimestamp("match_results")
}

func (orm *ORM) getMergeResultMaxTimestamp() int64 {
	return orm.getMaxTimestamp("match_results")

}

// nolint
type IKline1MDataSource interface {
	getDataSourceMinTimestamp() int64
	getMaxMinSumByGroupSQL(startTS, endTS int64) string
	getOpenClosePrice(startTS, endTS int64, product string) (float64, float64)
}

// nolint
type DealDataSource struct {
	orm *ORM
}

func (dm *DealDataSource) getDataSourceMinTimestamp() int64 {
	return dm.orm.getDealsMinTimestamp()
}

func (dm *DealDataSource) getMaxMinSumByGroupSQL(startTS, endTS int64) string {
	sql := fmt.Sprintf("select product, sum(Quantity) as quantity, max(Price) as high, min(Price) as low, count(price) as cnt from deals "+
		"where Timestamp >= %d and Timestamp < %d and Side = 'BUY' group by product", startTS, endTS)
	return sql
}

func (dm *DealDataSource) getOpenClosePrice(startTS, endTS int64, product string) (float64, float64) {
	openDeal, closeDeal := dm.orm.getOpenCloseDeals(startTS, endTS, product)
	return openDeal.Price, closeDeal.Price
}

// nolint
type MergeResultDataSource struct {
	Orm *ORM
}

func (dm *MergeResultDataSource) getDataSourceMinTimestamp() int64 {
	return dm.Orm.getMergeResultMinTimestamp()
}

func (dm *MergeResultDataSource) getMaxMinSumByGroupSQL(startTS, endTS int64) string {
	sql := fmt.Sprintf("select product, sum(Quantity) as quantity, max(Price) as high, min(Price) as low, count(price) as cnt from match_results "+
		"where Timestamp >= %d and Timestamp < %d group by product", startTS, endTS)
	return sql
}

func (dm *MergeResultDataSource) getOpenClosePrice(startTS, endTS int64, product string) (float64, float64) {
	openDeal, closeDeal := dm.Orm.getOpenCloseDeals(startTS, endTS, product)
	return openDeal.Price, closeDeal.Price
}

// CreateKline1M batch insert into Kline1M
func (orm *ORM) CreateKline1M(startTS, endTS int64, dataSource IKline1MDataSource) (
	anchorEndTS int64, newProductCnt int, newKlineInfo map[string][]types.KlineM1, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	// 1. Get anchor start time.
	if endTS <= startTS {
		return -1, 0, nil, fmt.Errorf("EndTimestamp %d <= StartTimestamp %d, somewhere goes wrong", endTS, startTS)
	}

	acTS := startTS
	maxTSPersistent := orm.getKlineMaxTimestamp(&types.KlineM1{})
	if maxTSPersistent > 0 && maxTSPersistent > startTS {
		acTS = maxTSPersistent + 60
	}

	if acTS == 0 {
		minDataSourceTS := dataSource.getDataSourceMinTimestamp()
		// No Deals to handle if minDataSourceTS == -1, anchorEndTS <-- startTS
		if minDataSourceTS == -1 {
			return startTS, 0, nil, errors.New("No Deals to handled, return without converting job.")
		} else {
			acTS = minDataSourceTS
		}
	}

	anchorTime := time.Unix(acTS, 0).UTC()
	anchorStartTime := time.Date(
		anchorTime.Year(), anchorTime.Month(), anchorTime.Day(), anchorTime.Hour(), anchorTime.Minute(), 0, 0, time.UTC)
	// 2. Collect product's kline by deals
	productKlines := map[string][]types.KlineM1{}
	nextTime := anchorStartTime.Add(time.Minute)
	nextTimeStamp := nextTime.Unix()
	for nextTimeStamp <= endTS {
		sql := dataSource.getMaxMinSumByGroupSQL(anchorStartTime.Unix(), nextTime.Unix())
		orm.Debug(fmt.Sprintf("CreateKline1M sql:%s", sql))
		rows, err := orm.db.Raw(sql).Rows()

		if rows != nil && err == nil {
			for rows.Next() {
				var product string
				var quantity, high, low float64
				var cnt int

				if err = rows.Scan(&product, &quantity, &high, &low, &cnt); err != nil {
					orm.Error(fmt.Sprintf("CreateKline1M failed to execute scan result, error:%s sql:%s", err.Error(), sql))
				}
				if cnt > 0 {
					openPrice, closePrice := dataSource.getOpenClosePrice(anchorStartTime.Unix(), nextTime.Unix(), product)

					b := types.BaseKline{
						Product: product, High: high, Low: low, Volume: quantity,
						Timestamp: anchorStartTime.Unix(), Open: openPrice, Close: closePrice}
					k1min := types.NewKlineM1(&b)

					klines := productKlines[product]
					if klines == nil {
						klines = []types.KlineM1{*k1min}
					} else {
						klines = append(klines, *k1min)
					}
					productKlines[product] = klines
				}
			}

			if err = rows.Close(); err != nil {
				orm.Error(fmt.Sprintf("CreateKline1M failed to execute close rows, error:%s", err.Error()))
			}
		}

		anchorStartTime = nextTime
		nextTime = anchorStartTime.Add(time.Minute)
		nextTimeStamp = nextTime.Unix()
	}

	// 3. Batch insert into Kline1Min
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)
	for _, klines := range productKlines {
		for _, kline := range klines {
			// TODO: it should be a replacement here.
			ret := tx.Create(&kline)
			if ret.Error != nil {
				orm.Error(fmt.Sprintf("CreateKline1M failed to create kline Error: %+v, kline: %s", ret.Error, kline.PrettyTimeString()))
			} else {
				orm.Debug(fmt.Sprintf("CreateKline1M success to create in %s, %s %s", kline.GetTableName(), types.TimeString(kline.Timestamp), kline.PrettyTimeString()))
			}

		}
	}
	tx.Commit()

	anchorEndTS = anchorStartTime.Unix()
	orm.Debug(fmt.Sprintf("CreateKline1M return klinesMap: %+v", productKlines))
	return anchorEndTS, len(productKlines), productKlines, nil
}

func (orm *ORM) deleteKlinesBefore(unixTS int64, kline interface{}) (err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(kline, " Timestamp < ? ", unixTS)
	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

func (orm *ORM) deleteKlinesAfter(unixTS int64, product string, kline interface{}) (err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(kline, " Timestamp >= ? and Product = ?", unixTS, product)
	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

// DeleteKlineBefore delete from kline
func (orm *ORM) DeleteKlineBefore(unixTS int64, kline interface{}) error {
	return orm.deleteKlinesBefore(unixTS, kline)
}

func (orm *ORM) deleteKlineM1Before(unixTS int64) error {
	return orm.DeleteKlineBefore(unixTS, &types.KlineM1{})
}

func (orm *ORM) getAllUpdatedProducts(anchorStartTS, anchorEndTS int64) ([]string, error) {
	midTS := anchorEndTS - int64(16*60)
	p1, e1 := orm.getAllUpdatedProductsFromTable(anchorStartTS, midTS, "kline_m15")
	if e1 != nil {
		return nil, e1
	}

	p2, e2 := orm.getAllUpdatedProductsFromTable(midTS, anchorEndTS, "match_results")
	if e2 != nil {
		return nil, e1
	}

	tmpMap := map[string]bool{}
	for _, p := range p1 {
		tmpMap[p] = true
	}

	for _, p := range p2 {
		tmpMap[p] = true
	}

	mergedKline := []string{}
	for k := range tmpMap {
		mergedKline = append(mergedKline, k)
	}

	return mergedKline, nil
}

func (orm *ORM) getAllUpdatedProductsFromTable(anchorStartTS, anchorEndTS int64, tb string) ([]string, error) {
	sql := fmt.Sprintf("select distinct(Product) from %s where Timestamp >= %d and Timestamp < %d",
		tb, anchorStartTS, anchorEndTS)

	rows, err := orm.db.Raw(sql).Rows()

	if err == nil {
		products := []string{}
		for rows.Next() {
			var product string
			if err := rows.Scan(&product); err != nil {
				orm.Error("failed to execute scan result, error:" + err.Error())
			}
			products = append(products, product)
		}
		if err = rows.Close(); err != nil {
			orm.Error("failed to execute close rows, error:" + err.Error())
		}
		return products, nil

	} else {
		return nil, err
	}
}

// nolint
func (orm *ORM) GetLatestKlinesByProduct(product string, limit int, anchorTS int64, klines interface{}) error {

	var r *gorm.DB
	if anchorTS > 0 {
		r = orm.db.Where("Timestamp < ? and Product = ?", anchorTS, product).Order("Timestamp desc").Limit(limit).Find(klines)
	} else {
		r = orm.db.Where("Product = ?", product).Order("Timestamp desc").Limit(limit).Find(klines)
	}

	return r.Error
}

func (orm *ORM) getKlinesByTimeRange(product string, startTS, endTS int64, klines interface{}) error {

	r := orm.db.Where("Timestamp >= ? and Timestamp < ? and Product = ?", startTS, endTS, product).
		Order("Timestamp desc").Find(klines)

	return r.Error
}

func (orm *ORM) getLatestKlineM1ByProduct(product string, limit int) (*[]types.KlineM1, error) {
	klines := []types.KlineM1{}
	if err := orm.GetLatestKlinesByProduct(product, limit, -1, &klines); err != nil {
		return nil, err
	} else {
		return &klines, nil
	}
}

// MergeKlineM1  merge KlineM1 data to KlineM*
func (orm *ORM) MergeKlineM1(startTS, endTS int64, destKline types.IKline) (
	anchorEndTS int64, newKlineTypeCnt int, newKlines map[string][]interface{}, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	klineM1 := types.MustNewKlineFactory("kline_m1", nil)
	// 0. destKline should not be KlineM1 & endTS should be greater than startTS
	if destKline.GetFreqInSecond() <= klineM1.(types.IKline).GetFreqInSecond() {
		return startTS, 0, nil, fmt.Errorf("destKline's updating Freq #%d# should be greater than 60", destKline.GetFreqInSecond())
	}
	if endTS <= startTS {
		return -1, 0, nil, fmt.Errorf("EndTimestamp %d <= StartTimestamp %d, somewhere goes wrong", endTS, startTS)
	}

	// 1. Get anchor start time.
	acTS := startTS
	maxTSPersistent := orm.getMergingKlineTimestamp(destKline.GetTableName(), startTS)
	if maxTSPersistent > 0 {
		acTS = maxTSPersistent
	}

	if acTS == 0 {
		minTS := orm.getKlineMinTimestamp(klineM1.(types.IKline))
		// No Deals to handle if minDealTS == -1, anchorEndTS <-- startTS
		if minTS == -1 {
			return startTS, 0, nil, errors.New("DestKline:" + destKline.GetTableName() + ". No KlineM1 to handled, return without converting job.")
		} else {
			acTS = minTS
		}
	}

	var anchorStartTime time.Time
	if maxTSPersistent > 0 {
		anchorTime := time.Unix(acTS, 0).UTC()
		anchorStartTime = time.Date(
			anchorTime.Year(), anchorTime.Month(), anchorTime.Day(), anchorTime.Hour(), anchorTime.Minute(), anchorTime.Second(), 0, time.UTC)
	} else {
		anchorStartTime = time.Unix(destKline.GetAnchorTimeTS(acTS), 0).UTC()
	}

	// 2. Get anchor end time.
	anchorEndTime := endTS
	orm.Debug(fmt.Sprintf("[backend] MergeKlineM1 KlinesMX-#%d# [%s, %s]",
		destKline.GetFreqInSecond(), types.TimeString(anchorStartTime.Unix()), types.TimeString(anchorEndTime)))

	// 3. Collect product's kline by deals
	productKlines := map[string][]interface{}{}
	interval := time.Duration(int(time.Second) * destKline.GetFreqInSecond())
	nextTime := anchorStartTime.Add(interval)
	for nextTime.Unix() <= anchorEndTime {

		sql := fmt.Sprintf("select %d, product, sum(volume) as volume, max(high) as high, min(low) as low, count(*) as cnt from %s "+
			"where Timestamp >= %d and Timestamp < %d group by product", anchorStartTime.Unix(), klineM1.(types.IKline).GetTableName(), anchorStartTime.Unix(), nextTime.Unix())
		orm.Debug(fmt.Sprintf("[backend] MergeKlineM1 KlinesMX-#%d# sql=%s",
			destKline.GetFreqInSecond(), sql))
		rows, err := orm.db.Raw(sql).Rows()

		if rows != nil && err == nil {
			for rows.Next() {
				var product string
				var quantity, high, low float64
				var cnt int
				var ts int64

				if err = rows.Scan(&ts, &product, &quantity, &high, &low, &cnt); err != nil {
					orm.Error("failed to execute scan result, error:" + err.Error())
				}
				if cnt > 0 {

					openKline := types.MustNewKlineFactory(klineM1.(types.IKline).GetTableName(), nil)
					closeKline := types.MustNewKlineFactory(klineM1.(types.IKline).GetTableName(), nil)
					err = orm.getOpenCloseKline(anchorStartTime.Unix(), nextTime.Unix(), product, openKline, closeKline)
					if err != nil {
						orm.Error(fmt.Sprintf("failed to get open and close kline, error: %s", err.Error()))
					}
					b := types.BaseKline{
						Product: product, High: high, Low: low, Volume: quantity, Timestamp: anchorStartTime.Unix(),
						Open: openKline.(types.IKline).GetOpen(), Close: closeKline.(types.IKline).GetClose()}

					newDestK := types.MustNewKlineFactory(destKline.GetTableName(), &b)

					klines := productKlines[product]
					if klines == nil {
						klines = []interface{}{newDestK}
					} else {
						klines = append(klines, newDestK)
					}
					productKlines[product] = klines
				}
			}
			if err = rows.Close(); err != nil {
				orm.Error("failed to execute close rows, error:" + err.Error())
			}
		}

		anchorStartTime = nextTime
		nextTime = anchorStartTime.Add(interval)
	}

	// 4. Batch insert into Kline1Min
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)
	for _, klines := range productKlines {
		for _, kline := range klines {
			ret := tx.Delete(kline).Create(kline)
			if ret.Error != nil {
				orm.Error(fmt.Sprintf("Error: %+v, kline: %s", ret.Error, kline.(types.IKline).PrettyTimeString()))
			} else {
				orm.Debug(fmt.Sprintf("%s %s", types.TimeString(kline.(types.IKline).GetTimestamp()), kline.(types.IKline).PrettyTimeString()))
			}
		}
	}
	tx.Commit()

	anchorEndTS = anchorStartTime.Unix()
	return anchorEndTS, len(productKlines), productKlines, nil
}

// RefreshTickers Latest 24H KlineM1 to Ticker
func (orm *ORM) RefreshTickers(startTS, endTS int64, productList []string) (m map[string]*types.Ticker, err error) {

	orm.Debug(fmt.Sprintf("[backend] entering RefreshTickers, expected TickerTimeRange: [%d, %d)=[%s, %s), expectedProducts: %+v",
		startTS, endTS, types.TimeString(startTS), types.TimeString(endTS), productList))
	orm.bufferLock.Lock()
	defer orm.bufferLock.Unlock()

	// 1. Get updated product by Deals & KlineM15 in latest 120 seconds
	km1 := types.MustNewKlineFactory("kline_m1", nil)
	km15 := types.MustNewKlineFactory("kline_m15", nil)
	if len(productList) == 0 {
		productList, err = orm.getAllUpdatedProducts(endTS-types.SecondsInADay*14, endTS)
		if err != nil {
			return nil, err
		}
	}

	if len(productList) == 0 {
		return nil, nil
	}

	// 2. Update Buffer.
	// 	2.1 For each product, get latest [anchorKM15TS-types.SECONDS_IN_A_DAY, anchorKM15TS) KlineM15 list
	anchorKM15TS := km15.(types.IKline).GetAnchorTimeTS(endTS-types.KlinexGoRoutineWaitInSecond) - 60*15
	if anchorKM15TS != orm.lastK15Timestamp {
		orm.lastK15Timestamp = anchorKM15TS
		orm.klineM15sBuffer = map[string][]types.KlineM15{}
	}

	finalStartTS := km15.(types.IKline).GetAnchorTimeTS(endTS) - types.SecondsInADay
	for _, p := range productList {
		existsKM15 := orm.klineM15sBuffer[p]
		if len(existsKM15) > 0 {
			continue
		}

		klineM15s := []types.KlineM15{}

		if err = orm.getKlinesByTimeRange(p, finalStartTS, anchorKM15TS, &klineM15s); err != nil {
			orm.Error(fmt.Sprintf("failed to get kline between %d and %d, error: %s", finalStartTS, anchorKM15TS, err.Error()))
		}
		if len(klineM15s) > 0 {
			orm.klineM15sBuffer[p] = klineM15s
		}
	}

	// 	2.2 For each product, get latest [anchorKM15TS, anchorKM1TS) KlineM1 list
	anchorKM1TS := (km1).(types.IKline).GetAnchorTimeTS(endTS-types.Kline1GoRoutineWaitInSecond) - 60
	if anchorKM1TS != orm.lastK1Timestamp {
		orm.lastK1Timestamp = anchorKM1TS
		orm.klineM1sBuffer = map[string][]types.KlineM1{}
	}

	for _, p := range productList {
		existsKM1 := orm.klineM1sBuffer[p]
		if len(existsKM1) > 0 {
			continue
		}

		klineM1s := []types.KlineM1{}
		if anchorKM1TS < anchorKM15TS {
			anchorKM1TS = anchorKM15TS
		}
		if err = orm.getKlinesByTimeRange(p, anchorKM15TS, anchorKM1TS, &klineM1s); err != nil {
			orm.Error(fmt.Sprintf("failed to get kline between %d and %d, error: %s", finalStartTS, anchorKM15TS, err.Error()))
		}
		if len(klineM1s) > 0 {
			orm.klineM1sBuffer[p] = klineM1s
		}
	}

	// 	2.3 For each product, get latest [anchorKM1TS, endTS) MatchResult list
	matchResultMap := make(map[string][]types.MatchResult)
	for _, product := range productList {
		matchResults, err := orm.getMatchResultsByTimeRange(product, anchorKM1TS, endTS)
		if err != nil {
			orm.Error(fmt.Sprintf("failed to GetMatchResultsByTimeRange, error: %s", err.Error()))
			continue
		}

		if len(matchResults) > 0 {
			matchResultMap[product] = matchResults
		}
	}

	// 3. For each updated product, generate new ticker by KlineM15 & KlineM1 & Deals in 24 Hours
	orm.Debug(fmt.Sprintf("RefreshTickers: Cache KlineM15[%s, %s), KlineM1[%s, %s), Deals[%s, %s)",
		types.TimeString(finalStartTS), types.TimeString(anchorKM15TS), types.TimeString(anchorKM15TS),
		types.TimeString(anchorKM1TS), types.TimeString(anchorKM1TS), types.TimeString(endTS)))

	orm.Debug(fmt.Sprintf("KlineM15: %+v\n KlineM1: %+v\n MatchResults: %+v\n", orm.klineM15sBuffer,
		orm.klineM1sBuffer, matchResultMap))
	tickerMap := map[string]*types.Ticker{}

	orm.Debug(fmt.Sprintf("RefreshTickers's final productList %+v", productList))
	for _, p := range productList {
		klinesM1 := orm.klineM1sBuffer[p]
		klinesM15 := orm.klineM15sBuffer[p]
		iklines := types.IKlinesDsc{}

		for idx := range klinesM1[:] {
			iklines = append(iklines, &klinesM1[idx])
		}
		for idx := range klinesM15[:] {
			iklines = append(iklines, &klinesM15[idx])
		}

		// [X] 3.1 No klinesM1 & klinesM15 found, continue
		// FLT. 20190411. Go ahead even if there's no klines.

		// 3.2 Do iklines sort desc by timestamp.

		allVolume, lowest, highest := 0.0, 0.0, 0.0
		matchResults := matchResultMap[p]

		if len(iklines) > 0 {
			sort.Sort(iklines)
			allVolume, lowest, highest = 0.0, iklines[0].GetLow(), iklines[0].GetHigh()
			for _, k := range iklines {
				orm.Debug(fmt.Sprintf("RefreshTickers, Handled Kline(%s): %s", k.GetTableName(), k.PrettyTimeString()))

				allVolume += k.GetVolume()
				if k.GetHigh() > highest {
					highest = k.GetHigh()
				}
				if k.GetLow() < lowest {
					lowest = k.GetLow()
				}
			}
		} else {
			if len(matchResults) > 0 {
				allVolume, lowest, highest = 0.0, matchResults[0].Price, matchResults[0].Price
			}
		}

		for _, match := range matchResults {
			allVolume += match.Quantity
			if match.Price > highest {
				highest = match.Price
			}
			if match.Price < lowest {
				lowest = match.Price
			}
		}

		if len(iklines) == 0 && len(matchResults) == 0 {
			latestMatches, err := orm.getLatestMatchResults(p, 1)
			if err != nil {
				orm.Debug(fmt.Sprintf("failed to GetLatestMatchResults, error: %s", err.Error()))
			}

			if len(latestMatches) == 1 {
				matchResults = latestMatches
				highest = matchResults[0].Price
				lowest = matchResults[0].Price
			} else {
				continue
			}
		}

		t := types.Ticker{}
		if len(iklines) > 0 {
			t.Open = iklines[len(iklines)-1].GetOpen()
			t.Close = iklines[0].GetClose()
		} else {
			t.Open = matchResults[len(matchResults)-1].Price
		}

		if len(matchResults) > 0 {
			t.Close = matchResults[0].Price
		}

		t.Volume = allVolume
		t.High = highest
		t.Low = lowest
		t.Symbol = p
		t.Product = p
		dClose := decimal.NewFromFloat(t.Close)
		dOpen := decimal.NewFromFloat(t.Open)
		dChange := dClose.Sub(dOpen)
		t.Change, _ = dChange.Float64()
		t.ChangePercentage = fmt.Sprintf("%.2f", t.Change*100/t.Open) + "%"
		t.Price = t.Close
		t.Timestamp = endTS
		tickerMap[p] = &t
	}

	//for k, v := range tickerMap {
	//	orm.Debug(fmt.Sprintf("RefreshTickers Ticker[%s] %s", k, v.PrettyString()))
	//}
	return tickerMap, nil
}

// AddFeeDetails insert into fees
func (orm *ORM) AddFeeDetails(feeDetails []*token.FeeDetail) (addedCnt int, err error) {

	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)
	cnt := 0

	for _, feeDetail := range feeDetails {
		if feeDetail != nil {
			ret := tx.Create(feeDetail)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

// nolint
func (orm *ORM) GetFeeDetails(address string, offset, limit int) ([]token.FeeDetail, int) {
	var feeDetails []token.FeeDetail
	query := orm.db.Model(token.FeeDetail{}).Where("address = ?", address)
	var total int
	query.Count(&total)
	if offset >= total {
		return feeDetails, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&feeDetails)
	return feeDetails, total
}

// AddOrders insert into orders
func (orm *ORM) AddOrders(orders []*types.Order) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	cnt := 0
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	for _, order := range orders {
		if order != nil {
			ret := tx.Create(order)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

// UpdateOrders return count of orders
func (orm *ORM) UpdateOrders(orders []*types.Order) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	cnt := 0
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	for _, order := range orders {
		if order != nil {
			ret := tx.Save(order)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

// nolint
func (orm *ORM) GetOrderList(address, product, side string, open bool, offset, limit int,
	startTS, endTS int64, hideNoFill bool) ([]types.Order, int) {
	var orders []types.Order

	if endTS == 0 {
		endTS = time.Now().Unix()
	}

	query := orm.db.Model(types.Order{}).Where("sender = ? AND timestamp >= ? AND timestamp < ?", address, startTS, endTS)
	if product != "" {
		query = query.Where("product = ?", product)
	}
	if open {
		query = query.Where("status = 0")
	} else {
		if hideNoFill {
			query = query.Where("status in (1, 4, 5)")
		} else {
			query = query.Where("status > 0")
		}
	}

	if side != "" {
		query = query.Where("side = ?", side)
	}

	var total int
	query.Count(&total)
	if offset >= total {
		return orders, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&orders)
	return orders, total
}

func (orm *ORM) GetAccountOrders(address string, startTS, endTS int64, offset, limit int) ([]types.Order, int) {
	var orders []types.Order

	if endTS == 0 {
		endTS = time.Now().Unix()
	}

	query := orm.db.Model(types.Order{}).Where("sender = ? AND timestamp >= ? AND timestamp < ?", address, startTS, endTS)
	var total int
	query.Count(&total)
	if offset >= total {
		return orders, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&orders)
	return orders, total
}

// AddTransactions insert into transactions, return count
func (orm *ORM) AddTransactions(transactions []*types.Transaction) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	cnt := 0
	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	for _, transaction := range transactions {
		if transaction != nil {
			ret := tx.Create(transaction)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

// nolint
func (orm *ORM) GetTransactionList(address string, txType, startTime, endTime int64, offset, limit int) ([]types.Transaction, int) {
	var txs []types.Transaction
	query := orm.db.Model(types.Transaction{}).Where("address = ?", address)
	if txType != 0 {
		query = query.Where("type = ?", txType)
	}
	if startTime > 0 {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("timestamp < ?", endTime)
	}

	var total int
	query.Count(&total)
	if offset >= total {
		return txs, total
	}

	query.Order("timestamp desc").Offset(offset).Limit(limit).Find(&txs)
	return txs, total
}

// BatchInsertOrUpdate return map mean success or fail
func (orm *ORM) BatchInsertOrUpdate(newOrders []*types.Order, updatedOrders []*types.Order, deals []*types.Deal, mrs []*types.MatchResult,
	feeDetails []*token.FeeDetail, trxs []*types.Transaction, swapInfos []*types.SwapInfo, claimInfos []*types.ClaimInfo) (resultMap map[string]int, err error) {

	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	trx := orm.db.Begin()
	defer orm.deferRollbackTx(trx, err)

	resultMap = map[string]int{}
	resultMap["newOrders"] = 0
	resultMap["updatedOrders"] = 0
	resultMap["deals"] = 0
	resultMap["feeDetails"] = 0
	resultMap["transactions"] = 0
	resultMap["matchResults"] = 0
	resultMap["swapInfos"] = 0
	resultMap["claimInfos"] = 0

	// 1. Batch Insert Orders.
	orderVItems := []string{}
	for _, order := range newOrders {
		vItem := fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%d','%s','%s','%d')",
			order.TxHash, order.OrderID, order.Sender, order.Product, order.Side, order.Price, order.Quantity,
			order.Status, order.FilledAvgPrice, order.RemainQuantity, order.Timestamp)
		orderVItems = append(orderVItems, vItem)

	}
	if len(orderVItems) > 0 {
		orderValueSQL := strings.Join(orderVItems, ", ")
		orderSQL := fmt.Sprintf("INSERT INTO `orders` (`tx_hash`,`order_id`,`sender`,`product`,`side`,`price`,"+
			"`quantity`,`status`,`filled_avg_price`,`remain_quantity`,`timestamp`) VALUES %s", orderValueSQL)
		ret := trx.Exec(orderSQL)
		if ret.Error != nil {
			return resultMap, ret.Error
		} else {
			resultMap["newOrders"] += len(orderVItems)
		}
	}

	for _, order := range updatedOrders {
		ret := trx.Save(&order)
		if ret.Error != nil {
			return resultMap, ret.Error
		} else {
			resultMap["updatedOrders"]++
		}
	}

	for _, mr := range mrs {
		ret := trx.Create(mr)
		if ret.Error != nil {
			return resultMap, ret.Error
		} else {
			resultMap["matchResults"]++
		}
	}

	// 2. Batch Insert Deals
	dealVItems := []string{}
	for _, d := range deals {
		vItem := fmt.Sprintf("('%d','%d','%s','%s','%s','%s','%f','%f','%s', '%s')",
			d.Timestamp, d.BlockHeight, d.OrderID, d.Sender, d.Product, d.Side, d.Price, d.Quantity, d.Fee, d.FeeReceiver)
		dealVItems = append(dealVItems, vItem)
	}
	if len(dealVItems) > 0 {
		dealsSQL := fmt.Sprintf("INSERT INTO `deals` (`timestamp`,`block_height`,`order_id`,`sender`,`product`,`side`,`price`,`quantity`,`fee`,`fee_receiver`) "+
			"VALUES %s", strings.Join(dealVItems, ","))
		ret := trx.Exec(dealsSQL)
		if ret.Error != nil {
			return resultMap, ret.Error
		} else {
			resultMap["deals"] += len(dealVItems)
		}
	}

	// 3. Batch Insert Transactions.
	trxVItems := []string{}
	for _, t := range trxs {
		vItem := fmt.Sprintf("('%s','%d','%s','%s','%d','%s','%s','%d')",
			t.TxHash, t.Type, t.Address, t.Symbol, t.Side, t.Quantity, t.Fee, t.Timestamp)
		trxVItems = append(trxVItems, vItem)
	}
	if len(trxVItems) > 0 {
		trxSQL := fmt.Sprintf("INSERT INTO `transactions` (`tx_hash`,`type`,`address`,`symbol`,`side`,`quantity`,`fee`,`timestamp`) "+
			"VALUES %s", strings.Join(trxVItems, ", "))
		ret := trx.Exec(trxSQL)
		if ret.Error != nil {
			return resultMap, ret.Error
		} else {
			resultMap["transactions"] += len(trxVItems)
		}

	}

	// 4. Batch Insert Fee Details.
	fdVItems := []string{}
	for _, fd := range feeDetails {
		vItem := fmt.Sprintf("('%s','%s','%s','%d')", fd.Address, fd.Fee, fd.FeeType, fd.Timestamp)
		fdVItems = append(fdVItems, vItem)
	}
	if len(fdVItems) > 0 {
		fdSQL := fmt.Sprintf("INSERT INTO `fee_details` (`address`,`fee`,`fee_type`,`timestamp`) VALUES %s", strings.Join(fdVItems, ","))
		ret := trx.Exec(fdSQL)
		if ret.Error != nil {
			return resultMap, ret.Error
		} else {
			resultMap["feeDetails"] += len(fdVItems)
		}
	}

	// 5. insert swap infos
	for _, swapInfo := range swapInfos {
		if swapInfo != nil {
			ret := trx.Create(swapInfo)
			if ret.Error != nil {
				return resultMap, ret.Error
			} else {
				resultMap["swapInfos"] += 1
			}
		}
	}

	// 6. insert claim infos
	for _, claimInfo := range claimInfos {
		if claimInfo != nil {
			ret := trx.Create(claimInfo)
			if ret.Error != nil {
				return resultMap, ret.Error
			} else {
				resultMap["claimInfos"] += 1
			}
		}
	}

	trx.Commit()

	return resultMap, nil
}

// nolint
func (orm *ORM) GetOrderListV2(instrumentID string, address string, side string, open bool, after string, before string, limit int) []types.Order {
	var orders []types.Order

	query := orm.db.Model(types.Order{})

	if instrumentID != "" {
		query = query.Where("product = ? ", instrumentID)
	}

	if after != "" {
		query = query.Where("timestamp > ? ", after)
	}

	if before != "" {
		query = query.Where("timestamp < ? ", before)
	}

	if address != "" {
		query = query.Where("sender = ? ", address)
	}

	if side != "" {
		query = query.Where("side = ? ", side)
	}

	if open {
		query = query.Where("status = 0")
	} else {
		query = query.Where("status > 0")
	}

	query.Order("timestamp desc").Limit(limit).Find(&orders)
	return orders
}

// nolint
func (orm *ORM) GetOrderByID(orderID string) *types.Order {
	var orders []types.Order

	query := orm.db.Model(types.Order{}).Where("order_id = ? ", orderID)

	query.Find(&orders)

	if len(orders) > 0 {
		return &orders[0]
	}
	return nil
}

// nolint
func (orm *ORM) GetMatchResultsV2(instrumentID string, after string, before string, limit int) []types.MatchResult {
	var matchResults []types.MatchResult
	query := orm.db.Model(types.MatchResult{})

	if instrumentID != "" {
		query = query.Where("product = ?", instrumentID)
	}

	if after != "" {
		query = query.Where("timestamp > ?", after)
	}
	if before != "" {
		query = query.Where("timestamp < ?", before)
	}

	query.Order("timestamp desc").Limit(limit).Find(&matchResults)
	return matchResults
}

// nolint
func (orm *ORM) GetFeeDetailsV2(address string, after string, before string, limit int) []token.FeeDetail {
	var feeDetails []token.FeeDetail
	query := orm.db.Model(token.FeeDetail{}).Where("address = ?", address)
	if after != "" {
		query = query.Where("timestamp > ?", after)
	}
	if before != "" {
		query = query.Where("timestamp < ?", before)
	}

	query.Order("timestamp desc").Limit(limit).Find(&feeDetails)
	return feeDetails
}

// nolint
func (orm *ORM) GetDealsV2(address, product, side string, after string, before string, limit int) []types.Deal {
	var deals []types.Deal
	query := orm.db.Model(types.Deal{})

	if address != "" {
		query = query.Where("sender = ?", address)
	}
	if product != "" {
		query = query.Where("product = ?", product)
	}
	if side != "" {
		query = query.Where("side = ?", side)
	}
	if after != "" {
		query = query.Where("timestamp > ?", after)
	}
	if before != "" {
		query = query.Where("timestamp < ?", before)
	}

	query.Order("timestamp desc").Limit(limit).Find(&deals)
	return deals
}

// nolint
func (orm *ORM) GetTransactionListV2(address string, txType int, after string, before string, limit int) []types.Transaction {
	var txs []types.Transaction
	query := orm.db.Model(types.Transaction{}).Where("address = ?", address)
	if txType != 0 {
		query = query.Where("type = ?", txType)
	}
	if after != "" {
		query = query.Where("timestamp > ?", after)
	}
	if before != "" {
		query = query.Where("timestamp < ?", before)
	}

	query.Order("timestamp desc").Limit(limit).Find(&txs)
	return txs
}
