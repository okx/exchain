package cli

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	govTypes "github.com/okex/okchain/x/gov/types"
	tokenUtils "github.com/okex/okchain/x/token/client/utils"
	"github.com/okex/okchain/x/token/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	From          = "from"
	To            = "to"
	Sign          = "sign"
	Amount        = "amount"
	TotalSupply   = "total-supply"
	Symbol        = "symbol"
	WholeName     = "whole-name"
	TokenDesc     = "desc"
	Mintable      = "mintable"
	Transfers     = "transfers"
	TransfersFile = "transfers-file"
)

const (
	TokenDescLenLimit = 256
)

var (
	errSymbolNotValid         = errors.New("symbol not valid")
	errTotalSupplyNotValid    = errors.New("total-supply not valid")
	errFromNotValid           = errors.New("from not valid")
	errAmountNotValid         = errors.New("amount not valid")
	errTokenDescNotValid      = errors.New("token-desc not valid")
	errTokenWholeNameNotValid = errors.New("token whole name not valid")
	errMintableNotValid       = errors.New("mintable not valid")
	errTransfersNotValid      = errors.New("transfers not valid")
	errTransfersFileNotValid  = errors.New("transfers file not valid")
	errSign                   = errors.New("sign not succeed")
	errParam                  = errors.New("can't get token desc or whole name")
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	distTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	distTxCmd.AddCommand(client.PostCommands(
		getCmdTokenIssue(cdc),
		getCmdTokenBurn(cdc),
		getCmdTokenMint(cdc),
		getCmdTokenMultiSend(cdc),
		getCmdTransferOwnership(cdc),
		getMultiSignsCmd(cdc),
		getCmdTokenEdit(cdc),
		getCmdTokenActive(cdc),
	)...)

	return distTxCmd
}

// getCmdTokenIssue is the CLI command for sending a IssueToken transaction
func getCmdTokenIssue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "issue a token",
		//Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}

			flags := cmd.Flags()

			// params check
			originalSymbol, err := flags.GetString(Symbol)
			originalSymbol = strings.ToLower(originalSymbol)
			if err != nil {
				return errSymbolNotValid
			}
			_, err = flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}
			totalSupply, err := flags.GetString(TotalSupply)
			if err != nil {
				return errTotalSupplyNotValid
			}
			tokenDesc, err := flags.GetString(TokenDesc)
			if err != nil || len(tokenDesc) > TokenDescLenLimit {
				return errTokenDescNotValid
			}

			wholeName, err := flags.GetString(WholeName)
			if err != nil {
				return errTokenWholeNameNotValid
			}
			// check wholeName
			var isValid bool
			wholeName, isValid = types.WholeNameCheck(wholeName)
			if !isValid {
				return errTokenWholeNameNotValid
			}

			mintable, err := flags.GetBool(Mintable)
			if err != nil {
				return errMintableNotValid
			}

			var symbol string

			// totalSupply int64 ,coins bigint
			msg := types.NewMsgTokenIssue(tokenDesc, symbol, originalSymbol, wholeName, totalSupply, cliCtx.FromAddress, mintable)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(Symbol, "s", "", "symbol of the new token")
	cmd.Flags().StringP(WholeName, "w", "", "whole name of the new token")
	cmd.Flags().String(TokenDesc, "", "describe of the token")
	cmd.Flags().StringP(TotalSupply, "n", "0", "total supply of the new token")
	cmd.Flags().Bool(Mintable, false, "whether the token can be minted")

	return cmd
}

// getCmdTokenBurn is the CLI command for sending a BurnToken transaction
func getCmdTokenBurn(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [amount]",
		Short: "burn some amount of token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}

			// params check
			flags := cmd.Flags()

			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}

			_, err = flags.GetString(From)
			if err != nil {
				return err
			}

			msg := types.NewMsgTokenBurn(amount, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

// getCmdTokenMint is the CLI command for sending a MintToken transaction
func getCmdTokenMint(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [amount]",
		Short: "mint tokens for an existing token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}
			flags := cmd.Flags()

			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}

			_, err = flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}

			msg := types.NewMsgTokenMint(amount, cliCtx.FromAddress)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	return cmd
}

// getCmdTokenMultiSend is the CLI command for sending a MultiSend transaction
func getCmdTokenMultiSend(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-send",
		Short: "Create and sign a multi send tx",
		//Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}
			flags := cmd.Flags()

			_, err := flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}
			transferStr, err := flags.GetString(Transfers)
			if err != nil {
				return errTransfersNotValid
			}

			transfersFile, err := flags.GetString(TransfersFile)
			if err != nil {
				return errTransfersFileNotValid
			}

			var transfers []types.TransferUnit
			if transferStr != "" {
				transfers, err = types.StrToTransfers(transferStr)
				if err != nil {
					return err
				}
			}

			if transfersFile != "" {
				transferBytes, err := ioutil.ReadFile(transfersFile)
				if err != nil {
					return err
				}
				transferStr = string(transferBytes)
				//return errors.New(transferStr)
				transfers, err = types.StrToTransfers(transferStr)
				if err != nil {
					return err
				}
			}

			for _, transfer := range transfers {
				if transfer.To.Equals(cliCtx.GetFromAddress()) {
					return errors.New(transfer.To.String())
					//return errors.New("can not transfer coins to yourself")
				}
			}

			msg := types.NewMsgMultiSend(cliCtx.FromAddress, transfers)

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(Transfers, "", `Transfers details, format: [{"to": "addr", "amount": "1okt,2btc"}, ...]`)
	cmd.Flags().String(TransfersFile, "", "File of transfers details, if transfers-file is not empty, --transfers will be ignore")
	//cmd.MarkFlagRequired(Amount)
	return cmd
}

// getCmdTransferOwnership is the CLI command for sending a ChangeOwner transaction
func getCmdTransferOwnership(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership",
		Short: "change the owner of the token",
		//Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}
			flags := cmd.Flags()

			symbol, err := flags.GetString(Symbol)
			if err != nil {
				return errSymbolNotValid
			}
			_, err = flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}
			to, err := flags.GetString(To)
			if err != nil {
				return errAmountNotValid
			}

			from := cliCtx.GetFromAddress()

			toBytes, err := sdk.AccAddressFromBech32(to)
			if err != nil {
				return errFromNotValid
			}

			msg := types.NewMsgTransferOwnership(from, toBytes, symbol)

			return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP("symbol", "s", "", "symbol of the token to be transferred")
	cmd.Flags().String("to", "", "the user to be transferred")
	return cmd
}

// nolint
func getMultiSignsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisigns",
		Short: "append signature to the chown unsignedtx file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			stdTx, err := utils.ReadStdTxFromFile(cdc, args[0])
			if err != nil {
				return err
			}

			if len(stdTx.Msgs) == 0 {
				return err
			}

			msg, ok := stdTx.Msgs[0].(types.MsgTransferOwnership)
			if !ok {
				// todo
				return errSign
			}

			flags := cmd.Flags()
			_, err = flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}

			//
			passphrase, err := keys.GetPassphrase(cliCtx.GetFromName())
			if err != nil {
				return err
			}
			ToSignature, _, err := txBldr.Keybase().Sign(cliCtx.GetFromName(), passphrase, msg.GetSignBytes())
			if err != nil {
				return errSign
			}
			info, err := txBldr.Keybase().Get(cliCtx.GetFromName())
			if err != nil {
				return err
			}
			stdSig := auth.StdSignature{
				PubKey:    info.GetPubKey(),
				Signature: ToSignature,
			}
			msg.ToSignature = stdSig

			return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg})

		},
	}
	return cmd
}

// SendTxCmd will create a transaction to send and sign
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_key_or_address] [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[1])
			}

			coins, err := sdk.ParseDecCoins(args[2])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgTokenSend(cliCtx.GetFromAddress(), to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// getCmdTokenEdit is the CLI command for sending a TokenEdit transaction
func getCmdTokenEdit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit a token's whole name and desc",
		//Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := authTypes.NewAccountRetriever(cliCtx).EnsureExists(cliCtx.FromAddress); err != nil {
				return err
			}

			flags := cmd.Flags()

			// params check
			symbol, err := flags.GetString(Symbol)
			if err != nil {
				return errSymbolNotValid
			}
			_, err = flags.GetString(From)
			if err != nil {
				return errFromNotValid
			}

			var isDescEdit, isWholeNameEdit bool
			var tokenDesc, wholeName string
			dcEditFlag := flags.Lookup(TokenDesc)
			if dcEditFlag != nil && dcEditFlag.Changed {
				isDescEdit = true
				tokenDesc, err = flags.GetString(TokenDesc)
				if err != nil || len(tokenDesc) > TokenDescLenLimit {
					return errTokenDescNotValid
				}
			}
			wnEditFlag := flags.Lookup(WholeName)
			if wnEditFlag != nil && wnEditFlag.Changed {
				isWholeNameEdit = true
				wholeName, err = flags.GetString(WholeName)
				if err != nil {
					return errTokenWholeNameNotValid
				}
				// check wholeName
				var isValid bool
				wholeName, isValid = types.WholeNameCheck(wholeName)
				if !isValid {
					return errTokenWholeNameNotValid
				}
			}
			if !isWholeNameEdit && !isDescEdit {
				return errParam
			}

			msg := types.NewMsgTokenModify(symbol, tokenDesc, wholeName, isDescEdit, isWholeNameEdit, cliCtx.FromAddress)
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(Symbol, "s", "", "symbol of the token")
	cmd.Flags().StringP(WholeName, "w", "", "whole name of the token")
	cmd.Flags().String(TokenDesc, "", "description of the token")

	return cmd
}

// GetCmdSubmitProposal implements a command handler for submitting a parameter change proposal transaction
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certified-token [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a token without suffix proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a certified token proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal certified-token <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Issue a token without suffix",
  "description": "Issue a token without suffix: btc",
  "token": {
     "description": "Bitcoin in testnet，1:1 anchoring with Bitcoin",
     "symbol": "tbtc",
     "whole_name": "Testnet Bitcoin",
     "total_supply": "21000000",
     "owner": "okchain170ydluqy9qnza6t7mt5rqtqw5p3esdrz00hvm2",
     "mintable": false
   },
  "deposit": [
    {
      "denom": %s,
      "amount": "10000"
    }
  ]
}
`, version.ClientName, sdk.DefaultBondDenom)),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := tokenUtils.ParseCertifiedTokenProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			from := cliCtx.GetFromAddress()
			content := types.NewCertifiedTokenProposal(
				proposal.Title,
				proposal.Description,
				proposal.Token,
			)

			msg := govTypes.NewMsgSubmitProposal(content, proposal.Deposit, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

// getCmdActive implements creating a new vote command.
func getCmdTokenActive(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "active [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Activate a certified token",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Activate a certified token. You can
find the proposal-id by running "%s query gov proposals".


Example:
$ %s tx token active 1 --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get voting address
			from := cliCtx.GetFromAddress()

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			msg := types.NewMsgTokenActive(proposalID, from)
			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
