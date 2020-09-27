package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okexchain/x/farm/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	farmTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmTxCmd.AddCommand(
		GetCmdCreatePool(cdc),
		GetCmdDestroyPool(cdc),
		GetCmdProvide(cdc),
		GetCmdLock(cdc),
		GetCmdUnlock(cdc),
		GetCmdClaim(cdc),
	)
	return farmTxCmd
}

func GetCmdCreatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-name] [lock-token] [yield-token]",
		Short: "create a farm pool with the name of pool, token to be locked in the pool and token to be yielded",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			poolName := args[0]
			lockToken := args[1]
			yieldToken := args[2]
			msg := types.NewMsgCreatePool(cliCtx.GetFromAddress(), poolName, lockToken, yieldToken)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdDestroyPool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy-pool [pool-name]",
		Short: "destroy a farm pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			poolName := args[0]
			msg := types.NewMsgDestroyPool(cliCtx.GetFromAddress(), poolName)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdProvide(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provide [pool-name] [amount] [yield-per-block] [start-height-to-yield]",
		Short: "provide yiled-token into a pool, and start mining the token after the specified height",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}

			yieldPerBlock, err := sdk.NewDecFromStr(args[2])
			if err != nil {
				return err
			}

			startHeightToYiled, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			poolName := args[0]
			msg := types.NewMsgProvide(poolName, cliCtx.GetFromAddress(), amount, yieldPerBlock, startHeightToYiled)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdLock(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [pool-name] [amount]",
		Short: "lock a number of coins for liquidity mining",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}

			poolName := args[0]
			msg := types.NewMsgLock(poolName, cliCtx.GetFromAddress(), amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdUnlock(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock [pool-name] [amount]",
		Short: "unlock a number of coins for mining reward",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}

			poolName := args[0]
			msg := types.NewMsgUnlock(poolName, cliCtx.GetFromAddress(), amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdClaim(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [pool-name]",
		Short: "claim all the mining rewards till this moment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			poolName := args[0]
			msg := types.NewMsgClaim(poolName, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
