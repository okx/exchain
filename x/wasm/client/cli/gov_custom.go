package cli

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okx/okbchain/x/gov"
	govcli "github.com/okx/okbchain/x/gov/client/cli"
	"github.com/okx/okbchain/x/wasm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ProposalUpdateDeploymentWhitelistCmd(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-wasm-deployment-whitelist [comma-separated addresses]",
		Short: "Submit an update wasm contract deployment whitelist proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cdc := cdcP.GetCdc()
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposalTitle, err := cmd.Flags().GetString(govcli.FlagTitle)
			if err != nil {
				return fmt.Errorf("proposal title: %s", err)
			}
			proposalDescr, err := cmd.Flags().GetString(govcli.FlagDescription)
			if err != nil {
				return fmt.Errorf("proposal description: %s", err)
			}
			addrs := strings.Split(strings.TrimSpace(args[0]), ",")
			sort.Strings(addrs)

			proposal := types.UpdateDeploymentWhitelistProposal{
				Title:            proposalTitle,
				Description:      proposalDescr,
				DistributorAddrs: addrs,
			}

			err = proposal.ValidateBasic()
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(viper.GetString(govcli.FlagDeposit))
			if err != nil {
				return err
			}

			msg := gov.NewMsgSubmitProposal(&proposal, deposit, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	// proposal flags
	cmd.Flags().String(govcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "Deposit of proposal")

	return cmd
}

const isDelete = "delete"

func ProposalUpdateWASMContractMethodBlockedListCmd(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-wasm-contract-method-blocked-list [contract address] [comma-separated methods]",
		Short: "Submit an update wasm contract deployment whitelist proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cdc := cdcP.GetCdc()
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			methods := strings.Split(strings.TrimSpace(args[1]), ",")
			sort.Strings(methods)
			var extraMethods []*types.Method
			for _, m := range methods {
				extraMethods = append(extraMethods, &types.Method{
					Name: m,
				})
			}

			proposal := types.UpdateWASMContractMethodBlockedListProposal{
				Title:       viper.GetString(govcli.FlagTitle),
				Description: viper.GetString(govcli.FlagDescription),
				BlockedMethods: &types.ContractMethods{
					ContractAddr: args[0],
					Methods:      extraMethods,
				},
				IsDelete: viper.GetBool(isDelete),
			}

			if err := proposal.ValidateBasic(); err != nil {
				return err
			}

			deposit, err := sdk.ParseCoins(viper.GetString(govcli.FlagDeposit))
			if err != nil {
				return err
			}

			msg := gov.NewMsgSubmitProposal(&proposal, deposit, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	// proposal flags
	cmd.Flags().String(govcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().Bool(isDelete, false, "True to delete methods and default to add")

	return cmd
}
