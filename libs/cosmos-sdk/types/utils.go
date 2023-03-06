package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"

	dbm "github.com/okx/okbchain/libs/tm-db"
)

const FlagDBBackend = "db_backend"

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

// NewDB instantiate a new DB instance according to db_backend.
func NewDB(name, dir string) (db dbm.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't create db: %v", r)
		}
	}()

	backend := viper.GetString(FlagDBBackend)
	if len(backend) == 0 {
		backend = string(dbm.GoLevelDBBackend)
	}
	db = dbm.NewDB(name, dbm.BackendType(backend), dir)
	return
}

type ParaMsg struct {
	UseCurrentState     bool
	HaveCosmosTxInBlock bool
	AnteErr             error
	RefundFee           Coins
	LogIndex            int
	HasRunEvmTx         bool
}

type FeeSplitInfo struct {
	Addr   AccAddress
	Fee    Coins
	HasFee bool
}

type TxWatcher struct {
	IWatcher
}
type WatchMessage interface {
	GetKey() []byte
	GetValue() string
	GetType() uint32
}

type IWatcher interface {
	Enabled() bool
	SaveContractCode(addr common.Address, code []byte, height uint64)
	SaveContractCodeByHash(hash []byte, code []byte)
	SaveAccount(account interface{})
	DeleteAccount(account interface{})
	SaveState(addr common.Address, key, value []byte)
	SaveContractBlockedListItem(addr interface{})
	SaveContractMethodBlockedListItem(addr interface{}, methods []byte)
	SaveContractDeploymentWhitelistItem(addr interface{})
	DeleteContractBlockedList(addr interface{})
	DeleteContractDeploymentWhitelist(addr interface{})
	Finalize()
	Destruct() []WatchMessage
}

type EmptyWatcher struct {
}

func (e EmptyWatcher) Enabled() bool                                                      { return false }
func (e EmptyWatcher) SaveContractCode(addr common.Address, code []byte, height uint64)   {}
func (e EmptyWatcher) SaveContractCodeByHash(hash []byte, code []byte)                    {}
func (e EmptyWatcher) SaveAccount(account interface{})                                    {}
func (e EmptyWatcher) DeleteAccount(account interface{})                                  {}
func (e EmptyWatcher) SaveState(addr common.Address, key, value []byte)                   {}
func (e EmptyWatcher) SaveContractBlockedListItem(addr interface{})                       {}
func (e EmptyWatcher) SaveContractMethodBlockedListItem(addr interface{}, methods []byte) {}
func (e EmptyWatcher) SaveContractDeploymentWhitelistItem(addr interface{})               {}
func (e EmptyWatcher) DeleteContractBlockedList(addr interface{})                         {}
func (e EmptyWatcher) DeleteContractDeploymentWhitelist(addr interface{})                 {}
func (e EmptyWatcher) Finalize()                                                          {}
func (e EmptyWatcher) Destruct() []WatchMessage                                           { return nil }
