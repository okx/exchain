package client

import "github.com/spf13/cobra"

const (
	FlagPersonalAPI = "personal-api"
	FlagFastQuery   = "fast-query"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagPersonalAPI, false, "Enable the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
	cmd.Flags().Bool(FlagFastQuery, false, "Enable the fast query mode for rpc queries")
}
