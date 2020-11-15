package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/okexchain/x/staking/types"
	"github.com/spf13/cobra"
)

// GetCmdDestroyValidator gets command for destroying a validator and unbonding the min-self-delegation
func GetCmdDestroyValidator(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "destroy-validator [flags]",
		Args:  cobra.NoArgs,
		Short: "deregister the validator from the OKExChain and unbond the min self delegation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deregister the validator from the OKExChain and unbond the min self delegation.

Example:
$ %s tx staking destroy-validator --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgDestroyValidator(delAddr)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})

		},
	}
}

// GetCmdDeposit gets command for deposit
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "deposit [amount]",
		Args: cobra.ExactArgs(1),
		Short: fmt.Sprintf("deposit an amount of %s to delegator account; deposited %s in delegator account is a prerequisite for adding shares",
			sdk.DefaultBondDenom, sdk.DefaultBondDenom),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Deposit an amount of %s to delegator account. Deposited %s in delegator account is a prerequisite for adding shares.

Example:
$ %s tx staking deposit 1000%s --from mykey
`,
				sdk.DefaultBondDenom, sdk.DefaultBondDenom, version.ClientName, sdk.DefaultBondDenom,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}

			delAddr := cliCtx.GetFromAddress()
			msg := types.NewMsgDeposit(delAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdWithdraw gets command for withdrawing the deposit
func GetCmdWithdraw(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [amount]",
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("withdraw an amount of %s and the corresponding shares from all validators", sdk.DefaultBondDenom),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw an amount of %s and the corresponding shares from all validators.

Example:
$ %s tx staking withdraw 1%s
`,
				sdk.DefaultBondDenom, version.ClientName, sdk.DefaultBondDenom,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := sdk.ParseDecCoin(args[0])
			if err != nil {
				return err
			}

			delAddr := cliCtx.GetFromAddress()
			msg := types.NewMsgWithdraw(delAddr, amount)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

// GetCmdAddShares gets command for multi voting
func GetCmdAddShares(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-shares [validator-addr1, validator-addr2, validator-addr3, ... validator-addrN] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("add shares to one or more validators by all deposited %s", sdk.DefaultBondDenom),
		Long: strings.TrimSpace(
			fmt.Sprintf("Add shares to one or more validators by all deposited %s.\n\nExample:\n$ %s tx staking add-shares "+
				"okexchainvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg,"+
				"okexchainvaloper1svzxp4ts5le2s4zugx34ajt6shz2hg42dnwst5,"+
				"okexchainvaloper10q0rk5qnyag7wfvvt7rtphlw589m7frshchly8,"+
				"okexchainvaloper1g7znsf24w4jc3xfca88pq9kmlyjdare6tr3mk6 --from mykey\n",
				sdk.DefaultBondDenom, version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()
			valAddrs, err := getValsSet(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddShares(delAddr, valAddrs)
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
		flags.PostCommands(
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgRegProxy(delAddr, true)
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgRegProxy(delAddr, false)
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
$ %s tx staking proxy bind okexchain1hw4r48aww06ldrfeuq2v438ujnl6alsz0685a0 --from mykey
`,
				version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()

			proxyAddress, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}
			msg := types.NewMsgBindProxy(delAddr, proxyAddress)
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delAddr := cliCtx.GetFromAddress()

			msg := types.NewMsgUnbindProxy(delAddr)
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
			return nil, fmt.Errorf("invalid target validator address: %s", addrs[i])
		}
	}
	return
}
