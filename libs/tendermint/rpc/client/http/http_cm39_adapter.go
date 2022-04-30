package http

import (
	"context"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	rpcclient "github.com/okex/exchain/libs/tendermint/rpc/client"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

var (
	_ rpcclient.Client = (*HttpCM39Adapter)(nil)
)

type HttpCM39Adapter struct {
	proxy *HTTP
}

func NewHttpCM39Adapter(proxy *HTTP) *HttpCM39Adapter {
	return &HttpCM39Adapter{proxy: proxy}
}

func (h *HttpCM39Adapter) Start() error {
	return h.proxy.Start()
}

func (h *HttpCM39Adapter) OnStart() error {
	return h.proxy.OnStart()
}

func (h *HttpCM39Adapter) Stop() error {
	return h.proxy.Stop()
}

func (h *HttpCM39Adapter) OnStop() {
	h.proxy.OnStop()
}

func (h *HttpCM39Adapter) Reset() error {
	return h.proxy.Reset()
}

func (h *HttpCM39Adapter) OnReset() error {
	return h.proxy.Reset()
}

func (h *HttpCM39Adapter) IsRunning() bool {
	return h.proxy.IsRunning()
}

func (h *HttpCM39Adapter) Quit() <-chan struct{} {
	return h.proxy.Quit()
}

func (h *HttpCM39Adapter) String() string {
	return h.proxy.String()
}

func (h *HttpCM39Adapter) SetLogger(logger log.Logger) {
	h.proxy.SetLogger(logger)
}

func (h *HttpCM39Adapter) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	return h.proxy.ABCIInfo()
}

func (h *HttpCM39Adapter) ABCIQuery(path string, data tmbytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	return h.proxy.ABCIQuery(path, data)
}

func (h *HttpCM39Adapter) ABCIQueryWithOptions(path string, data tmbytes.HexBytes, opts rpcclient.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	ret, err := h.proxy.ABCIQueryWithOptions(path, data, opts)
	if nil == err {
		return ret, err
	}
	return h.cm39ABCIQueryWithOptions(path, data, opts)
}

func (c *HttpCM39Adapter) cm39ABCIQueryWithOptions(
	path string,
	data tmbytes.HexBytes,
	opts rpcclient.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	result := new(ctypes.CM39ResultABCIQuery)
	_, err := c.proxy.caller.Call("abci_query",
		map[string]interface{}{"path": path, "data": data, "height": opts.Height, "prove": opts.Prove},
		result)
	if err != nil {
		return nil, errors.Wrap(err, "ABCIQuery")
	}
	return result.ToResultABCIQuery(), nil
}

func (h *HttpCM39Adapter) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return h.proxy.BroadcastTxCommit(tx)
}

func (h *HttpCM39Adapter) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return h.proxy.BroadcastTxAsync(tx)
}

func (h *HttpCM39Adapter) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return h.proxy.BroadcastTxSync(tx)
}

func (h *HttpCM39Adapter) Subscribe(ctx context.Context, subscriber, query string, outCapacity ...int) (out <-chan ctypes.ResultEvent, err error) {
	return h.proxy.Subscribe(ctx, subscriber, query, outCapacity...)
}

func (h *HttpCM39Adapter) Unsubscribe(ctx context.Context, subscriber, query string) error {
	return h.proxy.Unsubscribe(ctx, subscriber, query)
}

func (h *HttpCM39Adapter) UnsubscribeAll(ctx context.Context, subscriber string) error {
	return h.proxy.UnsubscribeAll(ctx, subscriber)
}

func (h *HttpCM39Adapter) Genesis() (*ctypes.ResultGenesis, error) {
	return h.proxy.Genesis()
}

func (h *HttpCM39Adapter) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	return h.proxy.BlockchainInfo(minHeight, maxHeight)
}

func (h *HttpCM39Adapter) NetInfo() (*ctypes.ResultNetInfo, error) {
	return h.proxy.NetInfo()
}

func (h *HttpCM39Adapter) DumpConsensusState() (*ctypes.ResultDumpConsensusState, error) {
	return h.proxy.DumpConsensusState()
}

func (h *HttpCM39Adapter) ConsensusState() (*ctypes.ResultConsensusState, error) {
	return h.proxy.ConsensusState()
}

func (h *HttpCM39Adapter) ConsensusParams(height *int64) (*ctypes.ResultConsensusParams, error) {
	return h.proxy.ConsensusParams(height)
}

func (h *HttpCM39Adapter) Health() (*ctypes.ResultHealth, error) {
	return h.proxy.Health()
}

func (h *HttpCM39Adapter) Block(height *int64) (*ctypes.ResultBlock, error) {
	return h.proxy.Block(height)
}

func (h *HttpCM39Adapter) BlockResults(height *int64) (*ctypes.ResultBlockResults, error) {
	return h.proxy.BlockResults(height)
}

func (h *HttpCM39Adapter) Commit(height *int64) (*ctypes.ResultCommit, error) {
	return h.proxy.Commit(height)
}

func (h *HttpCM39Adapter) Validators(height *int64, page, perPage int) (*ctypes.ResultValidators, error) {
	return h.proxy.Validators(height, page, perPage)
}

func (h *HttpCM39Adapter) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return h.proxy.Tx(hash, prove)
}

func (h *HttpCM39Adapter) TxSearch(query string, prove bool, page, perPage int, orderBy string) (*ctypes.ResultTxSearch, error) {
	return h.proxy.TxSearch(query, prove, page, perPage, orderBy)
}

func (h *HttpCM39Adapter) Status() (*ctypes.ResultStatus, error) {
	return h.proxy.Status()
}

func (h *HttpCM39Adapter) BroadcastEvidence(ev types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	return h.proxy.BroadcastEvidence(ev)
}

func (h *HttpCM39Adapter) UnconfirmedTxs(limit int) (*ctypes.ResultUnconfirmedTxs, error) {
	return h.proxy.UnconfirmedTxs(limit)
}

func (h *HttpCM39Adapter) NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	return h.proxy.NumUnconfirmedTxs()
}

func (h *HttpCM39Adapter) UserUnconfirmedTxs(address string, limit int) (*ctypes.ResultUserUnconfirmedTxs, error) {
	return h.proxy.UserUnconfirmedTxs(address, limit)
}

func (h *HttpCM39Adapter) UserNumUnconfirmedTxs(address string) (*ctypes.ResultUserUnconfirmedTxs, error) {
	return h.proxy.UserNumUnconfirmedTxs(address)
}

func (h *HttpCM39Adapter) GetUnconfirmedTxByHash(hash [32]byte) (types.Tx, error) {
	return h.proxy.GetUnconfirmedTxByHash(hash)
}

func (h *HttpCM39Adapter) GetAddressList() (*ctypes.ResultUnconfirmedAddresses, error) {
	return h.proxy.GetAddressList()
}

func (h *HttpCM39Adapter) GetPendingNonce(address string) (*ctypes.ResultPendingNonce, bool) {
	return h.proxy.GetPendingNonce(address)
}
