package txpool

import (
	"fmt"
	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/x/backend/types"

	"github.com/okex/exchain/app/rpc/backend"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
)

// PublicTxPoolAPI offers and API for the transaction pool. It only operates on data that is non confidential.
type PublicTxPoolAPI struct {
	clientCtx clientcontext.CLIContext
	logger    log.Logger
	backend backend.Backend
}

// NewPublicTxPoolAPI creates a new tx pool service that gives information about the transaction pool.
func NewAPI(clientCtx clientcontext.CLIContext, log log.Logger, backend backend.Backend) *PublicTxPoolAPI {
	api := &PublicTxPoolAPI{
		clientCtx: clientCtx,
		backend:   backend,
		logger:    log.With("module", "json-rpc", "namespace", "txpool"),
	}
	return api
}

// Content returns the transactions contained within the transaction pool.
func (s *PublicTxPoolAPI) Content() map[string]map[string]map[string]*rpctypes.Transaction {
	//todo get address
	address := "address"
	content := map[string]map[string]map[string]*rpctypes.Transaction{
		"queued":  make(map[string]map[string]*rpctypes.Transaction),
	}
	txs, err := s.backend.UserPendingTransactions(address, -1)
	if err != nil {
		s.logger.Error("txpool.Content err: ", err)
	}

	// Flatten the queued transactions
	dump := make(map[string]*rpctypes.Transaction)
	for _, tx := range txs {
		dump[fmt.Sprintf("%d", tx.Nonce)] = tx
	}
	content["queued"][address] = dump

	return content
}

// Status returns the number of pending and queued transaction in the pool.
func (s *PublicTxPoolAPI) Status() map[string]hexutil.Uint {
	numRes, err := s.backend.PendingTransactionCnt()
	if err != nil {
		s.logger.Error("txpool.Status err: ", err)
		return nil
	}
	return map[string]hexutil.Uint{
		"pending": 0,
		"queued":  hexutil.Uint(numRes),
	}
}

// Inspect retrieves the content of the transaction pool and flattens it into an
// easily inspectable list.
func (s *PublicTxPoolAPI) Inspect() map[string]map[string]map[string]string {
	//todo get address
	address := "address"
	content := map[string]map[string]map[string]string{
		"queued":  make(map[string]map[string]string),
	}
	txs, err := s.backend.UserPendingTransactions(address, -1)
	if err != nil {
		s.logger.Error("txpool.Content err: ", err)
	}

	// Define a formatter to flatten a transaction into a string
	var format = func(tx *rpctypes.Transaction) string {
		if to := tx.To; to != nil {
			return fmt.Sprintf("%s: %v wei + %v gas × %v wei", tx.To.Hex(), tx.Value, tx.Gas, tx.GasPrice)
		}
		return fmt.Sprintf("contract creation: %v wei + %v gas × %v wei", tx.Value, tx.Gas, tx.GasPrice)
	}

	// Flatten the queued transactions
	dump := make(map[string]string)
	for _, tx := range txs {
		dump[fmt.Sprintf("%d", tx.Nonce)] = format(tx)
	}
	content["queued"][address] = dump

	return content
}
