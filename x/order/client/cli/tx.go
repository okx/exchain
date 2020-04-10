package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/okex/okchain/x/order/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "order",
		Short: "Order transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		GetCmdNewOrder(cdc),
		GetCmdCancelOrder(cdc),
	)...)

	return txCmd
}

func GetCmdNewOrder(cdc *codec.Codec) *cobra.Command {
	// new order flags
	var product string
	var side string
	var price string
	var quantity string
	cmd := &cobra.Command{
		Use:   "new",
		Short: "place a new order",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(product) == 0 || len(side) == 0 || len(price) == 0 || len(quantity) == 0 {
				return errors.New("invalid param format")
			}
			if len(args) > 0 {
				return errors.New(`invalid param format. tips:use comma "," to place multi orders`)
			}
			productArr := strings.Split(product, ",")

			if len(productArr) == 0 {
				return errors.New("invalid param counts")
			}

			err := handleNewOrder(cdc, product, side, price, quantity)
			return err

		},
	}

	cmd.Flags().StringVarP(&product, "product", "", "", "Trading pair in full name of the tokens: ${baseAssetSymbol}_${quoteAssetSymbol}, for example \"mycoin_okt\".")
	cmd.Flags().StringVarP(&side, "side", "s", "", "BUY or SELL (default \"SELL\")")
	cmd.Flags().StringVarP(&price, "price", "p", "", "The price of the order")
	cmd.Flags().StringVarP(&quantity, "quantity", "q", "", "The quantity of the order")
	return cmd
}

func handleNewOrder(cdc *codec.Codec, product string, side string, price string, quantity string) error {
	var items []types.OrderItem
	productArr := strings.Split(product, ",")
	sideArr := strings.Split(side, ",")
	priceArr := strings.Split(price, ",")
	quantityArr := strings.Split(quantity, ",")
	if len(productArr) != len(sideArr) {
		return errors.New("invalid param side counts")
	}

	if len(productArr) != len(priceArr) {
		return errors.New("invalid param price counts")
	}

	if len(productArr) != len(quantityArr) {
		return errors.New("invalid param quantity counts")
	}

	for i := 0; i < len(productArr); i++ {
		product := productArr[i]
		side := sideArr[i]
		price, err := sdk.NewDecFromStr(priceArr[i])
		if err != nil {
			return errors.New(err.Error())
		}
		quantity, err := sdk.NewDecFromStr(quantityArr[i])
		if err != nil {
			return errors.New(err.Error())
		}
		items = append(items, types.OrderItem{
			Product:  product,
			Side:     side,
			Price:    price,
			Quantity: quantity,
		})
	}

	txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	msg := types.NewMsgNewOrders(cliCtx.GetFromAddress(), items)
	err := utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
	return err
}

func GetCmdCancelOrder(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel [order-id]",
		Short: "cancel order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orderIDs := strings.Split(args[0], ",")

			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgCancelOrders(cliCtx.GetFromAddress(), orderIDs)
			err := utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			if err != nil {
				fmt.Println(err)
			}
			return err
		},
	}
}
