package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "distribution"

	// StoreKey is the store key string for distribution
	StoreKey = ModuleName

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// QuerierRoute is the querier route for distribution
	QuerierRoute = ModuleName

	// ShortUseByCli added for okexchaincli
	ShortUseByCli = "distr"
)

// Keys for distribution store
// Items are stored with the following key: values
//
// - 0x01: sdk.ConsAddress
//
// - 0x03<accAddr_Bytes>: sdk.AccAddress
//
// - 0x07<valAddr_Bytes>: ValidatorCurrentRewards
var (
	FeePoolKey                           = []byte{0x00} // key for global distribution state
	ProposerKey                          = []byte{0x01} // key for the proposer operator address
	DelegatorWithdrawAddrPrefix          = []byte{0x03} // key for delegator withdraw address
	ValidatorAccumulatedCommissionPrefix = []byte{0x07} // key for accumulated validator commission
)

// GetDelegatorWithdrawInfoAddress returns an address from a delegator's withdraw info key
func GetDelegatorWithdrawInfoAddress(key []byte) (delAddr sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

//GetValidatorAccumulatedCommissionAddress returns the address from a validator's accumulated commission key
func GetValidatorAccumulatedCommissionAddress(key []byte) (valAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}

// GetDelegatorWithdrawAddrKey returns the key for a delegator's withdraw addr
func GetDelegatorWithdrawAddrKey(delAddr sdk.AccAddress) []byte {
	return append(DelegatorWithdrawAddrPrefix, delAddr.Bytes()...)
}

// GetValidatorAccumulatedCommissionKey returns the key for a validator's current commission
func GetValidatorAccumulatedCommissionKey(v sdk.ValAddress) []byte {
	return append(ValidatorAccumulatedCommissionPrefix, v.Bytes()...)
}
