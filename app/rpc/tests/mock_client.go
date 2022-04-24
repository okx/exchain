package tests

import (
	"crypto/sha256"
	"fmt"
	"time"

	"container/list"

	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	apptesting "github.com/okex/exchain/libs/ibc-go/testing"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmcfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/bytes"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
	"github.com/okex/exchain/libs/tendermint/mempool"
	mempl "github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/rpc/client"
	"github.com/okex/exchain/libs/tendermint/rpc/client/mock"
	rpccore "github.com/okex/exchain/libs/tendermint/rpc/core"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	sm "github.com/okex/exchain/libs/tendermint/state"
	tmstate "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/state/txindex/kv"
	"github.com/okex/exchain/libs/tendermint/state/txindex/null"
	"github.com/okex/exchain/libs/tendermint/store"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
)

type MockClient struct {
	mock.Client
	app          abci.Application
	chain        apptesting.TestChainI
	env          *rpccore.Environment
	state        tmstate.State
	blockExec    *sm.BlockExecutor
	priv         types.PrivValidator
	responseList list.List
}

func createAndStartProxyAppConns(clientCreator proxy.ClientCreator, logger log.Logger) (proxy.AppConns, error) {
	proxyApp := proxy.NewAppConns(clientCreator)
	proxyApp.SetLogger(logger.With("module", "proxy"))
	if err := proxyApp.Start(); err != nil {
		return nil, fmt.Errorf("error starting proxy app connections: %v", err)
	}
	return proxyApp, nil
}

func NewMockClient(chainId string, chain apptesting.TestChainI, app abci.Application) *MockClient {
	config := tmcfg.ResetTestRootWithChainID("blockchain_reactor_test", chainId)

	papp := proxy.NewLocalClientCreator(app)
	proxyApp, err := createAndStartProxyAppConns(papp, log.NewNopLogger())
	if err != nil {
		panic(err)
	}

	mc := &MockClient{
		app:   app,
		chain: chain,
		env: &rpccore.Environment{
			BlockStore: store.NewBlockStore(dbm.NewMemDB()),
			StateDB:    dbm.NewMemDB(),
			TxIndexer:  kv.NewTxIndex(dbm.NewMemDB()),
		},
	}
	mc.state, err = tmstate.LoadStateFromDBOrGenesisFile(mc.env.StateDB, config.GenesisFile())
	if err != nil {
		panic(err)
	}
	mempool := mempool.NewCListMempool(
		config.Mempool,
		proxyApp.Mempool(),
		mc.state.LastBlockHeight,
	)
	mc.env.Mempool = mempool
	mc.env.PubKey = chain.SenderAccount().GetPubKey()

	db := dbm.NewMemDB()
	mc.blockExec = sm.NewBlockExecutor(db, log.NewNopLogger(), proxyApp.Consensus(), mempool, sm.MockEvidencePool{})
	sm.SaveState(db, mc.state)
	return mc
}
func (c MockClient) makeBlock(height int64, state sm.State, lastCommit *types.Commit) *types.Block {
	tx := c.env.Mempool.ReapMaxTxs(1000)
	block, _ := state.MakeBlock(height, tx, lastCommit, nil, state.Validators.GetProposer().Address)
	return block
}
func (c *MockClient) CommitBlock() {
	//blockHeight := c.state.LastBlockHeight + 1
	if c.priv == nil {
		_, c.priv = types.RandValidator(false, 30)
	}
	blockHeight := c.state.LastBlockHeight + 1
	lastCommit := types.NewCommit(blockHeight-1, 0, types.BlockID{}, nil)
	if blockHeight > 1 {
		lastBlockMeta := c.env.BlockStore.LoadBlockMeta(blockHeight - 1)
		lastBlock := c.env.BlockStore.LoadBlock(blockHeight - 1)

		vote, err := types.MakeVote(
			lastBlock.Header.Height,
			lastBlockMeta.BlockID,
			c.state.Validators,
			c.priv,
			lastBlock.Header.ChainID,
			time.Now(),
		)
		if err != nil {
			panic(err)
		}
		lastCommit = types.NewCommit(vote.Height, vote.Round,
			lastBlockMeta.BlockID, []types.CommitSig{vote.CommitSig()})
	}

	thisBlock := c.makeBlock(blockHeight, c.state, lastCommit)

	thisParts := thisBlock.MakePartSet(types.BlockPartSizeBytes)
	blockID := types.BlockID{Hash: thisBlock.Hash(), PartsHeader: thisParts.Header()}

	//thisBlock.Height = state.LastBlockHeight + 1
	//thisBlock.LastCommit.Signatures = []types.CommitSig{}
	//thisBlock.Time = state.LastBlockTime
	state, _, err := c.blockExec.ApplyBlock(c.state, blockID, thisBlock)
	if err != nil {
		panic(errors.Wrap(err, "error apply block"))
	}
	//thisBlock.Height = state.LastBlockHeight + 1
	c.env.BlockStore.SaveBlock(thisBlock, thisParts, lastCommit)
	c.state = state
}

func (c *MockClient) SetResponse(response interface{}, clear bool) {
	if clear {
		c.responseList = *list.New()
	}
	c.responseList.PushBack(response)

}
func (c MockClient) ABCIQueryWithOptions(
	path string,
	data bytes.HexBytes,
	opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	resQuery := c.app.Query(abci.RequestQuery{
		Path:   path,
		Data:   data,
		Height: opts.Height,
		Prove:  opts.Prove,
	})
	return &ctypes.ResultABCIQuery{Response: resQuery}, nil
}
func (c MockClient) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	resCh := make(chan *abci.Response, 1)
	err := c.env.Mempool.CheckTx(tx, func(res *abci.Response) {
		resCh <- res
	}, mempl.TxInfo{})
	if err != nil {
		return nil, err
	}
	res := <-resCh
	r := res.GetCheckTx()
	return &ctypes.ResultBroadcastTx{
		Code:      r.Code,
		Data:      r.Data,
		Log:       r.Log,
		Codespace: r.Codespace,
		Hash:      tx.Hash(c.env.BlockStore.Height()),
	}, nil
}

// error if either min or max are negative or min > max
// if 0, use blockstore base for min, latest block height for max
// enforce limit.
func filterMinMax(base, height, min, max, limit int64) (int64, int64, error) {
	// filter negatives
	if min < 0 || max < 0 {
		return min, max, fmt.Errorf("heights must be non-negative")
	}

	// adjust for default values
	if max == 0 {
		max = height
	}

	// limit max to the height
	max = tmmath.MinInt64(height, max)

	// limit min to the base
	min = tmmath.MaxInt64(base, min)

	// limit min to within `limit` of max
	// so the total number of blocks returned will be `limit`
	min = tmmath.MaxInt64(min, max-limit+1)

	if min > max {
		return min, max, fmt.Errorf("min height %d can't be greater than max height %d", min, max)
	}
	return min, max, nil
}
func (c MockClient) Status() (*ctypes.ResultStatus, error) {
	var (
		earliestBlockHash     tmbytes.HexBytes
		earliestAppHash       tmbytes.HexBytes
		earliestBlockTimeNano int64

		earliestBlockHeight = c.env.BlockStore.Base()
	)

	if earliestBlockMeta := c.env.BlockStore.LoadBlockMeta(earliestBlockHeight); earliestBlockMeta != nil {
		earliestAppHash = earliestBlockMeta.Header.AppHash
		earliestBlockHash = earliestBlockMeta.BlockID.Hash
		earliestBlockTimeNano = earliestBlockMeta.Header.Time.UnixNano()
	}

	var (
		latestBlockHash     tmbytes.HexBytes
		latestAppHash       tmbytes.HexBytes
		latestBlockTimeNano int64

		latestHeight = c.env.BlockStore.Height()
	)

	if latestHeight != 0 {
		latestBlockMeta := c.env.BlockStore.LoadBlockMeta(latestHeight)
		if latestBlockMeta != nil {
			latestBlockHash = latestBlockMeta.BlockID.Hash
			latestAppHash = latestBlockMeta.Header.AppHash
			latestBlockTimeNano = latestBlockMeta.Header.Time.UnixNano()
		}
	}

	// Return the very last voting power, not the voting power of this validator
	// during the last block.
	var votingPower int64
	blockHeight := c.env.BlockStore.Height() + 1
	if val := c.validatorAtHeight(blockHeight); val != nil {
		votingPower = val.VotingPower
	}

	result := &ctypes.ResultStatus{
		//NodeInfo: c.env.P2PTransport.NodeInfo().(p2p.DefaultNodeInfo),
		SyncInfo: ctypes.SyncInfo{
			LatestBlockHash:     latestBlockHash,
			LatestAppHash:       latestAppHash,
			LatestBlockHeight:   latestHeight,
			LatestBlockTime:     time.Unix(0, latestBlockTimeNano),
			EarliestBlockHash:   earliestBlockHash,
			EarliestAppHash:     earliestAppHash,
			EarliestBlockHeight: earliestBlockHeight,
			EarliestBlockTime:   time.Unix(0, earliestBlockTimeNano),
			//CatchingUp:          c.env.ConsensusReactor.FastSync(),
		},
		ValidatorInfo: ctypes.ValidatorInfo{
			Address:     c.env.PubKey.Address(),
			PubKey:      c.env.PubKey,
			VotingPower: votingPower,
		},
	}

	return result, nil
}
func (c MockClient) validatorAtHeight(h int64) *types.Validator {
	vals, err := sm.LoadValidators(c.env.StateDB, h)
	if err != nil {
		return nil
	}
	_, val := vals.GetByIndex(0)
	return val
}

// latestHeight can be either latest committed or uncommitted (+1) height.
func (c MockClient) getHeight(latestHeight int64, heightPtr *int64) (int64, error) {
	if heightPtr != nil {
		height := *heightPtr
		if height <= 0 {
			return 0, fmt.Errorf("height must be greater than 0, but got %d", height)
		}
		if height > latestHeight {
			return 0, fmt.Errorf("height %d must be less than or equal to the current blockchain height %d",
				height, latestHeight)
		}
		base := c.env.BlockStore.Base()
		if height < base {
			return 0, fmt.Errorf("height %v is not available, blocks pruned at height %v",
				height, base)
		}
		return height, nil
	}
	return latestHeight, nil
}
func (c *MockClient) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	const limit int64 = 20
	var err error
	minHeight, maxHeight, err = filterMinMax(
		c.env.BlockStore.Base(),
		c.env.BlockStore.Height(),
		minHeight,
		maxHeight,
		limit)
	if err != nil {
		return nil, err
	}
	blockMetas := []*types.BlockMeta{}
	for height := maxHeight; height >= minHeight; height-- {
		blockMeta := c.env.BlockStore.LoadBlockMeta(height)
		blockMetas = append(blockMetas, blockMeta)
	}
	return &ctypes.ResultBlockchainInfo{
		LastHeight: c.env.BlockStore.Height(),
		BlockMetas: blockMetas}, nil
}
func (c *MockClient) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return &ctypes.ResultUnconfirmedTxs{
		Count:      c.env.Mempool.Size(),
		Total:      c.env.Mempool.Size(),
		TotalBytes: c.env.Mempool.TxsBytes()}, nil
}
func (c *MockClient) Block(heightPtr *int64) (*ctypes.ResultBlock, error) {
	if c.responseList.Len() != 0 {
		res, ok := c.responseList.Front().Value.(*ctypes.ResultBlock)
		if ok {
			c.responseList.Remove(c.responseList.Front())
			return res, nil
		}
	}
	height, err := c.getHeight(c.env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	block := c.env.BlockStore.LoadBlock(height)
	blockMeta := c.env.BlockStore.LoadBlockMeta(height)
	if blockMeta == nil {
		return &ctypes.ResultBlock{BlockID: types.BlockID{}, Block: block}, nil
	}
	return &ctypes.ResultBlock{BlockID: blockMeta.BlockID, Block: block}, nil
}
func (c *MockClient) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	if c.responseList.Len() != 0 {
		res, ok := c.responseList.Front().Value.(*ctypes.ResultTx)
		if ok {
			c.responseList.Remove(c.responseList.Front())
			return res, nil
		}
	}
	// if index is disabled, return error
	if _, ok := c.env.TxIndexer.(*null.TxIndex); ok {
		return nil, fmt.Errorf("transaction indexing is disabled")
	}

	r, err := c.env.TxIndexer.Get(hash)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, fmt.Errorf("tx (%X) not found", hash)
	}

	height := r.Height
	index := r.Index

	var proof types.TxProof
	if prove {
		block := c.env.BlockStore.LoadBlock(height)
		proof = block.Data.Txs.Proof(int(index), block.Height) // XXX: overflow on 32-bit machines
	}

	return &ctypes.ResultTx{
		Hash:     hash,
		Height:   height,
		Index:    index,
		TxResult: r.Result,
		Tx:       r.Tx,
		Proof:    proof,
	}, nil
}
func (c *MockClient) GetAddressList() (*ctypes.ResultUnconfirmedAddresses, error) {
	addressList := c.env.Mempool.GetAddressList()
	return &ctypes.ResultUnconfirmedAddresses{
		Addresses: addressList,
	}, nil
}
func (c *MockClient) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	txs := c.env.Mempool.ReapMaxTxs(limit)
	return &ctypes.ResultUnconfirmedTxs{
		Count:      len(txs),
		Total:      c.env.Mempool.Size(),
		TotalBytes: c.env.Mempool.TxsBytes(),
		Txs:        txs}, nil
}
func (c MockClient) GetUnconfirmedTxByHash(hash [sha256.Size]byte) (types.Tx, error) {
	return c.env.Mempool.GetTxByHash(hash)
}
func (c *MockClient) UserUnconfirmedTxs(address string, limit int) (*ctypes.ResultUserUnconfirmedTxs, error) {
	txs := c.env.Mempool.ReapUserTxs(address, limit)
	return &ctypes.ResultUserUnconfirmedTxs{
		Count: len(txs),
		Txs:   txs}, nil
}
func (c MockClient) UserNumUnconfirmedTxs(address string) (*ctypes.ResultUserUnconfirmedTxs, error) {
	nums := c.env.Mempool.ReapUserTxsCnt(address)
	return &ctypes.ResultUserUnconfirmedTxs{
		Count: nums}, nil
}
func (c MockClient) GetPendingNonce(address string) (*ctypes.ResultPendingNonce, bool) {
	nonce, ok := c.env.Mempool.GetPendingNonce(address)
	if !ok {
		return nil, false
	}
	return &ctypes.ResultPendingNonce{
		Nonce: nonce,
	}, true
}
