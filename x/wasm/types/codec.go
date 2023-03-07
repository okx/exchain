package types

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	cryptocodec "github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	txmsg "github.com/okx/okbchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/msgservice"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

// RegisterLegacyAminoCodec registers the account types and interface
func RegisterLegacyAminoCodec(cdc *codec.Codec) { //nolint:staticcheck
	cdc.RegisterConcrete(&MsgStoreCode{}, "wasm/MsgStoreCode", nil)
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "wasm/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(&MsgExecuteContract{}, "wasm/MsgExecuteContract", nil)
	cdc.RegisterConcrete(&MsgMigrateContract{}, "wasm/MsgMigrateContract", nil)
	cdc.RegisterConcrete(&MsgUpdateAdmin{}, "wasm/MsgUpdateAdmin", nil)
	cdc.RegisterConcrete(&MsgClearAdmin{}, "wasm/MsgClearAdmin", nil)

	cdc.RegisterConcrete(&PinCodesProposal{}, "wasm/PinCodesProposal", nil)
	cdc.RegisterConcrete(&UnpinCodesProposal{}, "wasm/UnpinCodesProposal", nil)
	cdc.RegisterConcrete(&StoreCodeProposal{}, "wasm/StoreCodeProposal", nil)
	cdc.RegisterConcrete(&InstantiateContractProposal{}, "wasm/InstantiateContractProposal", nil)
	cdc.RegisterConcrete(&MigrateContractProposal{}, "wasm/MigrateContractProposal", nil)
	cdc.RegisterConcrete(&SudoContractProposal{}, "wasm/SudoContractProposal", nil)
	cdc.RegisterConcrete(&ExecuteContractProposal{}, "wasm/ExecuteContractProposal", nil)
	cdc.RegisterConcrete(&UpdateAdminProposal{}, "wasm/UpdateAdminProposal", nil)
	cdc.RegisterConcrete(&ClearAdminProposal{}, "wasm/ClearAdminProposal", nil)
	cdc.RegisterConcrete(&UpdateInstantiateConfigProposal{}, "wasm/UpdateInstantiateConfigProposal", nil)
	cdc.RegisterConcrete(&UpdateDeploymentWhitelistProposal{}, "wasm/UpdateDeploymentWhitelistProposal", nil)
	cdc.RegisterConcrete(&UpdateWASMContractMethodBlockedListProposal{}, "wasm/UpdateWASMContractMethodBlockedListProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgStoreCode{},
		&MsgInstantiateContract{},
		&MsgExecuteContract{},
		&MsgMigrateContract{},
		&MsgUpdateAdmin{},
		&MsgClearAdmin{},
		&MsgIBCCloseChannel{},
		&MsgIBCSend{},
	)
	registry.RegisterImplementations((*txmsg.Msg)(nil),
		&MsgStoreCode{},
		&MsgInstantiateContract{},
		&MsgExecuteContract{},
		&MsgMigrateContract{},
		&MsgUpdateAdmin{},
		&MsgClearAdmin{},
		&MsgIBCCloseChannel{},
		&MsgIBCSend{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&StoreCodeProposal{},
		&InstantiateContractProposal{},
		&MigrateContractProposal{},
		&SudoContractProposal{},
		&ExecuteContractProposal{},
		&UpdateAdminProposal{},
		&ClearAdminProposal{},
		&PinCodesProposal{},
		&UnpinCodesProposal{},
		&UpdateInstantiateConfigProposal{},
		&UpdateDeploymentWhitelistProposal{},
		&UpdateWASMContractMethodBlockedListProposal{},
	)

	registry.RegisterInterface("ContractInfoExtension", (*ContractInfoExtension)(nil))

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.New()
)

func init() {
	RegisterLegacyAminoCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
}
