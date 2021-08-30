package cli

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	client "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/exchain/x/farm/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group farm queries under a subcommand
	farmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	farmQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryPool(queryRoute, cdc),
			GetCmdQueryPools(queryRoute, cdc),
			GetCmdQueryPoolNum(queryRoute, cdc),
			GetCmdQueryLockInfo(queryRoute, cdc),
			GetCmdQueryEarnings(queryRoute, cdc),
			GetCmdQueryAccount(queryRoute, cdc),
			GetCmdQueryAccountsLockedTo(queryRoute, cdc),
			GetCmdQueryWhitelist(queryRoute, cdc),
			GetCmdQueryParams(queryRoute, cdc),
		)...,
	)

	return farmQueryCmd
}

// GetCmdQueryPool gets the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool [pool-name]",
		Short: "query a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about the kind of coins as reward, the balance and the amount to farm in one block.

Example:
$ %s query farm pool pool-eth-xxb
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bytes, err := cdc.MarshalJSON(types.NewQueryPoolParams(args[0]))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryPool)
			resp, _, err := cliCtx.QueryWithData(route, bytes)
			if err != nil {
				return err
			}

			var pool types.FarmPool
			cdc.MustUnmarshalJSON(resp, &pool)
			return cliCtx.PrintOutput(pool)
		},
	}
}

// GetCmdQueryPools gets the pools query command.
func GetCmdQueryPools(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pools",
		Short: "query for all pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all pools.

Example:
$ %s query farm pools
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// fixed to all pools query
			jsonBytes, err := cdc.MarshalJSON(types.NewQueryPoolsParams(1, 0))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryPools)
			bz, _, err := cliCtx.QueryWithData(route, jsonBytes)
			if err != nil {
				return err
			}

			var pools types.FarmPools
			cdc.MustUnmarshalJSON(bz, &pools)
			return cliCtx.PrintOutput(pools)
		},
	}
}

// GetCmdQueryEarnings gets the earnings query command.
func GetCmdQueryEarnings(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "rewards [pool-name] [address]",
		Short: "query the current rewards of an account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query available rewards for an address.

Example:
$ %s query farm rewards pool-eth-xxb ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			jsonBytes, err := cdc.MarshalJSON(types.NewQueryPoolAccountParams(args[0], accAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryEarnings)
			bz, _, err := cliCtx.QueryWithData(route, jsonBytes)
			if err != nil {
				return err
			}

			var earnings types.Earnings
			cdc.MustUnmarshalJSON(bz, &earnings)
			return cliCtx.PrintOutput(earnings)
		},
	}
}

// GetCmdQueryParams gets the pools query command.
func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "query the current farm parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as farm parameters.

Example:
$ %s query farm params
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryWhitelist gets the whitelist query command.
func GetCmdQueryWhitelist(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "whitelist",
		Short: "query the whitelist of pools to farm okt",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the current whitelist of pools which are approved to farm okt.

Example:
$ %s query farm whitelist
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryWhitelist)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var whitelist types.PoolNameList
			cdc.MustUnmarshalJSON(bz, &whitelist)
			return cliCtx.PrintOutput(whitelist)
		},
	}
}

// GetCmdQueryAccount gets the account query command.
func GetCmdQueryAccount(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "account [address]",
		Short: "query the name of pools that an account has locked coins in",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the names of all pools that an account has locked coins in.

Example:
$ %s query farm account ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			jsonBytes, err := cdc.MarshalJSON(types.NewQueryAccountParams(accAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryAccount)
			bz, _, err := cliCtx.QueryWithData(route, jsonBytes)
			if err != nil {
				return err
			}

			var poolNameList types.PoolNameList
			cdc.MustUnmarshalJSON(bz, &poolNameList)
			return cliCtx.PrintOutput(poolNameList)
		},
	}
}

// GetCmdQueryPoolNum gets the pool number query command.
func GetCmdQueryPoolNum(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool-num",
		Short: "query the number of pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the number of pools that already exist.

Example:
$ %s query farm pool-num
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryPoolNum)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var poolNum types.PoolNum
			cdc.MustUnmarshalJSON(bz, &poolNum)
			return cliCtx.PrintOutput(poolNum)
		},
	}
}

// GetCmdQueryAccountsLockedTo gets all addresses of accounts that locked coins in a specific pool
func GetCmdQueryAccountsLockedTo(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "accounts-locked-to [pool-name]",
		Short: "query the addresses of accounts locked in a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the addresses of accounts that have locked coins in a specific pool.

Example:
$ %s query farm accounts-locked-to pool-eth-xxb
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			jsonBytes, err := cdc.MarshalJSON(types.NewQueryPoolParams(args[0]))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryAccountsLockedTo)
			bz, _, err := cliCtx.QueryWithData(route, jsonBytes)
			if err != nil {
				return err
			}

			var accAddrList types.AccAddrList
			cdc.MustUnmarshalJSON(bz, &accAddrList)
			return cliCtx.PrintOutput(accAddrList)
		},
	}
}

// GetCmdQueryLockInfo gets the lock info of an account's token locking on a specific pool
func GetCmdQueryLockInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "lock-info [pool-name] [address]",
		Short: "query the lock info of an account on a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the lock info of an account's token locking on a specific pool.

Example:
$ %s query farm lock-info pool-eth-xxb ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			accAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			jsonBytes, err := cdc.MarshalJSON(types.NewQueryPoolAccountParams(args[0], accAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryLockInfo)
			bz, _, err := cliCtx.QueryWithData(route, jsonBytes)
			if err != nil {
				return err
			}

			var lockInfo types.LockInfo
			cdc.MustUnmarshalJSON(bz, &lockInfo)
			return cliCtx.PrintOutput(lockInfo)
		},
	}
}
