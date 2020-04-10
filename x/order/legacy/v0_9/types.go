package v0_9

import (
	"github.com/okex/okchain/x/order/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	GenesisState struct {
		Params     types.Params   `json:"params"`
		OpenOrders []*types.Order `json:"open_orders"`
	}
)
