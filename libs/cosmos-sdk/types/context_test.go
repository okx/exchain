package types_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	dbm "github.com/okex/exchain/libs/tm-db"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"

	"github.com/okex/exchain/libs/cosmos-sdk/store"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
)

type MockLogger struct {
	logs *[]string
}

func NewMockLogger() MockLogger {
	logs := make([]string, 0)
	return MockLogger{
		&logs,
	}
}

func (l MockLogger) Debug(msg string, kvs ...interface{}) {
	*l.logs = append(*l.logs, msg)
}

func (l MockLogger) Info(msg string, kvs ...interface{}) {
	*l.logs = append(*l.logs, msg)
}

func (l MockLogger) Error(msg string, kvs ...interface{}) {
	*l.logs = append(*l.logs, msg)
}

func (l MockLogger) With(kvs ...interface{}) log.Logger {
	panic("not implemented")
}

func defaultContext(key types.StoreKey) types.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, types.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := types.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestCacheContext(t *testing.T) {
	key := types.NewKVStoreKey(t.Name())
	k1 := []byte("hello")
	v1 := []byte("world")
	k2 := []byte("key")
	v2 := []byte("value")

	ctx := defaultContext(key)
	store := ctx.KVStore(key)
	store.Set(k1, v1)
	require.Equal(t, v1, store.Get(k1))
	require.Nil(t, store.Get(k2))

	cctx, write := ctx.CacheContext()
	cstore := cctx.KVStore(key)
	require.Equal(t, v1, cstore.Get(k1))
	require.Nil(t, cstore.Get(k2))

	cstore.Set(k2, v2)
	require.Equal(t, v2, cstore.Get(k2))
	require.Nil(t, store.Get(k2))

	write()

	require.Equal(t, v2, store.Get(k2))
}

func TestLogContext(t *testing.T) {
	key := types.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	logger := NewMockLogger()
	ctx.SetLogger(logger)
	ctx.Logger().Debug("debug")
	ctx.Logger().Info("info")
	ctx.Logger().Error("error")
	require.Equal(t, *logger.logs, []string{"debug", "info", "error"})
}

type dummy int64 //nolint:unused

func (d dummy) Clone() interface{} {
	return d
}

// Testing saving/loading sdk type values to/from the context
func TestContextWithCustom(t *testing.T) {
	var ctx types.Context
	require.True(t, ctx.IsZero())

	header := abci.Header{}
	height := int64(1)
	chainid := "chainid"
	ischeck := true
	txbytes := []byte("txbytes")
	logger := NewMockLogger()
	voteinfos := []abci.VoteInfo{{}}
	meter := types.NewGasMeter(10000)
	minGasPrices := types.DecCoins{types.NewInt64DecCoin("feetoken", 1)}

	ctx = types.NewContext(nil, header, ischeck, logger)
	require.Equal(t, header, ctx.BlockHeader())

	ctx.
		SetBlockHeight(height).
		SetChainID(chainid).
		SetTxBytes(txbytes).
		SetVoteInfos(voteinfos).
		SetGasMeter(meter).
		SetMinGasPrices(minGasPrices)

	require.Equal(t, height, ctx.BlockHeight())
	require.Equal(t, chainid, ctx.ChainID())
	require.Equal(t, ischeck, ctx.IsCheckTx())
	require.Equal(t, txbytes, ctx.TxBytes())
	require.Equal(t, logger, ctx.Logger())
	require.Equal(t, voteinfos, ctx.VoteInfos())
	require.Equal(t, meter, ctx.GasMeter())
	require.Equal(t, minGasPrices, ctx.MinGasPrices())
}

// Testing saving/loading of header fields to/from the context
func TestContextHeader(t *testing.T) {
	var ctx types.Context

	height := int64(5)
	time := time.Now()
	addr := secp256k1.GenPrivKey().PubKey().Address()
	proposer := types.ConsAddress(addr)

	ctx = types.NewContext(nil, abci.Header{}, false, nil)

	ctx.
		SetBlockHeight(height).
		SetBlockTime(time).
		SetProposer(proposer)

	require.Equal(t, height, ctx.BlockHeight())
	require.Equal(t, height, ctx.BlockHeader().Height)
	require.Equal(t, time.UTC(), ctx.BlockHeader().Time)
	require.Equal(t, proposer.Bytes(), ctx.BlockHeader().ProposerAddress)
}

func TestContextHeaderClone(t *testing.T) {
	cases := map[string]struct {
		h abci.Header
	}{
		"empty": {
			h: abci.Header{},
		},
		"height": {
			h: abci.Header{
				Height: 77,
			},
		},
		"time": {
			h: abci.Header{
				Time: time.Unix(12345677, 12345),
			},
		},
		"zero time": {
			h: abci.Header{
				Time: time.Unix(0, 0),
			},
		},
		"many items": {
			h: abci.Header{
				Height:  823,
				Time:    time.Unix(9999999999, 0),
				ChainID: "silly-demo",
			},
		},
		"many items with hash": {
			h: abci.Header{
				Height:        823,
				Time:          time.Unix(9999999999, 0),
				ChainID:       "silly-demo",
				AppHash:       []byte{5, 34, 11, 3, 23},
				ConsensusHash: []byte{11, 3, 23, 87, 3, 1},
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctx := types.NewContext(nil, tc.h, false, nil)
			require.Equal(t, tc.h.Height, ctx.BlockHeight())
			require.Equal(t, tc.h.Time.UTC(), ctx.BlockTime())

			// update only changes one field
			var newHeight int64 = 17
			ctx.SetBlockHeight(newHeight)
			require.Equal(t, newHeight, ctx.BlockHeight())
			require.Equal(t, tc.h.Time.UTC(), ctx.BlockTime())
		})
	}
}

//go:noinline
func testFoo(ctx types.Context) int {
	return len(ctx.From())
}

func BenchmarkContextDuffCopy(b *testing.B) {
	ctx := types.NewContext(nil, abci.Header{}, false, nil)
	b.Run("1", func(b *testing.B) {
		b.Run("with", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ctx = ctx.WithIsCheckTx(true)
				testFoo(ctx)
			}
		})
		b.Run("set", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ctx.SetIsCheckTx(true)
				testFoo(ctx)
			}
		})
	})

	b.Run("2", func(b *testing.B) {
		b.Run("with", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newCtx := ctx.WithIsCheckTx(true)
				testFoo(newCtx)
			}
		})
		b.Run("set", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newCtx := ctx
				newCtx.SetIsCheckTx(true)
				testFoo(newCtx)
			}
		})
	})

	b.Run("3", func(b *testing.B) {
		b.Run("with", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testFoo(ctx.WithIsCheckTx(true))
			}
		})
		b.Run("set", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newCtx := ctx
				newCtx.SetIsCheckTx(true)
				testFoo(newCtx)
			}
		})
	})

	b.Run("4", func(b *testing.B) {
		b.Run("with", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testFoo(ctx.WithIsCheckTx(true).WithIsReCheckTx(false))
			}
		})
		b.Run("set", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newCtx := ctx
				newCtx.SetIsCheckTx(true).SetIsReCheckTx(false)
				testFoo(newCtx)
			}
		})
	})
}

func BenchmarkContextWrapAndUnwrap(b *testing.B) {
	key := types.NewKVStoreKey(b.Name())

	ctx := defaultContext(key)
	logger := NewMockLogger()
	ctx.SetLogger(logger)

	height := int64(1)
	chainid := "chainid"
	txbytes := []byte("txbytes")
	voteinfos := []abci.VoteInfo{{}}
	meter := types.NewGasMeter(10000)
	minGasPrices := types.DecCoins{types.NewInt64DecCoin("feetoken", 1)}
	ctx.
		SetBlockHeight(height).
		SetChainID(chainid).
		SetTxBytes(txbytes).
		SetVoteInfos(voteinfos).
		SetGasMeter(meter).
		SetMinGasPrices(minGasPrices).SetContext(context.Background())
	ctxC := types.WrapSDKContext(ctx)
	b.Run("UnwrapSDKContext", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			types.UnwrapSDKContext(ctxC)
		}
	})

	b.Run("WrapSDKContext", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			types.WrapSDKContext(ctx)
		}
	})
}
