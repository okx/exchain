package params

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/system"
	"github.com/okx/okbchain/x/params/types"
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
	cdc.RegisterConcrete(types.ParameterChangeProposal{}, system.Chain+"/params/ParameterChangeProposal", nil)
	cdc.RegisterConcrete(types.UpgradeProposal{}, system.Chain+"/params/UpgradeProposal", nil)
	cdc.RegisterConcrete(types.UpgradeInfo{}, system.Chain+"/params/UpgradeInfo", nil)
}
