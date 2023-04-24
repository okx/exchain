package types

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
)

// ViewKeeper provides read only operations
type ViewKeeper interface {
	GetContractHistory(ctx sdk.Context, contractAddr sdk.WasmAddress) []ContractCodeHistoryEntry
	QuerySmart(ctx sdk.Context, contractAddr sdk.WasmAddress, req []byte) ([]byte, error)
	QueryRaw(ctx sdk.Context, contractAddress sdk.WasmAddress, key []byte) []byte
	HasContractInfo(ctx sdk.Context, contractAddress sdk.WasmAddress) bool
	GetContractInfo(ctx sdk.Context, contractAddress sdk.WasmAddress) *ContractInfo
	IterateContractInfo(ctx sdk.Context, cb func(sdk.WasmAddress, ContractInfo) bool)
	IterateContractsByCode(ctx sdk.Context, codeID uint64, cb func(address sdk.WasmAddress) bool)
	IterateContractState(ctx sdk.Context, contractAddress sdk.WasmAddress, cb func(key, value []byte) bool)
	GetCodeInfo(ctx sdk.Context, codeID uint64) *CodeInfo
	IterateCodeInfos(ctx sdk.Context, cb func(uint64, CodeInfo) bool)
	GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, error)
	IsPinnedCode(ctx sdk.Context, codeID uint64) bool
	GetContractMethodBlockedList(ctx sdk.Context, contractAddr string) *ContractMethods
	GetParams(ctx sdk.Context) Params
	GetGasFactor(ctx sdk.Context) uint64
}

// ContractOpsKeeper contains mutable operations on a contract.
type ContractOpsKeeper interface {
	// Create uploads and compiles a WASM contract, returning a short identifier for the contract
	Create(ctx sdk.Context, creator sdk.WasmAddress, wasmCode []byte, instantiateAccess *AccessConfig) (codeID uint64, err error)

	// Instantiate creates an instance of a WASM contract
	Instantiate(ctx sdk.Context, codeID uint64, creator, admin sdk.WasmAddress, initMsg []byte, label string, deposit sdk.Coins) (sdk.WasmAddress, []byte, error)

	// Execute executes the contract instance
	Execute(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, msg []byte, coins sdk.Coins) ([]byte, error)

	// Migrate allows to upgrade a contract to a new code with data migration.
	Migrate(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, newCodeID uint64, msg []byte) ([]byte, error)

	// Sudo allows to call privileged entry point of a contract.
	Sudo(ctx sdk.Context, contractAddress sdk.WasmAddress, msg []byte) ([]byte, error)

	// UpdateContractAdmin sets the admin value on the ContractInfo. It must be a valid address (use ClearContractAdmin to remove it)
	UpdateContractAdmin(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress, newAdmin sdk.WasmAddress) error

	// ClearContractAdmin sets the admin value on the ContractInfo to nil, to disable further migrations/ updates.
	ClearContractAdmin(ctx sdk.Context, contractAddress sdk.WasmAddress, caller sdk.WasmAddress) error

	// PinCode pins the wasm contract in wasmvm cache
	PinCode(ctx sdk.Context, codeID uint64) error

	// UnpinCode removes the wasm contract from wasmvm cache
	UnpinCode(ctx sdk.Context, codeID uint64) error

	// SetContractInfoExtension updates the extension point data that is stored with the contract info
	SetContractInfoExtension(ctx sdk.Context, contract sdk.WasmAddress, extra ContractInfoExtension) error

	// SetAccessConfig updates the access config of a code id.
	SetAccessConfig(ctx sdk.Context, codeID uint64, config AccessConfig) error

	// UpdateUploadAccessConfig updates the access config of uploading code.
	UpdateUploadAccessConfig(ctx sdk.Context, config AccessConfig)

	// UpdateContractMethodBlockedList updates the blacklist of contract methods.
	UpdateContractMethodBlockedList(ctx sdk.Context, methods *ContractMethods, isDelete bool) error

	// GetParams get params from paramsubspace.
	GetParams(ctx sdk.Context) Params

	// InvokeExtraProposal invoke extra proposal
	InvokeExtraProposal(ctx sdk.Context, action string, extra string) error
}

// IBCContractKeeper IBC lifecycle event handler
type IBCContractKeeper interface {
	OnOpenChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCChannelOpenMsg,
	) (string, error)
	OnConnectChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCChannelConnectMsg,
	) error
	OnCloseChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCChannelCloseMsg,
	) error
	OnRecvPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCPacketReceiveMsg,
	) ([]byte, error)
	OnAckPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		acknowledgement wasmvmtypes.IBCPacketAckMsg,
	) error
	OnTimeoutPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg wasmvmtypes.IBCPacketTimeoutMsg,
	) error
	// ClaimCapability allows the transfer module to claim a capability
	// that IBC module passes to it
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
	AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool
}
