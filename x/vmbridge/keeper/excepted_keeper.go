package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

type EVMKeeper interface {
	GetChainConfig(ctx sdk.Context) (evmtypes.ChainConfig, bool)
	GenerateCSDBParams() evmtypes.CommitStateDBParams
	GetParams(ctx sdk.Context) evmtypes.Params
	AddInnerTx(...interface{})
	AddContract(...interface{})
}

type WASMKeeper interface {
	// Execute executes the contract instance
	Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins) ([]byte, error)
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	SetAccount(ctx sdk.Context, acc authexported.Account)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
}
