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

const (
	testContractKind = iota
	erc20ContractKind
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

// deployTestContract deploys a contract that emits an event in the constructor
func DeployTestContract(t *testing.T, senderAddr ethcmn.Address, kind int) (ethcmn.Hash, map[string]interface{}) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddr.Hex()
	param[0]["gasprice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	switch kind {
	case testContractKind:
		param[0]["data"] = "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029"
	case erc20ContractKind:
		// TODO
		param[0]["data"] = ""
	default:
		panic("unsupported contract kind")
	}

	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))

	receipt := WaitForReceipt(t, hash)
	require.NotNil(t, receipt, "transaction failed")
	require.Equal(t, "0x1", receipt["status"].(string))

	return hash, receipt
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
