package okex

/*
 OKEX websocket api constants
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import "errors"

const (
	WS_API_HOST = "okexcomreal.bafang.com:8443"
	WS_API_URL  = "wss://real.okex.com:8443/ws/v3"

	CHNL_FUTURES_TICKER          = "futures/ticker"          // 行情数据频道
	CHNL_FUTURES_CANDLE60S       = "futures/candle60s"       // 1分钟k线数据频道
	CHNL_FUTURES_CANDLE180S      = "futures/candle180s"      // 3分钟k线数据频道
	CHNL_FUTURES_CANDLE300S      = "futures/candle300s"      // 5分钟k线数据频道
	CHNL_FUTURES_CANDLE900S      = "futures/candle900s"      // 15分钟k线数据频道
	CHNL_FUTURES_CANDLE1800S     = "futures/candle1800s"     // 30分钟k线数据频道
	CHNL_FUTURES_CANDLE3600S     = "futures/candle3600s"     // 1小时k线数据频道
	CHNL_FUTURES_CANDLE7200S     = "futures/candle7200s"     // 2小时k线数据频道
	CHNL_FUTURES_CANDLE14400S    = "futures/candle14400s"    // 4小时k线数据频道
	CHNL_FUTURES_CANDLE21600     = "futures/candle21600"     // 6小时k线数据频道
	CHNL_FUTURES_CANDLE43200S    = "futures/candle43200s"    // 12小时k线数据频道
	CHNL_FUTURES_CANDLE86400S    = "futures/candle86400s"    // 1day k线数据频道
	CHNL_FUTURES_CANDLE604800S   = "futures/candle604800s"   // 1week k线数据频道
	CHNL_FUTURES_TRADE           = "futures/trade"           // 交易信息频道
	CHNL_FUTURES_ESTIMATED_PRICE = "futures/estimated_price" //获取预估交割价
	CHNL_FUTURES_PRICE_RANGE     = "futures/price_range"     // 限价范围频道
	CHNL_FUTURES_DEPTH           = "futures/depth"           // 深度数据频道，首次200档，后续增量
	CHNL_FUTURES_DEPTH5          = "futures/depth5"          // 深度数据频道，每次返回前5档
	CHNL_FUTURES_MARK_PRICE      = "futures/mark_price"      // 标记价格频道

	CHNL_FUTURES_ACCOUNT  = "futures/account"  // 用户账户信息频道
	CHNL_FUTURES_POSITION = "futures/position" // 用户持仓信息频道
	CHNL_FUTURES_ORDER    = "futures/order"    // 用户交易数据频道

	CHNL_SPOT_TICKER        = "spot/ticker"        // 行情数据频道
	CHNL_SPOT_CANDLE60S     = "spot/candle60s"     // 1分钟k线数据频道
	CHNL_SPOT_CANDLE180S    = "spot/candle180s"    // 3分钟k线数据频道
	CHNL_SPOT_CANDLE300S    = "spot/candle300s"    // 5分钟k线数据频道
	CHNL_SPOT_CANDLE900S    = "spot/candle900s"    // 15分钟k线数据频道
	CHNL_SPOT_CANDLE1800S   = "spot/candle1800s"   // 30分钟k线数据频道
	CHNL_SPOT_CANDLE3600S   = "spot/candle3600s"   // 1小时k线数据频道
	CHNL_SPOT_CANDLE7200S   = "spot/candle7200s"   // 2小时k线数据频道
	CHNL_SPOT_CANDLE14400S  = "spot/candle14400s"  // 4小时k线数据频道
	CHNL_SPOT_CANDLE21600   = "spot/candle21600"   // 6小时k线数据频道
	CHNL_SPOT_CANDLE43200S  = "spot/candle43200s"  // 12小时k线数据频道
	CHNL_SPOT_CANDLE86400S  = "spot/candle86400s"  // 1day k线数据频道
	CHNL_SPOT_CANDLE604800S = "spot/candle604800s" // 1week k线数据频道
	CHNL_SPOT_TRADE         = "spot/trade"         // 交易信息频道
	CHNL_SPOT_DEPTH         = "spot/depth"         // 深度数据频道，首次200档，后续增量
	CHNL_SPOT_DEPTH5        = "spot/depth5"        // 深度数据频道，每次返回前5档

	CHNL_SPOT_ACCOUNT        = "spot/account"        // 用户币币账户信息频道
	CHNL_SPOT_MARGIN_ACCOUNT = "spot/margin_account" // 用户杠杆账户信息频道
	CHNL_SPOT_ORDER          = "spot/order"          // 用户交易数据频道

	CHNL_SWAP_TICKER        = "swap/ticker"        // 行情数据频道
	CHNL_SWAP_CANDLE60S     = "swap/candle60s"     // 1分钟k线数据频道
	CHNL_SWAP_CANDLE180S    = "swap/candle180s"    // 3分钟k线数据频道
	CHNL_SWAP_CANDLE300S    = "swap/candle300s"    // 5分钟k线数据频道
	CHNL_SWAP_CANDLE900S    = "swap/candle900s"    // 15分钟k线数据频道
	CHNL_SWAP_CANDLE1800S   = "swap/candle1800s"   // 30分钟k线数据频道
	CHNL_SWAP_CANDLE3600S   = "swap/candle3600s"   // 1小时k线数据频道
	CHNL_SWAP_CANDLE7200S   = "swap/candle7200s"   // 2小时k线数据频道
	CHNL_SWAP_CANDLE14400S  = "swap/candle14400s"  // 4小时k线数据频道
	CHNL_SWAP_CANDLE21600   = "swap/candle21600"   // 6小时k线数据频道
	CHNL_SWAP_CANDLE43200S  = "swap/candle43200s"  // 12小时k线数据频道
	CHNL_SWAP_CANDLE86400S  = "swap/candle86400s"  // 1day
	CHNL_SWAP_CANDLE604800S = "swap/candle604800s" // 1week
	CHNL_SWAP_TRADE         = "swap/trade"         // 交易信息频道
	CHNL_SWAP_FUNDING_RATE  = "swap/funding_rate"  // 资金费率频道
	CHNL_SWAP_PRICE_RANGE   = "swap/price_range"   // 限价范围频道
	CHNL_SWAP_DEPTH         = "swap/depth"         // 深度数据频道，首次200档，后续增量
	CHNL_SWAP_DEPTH5        = "swap/depth5"        // 深度数据频道，每次返回前5档
	CHNL_SWAP_MARK_PRICE    = "swap/mark_price"    // 标记价格频道

	CHNL_SWAP_ACCOUNT  = "swap/account"  // 用户账户信息频道
	CHNL_SWAP_POSITION = "swap/position" // 用户持仓信息频道
	CHNL_SWAP_ORDER    = "swap/order"    // 用户交易数据频道

	CHNL_EVENT_SUBSCRIBE   = "subscribe"
	CHNL_EVENT_UNSUBSCRIBE = "unsubscribe"
)

var (
	ERR_WS_SUBSCRIOTION_PARAMS = errors.New(`ws subscription parameter error`)
	ERR_WS_CACHE_NOT_MATCH     = errors.New(`ws hot cache not matched`)
)

var (
	DefaultDataCallBack = defaultPrintData
)
