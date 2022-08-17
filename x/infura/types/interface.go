package types

// ***********************************
type IStreamEngine interface {
	Write(data IStreamData) bool
}

type IStreamData interface {
	ConvertEngineData() EngineData
}

// Distributed State Service Interface
type IDistributeStateService interface {
	GetLockerID() string
	GetDistState(stateKey string) string
	SetDistState(stateKey string, stateValue string) error
	FetchDistLock(lockKey string, locker string, expiredInMS int) (bool, error)
	ReleaseDistLock(lockKey string, locker string) (bool, error)
	UnlockDistLockWithState(lockKey string, locker string, stateKey string, stateValue string) (bool, error)
}
