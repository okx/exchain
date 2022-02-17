package main

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

const (
	FlagDisplayVersion      string = "display-version"
	FlagDisplayContractAddr string = "display-contract-address"
)

func displayStateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "display-state",
		Short: "display evm storage account's state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display state start ---------")
			displayState(ctx)
			log.Println("--------- display state end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

func displayState(ctx *server.Context) {
	dispApp := newDisplayApp(ctx)

	// load start version
	displayVersion := viper.GetInt64(FlagDisplayVersion)
	dispApp.EvmKeeper.SetTargetMptVersion(displayVersion)

	err := dispApp.LoadHeight(displayVersion)
	panicError(err)

	contractAddr := viper.GetString(FlagDisplayContractAddr)
	addr := ethcmn.HexToAddress(contractAddr)

	// init deliver state
	dispApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: displayVersion + 1}})

	_ = dispApp.EvmKeeper.ForEachStorage(dispApp.GetDeliverStateCtx(), addr, func(key, value ethcmn.Hash) bool {
		fmt.Println("Key is: ", key.String(), ", value is: ", value.String())
		return false
	})
}

func newDisplayApp(ctx *server.Context) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	return app.NewOKExChainApp(
		ctx.Logger,
		db,
		nil,
		false,
		map[int64]bool{},
		0,
	)
}
