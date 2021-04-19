package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	client "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/exchain/x/order/keeper"
	"github.com/okex/exchain/x/order/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group order queries under a subcommand
	queryCmd := &cobra.Command{
		Use:   "order",
		Short: "Querying commands for the order module",
	}

	queryCmd.AddCommand(client.GetCommands(
		GetCmdQueryOrder(queryRoute, cdc),
		GetCmdDepthBook(queryRoute, cdc),
		GetCmdQueryStore(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
	)...)

	queryCmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	return queryCmd
}

// GetCmdQueryOrder queries order info by orderID
func GetCmdQueryOrder(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "detail [order-id]",
		Short: "Query an order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			orderID := args[0]

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryOrderDetail, orderID),
				nil)
			if err != nil {
				fmt.Printf("order does not exist - %s \n", orderID)
				return nil
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdDepthBook queries order book about a product
func GetCmdDepthBook(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depthbook [product]",
		Short: "Query the depth book of a trading pair",
		Long: strings.TrimSpace(`Query the depth book of a trading pair:

$ exchaincli query depthbook mytoken_okt

The 'product' is a trading pair in full name of the tokens: ${base-asset-symbol}_${quote-asset-symbol}, for example 'mytoken_okt'.
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			product := args[0]
			size := viper.GetUint("size")
			params := keeper.NewQueryDepthBookParams(product, size)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDepthBook),
				bz)
			if err != nil {
				fmt.Printf("get depth book of %s failed: %v\n", product, err.Error())
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().Uint("size", keeper.DefaultBookSize, "depth book single-side size")
	return cmd
}

// GetCmdQueryStore queries store statistic
func GetCmdQueryStore(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "store",
		Short: "query store of order module",
		RunE: func(cmd *cobra.Command, args []string) error {
			dbFile := viper.GetString("app_dbpath")
			//dump := viper.GetBool("dump")
			if dbFile != "" {
				//// query through db file
				//mapp, db, err := filedb.GetMockApp(dbFile)
				//if err != nil {
				//	return err
				//}
				//defer filedb.CloseApp(db)
				//
				//kp := mapp.OrderKeeper
				//mapp.BeginBlock(abci.RequestBeginBlock{})
				//ctx := mapp.GetState(baseapp.RunTxModeDeliver()).Context()
				//if dump {
				//	keeper.DumpStore(ctx, kp)
				//	return nil
				//}
				//
				//ss := keeper.GetStoreStatistic(ctx, kp)
				//res, errRes := codec.MarshalJSONIndent(cdc, ss)
				//if errRes != nil {
				//	return sdk.ErrInternal(
				//		sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
				//}
				//fmt.Println(string(res))
			} else {
				cliCtx := context.NewCLIContext().WithCodec(cdc)
				res, _, err := cliCtx.QueryWithData(
					fmt.Sprintf("custom/order/%s", types.QueryStore), nil)
				if err != nil {
					fmt.Printf("query store failed: %v\n", err.Error())
					return nil
				}

				fmt.Println(string(res))
				return nil
			}
			return nil
		},
	}
	cmd.Flags().String("dbpath", "", "db path (if this path is given, query through local file)")
	cmd.Flags().Bool("dump", false, "dump all key-value constants of specified module")
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the order process",
		Long: strings.TrimSpace(`Query the all the parameters for the governance process:

$ exchaincli query order params
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
