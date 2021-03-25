package client

import (
	evmtypes "github.com/okex/okexchain/x/evm/types"
	"github.com/okex/okexchain/x/evm/watcher"
	"github.com/spf13/cobra"
)

const (
	FlagPersonalAPI        = "personal-api"
	FlagCloseMutex         = "close-mutex"
	FlagOSSEndpoint        = "oss-endpoint"
	FlagOSSAccessKeyID     = "oss-access-key-id"
	FlagOSSAccessKeySecret = "oss-access-key-secret"
	FlagOSSBucketName      = "oss-bucket-name"
	FlagOSSObjectPath      = "oss-object-path"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().Bool(FlagPersonalAPI, true, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(evmtypes.FlagEnableBloomFilter, false, "Enable bloom filter for event logs")
	cmd.Flags().Bool(FlagCloseMutex, false, "Close local client query mutex for better concurrency")
	cmd.Flags().String(FlagOSSEndpoint, "", "The OSS datacenter endpoint such as http://oss-cn-hangzhou.aliyuncs.com")
	cmd.Flags().String(FlagOSSAccessKeyID, "", "The OSS access key Id")
	cmd.Flags().String(FlagOSSAccessKeySecret, "", "The OSS access key secret")
	cmd.Flags().String(FlagOSSBucketName, "", "The OSS bucket name")
	cmd.Flags().String(FlagOSSObjectPath, "", "The OSS object path")
}
