package main

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	"github.com/spf13/viper"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/iavl"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mock"
	"github.com/okex/exchain/libs/tendermint/node"
	"github.com/okex/exchain/libs/tendermint/proxy"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/store"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/cobra"
	dbm "github.com/tendermint/tm-db"
)

const (
	FlagStartHeight string = "start-height"
)

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
	cmd.Flags().Bool(sm.FlagParalleledTx, false, "parallel execution for evm txs")
	cmd.Flags().Int64(FlagStartHeight, 0, "Set the start block height for repair")
	return cmd
}

type repairApp struct {
	db dbm.DB
	*app.OKExChainApp
}

func (app *repairApp) getLatestVersion() int64 {
	rs := rootmulti.NewStore(app.db)
	return rs.GetLatestVersion()
}

func repairState(ctx *server.Context) {
	// set ignore smb check
	sm.SetIgnoreSmbCheck(true)
	iavl.SetIgnoreVersionCheck(true)

	// load latest block height
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	latestBlockHeight := latestBlockHeight(dataDir)
	startBlockHeight := types.GetStartBlockHeight()
	if latestBlockHeight <= startBlockHeight+2 {
		panic(fmt.Sprintf("There is no need to repair data. The latest block height is %d, start block height is %d", latestBlockHeight, startBlockHeight))
	}

	// create proxy app
	proxyApp, repairApp, err := createRepairApp(ctx)
	panicError(err)

	// load state
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)
	genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	state, _, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	panicError(err)

	// load start version
	startVersion := viper.GetInt64(FlagStartHeight)
	if startVersion == 0 {
		latestVersion := repairApp.getLatestVersion()
		startVersion = latestVersion - 2
	}
	if startVersion == 0 {
		panic("height too low, please restart from height 0 with genesis file")
	}
	err = repairApp.LoadStartVersion(startVersion)
	panicError(err)

	// repair data by apply the latest two blocks
	doRepair(ctx, state, stateStoreDB, proxyApp, startVersion, latestBlockHeight, dataDir)
	repairApp.StopStore()
}

func createRepairApp(ctx *server.Context) (proxy.AppConns, *repairApp, error) {
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

func newRepairApp(logger tmlog.Logger, db dbm.DB, traceStore io.Writer) *repairApp {
	return &repairApp{db, app.NewOKExChainApp(
		logger,
		db,
		traceStore,
		false,
		map[int64]bool{},
		0,
	)}
}

func doRepair(ctx *server.Context, state sm.State, stateStoreDB dbm.DB,
	proxyApp proxy.AppConns, startHeight, latestHeight int64, dataDir string) {
	// height n with multiple rounds will have shifted the priority such that state.LoadValidators(n) will return an incorrect value.
	// so we should use stateCopy to correct the repaired state below.
	stateCopy := state.Copy()
	// construct state for repair
	state = constructStartState(state, stateStoreDB, startHeight)
	ctx.Logger.Debug("constructStartState", "state", fmt.Sprintf("%+v", state))
	var err error
	// repair state
	blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), mock.Mempool{}, sm.MockEvidencePool{})
	blockExec.SetIsAsyncDeliverTx(viper.GetBool(sm.FlagParalleledTx))
	for height := startHeight + 1; height <= latestHeight; height++ {
		repairBlock, repairBlockMeta := loadBlock(height, dataDir)
		state, _, err = blockExec.ApplyBlock(state, repairBlockMeta.BlockID, repairBlock, &types.Deltas{}, &types.WatchData{})
		panicError(err)
		// use stateCopy to correct the repaired state
		if state.LastBlockHeight == stateCopy.LastBlockHeight {
			state = stateCopy
			sm.SaveState(stateStoreDB, state)
		}
		ctx.Logger.Debug("repairedState", "state", fmt.Sprintf("%+v", state))
		res, err := proxyApp.Query().InfoSync(proxy.RequestInfo)
		panicError(err)
		repairedBlockHeight := res.LastBlockHeight
		repairedAppHash := res.LastBlockAppHash
		log.Println("Repaired block height", repairedBlockHeight)
		log.Println("Repaired app hash", fmt.Sprintf("%X", repairedAppHash))
	}
}

func constructStartState(state sm.State, stateStoreDB dbm.DB, startHeight int64) sm.State {
	stateCopy := state.Copy()
	validators, err := sm.LoadValidators(stateStoreDB, startHeight)
	panicError(err)
	lastValidators, err := sm.LoadValidators(stateStoreDB, startHeight-1)
	panicError(err)
	nextValidators, err := sm.LoadValidators(stateStoreDB, startHeight+1)
	panicError(err)
	consensusParams, err := sm.LoadConsensusParams(stateStoreDB, startHeight+1)
	panicError(err)
	stateCopy.Validators = validators
	stateCopy.LastValidators = lastValidators
	stateCopy.NextValidators = nextValidators
	stateCopy.ConsensusParams = consensusParams
	stateCopy.LastBlockHeight = startHeight
	return stateCopy
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
