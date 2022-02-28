package main

import (
	"os"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/server"
	svrcmd "github.com/okex/exchain/ibc-3rd/cosmos-v443/server/cmd"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp/simd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, simapp.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
