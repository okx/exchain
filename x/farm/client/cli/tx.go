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
	"github.com/okex/okexchain/x/gov"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	farmutils "github.com/okex/okexchain/x/farm/client/utils"
	"github.com/okex/okexchain/x/farm/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	farmTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmTxCmd.AddCommand(client.PostCommands(
		GetCmdCreatePool(cdc),
		GetCmdDestroyPool(cdc),
		GetCmdProvide(cdc),
		GetCmdLock(cdc),
		GetCmdUnlock(cdc),
		GetCmdClaim(cdc),
		GetCmdSetWhite(cdc),
	)...)
	return farmTxCmd
}

func GetCmdCreatePool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pool-name] [min-lock-amount] [yield-token]",
		Short: "create a farm pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a farm pool.

Example:
$ %s tx farm create-pool pool-airtoken1-eth 10eth xxb --from mykey
$ %s tx farm create-pool pool-airtoken1-eth_usdk 10ammswap_eth_usdk xxb --from mykey
`, version.ClientName, version.ClientName),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			poolName := args[0]
			minLockAmount, err := sdk.ParseDecCoin(args[1])
			if err != nil {
				return err
			}
			yieldToken := args[2]
			msg := types.NewMsgCreatePool(cliCtx.GetFromAddress(), poolName, minLockAmount, yieldToken)

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
			fmt.Sprintf(`Destroy a farm pool.

Example:
$ %s tx farm destroy-pool pool-airtoken1-eth --from mykey
`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
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
		Short: "provide a number of yield tokens into a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Provide a number of yield tokens into a pool.

Example:
$ %s tx farm provide pool-airtoken1-eth 1000xxb 5 10000 --from mykey
`, version.ClientName),
		),
		Args: cobra.ExactArgs(4),
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

			startHeightToYield, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			poolName := args[0]
			msg := types.NewMsgProvide(poolName, cliCtx.GetFromAddress(), amount, yieldPerBlock, startHeightToYield)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func GetCmdLock(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [pool-name] [amount]",
		Short: "lock a number of tokens for yield farming",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Lock a number of tokens for yield farming.

Example:
$ %s tx farm lock pool-airtoken1-eth 5eth --from mykey
`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
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
		Short: "unlock a number of tokens",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unlock a number of tokens.

Example:
$ %s tx farm unlock pool-airtoken1-eth 1eth --from mykey
`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
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
		Short: "claim yield farming rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Claim yield farming rewards.

Example:
$ %s tx farm claim --from mykey
`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
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

// GetCmdManageWhiteListProposal implements a command handler for submitting a farm manage white list proposal transaction
func GetCmdManageWhiteListProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "manage-white-list [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a manage white list proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a manage white list proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal manage-white-list <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
 "title": "manage white list with a pool name",
 "description": "add a pool name into the white list",
 "pool_name": "pool-airtoken1-eth",
 "is_added": true,
 "deposit": [
   {
     "denom": "%s",
     "amount": "100"
   }
 ]
}
`, version.ClientName, sdk.DefaultBondDenom,
			)),
		RunE: func(_ *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := farmutils.ParseManageWhiteListProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewManageWhiteListProposal(proposal.Title, proposal.Description, proposal.PoolName, proposal.IsAdded)
			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

//TODO: remove it later
func GetCmdSetWhite(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-white [pool-name]",
		Short: "add specified pool into white list to farm native token",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add pool into white list.

Example:
$ %s tx farm set-white [pool-name] --from mykey
`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			poolName := args[0]
			msg := types.NewMsgSetWhite(poolName, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
