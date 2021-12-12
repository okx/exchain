package state

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/delta"
	"github.com/okex/exchain/libs/tendermint/delta/redis-cgi"
	"time"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/fail"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	mempl "github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
	dbm "github.com/tendermint/tm-db"
)

//-----------------------------------------------------------------------------
// BlockExecutor handles block execution and state updates.
// It exposes ApplyBlock(), which validates & executes the block, updates state w/ ABCI responses,
// then commits and updates the mempool atomically, then saves state.

// BlockExecutor provides the context and accessories for properly executing a block.
type BlockExecutor struct {
	// save state, validators, consensus params, abci responses here
	db dbm.DB

	// download or upload data to dds
	deltaBroker   delta.DeltaBroker
	deltaCh       chan *types.Deltas
	deltaHeightCh chan int64

	// execute the app against this
	proxyApp proxy.AppConnConsensus

	// events
	eventBus types.BlockEventPublisher

	// manage the mempool lock during commit
	// and update both with block results after commit.
	mempool mempl.Mempool
	evpool  EvidencePool

	logger log.Logger
	metrics *Metrics
	isAsync bool

	proactivelyRunTx bool
	prerunChan chan *executionContext
	prerunResultChan chan *executionContext
	prerunIndex int64
	prerunContext *executionContext
}


type BlockExecutorOption func(executor *BlockExecutor)

func BlockExecutorWithMetrics(metrics *Metrics) BlockExecutorOption {
	return func(blockExec *BlockExecutor) {
		blockExec.metrics = metrics
	}
}

// NewBlockExecutor returns a new BlockExecutor with a NopEventBus.
// Call SetEventBus to provide one.
func NewBlockExecutor(
	db dbm.DB,
	logger log.Logger,
	proxyApp proxy.AppConnConsensus,
	mempool mempl.Mempool,
	evpool EvidencePool,
	options ...BlockExecutorOption,
) *BlockExecutor {
	res := &BlockExecutor{
		db:             db,
		proxyApp:       proxyApp,
		eventBus:       types.NopEventBus{},
		mempool:        mempool,
		evpool:         evpool,
		logger:         logger,
		metrics:        NopMetrics(),
		isAsync:        viper.GetBool(FlagParalleledTx),
		deltaCh:        make(chan *types.Deltas, 1),
		deltaHeightCh:  make(chan int64, 1),
		prerunChan:           make(chan *executionContext, 1),
		prerunResultChan:     make(chan *executionContext, 1),
	}

	for _, option := range options {
		option(res)
	}

	res.deltaBroker = redis_cgi.NewRedisClient(types.RedisUrl())
	if types.EnableDownloadDelta() {
		go res.GetDeltaFromDDS()
	}

	return res
}

func (blockExec *BlockExecutor) SetIsAsyncDeliverTx(sw bool) {
	blockExec.isAsync = sw

}
func (blockExec *BlockExecutor) DB() dbm.DB {
	return blockExec.db
}

// SetEventBus - sets the event bus for publishing block related events.
// If not called, it defaults to types.NopEventBus.
func (blockExec *BlockExecutor) SetEventBus(eventBus types.BlockEventPublisher) {
	blockExec.eventBus = eventBus
	blockExec.mempool.SetEventBus(eventBus)
}

// CreateProposalBlock calls state.MakeBlock with evidence from the evpool
// and txs from the mempool. The max bytes must be big enough to fit the commit.
// Up to 1/10th of the block space is allcoated for maximum sized evidence.
// The rest is given to txs, up to the max gas.
func (blockExec *BlockExecutor) CreateProposalBlock(
	height int64,
	state State, commit *types.Commit,
	proposerAddr []byte,
) (*types.Block, *types.PartSet) {

	maxBytes := state.ConsensusParams.Block.MaxBytes
	maxGas := state.ConsensusParams.Block.MaxGas

	// Fetch a limited amount of valid evidence
	maxNumEvidence, _ := types.MaxEvidencePerBlock(maxBytes)
	evidence := blockExec.evpool.PendingEvidence(maxNumEvidence)

	// Fetch a limited amount of valid txs
	maxDataBytes := types.MaxDataBytes(maxBytes, state.Validators.Size(), len(evidence))
	if cfg.DynamicConfig.GetMaxGasUsedPerBlock() > -1 {
		maxGas = cfg.DynamicConfig.GetMaxGasUsedPerBlock()
	}
	txs := blockExec.mempool.ReapMaxBytesMaxGas(maxDataBytes, maxGas)

	return state.MakeBlock(height, txs, commit, evidence, proposerAddr)
}

// ValidateBlock validates the given block against the given state.
// If the block is invalid, it returns an error.
// Validation does not mutate state, but does require historical information from the stateDB,
// ie. to verify evidence from a validator at an old height.
func (blockExec *BlockExecutor) ValidateBlock(state State, block *types.Block) error {
	if IgnoreSmbCheck {
		// debug only
		return nil
	}
	return validateBlock(blockExec.evpool, blockExec.db, state, block)
}

// ApplyBlock validates the block against the state, executes it against the app,
// fires the relevant events, commits the app, and saves the new state and responses.
// It returns the new state and the block height to retain (pruning older blocks).
// It's the only function that needs to be called
// from outside this package to process and commit an entire block.
// It takes a blockID to avoid recomputing the parts hash.
func (blockExec *BlockExecutor) ApplyBlock(
	state State, blockID types.BlockID, block *types.Block, deltas *types.Deltas) (State, int64, *types.Deltas, error) {
	if ApplyBlockPprofTime >= 0 {
		f, t := PprofStart()
		defer PprofEnd(int(block.Height), f, t)
	}
	trc := trace.NewTracer(trace.ApplyBlock)
	var inAbciRspLen, inDeltaLen, inWatchLen int
	if deltas == nil || deltas.Height != block.Height {
		deltas = &types.Deltas{}
	}

	defer func() {
		trace.GetElapsedInfo().AddInfo(trace.Height, fmt.Sprintf("%d", block.Height))
		trace.GetElapsedInfo().AddInfo(trace.Tx, fmt.Sprintf("%d", len(block.Data.Txs)))
		trace.GetElapsedInfo().AddInfo(trace.InDelta, fmt.Sprintf(
			"abciRspLen<%d>, deltaLen<%d>, watchLen<%d>", inAbciRspLen, inDeltaLen, inWatchLen))
		trace.GetElapsedInfo().AddInfo(trace.OutDelta, fmt.Sprintf(
			"abciRspLen<%d>, deltaLen<%d>, watchLen<%d>", len(deltas.ABCIRsp), len(deltas.DeltasBytes), len(deltas.WatchBytes)))
		trace.GetElapsedInfo().AddInfo(trace.RunTx, trc.Format())
		trace.GetElapsedInfo().SetElapsedTime(trc.GetElapsedTime())

		now := time.Now().UnixNano()
		blockExec.metrics.IntervalTime.Set(float64(now-blockExec.metrics.lastBlockTime) / 1e6)
		blockExec.metrics.lastBlockTime = now
	}()

	trc.Pin("ValidateBlock")
	if err := blockExec.ValidateBlock(state, block); err != nil {
		return state, 0, deltas, ErrInvalidBlock(err)
	}

	trc.Pin("GetDelta")

	fastQuery := types.IsFastQuery()
	applyDelta := types.EnableApplyP2PDelta()
	broadDelta := types.EnableBroadcastP2PDelta()
	downloadDelta := types.EnableDownloadDelta()
	uploadDelta := types.EnableUploadDelta()

	useDeltas, deltas := blockExec.prepareStateDelta(block, deltas)
	inAbciRspLen = len(deltas.ABCIRsp)
	inDeltaLen = len(deltas.DeltasBytes)
	inWatchLen = len(deltas.WatchBytes)

	trc.Pin(trace.Abci)

	startTime := time.Now().UnixNano()
	var abciResponses *ABCIResponses
	var err error
	if useDeltas {
		SetCenterBatch(deltas.WatchBytes)
		execBlockOnProxyAppWithDeltas(blockExec.proxyApp, block, blockExec.db)
		err = types.Json.Unmarshal(deltas.ABCIRsp, &abciResponses)
		if err != nil {
			panic(err)
		}
	} else {
		if blockExec.proactivelyRunTx {
			abciResponses, err = blockExec.getPrerunResult(blockExec.prerunContext)
			blockExec.prerunContext = nil
			blockExec.prerunIndex = 0
		} else {
			ctx := &executionContext{
				logger: blockExec.logger,
				block: block,
				db: blockExec.db,
				proxyApp: blockExec.proxyApp,
			}
			if blockExec.isAsync {
				abciResponses, err = execBlockOnProxyAppAsync(blockExec.logger, blockExec.proxyApp, block, blockExec.db)
			} else {
				abciResponses, err = execBlockOnProxyApp(ctx)
			}
		}

		if broadDelta || uploadDelta {
			bytes, err := types.Json.Marshal(abciResponses)
			if err != nil {
				panic(err)
			}
			deltas.ABCIRsp = bytes
		}
	}
	if err != nil {
		return state, 0, deltas, ErrProxyAppConn(err)
	}

	fail.Fail() // XXX

	trc.Pin(trace.SaveResp)

	// Save the results before we commit.
	SaveABCIResponses(blockExec.db, block.Height, abciResponses)

	fail.Fail() // XXX
	endTime := time.Now().UnixNano()
	blockExec.metrics.BlockProcessingTime.Observe(float64(endTime-startTime) / 1e6)
	blockExec.metrics.AbciTime.Set(float64(endTime-startTime) / 1e6)

	// validate the validator updates and convert to tendermint types
	abciValUpdates := abciResponses.EndBlock.ValidatorUpdates
	err = validateValidatorUpdates(abciValUpdates, state.ConsensusParams.Validator)
	if err != nil {
		return state, 0, deltas, fmt.Errorf("error in validator updates: %v", err)
	}
	validatorUpdates, err := types.PB2TM.ValidatorUpdates(abciValUpdates)
	if err != nil {
		return state, 0, deltas, err
	}
	if len(validatorUpdates) > 0 {
		blockExec.logger.Info("Updates to validators", "updates", types.ValidatorListString(validatorUpdates))
	}

	// Update the state with the block and responses.
	state, err = updateState(state, blockID, &block.Header, abciResponses, validatorUpdates)
	if err != nil {
		return state, 0, deltas, fmt.Errorf("commit failed for application: %v", err)
	}

	trc.Pin(trace.Persist)
	startTime = time.Now().UnixNano()

	// Lock mempool, commit app state, update mempoool.
	appHash, retainHeight, err := blockExec.Commit(state, block, abciResponses.DeliverTxs, deltas)
	endTime = time.Now().UnixNano()
	blockExec.metrics.CommitTime.Set(float64(endTime-startTime) / 1e6)
	if err != nil {
		return state, 0, deltas, fmt.Errorf("commit failed for application: %v", err)
	}

	if !useDeltas {
		// get deliverTx WatchData and let wd = it
		deltas.WatchBytes = GetWatchData()
	} else {
		// commitBatch with wd in exchain
		UseWatchData(deltas.WatchBytes)
	}

	trc.Pin("evpool")
	// Update evpool with the block and state.
	blockExec.evpool.Update(block, state)

	fail.Fail() // XXX

	trc.Pin(trace.SaveState)

	// Update the app hash and save the state.
	state.AppHash = appHash
	SaveState(blockExec.db, state)
	blockExec.logger.Debug("SaveState", "state", fmt.Sprintf("%+v", state))
	fail.Fail() // XXX

	// Events are fired after everything else.
	// NOTE: if we crash between Commit and Save, events wont be fired during replay
	fireEvents(blockExec.logger, blockExec.eventBus, block, abciResponses, validatorUpdates)

	if broadDelta || uploadDelta {
		deltas.Height = block.Height
	}
	if types.EnableUploadDelta() {
		go blockExec.uploadData(block, deltas)
	}

	blockExec.logger.Info("Begin abci", "len(deltas)", deltas.Size(),
		"applyDelta", applyDelta, "downloadDelta", downloadDelta, "uploadDelta", uploadDelta, "broadDelta", broadDelta,
		"fastQuery", fastQuery, "FlagUseDelta", useDeltas)

	return state, retainHeight, deltas, nil
}

func (blockExec *BlockExecutor) uploadData(block *types.Block, deltas *types.Deltas) {
	if err := blockExec.deltaBroker.SetDeltas(deltas); err != nil {
		blockExec.logger.Error("uploadData err:", err)
		return
	}
	blockExec.logger.Info("uploadData",
		"height", block.Height,
		"blockLen", block.Size(),
		"abciRspLen", len(deltas.ABCIRsp),
		"deltaLen", len(deltas.DeltasBytes),
		"watchLen", len(deltas.WatchBytes))
}

func (blockExec *BlockExecutor) prepareStateDelta(block *types.Block, deltas *types.Deltas) (bool, *types.Deltas) {
	fastQuery := types.IsFastQuery()
	applyDelta := types.EnableApplyP2PDelta()
	downloadDelta := types.EnableDownloadDelta()

	// not use delta, exe abci itself
	if !applyDelta && !downloadDelta {
		return false, deltas
	}

	// get watchData and Delta from p2p
	if applyDelta {
		if len(deltas.ABCIRsp) > 0 && len(deltas.DeltasBytes) > 0 {
			if !fastQuery || len(deltas.WatchBytes) > 0 {
				return true, deltas
			}
		}
	}

	if !downloadDelta {
		return false, deltas
	}

	var directDelta *types.Deltas
	var err error
	needDDS := true
	select {
	case directDelta = <-blockExec.deltaCh:
		if directDelta.Height == block.Height {
			needDDS = false
		}
		// already get delta of height, then request delta of height+1
		blockExec.deltaHeightCh <- block.Height + 1
	default:
		// can't get delta of height, request delta of height+1 and return
		blockExec.deltaHeightCh <- block.Height + 1
	}

	if needDDS {
		// request watchData and Delta from dds
		directDelta, err = blockExec.deltaBroker.GetDeltas(block.Height)
	}

	// can't get data from dds
	if directDelta == nil {
		if err != nil {
			blockExec.logger.Error("Download Delta err:", err)
		}
		return false, deltas
	}

	//// get watchData from dds
	if !fastQuery || len(directDelta.WatchBytes) > 0 {
		// get Delta from dds
		if len(directDelta.ABCIRsp) > 0 && len(directDelta.DeltasBytes) > 0 {
			return true, directDelta
		}
		// get Delta from p2p
		if len(deltas.ABCIRsp) > 0 && len(deltas.DeltasBytes) > 0 {
			deltas.WatchBytes = directDelta.WatchBytes
			return true, deltas
		}
		// can't get Delta
		return false, deltas
	}

	//// can't get watchData from dds
	{
		if len(deltas.WatchBytes) <= 0 {
			// can't get watchData
			return false, deltas
		}

		// get Delta from dds
		if len(directDelta.ABCIRsp) > 0 && len(directDelta.DeltasBytes) > 0 {
			directDelta.WatchBytes = deltas.WatchBytes
			return true, directDelta
		}
	}

	return false, deltas
}

func (blockExec *BlockExecutor) GetDeltaFromDDS() {
	flag := false
	var height int64 = 0
	tryGetDDSTicker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-tryGetDDSTicker.C:
			if flag {
				directDelta, _ := blockExec.deltaBroker.GetDeltas(height)
				if directDelta != nil {
					flag = false
					blockExec.deltaCh <- directDelta
				}
			}

		case height = <-blockExec.deltaHeightCh:
			flag = true
		}
	}
}

// Commit locks the mempool, runs the ABCI Commit message, and updates the
// mempool.
// It returns the result of calling abci.Commit (the AppHash) and the height to retain (if any).
// The Mempool must be locked during commit and update because state is
// typically reset on Commit and old txs must be replayed against committed
// state before new txs are run in the mempool, lest they be invalid.
func (blockExec *BlockExecutor) Commit(
	state State,
	block *types.Block,
	deliverTxResponses []*abci.ResponseDeliverTx,
	deltas *types.Deltas,
) ([]byte, int64, error) {
	blockExec.mempool.Lock()
	defer func() {
		blockExec.mempool.Unlock()
		// Forced flushing mempool
		if cfg.DynamicConfig.GetMempoolFlush() {
			blockExec.mempool.Flush()
		}
	}()

	// while mempool is Locked, flush to ensure all async requests have completed
	// in the ABCI app before Commit.
	err := blockExec.mempool.FlushAppConn()
	if err != nil {
		blockExec.logger.Error("Client error during mempool.FlushAppConn", "err", err)
		return nil, 0, err
	}

	// Commit block, get hash back
	res, err := blockExec.proxyApp.CommitSync(abci.RequestCommit{Deltas: &abci.Deltas{DeltasByte: deltas.DeltasBytes}})
	if err != nil {
		blockExec.logger.Error(
			"Client error during proxyAppConn.CommitSync",
			"err", err,
		)
		return nil, 0, err
	}
	if res.Deltas == nil {
		res.Deltas = &abci.Deltas{}
	}

	// ResponseCommit has no error code - just data

	blockExec.logger.Info(
		"Committed state",
		"height", block.Height,
		"txs", len(block.Txs),
		"appHash", fmt.Sprintf("%X", res.Data),
		"blockLen", block.Size(),
		"inDeltasLen", len(deltas.DeltasBytes),
		"outDeltasLen", len(res.Deltas.DeltasByte),
	)

	deltas.DeltasBytes = res.Deltas.DeltasByte

	// Update mempool.
	err = blockExec.mempool.Update(
		block.Height,
		block.Txs,
		deliverTxResponses,
		TxPreCheck(state),
		TxPostCheck(state),
	)

	if !cfg.DynamicConfig.GetMempoolRecheck() && block.Height%cfg.DynamicConfig.GetMempoolForceRecheckGap() == 0 {
		proxyCb := func(req *abci.Request, res *abci.Response) {

		}
		blockExec.proxyApp.SetResponseCallback(proxyCb)
		// reset checkState
		blockExec.proxyApp.SetOptionAsync(abci.RequestSetOption{
			Key: "ResetCheckState",
		})
	}

	return res.Data, res.RetainHeight, err
}

func transTxsToBytes(txs types.Txs) [][]byte {
	ret := make([][]byte, 0)
	for _, v := range txs {
		ret = append(ret, v)
	}
	return ret
}

//---------------------------------------------------------
// Helper functions for executing blocks and updating state

// Executes block's transactions on proxyAppConn.
// Returns a list of transaction results and updates to the validator set
func execBlockOnProxyApp(context *executionContext) (*ABCIResponses, error) {
	block := context.block
	proxyAppConn := context.proxyApp
	stateDB := context.db
	logger := context.logger

	var validTxs, invalidTxs = 0, 0

	txIndex := 0
	abciResponses := NewABCIResponses(block)

	// Execute transactions and get hash.
	proxyCb := func(req *abci.Request, res *abci.Response) {
		if r, ok := res.Value.(*abci.Response_DeliverTx); ok {
			// TODO: make use of res.Log
			// TODO: make use of this info
			// Blocks may include invalid txs.
			txRes := r.DeliverTx
			if txRes.Code == abci.CodeTypeOK {
				validTxs++
			} else {
				logger.Debug("Invalid tx", "code", txRes.Code, "log", txRes.Log)
				invalidTxs++
			}
			abciResponses.DeliverTxs[txIndex] = txRes
			txIndex++
		}
	}
	proxyAppConn.SetResponseCallback(proxyCb)

	commitInfo, byzVals := getBeginBlockValidatorInfo(block, stateDB)

	// Begin block
	var err error
	abciResponses.BeginBlock, err = proxyAppConn.BeginBlockSync(abci.RequestBeginBlock{
		Hash:                block.Hash(),
		Header:              types.TM2PB.Header(&block.Header),
		LastCommitInfo:      commitInfo,
		ByzantineValidators: byzVals,
	})
	if err != nil {
		logger.Error("Error in proxyAppConn.BeginBlock", "err", err)
		return nil, err
	}

	// Run txs of block.
	for count, tx := range block.Txs {
		proxyAppConn.DeliverTxAsync(abci.RequestDeliverTx{Tx: tx})
		if err := proxyAppConn.Error(); err != nil {
			return nil, err
		}

		if context != nil && context.stopped {
			context.dump(fmt.Sprintf("Prerun stopped, %d/%d tx executed", count+1, len(block.Txs)))
			return nil, fmt.Errorf("Prerun stopped")
		}
	}

	// End block.
	abciResponses.EndBlock, err = proxyAppConn.EndBlockSync(abci.RequestEndBlock{Height: block.Height})
	if err != nil {
		logger.Error("Error in proxyAppConn.EndBlock", "err", err)
		return nil, err
	}

	logger.Info("Executed block", "height", block.Height, "validTxs", validTxs, "invalidTxs", invalidTxs)
	trace.GetElapsedInfo().AddInfo(trace.InvalidTxs, fmt.Sprintf("%d", invalidTxs))

	return abciResponses, nil
}

func execBlockOnProxyAppWithDeltas(
	proxyAppConn proxy.AppConnConsensus,
	block *types.Block,
	stateDB dbm.DB,
) {
	proxyCb := func(req *abci.Request, res *abci.Response) {
	}
	proxyAppConn.SetResponseCallback(proxyCb)

	commitInfo, byzVals := getBeginBlockValidatorInfo(block, stateDB)
	_, _ = proxyAppConn.BeginBlockSync(abci.RequestBeginBlock{
		Hash:                block.Hash(),
		Header:              types.TM2PB.Header(&block.Header),
		LastCommitInfo:      commitInfo,
		ByzantineValidators: byzVals,
		UseDeltas:           true,
	})
}

func getBeginBlockValidatorInfo(block *types.Block, stateDB dbm.DB) (abci.LastCommitInfo, []abci.Evidence) {
	voteInfos := make([]abci.VoteInfo, block.LastCommit.Size())
	// block.Height=1 -> LastCommitInfo.Votes are empty.
	// Remember that the first LastCommit is intentionally empty, so it makes
	// sense for LastCommitInfo.Votes to also be empty.
	if block.Height > types.GetStartBlockHeight()+1 {
		lastValSet, err := LoadValidators(stateDB, block.Height-1)
		if err != nil {
			panic(err)
		}

		// Sanity check that commit size matches validator set size - only applies
		// after first block.
		var (
			commitSize = block.LastCommit.Size()
			valSetLen  = len(lastValSet.Validators)
		)
		if commitSize != valSetLen {
			panic(fmt.Sprintf("commit size (%d) doesn't match valset length (%d) at height %d\n\n%v\n\n%v",
				commitSize, valSetLen, block.Height, block.LastCommit.Signatures, lastValSet.Validators))
		}

		for i, val := range lastValSet.Validators {
			commitSig := block.LastCommit.Signatures[i]
			voteInfos[i] = abci.VoteInfo{
				Validator:       types.TM2PB.Validator(val),
				SignedLastBlock: !commitSig.Absent(),
			}
		}
	}

	byzVals := make([]abci.Evidence, len(block.Evidence.Evidence))
	for i, ev := range block.Evidence.Evidence {
		// We need the validator set. We already did this in validateBlock.
		// TODO: Should we instead cache the valset in the evidence itself and add
		// `SetValidatorSet()` and `ToABCI` methods ?
		valset, err := LoadValidators(stateDB, ev.Height())
		if err != nil {
			panic(err)
		}
		byzVals[i] = types.TM2PB.Evidence(ev, valset, block.Time)
	}

	return abci.LastCommitInfo{
		Round: int32(block.LastCommit.Round),
		Votes: voteInfos,
	}, byzVals
}

func validateValidatorUpdates(abciUpdates []abci.ValidatorUpdate,
	params types.ValidatorParams) error {
	for _, valUpdate := range abciUpdates {
		if valUpdate.GetPower() < 0 {
			return fmt.Errorf("voting power can't be negative %v", valUpdate)
		} else if valUpdate.GetPower() == 0 {
			// continue, since this is deleting the validator, and thus there is no
			// pubkey to check
			continue
		}

		// Check if validator's pubkey matches an ABCI type in the consensus params
		thisKeyType := valUpdate.PubKey.Type
		if !params.IsValidPubkeyType(thisKeyType) {
			return fmt.Errorf("validator %v is using pubkey %s, which is unsupported for consensus",
				valUpdate, thisKeyType)
		}
	}
	return nil
}

// updateState returns a new State updated according to the header and responses.
func updateState(
	state State,
	blockID types.BlockID,
	header *types.Header,
	abciResponses *ABCIResponses,
	validatorUpdates []*types.Validator,
) (State, error) {

	// Copy the valset so we can apply changes from EndBlock
	// and update s.LastValidators and s.Validators.
	nValSet := state.NextValidators.Copy()

	// Update the validator set with the latest abciResponses.
	lastHeightValsChanged := state.LastHeightValidatorsChanged
	if len(validatorUpdates) > 0 {
		err := nValSet.UpdateWithChangeSet(validatorUpdates)
		if err != nil {
			return state, fmt.Errorf("error changing validator set: %v", err)
		}
		// Change results from this height but only applies to the next next height.
		lastHeightValsChanged = header.Height + 1 + 1
	}

	// Update validator proposer priority and set state variables.
	nValSet.IncrementProposerPriority(1)

	// Update the params with the latest abciResponses.
	nextParams := state.ConsensusParams
	lastHeightParamsChanged := state.LastHeightConsensusParamsChanged
	if abciResponses.EndBlock.ConsensusParamUpdates != nil {
		// NOTE: must not mutate s.ConsensusParams
		nextParams = state.ConsensusParams.Update(abciResponses.EndBlock.ConsensusParamUpdates)
		err := nextParams.Validate()
		if err != nil {
			return state, fmt.Errorf("error updating consensus params: %v", err)
		}
		// Change results from this height but only applies to the next height.
		lastHeightParamsChanged = header.Height + 1
	}

	// TODO: allow app to upgrade version
	nextVersion := state.Version

	// NOTE: the AppHash has not been populated.
	// It will be filled on state.Save.
	return State{
		Version:                          nextVersion,
		ChainID:                          state.ChainID,
		LastBlockHeight:                  header.Height,
		LastBlockID:                      blockID,
		LastBlockTime:                    header.Time,
		NextValidators:                   nValSet,
		Validators:                       state.NextValidators.Copy(),
		LastValidators:                   state.Validators.Copy(),
		LastHeightValidatorsChanged:      lastHeightValsChanged,
		ConsensusParams:                  nextParams,
		LastHeightConsensusParamsChanged: lastHeightParamsChanged,
		LastResultsHash:                  abciResponses.ResultsHash(),
		AppHash:                          nil,
	}, nil
}

// Fire NewBlock, NewBlockHeader.
// Fire TxEvent for every tx.
// NOTE: if Tendermint crashes before commit, some or all of these events may be published again.
func fireEvents(
	logger log.Logger,
	eventBus types.BlockEventPublisher,
	block *types.Block,
	abciResponses *ABCIResponses,
	validatorUpdates []*types.Validator,
) {
	eventBus.PublishEventNewBlock(types.EventDataNewBlock{
		Block:            block,
		ResultBeginBlock: *abciResponses.BeginBlock,
		ResultEndBlock:   *abciResponses.EndBlock,
	})
	eventBus.PublishEventNewBlockHeader(types.EventDataNewBlockHeader{
		Header:           block.Header,
		NumTxs:           int64(len(block.Txs)),
		ResultBeginBlock: *abciResponses.BeginBlock,
		ResultEndBlock:   *abciResponses.EndBlock,
	})

	for i, tx := range block.Data.Txs {
		eventBus.PublishEventTx(types.EventDataTx{TxResult: types.TxResult{
			Height: block.Height,
			Index:  uint32(i),
			Tx:     tx,
			Result: *(abciResponses.DeliverTxs[i]),
		}})
	}

	if len(validatorUpdates) > 0 {
		eventBus.PublishEventValidatorSetUpdates(
			types.EventDataValidatorSetUpdates{ValidatorUpdates: validatorUpdates})
	}
}

//----------------------------------------------------------------------------------------------------
// Execute block without state. TODO: eliminate

// ExecCommitBlock executes and commits a block on the proxyApp without validating or mutating the state.
// It returns the application root hash (result of abci.Commit).
func ExecCommitBlock(
	appConnConsensus proxy.AppConnConsensus,
	block *types.Block,
	logger log.Logger,
	stateDB dbm.DB,
) ([]byte, error) {

	ctx := &executionContext{
		logger: logger,
		block: block,
		db: stateDB,
		proxyApp: appConnConsensus,
	}

	_, err := execBlockOnProxyApp(ctx)
	if err != nil {
		logger.Error("Error executing block on proxy app", "height", block.Height, "err", err)
		return nil, err
	}
	// Commit block, get hash back
	res, err := appConnConsensus.CommitSync(abci.RequestCommit{})
	if err != nil {
		logger.Error("Client error during proxyAppConn.CommitSync", "err", res)
		return nil, err
	}
	// ResponseCommit has no error or log, just data
	return res.Data, nil
}
