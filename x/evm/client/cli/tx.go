package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	emint "github.com/okex/okexchain/app/types"
	evmutils "github.com/okex/okexchain/x/evm/client/utils"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/okex/okexchain/x/gov"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// GetTxCmd defines the CLI commands regarding evm module transactions
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	evmTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "EVM transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	evmTxCmd.AddCommand(flags.PostCommands(
		GetCmdSendTx(cdc),
		GetCmdGenCreateTx(cdc),
	)...)

	return evmTxCmd
}

func cosmosAddressFromArg(addr string) (sdk.AccAddress, error) {
	if strings.HasPrefix(addr, sdk.GetConfig().GetBech32AccountAddrPrefix()) {
		// Check to see if address is Cosmos bech32 formatted
		toAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid bech32 formatted address")
		}
		return toAddr, nil
	}

	// Strip 0x prefix if exists
	addr = strings.TrimPrefix(addr, "0x")

	return sdk.AccAddressFromHex(addr)
}

// GetCmdSendTx generates an Ethermint transaction (excludes create operations)
func GetCmdSendTx(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [to_address] [amount (in aphotons)] [<data>]",
		Short: "send transaction to address (call operations included)",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			toAddr, err := cosmosAddressFromArg(args[0])
			if err != nil {
				return errors.Wrap(err, "must provide a valid Bech32 address for to_address")
			}

			// Ambiguously decode amount from any base
			amount, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return err
			}

			var data []byte
			if len(args) > 2 {
				payload := args[2]
				if !strings.HasPrefix(payload, "0x") {
					payload = "0x" + payload
				}

				data, err = hexutil.Decode(payload)
				if err != nil {
					return err
				}
			}

			from := cliCtx.GetFromAddress()

			_, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
			if err != nil {
				return errors.Wrap(err, "Could not retrieve account sequence")
			}

			// TODO: Potentially allow overriding of gas price and gas limit
			msg := types.NewMsgEthermint(seq, &toAddr, sdk.NewIntFromBigInt(amount.Int), txBldr.Gas(),
				sdk.NewInt(emint.DefaultGasPrice), data, from)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdGenCreateTx generates an Ethermint transaction (excludes create operations)
func GetCmdGenCreateTx(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [contract bytecode] [<amount (in aphotons)>]",
		Short: "create contract through the evm using compiled bytecode",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			payload := args[0]
			if !strings.HasPrefix(payload, "0x") {
				payload = "0x" + payload
			}

			data, err := hexutil.Decode(payload)
			if err != nil {
				return err
			}

			amount := sdk.ZeroDec()
			if len(args) > 1 {
				// Ambiguously decode amount from any base
				amount, err = sdk.NewDecFromStr(args[1])
				if err != nil {
					return errors.Wrap(err, "invalid amount")
				}
			}

			from := cliCtx.GetFromAddress()

			_, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
			if err != nil {
				return errors.Wrap(err, "Could not retrieve account sequence")
			}

			// TODO: Potentially allow overriding of gas price and gas limit
			msg := types.NewMsgEthermint(seq, nil, sdk.NewIntFromBigInt(amount.Int), txBldr.Gas(),
				sdk.NewInt(emint.DefaultGasPrice), data, from)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			if err = authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}); err != nil {
				return err
			}

			contractAddr := ethcrypto.CreateAddress(common.BytesToAddress(from.Bytes()), seq)
			fmt.Printf(
				"Contract will be deployed to: \nHex: %s\nCosmos Address: %s\n",
				contractAddr.Hex(),
				sdk.AccAddress(contractAddr.Bytes()),
			)
			return nil
		},
	}
}

// GetCmdManageContractDeploymentWhitelistProposal implements a command handler for submitting a manage contract deployment
// whitelist proposal transaction
func GetCmdManageContractDeploymentWhitelistProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "manage-contract-deployment-whitelist [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a manage contract deployment whitelist proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a manage contract deployment whitelist proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal manage-contract-deployment-whitelist <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
 "title": "manage contract proposal whitelist with a deployer address",
 "description": "add a deployer address into the whitelist",
 "deployer_address": "okexchain1hw4r48aww06ldrfeuq2v438ujnl6alsz0685a0",
 "is_added": true,
 "deposit": [
   {
     "denom": "%s",
     "amount": "100"
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

			deployerAddr, err := sdk.AccAddressFromBech32(proposal.DeployerAddr)
			if err != nil {
				return err
			}

			content := types.NewManageContractDeploymentWhitelistProposal(
				proposal.Title,
				proposal.Description,
				deployerAddr,
				proposal.IsAdded,
			)

			msg := gov.NewMsgSubmitProposal(content, proposal.Deposit, cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
