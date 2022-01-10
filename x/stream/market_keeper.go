package stream

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/stream/pushservice/conn"
)

type BaseMarketKeeper struct {
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
	key := fmt.Sprintf("%d_%d", productID, granularity)
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
			return nil, err
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
		key := product
		k.logger.Debug("GetTickerByInstruments", "key", key)
		r, err := k.client.Get(key)
		if err != nil {
			return tickers, err
		}
		ticker := map[string]string{}
		if len(r) > 0 {
			err := json.Unmarshal([]byte(r), &ticker)
			if err == nil {
				tickers = append(tickers, ticker)
			} else {
				return tickers, nil
			}
		}
	}

	return tickers, nil
}
