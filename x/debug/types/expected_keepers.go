package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//type SupplyKeeper interface {
//	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
//}
//
//type TokenKeeper interface {
//	TokenExist(ctx sdk.Context, symbol string) bool
//}
//
//type StakingKeeper interface {
//	GetParams(ctx sdk.Context) stakingTypes.Params
//	SetParams(ctx sdk.Context, params stakingTypes.Params)
//}

type OrderKeeper interface {
	DumpStore(ctx sdk.Context)
}
