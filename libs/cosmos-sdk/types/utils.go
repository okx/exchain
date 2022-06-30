package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	tmtypes "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/ethereum/go-ethereum/common"

	dbm "github.com/okex/exchain/libs/tm-db"
)

var (
	// This is set at compile time. Could be cleveldb, defaults is goleveldb.
	DBBackend = ""
	backend   = dbm.GoLevelDBBackend
)

func init() {
	if len(DBBackend) != 0 {
		backend = dbm.BackendType(DBBackend)
	}
}

// SortedJSON takes any JSON and returns it sorted by keys. Also, all white-spaces
// are removed.
// This method can be used to canonicalize JSON to be returned by GetSignBytes,
// e.g. for the ledger integration.
// If the passed JSON isn't valid it will return an error.
func SortJSON(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// MustSortJSON is like SortJSON but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSortJSON(toSortJSON []byte) []byte {
	js, err := SortJSON(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}

// Uint64ToBigEndian - marshals uint64 to a bigendian byte slice so it can be sorted
func Uint64ToBigEndian(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// BigEndianToUint64 returns an uint64 from big endian encoded bytes. If encoding
// is empty, zero is returned.
func BigEndianToUint64(bz []byte) uint64 {
	if len(bz) == 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// Slight modification of the RFC3339Nano but it right pads all zeros and drops the time zone info
const SortableTimeFormat = "2006-01-02T15:04:05.000000000"

// Formats a time.Time into a []byte that can be sorted
func FormatTimeBytes(t time.Time) []byte {
	return []byte(t.UTC().Round(0).Format(SortableTimeFormat))
}

// Parses a []byte encoded using FormatTimeKey back into a time.Time
func ParseTimeBytes(bz []byte) (time.Time, error) {
	str := string(bz)
	t, err := time.Parse(SortableTimeFormat, str)
	if err != nil {
		return t, err
	}
	return t.UTC().Round(0), nil
}

// NewLevelDB instantiate a new LevelDB instance according to DBBackend.
func NewLevelDB(name, dir string) (db dbm.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't create db: %v", r)
		}
	}()
	return dbm.NewDB(name, backend, dir), err
}

type ParaMsg struct {
	HaveCosmosTxInBlock bool
	AnteErr             error
	RefundFee           Coins
	LogIndex            int
	HasRunEvmTx         bool
}

type IWatcher interface {
	SetFirstUse(v bool)
	Used()
	Enabled() bool
	GetEvmTxIndex() uint64
	NewHeight(height uint64, blockHash common.Hash, header tmtypes.Header)
	SaveContractCode(addr common.Address, code []byte)
	SaveContractCodeByHash(hash []byte, code []byte)
	SaveTransactionReceipt(status uint32, msg interface{}, txHash common.Hash, txIndex uint64, data interface{}, gasUsed uint64)
	UpdateCumulativeGas(txIndex, gasUsed uint64)
	SaveAccount(account interface{}, isDirectly bool)
	AddDelAccMsg(account interface{}, isDirectly bool)
	DeleteAccount(addr interface{})
	DelayEraseKey()
	ExecuteDelayEraseKey(delayEraseKey [][]byte)
	SaveState(addr common.Address, key, value []byte)
	SaveBlock(bloom ethtypes.Bloom)
	SaveLatestHeight(height uint64)
	SaveParams(params interface{})
	SaveContractBlockedListItem(addr interface{})
	SaveContractMethodBlockedListItem(addr interface{}, methods []byte)
	SaveContractDeploymentWhitelistItem(addr interface{})
	DeleteContractBlockedList(addr interface{})
	DeleteContractDeploymentWhitelist(addr interface{})
	Finalize()
	CommitCodeHashToDb(hash []byte, code []byte)
	Reset()
	Commit()
	CommitWatchData(data interface{})
}

type EmptyWatcher struct {
}

func (e EmptyWatcher) SetFirstUse(v bool) {}
func (e EmptyWatcher) Used()              {}
func (e EmptyWatcher) Enabled() bool {
	return false
}
func (e EmptyWatcher) GetEvmTxIndex() uint64 {
	return 0
}
func (e EmptyWatcher) NewHeight(height uint64, blockHash common.Hash, header tmtypes.Header) {}
func (e EmptyWatcher) SaveContractCode(addr common.Address, code []byte)                     {}
func (e EmptyWatcher) SaveContractCodeByHash(hash []byte, code []byte)                       {}
func (e EmptyWatcher) SaveTransactionReceipt(status uint32, msg interface{}, txHash common.Hash, txIndex uint64, data interface{}, gasUsed uint64) {
}
func (e EmptyWatcher) UpdateCumulativeGas(txIndex, gasUsed uint64)                        {}
func (e EmptyWatcher) SaveAccount(account interface{}, isDirectly bool)                   {}
func (e EmptyWatcher) AddDelAccMsg(account interface{}, isDirectly bool)                  {}
func (e EmptyWatcher) DeleteAccount(addr interface{})                                     {}
func (e EmptyWatcher) DelayEraseKey()                                                     {}
func (e EmptyWatcher) ExecuteDelayEraseKey(delayEraseKey [][]byte)                        {}
func (e EmptyWatcher) SaveState(addr common.Address, key, value []byte)                   {}
func (e EmptyWatcher) SaveBlock(bloom ethtypes.Bloom)                                     {}
func (e EmptyWatcher) SaveLatestHeight(height uint64)                                     {}
func (e EmptyWatcher) SaveParams(params interface{})                                      {}
func (e EmptyWatcher) SaveContractBlockedListItem(addr interface{})                       {}
func (e EmptyWatcher) SaveContractMethodBlockedListItem(addr interface{}, methods []byte) {}
func (e EmptyWatcher) SaveContractDeploymentWhitelistItem(addr interface{})               {}
func (e EmptyWatcher) DeleteContractBlockedList(addr interface{})                         {}
func (e EmptyWatcher) DeleteContractDeploymentWhitelist(addr interface{})                 {}
func (e EmptyWatcher) Finalize()                                                          {}
func (e EmptyWatcher) CommitCodeHashToDb(hash []byte, code []byte)                        {}
func (e EmptyWatcher) Reset()                                                             {}
func (e EmptyWatcher) Commit()                                                            {}
func (e EmptyWatcher) CommitWatchData(data interface{})                                   {}
