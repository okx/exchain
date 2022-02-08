package main

import (
	"log"
	"path/filepath"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

func exportAppCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-app",
		Short: "export current latest version of application db to new db into export dir",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- export start ---------")
			export(ctx)
			log.Println("--------- export success ---------")
		},
	}
	return cmd
}

// export current latest version of application db to new db into export dir
func export(ctx *server.Context) {
	fromApp := createApp(ctx, "data")
	toApp := createApp(ctx, "export")

	version := fromApp.LastCommitID().Version
	log.Println("export app version ", version)

	err := fromApp.Export(toApp.BaseApp, version)
	if err != nil {
		panicError(err)
	}
}

func createApp(ctx *server.Context, dataPath string) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, dataPath)
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	exapp := newApp(ctx.Logger, db, nil)
	return exapp.(*app.OKExChainApp)
}
