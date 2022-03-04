package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/x/common/analyzer"

	"github.com/okex/exchain/app/config"
	okexchain "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmiavl "github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/system"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/mock"
	"github.com/okex/exchain/libs/tendermint/node"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/state"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/store"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	replayedBlockDir = "replayed_block_dir"
	applicationDB    = "application"
	blockStoreDB     = "blockstore"
	stateDB          = "state"

	pprofAddrFlag       = "pprof_addr"
	runWithPprofFlag    = "gen_pprof"
	runWithPprofMemFlag = "gen_pprof_mem"

	saveBlock = "save_block"

	defaulPprofFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	defaultPprofFilePerm = 0644
)

func replayCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay",
		Short: "Replay blocks from local db",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// set external package flags
			setExternalPackageValue(cmd)
			types.InitSignatureCache()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- replay start ---------")
			pprofAddress := viper.GetString(pprofAddrFlag)
			go func() {
				err := http.ListenAndServe(pprofAddress, nil)
				if err != nil {
					fmt.Println(err)
				}
			}()

			dataDir := viper.GetString(replayedBlockDir)
			replayBlock(ctx, dataDir)
			log.Println("--------- replay success ---------")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(runWithPprofMemFlag) {
				log.Println("--------- gen pprof mem start ---------")
				err := dumpMemPprof()
				if err != nil {
					log.Println(err)
				} else {
					log.Println("--------- gen pprof mem success ---------")
				}
			}
		},
	}
	cmd.Flags().StringP(replayedBlockDir, "d", ".exchaind/data", "Directory of block data to be replayed")
	cmd.Flags().StringP(pprofAddrFlag, "p", "0.0.0.0:26661", "Address and port of pprof HTTP server listening")
	cmd.Flags().BoolVarP(&state.IgnoreSmbCheck, "ignore-smb", "i", false, "ignore state machine broken")
	cmd.Flags().Bool(types.FlagDownloadDDS, false, "get delta from dc/redis or not")
	cmd.Flags().Bool(types.FlagUploadDDS, false, "send delta to dc/redis or not")
	cmd.Flags().String(types.FlagRedisUrl, "localhost:6379", "redis url")
	cmd.Flags().String(types.FlagRedisAuth, "", "redis auth")
	cmd.Flags().Int(types.FlagRedisExpire, 300, "delta expiration time. unit is second")
	cmd.Flags().Int(types.FlagRedisDB, 0, "delta db num")
	cmd.Flags().Int(types.FlagDeltaVersion, types.DeltaVersion, "Specify delta version")
	cmd.Flags().Bool(types.FlagFastQuery, false, "enable watch db or not")

	cmd.Flags().String(server.FlagPruning, storetypes.PruningOptionNothing, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(server.FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(config.FlagPprofAutoDump, false, "Enable auto dump pprof")
	cmd.Flags().String(config.FlagPprofCollectInterval, "30s", "Interval for pprof dump loop")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentMin, 45, "TriggerPercentMin of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentDiff, 50, "TriggerPercentDiff of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentAbs, 50, "TriggerPercentAbs of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentMin, 70, "TriggerPercentMin of mem to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentDiff, 50, "TriggerPercentDiff of mem to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentAbs, 75, "TriggerPercentAbs of cpu mem dump pprof")
	cmd.Flags().String(config.FlagPprofCoolDown, "3m", "The cool down time after every type of pprof dump")
	cmd.Flags().Int64(config.FlagPprofAbciElapsed, 5000, "Elapsed time of abci in millisecond for pprof dump")
	cmd.Flags().Bool(config.FlagPprofUseCGroup, false, "Use cgroup when exchaind run in docker")
	cmd.Flags().IntVar(&iavl.IavlCacheSize, iavl.FlagIavlCacheSize, 1000000, "Max size of iavl cache")
	cmd.Flags().StringToIntVar(&tmiavl.OutputModules, tmiavl.FlagOutputModules, map[string]int{}, "decide which module in iavl to be printed")
	cmd.Flags().Int64Var(&tmiavl.CommitIntervalHeight, tmiavl.FlagIavlCommitIntervalHeight, 100, "Max interval to commit node cache into leveldb")
	cmd.Flags().Int64Var(&tmiavl.MinCommitItemCount, tmiavl.FlagIavlMinCommitItemCount, 500000, "Min nodes num to triggle node cache commit")
	cmd.Flags().IntVar(&tmiavl.HeightOrphansCacheSize, tmiavl.FlagIavlHeightOrphansCacheSize, 8, "Max orphan version to cache in memory")
	cmd.Flags().IntVar(&tmiavl.MaxCommittedHeightNum, tmiavl.FlagIavlMaxCommittedHeightNum, 8, "Max committed version to cache in memory")
	cmd.Flags().BoolVar(&tmiavl.EnableAsyncCommit, tmiavl.FlagIavlEnableAsyncCommit, false, "Enable cache iavl node data to optimization leveldb pruning process")
	cmd.Flags().BoolVar(&system.EnableGid, system.FlagEnableGid, false, "Display goroutine id in log")
	cmd.Flags().Bool(runWithPprofFlag, false, "Dump the pprof of the entire replay process")
	cmd.Flags().Bool(runWithPprofMemFlag, false, "Dump the mem profile of the entire replay process")
	cmd.Flags().Bool(sm.FlagParalleledTx, false, "pall Tx")
	cmd.Flags().Bool(saveBlock, false, "save block when replay")
	cmd.Flags().Int64(config.FlagMaxGasUsedPerBlock, -1, "Maximum gas used of transactions in a block")
	cmd.Flags().Bool(sdk.FlagMultiCache, false, "Enable multi cache")
	cmd.Flags().Int(sdk.MaxAccInMultiCache, 0, "max acc in multi cache")
	cmd.Flags().Int(sdk.MaxStorageInMultiCache, 0, "max storage in multi cache")
	cmd.Flags().Bool(flatkv.FlagEnable, false, "Enable flat kv storage for read performance")
	cmd.Flags().String(app.Elapsed, app.DefaultElapsedSchemas, "schemaName=1|0,,,")
	cmd.Flags().Bool(analyzer.FlagEnableAnalyzer, true, "Enable auto open log analyzer")
	cmd.Flags().Int(types.FlagSigCacheSize, 200000, "Maximum number of signatures in the cache")
	return cmd
}

// setExternalPackageValue set external package config value.
func setExternalPackageValue(cmd *cobra.Command) {
	types.DownloadDelta = viper.GetBool(types.FlagDownloadDDS)
	types.UploadDelta = viper.GetBool(types.FlagUploadDDS)
	types.FastQuery = viper.GetBool(types.FlagFastQuery)
	types.DeltaVersion = viper.GetInt(types.FlagDeltaVersion)
}

// replayBlock replays blocks from db, if something goes wrong, it will panic with error message.
func replayBlock(ctx *server.Context, originDataDir string) {
	config.RegisterDynamicConfig(ctx.Logger.With("module", "config"))
	proxyApp, err := createProxyApp(ctx)
	panicError(err)

	res, err := proxyApp.Query().InfoSync(proxy.RequestInfo)
	panicError(err)
	currentBlockHeight := res.LastBlockHeight
	currentAppHash := res.LastBlockAppHash
	log.Println("current block height", "height", currentBlockHeight)
	log.Println("current app hash", "appHash", fmt.Sprintf("%X", currentAppHash))

	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)

	genesisDocProvider := node.DefaultGenesisDocProviderFunc(ctx.Config)
	state, genDoc, err := node.LoadStateFromDBOrGenesisDocProvider(stateStoreDB, genesisDocProvider)
	panicError(err)

	// If startBlockHeight == 0 it means that we are at genesis and hence should initChain.
	if currentBlockHeight == types.GetStartBlockHeight() {
		err := initChain(state, stateStoreDB, genDoc, proxyApp)
		panicError(err)
		state = sm.LoadState(stateStoreDB)
	}
	//cache chain epoch
	err = okexchain.SetChainId(genDoc.ChainID)
	if err != nil {
		panicError(err)
	}
	// replay
	doReplay(ctx, state, stateStoreDB, proxyApp, originDataDir, currentAppHash, currentBlockHeight, int8(0))
	if viper.GetBool(sm.FlagParalleledTx) {
		baseapp.ParaLog.PrintLog()
	}
}

// panic if error is not nil
func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func openDB(dbName string, dataDir string) (db dbm.DB, err error) {
	return sdk.NewLevelDB(dbName, dataDir)
}

func createProxyApp(ctx *server.Context) (proxy.AppConns, error) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	app := newApp(ctx.Logger, db, nil)
	clientCreator := proxy.NewLocalClientCreator(app)
	return createAndStartProxyAppConns(clientCreator)
}

func createAndStartProxyAppConns(clientCreator proxy.ClientCreator) (proxy.AppConns, error) {
	proxyApp := proxy.NewAppConns(clientCreator)
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("error starting proxy app connections: %v", err)
	}
	return proxyApp, nil
}

func initChain(state sm.State, stateDB dbm.DB, genDoc *types.GenesisDoc, proxyApp proxy.AppConns) error {
	validators := make([]*types.Validator, len(genDoc.Validators))
	for i, val := range genDoc.Validators {
		validators[i] = types.NewValidator(val.PubKey, val.Power)
	}
	validatorSet := types.NewValidatorSet(validators)
	nextVals := types.TM2PB.ValidatorUpdates(validatorSet)
	csParams := types.TM2PB.ConsensusParams(genDoc.ConsensusParams)
	req := abci.RequestInitChain{
		Time:            genDoc.GenesisTime,
		ChainId:         genDoc.ChainID,
		ConsensusParams: csParams,
		Validators:      nextVals,
		AppStateBytes:   genDoc.AppState,
	}
	res, err := proxyApp.Consensus().InitChainSync(req)
	if err != nil {
		return err
	}

	if state.LastBlockHeight == types.GetStartBlockHeight() { //we only update state when we are in initial state
		// If the app returned validators or consensus params, update the state.
		if len(res.Validators) > 0 {
			vals, err := types.PB2TM.ValidatorUpdates(res.Validators)
			if err != nil {
				return err
			}
			state.Validators = types.NewValidatorSet(vals)
			state.NextValidators = types.NewValidatorSet(vals)
		} else if len(genDoc.Validators) == 0 {
			// If validator set is not set in genesis and still empty after InitChain, exit.
			return fmt.Errorf("validator set is nil in genesis and still empty after InitChain")
		}

		if res.ConsensusParams != nil {
			state.ConsensusParams = state.ConsensusParams.Update(res.ConsensusParams)
		}
		sm.SaveState(stateDB, state)
	}
	return nil
}

var (
	alreadyInit  bool
	stateStoreDb *store.BlockStore
)

// TODO need delete
func SaveBlock(ctx *server.Context, originDB *store.BlockStore, height int64) {
	if !alreadyInit {
		alreadyInit = true
		dataDir := filepath.Join(ctx.Config.RootDir, "data")
		blockStoreDB, err := openDB(blockStoreDB, dataDir)
		panicError(err)
		stateStoreDb = store.NewBlockStore(blockStoreDB)
	}

	block := originDB.LoadBlock(height)
	meta := originDB.LoadBlockMeta(height)
	seenCommit := originDB.LoadSeenCommit(height)

	ps := types.NewPartSetFromHeader(meta.BlockID.PartsHeader)
	for index := 0; index < ps.Total(); index++ {
		ps.AddPart(originDB.LoadBlockPart(height, index))
	}

	stateStoreDb.SaveBlock(block, ps, seenCommit)
}

func doReplay(ctx *server.Context, state sm.State, stateStoreDB dbm.DB,
	proxyApp proxy.AppConns, originDataDir string, lastAppHash []byte, lastBlockHeight int64, deliverTxsMode int8) {
	originBlockStoreDB, err := openDB(blockStoreDB, originDataDir)
	panicError(err)
	originBlockStore := store.NewBlockStore(originBlockStoreDB)
	originLatestBlockHeight := originBlockStore.Height()
	log.Println("origin latest block height", "height", originLatestBlockHeight)

	haltheight := viper.GetInt64(server.FlagHaltHeight)
	if haltheight == 0 {
		haltheight = originLatestBlockHeight
	}
	if haltheight <= lastBlockHeight+1 {
		panic("haltheight <= startBlockHeight please check data or height")
	}

	log.Println("replay stop block height", "height", haltheight, "lastBlockHeight", lastBlockHeight, "state.LastBlockHeight", state.LastBlockHeight)

	// Replay blocks up to the latest in the blockstore.
	if lastBlockHeight == state.LastBlockHeight+1 {
		global.SetGlobalHeight(lastBlockHeight)
		abciResponses, err := sm.LoadABCIResponses(stateStoreDB, lastBlockHeight)
		panicError(err)
		mockApp := newMockProxyApp(lastAppHash, abciResponses)
		block := originBlockStore.LoadBlock(lastBlockHeight)
		meta := originBlockStore.LoadBlockMeta(lastBlockHeight)
		blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, mockApp, mock.Mempool{}, sm.MockEvidencePool{}, int8(0))
		blockExec.SetIsAsyncDeliverTx(false) // mockApp not support parallel tx
		state, _, err = blockExec.ApplyBlock(state, meta.BlockID, block)
		panicError(err)
	}

	blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), mock.Mempool{}, sm.MockEvidencePool{}, deliverTxsMode)
	if viper.GetBool(runWithPprofFlag) {
		startDumpPprof()
		defer stopDumpPprof()
	}

	baseapp.SetGlobalMempool(mock.Mempool{}, ctx.Config.Mempool.SortTxByGp, ctx.Config.Mempool.EnablePendingPool)
	needSaveBlock := viper.GetBool(saveBlock)
	global.SetGlobalHeight(lastBlockHeight + 1)
	blockExec.SetDeliverTxsMode(int8(1))
	//for height := lastBlockHeight + 1; height <= haltheight; height++ {
	height := lastBlockHeight + 1
	for i := 0; i <= 37; i++ {
		log.Println("replaying ", height)
		block := originBlockStore.LoadBlock(height)
		meta := originBlockStore.LoadBlockMeta(height)
		blockExec.SetIsAsyncDeliverTx(viper.GetBool(sm.FlagParalleledTx))
		state, _, err = blockExec.ApplyBlock(state, meta.BlockID, block)
		panicError(err)
		if needSaveBlock {
			SaveBlock(ctx, originBlockStore, height)
		}
		height++
	}
}

func dumpMemPprof() error {
	fileName := fmt.Sprintf("replay_pprof_%s.mem.bin", time.Now().Format("20060102150405"))
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create mem pprof file %s error: %w", fileName, err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err = pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("could not write memory profile: %w", err)
	}
	return nil
}

func startDumpPprof() {
	var (
		binarySuffix = time.Now().Format("20060102150405") + ".bin"
	)
	fileName := fmt.Sprintf("replay_pprof_%s", binarySuffix)
	bf, err := os.OpenFile(fileName, defaulPprofFileFlags, defaultPprofFilePerm)
	if err != nil {
		fmt.Printf("open pprof file(%s) error:%s\n", fileName, err.Error())
		return
	}

	err = pprof.StartCPUProfile(bf)
	if err != nil {
		fmt.Printf("dump pprof StartCPUProfile error:%s\n", err.Error())
		return
	}
	fmt.Printf("start to dump pprof file(%s)\n", fileName)
}

func stopDumpPprof() {
	pprof.StopCPUProfile()
	fmt.Printf("dump pprof successfully\n")
}

func newMockProxyApp(appHash []byte, abciResponses *sm.ABCIResponses) proxy.AppConnConsensus {
	clientCreator := proxy.NewLocalClientCreator(&mockProxyApp{
		appHash:       appHash,
		abciResponses: abciResponses,
	})
	cli, _ := clientCreator.NewABCIClient()
	err := cli.Start()
	if err != nil {
		panic(err)
	}
	return proxy.NewAppConnConsensus(cli)
}

type mockProxyApp struct {
	abci.BaseApplication

	appHash       []byte
	txCount       int
	abciResponses *sm.ABCIResponses
}

func (mock *mockProxyApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	r := mock.abciResponses.DeliverTxs[mock.txCount]
	mock.txCount++
	if r == nil { //it could be nil because of amino unMarshall, it will cause an empty ResponseDeliverTx to become nil
		return abci.ResponseDeliverTx{}
	}
	return *r
}

func (mock *mockProxyApp) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	mock.txCount = 0
	return *mock.abciResponses.EndBlock
}

func (mock *mockProxyApp) Commit(req abci.RequestCommit) abci.ResponseCommit {
	return abci.ResponseCommit{Data: mock.appHash}
}
