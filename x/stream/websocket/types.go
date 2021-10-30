package websocket

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/stream/common"
	pushservice "github.com/okex/exchain/x/stream/pushservice/types"
	"github.com/okex/exchain/x/stream/types"
)

type PushData struct {
	*pushservice.RedisBlock
	eventMgr *sdk.EventManager
}

func NewPushData() *PushData {
	baseData := pushservice.NewRedisBlock()
	pd := PushData{RedisBlock: baseData, eventMgr: nil}
	return &pd
}

func (data *PushData) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper,
	dexKeeper types.DexKeeper, swapKeeper types.SwapKeeper, cache *common.Cache) {
	data.eventMgr = ctx.EventManager()
	data.RedisBlock.SetData(ctx, orderKeeper, tokenKeeper, dexKeeper, swapKeeper, cache)

	// update depthBook cache
	products := orderKeeper.GetUpdatedDepthbookKeys()
	for _, product := range products {
		depthBook := orderKeeper.GetDepthBookCopy(product)
		bookRes := pushservice.ConvertBookRes(product, orderKeeper, depthBook, 200)
		UpdateDepthBookCache(product, bookRes)
	}

}

func (data PushData) DataType() types.StreamDataKind {
	return types.StreamDataWebSocketKind
}

type EventResponse struct {
	Event   string `json:"event"`
	Success string `json:"success"`
	Channel string `json:"Channel"`
}

func (r *EventResponse) Valid() bool {
	return (len(r.Event) > 0 && len(r.Channel) > 0) || r.Event == "login"
}

type TableResponse struct {
	Table  string        `json:"table"`
	Action string        `json:"action"`
	Data   []interface{} `json:"data"`
}

func (r *TableResponse) Valid() bool {
	return (len(r.Table) > 0 || len(r.Action) > 0) && len(r.Data) > 0
}

type ErrorResponse struct {
	Event     string `json:"event"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

func (r *ErrorResponse) Valid() bool {
	return len(r.Event) > 0 && len(r.Message) > 0 && r.ErrorCode >= 30000
}

type BaseOp struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}
