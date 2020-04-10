package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/okchain/x/staking/types"
	"github.com/spf13/cobra"
)

// GetCmdDestroyValidator gets command for destroying a validator and unbonding the min-self-delegation
func GetCmdDestroyValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy-validator [flags]",
		Args:  cobra.NoArgs,
		Short: "deregister the validator from the OKChain and unbond the min self delegation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deregister the validator from the OKChain and unbond the min self delegation.

Example:
$ %s tx staking destroy-validator --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			voterAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgDestroyValidator(voterAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

		},
	}
}

// GetCmdDelegate gets command for delegating
func GetCmdDelegate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "delegate an amount of okt.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Delegate an amount of okt.

Example:
$ %s tx staking delegate 1000okt --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}

			delAddr := cliCtx.GetFromAddress()
			msg := types.NewMsgDelegate(delAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdUndelegate gets command for undelegating
func GetCmdUndelegate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unbond [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "unbond shares and withdraw the same amount of votes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbond shares and withdraw the same amount of votes.

Example:
$ %s tx staking unbond 1okt
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}

			delAddr := cliCtx.GetFromAddress()
			msg := types.NewMsgUndelegate(delAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdVote gets command for multi voting
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [validator-addr1, validator-addr2, validator-addr3, ... validator-addrN] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "vote on validators",
		Long: strings.TrimSpace(
			fmt.Sprintf("Vote on one or more validator(s).\n\nExample:\n$ %s tx staking vote "+
				"okchainvaloper1alq9na49n9yycysh889rl90g9nhe58lcs50wu5,"+
				"okchainvaloper1svzxp4ts5le2s4zugx34ajt6shz2hg42a3gl7g,"+
				"okchainvaloper10q0rk5qnyag7wfvvt7rtphlw589m7frs863s3m,"+
				"okchainvaloper1g7znsf24w4jc3xfca88pq9kmlyjdare6mph5rx --from mykey\n",
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			voterAddr := cliCtx.GetFromAddress()
			valAddrs, err := getValsSet(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgVote(voterAddr, valAddrs)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

		},
	}
}

// GetCmdProxy gets subcommands for proxy voting
func GetCmdProxy(cdc *codec.Codec) *cobra.Command {

	proxyCmd := &cobra.Command{
		Use:   "proxy",
		Short: "proxy subcommands",
	}

	proxyCmd.AddCommand(
		client.PostCommands(
			GetCmdRegProxy(cdc),
			GetCmdUnregProxy(cdc),
			GetCmdBindProxy(cdc),
			GetCmdUnbindProxy(cdc),
		)...)

	return proxyCmd
}

// GetCmdRegProxy gets command for proxy registering
func GetCmdRegProxy(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reg [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "become a proxy by registration",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Become a proxy by registration.

Example:
$ %s tx staking proxy reg --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			voterAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgRegProxy(voterAddr, true)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUnregProxy gets command for proxy unregistering
func GetCmdUnregProxy(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unreg [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "unregister the proxy identity",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unregister the proxy identity.

Example:
$ %s tx staking proxy unreg --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			voterAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgRegProxy(voterAddr, false)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBindProxy gets command for binding proxy
func GetCmdBindProxy(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bind [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "bind proxy relationship",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Bind proxy relationship.

Example:
$ %s tx staking proxy bind okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			voterAddr := cliCtx.GetFromAddress()

			proxyAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}
			msg := types.NewMsgBindProxy(voterAddr, proxyAddress)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdUnbindProxy gets command for unbinding proxy
func GetCmdUnbindProxy(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unbind [flags]",
		Args:  cobra.ExactArgs(0),
		Short: "unbind proxy relationship",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unbind proxy relationship.

Example:
$ %s tx staking proxy unbind --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			voterAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgUnbindProxy(voterAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// getValsSet gets validator set from client args
func getValsSet(address string) (valAddrs []sdk.ValAddress, err error) {
	addrs := strings.Split(strings.TrimSpace(address), ",")
	lenVals := len(addrs)
	valAddrs = make([]sdk.ValAddress, lenVals)
	for i := 0; i < lenVals; i++ {
		valAddrs[i], err = sdk.ValAddressFromBech32(addrs[i])
		if err != nil {
			return nil, err
		}
	}
	return
}
