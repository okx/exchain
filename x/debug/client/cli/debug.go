package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/debug/types"
	"github.com/spf13/cobra"
)

// GetDebugCmd returns the cli query commands for this module
func GetDebugCmd(cdc *codec.Codec) *cobra.Command {

	queryRoute := types.ModuleName

	queryCmd := &cobra.Command{
		Use:   "debug",
		Short: "Debugging subcommands",
	}

	queryCmd.AddCommand(client.GetCommands(
		CmdSetLogLevel(queryRoute, cdc),
		CmdDumpStore(queryRoute, cdc),
		CmdSanityCheck(queryRoute, cdc),
	)...)

	return queryCmd
}

// CmdDumpStore implements the query params command.
func CmdDumpStore(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "dump [module]",
		Args:  cobra.ExactArgs(1),
		Short: "Dump the data of kv-stores by a module name",
		RunE: func(cmd *cobra.Command, args []string) error {

			module := "all"
			if len(args) == 1 {
				module = args[0]
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			bz, err := cdc.MarshalJSON(types.DumpInfoParams{Module: module})
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.DumpStore), bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// CmdSetLogLevel sets log level dynamically
func CmdSetLogLevel(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-loglevel",
		Args:  cobra.ExactArgs(1),
		Short: "Set the okchaind log level",
		Long: strings.TrimSpace(`
$ okchaincli debug set-loglevel "main:info,state:info"

$ okchaincli debug set-loglevel "upgrade:error"
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			if len(args) > 1 {
				return fmt.Errorf("wrong number of arguments for set-loglevel")
			}

			if _, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.SetLogLevel, args[0]), nil); err != nil {
				return err
			}

			fmt.Println("Succeed to set the okchaind log level.")
			return nil
		},
	}
}

// CmdSanityCheck does sanity check
func CmdSanityCheck(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sanity-check",
		Short: "sanity check for all modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.SanityCheck), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}
