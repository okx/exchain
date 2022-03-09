package backend

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLruCache_AddOrUpdateBlock(t *testing.T) {
	type args struct {
		block *watcher.Block
	}
	type result struct {
		blockCount int
		block      *watcher.Block
		txCount    int
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			name: "cache empty Block",
			args: args{
				block: &watcher.Block{
					Number:       hexutil.Uint64(0x10),
					Hash:         common.HexToHash("0x6b2cfa0a20e291ca0bb58b2112086f247026bb94a65133e87ee3aaa4658399e5"),
					Transactions: []*watcher.Transaction{},
				},
			},
			result: result{
				blockCount: 1,
				block: &watcher.Block{
					Number:       hexutil.Uint64(0x10),
					Hash:         common.HexToHash("0x6b2cfa0a20e291ca0bb58b2112086f247026bb94a65133e87ee3aaa4658399e5"),
					Transactions: []*watcher.Transaction{},
				},
				txCount: 0,
			},
		},
		{
			name: "duplicate Block",
			args: args{
				block: &watcher.Block{
					Number:       hexutil.Uint64(0x10),
					Hash:         common.HexToHash("0x6b2cfa0a20e291ca0bb58b2112086f247026bb94a65133e87ee3aaa4658399e5"),
					Transactions: []*watcher.Transaction{},
				},
			},
			result: result{
				blockCount: 1,
				block: &watcher.Block{
					Number:       hexutil.Uint64(0x10),
					Hash:         common.HexToHash("0x6b2cfa0a20e291ca0bb58b2112086f247026bb94a65133e87ee3aaa4658399e5"),
					Transactions: []*watcher.Transaction{},
				},
				txCount: 0,
			},
		},
		{
			name: "Block with txs",
			args: args{
				block: &watcher.Block{
					Number: hexutil.Uint64(0x11),
					Hash:   common.HexToHash("0x3bb254ed105476b94583eec8375c5d2fc0a5cf50047c5912b4337ba43a837b88"),
					Transactions: []*watcher.Transaction{
						{
							From: common.HexToAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f"),
							Hash: common.HexToHash("0xb4a40e844ee4c012d4a6d9e16d4ee8dcf52ef5042da491dbc73574f6764e17d1"),
						},
					},
				},
			},
			result: result{
				blockCount: 2,
				txCount:    1,
				block: &watcher.Block{
					Number: hexutil.Uint64(0x11),
					Hash:   common.HexToHash("0x3bb254ed105476b94583eec8375c5d2fc0a5cf50047c5912b4337ba43a837b88"),
					Transactions: []*watcher.Transaction{
						{
							From: common.HexToAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f"),
							Hash: common.HexToHash("0xb4a40e844ee4c012d4a6d9e16d4ee8dcf52ef5042da491dbc73574f6764e17d1"),
						},
					},
				},
			},
		},
	}
	viper.Set(FlagApiBackendLru, 100) // must be 3
	alc := NewLruCache()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alc.AddOrUpdateBlock(tt.args.block.Hash, tt.args.block)
			blockLru := alc.lruBlock
			require.NotNil(t, blockLru)
			require.Equal(t, tt.result.blockCount, blockLru.Len())

			block, err := alc.GetBlockByHash(tt.result.block.Hash)
			require.Nil(t, err)
			require.NotNil(t, block)
			require.Equal(t, tt.result.block.Hash, block.Hash)

			//must update tx in block
			txLru := alc.lruTx
			require.NotNil(t, txLru)
			require.Equal(t, tt.result.txCount, txLru.Len())
		})
	}
}

func TestLruCache_AddOrUpdateTransaction(t *testing.T) {
	type result struct {
		tx      *watcher.Transaction
		txCount int
	}
	type args struct {
		tx *watcher.Transaction
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			name: "cache tx",
			args: args{
				tx: &watcher.Transaction{
					From: common.HexToAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f"),
					Hash: common.HexToHash("0xb4a40e844ee4c012d4a6d9e16d4ee8dcf52ef5042da491dbc73574f6764e17d1"),
				},
			},
			result: result{
				txCount: 1,
				tx: &watcher.Transaction{
					From: common.HexToAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f"),
					Hash: common.HexToHash("0xb4a40e844ee4c012d4a6d9e16d4ee8dcf52ef5042da491dbc73574f6764e17d1"),
				},
			},
		},
		{
			name: "duplicate tx",
			args: args{
				tx: &watcher.Transaction{
					From: common.HexToAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f"),
					Hash: common.HexToHash("0xb4a40e844ee4c012d4a6d9e16d4ee8dcf52ef5042da491dbc73574f6764e17d1"),
				},
			},
			result: result{
				txCount: 1,
				tx: &watcher.Transaction{
					From: common.HexToAddress("0xbbe4733d85bc2b90682147779da49cab38c0aa1f"),
					Hash: common.HexToHash("0xb4a40e844ee4c012d4a6d9e16d4ee8dcf52ef5042da491dbc73574f6764e17d1"),
				},
			},
		},
	}
	viper.Set(FlagApiBackendLru, 100)
	alc := NewLruCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alc.AddOrUpdateTransaction(tt.args.tx.Hash, tt.args.tx)
			txLru := alc.lruTx
			require.NotNil(t, txLru)
			require.Equal(t, tt.result.txCount, txLru.Len())

			tx, err := alc.GetTransaction(tt.result.tx.Hash)
			require.Nil(t, err)
			require.NotNil(t, tx)
			require.Equal(t, tt.result.tx.Hash, tx.Hash)
		})
	}
}

func TestLruCache_GetBlockByNumber(t *testing.T) {
	type args struct {
		block *watcher.Block
	}
	type result struct {
		blockCount int
		block      *watcher.Block
		txCount    int
	}
	tests := []struct {
		name   string
		args   args
		result result
	}{
		{
			name: "Get Block by Number",
			args: args{
				block: &watcher.Block{
					Number:       hexutil.Uint64(0x10),
					Hash:         common.HexToHash("0x6b2cfa0a20e291ca0bb58b2112086f247026bb94a65133e87ee3aaa4658399e5"),
					Transactions: []*watcher.Transaction{},
				},
			},
			result: result{
				blockCount: 1,
				block: &watcher.Block{
					Number:       hexutil.Uint64(0x10),
					Hash:         common.HexToHash("0x6b2cfa0a20e291ca0bb58b2112086f247026bb94a65133e87ee3aaa4658399e5"),
					Transactions: []*watcher.Transaction{},
				},
				txCount: 0,
			},
		},
	}
	viper.Set(FlagApiBackendLru, 100) // must be 3
	alc := NewLruCache()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alc.AddOrUpdateBlock(tt.args.block.Hash, tt.args.block)
			alc.AddOrUpdateBlockHash(uint64(tt.args.block.Number), tt.args.block.Hash)

			blockLru := alc.lruBlock
			require.NotNil(t, blockLru)
			require.Equal(t, tt.result.blockCount, blockLru.Len())

			blockInfoLru := alc.lruBlockInfo
			require.NotNil(t, blockInfoLru)
			require.Equal(t, tt.result.blockCount, blockInfoLru.Len())

			block, err := alc.GetBlockByNumber(uint64(tt.result.block.Number))
			require.Nil(t, err)
			require.NotNil(t, block)
			require.Equal(t, tt.result.block.Hash, block.Hash)
		})
	}
}
