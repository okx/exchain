package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a farm pool with the name of pool, token to be locked in the pool and token to be yielded.

Example:
$ %s tx farm create-pool pool-airtoken1-eth eth xxb --from mykey
$ %s tx farm create-pool pool-airtoken1-eth_usdk ammswap_eth_usdk xxb --from mykey
`, version.ClientName, version.ClientName),
		),
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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Destroy a specific farm pool.

Example:
$ %s tx farm destroy-pool pool-airtoken1-eth --from mykey
`, version.ClientName),
		),
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
		Short: "provide yield-token into a pool, and start mining the coin after the specific height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Provide yield-token into a pool, and start mining the coin after the specific height.

Example:
$ %s tx farm provide pool-airtoken1-eth 1000xxb 5 10000 --from mykey
`, version.ClientName),
		),
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
		Short: "lock a number of tokens for liquidity mining",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Lock a number of tokens for liquidity mining.

Example:
$ %s tx farm lock pool-airtoken1-eth 5eth --from mykey
`, version.ClientName),
		),
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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unlock a number of coins for mining reward.

Example:
$ %s tx farm unlock pool-airtoken1-eth 1eth --from mykey
`, version.ClientName),
		),
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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Claim all the mining rewards till this moment.

Example:
$ %s tx farm claim --from mykey
`, version.ClientName),
		),
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
