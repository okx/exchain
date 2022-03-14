package cli

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for IBC clients
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "IBC client query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
	//GetCmdQueryClientStates(cdc),
	//GetCmdQueryClientState(cdc),
	//GetCmdQueryConsensusStates(cdc),
	//GetCmdQueryConsensusState(cdc),
	//GetCmdQueryHeader(cdc),
	//GetCmdNodeConsensusState(cdc),
	//GetCmdParams(cdc),
	)

	return queryCmd
}
