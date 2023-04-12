package params

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdkparams "github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/x/params/types"
)

// ModuleCdc is the codec of module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers all necessary param module types with a given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(types.ParameterChangeProposal{}, "okexchain/params/ParameterChangeProposal", nil)
	cdc.RegisterConcrete(sdkparams.ParameterChangeProposal{}, "cosmos-sdk/params/ParameterChangeProposal", nil)
	cdc.RegisterConcrete(types.UpgradeProposal{}, "okexchain/params/UpgradeProposal", nil)
	cdc.RegisterConcrete(types.UpgradeInfo{}, "okexchain/params/UpgradeInfo", nil)
}
