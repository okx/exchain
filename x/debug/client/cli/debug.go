package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/exchain/x/debug/types"
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
		CmdInvariantCheck(queryRoute, cdc),
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
		Short: "Set the exchaind log level",
		Long: strings.TrimSpace(`
$ exchaincli debug set-loglevel "main:info,state:info"

$ exchaincli debug set-loglevel "upgrade:error"
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

			fmt.Println("Succeed to set the exchaind log level.")
			return nil
		},
	}
}

// CmdSanityCheck does sanity check
func CmdSanityCheck(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sanity-check-shares",
		Short: "check the total share of validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.SanityCheckShares), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// CmdInvariantCheck does invariants check
func CmdInvariantCheck(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "invariant-check",
		Short: "check the invariant of all module",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.InvariantCheck), nil)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}
