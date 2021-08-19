package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

func dataCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "modify data or query data in database",
	}

	cmd.AddCommand(pruningCmd(ctx), queryCmd(ctx))

	return cmd
}
