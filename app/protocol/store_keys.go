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
		auth.StoreKey,

		// for staking/distr rollback to cosmos-sdk
		//staking.StoreKey, staking.DelegatorPoolKey, staking.RedelegationActonKey, staking.RedelegationKeyM, staking.UnbondingKey,
		staking.StoreKey,

		supply.StoreKey,
		mint.StoreKey,
		slashing.StoreKey,

		// for staking/distr rollback to cosmos-sdk
		//distr.StoreKey, distr.ValidatorsSnapshotKey, distr.DelegationSnapshotKey,
		distr.StoreKey,

		gov.StoreKey,
		params.StoreKey,
		token.StoreKey, token.KeyMint, token.KeyLock, token.KeyFreeze,
		order.OrderStoreKey,
		upgrade.StoreKey,
		// for test
		debug.StoreKey,
		baseapp.MainStoreKey,
		dex.StoreKey, dex.TokenPairStoreKey,
	)

	transientStoreKeysMap = sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)
)

func GetMainStoreKey() *sdk.KVStoreKey {
	return kvStoreKeysMap[baseapp.MainStoreKey]
}

func GetKVStoreKeysMap() map[string]*sdk.KVStoreKey {
	return kvStoreKeysMap
}

func GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey {
	return transientStoreKeysMap
}
