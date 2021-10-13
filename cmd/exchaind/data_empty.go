// +build !rocksdb

package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

func dataCmd(ctx *server.Context) *cobra.Command {
	return &cobra.Command{}
}
