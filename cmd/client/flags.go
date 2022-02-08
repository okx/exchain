package client

import (
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/rpc"
	"github.com/okex/exchain/app/rpc/namespaces/eth"
	"github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/app/utils/sanity"
	"github.com/okex/exchain/libs/tendermint/consensus"
	"github.com/okex/exchain/libs/tendermint/libs/automation"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/common/analyzer"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/token"
	"github.com/spf13/cobra"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().Int(watcher.FlagFastQueryLru, 1000, "Set the size of LRU cache under fast-query mode")
	cmd.Flags().Bool(watcher.FlagCheckWd, false, "Enable check watchDB in log")
	cmd.Flags().Bool(rpc.FlagPersonalAPI, true, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(rpc.FlagDebugAPI, false, "Enable the debug_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(evmtypes.FlagEnableBloomFilter, false, "Enable bloom filter for event logs")
	cmd.Flags().Int64(filters.FlagGetLogsHeightSpan, 2000, "config the block height span for get logs")
	// register application rpc to nacos
	cmd.Flags().String(rpc.FlagRestApplicationName, "", "rest application name in  nacos")
	cmd.Flags().MarkHidden(rpc.FlagRestApplicationName)
	cmd.Flags().String(rpc.FlagRestNacosUrls, "", "nacos server urls for discovery service of rest api")
	cmd.Flags().MarkHidden(rpc.FlagRestNacosUrls)
	cmd.Flags().String(rpc.FlagRestNacosNamespaceId, "", "nacos namepace id for discovery service of rest api")
	cmd.Flags().MarkHidden(rpc.FlagRestNacosNamespaceId)
	cmd.Flags().String(rpc.FlagExternalListenAddr, "127.0.0.1:26659", "Set the rest-server external ip and port, when it is launched by Docker")
	// register tendermint rpc to nacos
	cmd.Flags().String(rpc.FlagNacosTmrpcUrls, "", "nacos server urls for discovery service of tendermint rpc")
	cmd.Flags().MarkHidden(rpc.FlagNacosTmrpcUrls)
	cmd.Flags().String(rpc.FlagNacosTmrpcNamespaceID, "", "nacos namepace id for discovery service of tendermint rpc")
	cmd.Flags().MarkHidden(rpc.FlagNacosTmrpcNamespaceID)
	cmd.Flags().String(rpc.FlagNacosTmrpcAppName, "", " tendermint rpc name in nacos")
	cmd.Flags().MarkHidden(rpc.FlagNacosTmrpcAppName)
	cmd.Flags().String(rpc.FlagRpcExternalAddr, "127.0.0.1:26657", "Set the rpc-server external ip and port, when it is launched by Docker (default \"127.0.0.1:26657\")")

	cmd.Flags().String(rpc.FlagRateLimitAPI, "", "Set the RPC API to be controlled by the rate limit policy, such as \"eth_getLogs,eth_newFilter,eth_newBlockFilter,eth_newPendingTransactionFilter,eth_getFilterChanges\"")
	cmd.Flags().Int(rpc.FlagRateLimitCount, 0, "Set the count of requests allowed per second of rpc rate limiter")
	cmd.Flags().Int(rpc.FlagRateLimitBurst, 1, "Set the concurrent count of requests allowed of rpc rate limiter")
	cmd.Flags().Uint64(config.FlagGasLimitBuffer, 50, "Percentage to increase gas limit")
	cmd.Flags().String(rpc.FlagDisableAPI, "", "Set the RPC API to be disabled, such as \"eth_getLogs,eth_newFilter,eth_newBlockFilter,eth_newPendingTransactionFilter,eth_getFilterChanges\"")
	cmd.Flags().Int(config.FlagDynamicGpWeight, 80, "The recommended weight of dynamic gas price [1,100])")
	cmd.Flags().Bool(config.FlagEnableDynamicGp, true, "Enable node to dynamic support gas price suggest")
	cmd.Flags().Bool(eth.FlagEnableMultiCall, false, "Enable node to support the eth_multiCall RPC API")

	cmd.Flags().Bool(token.FlagOSSEnable, false, "Enable the function of exporting account data and uploading to oss")
	cmd.Flags().String(token.FlagOSSEndpoint, "", "The OSS datacenter endpoint such as http://oss-cn-hangzhou.aliyuncs.com")
	cmd.Flags().String(token.FlagOSSAccessKeyID, "", "The OSS access key Id")
	cmd.Flags().String(token.FlagOSSAccessKeySecret, "", "The OSS access key secret")
	cmd.Flags().String(token.FlagOSSBucketName, "", "The OSS bucket name")
	cmd.Flags().String(token.FlagOSSObjectPath, "", "The OSS object path")

	cmd.Flags().Bool(eth.FlagEnableTxPool, false, "Enable the function of txPool to support concurrency call eth_sendRawTransaction")
	cmd.Flags().Uint64(eth.TxPoolCap, 10000, "Set the txPool slice max length")
	cmd.Flags().Int(eth.BroadcastPeriodSecond, 10, "every BroadcastPeriodSecond second check the txPool, and broadcast when it's eligible")

	cmd.Flags().Bool(rpc.FlagEnableMonitor, false, "Enable the rpc monitor and register rpc metrics to prometheus")

	cmd.Flags().String(rpc.FlagKafkaAddr, "", "The address of kafka cluster to consume pending txs")
	cmd.Flags().String(rpc.FlagKafkaTopic, "", "The topic that the kafka writer will produce messages to")

	cmd.Flags().Bool(config.FlagEnableDynamic, false, "Enable dynamic configuration for nodes")
	cmd.Flags().String(config.FlagApollo, "", "Apollo connection config(IP|AppID|NamespaceName) for dynamic configuration")

	// flags for evm trace
	cmd.Flags().Bool(evmtypes.FlagEnableTraces, false, "Enable traces db to save evm transaction trace")
	cmd.Flags().String(evmtypes.FlagTraceSegment, "1-1-0", "Parameters for segmented execution of evm trace, such as \"step-total-num\"")
	cmd.Flags().String(evmtypes.FlagTraceFromAddrs, "", "Generate traces for transactions at specified from addresses (comma separated)")
	cmd.Flags().String(evmtypes.FlagTraceToAddrs, "", "Generate traces for transactions at specified to addresses (comma separated)")
	cmd.Flags().Bool(evmtypes.FlagTraceDisableMemory, false, "Disable memory output for evm trace")
	cmd.Flags().Bool(evmtypes.FlagTraceDisableStack, false, "Disable stack output for evm trace")
	cmd.Flags().Bool(evmtypes.FlagTraceDisableStorage, false, "Disable storage output for evm trace")
	cmd.Flags().Bool(evmtypes.FlagTraceDisableReturnData, false, "Disable return data output for evm trace")
	cmd.Flags().Bool(evmtypes.FlagTraceDebug, false, "Output full trace logs for evm")

	cmd.Flags().Bool(config.FlagPprofAutoDump, false, "Enable auto dump pprof")
	cmd.Flags().String(config.FlagPprofCollectInterval, "5s", "Interval for pprof dump loop")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentMin, 45, "TriggerPercentMin of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentDiff, 50, "TriggerPercentDiff of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentAbs, 50, "TriggerPercentAbs of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentMin, 70, "TriggerPercentMin of mem to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentDiff, 50, "TriggerPercentDiff of mem to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentAbs, 75, "TriggerPercentAbs of cpu mem dump pprof")

	cmd.Flags().String(app.Elapsed, app.DefaultElapsedSchemas, "schemaName=1|0,,,")

	cmd.Flags().String(config.FlagPprofCoolDown, "3m", "The cool down time after every type of pprof dump")
	cmd.Flags().Int64(config.FlagPprofAbciElapsed, 5000, "Elapsed time of abci in millisecond for pprof dump")
	cmd.Flags().Bool(config.FlagPprofUseCGroup, false, "Use cgroup when exchaind run in docker")

	cmd.Flags().String(tmdb.FlagGoLeveldbOpts, "", "Options of goleveldb. (cache_size=128MB,handlers_num=1024)")
	cmd.Flags().String(tmdb.FlagRocksdbOpts, "", "Options of rocksdb. (block_size=4KB,block_cache=1GB,statistics=true,allow_mmap_reads=true,max_open_files=-1)")
	cmd.Flags().String(types.FlagNodeMode, "", "Node mode (rpc|validator|archive) is used to manage flags")

	cmd.Flags().Bool(consensus.EnablePrerunTx, false, "enable proactively runtx mode, default close")
	cmd.Flags().String(automation.ConsensusRole, "", "consensus role")
	cmd.Flags().String(automation.ConsensusTestcase, "", "consensus test case file")

	cmd.Flags().Bool(app.FlagEnableRepairState, false, "Enable auto repair state on start")

	cmd.Flags().Bool(analyzer.FlagEnableAnalyzer, true, "Enable auto open log analyzer")
	cmd.Flags().Bool(sanity.FlagDisableSanity, false, "Disable sanity check")
}
