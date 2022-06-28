// nolint
package cli

import (
	"bufio"
	"fmt"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/x/distribution/client/common"
	"github.com/spf13/viper"
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

var (
	flagCommission       = "commission"
	flagMaxMessagesPerTx = "max-msgs"
)

const (
	MaxMessagesPerTxDefault = 5
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
		GetCmdWithdrawAllRewards(cdc, storeKey),
	)...)

	return distTxCmd
}

type generateOrBroadcastFunc func(context.CLIContext, auth.TxBuilder, []sdk.Msg) error

func splitAndApply(
	generateOrBroadcast generateOrBroadcastFunc,
	cliCtx context.CLIContext,
	txBldr auth.TxBuilder,
	msgs []sdk.Msg,
	chunkSize int,
) error {

	if chunkSize == 0 {
		return generateOrBroadcast(cliCtx, txBldr, msgs)
	}

	// split messages into slices of length chunkSize
	totalMessages := len(msgs)
	for i := 0; i < len(msgs); i += chunkSize {

		sliceEnd := i + chunkSize
		if sliceEnd > totalMessages {
			sliceEnd = totalMessages
		}

		msgChunk := msgs[i:sliceEnd]
		if err := generateOrBroadcast(cliCtx, txBldr, msgChunk); err != nil {
			return err
		}
	}

	return nil
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
		Short: "withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
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

			delAddr := cliCtx.GetFromAddress()
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msgs := []sdk.Msg{}
			if viper.GetBool(flagCommission) {
				msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(valAddr))
			} else {
				msgs = append(msgs, types.NewMsgWithdrawDelegatorReward(delAddr, valAddr))
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
	cmd.Flags().Bool(flagCommission, false, "also withdraw validator's commission")
	return cmd
}

// GetCmdWithdrawAllRewards command to withdraw all rewards
func GetCmdWithdrawAllRewards(cdc *codec.Codec, queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-all-rewards",
		Short: "withdraw all delegations rewards for a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw all rewards for a single delegator.

Example:
$ %s tx distribution withdraw-all-rewards --from mykey
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()

			// The transaction cannot be generated offline since it requires a query
			// to get all the validators.
			if cliCtx.GenerateOnly {
				return fmt.Errorf("command disabled with the provided flag: %s", flags.FlagGenerateOnly)
			}

			msgs, err := common.WithdrawAllDelegatorRewards(cliCtx, queryRoute, delAddr)
			if err != nil {
				return err
			}

			chunkSize := viper.GetInt(flagMaxMessagesPerTx)
			return splitAndApply(utils.GenerateOrBroadcastMsgs, cliCtx, txBldr, msgs, chunkSize)
		},
	}

	cmd.Flags().Int(flagMaxMessagesPerTx, MaxMessagesPerTxDefault, "Limit the number of messages per tx (0 for unlimited)")
	return cmd
}

// GetCmdCommunityPoolSpendProposal implements the command to submit a community-pool-spend proposal
func GetCmdCommunityPoolSpendProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
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
			cdc := cdcP.GetCdc()
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

// GetChangeDistributionTypeProposal implements the command to submit a change-distr-type proposal
func GetChangeDistributionTypeProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-distr-type [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a change distribution type proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a change distribution type proposal with the specified value, 0: offchain model, 1:onchain model

Example:
$ %s tx gov submit-proposal change-distr-type <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Change Distribution Type",
  "description": "Will change the distribution type",
  "model": 0,
  "deposit": [
    {
      "denom": "%s",
      "amount": "100.000000000000000000"
    }
  ]
}
`,
				version.ClientName, sdk.DefaultBondDenom,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cdc := cdcP.GetCdc()
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := ParseChangeDistributionTypeProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewChangeDistributionTypeProposal(proposal.Title, proposal.Description, proposal.Type)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
