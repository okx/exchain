package cli

import (
	"bufio"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	utils2 "github.com/okex/exchain/libs/cosmos-sdk/x/mint/client/utils"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/okex/exchain/x/gov"
	"github.com/spf13/cobra"
	"strings"
)

// GetCmdManageTreasuresProposal implements a command handler for submitting a manage treasures proposal transaction
func GetCmdManageTreasuresProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "treasures [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update treasures proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update treasures proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal treasures <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
    "title":"update treasures",
    "description":"update treasures",
    "treasures":[
        {
            "address": "0xA6931Ac6b58E3Db85DFbE1aD408F5096c9736fAE",
            "proportion":"0.1000000000000000"
        }, {
            "address": "0xA6931Ac6b58E3Db85DFbE1aD408F5096c9736fAE",
            "proportion":"0.2000000000000000"
        }，{
            "address": "0xA6931Ac6b58E3Db85DFbE1aD408F5096c9736fAE",
            "proportion":"0.2000000000000000"
        }
    ],
    "is_added":true,
    "deposit":[
        {
            "denom":"okt",
            "amount":"100.000000000000000000"
        }
    ]
}
`, version.ClientName, sdk.DefaultBondDenom,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := utils2.ParseManageTreasuresProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewManageTreasuresProposal(
				proposal.Title,
				proposal.Description,
				proposal.Treasures,
				proposal.IsAdded,
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
