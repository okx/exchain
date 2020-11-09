package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/okexchain/x/common/proto"
	"github.com/okex/okexchain/x/upgrade/types"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/gov"
)

// GetCmdSubmitProposal implements a command handler for submitting a dex list proposal transaction
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "submit a app upgrade proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a app upgrade proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal upgrade <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "app upgrade",
  "description": "Update max validators", 
  "protocol_definition": {
    "version": "1",
    "software": "http://github.com/okex/okexchain/v1",
    "height": "1000",
    "threshold": "0.8",
  }
  "deposit": [
    {
      "denom": common.NativeToken,
      "amount": "100"
    }
  ],
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := parseDexListProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewAppUpgradeProposal(proposal.Title, proposal.Description, proposal.ProtocolDefinition)
			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// UpgradeProposalJSON defines a UpgradeProposal with a deposit used to parse app upgrade proposals from a JSON file
type UpgradeProposalJSON struct {
	Title              string                   `json:"title" yaml:"title"`
	Description        string                   `json:"description" yaml:"description"`
	ProtocolDefinition proto.ProtocolDefinition `json:"protocol_definition" yaml:"protocol_definition"`
	Deposit            sdk.SysCoins             `json:"deposit" yaml:"deposit"`
}

func parseDexListProposalJSON(cdc *codec.Codec, proposalFile string) (UpgradeProposalJSON, error) {
	proposal := UpgradeProposalJSON{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
