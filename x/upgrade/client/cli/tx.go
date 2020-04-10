package cli

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/upgrade/types"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/gov"
)

// GetCmdSubmitProposal implements a command handler for submitting a dex list proposal transaction.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a app upgrade proposal",
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
    "software": "http://github.com/okex/okchain/v1",
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

			proposal, err := ParseDexListProposalJSON(cdc, args[0])
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

// UpgradeProposalJSON defines a UpgradeProposal with a deposit used
// to parse app upgrade proposals from a JSON file.
type UpgradeProposalJSON struct {
	Title              string                   `json:"title" yaml:"title"`
	Description        string                   `json:"description" yaml:"description"`
	ProtocolDefinition proto.ProtocolDefinition `json:"protocol_definition" yaml:"protocol_definition"`
	Deposit            sdk.DecCoins             `json:"deposit" yaml:"deposit"`
}

func ParseDexListProposalJSON(cdc *codec.Codec, proposalFile string) (UpgradeProposalJSON, error) {
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

//just 4 test
//////////////////////////////////////////////////////////////////////////////////
func GetCmdUpgradeConfig(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [proposal-id] [version] [height] [software]",
		Args:  cobra.ExactArgs(4),
		Short: "config app upgrade",
		Long: strings.TrimSpace(`Config an app upgrade:

$ okchaincli tx upgrade config 0 1 10 http://web.abc
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			proposalID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			ver, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			height, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			software := args[3]

			msg := types.NewMsgUpgradeConfig(uint64(proposalID), uint64(ver), uint64(height), software, cliCtx.FromAddress)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	_ = cmd.MarkFlagRequired(client.FlagFrom)
	return cmd
}

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "[test] Upgrade transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdUpgradeConfig(cdc),
	)...)

	return cmd
}

//////////////////////////////////////////////////////////////////////////////////
