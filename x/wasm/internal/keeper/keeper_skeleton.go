package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/okex/okchain/x/wasm/internal/types"
)

type MessageEncoders struct {}

type QueryPlugins struct {}

type Keeper struct{}

// NewKeeper creates a new contract Keeper instance
// If customEncoders is non-nil, we can use this to override some of the message handler, especially custom
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper,
	router sdk.Router, homeDir string, wasmConfig types.WasmConfig, supportedFeatures string,
	customEncoders *MessageEncoders, customPlugins *QueryPlugins) Keeper {
	return Keeper{}
}
