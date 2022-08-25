package cli

import (
	"bufio"
	"fmt"
	"github.com/okex/exchain/x/wasm/types"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	govcli "github.com/okex/exchain/libs/cosmos-sdk/x/gov/client/cli"
	"github.com/okex/exchain/x/gov"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	isAdded = "add"
)

func ProposalUpdateDeploymentWhitelistCmd(cdcP *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-wasm-deployment-whitelist [comma-separated addresses] [add | delete]",
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
			depositArg, err := cmd.Flags().GetString(govcli.FlagDeposit)
			if err != nil {
				return err
			}
			deposit, err := sdk.ParseCoins(depositArg)
			if err != nil {
				return err
			}
			fmt.Println("deposit:", deposit)
			addOrNot, err := cmd.Flags().GetBool(isAdded)
			if err != nil {
				return err
			}
			addrs := strings.Split(strings.TrimSpace(args[0]), ",")

			proposal := types.UpdateDeploymentWhitelistProposal{
				Title:            proposalTitle,
				Description:      proposalDescr,
				DistributorAddrs: addrs,
				IsAdded:          addOrNot,
			}

			fmt.Println("proposal:", proposal)

			err = proposal.ValidateBasic()
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(viper.GetString(govcli.FlagDeposit))
			if err != nil {
				return err
			}

			fmt.Println("amount:", amount)

			msg := gov.NewMsgSubmitProposal(&proposal, amount, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	// proposal flags
	cmd.Flags().String(govcli.FlagTitle, "", "Title of proposal")
	cmd.Flags().String(govcli.FlagDescription, "", "Description of proposal")
	cmd.Flags().String(govcli.FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().Bool(isAdded, true, "True to add and false to delete")

	return cmd
}
