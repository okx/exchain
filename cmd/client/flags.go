package client

import (
	"github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/rpc"
	"github.com/okex/exchain/app/rpc/namespaces/eth"
	"github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/stream"
	"github.com/okex/exchain/x/token"
	"github.com/spf13/cobra"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().Int(watcher.FlagFastQueryLru, 1000, "Set the size of LRU cache under fast-query mode")
	cmd.Flags().Bool(rpc.FlagPersonalAPI, true, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(evmtypes.FlagEnableBloomFilter, false, "Enable bloom filter for event logs")
	cmd.Flags().Int64(filters.FlagGetLogsHeightSpan, 2000, "config the block height span for get logs")
	cmd.Flags().String(stream.NacosTmrpcUrls, "", "Stream plugin`s nacos server urls for discovery service of tendermint rpc")
	cmd.Flags().String(stream.NacosTmrpcNamespaceID, "", "Stream plugin`s nacos namepace id for discovery service of tendermint rpc")
	cmd.Flags().String(stream.NacosTmrpcAppName, "", "Stream plugin`s tendermint rpc name in eureka or nacos")
	cmd.Flags().String(stream.RpcExternalAddr, "127.0.0.1:26657", "Set the rpc-server external ip and port, when it is launched by Docker (default \"127.0.0.1:26657\")")

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

	cmd.Flags().Bool(evmtypes.FlagEnableTraces, false, "enable traces db to save evm transaction trace")
}
