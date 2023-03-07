package watcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"reflect"
	"sync"
	"testing"
)

func TestWatcher_getRealTx(t *testing.T) {
	type fields struct {
		store          *WatchStore
		height         uint64
		blockHash      common.Hash
		header         types.Header
		batch          []WatchMessage
		cumulativeGas  map[uint64]uint64
		gasUsed        uint64
		blockTxs       []common.Hash
		blockStdTxs    []common.Hash
		enable         bool
		firstUse       bool
		delayEraseKey  [][]byte
		eraseKeyFilter map[string][]byte
		log            log.Logger
		watchData      *WatchData
		jobChan        chan func()
		jobDone        *sync.WaitGroup
		evmTxIndex     uint64
		checkWd        bool
		filterMap      map[string]struct{}
		InfuraKeeper   InfuraKeeper
		delAccountMtx  sync.Mutex
	}
	type args struct {
		tx        types.TxEssentials
		txDecoder types.TxDecoder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Tx
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Watcher{
				store:          tt.fields.store,
				height:         tt.fields.height,
				blockHash:      tt.fields.blockHash,
				header:         tt.fields.header,
				batch:          tt.fields.batch,
				cumulativeGas:  tt.fields.cumulativeGas,
				gasUsed:        tt.fields.gasUsed,
				blockTxs:       tt.fields.blockTxs,
				blockStdTxs:    tt.fields.blockStdTxs,
				enable:         tt.fields.enable,
				firstUse:       tt.fields.firstUse,
				delayEraseKey:  tt.fields.delayEraseKey,
				eraseKeyFilter: tt.fields.eraseKeyFilter,
				log:            tt.fields.log,
				watchData:      tt.fields.watchData,
				jobChan:        tt.fields.jobChan,
				jobDone:        tt.fields.jobDone,
				evmTxIndex:     tt.fields.evmTxIndex,
				checkWd:        tt.fields.checkWd,
				filterMap:      tt.fields.filterMap,
				InfuraKeeper:   tt.fields.InfuraKeeper,
				delAccountMtx:  tt.fields.delAccountMtx,
			}
			got, err := w.getRealTx(tt.args.tx, tt.args.txDecoder)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRealTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRealTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
