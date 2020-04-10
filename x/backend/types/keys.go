// nolint
package types

const (
	// ModuleName is the name of the backend module
	ModuleName = "backend"
	// QuerierRoute is the querier route for the backend module
	QuerierRoute = ModuleName
	// RouterKey is the msg router key for the backend module
	RouterKey = ""

	// query endpoints supported by the backend querier
	QueryMatchResults = "matches"
	QueryDealList     = "deals"
	QueryFeeDetails   = "fees"
	QueryOrderList    = "orders"
	QueryTxList       = "txs"
	QueryCandleList   = "candles"
	QueryTickerList   = "tickers"

	// v2
	QueryTickerListV2   = "tickerListV2"
	QueryTickerV2       = "tickerV2"
	QueryInstrumentsV2  = "instrumentsV2"
	QueryOrderListV2    = "orderListV2"
	QueryOrderV2        = "orderV2"
	QueryCandleListV2   = "candlesV2"
	QueryMatchResultsV2 = "matchesV2"
	QueryFeeDetailsV2   = "feesV2"
	QueryDealListV2     = "dealsV2"
	QueryTxListV2       = "txsV2"

	// kline const

	Kline1GoRoutineWaitInSecond = 5
	KlinexGoRoutineWaitInSecond = 10

	SecondsInADay = 24 * 60 * 60
)
