package watcher

import (
	"math/big"
	"time"

	clientcontext "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

// NewTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func NewTransaction(tx *evmtypes.MsgEthereumTx, txHash, blockHash common.Hash, blockNumber, index uint64) (*Transaction, error) {
	// Verify signature and retrieve sender address
	err := tx.VerifySig(tx.ChainID(), int64(blockNumber))
	if err != nil {
		return nil, err
	}

	rpcTx := &Transaction{
		From:     common.HexToAddress(tx.GetFrom()),
		Gas:      hexutil.Uint64(tx.Data.GasLimit),
		GasPrice: (*hexutil.Big)(tx.Data.Price),
		Hash:     txHash,
		Input:    hexutil.Bytes(tx.Data.Payload),
		Nonce:    hexutil.Uint64(tx.Data.AccountNonce),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Data.Amount),
		V:        (*hexutil.Big)(tx.Data.V),
		R:        (*hexutil.Big)(tx.Data.R),
		S:        (*hexutil.Big)(tx.Data.S),
	}

	if blockHash != (common.Hash{}) {
		rpcTx.BlockHash = &blockHash
		rpcTx.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		rpcTx.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	return rpcTx, nil
}

func RawTxResultToStdResponse(clientCtx clientcontext.CLIContext,
	tr *ctypes.ResultTx, tx sdk.Tx, timestamp time.Time) (*TransactionResult, error) {

	var err error
	if tx == nil {
		tx, err = evmtypes.TxDecoder(clientCtx.CodecProy)(tr.Tx, tr.Height)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}
	}

	var realTx *authtypes.StdTx
	switch tx.(type) {
	case *authtypes.IbcTx:
		realTx, err = authtypes.FromProtobufTx(clientCtx.CodecProy, tx.(*authtypes.IbcTx))
		if nil != err {
			return nil, err
		}
	default:
		err = clientCtx.Codec.UnmarshalBinaryLengthPrefixed(tr.Tx, &realTx)
		if err != nil {
			return nil, err
		}
	}

	response := sdk.NewResponseResultTx(tr, realTx, timestamp.Format(time.RFC3339))
	wrappedR := &WrappedResponseWithCodec{Response: response, Codec: clientCtx.Codec}

	return &TransactionResult{TxType: hexutil.Uint64(StdResponse), Response: wrappedR}, nil
}
