package protocol

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap"
	"github.com/okex/okexchain/x/debug"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/farm"
	"github.com/okex/okexchain/x/staking"

	distr "github.com/okex/okexchain/x/distribution"
	"github.com/okex/okexchain/x/gov"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/params"

	"github.com/okex/okexchain/x/token"
)

// store keys used in all modules
var (
	kvStoreKeysMap = sdk.NewKVStoreKeys(
		baseapp.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		mint.StoreKey,
		slashing.StoreKey,
		distr.StoreKey,
		gov.StoreKey,
		params.StoreKey,
		token.StoreKey, token.KeyMint, token.KeyLock,
		order.OrderStoreKey,
		dex.StoreKey, dex.TokenPairStoreKey,
		debug.StoreKey,
		ammswap.StoreKey,
		farm.StoreKey,
	)

	transientStoreKeysMap = sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)
)

// GetMainStoreKey gets the main store key
func GetMainStoreKey() *sdk.KVStoreKey {
	return kvStoreKeysMap[baseapp.MainStoreKey]
}

// GetKVStoreKeysMap gets the map of all kv store keys
func GetKVStoreKeysMap() map[string]*sdk.KVStoreKey {
	return kvStoreKeysMap
}

// GetTransientStoreKeysMap gets the map of all transient store keys
func GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey {
	return transientStoreKeysMap
}
