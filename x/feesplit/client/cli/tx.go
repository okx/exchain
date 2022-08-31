package cli

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/x/feesplit/types"
)

// GetTxCmd returns a root CLI command handler for certain modules/feesplit
// transaction commands.
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "feesplit subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(flags.PostCommands(
		GetRegisterFeeSplit(cdc),
		GetCancelFeeSplit(cdc),
		GetUpdateFeeSplit(cdc),
	)...)
	return cmd
}

// GetRegisterFeeSplit returns a CLI command handler for registering a
// contract for fee distribution
func GetRegisterFeeSplit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [contract_hex] [nonces] [withdraw_bech32]",
		Short: "Register a contract for fee distribution. **NOTE** Please ensure, that the deployer of the contract (or the factory that deployes the contract) is an account that is owned by your project, to avoid that an individual deployer who leaves your project becomes malicious.",
		Long:  "Register a contract for fee distribution.\nOnly the contract deployer can register a contract.\nProvide the account nonce(s) used to derive the contract address. E.g.: you have an account nonce of 4 when you send a deployment transaction for a contract A; you use this contract as a factory, to create another contract B. If you register A, the nonces value is \"4\". If you register B, the nonces value is \"4,1\" (B is the first contract created by A). \nThe withdraw address defaults to the deployer address if not provided.",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var withdraw string
			deployer := cliCtx.GetFromAddress()

			contract := args[0]
			if err := types.ValidateNonZeroAddress(contract); err != nil {
				return fmt.Errorf("invalid contract hex address %w", err)
			}

			var nonces []uint64
			if err := json.Unmarshal([]byte("["+args[1]+"]"), &nonces); err != nil {
				return fmt.Errorf("invalid nonces %w", err)
			}

			if len(args) == 3 {
				withdraw = args[2]
				if _, err := sdk.AccAddressFromBech32(withdraw); err != nil {
					return fmt.Errorf("invalid withdraw bech32 address %w", err)
				}
			}

			// If withdraw address is the same as contract deployer, remove the
			// field for avoiding storage bloat
			if deployer.String() == withdraw {
				withdraw = ""
			}

			msg := &types.MsgRegisterFeeSplit{
				ContractAddress:   contract,
				DeployerAddress:   deployer.String(),
				WithdrawerAddress: withdraw,
				Nonces:            nonces,
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetCancelFeeSplit returns a CLI command handler for canceling a
// contract for fee distribution
func GetCancelFeeSplit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel [contract_hex]",
		Short: "Cancel a contract from fee distribution",
		Long:  "Cancel a contract from fee distribution. The deployer will no longer receive fees from users interacting with the contract. \nOnly the contract deployer can cancel a contract.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			deployer := cliCtx.GetFromAddress()

			contract := args[0]
			if err := types.ValidateNonZeroAddress(contract); err != nil {
				return fmt.Errorf("invalid contract hex address %w", err)
			}

			msg := &types.MsgCancelFeeSplit{
				ContractAddress: contract,
				DeployerAddress: deployer.String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// GetUpdateFeeSplit returns a CLI command handler for updating the withdraw
// address of a contract for fee distribution
func GetUpdateFeeSplit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [contract_hex] [withdraw_bech32]",
		Short: "Update withdraw address for a contract registered for fee distribution.",
		Long:  "Update withdraw address for a contract registered for fee distribution. \nOnly the contract deployer can update the withdraw address.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			deployer := cliCtx.GetFromAddress()

			contract := args[0]
			if err := types.ValidateNonZeroAddress(contract); err != nil {
				return fmt.Errorf("invalid contract hex address %w", err)
			}

			withdraw := args[1]
			if _, err := sdk.AccAddressFromBech32(withdraw); err != nil {
				return fmt.Errorf("invalid withdraw bech32 address %w", err)
			}

			msg := &types.MsgUpdateFeeSplit{
				ContractAddress:   contract,
				DeployerAddress:   deployer.String(),
				WithdrawerAddress: withdraw,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	
	return cmd
}
