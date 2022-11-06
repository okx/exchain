package types

const (
	// ModuleName is the name of the contract module
	ModuleName = "vmbridge"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the wasm module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the wasm module
	RouterKey = ModuleName
)
