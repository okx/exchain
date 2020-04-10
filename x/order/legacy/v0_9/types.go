// nolint
package v0_9

import (
	"github.com/okex/okchain/x/order/types"
)

// nolint
const (
	ModuleName = types.ModuleName
)

// nolint
type (
	GenesisState struct {
		Params     types.Params   `json:"params"`
		OpenOrders []*types.Order `json:"open_orders"`
	}
)
