package quoteslite

import (
	"encoding/json"
	"fmt"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	okex "github.com/okex/okchain/x/stream/quoteslite/okwebsocket"
	"github.com/okex/okchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

// ***********************************
// WebsocketEngine
type WebSocketEngine struct {
	url    string
	logger log.Logger
}

func NewWebSocketEngine(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	engine := &WebSocketEngine{url: url, logger: logger}
	return engine, nil
}

func (engine *WebSocketEngine) Url() string {
	return engine.url
}

func (engine *WebSocketEngine) NewEvent(channel string, data interface{}) (sdk.Event, error) {
	eventData, err := json.Marshal(data)
	if err != nil {
		return sdk.Event{}, err
	}
	return sdk.NewEvent(
		EventTypeBackend,
		sdk.NewAttribute("channel", channel),
		sdk.NewAttribute("data", string(eventData)),
	), nil
}

func (engine *WebSocketEngine) Write(data types.IStreamData, success *bool) {
	defer func() {
		if e := recover(); e != nil {
			*success = false
			engine.logger.Error("WebSocketEngine Write", "err", e)
		}
	}()

	wsData := data.(*WebSocketPushData)
	engine.logger.Debug(fmt.Sprintf("WebSocketEngine Write data:%v", wsData.RedisBlock))
	events := sdk.Events{}

	// 1. collect dex_spot/account events
	for key, value := range wsData.AccountsMap {
		// dex_spot/account:okt:okchain10q0rk...
		channel := fmt.Sprintf("%s:%s", okex.DexSpotAccount, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	// 2. collect dex_spot/order events
	for key, value := range wsData.OrdersMap {
		// dex_spot/order:xxb_okt:okchain10q0rk...
		channel := fmt.Sprintf("%s:%s", okex.DexSpotOrder, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	// 3. collect dex_spot/matches events
	for key, value := range wsData.MatchesMap {
		// dex_spot/matches:xxb_okt
		channel := fmt.Sprintf("%s:%s", okex.DexSpotMatch, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	// 4. collect dex_spot/optimized_depth events
	for key, value := range wsData.DepthBooksMap {
		// dex_spot/optimized_depth:xxb_okt
		channel := fmt.Sprintf("%s:%s", okex.DexSpotDepthBook, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	wsData.eventMgr.EmitEvents(events)
	*success = true

	//engine.logger.Debug("stream.Events", "size", len(wsData.eventMgr.Events()))
	//for i, e := range wsData.eventMgr.ABCIEvents() {
	//	engine.logger.Debug("stream.Event", i, e.Type, "attrs", e.Attributes[0].String())
	//}

}
