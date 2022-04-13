package temp

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

type AccountI interface {
	GetAddress() sdk.AccAddress
}
