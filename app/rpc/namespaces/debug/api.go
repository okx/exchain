package debug

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/exchain/app/rpc/namespaces/eth/simulation"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	evmtypes "github.com/okex/exchain/x/evm/types"

	"github.com/okex/exchain/app/rpc/backend"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

// PublicTxPoolAPI offers and API for the transaction pool. It only operates on data that is non confidential.
type PublicDebugAPI struct {
	clientCtx  clientcontext.CLIContext
	logger     log.Logger
	backend    backend.Backend
	evmFactory simulation.EvmFactory
}

// NewPublicTxPoolAPI creates a new tx pool service that gives information about the transaction pool.
func NewAPI(clientCtx clientcontext.CLIContext, log log.Logger, backend backend.Backend) *PublicDebugAPI {
	api := &PublicDebugAPI{
		clientCtx: clientCtx,
		backend:   backend,
		logger:    log.With("module", "json-rpc", "namespace", "debug"),
	}
	return api
}

// GetAccount returns the provided account's balance up to the provided block number.
func (api *PublicDebugAPI) GetAccount(address common.Address) (*ethermint.EthAccount, error) {
	clientCtx := api.clientCtx

	bs, err := api.clientCtx.Codec.MarshalJSON(auth.NewQueryAccountParams(address.Bytes()))
	if err != nil {
		return nil, err
	}

	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs)
	if err != nil {
		return nil, err
	}

	var account ethermint.EthAccount
	if err := api.clientCtx.Codec.UnmarshalJSON(res, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (api *PublicDebugAPI) GetCodeByHash(hash common.Hash) (hexutil.Bytes, error) {

	clientCtx := api.clientCtx
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryCodeByHash, hash.Hex()), nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResCode
	api.clientCtx.Codec.MustUnmarshalJSON(res, &out)

	return out.Code, nil
}

// GetStorageAtInternal returns the contract storage at the given address, block number, and key.
func (api *PublicDebugAPI) GetStorageAtInternal(address common.Address, key []byte) (hexutil.Bytes, error) {
	return api.getStorageAt(address, key, 0, true)
}
func (api *PublicDebugAPI) getStorageAt(address common.Address, key []byte, blockNum rpctypes.BlockNumber, directlyKey bool) (hexutil.Bytes, error) {
	clientCtx := api.clientCtx.WithHeight(blockNum.Int64())

	var queryStr = ""
	if !directlyKey {
		queryStr = fmt.Sprintf("custom/%s/storage/%s/%X", evmtypes.ModuleName, address.Hex(), key)
	} else {
		queryStr = fmt.Sprintf("custom/%s/storageKey/%s/%X", evmtypes.ModuleName, address.Hex(), key)
	}

	res, _, err := clientCtx.QueryWithData(queryStr, nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResStorage
	api.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Value, nil
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *PublicDebugAPI) TraceTransaction(txHash common.Hash) (hexutil.Bytes, error) {
	// Get transaction by hash
	tx, err := api.clientCtx.Client.Tx(txHash.Bytes(), false)
	if err != nil {
		//to keep consistent with rpc of ethereum, should be return nil
		return nil, nil
	}

	// check if block number is 0
	if tx.Height == 0 {
		return nil, errors.New("genesis is not traceable")
	}

	// Can either cache or just leave this out if not necessary
	res, err := api.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}
	block := res.Block
	// check tx index is not out of bound
	if uint32(len(block.Txs)) < tx.Index {
		return nil, fmt.Errorf("transaction not included in block %v", block.Height)
	}

	var predecessors []*evmtypes.MsgEthereumTx
	for _, txBz := range block.Txs[:tx.Index] {
		tx, err := rpctypes.RawTxToEthTx(api.clientCtx, txBz)
		if err != nil {
			return nil, err
		}
		predecessors = append(predecessors, tx)
	}

	ethMessage, err := rpctypes.RawTxToEthTx(api.clientCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	// minus one to get the context of block beginning
	contextHeight := tx.Height - 1
	if contextHeight < 1 {
		// 0 is a special value in `ContextWithHeight`
		contextHeight = 1
	}

	sim := api.evmFactory.BuildSimulatorForSpecificBlock(api, res.BlockID, res.Block)
	if sim == nil {
		return nil, err
	}
	resTrace, err := sim.TraceTx(ethMessage, predecessors)
	if err != nil {
		return nil, err
	}
	return resTrace.Data, nil
}
