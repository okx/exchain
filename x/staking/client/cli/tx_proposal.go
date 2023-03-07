package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okx/okbchain/x/gov"
	"github.com/okx/okbchain/x/staking/types"
	"github.com/spf13/cobra"
)

// GetCmdEditValidatorCommissionRate gets the edit validator commission rate command
func GetCmdEditValidatorCommissionRate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator-commission-rate [commission-rate]",
		Args:  cobra.ExactArgs(1),
		Short: "edit an existing validator commission rate",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddr := cliCtx.GetFromAddress()

			rate, err := sdk.NewDecFromStr(args[0])
			if err != nil {
				return fmt.Errorf("invalid new commission rate: %v", err)
			}

			msg := types.NewMsgEditValidatorCommissionRate(sdk.ValAddress(valAddr), rate)

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCmdProposeValidatorProposal implements a command handler for submitting propose validator proposal transaction
func GetCmdProposeValidatorProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return &cobra.Command{
		Use:   "propose-validator [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a proposal for proposing validator when consensus type is PoA.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal for proposing validator along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal propose-validator <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "propose-validator",
  "description": "propose a validator",
  "is_add": true,
  "deposit": [
    {
      "denom": "okb",
      "amount": "100.000000000000000000"
    }
  ],
  "validator": {
    "delegator_address": "ex1ve4mwgq9967gk338yptsg2fheur4ke32u0gqh3",
    "description": {
      "details": "",
      "identity": "",
      "moniker": "node4",
      "website": ""
    },
    "min_self_delegation": {
      "amount": "0.000000000000000000",
      "denom": "%s"
    },
    "pubkey": "exvalconspub1zcjduepqc4l9dy4g3ghtc6g2wdy0m24tmjju2lggfd0wjpl055tx4knq82squ8ukzn=",
    "validator_address": "exvaloper1tfwvmtfkfzrla52w0u0u07gadkegre9gqk9nel"
  }
}
`, version.ClientName, sdk.DefaultBondDenom,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cdc := cdcP.GetCdc()
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := parseProposeValidatorJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewProposeValidatorProposal(
				proposal.Title,
				proposal.Description,
				proposal.IsAdd,
				proposal.Validator,
			)

			err = content.ValidateBasic()
			if err != nil {
				return err
			}

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

type ProposeValidatorJSON struct {
	Title       string                 `json:"title" yaml:"title"`
	Description string                 `json:"description" yaml:"description"`
	IsAdd       bool                   `json:"is_add" yaml:"is_add"`
	Deposit     sdk.SysCoins           `json:"deposit" yaml:"deposit"`
	Validator   types.ProposeValidator `json:"validator" yaml:"validator"`
}

// ProposeValidatorJSON parses json from proposal file to ProposeValidatorJSON struct
func parseProposeValidatorJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ProposeValidatorJSON, err error) {
	contents, err := os.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
