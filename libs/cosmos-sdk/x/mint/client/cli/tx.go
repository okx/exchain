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
	types2 "github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/okex/exchain/x/gov"
	"github.com/spf13/cobra"
	"strings"
)

// GetCmdManageContractMethodBlockedListProposal implements a command handler for submitting a manage contract blocked list
// proposal transaction
func GetCmdManageTreasuresProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "treasures [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update contract method blocked list proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update method contract blocked list proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal update-contract-blocked-list <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
    "title":"update contract blocked list proposal with a contract address list",
    "description":"add a contract address list into the blocked list",
    "contract_addresses":[
        {
            "address":"ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc",
            "block_methods": [
                {
                    "Name": "0x371303c0",
                    "Extra": "inc()"
                },
                {
                    "Name": "0x579be378",
                    "Extra": "onc()"
                }
            ]
        },
        {
            "address":"ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9",
            "block_methods": [
                {
                    "Name": "0x371303c0",
                    "Extra": "inc()"
                },
                {
                    "Name": "0x579be378",
                    "Extra": "onc()"
                }
            ]
        }
    ],
    "is_added":true,
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := utils2.ParseManageTreasureProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types2.NewManageTreasuresProposal(
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
