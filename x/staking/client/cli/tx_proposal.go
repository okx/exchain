package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/staking/types"
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
    "title":"propose-validator",
    "description":"propose a validator",
	"isAdd": true,
    "block_num":123456,
    "deposit":[
        {
            "denom":"%s",
            "amount":"100.000000000000000000"
        }
    ]
}
`, version.ClientName, sdk.DefaultBondDenom,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cdc := cdcP.GetCdc()
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := parseProposeValidatorProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewProposeValidatorProposal(
				proposal.Title,
				proposal.Description,
				proposal.IsAdd,
				proposal.BlockNum,
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

// ManageTreasuresProposalJSON defines a ManageTreasureProposal with a deposit used to parse
// manage treasures proposals from a JSON file.
type ManageTreasuresProposalJSON struct {
	Title       string                 `json:"title" yaml:"title"`
	Description string                 `json:"description" yaml:"description"`
	IsAdd       bool                   `json:"is_add" yaml:"is_add"`
	BlockNum    uint64                 `json:"block_num" yaml:"block_num"`
	Deposit     sdk.SysCoins           `json:"deposit" yaml:"deposit"`
	Validator   types.ProposeValidator `json:"validator" yaml:"validator"`
}

// parseProposeValidatorProposalJSON parses json from proposal file to ProposeValidatorProposalJSON struct
func parseProposeValidatorProposalJSON(cdc *codec.Codec, proposalFilePath string) (
	proposal ManageTreasuresProposalJSON, err error) {
	contents, err := os.ReadFile(proposalFilePath)
	if err != nil {
		return
	}

	cdc.MustUnmarshalJSON(contents, &proposal)
	return
}
