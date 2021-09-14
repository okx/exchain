package main

import (
	"fmt"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/okex/exchain/app/config"
	"github.com/tendermint/tendermint/state"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/mock"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	dataDirFlag   = "data_dir"
	applicationDB = "application"
	blockStoreDB  = "blockstore"
	stateDB       = "state"
	pprofAddrFlag = "pprof_addr"
	//FlagHaltHeight = "halt-height"
)

func replayCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay",
		Short: "Replay blocks from local db",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- replay start ---------")
			pprofAddress := viper.GetString(pprofAddrFlag)
			go func() {
				err := http.ListenAndServe(pprofAddress, nil)
				if err != nil {
					fmt.Println(err)
				}
			}()

			dataDir := viper.GetString(dataDirFlag)
			replayBlock(ctx, dataDir)
			log.Println("--------- replay success ---------")
		},
	}
	cmd.Flags().StringP(dataDirFlag, "d", ".exchaind/data", "Directory of block data for replaying")
	cmd.Flags().StringP(pprofAddrFlag, "p", "0.0.0.0:26661", "Address and port of pprof HTTP server listening")
	cmd.Flags().BoolVarP(&state.IgnoreSmbCheck, "ignore-smb", "i", false, "ignore state machine broken")
	cmd.Flags().String(server.FlagPruning, storetypes.PruningOptionNothing, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(server.FlagHaltHeight, 0, "Block height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Bool(config.FlagPprofAutoDump, false, "Enable auto dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentMin, 45, "TriggerPercentMin of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentDiff, 50, "TriggerPercentDiff of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofCpuTriggerPercentAbs, 50, "TriggerPercentAbs of cpu to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentMin, 70, "TriggerPercentMin of mem to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentDiff, 50, "TriggerPercentDiff of mem to dump pprof")
	cmd.Flags().Int(config.FlagPprofMemTriggerPercentAbs, 75, "TriggerPercentAbs of cpu mem dump pprof")
	return cmd
}

// replayBlock replays blocks from db, if something goes wrong, it will panic with error message.
func replayBlock(ctx *server.Context, originDataDir string) {
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

	// replay
	startBlockHeight := currentBlockHeight + 1
	//doReplay(ctx, state, stateStoreDB, proxyApp, originDataDir, startBlockHeight)
	haltBlockHeight := viper.GetInt64(server.FlagHaltHeight)
	doReplay(ctx, state, stateStoreDB, proxyApp, originDataDir, startBlockHeight, haltBlockHeight)
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

func doReplay(ctx *server.Context, state sm.State, stateStoreDB dbm.DB,
	proxyApp proxy.AppConns, originDataDir string, startBlockHeight int64, haltBlockHeight int64) {
	originBlockStoreDB, err := openDB(blockStoreDB, originDataDir)
	panicError(err)
	originBlockStore := store.NewBlockStore(originBlockStoreDB)
	originLatestBlockHeight := originBlockStore.Height()
	log.Println("origin latest block height", "height", originLatestBlockHeight)

	haltheight := haltBlockHeight
	if haltheight == 0 {
		haltheight = originLatestBlockHeight
	}
	if haltheight <= startBlockHeight {
		panic("haltheight <= startBlockHeight please check data or height")
	}

	log.Println("replay stop block height", "height", haltheight)

	for height := startBlockHeight; height <= haltheight; height++ {
		log.Println("replaying ", height)
		block := originBlockStore.LoadBlock(height)
		meta := originBlockStore.LoadBlockMeta(height)

		blockExec := sm.NewBlockExecutor(stateStoreDB, ctx.Logger, proxyApp.Consensus(), mock.Mempool{}, sm.MockEvidencePool{})
		state, _, err = blockExec.ApplyBlock(state, meta.BlockID, block)
		panicError(err)
	}
}
