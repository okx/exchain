package cli

import (
	"fmt"
	"os"

	"github.com/okex/okexchain/x/common"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/okexchain/x/staking/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		client.PostCommands(
			GetCmdCreateValidator(cdc),
			GetCmdDestroyValidator(cdc),
			GetCmdEditValidator(cdc),
			GetCmdDeposit(cdc),
			GetCmdWithdraw(cdc),
			GetCmdAddShares(cdc),
		)...)

	stakingTxCmd.AddCommand(GetCmdProxy(cdc))

	return stakingTxCmd
}

// GetCmdCreateValidator gets the create validator command handler
func GetCmdCreateValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr, msg, err := BuildCreateValidatorMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(FsPk)
	//cmd.Flags().AddFlagSet(FsAmount)
	cmd.Flags().AddFlagSet(fsDescriptionCreate)
	//cmd.Flags().AddFlagSet(FsCommissionCreate)
	//cmd.Flags().AddFlagSet(FsMinSelfDelegation)

	cmd.Flags().String(FlagIP, "",
		fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", client.FlagGenerateOnly))
	cmd.Flags().String(FlagNodeID, "", "The node's ID")

	cmd.MarkFlagRequired(client.FlagFrom)
	cmd.MarkFlagRequired(FlagPubKey)
	cmd.MarkFlagRequired(FlagMoniker)

	return cmd
}

// GetCmdEditValidator gets the create edit validator command
// TODO: add full description
func GetCmdEditValidator(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-validator",
		Short: "edit an existing validator account",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddr := cliCtx.GetFromAddress()
			description := types.Description{
				Moniker:  viper.GetString(FlagMoniker),
				Identity: viper.GetString(FlagIdentity),
				Website:  viper.GetString(FlagWebsite),
				Details:  viper.GetString(FlagDetails),
			}

			// TODO: recover the msd modification later
			//var newMinSelfDelegation *sdk.Int
			//
			//minSelfDelegationString := viper.GetString(FlagMinSelfDelegation)
			//if minSelfDelegationString != "" {
			//	msb, ok := sdk.NewIntFromString(minSelfDelegationString)
			//	if !ok {
			//		return fmt.Errorf(types.ErrMinSelfDelegationInvalid(types.DefaultCodespace).Error())
			//	}
			//	/* required by okexchain */
			//	msb = msb.StandardizeAsc()
			//
			//	newMinSelfDelegation = &msb
			//}
			//
			//msg := types.NewMsgEditValidator(sdk.ValAddress(valAddr), description, newRate, newMinSelfDelegation)
			msg := types.NewMsgEditValidator(sdk.ValAddress(valAddr), description)

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsDescriptionEdit)
	//cmd.Flags().AddFlagSet(fsCommissionUpdate)

	return cmd
}

// CreateValidatorMsgHelpers returns the flagset, particular flags, and a description of defaults
// this is anticipated to be used with the gen-tx
func CreateValidatorMsgHelpers(ipDefault string) (fs *flag.FlagSet, nodeIDFlag, pubkeyFlag, amountFlag,
	defaultsDesc string) {
	fsCreateValidator := flag.NewFlagSet("", flag.ContinueOnError)
	fsCreateValidator.String(FlagIP, ipDefault, "The node's public IP")
	fsCreateValidator.String(FlagNodeID, "", "The node's NodeID")
	fsCreateValidator.String(FlagWebsite, "", "The validator's (optional) website")
	fsCreateValidator.String(FlagDetails, "", "The validator's (optional) details")
	fsCreateValidator.String(FlagIdentity, "", "The (optional) identity signature (ex. UPort or Keybase)")
	//fsCreateValidator.AddFlagSet(FsCommissionCreate)
	//fsCreateValidator.AddFlagSet(FsMinSelfDelegation)
	//fsCreateValidator.AddFlagSet(FsAmount)
	fsCreateValidator.AddFlagSet(FsPk)

	return fsCreateValidator, FlagNodeID, FlagPubKey, "", ""
}

// PrepareFlagsForTxCreateValidator prepares flags in config
func PrepareFlagsForTxCreateValidator(
	config *cfg.Config, nodeID, chainID string, valPubKey crypto.PubKey,
) {

	ip := viper.GetString(FlagIP)
	if ip == "" {
		fmt.Fprintf(os.Stderr, "couldn't retrieve an external IP; "+
			"the tx's memo field will be unset")
	}

	website := viper.GetString(FlagWebsite)
	details := viper.GetString(FlagDetails)
	identity := viper.GetString(FlagIdentity)

	viper.Set(client.FlagChainID, chainID)
	viper.Set(client.FlagFrom, viper.GetString(client.FlagName))
	viper.Set(FlagNodeID, nodeID)
	viper.Set(FlagIP, ip)
	viper.Set(FlagPubKey, sdk.MustBech32ifyConsPub(valPubKey))
	viper.Set(FlagMoniker, config.Moniker)
	viper.Set(FlagWebsite, website)
	viper.Set(FlagDetails, details)
	viper.Set(FlagIdentity, identity)

	if config.Moniker == "" {
		viper.Set(FlagMoniker, viper.GetString(client.FlagName))
	}
	//if viper.GetString(FlagAmount) == "" {
	//	viper.Set(FlagAmount, defaultAmount)
	//}
	//if viper.GetString(FlagCommissionRate) == "" {
	//	viper.Set(FlagCommissionRate, defaultCommissionRate)
	//}
	//if viper.GetString(FlagCommissionMaxRate) == "" {
	//	viper.Set(FlagCommissionMaxRate, defaultCommissionMaxRate)
	//}
	//if viper.GetString(FlagCommissionMaxChangeRate) == "" {
	//	viper.Set(FlagCommissionMaxChangeRate, defaultCommissionMaxChangeRate)
	//}
	// if viper.GetString(FlagMinSelfDelegation) == "" {
	//	viper.Set(FlagMinSelfDelegation, defaultMinSelfDelegation)
	//}
}

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func BuildCreateValidatorMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {

	valAddr := cliCtx.GetFromAddress()
	pkStr := viper.GetString(FlagPubKey)
	pk, err := sdk.GetConsPubKeyBech32(pkStr)
	if err != nil {
		return txBldr, nil, err
	}

	description := types.NewDescription(
		viper.GetString(FlagMoniker),
		viper.GetString(FlagIdentity),
		viper.GetString(FlagWebsite),
		viper.GetString(FlagDetails),
	)

	// get the initial validator min self delegation
	minSelfDelegation := sdk.DecCoin{
		Amount: types.DefaultMinSelfDelegation,
		Denom:  common.NativeToken,
	}

	msg := types.NewMsgCreateValidator(
		sdk.ValAddress(valAddr),
		pk,
		description,
		minSelfDelegation,
	)

	if viper.GetBool(client.FlagGenerateOnly) {
		ip := viper.GetString(FlagIP)
		nodeID := viper.GetString(FlagNodeID)
		if nodeID != "" && ip != "" {
			txBldr = txBldr.WithMemo(fmt.Sprintf("%s@%s:26656", nodeID, ip))
		}
	}

	return txBldr, msg, nil
}
