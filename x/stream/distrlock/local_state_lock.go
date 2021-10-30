package distrlock

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/okex/exchain/dependence/tendermint/libs/log"
)

// LocalStateService is designed to save stream state info into a local file.
// It's not supported to satisfy HA requirement.
// It mainly works when paired with LocalWebSocketEngine.
type LocalStateService struct {
	logger      log.Logger
	lockerID    string // unique identifier of locker
	lockFileDir string
	mutex       *sync.Mutex
}

func NewLocalStateService(logger log.Logger, lockerID string, lockFileDir string) (s *LocalStateService, err error) {

	_, err = os.Stat(lockFileDir)
	if err != nil {
		err = os.MkdirAll(lockFileDir, os.ModePerm)
	}

	if err == nil {
		s = &LocalStateService{
			logger:      logger,
			lockerID:    lockerID,
			lockFileDir: lockFileDir,
			mutex:       &sync.Mutex{},
		}
	}
	logger.Debug(fmt.Sprintf("NewLocalStateService lockerId: %s lockFileDir: %s", lockerID, lockFileDir))
	return s, err
}

func (s *LocalStateService) RemoveStateFile(stateKey string) error {
	path := s.getFullPath(stateKey)
	return os.Remove(path)
}

func (s *LocalStateService) getFullPath(stateName string) string {
	return s.lockFileDir + string(os.PathSeparator) + s.lockerID + "." + stateName
}

func (s *LocalStateService) GetLockerID() string {
	return s.lockerID
}

func (s *LocalStateService) GetDistState(stateKey string) (state string, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stateFilePath := s.getFullPath(stateKey)
	_, err = os.Stat(stateFilePath)
	if os.IsNotExist(err) {
		return "", nil
	}

	bytes, err := ioutil.ReadFile(stateFilePath)
	if err == nil {
		state = string(bytes)
	}

	return state, err
}

func (s *LocalStateService) SetDistState(stateKey string, stateValue string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	stateFilePath := s.getFullPath(stateKey)
	err := ioutil.WriteFile(stateFilePath, []byte(stateValue), 0600)

	return err
}

// DiskLock is not supported in LocalStateService, in the other word,
// FetchDistLock and ReleaseDistLock will always be success.
func (s *LocalStateService) FetchDistLock(lockKey string, locker string, expiredInMS int) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return true, nil
}

func (s *LocalStateService) ReleaseDistLock(lockKey string, locker string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return true, nil
}

func (s *LocalStateService) UnlockDistLockWithState(
	lockKey string, locker string, stateKey string, stateValue string) (bool, error) {
	err := s.SetDistState(stateKey, stateValue)
	return true, err
}
