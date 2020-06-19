package poolswap

import (
	"github.com/okex/okchain/x/poolswap/keeper"
	"github.com/okex/okchain/x/poolswap/types"
)

const (
	// nolint
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	// nolint
	NewKeeper          = keeper.NewKeeper
	NewQuerier         = keeper.NewQuerier
	RegisterCodec      = types.RegisterCodec
	NewMsgAddLiquidity = types.NewMsgAddLiquidity

	// variable aliases
	// nolint
	ModuleCdc = types.ModuleCdc
)

type (
	// nolint
	Keeper = keeper.Keeper
	Params = types.Params

	// nolint
	SwapTokenPair = types.SwapTokenPair
)
