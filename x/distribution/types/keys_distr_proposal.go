package types

import (
	"encoding/binary"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

var (
	ValidatorOutstandingRewardsPrefix = []byte{0x02} // key for outstanding rewards
	DelegatorStartingInfoPrefix       = []byte{0x04} // key for delegator starting info
	ValidatorHistoricalRewardsPrefix  = []byte{0x05} // key for historical validators rewards / stake
	ValidatorCurrentRewardsPrefix     = []byte{0x06} // key for current validator rewards
)

// gets an address from a validator's outstanding rewards key
func GetValidatorOutstandingRewardsAddress(key []byte) (valAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}

// gets the addresses from a delegator starting info key
func GetDelegatorStartingInfoAddresses(key []byte) (valAddr sdk.ValAddress, delAddr sdk.AccAddress) {
	addr := key[1 : 1+sdk.AddrLen]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	valAddr = sdk.ValAddress(addr)
	addr = key[1+sdk.AddrLen:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	delAddr = sdk.AccAddress(addr)
	return
}

// gets the address & period from a validator's historical rewards key
func GetValidatorHistoricalRewardsAddressPeriod(key []byte) (valAddr sdk.ValAddress, period uint64) {
	addr := key[1 : 1+sdk.AddrLen]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	valAddr = sdk.ValAddress(addr)
	b := key[1+sdk.AddrLen:]
	if len(b) != 8 {
		panic("unexpected key length")
	}
	period = binary.LittleEndian.Uint64(b)
	return
}

// gets the address from a validator's current rewards key
func GetValidatorCurrentRewardsAddress(key []byte) (valAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}

// gets the outstanding rewards key for a validator
func GetValidatorOutstandingRewardsKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorOutstandingRewardsPrefix, valAddr.Bytes()...)
}

// gets the key for a delegator's starting info
func GetDelegatorStartingInfoKey(v sdk.ValAddress, d sdk.AccAddress) []byte {
	return append(append(DelegatorStartingInfoPrefix, v.Bytes()...), d.Bytes()...)
}

// gets the prefix key for a validator's historical rewards
func GetValidatorHistoricalRewardsPrefix(v sdk.ValAddress) []byte {
	return append(ValidatorHistoricalRewardsPrefix, v.Bytes()...)
}

// gets the key for a validator's historical rewards
func GetValidatorHistoricalRewardsKey(v sdk.ValAddress, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(ValidatorHistoricalRewardsPrefix, v.Bytes()...), b...)
}

// gets the key for a validator's current rewards
func GetValidatorCurrentRewardsKey(v sdk.ValAddress) []byte {
	return append(ValidatorCurrentRewardsPrefix, v.Bytes()...)
}
