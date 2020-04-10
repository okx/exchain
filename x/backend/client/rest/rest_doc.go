// Package rest API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     BasePath: /
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
// swagger:meta
package rest

import (
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token"
)

// swagger:route GET /candles/{instrumentId} backend getCandles
//
// Get lastest k1min candles by instrument_id
//
//     Schemes: http, https
//     Responses:
//       200: CandleResponse

// swagger:parameters getCandles
type CandleParam struct {
	// instrument or product name,
	// Required: true
	// in: path
	InstrumentId string `json:"instrument_id"`

	// 时间颗粒度，时间粒度，以秒为单位，如[60/180/300/900/1800/3600/7200/14400/21600/43200/86400/604800]
	// Required: true
	// in: query
	Granularity string `json:"granularity"`

	// 获取k线数据的数量，最多1000条
	// Required: true
	// in: query
	Size string `json:"size"`
}

// swagger:response CandleResponse
type CandleResponse struct {
	// in: body
	Data      [][]string `json:"data"`
	Code      int        `json:"code"`
	DetailMsg string     `json:"detail_msg"`
	Msg       string     `json:"msg"`
}

// swagger:route GET /tickers backend getTickers
//
// 获取所有的交易行情信息
//
//     Schemes: http, https
//     Responses:
//       200: TickersResponse

// swagger:parameters getTickers
type TickersParam struct {
	// 获取所有的交易行情信息
	// Required: true
	// in: path
	Count int `json:"count"`

	// 1表示升序，0表示降序
	// Required: false
	// in: query
	Sort string `json:"sort"`
}

// Tickers Response
// swagger:response TickersResponse
type TickersResponse struct {
	// in: body
	Data []types.Ticker `json:"data"`
}

// swagger:route GET /order/list/{openOrClosed} backend getOrderList
//
// Get order list
//
//     Schemes: http, https
//     Responses:
//       200: OrderListResponse

// swagger:parameters getOrderList
type OrderListParam struct {
	// open/closed, get open/closed orders
	// Required: true
	// in: path
	OpenOrClosed string `json:"open_or_closed"`
	// user address
	// Required: true
	// in: query
	Address string `json:"address"`
	// token pair string
	// Required: false
	// in: query
	Product string `json:"product"`
	//side param,BUY or SELL
	//Required: false
	//in: query
	Side string `json:"side"`
	// page param
	// Required: false
	// in: query
	Page int `json:"page"`
	// PerPage param
	// Required: false
	// in: query
	PerPage int `json:"per_page"`

	// Start Timestamp
	// Required: false
	// in: query
	Start int64 `json:"start"`

	// End Timestamp
	// Required: false
	// in: query
	End int64 `json:"end"`
	// whether hide orders that have no fills, 0/1
	// Required: false
	// in: query
	HideNoFill int `json:"hide_no_fill"`
}

type ListResponseWithOrder struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detail_msg"`
	Data      struct {
		Data      []types.Order    `json:"data"`
		ParamPage common.ParamPage `json:"param_page"`
	} `json:"data"`
}

// Order list
// swagger:response OrderListResponse
type OrderListResponse struct {
	// in: body
	Body ListResponseWithOrder
}

// swagger:route GET /transactions backend getTxList
//
// Get tx list
//
//     Schemes: http, https
//     Responses:
//       200: TxListResponse

// swagger:parameters getTxList
type TxListParam struct {
	// user address
	// Required: true
	// in: query
	Address string `json:"address"`
	// tx type: 1:Transfer, 2:NewOrder, 3:CancelOrder
	// Required: false
	// in: query
	Type int `json:"type"`
	// start timestamp
	// Required: false
	// in: query
	Start int `json:"start"`
	// end timestamp
	// Required: false
	// in: query
	End int `json:"end"`
	// page param
	// Required: false
	// in: query
	Page int `json:"page"`
	// page param
	// Required: false
	// in: query
	PerPage int `json:"per_page"`
}

type ListResponseWithTx struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detail_msg"`
	Data      struct {
		Data      []types.Transaction `json:"data"`
		ParamPage common.ParamPage    `json:"param_page"`
	} `json:"data"`
}

// Transaction list
// swagger:response TxListResponse
type TxListResponse struct {
	// in: body
	Body ListResponseWithTx
}

// swagger:route GET /matches backend getMatchResults
//
// Get match result list
//
//     Schemes: http, https
//     Responses:
//       200: MatchListResponse

// swagger:parameters getMatchResults
type MatchesParam struct {
	// token pair string
	// Required: false
	// in: query
	Product string `json:"product"`
	// start timestamp
	// Required: false
	// in: query
	Start int `json:"start"`
	// end timestamp
	// Required: false
	// in: query
	End int `json:"end"`
	// page param
	// Required: false
	// in: query
	Page int `json:"page"`
	// page param
	// Required: false
	// in: query
	PerPage int `json:"per_page"`
}

type ListResponseWithMatch struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detail_msg"`
	Data      struct {
		Data      []types.MatchResult `json:"data"`
		ParamPage common.ParamPage    `json:"param_page"`
	} `json:"data"`
}

// Match Result List
// swagger:response MatchListResponse
type MatchListResponse struct {
	// in: body
	Body ListResponseWithMatch
}

// swagger:route GET /deals backend getDealList
//
// Get deal list
//
//     Schemes: http, https
//     Responses:
//       200: DealListResponse

// swagger:parameters getDealList
type DealsParam struct {
	// user address
	// Required: false
	// in: query
	Address string `json:"address"`
	// token pair string
	// Required: false
	// in: query
	Product string `json:"product"`
	//side param,BUY or SELL
	//Required: false
	//in: query
	Side string `json:"side"`
	// start timestamp
	// Required: false
	// in: query
	Start int `json:"start"`
	// end timestamp
	// Required: false
	// in: query
	End int `json:"end"`
	// page param
	// Required: false
	// in: query
	Page int `json:"page"`
	// page param
	// Required: false
	// in: query
	PerPage int `json:"per_page"`
}

type ListResponseWithDeal struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detail_msg"`
	Data      struct {
		Data      []types.Deal     `json:"data"`
		ParamPage common.ParamPage `json:"param_page"`
	} `json:"data"`
}

// Deal list
// swagger:response DealListResponse
type DealListResponse struct {
	// in: body
	Body ListResponseWithDeal
}

// swagger:route GET /fees backend getFeeDetails
//
// Get fee detail list
//
//     Schemes: http, https
//     Responses:
//       200: FeaDetailListResponse

// swagger:parameters getFeeDetails
type FeeDetailsParam struct {
	// user address
	// Required: true
	// in: query
	Address string `json:"address"`
	// page param
	// Required: false
	// in: query
	Page int `json:"page"`
	// page param
	// Required: false
	// in: query
	PerPage int `json:"per_page"`
}

type ListResponseWithFeeDetail struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	DetailMsg string `json:"detail_msg"`
	Data      struct {
		Data      []token.FeeDetail `json:"data"`
		ParamPage common.ParamPage  `json:"param_page"`
	} `json:"data"`
}

// Fee detail list
// swagger:response FeaDetailListResponse
type FeaDetailListResponse struct {
	// in: body
	Body ListResponseWithFeeDetail
}

// swagger:route GET /block_tx_hashes/{blockHeight} backend blockTxHashes
//
// Get tx hashes in the block of the blockHeight
//
//     Schemes: http, https
//     Responses:
//       200: BlockTxHashesResponse

// swagger:parameters blockTxHashes
type BlockTxHashesParam struct {
	// height of block
	// Required: true
	// in: path
	BlockHeight int `json:"block_height"`
}

// BlockTxHashes Response
// swagger:response BlockTxHashesResponse
type BlockTxHashesResponse struct {
	// in: body
	Body []string
}

// swagger:route GET /tickers/{instrumentId} backend getTicker
//
// Get candles of the instrument
//
//     Schemes: http, https
//     Responses:
//       200: TickerResponse

// swagger:parameters getTicker
type TickerParam struct {
	// instrument or product name,
	// Required: true
	// in: path
	InstrumentId string `json:"instrument_id"`
}

// Ticker Response
// swagger:response TickerResponse
type TickerResponse struct {
	// in: body
	Data []types.Ticker `json:"data"`
}
