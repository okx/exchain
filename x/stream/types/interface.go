package types

import (
	"sort"
)

type EngineEnqueue func(data IStreamData) bool

// ***********************************
type IStreamEngine interface {
	//Enqueue(scheduler IScheduler)
	Url() string
	Write(data IStreamData, success *bool)
}

type StreamDataKind byte

const (
	StreamDataNilKind       StreamDataKind = 0x00
	StreamDataAnalysisKind  StreamDataKind = 0x01
	StreamDataNotifyKind    StreamDataKind = 0x02
	StreamDataKlineKind     StreamDataKind = 0x03
	StreamDataWebSocketKind StreamDataKind = 0x04
)

// ***********************************
type IStreamData interface {
	BlockHeight() int64
	DataType() StreamDataKind
}

type IStreamDatas []IStreamData

//nolint
func (datas IStreamDatas) Len() int           { return len(datas) }
func (datas IStreamDatas) Less(i, j int) bool { return datas[i].BlockHeight() < datas[j].BlockHeight() }
func (datas IStreamDatas) Swap(i, j int)      { datas[i], datas[j] = datas[j], datas[i] }

// Sort is a helper function to sort the set of IStreamDatas in-place.
func (datas IStreamDatas) Sort() IStreamDatas {
	sort.Sort(datas)
	return datas
}

// Distributed State Service Interface
type IDistributeStateService interface {
	GetLockerId() string
	GetDistState(stateKey string) (string, error)
	SetDistState(stateKey string, stateValue string) error
	FetchDistLock(lockKey string, locker string, expiredInMS int) (bool, error)
	ReleaseDistLock(lockKey string, locker string) (bool, error)
	UnlockDistLockWithState(lockKey string, locker string, stateKey string, stateValue string) (bool, error)
}
