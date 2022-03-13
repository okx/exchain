package types

const (
	// ModuleName string name of module
	ModuleName = "erc20"

	// StoreKey key for ethereum storage data, account code (StateDB) or block
	// related data for Web3.
	// The erc20 module should use a prefix store.
	StoreKey = ModuleName

	// RouterKey uses module name for routing
	RouterKey = ModuleName
)

// KVStore key prefixes
var (
	KeyPrefixContractToDenom         = []byte{0x01}
	KeyPrefixDenomToExternalContract = []byte{0x02}
	KeyPrefixDenoToAutoContract      = []byte{0x03}
)

// ContractToDenomKey defines the store key for contract to denom reverse index
func ContractToDenomKey(contract []byte) []byte {
	return append(KeyPrefixContractToDenom, contract...)
}

// DenomToExternalContractKey defines the store key for denom to external contract mapping
func DenomToExternalContractKey(denom string) []byte {
	return append(KeyPrefixDenomToExternalContract, denom...)
}

// DenomToAutoContractKey defines the store key for denom to auto contract mapping
func DenomToAutoContractKey(denom string) []byte {
	return append(KeyPrefixDenoToAutoContract, denom...)
}
