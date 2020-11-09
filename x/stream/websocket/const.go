package websocket

import (
	"errors"
	"time"
)

const (
	// time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// rpc Websocket timeout
	maxRPCContextTimeout = writeWait

	eventTypeBackend  = "backend"
	rpcChannelKey     = "backend.channel"
	rpcChannelDataKey = "backend.data"

	DexSpotAccount     = "dex_spot/account"
	DexSpotOrder       = "dex_spot/order"
	DexSpotMatch       = "dex_spot/matches"
	DexSpotAllTicker3s = "dex_spot/all_ticker_3s"
	DexSpotTicker      = "dex_spot/ticker"
	DexSpotDepthBook   = "dex_spot/optimized_depth"

	eventSubscribe   = "subscribe"
	eventUnsubscribe = "unsubscribe"
	eventLogin       = "dex_jwt"
)

var (
	errSubscribeParams = errors.New(`ws subscription parameter error`)
)
