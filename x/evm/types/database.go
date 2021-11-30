package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
)

type Watcher interface {
	SaveAccount(account auth.Account, isDirectly bool)
	SaveState(addr ethcmn.Address, key, value []byte)
	Enabled() bool
	SaveContractBlockedListItem(addr sdk.AccAddress)
	SaveContractDeploymentWhitelistItem(addr sdk.AccAddress)
	DeleteContractBlockedList(addr sdk.AccAddress)
	DeleteContractDeploymentWhitelist(addr sdk.AccAddress)
}

type DefaultPrefixDb struct {
}

func (d DefaultPrefixDb) NewStore(parent types.KVStore, Prefix []byte) StoreProxy {
	return prefix.NewStore(parent, Prefix)
}

type StoreProxy interface {
	Set(key, value []byte)
	Get(key []byte) []byte
	Delete(key []byte)
	Has(key []byte) bool
}

type DbAdapter interface {
	NewStore(parent types.KVStore, prefix []byte) StoreProxy
}
