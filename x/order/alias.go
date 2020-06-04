package order

import (
	"github.com/okex/okchain/x/order/client/cli"
	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
)

// nolint
// const params aliases
const (
	ModuleName              = types.ModuleName
	RouterKey               = types.RouterKey
	QuerierRoute            = types.QuerierRoute
	DefaultParamspace       = types.DefaultParamspace
	DefaultCodespace        = types.DefaultCodespace
	OrderStoreKey           = types.OrderStoreKey
	OrdinaryOrder           = types.OrdinaryOrder
	MarginOrder             = types.MarginOrder
	BuyOrder                = types.BuyOrder
	SellOrder               = types.SellOrder
	DefaultNewOrderFeeRatio = types.DefaultNewOrderFeeRatio
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
	MarginKeeper     = keeper.MarginKeeper
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
	GetCmdNewOrder    = cli.GetCmdNewOrder
	GetCmdCancelOrder = cli.GetCmdCancelOrder
)
