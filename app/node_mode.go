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
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mempool"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
)

func SetNodeConfig(ctx *server.Context) {
	nodeMode := viper.GetString(types.FlagNodeMode)
	switch types.NodeMode(nodeMode) {
	case types.RpcNode:
		setRpcConfig(ctx)
	case types.ValidatorNode:
		setValidatorConfig(ctx)
	case types.ArchiveNode:
		setArchiveConfig(ctx)
	case "":
		ctx.Logger.Info("The node mode is not set for this node")
	default:
		ctx.Logger.Error(
			fmt.Sprintf("Wrong value (%s) is set for %s, the correct value should be one of %s, %s, and %s",
				nodeMode, types.FlagNodeMode, types.RpcNode, types.ValidatorNode, types.ArchiveNode))
	}
}

func setRpcConfig(ctx *server.Context) {
	viper.SetDefault(abcitypes.FlagDisableCheckTxMutex, true)
	viper.SetDefault(abcitypes.FlagDisableQueryMutex, true)
	viper.SetDefault(evmtypes.FlagEnableBloomFilter, true)
	viper.SetDefault(watcher.FlagFastQueryLru, 10000)
	viper.SetDefault(watcher.FlagFastQuery, true)
	viper.SetDefault(iavl.FlagIavlEnableAsyncCommit, true)
	viper.SetDefault(flags.FlagMaxOpenConnections, 20000)
	viper.SetDefault(mempool.FlagEnablePendingPool, true)
	viper.SetDefault(server.FlagCORS, "*")
	ctx.Logger.Info(fmt.Sprintf(
		"Set --%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v by rpc node mode",
		abcitypes.FlagDisableCheckTxMutex, true, abcitypes.FlagDisableQueryMutex, true,
		evmtypes.FlagEnableBloomFilter, true, watcher.FlagFastQueryLru, 10000,
		watcher.FlagFastQuery, true, iavl.FlagIavlEnableAsyncCommit, true,
		flags.FlagMaxOpenConnections, 20000, mempool.FlagEnablePendingPool, true,
		server.FlagCORS, "*"))
}

func setValidatorConfig(ctx *server.Context) {
	viper.SetDefault(abcitypes.FlagDisableCheckTxMutex, true)
	viper.SetDefault(abcitypes.FlagDisableQueryMutex, true)
	viper.SetDefault(appconfig.FlagEnableDynamicGp, false)
	viper.SetDefault(iavl.FlagIavlEnableAsyncCommit, true)
	viper.SetDefault(store.FlagIavlCacheSize, 10000000)
	viper.SetDefault(server.FlagPruning, "everything")
	ctx.Logger.Info(fmt.Sprintf("Set --%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v by validator node mode",
		abcitypes.FlagDisableCheckTxMutex, true, abcitypes.FlagDisableQueryMutex, true,
		appconfig.FlagEnableDynamicGp, false, iavl.FlagIavlEnableAsyncCommit, true,
		store.FlagIavlCacheSize, 10000000, server.FlagPruning, "everything"))
}

func setArchiveConfig(ctx *server.Context) {
	viper.SetDefault(server.FlagPruning, "nothing")
	viper.SetDefault(abcitypes.FlagDisableCheckTxMutex, true)
	viper.SetDefault(abcitypes.FlagDisableQueryMutex, true)
	viper.SetDefault(evmtypes.FlagEnableBloomFilter, true)
	viper.SetDefault(watcher.FlagFastQueryLru, 10000)
	viper.SetDefault(watcher.FlagFastQuery, true)
	viper.SetDefault(iavl.FlagIavlEnableAsyncCommit, true)
	viper.SetDefault(flags.FlagMaxOpenConnections, 20000)
	viper.SetDefault(server.FlagCORS, "*")
	ctx.Logger.Info(fmt.Sprintf(
		"Set --%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v\n--%s=%v by rpc archive mode",
		server.FlagPruning, "nothing", abcitypes.FlagDisableCheckTxMutex, true,
		abcitypes.FlagDisableQueryMutex, true, evmtypes.FlagEnableBloomFilter, true,
		watcher.FlagFastQueryLru, 10000, watcher.FlagFastQuery, true,
		iavl.FlagIavlEnableAsyncCommit, true, flags.FlagMaxOpenConnections, 20000,
		server.FlagCORS, "*"))
}

func logStartingFlags(logger log.Logger) {
	flagMap := map[string]interface{}{
		server.FlagPruning:                viper.GetString(server.FlagPruning),
		abcitypes.FlagDisableCheckTxMutex: viper.GetBool(abcitypes.FlagDisableCheckTxMutex),
		abcitypes.FlagDisableQueryMutex:   viper.GetBool(abcitypes.FlagDisableQueryMutex),
		evmtypes.FlagEnableBloomFilter:    viper.GetBool(evmtypes.FlagEnableBloomFilter),
		watcher.FlagFastQueryLru:          viper.GetInt(watcher.FlagFastQueryLru),
		watcher.FlagFastQuery:             viper.GetBool(watcher.FlagFastQuery),
		iavl.FlagIavlEnableAsyncCommit:    viper.GetBool(iavl.FlagIavlEnableAsyncCommit),
		flags.FlagMaxOpenConnections:      viper.GetInt(flags.FlagMaxOpenConnections),
		server.FlagCORS:                   viper.GetString(server.FlagCORS),
		appconfig.FlagEnableDynamicGp:     viper.GetBool(appconfig.FlagEnableDynamicGp),
		store.FlagIavlCacheSize:           viper.GetInt(store.FlagIavlCacheSize),
		mempool.FlagEnablePendingPool:     viper.GetBool(mempool.FlagEnablePendingPool),
		appconfig.FlagMaxGasUsedPerBlock:  viper.GetInt64(appconfig.FlagMaxGasUsedPerBlock),
	}
	msg := "starting flags:"
	for k, v := range flagMap {
		msg += fmt.Sprintf("\n%s=%v", k, v)
	}

	logger.Info(msg)
}
