package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/rpc/core"
	tm "github.com/tendermint/tendermint/types"
	//sdk "github.com/cosmos/cosmos-sdk/types"
)

// Backend event type for EventBus
const (
	EventTypeBackend = "backend"
)

var (
	RPCCdcRegistered  = false
	EventQueryBackend = tm.QueryForEvent(EventTypeBackend)
)

type EventDataTicker struct {
	Symbol    string `json:"symbol"`
	Product   string `json:"product"`
	Timestamp int64  `json:"timestamp"`
}

type EventDataBackend struct {
	Timestamp string `json:"timestamp"`
}

func QueryForEvent(eventType string) string {
	return fmt.Sprintf("%s='%s'", EventTypeBackend, eventType)
}

func RegisterEventDatas(cdc *codec.Codec) {
	//cdc.RegisterInterface((*tm.TMEventData)(nil), nil)
	//cdc.RegisterConcrete(EventDataTicker{}, "okchain/event/Ticker", nil)
	cdc.RegisterConcrete(EventDataBackend{}, "tendermint/event/Backend", nil)
}

func PublishBackend(backend *EventDataBackend) error {

	// websocket extend, register event type to rpc cdc
	if !RPCCdcRegistered && core.GetCoreCdc() != nil {
		RegisterCodec(core.GetCoreCdc())
		RPCCdcRegistered = true
	}
	eventBus := core.GetEventBus()
	if eventBus != nil {
		return eventBus.Publish(EventTypeBackend, *backend)
	}
	return nil
}
