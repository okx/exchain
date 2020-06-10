package protocol

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okchain/x/debug"
	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/wasm"

	//distr "github.com/okex/okchain/x/distribution"
	distr "github.com/okex/okchain/x/distribution"
	"github.com/okex/okchain/x/gov"
	"github.com/okex/okchain/x/order"
	"github.com/okex/okchain/x/params"

	//"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/token"
	"github.com/okex/okchain/x/upgrade"
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
		wasm.StoreKey,
		upgrade.StoreKey,
		dex.StoreKey, dex.TokenPairStoreKey,
		debug.StoreKey,
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
