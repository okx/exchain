package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/okex/okexchain/x/dex/types"
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
		GetCmdQueryProductRank(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryProductsUnderDelisting(queryRoute, cdc),
		GetCmdQueryOperator(queryRoute, cdc),
		GetCmdQueryOperators(queryRoute, cdc),
	)...)

	return queryCmd
}

// GetCmdQueryProducts queries products info
func GetCmdQueryProducts(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Query the list of token pairs",
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			ownerAddress := viper.GetString("owner")
			page := viper.GetUint("page-number")
			perPage := viper.GetUint("items-per-page")
			queryParams := types.NewQueryDexInfoParams(ownerAddress, int(page), int(perPage))
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
	cmd.Flags().UintP("page-number", "p", types.DefaultPage, "page num")
	cmd.Flags().UintP("items-per-page", "i", types.DefaultPerPage, "items per page")

	return cmd
}

// GetCmdQueryDeposits queries deposits about address
func GetCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits [account-addr]",
		Short: "Query product deposits",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			ownerAddress := args[0]
			baseAsset := viper.GetString("base-asset")
			quoteAsset := viper.GetString("quote-asset")
			page := viper.GetUint("page-number")
			perPage := viper.GetUint("items-per-page")
			queryParams := types.NewQueryDepositParams(ownerAddress, baseAsset, quoteAsset, int(page), int(perPage))

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
	cmd.Flags().String("base-asset", "", "base asset")
	cmd.Flags().String("quote-asset", "", "quote asset")
	cmd.Flags().UintP("page-number", "p", types.DefaultPage, "page num")
	cmd.Flags().UintP("items-per-page", "i", types.DefaultPerPage, "items per page")
	return cmd
}

// GetCmdQueryProductRank queries products ranked by deposits
func GetCmdQueryProductRank(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product-rank",
		Short: "Query the rank of token pairs",
		RunE: func(_ *cobra.Command, _ []string) error {
			page := viper.GetUint("page-number")
			perPage := viper.GetUint("items-per-page")
			queryParams := types.NewQueryDexInfoParams("", int(page), int(perPage))
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
	cmd.Flags().UintP("page-number", "p", types.DefaultPage, "page num")
	cmd.Flags().UintP("items-per-page", "i", types.DefaultPerPage, "items per page")
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query all the modifiable parameters of gov proposal",
		Long: strings.TrimSpace(`Query the all the parameters for the governance process:

$ okexchaincli query dex params
`),
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
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

// GetCmdQueryProductsUnderDelisting queries the token pairs involved in dex delisting
func GetCmdQueryProductsUnderDelisting(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "products-delisting",
		Short: "Query the products under dex delisting",
		Long: strings.TrimSpace(`
		Query all the products' names involved in dex delisting:

$ okexchaincli query dex products-delisting`),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

// GetCmdQueryOperator queries operator info
func GetCmdQueryOperator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator [operator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the operator of the account",
		RunE: func(_ *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return sdk.ErrInvalidAddress(fmt.Sprintf("invalid address：%s", args[0]))
			}

			params := types.NewQueryDexOperatorParams(addr)
			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryOperator), bz)
			if err != nil {
				return err
			}
			var operator types.DEXOperator
			cdc.MustUnmarshalJSON(res, &operator)
			return cliCtx.PrintOutput(operator)
		},
	}

	return cmd
}

// GetCmdQueryOperators queries all operator info
func GetCmdQueryOperators(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operators",
		Short: "Query all operator",
		RunE: func(_ *cobra.Command, _ []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryOperators), nil)
			if err != nil {
				return err
			}
			var operators types.DEXOperators
			cdc.MustUnmarshalJSON(res, &operators)
			return cliCtx.PrintOutput(operators)
		},
	}

	return cmd
}

// Strings is just for the object of []string could be inputted into cliCtx.PrintOutput(...)
type Strings []string

func (strs Strings) String() string {
	return strings.Join(strs, "\n")
}
