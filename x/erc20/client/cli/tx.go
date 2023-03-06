package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	govcli "github.com/okx/okbchain/libs/cosmos-sdk/x/gov/client/cli"
	"github.com/okx/okbchain/x/erc20/types"
	"github.com/okx/okbchain/x/gov"
	"github.com/spf13/cobra"
)

// GetCmdTokenMappingProposal returns a CLI command handler for creating
// a token mapping proposal governance transaction.
func GetCmdTokenMappingProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-mapping [denom] [contract]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a token mapping proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a token mapping proposal.

Example:
$ %s tx gov submit-proposal token-mapping xxb 0x0000...0000 --from=<key_or_address>
`, version.ClientName,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cdc := cdcP.GetCdc()
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
				if common.IsHexAddress(args[1]) {
					addr := common.HexToAddress(args[1])
					contract = &addr
				} else {
					return fmt.Errorf("invalid contract address %s", args[1])
				}
			}

			content := types.NewTokenMappingProposal(
				title, description, args[0], contract,
			)
			if err := content.ValidateBasic(); err != nil {
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

func GetCmdProxyContractRedirectProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cdc := cdcP.GetCdc()
	cmd := &cobra.Command{
		Use:   "contract-redirect [denom] [tp] [contract|owner]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a proxy contract redirect proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proxy contract redirect proposal.
  tp: 
	0	implementation	implementation address
	1	owner			owner address
Example:
$ %s tx gov submit-proposal contract-redirect xxb 0 0xffffffffffffffffffff
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

			tp, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid redirect tp,only support 0(change contract) or 1(change owner),but give %s", args[1])
			}
			var redirectaddr common.Address
			if len(args[1]) > 0 {
				if common.IsHexAddress(args[2]) {
					redirectaddr = common.HexToAddress(args[2])
				} else {
					return fmt.Errorf("invalid contract address %s", args[1])
				}
			}

			content := types.NewProxyContractRedirectProposal(
				title, description, args[0], types.RedirectType(tp), &redirectaddr,
			)
			if err := content.ValidateBasic(); err != nil {
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

func SetContractTemplateProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cdc := cdcP.GetCdc()
	cmd := &cobra.Command{
		Use:   "contract-template  [file-path] proxy/implement",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a new bytecode template contract proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a contract template proposal.
Example:
$ %s tx gov submit-proposal contract-template ~/template.json proxy/implement --from=<key_or_address>
the template.json should be like : 
{
	"abi":[xxxxx],
	"bin":"xxxx"
}
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

			data, err := ioutil.ReadFile(args[0])
			if nil != err {
				return err
			}

			_, err = types.UnmarshalCompileContract(data)
			if nil != err {
				return err
			}

			content := types.NewContractTemplateProposal(
				title, description, args[1], string(data),
			)
			if err := content.ValidateBasic(); err != nil {
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
