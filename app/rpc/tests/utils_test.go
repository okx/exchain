package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/okexchain/app/crypto/hd"
	"github.com/okex/okexchain/app/rpc/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	// keys that provided on node (from test.sh)
	mnemo1          = "plunge silk glide glass curve cycle snack garbage obscure express decade dirt"
	mnemo2          = "lazy cupboard wealth canoe pumpkin gasp play dash antenna monitor material village"
	defaultPassWd   = "12345678"
	defaultCoinType = 60
)

var (
	keyInfo1, keyInfo2 keys.Info
	Kb                 = keys.NewInMemory(hd.EthSecp256k1Options()...)
	hexAddr1, hexAddr2 ethcmn.Address
	addrCounter        = 2
	defaultGasPrice    sdk.SysCoin
)

func init() {
	config := sdk.GetConfig()
	config.SetCoinType(defaultCoinType)

	keyInfo1, _ = createAccountWithMnemo(mnemo1, "alice", defaultPassWd)
	keyInfo2, _ = createAccountWithMnemo(mnemo2, "bob", defaultPassWd)
	hexAddr1 = ethcmn.BytesToAddress(keyInfo1.GetAddress().Bytes())
	hexAddr2 = ethcmn.BytesToAddress(keyInfo2.GetAddress().Bytes())
	defaultGasPrice, _ = sdk.ParseDecCoin(defaultMinGasPrice)
}

func TestGetAddress(t *testing.T) {
	addr, err := GetAddress()
	require.NoError(t, err)
	require.True(t, bytes.Equal(addr, hexAddr1[:]))
}

func createAccountWithMnemo(mnemonic, name, passWd string) (info keys.Info, err error) {
	hdPath := keys.CreateHDPath(0, 0).String()
	info, err = Kb.CreateAccount(name, mnemonic, "", passWd, hdPath, hd.EthSecp256k1)
	if err != nil {
		return info, fmt.Errorf("failed. Kb.CreateAccount err : %s", err.Error())
	}

	return info, err
}

// sendTestTransaction sends a dummy transaction
func sendTestTransaction(t *testing.T, senderAddr, receiverAddr ethcmn.Address, value uint) ethcmn.Hash {
	fromAddrStr, toAddrStr := senderAddr.Hex(), receiverAddr.Hex()
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = fromAddrStr
	param[0]["to"] = toAddrStr
	param[0]["value"] = hexutil.Uint(value).String()
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	t.Logf("%s transfers %d to %s successfully\n", fromAddrStr, value, toAddrStr)
	return hash
}

func getBlockHeightFromTxHash(t *testing.T, hash ethcmn.Hash) hexutil.Uint64 {
	rpcRes := Call(t, "eth_getTransactionByHash", []interface{}{hash})
	var transaction types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transaction))

	if transaction.BlockNumber == nil {
		return hexutil.Uint64(0)
	}

	return hexutil.Uint64(transaction.BlockNumber.ToInt().Uint64())
}

func getBlockHashFromTxHash(t *testing.T, hash ethcmn.Hash) *ethcmn.Hash {
	rpcRes := Call(t, "eth_getTransactionByHash", []interface{}{hash})
	var transaction types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transaction))

	return transaction.BlockHash
}

func assertNullFromJSONResponse(t *testing.T, jrm json.RawMessage) {
	require.True(t, bytes.Equal([]byte("null"), jrm))
}
