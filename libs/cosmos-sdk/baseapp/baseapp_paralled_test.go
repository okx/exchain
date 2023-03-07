package baseapp_test

import (
	"encoding/json"
	"github.com/okx/okbchain/libs/system"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	types3 "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/simapp/helpers"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	simapp2 "github.com/okx/okbchain/libs/ibc-go/testing/simapp"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	types2 "github.com/okx/okbchain/libs/tendermint/types"
	"github.com/okx/okbchain/x/evm/types"
	types4 "github.com/okx/okbchain/x/token/types"
	"github.com/stretchr/testify/require"
)

type Env struct {
	priv []ethsecp256k1.PrivKey
	addr []sdk.AccAddress
}
type Chain struct {
	app          *simapp2.SimApp
	priv         []ethsecp256k1.PrivKey
	addr         []sdk.AccAddress
	acc          []*types3.EthAccount
	seq          []uint64
	num          []uint64
	chainIdStr   string
	chainIdInt   *big.Int
	ContractAddr []byte
}

func NewChain(env *Env) *Chain {

	chain := new(Chain)
	chain.acc = make([]*types3.EthAccount, 10)
	chain.priv = make([]ethsecp256k1.PrivKey, 10)
	chain.addr = make([]sdk.AccAddress, 10)
	chain.seq = make([]uint64, 10)
	chain.num = make([]uint64, 10)
	genAccs := []authexported.GenesisAccount{}
	for i := 0; i < 10; i++ {
		chain.acc[i] = &types3.EthAccount{
			BaseAccount: &auth.BaseAccount{
				Address: env.addr[i],
				Coins:   sdk.Coins{sdk.NewInt64Coin(system.Currency, 1000000)},
			},
			//CodeHash: []byte{1, 2},
		}
		genAccs = append(genAccs, chain.acc[i])
		chain.priv[i] = env.priv[i]
		chain.addr[i] = env.addr[i]
		chain.seq[i] = 0
		chain.num[i] = uint64(i)
	}
	chain.chainIdStr = "ethermint-3"
	chain.chainIdInt = big.NewInt(3)

	chain.app = simapp2.SetupWithGenesisAccounts(genAccs, sdk.NewDecCoins(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDecWithPrec(1000000, 0))))
	//header := abci.Header{Height: app.LastBlockHeight() + 1, ChainID: chainIdStr}

	chain.app.BaseApp.Commit(abci.RequestCommit{})
	return chain
}

func createEthTx(t *testing.T, chain *Chain, i int) []byte {
	amount, gasPrice, gasLimit := int64(1024), int64(2048), uint64(100000)
	addrTo := ethcmn.BytesToAddress(chain.priv[i+1].PubKey().Address().Bytes())
	msg := types.NewMsgEthereumTx(chain.seq[i], &addrTo, big.NewInt(amount), gasLimit, big.NewInt(gasPrice), []byte("test"))
	chain.seq[i]++
	err := msg.Sign(chain.chainIdInt, chain.priv[i].ToECDSA())
	require.NoError(t, err)
	rawtx, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)

	return rawtx
}

// contract Storage {
// uint256 number;
// /**
// * @dev Store value in variable
// * @param num value to store
// */
// function store(uint256 num) public {
// number = num;
// }
// function add() public {
// number += 1;
// }
// /**
// * @dev Return value
// * @return value of 'number'
// */
// function retrieve() public view returns (uint256){
// return number;
// }
// }
var abiStr = `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"add","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func deployContract(t *testing.T, chain *Chain, i int) []byte {
	// Deploy contract - Owner.sol
	gasLimit := uint64(10000000000000)
	gasPrice := big.NewInt(10000)

	sender := ethcmn.HexToAddress(chain.priv[i].PubKey().Address().String())

	bytecode := ethcmn.FromHex("608060405234801561001057600080fd5b50610217806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631003e2d2146100465780632e64cec1146100625780636057361d14610080575b600080fd5b610060600480360381019061005b9190610105565b61009c565b005b61006a6100b7565b6040516100779190610141565b60405180910390f35b61009a60048036038101906100959190610105565b6100c0565b005b806000808282546100ad919061018b565b9250508190555050565b60008054905090565b8060008190555050565b600080fd5b6000819050919050565b6100e2816100cf565b81146100ed57600080fd5b50565b6000813590506100ff816100d9565b92915050565b60006020828403121561011b5761011a6100ca565b5b6000610129848285016100f0565b91505092915050565b61013b816100cf565b82525050565b60006020820190506101566000830184610132565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610196826100cf565b91506101a1836100cf565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156101d6576101d561015c565b5b82820190509291505056fea2646970667358221220318e29d6b4806f219eedd0cc861e82c13e28eb7f42161f2c780dc539b0e32b4e64736f6c634300080a0033")
	msg := types.NewMsgEthereumTx(chain.seq[i], &sender, big.NewInt(0), gasLimit, gasPrice, bytecode)
	err := msg.Sign(big.NewInt(3), chain.priv[i].ToECDSA())
	require.NoError(t, err)
	chain.seq[i]++
	rawTx, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	return rawTx
}

var contractJson = `{"abi":[{"inputs":[],"name":"add","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}],"bin":"608060405234801561001057600080fd5b50610205806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80632e64cec1146100465780634f2be91f146100645780636057361d1461006e575b600080fd5b61004e61008a565b60405161005b91906100d1565b60405180910390f35b61006c610093565b005b6100886004803603810190610083919061011d565b6100ae565b005b60008054905090565b60016000808282546100a59190610179565b92505081905550565b8060008190555050565b6000819050919050565b6100cb816100b8565b82525050565b60006020820190506100e660008301846100c2565b92915050565b600080fd5b6100fa816100b8565b811461010557600080fd5b50565b600081359050610117816100f1565b92915050565b600060208284031215610133576101326100ec565b5b600061014184828501610108565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610184826100b8565b915061018f836100b8565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156101c4576101c361014a565b5b82820190509291505056fea2646970667358221220742b7232e733bee3592cb9e558bdae3fbd0006bcbdba76abc47b6020744037b364736f6c634300080a0033"}`

type CompiledContract struct {
	ABI abi.ABI
	Bin string
}

func UnmarshalContract(t *testing.T) *CompiledContract {
	cc := new(CompiledContract)
	err := json.Unmarshal([]byte(contractJson), cc)
	require.NoError(t, err)
	return cc
}

func callContract(t *testing.T, chain *Chain, i int) []byte {
	gasLimit := uint64(10000000000000)
	gasPrice := big.NewInt(10000)
	//to := ethcmn.HexToAddress(chain.priv[i].PubKey().Address().String())
	to := ethcmn.BytesToAddress(chain.ContractAddr)
	cc := UnmarshalContract(t)
	data, err := cc.ABI.Pack("add")
	require.NoError(t, err)
	msg := types.NewMsgEthereumTx(chain.seq[i], &to, big.NewInt(0), gasLimit, gasPrice, data)
	err = msg.Sign(big.NewInt(3), chain.priv[i].ToECDSA())
	require.NoError(t, err)
	chain.seq[i]++
	rawTx, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	return rawTx
}

func createCosmosTx(t *testing.T, chain *Chain, i int) []byte {
	msg := types4.NewMsgTokenSend(chain.addr[i], chain.addr[i+1], sdk.Coins{sdk.NewInt64Coin(system.Currency, 10)})

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 1)},
		helpers.DefaultGenTxGas,
		chain.chainIdStr,
		[]uint64{chain.num[i]},
		[]uint64{chain.seq[i]},
		chain.priv[i],
	)
	chain.seq[i]++

	txBytes, err := chain.app.Codec().MarshalBinaryLengthPrefixed(tx)
	require.Nil(t, err)
	return txBytes
}

func runtxs(chain *Chain, rawTxs [][]byte, isParalle bool) []*abci.ResponseDeliverTx {
	header := abci.Header{Height: chain.app.LastBlockHeight() + 1, ChainID: chain.chainIdStr}
	chain.app.BaseApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	ret := []*abci.ResponseDeliverTx{}
	if isParalle {
		ret = chain.app.BaseApp.ParallelTxs(rawTxs, false)
	} else {
		for _, tx := range rawTxs {
			r := chain.app.BaseApp.DeliverTx(abci.RequestDeliverTx{Tx: tx})
			ret = append(ret, &r)
		}
	}
	chain.app.BaseApp.EndBlock(abci.RequestEndBlock{})
	chain.app.BaseApp.Commit(abci.RequestCommit{})

	return ret
}

func DeployContractAndGetContractAddress(t *testing.T, chain *Chain) {
	rawTxs := [][]byte{}
	rawTxs = append(rawTxs, deployContract(t, chain, 0))
	r := runtxs(chain, rawTxs, false)

	for _, e := range r[0].Events {
		for _, v := range e.Attributes {
			if string(v.Key) == "recipient" {
				chain.ContractAddr = v.Value
			}
		}
	}
}

func TestParalledTxs(t *testing.T) {
	env := new(Env)
	accountNum := 10
	env.priv = make([]ethsecp256k1.PrivKey, 10)
	env.addr = make([]sdk.AccAddress, 10)
	for i := 0; i < accountNum; i++ {
		priv, _ := ethsecp256k1.GenerateKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		env.priv[i] = priv
		env.addr[i] = addr
	}

	chainA, chainB := NewChain(env), NewChain(env)
	//deploy contract on chainA and chainB
	DeployContractAndGetContractAddress(t, chainA)
	DeployContractAndGetContractAddress(t, chainB)

	testCases := []struct {
		name     string
		malleate func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte)
	}{
		{
			"five txs one group:a->b b->c c->d d->e e->f",
			func(t *testing.T, chain *Chain, isParalled bool) ([]byte, []byte) {

				rawTxs := [][]byte{}
				for i := 0; i < 5; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}

				header := abci.Header{Height: chain.app.LastBlockHeight() + 1, ChainID: chain.chainIdStr}
				chain.app.BaseApp.BeginBlock(abci.RequestBeginBlock{Header: header})
				ret := runtxs(chain, rawTxs, isParalled)
				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs two group, no conflict:a->b b->c / d->e e->f f->g",
			func(t *testing.T, chain *Chain, isParalled bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				//one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runtxs(chain, rawTxs, isParalled)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs two group, has conflict",
			func(t *testing.T, chain *Chain, isParalled bool) ([]byte, []byte) {
				rawTxs := [][]byte{}

				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, callContract(t, chain, i))
				}
				////one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runtxs(chain, rawTxs, isParalled)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs one group with cosmos tx",
			func(t *testing.T, chain *Chain, isParalled bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 2; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				//cosmostx
				rawTxs = append(rawTxs, createCosmosTx(t, chain, 2))
				//one group 2txs
				for i := 3; i < 5; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runtxs(chain, rawTxs, isParalled)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs two group, no conflict with cosmos tx",
			func(t *testing.T, chain *Chain, isParalle bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs(2eth and 1 cosmos)
				for i := 0; i < 2; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				//cosmos tx
				rawTxs = append(rawTxs, createCosmosTx(t, chain, 2))
				//one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runtxs(chain, rawTxs, isParalle)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs two group, has conflict with cosmos tx",
			func(t *testing.T, chain *Chain, isParalled bool) ([]byte, []byte) {
				rawTxs := [][]byte{}

				//one group 3txs:2 evm tx with conflict, one cosmos tx
				for i := 0; i < 2; i++ {
					rawTxs = append(rawTxs, callContract(t, chain, i))
				}
				rawTxs = append(rawTxs, createCosmosTx(t, chain, 2))
				////one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createCosmosTx(t, chain, i))
				}
				ret := runtxs(chain, rawTxs, isParalled)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
	}
	for _, tc := range testCases {
		resultHashA, appHashA := tc.malleate(t, chainA, true)
		resultHashB, appHashB := tc.malleate(t, chainB, false)
		require.True(t, reflect.DeepEqual(resultHashA, resultHashB))
		require.True(t, reflect.DeepEqual(appHashA, appHashB))
	}

}

func resultHash(txs []*abci.ResponseDeliverTx) []byte {
	results := types2.NewResults(txs)
	return results.Hash()
}
