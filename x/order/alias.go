package order

import (
	"github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/order/types"
)

// nolint
// const params aliases
const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultParamspace = types.DefaultParamspace
	DefaultCodespace  = types.DefaultCodespace
	OrderStoreKey     = types.OrderStoreKey
)

// nolint
// types aliases
type (
	Keeper           = keeper.Keeper
	Order            = types.Order
	DepthBook        = types.DepthBook
	MatchResult      = types.MatchResult
	Deal             = types.Deal
	Params           = types.Params
	MsgNewOrder      = types.MsgNewOrder
	MsgCancelOrder   = types.MsgCancelOrder
	MsgNewOrders     = types.MsgNewOrders
	MsgCancelOrders  = types.MsgCancelOrders
	BlockMatchResult = types.BlockMatchResult
)

// nolint
// functions aliases
var (
	RegisterCodec     = types.RegisterCodec
	DefaultParams     = types.DefaultParams
	NewMsgNewOrder    = types.NewMsgNewOrder
	NewMsgCancelOrder = types.NewMsgCancelOrder
	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	FormatOrderIDsKey = types.FormatOrderIDsKey
)
