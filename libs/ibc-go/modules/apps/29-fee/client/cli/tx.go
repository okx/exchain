package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	"github.com/spf13/cobra"
)

const (
	flagRecvFee    = "recv-fee"
	flagAckFee     = "ack-fee"
	flagTimeoutFee = "timeout-fee"
)

// NewRegisterPayeeCmd returns the command to create a MsgRegisterPayee
func NewRegisterPayeeCmd(codecProxy *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "register-payee [port-id] [channel-id] [relayer] [payee] ",
		Short:   "Register a payee on a given channel.",
		Long:    strings.TrimSpace(`Register a payee address on a given channel.`),
		Example: fmt.Sprintf("%s tx ibc-fee register-payee transfer channel-0 cosmos1rsp837a4kvtgp2m4uqzdge0zzu6efqgucm0qdh cosmos153lf4zntqt33a4v0sm5cytrxyqn78q7kz8j8x5", version.ServerName),
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(codecProxy.GetCdc()))
			clientCtx := context.NewCLIContext().WithInterfaceRegistry(reg)

			msg := types.NewMsgRegisterPayee(args[0], args[1], args[2], args[3])

			return utils.GenerateOrBroadcastMsgs(clientCtx, txBldr, []sdk.Msg{msg})
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewRegisterCounterpartyPayeeCmd returns the command to create a MsgRegisterCounterpartyPayee
func NewRegisterCounterpartyPayeeCmd(codecProxy *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "register-counterparty-payee [port-id] [channel-id] [relayer] [counterparty-payee] ",
		Short:   "Register a counterparty payee address on a given channel.",
		Long:    strings.TrimSpace(`Register a counterparty payee address on a given channel.`),
		Example: fmt.Sprintf("%s tx ibc-fee register-counterparty-payee transfer channel-0 cosmos1rsp837a4kvtgp2m4uqzdge0zzu6efqgucm0qdh osmo1v5y0tz01llxzf4c2afml8s3awue0ymju22wxx2", version.ServerName),
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(codecProxy.GetCdc()))
			clientCtx := context.NewCLIContext().WithInterfaceRegistry(reg)

			msg := types.NewMsgRegisterCounterpartyPayee(args[0], args[1], args[2], args[3])

			return utils.GenerateOrBroadcastMsgs(clientCtx, txBldr, []sdk.Msg{msg})
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewPayPacketFeeAsyncTxCmd returns the command to create a MsgPayPacketFeeAsync
func NewPayPacketFeeAsyncTxCmd(codecProxy *codec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pay-packet-fee [src-port] [src-channel] [sequence]",
		Short:   "Pay a fee to incentivize an existing IBC packet",
		Long:    strings.TrimSpace(`Pay a fee to incentivize an existing IBC packet.`),
		Example: fmt.Sprintf("%s tx ibc-fee pay-packet-fee transfer channel-0 1 --recv-fee 10stake --ack-fee 10stake --timeout-fee 10stake", version.ServerName),
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(codecProxy.GetCdc()))
			clientCtx := context.NewCLIContext().WithInterfaceRegistry(reg)

			// NOTE: specifying non-nil relayers is currently unsupported
			var relayers []string

			sender := clientCtx.GetFromAddress().String()
			seq, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			packetID := channeltypes.NewPacketId(args[0], args[1], seq)

			recvFeeStr, err := cmd.Flags().GetString(flagRecvFee)
			if err != nil {
				return err
			}

			recvFee, err := sdk.ParseCoinsNormalized(recvFeeStr)
			if err != nil {
				return err
			}

			ackFeeStr, err := cmd.Flags().GetString(flagAckFee)
			if err != nil {
				return err
			}

			ackFee, err := sdk.ParseCoinsNormalized(ackFeeStr)
			if err != nil {
				return err
			}

			timeoutFeeStr, err := cmd.Flags().GetString(flagTimeoutFee)
			if err != nil {
				return err
			}

			timeoutFee, err := sdk.ParseCoinsNormalized(timeoutFeeStr)
			if err != nil {
				return err
			}

			fee := types.Fee{
				RecvFee:    utils.CliConvertCoinToCoinAdapters(recvFee),
				AckFee:     utils.CliConvertCoinToCoinAdapters(ackFee),
				TimeoutFee: utils.CliConvertCoinToCoinAdapters(timeoutFee),
			}
			fmt.Println(fee.String())

			packetFee := types.NewPacketFee(fee, sender, relayers)
			msg := types.NewMsgPayPacketFeeAsync(packetID, packetFee)

			return utils.GenerateOrBroadcastMsgs(clientCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagRecvFee, "", "Fee paid to a relayer for relaying a packet receive.")
	cmd.Flags().String(flagAckFee, "", "Fee paid to a relayer for relaying a packet acknowledgement.")
	cmd.Flags().String(flagTimeoutFee, "", "Fee paid to a relayer for relaying a packet timeout.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
