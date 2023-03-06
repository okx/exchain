package cli

import (
	"fmt"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/spf13/cobra"

	"github.com/okx/okbchain/libs/cosmos-sdk/client"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"
	"github.com/okx/okbchain/x/feesplit/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(moduleName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(flags.GetCommands(
		GetCmdQueryFeeSplits(moduleName, cdc),
		GetCmdQueryFeeSplit(moduleName, cdc),
		GetCmdQueryParams(moduleName, cdc),
		GetCmdQueryDeployerFeeSplits(moduleName, cdc),
		GetCmdQueryWithdrawerFeeSplits(moduleName, cdc),
	)...)

	return cmd
}

// GetCmdQueryFeeSplits implements a command to return all registered contracts
// for fee distribution
func GetCmdQueryFeeSplits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contracts",
		Short: "Query all fee splits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryFeeSplitsRequest{Pagination: pageReq}
			data, err := cliCtx.Codec.MarshalJSON(req)
			if err != nil {
				return err
			}

			// Query store
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryFeeSplits)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var resp types.QueryFeeSplitsResponse
			cdc.MustUnmarshalJSON(bz, &resp)
			return cliCtx.PrintOutput(resp)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "fee splits")
	return cmd
}

// GetCmdQueryFeeSplit implements a command to return a registered contract for fee
// distribution
func GetCmdQueryFeeSplit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contract [contract-address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query a registered contract for fee distribution by hex address",
		Long:    "Query a registered contract for fee distribution by hex address",
		Example: fmt.Sprintf("%s query feesplit contract <contract-address>", version.ClientName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			req := &types.QueryFeeSplitRequest{ContractAddress: args[0]}
			data, err := cliCtx.Codec.MarshalJSON(req)
			if err != nil {
				return err
			}

			// Query store
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryFeeSplit)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var resp types.QueryFeeSplitResponse
			cdc.MustUnmarshalJSON(bz, &resp)
			return cliCtx.PrintOutput(resp)
		},
	}

	return cmd
}

// GetCmdQueryParams implements a command to return the current feesplit
// parameters.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current feesplit module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.QueryParamsResponse
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}

	return cmd
}

// GetCmdQueryDeployerFeeSplits implements a command that returns all contracts
// that a deployer has registered for fee distribution
func GetCmdQueryDeployerFeeSplits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "deployer-contracts [deployer-address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query all contracts that a given deployer has registered for fee distribution",
		Long:    "Query all contracts that a given deployer has registered for fee distribution",
		Example: fmt.Sprintf("%s query feesplit deployer-contracts <deployer-address>", version.ClientName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &types.QueryDeployerFeeSplitsRequest{
				DeployerAddress: args[0],
				Pagination:      pageReq,
			}
			data, err := cliCtx.Codec.MarshalJSON(req)
			if err != nil {
				return err
			}

			// Query store
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDeployerFeeSplits)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var resp types.QueryDeployerFeeSplitsResponse
			cdc.MustUnmarshalJSON(bz, &resp)
			return cliCtx.PrintOutput(resp)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "deployer contracts")
	return cmd
}

// GetCmdQueryWithdrawerFeeSplits implements a command that returns all
// contracts that have registered for fee distribution with a given withdraw
// address
func GetCmdQueryWithdrawerFeeSplits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdrawer-contracts [withdrawer-address]",
		Args:    cobra.ExactArgs(1),
		Short:   "Query all contracts that have been registered for fee distribution with a given withdrawer address",
		Long:    "Query all contracts that have been registered for fee distribution with a given withdrawer address",
		Example: fmt.Sprintf("%s query feesplit withdrawer-contracts <withdrawer-address>", version.ClientName),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			req := &types.QueryWithdrawerFeeSplitsRequest{
				WithdrawerAddress: args[0],
				Pagination:        pageReq,
			}
			data, err := cliCtx.Codec.MarshalJSON(req)
			if err != nil {
				return err
			}

			// Query store
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryWithdrawerFeeSplits)
			bz, _, err := cliCtx.QueryWithData(route, data)
			if err != nil {
				return err
			}

			var resp types.QueryWithdrawerFeeSplitsResponse
			cdc.MustUnmarshalJSON(bz, &resp)
			return cliCtx.PrintOutput(resp)
		},
	}
	flags.AddPaginationFlagsToCmd(cmd, "withdrawer contracts")
	return cmd
}
