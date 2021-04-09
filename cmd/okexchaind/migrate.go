package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/app"
	evmtypes "github.com/okex/okexchain/x/evm/types"
	stakingtypes "github.com/okex/okexchain/x/staking/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	tmstate "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
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

	dataDir := filepath.Join(ctx.Config.RootDir, "data")
	blockStoreDB, err := openDB(blockStoreDB, dataDir)
	panicError(err)

	// update latest block according to app version
	blockState := store.BlockStoreStateJSON{
		Base:   version - 1,
		Height: version,
	}
	blockState.Save(blockStoreDB)

	blockStore := store.NewBlockStore(blockStoreDB)
	latestBlockHeight := blockStore.Height()
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
	chainApp.EvmKeeper.SetParams(deliverCtx, evmtypes.DefaultParams())

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.MaxValidators = uint16(1)
	stakingParams.Epoch = 1<<15 - 1
	chainApp.StakingKeeper.SetParams(deliverCtx, stakingParams)

	//TODO: just for test
	var lastValidatorPowers []stakingtypes.LastValidatorPower
	var valsUpdate abci.ValidatorUpdates
	chainApp.StakingKeeper.IterateLastValidatorPowers(deliverCtx, func(addr sdk.ValAddress, power int64) (stop bool) {
		lastValidatorPowers = append(lastValidatorPowers, stakingtypes.LastValidatorPower{Address: addr, Power: power})
		return false
	})
	for _, lv := range lastValidatorPowers[:1] {
		log.Println(lv.Address.String())
		chainApp.StakingKeeper.SetLastValidatorPower(deliverCtx, lv.Address, lv.Power)
		validator, found := chainApp.StakingKeeper.GetValidator(deliverCtx, lv.Address)
		if !found {
			panic(fmt.Sprintf("validator %s not found", lv.Address))
		}
		update := validator.ABCIValidatorUpdate()
		update.Power = lv.Power // keep the next-val-set offset, use the last power for the first block
		valsUpdate = append(valsUpdate, update)
	}

	if err != nil {
		panicError(err)
	}
	commitID := chainApp.MigrateCommit()

	updateState(dataDir, valsUpdate, commitID.Hash, version)
}

func createApp(ctx *server.Context, dataPath string) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, dataPath)
	db, err := openDB(applicationDB, dataDir)
	panicError(err)
	exapp := newApp(ctx.Logger, db, nil)
	return exapp.(*app.OKExChainApp)
}

//TODO: just for test
func updateState(dataDir string, valsUpdate abci.ValidatorUpdates, appHash []byte, height int64) {
	stateStoreDB, err := openDB(stateDB, dataDir)
	panicError(err)
	state := tmstate.LoadState(stateStoreDB)

	if len(valsUpdate) > 0 {
		vals, err := types.PB2TM.ValidatorUpdates(valsUpdate)
		panicError(err)
		state.Validators = types.NewValidatorSet(vals)
		state.NextValidators = types.NewValidatorSet(vals)
	}

	state.AppHash = appHash

	err = stateStoreDB.SetSync([]byte("stateKey"), state.Bytes())
	panicError(err)

	valInfo := &tmstate.ValidatorsInfo{
		LastHeightChanged: height + 1,
		ValidatorSet:      state.Validators,
	}

	err = stateStoreDB.Set(calcValidatorsKey(height+1), valInfo.Bytes())
	panicError(err)
}

func calcValidatorsKey(height int64) []byte {
	return []byte(fmt.Sprintf("validatorsKey:%v", height))
}
