package benchmarks

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/rlp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	types3 "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	types2 "github.com/okex/exchain/x/evm/types"
	token "github.com/okex/exchain/x/token/types"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"math/big"
	"testing"
	"time"
)

func BenchmarkTxSending(b *testing.B) {
	cases := map[string]struct {
		db          func(*testing.B) dbm.DB
		txBuilder   func(int, *AppInfo) []types3.Tx
		blockSize   int
		numAccounts int
	}{
		"basic send - memdb": {
			db:          buildMemDB,
			blockSize:   20,
			txBuilder:   buildTxFromMsg(bankSendMsg),
			numAccounts: 50,
		},
		"oip20 transfer - memdb": {
			db:          buildMemDB,
			blockSize:   20,
			txBuilder:   buildOip20Transfer,
			numAccounts: 50,
		},
		"basic send - leveldb": {
			db:          buildLevelDB,
			blockSize:   20,
			txBuilder:   buildTxFromMsg(bankSendMsg),
			numAccounts: 50,
		},
		"oip20 transfer - leveldb": {
			db:          buildLevelDB,
			blockSize:   20,
			txBuilder:   buildOip20Transfer,
			numAccounts: 50,
		},
		"basic send - leveldb - 8k accounts": {
			db:          buildLevelDB,
			blockSize:   20,
			txBuilder:   buildTxFromMsg(bankSendMsg),
			numAccounts: 8000,
		},
		"oip20 transfer - leveldb - 8k accounts": {
			db:          buildLevelDB,
			blockSize:   20,
			txBuilder:   buildOip20Transfer,
			numAccounts: 8000,
		},
		"basic send - leveldb - 8k accounts - huge blocks": {
			db:          buildLevelDB,
			blockSize:   1000,
			txBuilder:   buildTxFromMsg(bankSendMsg),
			numAccounts: 8000,
		},
		"oip20 transfer - leveldb - 8k accounts - huge blocks": {
			db:          buildLevelDB,
			blockSize:   1000,
			txBuilder:   buildOip20Transfer,
			numAccounts: 8000,
		},
		"basic send - leveldb - 80k accounts": {
			db:          buildLevelDB,
			blockSize:   20,
			txBuilder:   buildTxFromMsg(bankSendMsg),
			numAccounts: 80000,
		},
		"oip20 transfer - leveldb - 80k accounts": {
			db:          buildLevelDB,
			blockSize:   20,
			txBuilder:   buildOip20Transfer,
			numAccounts: 80000,
		},
	}

	for name, tc := range cases {
		b.Run(name, func(b *testing.B) {
			db := tc.db(b)
			defer db.Close()
			appInfo := InitializeOKXApp(b, db, tc.numAccounts)
			deployOip20(&appInfo)
			txs := tc.txBuilder(b.N, &appInfo)

			// number of Tx per block for the benchmarks
			blockSize := tc.blockSize
			height := int64(3)

			b.ResetTimer()

			for i := 0; i < b.N; {
				if i%blockSize == 0 {
					appInfo.App.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{ChainID: "exchain-67", Height: height, Time: time.Now()}})
				}
				//res := appInfo.App.CheckTx(abci.RequestCheckTx{
				//	Tx: txs[idx],
				//})
				//require.True(b, res.IsOK())

				res2 := appInfo.App.DeliverTx(abci.RequestDeliverTx{
					Tx: txs[i],
				})
				require.True(b, res2.IsOK())
				i++
				if i%blockSize == 0 {
					appInfo.App.EndBlock(abci.RequestEndBlock{Height: height})
					appInfo.App.Commit(abci.RequestCommit{})
					height++
				}
			}
		})
	}
}

func bankSendMsg(info *AppInfo) ([]sdk.Msg, error) {
	// Precompute all txs
	return tokenSendMsg(info)
}

func tokenSendMsg(info *AppInfo) ([]sdk.Msg, error) {
	// Precompute all txs
	rcpt := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	coins := sdk.Coins{sdk.NewInt64Coin(info.Denom, 1)}
	sendMsg := token.NewMsgTokenSend(info.MinterAddr, rcpt, coins)
	return []sdk.Msg{sendMsg}, nil
}

func buildTxFromMsg(builder func(info *AppInfo) ([]sdk.Msg, error)) func(n int, info *AppInfo) []types3.Tx {
	return func(n int, info *AppInfo) []types3.Tx {
		return GenSequenceOfTxs(info, builder, n)
	}
}

func buildOip20Transfer(n int, info *AppInfo) []types3.Tx {
	txs := make([]types3.Tx, n)
	// call oip20 transfer
	OipBytes, err := hex.DecodeString(Oip20TransferPayload)
	if err != nil {
		panic(err)
	}
	for i := range txs {
		oipTransferTx := types2.NewMsgEthereumTx(info.Nonce, &info.ContractAddr, nil, GasLimit, big.NewInt(GasPrice), OipBytes)
		if err := oipTransferTx.Sign(big.NewInt(ChainId), info.evmMintKey); err != nil {
			panic(err)
		}
		info.Nonce++
		tx, err := rlp.EncodeToBytes(oipTransferTx)
		if err != nil {
			panic(err)
		}
		txs[i] = tx
	}

	return txs
}

func buildMemDB(b *testing.B) dbm.DB {
	return dbm.NewMemDB()
}

func buildLevelDB(b *testing.B) dbm.DB {
	levelDB, err := dbm.NewGoLevelDBWithOpts("testing", b.TempDir(), &opt.Options{BlockCacher: opt.NoCacher})
	require.NoError(b, err)
	return levelDB
}
