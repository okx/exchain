package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/spf13/cobra"
)

func extractNodeKey(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-node-key [format] [filename] ",
		Short: "extract current node key or from specificed file",
		RunE: func(cmd *cobra.Command, args []string) error {
			format := "hex"
			if len(args) >= 1 {
				format = args[0]
			}
			filename := ctx.Config.NodeKeyFile()
			if len(args) >= 2 {
				filename = args[1]
			}
			nodekey, err := p2p.LoadNodeKey(filename)
			if err != nil {
				return err
			}
			if format == "base64" {
				fmt.Printf("Node Public Key: %s\n", base64.StdEncoding.EncodeToString(nodekey.PubKey().Bytes()))

			} else {
				fmt.Printf("Node Public Key: %s\n", hex.EncodeToString(nodekey.PubKey().Bytes()))
			}

			return nil
		},
	}
	return cmd
}
