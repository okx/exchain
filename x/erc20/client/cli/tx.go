package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	govcli "github.com/okex/exchain/libs/cosmos-sdk/x/gov/client/cli"
	"github.com/okex/exchain/x/erc20/types"
	"github.com/okex/exchain/x/gov"
	"github.com/spf13/cobra"
)

// GetCmdTokenMappingProposal returns a CLI command handler for creating
// a token mapping change proposal governance transaction.
func GetCmdTokenMappingProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-mapping [denom] [contract]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a token mapping change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a token mapping change proposal.

Example:
$ %s tx gov submit-proposal token-mapping-change xxb 0x0000...0000 --from=<key_or_address>
`, version.ClientName,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			title, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return err
			}

			var contract *common.Address
			if len(args[1]) > 0 {
				addr := common.HexToAddress(args[1])
				contract = &addr
			}

			content := types.NewTokenMappingChangeProposal(
				title, description, args[0], contract,
			)

			err = content.ValidateBasic()
			if err != nil {
				return err
			}

			strDeposit, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(strDeposit)
			if err != nil {
				return err
			}

			msg := gov.NewMsgSubmitProposal(content, deposit, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(govcli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "deposit of proposal")

	return cmd
}
