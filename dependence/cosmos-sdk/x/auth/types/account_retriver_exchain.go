package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (ar AccountRetriever) GetAccountNonce(address string) uint64 {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return 0
	}
	acc, err := ar.GetAccount(addr)
	if err != nil {
		return 0
	}
	return acc.GetSequence()
}
