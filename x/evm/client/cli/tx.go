package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	evmutils "github.com/okex/exchain/x/evm/client/utils"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/gov"
	"github.com/spf13/cobra"
)

// GetCmdManageContractDeploymentWhitelistProposal implements a command handler for submitting a manage contract deployment
// whitelist proposal transaction
func GetCmdManageContractDeploymentWhitelistProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-contract-deployment-whitelist [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update contract deployment whitelist proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update contract deployment whitelist proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal update-contract-deployment-whitelist <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "update contract proposal whitelist with a distributor address list",
  "description": "add a distributor address list into the whitelist",
  "distributor_addresses": [
    "ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02",
    "ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc"
  ],
  "is_added": true,
  "deposit": [
    {
      "denom": "%s",
      "amount": "100.000000000000000000"
    }
  ]
}
`, version.ClientName, sdk.DefaultBondDenom,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := evmutils.ParseManageContractDeploymentWhitelistProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewManageContractDeploymentWhitelistProposal(
				proposal.Title,
				proposal.Description,
				proposal.DistributorAddrs,
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

// GetCmdManageContractBlockedListProposal implements a command handler for submitting a manage contract blocked list
// proposal transaction
func GetCmdManageContractBlockedListProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-contract-blocked-list [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update contract blocked list proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update contract blocked list proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal update-contract-blocked-list <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "update contract blocked list proposal with a contract address list",
  "description": "add a contract address list into the blocked list",
  "contract_addresses": [
    "ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02",
    "ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc"
  ],
  "is_added": true,
  "deposit": [
    {
      "denom": "%s",
      "amount": "100.000000000000000000"
    }
  ]
}
`, version.ClientName, sdk.DefaultBondDenom,
			)),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := evmutils.ParseManageContractBlockedListProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewManageContractBlockedListProposal(
				proposal.Title,
				proposal.Description,
				proposal.ContractAddrs,
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

// GetCmdManageContractMethodBlockedListProposal implements a command handler for submitting a manage contract blocked list
// proposal transaction
func GetCmdManageContractMethodBlockedListProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-contract-method-blocked-list [proposal-file]",
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

			proposal, err := evmutils.ParseManageContractMethodBlockedListProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewManageContractMethodBlockedListProposal(
				proposal.Title,
				proposal.Description,
				proposal.ContractList,
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
