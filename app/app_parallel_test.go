package app

import (
	"encoding/json"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	apptypes "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/simapp/helpers"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	tokentypes "github.com/okex/exchain/x/token/types"
)

type Env struct {
	priv []ethsecp256k1.PrivKey
	addr []sdk.AccAddress
}

type Chain struct {
	app          *OKExChainApp
	codec        *codec.Codec
	priv         []ethsecp256k1.PrivKey
	addr         []sdk.AccAddress
	acc          []apptypes.EthAccount
	seq          []uint64
	num          []uint64
	chainIdStr   string
	chainIdInt   *big.Int
	ContractAddr []byte
}

func NewChain(env *Env) *Chain {
	tmtypes.UnittestOnlySetMilestoneVenusHeight(-1)
	tmtypes.UnittestOnlySetMilestoneVenus1Height(1)
	tmtypes.UnittestOnlySetMilestoneVenus2Height(1)
	tmtypes.UnittestOnlySetMilestoneEarthHeight(1)
	tmtypes.UnittestOnlySetMilestoneVenus6Height(1)
	chain := new(Chain)
	chain.acc = make([]apptypes.EthAccount, 10)
	chain.priv = make([]ethsecp256k1.PrivKey, 10)
	chain.addr = make([]sdk.AccAddress, 10)
	chain.seq = make([]uint64, 10)
	chain.num = make([]uint64, 10)
	chain.chainIdStr = "ethermint-3"
	chain.chainIdInt = big.NewInt(3)
	// initialize account
	genAccs := make([]authexported.GenesisAccount, 0)
	for i := 0; i < 10; i++ {
		chain.acc[i] = apptypes.EthAccount{
			BaseAccount: &auth.BaseAccount{
				Address: env.addr[i],
				Coins:   sdk.Coins{sdk.NewInt64Coin("okt", 1000000)},
			},
			//CodeHash: []byte{1, 2},
		}
		genAccs = append(genAccs, chain.acc[i])
		chain.priv[i] = env.priv[i]
		chain.addr[i] = env.addr[i]
		chain.seq[i] = 0
		chain.num[i] = uint64(i)
	}

	chain.app = SetupWithGenesisAccounts(false, genAccs, WithChainId(chain.chainIdStr))
	chain.codec = chain.app.Codec()

	params := evmtypes.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	chain.app.EvmKeeper.SetParams(chain.Ctx(), params)

	chain.app.BaseApp.Commit(abci.RequestCommit{})
	return chain
}

func (chain *Chain) Ctx() sdk.Context {
	return chain.app.BaseApp.GetDeliverStateCtx()
}

func DeployContractAndGetContractAddress(t *testing.T, chain *Chain) {
	var rawTxs [][]byte
	rawTxs = append(rawTxs, deployContract(t, chain, 0))
	r := runTxs(chain, rawTxs, false)

	log := r[0].Log[1 : len(r[0].Log)-1]
	logMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(log), &logMap)
	require.NoError(t, err)

	logs := strings.Split(logMap["log"].(string), ";")
	require.True(t, len(logs) == 3)
	contractLog := strings.Split(logs[2], " ")
	require.True(t, len(contractLog) == 4)
	chain.ContractAddr = []byte(contractLog[3])
}

func createEthTx(t *testing.T, chain *Chain, addressIdx int) []byte {
	amount, gasPrice, gasLimit := int64(1024), int64(2048), uint64(100000)
	addrTo := ethcmn.BytesToAddress(chain.priv[addressIdx+1].PubKey().Address().Bytes())
	msg := evmtypes.NewMsgEthereumTx(chain.seq[addressIdx], &addrTo, big.NewInt(amount), gasLimit, big.NewInt(gasPrice), []byte{})
	chain.seq[addressIdx]++
	err := msg.Sign(chain.chainIdInt, chain.priv[addressIdx].ToECDSA())
	require.NoError(t, err)
	rawTx, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)

	return rawTx
}

func createFailedEthTx(t *testing.T, chain *Chain, addressIdx int) []byte {
	amount, gasPrice, gasLimit := int64(1024), int64(2048), uint64(1)
	addrTo := ethcmn.BytesToAddress(chain.priv[addressIdx+1].PubKey().Address().Bytes())
	msg := evmtypes.NewMsgEthereumTx(chain.seq[addressIdx], &addrTo, big.NewInt(amount), gasLimit, big.NewInt(gasPrice), []byte{})
	chain.seq[addressIdx]++
	err := msg.Sign(chain.chainIdInt, chain.priv[addressIdx].ToECDSA())
	require.NoError(t, err)
	rawTx, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)

	return rawTx
}

func createTokenSendTx(t *testing.T, chain *Chain, i int) []byte {
	msg := tokentypes.NewMsgTokenSend(chain.addr[i], chain.addr[i+1], sdk.Coins{sdk.NewInt64Coin("okt", 10)})

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

func runTxs(chain *Chain, rawTxs [][]byte, isParallel bool) []*abci.ResponseDeliverTx {
	header := abci.Header{Height: chain.app.LastBlockHeight() + 1, ChainID: chain.chainIdStr}
	chain.app.BaseApp.BeginBlock(abci.RequestBeginBlock{Header: header})
	var ret []*abci.ResponseDeliverTx
	if isParallel {
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

func TestParallelTxs(t *testing.T) {

	env := new(Env)
	env.priv = make([]ethsecp256k1.PrivKey, 10)
	env.addr = make([]sdk.AccAddress, 10)
	for i := 0; i < 10; i++ {
		priv, _ := ethsecp256k1.GenerateKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		env.priv[i] = priv
		env.addr[i] = addr
	}
	chainA, chainB := NewChain(env), NewChain(env)
	DeployContractAndGetContractAddress(t, chainA)
	DeployContractAndGetContractAddress(t, chainB)

	testCases := []struct {
		title      string
		executeTxs func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte)
	}{
		{
			"five evm txs, one group:a->b b->c c->d d->e e->f",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {

				var rawTxs [][]byte
				for i := 0; i < 5; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)
				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five failed evm txs, one group:a->b b->c c->d d->e e->f",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {

				var rawTxs [][]byte
				for i := 0; i < 5; i++ {
					rawTxs = append(rawTxs, createFailedEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)
				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five evm txs, two group:a->b b->c / d->e e->f f->g",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				//one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five failed evm txs, two group:a->b b->c / d->e e->f f->g",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, createFailedEthTx(t, chain, i))
				}
				//one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createFailedEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"three evm txs and two failed evm txs, two group:a->b b->c / d->e e->f f->g",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, createFailedEthTx(t, chain, i))
				}
				//one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"three evm txs and two failed evm txs, two group:a->b b->c / d->e e->f f->g",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, createFailedEthTx(t, chain, i))
				}
				//one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"three contract txs and two normal evm txs, two group",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				var rawTxs [][]byte

				//one group 3txs
				for i := 0; i < 3; i++ {
					rawTxs = append(rawTxs, callContract(t, chain, i))
				}
				////one group 2txs
				for i := 8; i > 6; i-- {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs one group with cosmos tx",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				rawTxs := [][]byte{}
				//one group 3txs
				for i := 0; i < 2; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				//cosmostx
				rawTxs = append(rawTxs, createTokenSendTx(t, chain, 7))
				//one group 2txs
				for i := 3; i < 5; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
		{
			"five txs two group, has conflict with cosmos tx",
			func(t *testing.T, chain *Chain, isParallel bool) ([]byte, []byte) {
				rawTxs := [][]byte{}

				//one group 2txs
				for i := 0; i < 2; i++ {
					rawTxs = append(rawTxs, callContract(t, chain, i))
				}
				////one group 3txs
				rawTxs = append(rawTxs, createTokenSendTx(t, chain, 8))
				for i := 8; i < 6; i++ {
					rawTxs = append(rawTxs, createEthTx(t, chain, i))
				}
				ret := runTxs(chain, rawTxs, isParallel)

				return resultHash(ret), chain.app.BaseApp.LastCommitID().Hash
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			resultHashA, appHashA := tc.executeTxs(t, chainA, true)
			resultHashB, appHashB := tc.executeTxs(t, chainB, false)
			require.True(t, reflect.DeepEqual(resultHashA, resultHashB))
			require.True(t, reflect.DeepEqual(appHashA, appHashB))
		})
	}
}

func resultHash(txs []*abci.ResponseDeliverTx) []byte {
	results := tmtypes.NewResults(txs)
	return results.Hash()
}

//contract Storage {
//uint256 number;
///**
// * @dev Store value in variable
// * @param num value to store
// */
//function store(uint256 num) public {
//number = num;
//}
//function add() public {
//number += 1;
//}
///**
// * @dev Return value
// * @return value of 'number'
// */
//function retrieve() public view returns (uint256){
//return number;
//}
//}
var abiStr = `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"add","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func deployContract(t *testing.T, chain *Chain, i int) []byte {
	// Deploy contract - Owner.sol
	gasLimit := uint64(30000000)
	gasPrice := big.NewInt(100000000)

	//sender := ethcmn.HexToAddress(chain.priv[i].PubKey().Address().String())

	bytecode := ethcmn.FromHex("608060405234801561001057600080fd5b50610217806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80631003e2d2146100465780632e64cec1146100625780636057361d14610080575b600080fd5b610060600480360381019061005b9190610105565b61009c565b005b61006a6100b7565b6040516100779190610141565b60405180910390f35b61009a60048036038101906100959190610105565b6100c0565b005b806000808282546100ad919061018b565b9250508190555050565b60008054905090565b8060008190555050565b600080fd5b6000819050919050565b6100e2816100cf565b81146100ed57600080fd5b50565b6000813590506100ff816100d9565b92915050565b60006020828403121561011b5761011a6100ca565b5b6000610129848285016100f0565b91505092915050565b61013b816100cf565b82525050565b60006020820190506101566000830184610132565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610196826100cf565b91506101a1836100cf565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156101d6576101d561015c565b5b82820190509291505056fea2646970667358221220318e29d6b4806f219eedd0cc861e82c13e28eb7f42161f2c780dc539b0e32b4e64736f6c634300080a0033")
	msg := evmtypes.NewMsgEthereumTx(chain.seq[i], nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
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
	gasLimit := uint64(30000000)
	gasPrice := big.NewInt(100000000)
	//to := ethcmn.HexToAddress(chain.priv[i].PubKey().Address().String())
	to := ethcmn.BytesToAddress(chain.ContractAddr)
	cc := UnmarshalContract(t)
	data, err := cc.ABI.Pack("add")
	require.NoError(t, err)
	msg := evmtypes.NewMsgEthereumTx(chain.seq[i], &to, big.NewInt(0), gasLimit, gasPrice, data)
	err = msg.Sign(big.NewInt(3), chain.priv[i].ToECDSA())
	require.NoError(t, err)
	chain.seq[i]++
	rawTx, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err)
	return rawTx
}
