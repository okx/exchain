package app

import (
	"fmt"

	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	store "github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	"github.com/okex/exchain/libs/iavl"
	abcitypes "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/mempool"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

func SetNodeConfig(context *server.Context) {
	nodeMode := viper.GetString(types.FlagNodeMode)
	switch types.NodeMode(nodeMode) {
	case types.RpcNode:
		setRpcConfig()
	case types.ValidatorNode:
		setValidatorConfig()
	case types.ArchiveNode:
		setArchiveConfig()
	case "":
	default:
		context.Logger.Error(
			fmt.Sprintf("Wrong value (%s) is set for %s, the correct value should be one of %s, %s, and %s",
				nodeMode, types.FlagNodeMode, types.RpcNode, types.ValidatorNode, types.ArchiveNode))
	}
}

func setRpcConfig() {
	viper.SetDefault(abcitypes.FlagDisableCheckTxMutex, true)
	viper.SetDefault(abcitypes.FlagDisableQueryMutex, true)
	viper.SetDefault(evmtypes.FlagEnableBloomFilter, true)
	viper.SetDefault(watcher.FlagFastQueryLru, 10000)
	viper.SetDefault(watcher.FlagFastQuery, true)
	viper.SetDefault(iavl.FlagIavlEnableAsyncCommit, true)
	viper.SetDefault(flags.FlagMaxOpenConnections, 20000)
	viper.SetDefault(mempool.FlagEnablePendingPool, true)
	viper.SetDefault(server.FlagCORS, "*")
}

func setValidatorConfig() {
	viper.SetDefault(abcitypes.FlagDisableCheckTxMutex, true)
	viper.SetDefault(abcitypes.FlagDisableQueryMutex, true)
	viper.SetDefault(appconfig.FlagEnableDynamicGp, false)
	viper.SetDefault(iavl.FlagIavlEnableAsyncCommit, true)
	viper.SetDefault(store.FlagIavlCacheSize, 10000000)
	viper.SetDefault(server.FlagPruning, "everything")
}

func setArchiveConfig() {
	viper.SetDefault(server.FlagPruning, "nothing")
	viper.SetDefault(abcitypes.FlagDisableCheckTxMutex, true)
	viper.SetDefault(abcitypes.FlagDisableQueryMutex, true)
	viper.SetDefault(evmtypes.FlagEnableBloomFilter, true)
	viper.SetDefault(watcher.FlagFastQueryLru, 10000)
	viper.SetDefault(watcher.FlagFastQuery, true)
	viper.SetDefault(iavl.FlagIavlEnableAsyncCommit, true)
	viper.SetDefault(flags.FlagMaxOpenConnections, 20000)
	viper.SetDefault(server.FlagCORS, "*")
}
