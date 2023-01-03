// nolint
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/gov"
	staking "github.com/okex/exchain/x/staking/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			fmt.Sprintf(`Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.
If you are a validator, you will withdraw validator commission without param "--commission", and if it is also a delegator, it will withdraw delegator reward.

Example:
$ %s tx distr withdraw-rewards exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg --from mykey(delegator)			# withdraw mykey's reward only
$ %s tx distr withdraw-rewards exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg --from mykey(validator and delegator)	# withdraw mykey's reward only
$ %s tx distr withdraw-rewards exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg --from mykey(validator)	 --commission	# withdraw mykey's commission only
$ %s tx distr withdraw-rewards exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg --from mykey(validator)			# withdraw mykey's commission only
`,
				version.ClientName, version.ClientName, version.ClientName, version.ClientName,
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

			isVal := isValidator(cliCtx, cdc, sdk.ValAddress(delAddr))
			isDel := isDelegator(cliCtx, cdc, delAddr)

			msgs := []sdk.Msg{}
			if viper.GetBool(flagCommission) || (isVal && !isDel) {
				msgs = append(msgs, types.NewMsgWithdrawValidatorCommission(valAddr))
			} else {
				msgs = append(msgs, types.NewMsgWithdrawDelegatorReward(delAddr, valAddr))
				if isVal && isDel {
					fmt.Fprintf(os.Stdout, "%s\n", "\nWarning, check that you are a validator and can add the --commission flag to withdraw your validator commission.")
				}
			}

			return transErrInfo(utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs))
		},
	}
	cmd.Flags().Bool(flagCommission, false, "withdraw validator's commission")
	return cmd
}

func isValidator(cliCtx context.CLIContext, cdc *codec.Codec, valAddress sdk.ValAddress) bool {
	resKVs, _, err := cliCtx.QuerySubspace(staking.ValidatorsKey, staking.StoreKey)
	if err != nil {
		return false
	}

	for _, kv := range resKVs {
		if staking.MustUnmarshalValidator(cdc, kv.Value).GetOperator().Equals(valAddress) {
			return true
		}
	}

	return false
}

func isDelegator(cliCtx context.CLIContext, cdc *codec.Codec, delAddr sdk.AccAddress) bool {
	delegator := staking.NewDelegator(delAddr)
	resp, _, err := cliCtx.QueryStore(staking.GetDelegatorKey(delAddr), staking.StoreKey)
	if err != nil {
		return false
	}
	if len(resp) == 0 {
		return false
	}
	cdc.MustUnmarshalBinaryLengthPrefixed(resp, &delegator)

	if delegator.Tokens.IsZero() {
		return false
	}

	return true
}

func transErrInfo(err error) error {
	check := func(code uint32, describe string) {
		if err == nil {
			return
		}
		if strings.Contains(err.Error(), fmt.Sprint(abci.CodeTypeNonceInc+code)) {
			fmt.Fprintf(os.Stderr, "\nWaring: %s\n", describe)
		}
	}
	check(types.CodeEmptyDelegationDistInfo, "your account(--from) is not a delegator, please check it first.")
	check(types.CodeNoValidatorCommission, "your account(--from) is not a validator, please check it first.")
	check(types.CodeEmptyValidatorDistInfo, "your validator address is error, it's not a validator, please check it first.")
	check(types.CodeEmptyDelegationVoteValidator, "your validator address is error, you haven't voted the validator, please check it first.")

	return err
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
