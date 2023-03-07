package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	govTypes "github.com/okx/okbchain/x/gov/types"
	"github.com/okx/okbchain/x/params/types"
	"github.com/spf13/cobra"
)

type UpgradeProposalJSON struct {
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Deposit     sdk.SysCoins `json:"deposit" yaml:"deposit"`

	Name         string `json:"name" yaml:"name"`
	ExpectHeight uint64 `json:"expectHeight" yaml:"expectHeight"`
	Config       string `json:"config,omitempty" yaml:"config,omitempty"`
}

func ParseUpgradeProposalJSON(cdc *codec.Codec, proposalFile string) (UpgradeProposalJSON, error) {
	var proposal UpgradeProposalJSON

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

func GetCmdSubmitUpgradeProposal(cdcP *codec.CodecProxy, _ interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a upgrade proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a upgrade proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. proposal name is unique, so if a
proposal's name has been exist, proposal will not be commit successful.

Besides set upgrade's' take effect height, you can also set others config in 'config' field,
which must be a string. You can also omit it if you don't care about it.

IMPORTANT: Only validators or delagators can submit a upgrade proposal. 

Example:
$ %s tx gov submit-proposal upgrade <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "upgrade title",
  "description": "upgrade description",
  "deposit": [
    {
      "denom": "%s",
      "amount": "10000"
    }
  ],
  "name": "YourUpgradeName",
  "expectHeight": "1000",
  "config": "your config string or empty"
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

			proposal, err := ParseUpgradeProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewUpgradeProposal(
				proposal.Title,
				proposal.Description,
				proposal.Name,
				proposal.ExpectHeight,
				proposal.Config,
			)

			msg := govTypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
