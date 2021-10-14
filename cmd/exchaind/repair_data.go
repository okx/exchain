package main

import (
	"fmt"
	cfg "github.com/tendermint/tendermint/config"
	"io"
	"log"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/okex/exchain/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/iavl"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/mock"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var commitInterval int64

const FlagCommitInterval string = "commit-interval"

func repairStateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-state",
		Short: "Repair the SMB(state machine broken) data of node",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- repair data start ---------")

			repairState(ctx)
			log.Println("--------- repair data success ---------")
		},
	}
	cmd.Flags().Int64Var(&commitInterval, FlagCommitInterval, 1, "The number of interval heights for submitting Commit")
	return cmd
}

func repairState(ctx *server.Context) {
	// set ignore smb check
	sm.SetIgnoreSmbCheck(true)
	iavl.SetIgnoreVersionCheck(true)

	// load latest block height
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	latestBlockHeight := latestBlockHeight(dataDir)
	fmt.Println("???????", dataDir, latestBlockHeight)
	startBlockHeight := types.GetStartBlockHeight()
	if latestBlockHeight <= startBlockHeight+2 {
		panic(fmt.Sprintf("There is no need to repair data. The latest block height is %d, start block height is %d", latestBlockHeight, startBlockHeight))
	}

	// create proxy app
	proxyApp, repairApp, err := createRepairApp(ctx)
	panicError(err)
	// load start version
	startVersion := latestBlockHeight - ((latestBlockHeight - 2) % commitInterval) - 2
	if startVersion == 0 {
		panic("height too low, please restart from height 0 with genesis file")
	}
	fmt.Println("startVersion", latestBlockHeight, startVersion)
	err = repairApp.LoadStartVersion(startVersion)
	panicError(err)

	// load state
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	state, _, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	panicError(err)

	// repair data by apply the latest two blocks
	doRepair(ctx, state, stateStoreDB, proxyApp, startVersion, latestBlockHeight, dataDir)
	repairApp.StopStore()
}

func createRepairApp(ctx *server.Context) (proxy.AppConns, *app.OKExChainApp, error) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	repairApp := newRepairApp(ctx.Logger, db, nil)

	clientCreator := proxy.NewLocalClientCreator(repairApp)
	// Create the proxyApp and establish connections to the ABCI app (consensus, mempool, query).
	proxyApp, err := createAndStartProxyAppConns(clientCreator)
	return proxyApp, repairApp, err
}

func newRepairApp(logger tmlog.Logger, db dbm.DB, traceStore io.Writer) *app.OKExChainApp {
	return app.NewOKExChainApp(
		logger,
		db,
		traceStore,
		false,
		map[int64]bool{},
		0,
	)
}

func doRepair(ctx *server.Context, state sm.State, stateStoreDB dbm.DB,
	proxyApp proxy.AppConns, startHeight, latestHeight int64, dataDir string) {
	var err error
	for height := startHeight + 1; height <= latestHeight; height++ {
		repairBlock, repairBlockMeta := loadBlock(height, dataDir)
		blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), mock.Mempool{}, sm.MockEvidencePool{})
		blockExec.SetIsAsyncDeliverTx(cfg.IsAsyncDeliverTx())
		state, _, err = blockExec.ApplyBlock(state, repairBlockMeta.BlockID, repairBlock)
		panicError(err)
		res, err := proxyApp.Query().InfoSync(proxy.RequestInfo)
		panicError(err)
		repairedBlockHeight := res.LastBlockHeight
		repairedAppHash := res.LastBlockAppHash
		log.Println("Repaired block height", repairedBlockHeight)
		log.Println("Repaired app hash", fmt.Sprintf("%X", repairedAppHash))
		panic("test-----")

	}

}

func loadBlock(height int64, dataDir string) (*types.Block, *types.BlockMeta) {
	//rootDir := ctx.Config.RootDir
	//dataDir := filepath.Join(rootDir, "data")
	storeDB, err := openDB(blockStoreDB, dataDir)
	defer storeDB.Close()
	blockStore := store.NewBlockStore(storeDB)
	panicError(err)
	block := blockStore.LoadBlock(height)
	meta := blockStore.LoadBlockMeta(height)
	return block, meta
}

func latestBlockHeight(dataDir string) int64 {
	storeDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)
	defer storeDB.Close()
	blockStore := store.NewBlockStore(storeDB)
	return blockStore.Height()
}
