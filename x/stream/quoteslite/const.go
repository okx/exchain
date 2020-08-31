package quoteslite

import (
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// Rpc Websocket timeout
	maxRpcContextTimeout = writeWait
)

const (
	EventTypeBackend   = "backend"
	rpcChannelKey      = "backend.channel"
	rpcChannelDataKey  = "backend.data"
	dexSpotAll3Tickers = "dex_spot/all_tickers_3s"
)
