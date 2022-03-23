package client

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const (
	okexPrefix = "okexchain"
	exPrefix   = "ex"
	rawPrefix  = "0x"
)

type accAddrToPrefixFunc func(sdk.AccAddress, string) string

// AddrCommands registers a sub-tree of commands to interact with okc address
func AddrCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr",
		Short: "opreate all kind of address in the OKC network",
		Long: ` Address is a identification for join in the OKC network.

	The address in OKC network begins with "okexchain","ex" or "0x"`,
	}
	cmd.AddCommand(convertCommand())
	return cmd

}

func convertCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "convert [sourceAddr]",
		Short: "convert source address to all kind of address in the OKC network",
		Long: `sourceAddr must be begin with "okexchain","ex" or "0x".
	
	When input one of these address, we will convert to the other kinds.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			addrList := make(map[string]string)
			targetPrefix := []string{okexPrefix, exPrefix, rawPrefix}
			srcAddr := args[0]

			// register func to encode account address to prefix address.
			toPrefixFunc := map[string]accAddrToPrefixFunc{
				okexPrefix: bech32FromAccAddr,
				exPrefix:   bech32FromAccAddr,
				rawPrefix:  hexFromAccAddr,
			}

			// prefix is "okexchain","ex" or "0x"
			// convert srcAddr to accAddr
			var accAddr sdk.AccAddress
			var err error
			switch {
			case strings.HasPrefix(srcAddr, okexPrefix):
				//source address parse to account address
				addrList[okexPrefix] = srcAddr
				accAddr, err = bech32ToAccAddr(okexPrefix, srcAddr)

			case strings.HasPrefix(srcAddr, exPrefix):
				//source address parse to account address
				addrList[exPrefix] = srcAddr
				accAddr, err = bech32ToAccAddr(exPrefix, srcAddr)

			case strings.HasPrefix(srcAddr, rawPrefix):
				addrList[rawPrefix] = srcAddr
				accAddr, err = hexToAccAddr(rawPrefix, srcAddr)

			default:
				return fmt.Errorf("unsupported prefix to convert")
			}

			// check account address
			if err != nil {
				fmt.Printf("Parse bech32 address error: %s", err)
				return err
			}

			// fill other kinds of prefix address out
			for _, pfx := range targetPrefix {
				if _, ok := addrList[pfx]; !ok {
					addrList[pfx] = toPrefixFunc[pfx](accAddr, pfx)
				}
			}

			//show all kinds of prefix address out
			for _, pfx := range targetPrefix {
				addrType := "Bech32"
				if pfx == "0x" {
					addrType = "Hex"
				}
				fmt.Printf("%s format with prefix <%s>: %5s\n", addrType, pfx, addrList[pfx])
			}

			return nil
		},
	}
}

// bech32ToAccAddr convert a hex string which begins with 'prefix' to an account address
func bech32ToAccAddr(prefix string, srcAddr string) (sdk.AccAddress, error) {
	return sdk.AccAddressFromBech32ByPrefix(srcAddr, prefix)
}

// bech32FromAccAddr create a hex string which begins with 'prefix' to from account address
func bech32FromAccAddr(accAddr sdk.AccAddress, prefix string) string {
	return accAddr.Bech32String(prefix)
}

// hexToAccAddr convert a hex string to an account address
func hexToAccAddr(prefix string, srcAddr string) (sdk.AccAddress, error) {
	srcAddr = strings.TrimPrefix(srcAddr, prefix)
	return sdk.AccAddressFromHex(srcAddr)
}

// hexFromAccAddr create a hex string from an account address
func hexFromAccAddr(accAddr sdk.AccAddress, prefix string) string {
	return prefix + hex.EncodeToString(accAddr.Bytes())
}
