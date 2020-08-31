package stream

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/stream/pulsarclient"
	"github.com/okex/okchain/x/stream/pushservice/conn"
	"github.com/pkg/errors"
)

// expected market keeper which would get data from pulsar & redis
//type MarketKeeper interface {
//	GetTickerByInstrument(instrument string) *backend.Ticker
//	GetTickerByInstruments(instruments []string) map[string]backend.Ticker
//	GetKlineByInstrument(instrument string, granularity, size int) []backend.BaseKline
//}
type MarketKeeper backend.MarketKeeper

type BaseMarketKeeper struct {
}

// QUO:OPT_KLINE:${BIZ_TYPE}_${MARKET_ID}_${MARKET_TYPE}:${GRANULARITY}
func (k *BaseMarketKeeper) getLatestCandlesKey(productID uint64, granularity int) string {
	ptn := fmt.Sprintf("QUO:OPT_KLINE:%d_%d_%d:%d", pulsarclient.MARKET_CAL_SERVICE_DEX_SPOT_BIZ_TYPE,
		productID, pulsarclient.MARKET_CAL_SERVICE_DEX_SPOT_MARKET_TYPE, granularity)
	return ptn
}

// key: P3C:dex_spot:ticker:xxb_okb:
func (k *BaseMarketKeeper) getTickerCacheKey(instrumentName string) (key string) {
	key = fmt.Sprintf("P3C:dex_spot:ticker:%s:", instrumentName)
	return key
}

type RedisMarketKeeper struct {
	*BaseMarketKeeper
	client *conn.Client
	logger log.Logger
}

func NewRedisMarketKeeper(client *conn.Client, logger log.Logger) *RedisMarketKeeper {
	k := RedisMarketKeeper{}
	k.BaseMarketKeeper = &BaseMarketKeeper{}
	k.client = client
	k.logger = logger
	return &k
}

func (k *RedisMarketKeeper) GetKlineByProductID(productID uint64, granularity, size int) ([][]string, error) {
	key := k.getLatestCandlesKey(productID, granularity)
	k.logger.Debug("GetKlineByInstrument", "productID", productID, "key", key)
	r, err := k.client.HGetAll(key)
	k.logger.Debug("GetKlineByInstrument", "values", r, "error", err)
	klines := make([][]string, 0, len(r))
	if len(r) == 0 {
		return klines, nil
	}

	fieldList := make([]string, 0, len(r))
	for k := range r {
		fieldList = append(fieldList, k)
	}
	// sorts  fieldList in increasing order.
	sort.Strings(fieldList)

	for _, field := range fieldList {
		timeInt, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("server error: %s, key=%s, can not convert timestamp %s", err.Error(), key, field))
		}

		values := strings.Split(r[field], "|")
		// timeInt is millisecond
		values = append([]string{time.Unix(timeInt/1000, 0).UTC().Format("2006-01-02T15:04:05.000Z")}, values...)

		klines = append(klines, values)
	}

	end := len(klines)
	if end > size {
		return klines[end-size : end], nil
	}

	return klines, nil
}

func (k *RedisMarketKeeper) GetTickerByProducts(products []string) ([]map[string]string, error) {
	var tickers []map[string]string
	k.logger.Debug("GetTickerByInstruments", "instruments", products)
	for _, product := range products {
		key := k.getTickerCacheKey(product)
		k.logger.Debug("GetTickerByInstruments", "key", key)
		r, _ := k.client.Get(key)
		ticker := map[string]string{}
		if len(r) > 0 {
			err := json.Unmarshal([]byte(r), &ticker)
			if err == nil {
				tickers = append(tickers, ticker)
			} else {
				return tickers, errors.New(fmt.Sprintf("No value found for key: %s", key))
			}
		}
	}

	return tickers, nil
}
