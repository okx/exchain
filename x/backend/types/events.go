package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
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
	cdc.RegisterConcrete(EventDataBackend{}, "tendermint/event/Backend", nil)
}