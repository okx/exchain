package keeper

import (
	"encoding/json"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/okex/okchain/x/wasm/internal/types"
)

type BankEncoder func(sender sdk.AccAddress, msg *wasmTypes.BankMsg) ([]sdk.Msg, error)
type CustomEncoder func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error)
type StakingEncoder func(sender sdk.AccAddress, msg *wasmTypes.StakingMsg) ([]sdk.Msg, error)
type WasmEncoder func(sender sdk.AccAddress, msg *wasmTypes.WasmMsg) ([]sdk.Msg, error)

type MessageEncoders struct {
	Bank    BankEncoder
	Custom  CustomEncoder
	Staking StakingEncoder
	Wasm    WasmEncoder
}

type CustomQuerier func(ctx sdk.Context, request json.RawMessage) ([]byte, error)

type QueryPlugins struct {
	Bank    func(ctx sdk.Context, request *wasmTypes.BankQuery) ([]byte, error)
	Custom  CustomQuerier
	Staking func(ctx sdk.Context, request *wasmTypes.StakingQuery) ([]byte, error)
	Wasm    func(ctx sdk.Context, request *wasmTypes.WasmQuery) ([]byte, error)
}

type Keeper struct{}

// NewKeeper creates a new contract Keeper instance
// If customEncoders is non-nil, we can use this to override some of the message handler, especially custom
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper,
	router sdk.Router, homeDir string, wasmConfig types.WasmConfig, supportedFeatures string,
	customEncoders *MessageEncoders, customPlugins *QueryPlugins) Keeper {
	return Keeper{}
}
