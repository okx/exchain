// nolint
package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"

	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/gov"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:                        types.ShortUseByCli,
		Short:                      "Distribution transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	distTxCmd.AddCommand(flags.PostCommands(
		GetCmdWithdrawRewards(cdc),
		GetCmdSetWithdrawAddr(cdc),
	)...)

	return distTxCmd
}

// command to replace a delegator's withdrawal address
func GetCmdSetWithdrawAddr(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-withdraw-addr [withdraw-addr]",
		Short: "change the default withdraw address for rewards associated with an address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Set the withdraw address for rewards associated with a delegator address.

Example:
$ %s tx distr set-withdraw-addr ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02 --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()
			withdrawAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}

			msg := types.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdWithdrawRewards command to withdraw rewards
func GetCmdWithdrawRewards(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-rewards [validator-addr]",
		Short: "withdraw validator rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx distr withdraw-rewards exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg --from mykey 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// only withdraw commission of validator
			msg := types.NewMsgWithdrawValidatorCommission(valAddr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdSubmitProposal implements the command to submit a community-pool-spend proposal
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-pool-spend [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a community pool spend proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a community pool spend proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal community-pool-spend <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Community Pool Spend",
  "description": "Pay me some %s!",
  "recipient": "ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02",
  "amount": [
    {
      "denom": "%s",
      "amount": "10000"
    }
  ],
  "deposit": [
    {
      "denom": "%s",
      "amount": "10000"
    }
  ]
}
`,
				version.ClientName, sdk.DefaultBondDenom, sdk.DefaultBondDenom, sdk.DefaultBondDenom,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseCommunityPoolSpendProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, proposal.Recipient, proposal.Amount)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
