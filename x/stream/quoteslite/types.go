package quoteslite

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/stream/common"
	pushservice "github.com/okex/okchain/x/stream/pushservice/types"
	types2 "github.com/okex/okchain/x/stream/pushservice/types"
	"github.com/okex/okchain/x/stream/types"
)

type WebSocketPushData struct {
	*types2.RedisBlock
	eventMgr *sdk.EventManager
}

func NewWebSocketPushData() *WebSocketPushData {
	baseData := types2.NewRedisBlock()
	pd := WebSocketPushData{RedisBlock: baseData, eventMgr: nil}
	return &pd
}

func (data *WebSocketPushData) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper,
	tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, cache *common.Cache) {
	data.eventMgr = ctx.EventManager()
	data.RedisBlock.SetData(ctx, orderKeeper, tokenKeeper, dexKeeper, cache)

	// update depthBook cache
	products := orderKeeper.GetUpdatedDepthbookKeys()
	for _, product := range products {
		depthBook := orderKeeper.GetDepthBookCopy(product)
		bookRes := pushservice.ConvertBookRes(product, orderKeeper, depthBook, 200)
		UpdateDepthBookCache(product, bookRes)
	}

}

func (data WebSocketPushData) DataType() types.StreamDataKind {
	return types.StreamDataWebSocketKind
}
