package client

import (
	"fmt"
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const (
	okexPrefix = "okexchain"
	exPrefix   = "ex"
	zxPrefix   = "0x"
)

// AddrCommands registers a sub-tree of commands to interact with oec address
func AddrCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr",
		Short: "opreate all kind of address in the OEC network",
		Long: ` Address is a identification for join in the OEC network.

	The address in OEC network begins with "okexchain","ex" or "0x"`,
	}
	cmd.AddCommand(convertCommand())
	return cmd

}

func convertCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "convert [sourceAddr]",
		Short: "convert source address to all kind of address in the OEC network",
		Long: `sourceAddr must be begin with "okexchain","ex" or "0x".
	
	When input one of these address, we will convert to the other kinds.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrList := make(map[string]string)
			targetPrefix := []string{okexPrefix, exPrefix, zxPrefix}
			srcAddr := args[0]

			// read previous config
			config := sdk.GetConfig()
			pfxAddrPre := config.GetBech32AccountAddrPrefix()
			pfxPubPre := config.GetBech32AccountPubPrefix()
			config.Unseal()
			defer config.RecoverPrefixForAcc(pfxAddrPre, pfxPubPre)

			//prefix is "okexchain","ex" or "0x"
			var srcPrefix string
			switch {
			case strings.HasPrefix(srcAddr, okexPrefix):
				srcPrefix = okexPrefix

			case strings.HasPrefix(srcAddr, exPrefix):
				srcPrefix = exPrefix

			case strings.HasPrefix(srcAddr, zxPrefix):
				srcPrefix = zxPrefix

			default:
				return fmt.Errorf("unsupported prefix to convert")
			}

			//source address parse to account address
			addrList[srcPrefix] = srcAddr
			config.SetBech32PrefixForAccount(srcPrefix, fmt.Sprintf("%s%s", srcPrefix, sdk.PrefixPublic))
			accAddr, err := sdk.AccAddressFromBech32(srcAddr)
			if err != nil {
				fmt.Printf("Parse bech32 address error: %s", err)
				return err
			}

			// fill other kinds of prefix address out
			for _, pfx := range targetPrefix {
				if _, ok := addrList[pfx]; !ok {
					config.SetBech32PrefixForAccount(pfx, fmt.Sprintf("%s%s", pfx, sdk.PrefixPublic))
					addrList[pfx] = accAddr.String()
				}
			}

			//show all kinds of prefix address out
			for _, pfx := range targetPrefix {
				fmt.Printf("prefix: %s, its complete address is %s \n", pfx, addrList[pfx])
			}

			return nil
		},
	}
}
