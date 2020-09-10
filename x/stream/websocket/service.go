package websocket

import (
	"encoding/json"
	"fmt"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Engine
type Engine struct {
	url    string
	logger log.Logger
}

func NewEngine(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	engine := &Engine{url: url, logger: logger}
	return engine, nil
}

func (engine *Engine) URL() string {
	return engine.url
}

func (engine *Engine) NewEvent(channel string, data interface{}) (sdk.Event, error) {
	eventData, err := json.Marshal(data)
	if err != nil {
		return sdk.Event{}, err
	}
	return sdk.NewEvent(
		eventTypeBackend,
		sdk.NewAttribute("channel", channel),
		sdk.NewAttribute("data", string(eventData)),
	), nil
}

func (engine *Engine) Write(data types.IStreamData, success *bool) {
	defer func() {
		if e := recover(); e != nil {
			*success = false
			engine.logger.Error("error: WebSocketEngine Write", "err", e)
		}
	}()

	wsData := data.(*PushData)
	engine.logger.Debug(fmt.Sprintf("error: WebSocketEngine Write data:%v", wsData.RedisBlock))
	events := sdk.Events{}

	// 1. collect dex_spot/account events
	for key, value := range wsData.AccountsMap {
		// dex_spot/account:okt:okchain10q0rk...
		channel := fmt.Sprintf("%s:%s", DexSpotAccount, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	// 2. collect dex_spot/order events
	for key, value := range wsData.OrdersMap {
		// dex_spot/order:xxb_okt:okchain10q0rk...
		channel := fmt.Sprintf("%s:%s", DexSpotOrder, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	// 3. collect dex_spot/matches events
	for key, value := range wsData.MatchesMap {
		// dex_spot/matches:xxb_okt
		channel := fmt.Sprintf("%s:%s", DexSpotMatch, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	// 4. collect dex_spot/optimized_depth events
	for key, value := range wsData.DepthBooksMap {
		// dex_spot/optimized_depth:xxb_okt
		channel := fmt.Sprintf("%s:%s", DexSpotDepthBook, key)
		event, err := engine.NewEvent(channel, value)
		if err != nil {
			panic(err)
		}
		events = append(events, event)
	}

	wsData.eventMgr.EmitEvents(events)
	*success = true
}
