package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/okex/okexchain/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/store"
	"log"
	"path/filepath"
)

func readDBCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "readDB",
		Short: "readDB scheme for application db",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- readDB start ---------")
			readDB(ctx)
			log.Println("--------- readDB success ---------")
		},
	}
	return cmd
}

func readDB(ctx *server.Context) {
	chainApp := createApp(ctx, "data")
	version := chainApp.LastCommitID().Version
	log.Println("latest app height", version)

	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	blockStoreDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)

	blockStore := store.NewBlockStore(blockStoreDB)
	latestBlockHeight := blockStore.Height()

	log.Println("latest block height", latestBlockHeight)
}

func createApp(ctx *server.Context, dataPath string) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, dataPath)
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	exapp := newApp(ctx.Logger, db, nil)
	return exapp.(*app.OKExChainApp)
}
