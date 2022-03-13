package evm

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
)

// nolint
const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

// nolint
var (
	NewKeeper            = keeper.NewKeeper
	TxDecoder            = types.TxDecoder
	NewSimulateKeeper    = keeper.NewSimulateKeeper
  SetEvmParamsNeedUpdate = types.SetEvmParamsNeedUpdate
	NewLogProcessEvmHook = keeper.NewLogProcessEvmHook
)

//nolint
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
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
