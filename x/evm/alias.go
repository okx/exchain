package evm

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/evm/keeper"
	"github.com/okx/okbchain/x/evm/types"
)

// nolint
const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

// nolint
var (
	NewKeeper            = keeper.NewKeeper
	TxDecoder            = types.TxDecoder
	NewSimulateKeeper    = keeper.NewSimulateKeeper
	NewLogProcessEvmHook = keeper.NewLogProcessEvmHook
	NewMultiEvmHooks     = keeper.NewMultiEvmHooks
)

//nolint
type (
	Keeper        = keeper.Keeper
	GenesisState  = types.GenesisState
	EvmLogHandler = types.EvmLogHandler
)

func WithMoreDeocder(cdc *codec.Codec, cc sdk.TxDecoder) sdk.TxDecoder {
	return func(txBytes []byte, height ...int64) (sdk.Tx, error) {
		ret, err := cc(txBytes, height...)
		if nil == err {
			return ret, nil
		}
		return ret, nil
	}
}
