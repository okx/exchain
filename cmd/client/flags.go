package client

import "github.com/spf13/cobra"

const (
	FlagPersonalAPI        = "personal-api"
	FlagOSSEndpoint        = "oss-endpoint"
	FlagOSSAccessKeyID     = "oss-access-key-id"
	FlagOSSAccessKeySecret = "oss-access-key-secret"
	FlagOSSBucketName      = "oss-bucket-name"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagPersonalAPI, true, "Enable the the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().String(FlagOSSEndpoint, "", "The OSS datacenter endpoint such as http://oss-cn-hangzhou.aliyuncs.com")
	cmd.Flags().String(FlagOSSAccessKeyID, "", "The OSS access key Id")
	cmd.Flags().String(FlagOSSAccessKeySecret, "", "The OSS access key secret")
	cmd.Flags().String(FlagOSSBucketName, "", "The bucket name")
}
