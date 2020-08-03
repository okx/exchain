package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/backend/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "backend",
		Short: "Querying commands for the backend module",
	}

	queryCmd.AddCommand(client.GetCommands(
		GetCmdMatches(queryRoute, cdc),
		GetCmdDeals(queryRoute, cdc),
		GetCmdFeeDetails(queryRoute, cdc),
		GetCmdOrderList(queryRoute, cdc),
		GetCmdCandles(queryRoute, cdc),
		GetCmdTickers(queryRoute, cdc),
		GetCmdTxList(queryRoute, cdc),
		GetBlockTxHashesCommand(queryRoute, cdc),
	)...)

	return queryCmd
}

// GetCmdMatches queries match result of a product
func GetCmdMatches(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matches",
		Short: "get match result list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			flags := cmd.Flags()
			product, errProduct := flags.GetString("product")
			startTime, errST := flags.GetInt64("start")
			endTime, errET := flags.GetInt64("end")
			page, errPage := flags.GetInt("page")
			perPage, errPerPage := flags.GetInt("per-page")

			mError := types.NewErrorsMerged(errProduct, errST, errET, errPage, errPerPage)
			if mError != nil {
				return mError
			}

			params := types.NewQueryMatchParams(product, startTime, endTime, page, perPage)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryMatchResults), bz)
			if err != nil {
				fmt.Printf("failed to get matches: %v\n", err)
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().StringP("product", "", "", "filter deals by product")
	cmd.Flags().Int64P("start", "", 0, "filter deals by >= start timestamp")
	cmd.Flags().Int64P("end", "", 0, "filter deals by < end timestamp")
	cmd.Flags().IntP("page", "", 1, "page num")
	cmd.Flags().IntP("per-page", "", 50, "items per page")
	return cmd
}

// GetCmdDeals queries deals
func GetCmdDeals(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deals",
		Short: "get deal list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			flags := cmd.Flags()
			addr, errAddr := flags.GetString("address")
			product, errProduct := flags.GetString("product")
			startTime, errST := flags.GetInt64("start")
			endTime, errET := flags.GetInt64("end")
			page, errPage := flags.GetInt("page")
			perPage, errPerPage := flags.GetInt("per-page")
			side, errSide := flags.GetString("side")

			mError := types.NewErrorsMerged(errAddr, errProduct, errST, errET, errPage, errPerPage, errSide)
			if mError != nil {
				return mError
			}

			params := types.NewQueryDealsParams(addr, product, startTime, endTime, page, perPage, side)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDealList), bz)
			if err != nil {
				fmt.Printf("failed to get deals failed: %v\n", err)
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().StringP("address", "", "", "filter deals by address")
	cmd.Flags().StringP("product", "", "", "filter deals by product")
	cmd.Flags().Int64P("start", "", 0, "filter deals by >= start timestamp")
	cmd.Flags().Int64P("end", "", 0, "filter deals by < end timestamp")
	cmd.Flags().IntP("page", "", 1, "page num")
	cmd.Flags().IntP("per-page", "", 50, "items per page")
	cmd.Flags().StringP("side", "", "", "filter deals by side, support SELL|BUY|ALL, default for empty string means all")
	return cmd
}

// GetCmdCandles queries kline list
func GetCmdCandles(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "klines",
		Short: "get kline list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			flags := cmd.Flags()
			granularity, errGranularity := flags.GetInt("granularity")
			product, errProduct := flags.GetString("product")
			size, errSide := flags.GetInt("limit")

			mError := types.NewErrorsMerged(errGranularity, errProduct, errSide)
			if mError != nil {
				return mError
			}

			params := types.NewQueryKlinesParams(product, granularity, size)
			bz, err := cdc.MarshalJSON(params)
			var out bytes.Buffer
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCandleList), bz)
			if err != nil {
				fmt.Printf("failed to get klines: %v\n", err)
				return nil
			} else {
				if err = json.Indent(&out, res, "", "  "); err != nil {
					fmt.Printf("failed to format by JSON : %v\n", err)
					return nil
				}
			}

			fmt.Println(out.String())
			return nil
		},
	}
	cmd.Flags().IntP("granularity", "g", 60, "[60/180/300/900/1800/3600/7200/14400/21600/43200/86400/604800], second in unit")
	cmd.Flags().StringP("product", "p", "", "name of token pair")
	cmd.Flags().IntP("limit", "", 1, "at most 1000")
	return cmd
}

// GetCmdTickers queries latest ticker list
func GetCmdTickers(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tickers",
		Short: "get latest ticker list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			flags := cmd.Flags()
			count, errCnt := flags.GetInt("limit")
			sort, errSort := flags.GetBool("sort")
			product, errProduct := flags.GetString("product")

			mError := types.NewErrorsMerged(errCnt, errSort, errProduct)
			if mError != nil {
				return mError
			}

			params := types.QueryTickerParams{
				Product: product,
				Count:   count,
				Sort:    sort,
			}

			bz, err := cdc.MarshalJSON(params)
			var out bytes.Buffer
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTickerList), bz)
			if err != nil {
				fmt.Printf("failed to get tickers: %v\n", err)
				return nil
			} else {
				if err = json.Indent(&out, res, "", "  "); err != nil {
					fmt.Printf("failed to format by JSON : %v\n", err)
					return nil
				}
			}

			fmt.Println(out.String())
			return nil
		},
	}
	cmd.Flags().IntP("limit", "", 10, "ticker count")
	cmd.Flags().StringP("product", "p", "", "name of token pair")
	cmd.Flags().BoolP("sort", "s", true, "true or false")
	return cmd
}

// GetCmdFeeDetails queries fee details of a user
func GetCmdFeeDetails(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fees [addr]",
		Short: "get fee detail list of a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr := args[0]
			flags := cmd.Flags()
			page, errPage := flags.GetInt("page")
			perPage, errPerPage := flags.GetInt("per-page")

			mError := types.NewErrorsMerged(errPage, errPerPage)
			if mError != nil {
				return mError
			}

			params := types.NewQueryFeeDetailsParams(addr, page, perPage)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryFeeDetails, args[0]), bz)
			if err != nil {
				fmt.Printf("failed to get fee detail list of %s :%v\n", addr, err)
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().IntP("page", "", 1, "page num")
	cmd.Flags().IntP("per-page", "", 50, "items per page")
	return cmd
}

// GetCmdOrderList queries user's order list
func GetCmdOrderList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orders [open/closed] [addr]",
		Short: "get order list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			if args[0] != "open" && args[0] != "closed" {
				return fmt.Errorf(fmt.Sprintf("order status should be open/closed"))
			}
			addr := args[1]
			flags := cmd.Flags()
			product, errProduct := flags.GetString("product")
			page, errPage := flags.GetInt("page")
			perPage, errPerPage := flags.GetInt("per-page")
			start, errST := flags.GetInt64("start")
			end, errET := flags.GetInt64("end")
			side, errSide := flags.GetString("side")
			hideNoFill, errHide := flags.GetBool("hideNoFill")

			mError := types.NewErrorsMerged(errProduct, errST, errET, errPage, errPerPage, errSide, errHide)
			if mError != nil {
				return mError
			}

			params := types.NewQueryOrderListParams(
				addr, product, side, page, perPage, start, end, hideNoFill)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryOrderList, args[0]), bz)
			if err != nil {
				fmt.Printf("failed to get %s order list of %s :%v\n", args[0], addr, err)
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().StringP("product", "p", "", "filter orders by product")
	cmd.Flags().IntP("page", "", 1, "page num")
	cmd.Flags().IntP("per-page", "", 50, "items per page")
	cmd.Flags().Int64P("start", "", 0, "start timestamp. if start and end is set to 0, it means ignoring time condition.")
	cmd.Flags().Int64P("end", "", 0, "end timestamp. if start and end is set to 0, it means ignoring time condition.")
	cmd.Flags().StringP("side", "", "", "filter deals by side, support SELL|BUY, default for empty string means all")
	cmd.Flags().Bool("hideNoFill", false, "hide orders that have no fills")
	return cmd
}

// GetCmdTxList queries user's transaction history
func GetCmdTxList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs [addr]",
		Short: "get tx list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addr := args[0]
			flags := cmd.Flags()
			txType, errTxType := flags.GetInt64("type")
			startTime, errST := flags.GetInt64("start")
			endTime, errET := flags.GetInt64("end")
			page, errPage := flags.GetInt("page")
			perPage, errPerPage := flags.GetInt("per-page")

			mError := types.NewErrorsMerged(errTxType, errST, errET, errPage, errPerPage)
			if mError != nil {
				return mError
			}

			params := types.NewQueryTxListParams(addr, txType, startTime, endTime, page, perPage)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryTxList), bz)
			if err != nil {
				fmt.Printf("failed to get %s order list of %s :%v\n", args[0], addr, err)
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
	cmd.Flags().Int64P("type", "", 0, "filter txs by txType")
	cmd.Flags().Int64P("start", "", 0, "filter txs by start timestamp")
	cmd.Flags().Int64P("end", "", 0, "filter txs by end timestamp")
	cmd.Flags().IntP("page", "", 1, "page num")
	cmd.Flags().IntP("per-page", "", 50, "items per page")
	return cmd
}

//GetBlockTxHashesCommand queries the tx hashes in the block of the given height
func GetBlockTxHashesCommand(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-tx-hashes [height]",
		Short: "Get txs hash list for a the block at given height",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			txHashes, err := GetBlockTxHashes(cliCtx, &height)
			if err != nil {
				return err
			}

			res, err := json.Marshal(txHashes)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
	return cmd
}

// GetBlockTxHashes return tx hashes in the block of the given height
func GetBlockTxHashes(cliCtx context.CLIContext, height *int64) ([]string, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(height)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(res.Block.Height)
		if err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlockMeta(res.BlockMeta, check)
		if err != nil {
			return nil, err
		}

		err = tmliteProxy.ValidateBlock(res.Block, check)
		if err != nil {
			return nil, err
		}
	}

	txs := res.Block.Txs
	txLen := len(txs)
	txHashes := make([]string, txLen)
	for i, txBytes := range txs {
		txHashes[i] = fmt.Sprintf("%X", tmhash.Sum(txBytes))
	}
	return txHashes, nil
}
