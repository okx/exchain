package keeperadapter

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	supplyexported "github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
)

type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
}
