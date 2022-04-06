package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/rlp"
)

// TokenInitInfo ...
func (api *PublicEthereumAPI) TokenInitInfo(contractAddr common.Address) vm.TokenInitInfo {
	tokenInfoByte := vm.ReadToken([]byte(contractAddr.Hex()))
	var tokenInfo vm.TokenInitInfo
	if err := rlp.DecodeBytes(tokenInfoByte, &tokenInfo); err != nil {
		return tokenInfo
	}
	return tokenInfo
}

// GetInternalTransactions ...
func (api *PublicEthereumAPI) GetInternalTransactions(txHash string) []vm.InnerTx {
	return vm.GetFromDB(txHash)
}

// GetBlockInternalTransactions ...
func (api *PublicEthereumAPI) GetBlockInternalTransactions(blockHash string) (map[string][]vm.InnerTx, error) {
	var rtn = make(map[string][]vm.InnerTx)
	txHashes := vm.GetBlockDB(blockHash)
	if len(txHashes) > 0 {
		for _, txHash := range txHashes {
			inners := vm.GetFromDB(txHash)
			rtn[txHash] = inners
		}
		return rtn, nil
	} else {
		return nil, fmt.Errorf("no transaction found with hash %s, maybe this node has many blocks behind the tip of the chain", blockHash)
	}
}
