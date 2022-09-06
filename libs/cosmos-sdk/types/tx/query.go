package tx

////
//import (
//	"encoding/hex"
//	"errors"
//	"strings"
//
//	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
//
//	"github.com/okex/exchain/libs/cosmos-sdk/client"
//
//	types2 "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
//
//	cliContext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
//	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
//	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
//)
//
//// QueryTxsByEvents performs a search for transactions for a given set of events
//// via the Tendermint RPC. An event takes the form of:
//// "{eventAttribute}.{attributeKey} = '{attributeValue}'". Each event is
//// concatenated with an 'AND' operand. It returns a slice of Info object
//// containing txs and metadata. An error is returned if the query fails.
//// If an empty string is provided it will order txs by asc
//func QueryTxsByEvents(clientCtx cliContext.CLIContext, events []string, page, limit int, orderBy string) (*types2.SearchTxsResult, error) {
//	if len(events) == 0 {
//		return nil, errors.New("must declare at least one event to search")
//	}
//
//	if page <= 0 {
//		return nil, errors.New("page must greater than 0")
//	}
//
//	if limit <= 0 {
//		return nil, errors.New("limit must greater than 0")
//	}
//
//	// XXX: implement ANY
//	query := strings.Join(events, " AND ")
//
//	node, err := clientCtx.GetNode()
//	if err != nil {
//		return nil, err
//	}
//
//	// TODO: this may not always need to be proven
//	// https://github.com/cosmos/cosmos-sdk/issues/6807
//	resTxs, err := node.TxSearch(query, true, page, limit, orderBy)
//	if err != nil {
//		return nil, err
//	}
//
//	resBlocks, err := getBlocksForTxResults(clientCtx, resTxs.Txs)
//	if err != nil {
//		return nil, err
//	}
//	pbtxCfg := utils.NewPbTxConfig(clientCtx.InterfaceRegistry)
//	txs, err := formatTxResults(pbtxCfg, resTxs.Txs, resBlocks)
//	if err != nil {
//		return nil, err
//	}
//
//	result := types2.NewSearchTxsResult(uint64(resTxs.TotalCount), uint64(len(txs)), uint64(page), uint64(limit), txs)
//
//	return result, nil
//}
//
//// QueryTx queries for a single transaction by a hash string in hex format. An
//// error is returned if the transaction does not exist or cannot be queried.
//func QueryTx(clientCtx cliContext.CLIContext, hashHexStr string) (*types2.TxResponse, error) {
//	hash, err := hex.DecodeString(hashHexStr)
//	if err != nil {
//		return nil, err
//	}
//
//	node, err := clientCtx.GetNode()
//	if err != nil {
//		return nil, err
//	}
//
//	//TODO: this may not always need to be proven
//	// https://github.com/cosmos/cosmos-sdk/issues/6807
//	resTx, err := node.Tx(hash, true)
//	if err != nil {
//		return nil, err
//	}
//
//	resBlocks, err := getBlocksForTxResults(clientCtx, []*ctypes.ResultTx{resTx})
//	if err != nil {
//		return nil, err
//	}
//
//	pbtxCfg := utils.NewPbTxConfig(clientCtx.InterfaceRegistry)
//	out, err := mkTxResult(pbtxCfg, resTx, resBlocks[resTx.Height])
//	if err != nil {
//		return out, err
//	}
//
//	return out, nil
//}
//
//// formatTxResults parses the indexed txs into a slice of TxResponse objects.
//func formatTxResults(txConfig client.TxConfig, resTxs []*ctypes.ResultTx, resBlocks map[int64]*ctypes.ResultBlock) ([]*types2.TxResponse, error) {
//	var err error
//	out := make([]*types2.TxResponse, len(resTxs))
//	for i := range resTxs {
//		out[i], err = mkTxResult(txConfig, resTxs[i], resBlocks[resTxs[i].Height])
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return out, nil
//}
//
//func getBlocksForTxResults(clientCtx cliContext.CLIContext, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
//	node, err := clientCtx.GetNode()
//	if err != nil {
//		return nil, err
//	}
//
//	resBlocks := make(map[int64]*ctypes.ResultBlock)
//
//	for _, resTx := range resTxs {
//		if _, ok := resBlocks[resTx.Height]; !ok {
//			resBlock, err := node.Block(&resTx.Height)
//			if err != nil {
//				return nil, err
//			}
//
//			resBlocks[resTx.Height] = resBlock
//		}
//	}
//
//	return resBlocks, nil
//}
//
//func mkTxResult(txConfig client.TxConfig, resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (*types2.TxResponse, error) {
//	return nil, nil
//	//txb, err := txConfig.TxDecoder()(resTx.Tx)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//p, ok := txb.(intoAny)
//	//if !ok {
//	//	return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
//	//}
//	//any := p.AsAny()
//	//return types2.NewResponseResultTx(resTx, any, resBlock.Block.Time.Format(time.RFC3339)), nil
//}
//
//// Deprecated: this interface is used only internally for scenario we are
//// deprecating (StdTxConfig support)
//type intoAny interface {
//	AsAny() *codectypes.Any
//}
