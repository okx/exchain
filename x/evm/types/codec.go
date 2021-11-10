package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
)

// ModuleCdc defines the evm module's codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for the
// evm module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgEthereumTx{}, "ethermint/MsgEthereumTx", nil)
	cdc.RegisterConcrete(MsgEthermint{}, "ethermint/MsgEthermint", nil)
	cdc.RegisterConcrete(TxData{}, "ethermint/TxData", nil)
	cdc.RegisterConcrete(ChainConfig{}, "ethermint/ChainConfig", nil)
	cdc.RegisterConcrete(ManageContractDeploymentWhitelistProposal{}, "okexchain/evm/ManageContractDeploymentWhitelistProposal", nil)
	cdc.RegisterConcrete(ManageContractBlockedListProposal{}, "okexchain/evm/ManageContractBlockedListProposal", nil)
	cdc.RegisterConcrete(ManageContractMethodBlockedListProposal{}, "okexchain/evm/ManageContractMethodBlockedListProposal", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
