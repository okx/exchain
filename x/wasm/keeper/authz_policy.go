package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/x/wasm/types"
)

type AuthorizationPolicy interface {
	CanCreateCode(c types.AccessConfig, creator sdk.WasmAddress) bool
	CanInstantiateContract(c types.AccessConfig, actor sdk.WasmAddress) bool
	CanModifyContract(admin, actor sdk.WasmAddress) bool
}

type DefaultAuthorizationPolicy struct{}

func (p DefaultAuthorizationPolicy) CanCreateCode(config types.AccessConfig, actor sdk.WasmAddress) bool {
	return config.Allowed(actor)
}

func (p DefaultAuthorizationPolicy) CanInstantiateContract(config types.AccessConfig, actor sdk.WasmAddress) bool {
	return config.Allowed(actor)
}

func (p DefaultAuthorizationPolicy) CanModifyContract(admin, actor sdk.WasmAddress) bool {
	return admin != nil && admin.Equals(actor)
}

type GovAuthorizationPolicy struct{}

func (p GovAuthorizationPolicy) CanCreateCode(types.AccessConfig, sdk.WasmAddress) bool {
	return true
}

func (p GovAuthorizationPolicy) CanInstantiateContract(types.AccessConfig, sdk.WasmAddress) bool {
	return true
}

func (p GovAuthorizationPolicy) CanModifyContract(sdk.WasmAddress, sdk.WasmAddress) bool {
	return true
}
