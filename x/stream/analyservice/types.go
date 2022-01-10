package analyservice

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/order/keeper"
	"github.com/okex/exchain/x/stream/common"
	"github.com/okex/exchain/x/stream/types"
	"github.com/okex/exchain/x/token"
)

// the data enqueue to mysql
type DataAnalysis struct {
	Height     int64                   `json:"height"`
	FeeDetails []*token.FeeDetail      `json:"fee_details"`
	DepthBook  keeper.BookRes          `json:"depth_book"`
	AccStates  []token.AccountResponse `json:"account_states"`
}

func (d *DataAnalysis) Empty() bool {
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
	d.FeeDetails = tokenKeeper.GetFeeDetailList()
}
