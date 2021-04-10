package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/okex/okexchain/app"
	evmtypes "github.com/okex/okexchain/x/evm/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	tmstate "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	"log"
	"path/filepath"
)

func migrateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrate scheme for application db",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- migrate start ---------")
			migrate(ctx)
			log.Println("--------- migrate success ---------")
		},
	}
	return cmd
}

func migrate(ctx *server.Context) {
	chainApp := createApp(ctx, "data")
	version := chainApp.LastCommitID().Version
	log.Println("latest app height", version)

	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	blockStoreDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)

	blockStore := store.NewBlockStore(blockStoreDB)
	latestBlockHeight := version
	if version != latestBlockHeight {
		panicError(fmt.Errorf("app version %d not equal to blockstore height %d", version, latestBlockHeight))
	}

	log.Println("latest block height", latestBlockHeight)
	block := blockStore.LoadBlock(latestBlockHeight)
	req := abci.RequestBeginBlock{
		Hash:   block.Hash(),
		Header: types.TM2PB.Header(&block.Header),
	}

	deliverCtx := chainApp.DeliverStateCtx(req)
	evmParams := evmtypes.DefaultParams()
	evmParams.EnableCall = true
	evmParams.EnableCreate = true
	log.Println("set evm params:\n", evmParams)
	chainApp.EvmKeeper.SetParams(deliverCtx, evmParams)

	stakingParams := stakingtypes.DefaultParams()
	log.Println("set staking params: \n", stakingParams)
	chainApp.StakingKeeper.SetParams(deliverCtx, stakingParams)

	commitID := chainApp.MigrateCommit()

	evmParams = chainApp.EvmKeeper.GetParams(deliverCtx)
	log.Println("get evm params after set: \n", evmParams)

	stakingParams = chainApp.StakingKeeper.GetParams(deliverCtx)
	log.Println("get staking params after set: \n", stakingParams)

	updateState(dataDir, nil, commitID.Hash, version)
}

func createApp(ctx *server.Context, dataPath string) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, dataPath)
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	exapp := newApp(ctx.Logger, db, nil)
	return exapp.(*app.OKExChainApp)
}

func updateState(dataDir string, valsUpdate abci.ValidatorUpdates, appHash []byte, height int64) {
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)
	state := tmstate.LoadState(stateStoreDB)

	state.AppHash = appHash

	err = stateStoreDB.SetSync([]byte("stateKey"), state.Bytes())
	panicError(err)
}
