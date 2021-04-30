package client

import (
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/stream"
	"github.com/okex/exchain/x/token"
	"github.com/spf13/cobra"
)

const (
	FlagPersonalAPI       = "personal-api"
	FlagGetLogsHeightSpan = "logs-height-span"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().Bool(FlagPersonalAPI, true, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(evmtypes.FlagEnableBloomFilter, false, "Enable bloom filter for event logs")
	cmd.Flags().Int64(FlagGetLogsHeightSpan, -1, "config the block height span for get logs")
	cmd.Flags().String(stream.NacosTmrpcUrls, "", "Stream plugin`s nacos server urls for discovery service of tendermint rpc")
	cmd.Flags().String(stream.NacosTmrpcNamespaceID, "", "Stream plugin`s nacos namepace id for discovery service of tendermint rpc")
	cmd.Flags().String(stream.NacosTmrpcAppName, "", "Stream plugin`s tendermint rpc name in eureka or nacos")
	cmd.Flags().String(stream.RpcExternalAddr, "127.0.0.1:26657", "Set the rpc-server external ip and port, when it is launched by Docker (default \"127.0.0.1:26657\")")

	cmd.Flags().Bool(token.FlagOSSEnable, false, "Enable the function of exporting account data and uploading to oss")
	cmd.Flags().String(token.FlagOSSEndpoint, "", "The OSS datacenter endpoint such as http://oss-cn-hangzhou.aliyuncs.com")
	cmd.Flags().String(token.FlagOSSAccessKeyID, "", "The OSS access key Id")
	cmd.Flags().String(token.FlagOSSAccessKeySecret, "", "The OSS access key secret")
	cmd.Flags().String(token.FlagOSSBucketName, "", "The OSS bucket name")
	cmd.Flags().String(token.FlagOSSObjectPath, "", "The OSS object path")
}
