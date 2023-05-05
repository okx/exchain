package keeper

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/wasm/types"
)

var _ types.ContractOpsKeeper = PermissionedKeeper{}

// decoratedKeeper contains a subset of the wasm keeper that are already or can be guarded by an authorization policy in the future
type decoratedKeeper interface {
	create(ctx sdk.Context, creator sdk.WasmAddress, wasmCode []byte, instantiateAccess *types.AccessConfig, authZ AuthorizationPolicy) (codeID uint64, err error)
	instantiate(ctx sdk.Context, codeID uint64, creator, admin sdk.WasmAddress, initMsg []byte, label string, deposit sdk.Coins, authZ AuthorizationPolicy) (sdk.WasmAddress, []byte, error)
	migrate(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, newCodeID uint64, msg []byte, authZ AuthorizationPolicy) ([]byte, error)
	setContractAdmin(ctx sdk.Context, contractAddress, caller, newAdmin sdk.WasmAddress, authZ AuthorizationPolicy) error
	pinCode(ctx sdk.Context, codeID uint64) error
	unpinCode(ctx sdk.Context, codeID uint64) error
	execute(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, msg []byte, coins sdk.Coins) ([]byte, error)
	Sudo(ctx sdk.Context, contractAddress sdk.WasmAddress, msg []byte) ([]byte, error)
	setContractInfoExtension(ctx sdk.Context, contract sdk.WasmAddress, extra types.ContractInfoExtension) error
	setAccessConfig(ctx sdk.Context, codeID uint64, config types.AccessConfig) error
	updateUploadAccessConfig(ctx sdk.Context, config types.AccessConfig)
	updateContractMethodBlockedList(ctx sdk.Context, blockedMethods *types.ContractMethods, isDelete bool) error

	GetParams(ctx sdk.Context) types.Params
	newQueryHandler(ctx sdk.Context, contractAddress sdk.WasmAddress) QueryHandler
	runtimeGasForContract(ctx sdk.Context) uint64
	InvokeExtraProposal(ctx sdk.Context, action string, extra string) error
}

type PermissionedKeeper struct {
	authZPolicy AuthorizationPolicy
	nested      decoratedKeeper
}

func NewPermissionedKeeper(nested decoratedKeeper, authZPolicy AuthorizationPolicy) *PermissionedKeeper {
	return &PermissionedKeeper{authZPolicy: authZPolicy, nested: nested}
}

func NewGovPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested, GovAuthorizationPolicy{})
}

func NewDefaultPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested, DefaultAuthorizationPolicy{})
}

func (p PermissionedKeeper) Create(ctx sdk.Context, creator sdk.WasmAddress, wasmCode []byte, instantiateAccess *types.AccessConfig) (codeID uint64, err error) {
	return p.nested.create(ctx, creator, wasmCode, instantiateAccess, p.authZPolicy)
}

func (p PermissionedKeeper) Instantiate(ctx sdk.Context, codeID uint64, creator, admin sdk.WasmAddress, initMsg []byte, label string, deposit sdk.Coins) (sdk.WasmAddress, []byte, error) {
	return p.nested.instantiate(ctx, codeID, creator, admin, initMsg, label, deposit, p.authZPolicy)
}

func (p PermissionedKeeper) Execute(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, msg []byte, coins sdk.Coins) ([]byte, error) {
	return p.nested.execute(ctx, contractAddress, caller, msg, coins)
}

func (p PermissionedKeeper) Migrate(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, newCodeID uint64, msg []byte) ([]byte, error) {
	return p.nested.migrate(ctx, contractAddress, caller, newCodeID, msg, p.authZPolicy)
}

func (p PermissionedKeeper) Sudo(ctx sdk.Context, contractAddress sdk.WasmAddress, msg []byte) ([]byte, error) {
	return p.nested.Sudo(ctx, contractAddress, msg)
}

func (p PermissionedKeeper) UpdateContractAdmin(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, newAdmin sdk.WasmAddress) error {
	return p.nested.setContractAdmin(ctx, contractAddress, caller, newAdmin, p.authZPolicy)
}

func (p PermissionedKeeper) ClearContractAdmin(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress) error {
	return p.nested.setContractAdmin(ctx, contractAddress, caller, nil, p.authZPolicy)
}

func (p PermissionedKeeper) PinCode(ctx sdk.Context, codeID uint64) error {
	return p.nested.pinCode(ctx, codeID)
}

func (p PermissionedKeeper) UnpinCode(ctx sdk.Context, codeID uint64) error {
	return p.nested.unpinCode(ctx, codeID)
}

// SetExtraContractAttributes updates the extra attributes that can be stored with the contract info
func (p PermissionedKeeper) SetContractInfoExtension(ctx sdk.Context, contract sdk.WasmAddress, extra types.ContractInfoExtension) error {
	return p.nested.setContractInfoExtension(ctx, contract, extra)
}

// SetAccessConfig updates the access config of a code id.
func (p PermissionedKeeper) SetAccessConfig(ctx sdk.Context, codeID uint64, config types.AccessConfig) error {
	return p.nested.setAccessConfig(ctx, codeID, config)
}

func (p PermissionedKeeper) UpdateUploadAccessConfig(ctx sdk.Context, config types.AccessConfig) {
	p.nested.updateUploadAccessConfig(ctx, config)
}

func (p PermissionedKeeper) UpdateContractMethodBlockedList(ctx sdk.Context, blockedMethods *types.ContractMethods, isDelete bool) error {
	return p.nested.updateContractMethodBlockedList(ctx, blockedMethods, isDelete)
}

func (p PermissionedKeeper) GetParams(ctx sdk.Context) types.Params {
	return p.nested.GetParams(ctx)
}

func (p PermissionedKeeper) NewQueryHandler(ctx sdk.Context, contractAddress sdk.WasmAddress) wasmvmtypes.Querier {
	return p.nested.newQueryHandler(ctx, contractAddress)
}

func (p PermissionedKeeper) RuntimeGasForContract(ctx sdk.Context) uint64 {
	return p.nested.runtimeGasForContract(ctx)
}

func (p PermissionedKeeper) InvokeExtraProposal(ctx sdk.Context, action string, extra string) error {
	return p.nested.InvokeExtraProposal(ctx, action, extra)
}
