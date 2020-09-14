package analyservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/stream/common"
	"github.com/okex/okexchain/x/stream/types"
	"github.com/okex/okexchain/x/token"
)

// the data enqueue to mysql
type DataAnalysis struct {
	Height        int64                   `json:"height"`
	Deals         []*backend.Deal         `json:"deals"`
	FeeDetails    []*token.FeeDetail      `json:"feeDetails"`
	NewOrders     []*backend.Order        `json:"newOrders"`
	UpdatedOrders []*backend.Order        `json:"updatedOrders"`
	Trans         []*backend.Transaction  `json:"trans"`
	MatchResults  []*backend.MatchResult  `json:"matchResults"`
	DepthBook     keeper.BookRes          `json:"depthBook"`
	AccStates     []token.AccountResponse `json:"accStates"`
}

func (d *DataAnalysis) Empty() bool {
	if len(d.Deals) == 0 && len(d.FeeDetails) == 0 && len(d.NewOrders) == 0 &&
		len(d.UpdatedOrders) == 0 && len(d.Trans) == 0 && len(d.MatchResults) == 0 &&
		len(d.DepthBook.Asks) == 0 && len(d.DepthBook.Bids) == 0 && len(d.AccStates) == 0 {
		return true
	}
	return false
}

func (d *DataAnalysis) BlockHeight() int64 {
	return d.Height
}

func (d *DataAnalysis) DataType() types.StreamDataKind {
	return types.StreamDataAnalysisKind
}

func NewDataAnalysis() *DataAnalysis {
	return &DataAnalysis{}
}

// nolint
func (d *DataAnalysis) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper,
	tokenKeeper types.TokenKeeper, cache *common.Cache) {
	d.Height = ctx.BlockHeight()
	var err error
	d.Deals, d.MatchResults, err = common.GetDealsAndMatchResult(ctx, orderKeeper)
	if err != nil {
		ctx.Logger().Error("stream SetData error", "msg", err.Error())
	}
	d.NewOrders = common.GetNewOrders(ctx, orderKeeper)
	d.UpdatedOrders = backend.GetUpdatedOrdersAtEndBlock(ctx, orderKeeper)
	d.FeeDetails = tokenKeeper.GetFeeDetailList()
	d.Trans = cache.GetTransactions()
}
