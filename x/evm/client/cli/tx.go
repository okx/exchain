package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	evmutils "github.com/okx/okbchain/x/evm/client/utils"
	"github.com/okx/okbchain/x/evm/types"
	"github.com/okx/okbchain/x/gov"
	"github.com/spf13/cobra"
)

// GetCmdManageContractDeploymentWhitelistProposal implements a command handler for submitting a manage contract deployment
// whitelist proposal transaction
func GetCmdManageContractDeploymentWhitelistProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
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
			cdc := cdcP.GetCdc()
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
func GetCmdManageContractBlockedListProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
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
			cdc := cdcP.GetCdc()
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
func GetCmdManageContractMethodBlockedListProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
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
                    "sign": "0x371303c0",
                    "extra": "inc()"
                },
                {
                    "sign": "0x579be378",
                    "extra": "onc()"
                }
            ]
        },
        {
            "address":"ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9",
            "block_methods": [
                {
                    "sign": "0x371303c0",
                    "extra": "inc()"
                },
                {
                    "sign": "0x579be378",
                    "extra": "onc()"
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
			cdc := cdcP.GetCdc()
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

// GetCmdManageSysContractAddressProposal implements a command handler for submitting a manage system contract address
// proposal transaction
func GetCmdManageSysContractAddressProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return &cobra.Command{
		Use:   "system-contract-address [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a system contract address proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a system contract address proposal.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal system-contract-address <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title":"Update system contract address",
  "description":"Will change the system contract address",
  "contract_addresses": "0x1033796B018B2bf0Fc9CB88c0793b2F275eDB624",
  "is_added":true,
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
			cdc := cdcP.GetCdc()
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := evmutils.ParseManageSysContractAddressProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewManageSysContractAddressProposal(
				proposal.Title,
				proposal.Description,
				proposal.ContractAddr,
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

// GetCmdManageContractMethodGuFactorProposal implements a command handler for submitting a manage contract method gu-factor proposal transaction
func GetCmdManageContractMethodGuFactorProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return &cobra.Command{
		Use:   "update-contract-method-gu-factor [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update contract method gu-factor proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update method contract gu-factor proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal update-contract-method-gu-factor <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
    "title":"update contract method gu-factor proposal with a contract address list",
    "description":"add a contract method gu-factor list into chain",
    "contract_addresses":[
        {
            "address":"ex1k0wwsg7xf9tjt3rvxdewz42e74sp286agrf9qc",
            "block_methods": [
                {
                    "sign": "0x371303c0",
                    "extra": "{\"gu_factor\":\"10.000000000000000000\"}"
                },
                {
                    "sign": "0x579be378",
                    "extra": "{\"gu_factor\":\"20.000000000000000000\"}"
                }
            ]
        },
        {
            "address":"ex1s0vrf96rrsknl64jj65lhf89ltwj7lksr7m3r9",
            "block_methods": [
                {
                    "sign": "0x371303c0",
                    "extra": "{\"gu_factor\":\"30.000000000000000000\"}"
                },
                {
                    "sign": "0x579be378",
                    "extra": "{\"gu_factor\":\"40.000000000000000000\"}"
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
			cdc := cdcP.GetCdc()
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

// GetCmdManageContractByteCodeProposal implements a command handler for submitting a manage contract bytecode proposal transaction
func GetCmdManageContractByteCodeProposal(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	return &cobra.Command{
		Use:   "update-contract-bytecode [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an update contract bytecode proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit an update contract bytecode proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal update-contract-bytecode <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
    "title":"update contract bytecode",
    "description":"update contract bytecode",
    "contract":"0x9a59ae3Fc0948717F94242fc170ac1d5dB3f0D5D",
    "substitute_contract":"0xFc0b06f1C1e82eFAdC0E5c226616B092D2cb97fF",
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

			proposal, err := evmutils.ParseManageContractBytecodeProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			content := types.NewManageContractByteCodeProposal(
				proposal.Title,
				proposal.Description,
				proposal.Contract,
				proposal.SubstituteContract,
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
