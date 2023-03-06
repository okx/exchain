// nolint
package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/gov"
	"github.com/spf13/cobra"
)

var (
	flagCommission = "commission"
)

// GetCmdWithdrawAllRewards command to withdraw all rewards
func GetCmdWithdrawAllRewards(cdc *codec.Codec, queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-all-rewards",
		Short: "withdraw all delegations rewards for a delegator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw all rewards for a single delegator.

Example:
$ %s tx distr withdraw-all-rewards --from mykey
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
			msg := types.NewMsgWithdrawDelegatorAllRewards(delAddr)
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
  "type": 0,
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

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetWithdrawRewardEnabledProposal implements the command to submit a withdraw-reward-enabled proposal
func GetWithdrawRewardEnabledProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-reward-enabled [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a withdraw reward enabled or disabled proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a withdraw reward enabled or disabled proposal with the specified value, true: enabled, false: disabled

Example:
$ %s tx gov submit-proposal withdraw-reward-enabled <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Withdraw Reward Enabled | Disabled",
  "description": "Will set withdraw reward enabled | disabled",
  "enabled": false,
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

			proposal, err := ParseWithdrawRewardEnabledProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewWithdrawRewardEnabledProposal(proposal.Title, proposal.Description, proposal.Enabled)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetRewardTruncatePrecisionProposal implements the command to submit a reward-truncate-precision proposal
func GetRewardTruncatePrecisionProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-truncate-precision [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a reward truncated precision proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a reward truncated precision proposal with the specified value,

Example:
$ %s tx gov submit-proposal reward-truncate-precision <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
	"title": "Set reward truncated precision",
	"description": "Set distribution reward truncated precision",
	"deposit": [{
		"denom": "%s",
		"amount": "100.000000000000000000"
	}],
	"precision": "0"
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

			proposal, err := ParseRewardTruncatePrecisionProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewRewardTruncatePrecisionProposal(proposal.Title, proposal.Description, proposal.Precision)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
