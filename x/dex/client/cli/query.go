package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/okex/okchain/x/dex/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "dex",
		Short: "Querying commands for the dex module",
	}

	queryCmd.AddCommand(client.GetCommands(
		GetCmdQueryProducts(queryRoute, cdc),
		GetCmdQueryDeposits(queryRoute, cdc),
		GetCmdQueryMatchOrder(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryProductsUnderDelisting(queryRoute, cdc),
	)...)

	return queryCmd
}

// GetCmdQueryProducts queries products info
func GetCmdQueryProducts(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Query the list of token pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			ownerAddress := viper.GetString("owner")
			page := viper.GetInt("page-number")
			perPage := viper.GetInt("items-per-page")
			queryParams, err := types.NewQueryDexInfoParams(ownerAddress, page, perPage)
			if err != nil {
				return err
			}
			bz, err := cdc.MarshalJSON(queryParams)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryProducts), bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().StringP("owner", "", "", "address of the product owner")
	cmd.Flags().IntP("page-number", "p", types.DefaultPage, "page num")
	cmd.Flags().IntP("items-per-page", "i", types.DefaultPerPage, "items per page")

	return cmd
}

// GetCmdQueryDeposits queries deposits about address
func GetCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits [account-addr]",
		Short: "Query product deposits",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ownerAddress := args[0]
			page := viper.GetInt("page-number")
			perPage := viper.GetInt("items-per-page")
			queryParams, err := types.NewQueryDexInfoParams(ownerAddress, page, perPage)
			if err != nil {
				return err
			}
			bz, err := cdc.MarshalJSON(queryParams)
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDeposits), bz)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().IntP("page-number", "p", types.DefaultPage, "page num")
	cmd.Flags().IntP("items-per-page", "i", types.DefaultPerPage, "items per page")
	return cmd
}

// GetCmdQueryDeposits queries match order of products
func GetCmdQueryMatchOrder(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "match-order",
		Short: "Query the match order of token pairs",
		RunE: func(cmd *cobra.Command, args []string) error {
			page := viper.GetInt("page-number")
			perPage := viper.GetInt("items-per-page")
			queryParams, err := types.NewQueryDexInfoParams("", page, perPage)
			if err != nil {
				return err
			}
			bz, err := cdc.MarshalJSON(queryParams)
			if err != nil {
				return err
			}
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryMatchOrder), bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().IntP("page-number", "p", types.DefaultPage, "page num")
	cmd.Flags().IntP("items-per-page", "i", types.DefaultPerPage, "items per page")
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query all the modifiable parameters of gov proposal",
		Long: strings.TrimSpace(`Query the all the parameters for the governance process:

$ okchaincli query dex params
`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters)
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

// GetCmdQueryTokenpairUnderDelisting queries the token pairs involved in dex delisting
func GetCmdQueryProductsUnderDelisting(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "products-delisting",
		Short: "Query the products under dex delisting",
		Long: strings.TrimSpace(`
		Query all the products' names involved in dex delisting:

$ okchaincli query dex products-delisting`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryProductsDelisting), nil)
			if err != nil {
				return err
			}

			var tokenPairNames Strings
			if err := cdc.UnmarshalJSON(res, &tokenPairNames); err != nil {
				return err
			}
			return cliCtx.PrintOutput(tokenPairNames)
		},
	}
}

// just for the object of []string could be inputted into cliCtx.PrintOutput(...)
type Strings []string

func (strs Strings) String() string {
	return strings.Join(strs, "\n")
}
