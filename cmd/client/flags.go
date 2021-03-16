package client

import "github.com/spf13/cobra"

const (
	FlagPersonalAPI = "personal-api"
)

func RegisterAppFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(FlagPersonalAPI, true, "Enable the the personal_ prefixed set of APIs in the Web3 JSON-RPC spec")
}
