package client

import (
	"github.com/spf13/cobra"

	evmtypes "github.com/okex/okexchain/x/evm/types"
)

const (
	FlagPersonalAPI = "personal-api"
	FlagFastQuery   = "fast-query"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().Bool(FlagPersonalAPI, true, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(evmtypes.FlagEnableBloomFilter, false, "enable bloom filter for logs")
}
