package client

import (
	evmtypes "github.com/okex/okexchain/x/evm/types"
	"github.com/okex/okexchain/x/evm/watcher"
	"github.com/spf13/cobra"
)

const (
	FlagPersonalAPI = "personal-api"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().Bool(FlagPersonalAPI, true, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(evmtypes.FlagEnableBloomFilter, false, "Enable bloom filter for event logs")
}
