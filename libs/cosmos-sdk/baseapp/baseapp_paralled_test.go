package baseapp_test

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	types3 "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/simapp/helpers"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	simapp2 "github.com/okex/exchain/libs/ibc-go/testing/simapp"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types"
	types4 "github.com/okex/exchain/x/token/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

type Env struct {
	priv []ethsecp256k1.PrivKey
	addr []sdk.AccAddress
}
type Chain struct {
	app        *simapp2.SimApp
	priv       []ethsecp256k1.PrivKey
	addr       []sdk.AccAddress
	acc        []*types3.EthAccount
	seq        []uint64
	num        []uint64
	chainIdStr string
	chainIdInt *big.Int
}

func NewChain(env *Env) *Chain {
	types2.UnittestOnlySetMilestoneVenusHeight(-1)
	types2.UnittestOnlySetMilestoneVenus1Height(-1)
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

func createCosmosTx(t *testing.T, chain *Chain, i int) []byte {
	msg := types4.NewMsgTokenSend(chain.addr[i], chain.addr[i+1], sdk.Coins{sdk.NewInt64Coin("okt", 10)})

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
		//{
		//	"five txs two group, has conflict",
		//	func(t *testing.T, chain *Chain, isParalle bool) ([]byte, []byte) {},
		//	true,
		//},
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
		//{
		//	"five txs two group, has conflict with cosmos tx",
		//	func(t *testing.T, chain *Chain, isParalle bool) ([]byte, []byte) {}
		//	true,
		//},
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
