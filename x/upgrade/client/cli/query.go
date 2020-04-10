package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/upgrade/keeper"
	"github.com/okex/okchain/x/upgrade/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	upgradeQueryCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the upgrade module",
	}

	upgradeQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryUpgradeConfig(queryRoute, cdc),
		GetCmdQueryUpgradeVersion(queryRoute, cdc),
		GetCmdQueryUpgradeFailedVersion(queryRoute, cdc))...)

	return upgradeQueryCmd
}

func GetCmdQueryUpgradeConfig(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Query app upgrade config",
		Long: strings.TrimSpace(`Query details about app upgrade config:

$ okchaincli query upgrade config
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, keeper.QueryUpgradeConfig)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var appUpgradeConfig proto.AppUpgradeConfig
			cdc.MustUnmarshalJSON(bz, &appUpgradeConfig)
			return cliCtx.PrintOutput(appUpgradeConfig)
		},
	}
}

func GetCmdQueryUpgradeVersion(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Query app upgrade version",
		Long: strings.TrimSpace(`Query details about current app version:

$ okchaincli query upgrade version
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, keeper.QueryUpgradeVersion)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var version types.QueryVersion
			cdc.MustUnmarshalJSON(bz, &version)
			return cliCtx.PrintOutput(version)
		},
	}
}

func GetCmdQueryUpgradeFailedVersion(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "failed-version",
		Short: "Query app upgrade failed-version",
		Long: strings.TrimSpace(`Query details about last failed app version:

$ okchaincli query upgrade failed-version
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, keeper.QueryUpgradeFailedVersion)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var version types.QueryVersion
			cdc.MustUnmarshalJSON(bz, &version)
			return cliCtx.PrintOutput(version)
		},
	}
}
