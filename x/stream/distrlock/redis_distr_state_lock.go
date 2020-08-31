package distrlock

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/okex/okchain/x/stream/common"
	"github.com/tendermint/tendermint/libs/log"
)

var unlockScript = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

var unlockScriptWithState = redis.NewScript(1, `
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		redis.call("set", ARGV[2], ARGV[3])
		return redis.call("del", KEYS[1])
	else
		return -1
	end
`)

type RedisDistributeStateService struct {
	pool     *redis.Pool
	logger   log.Logger
	lockerId string // unique identifier of locker
}

func NewRedisDistributeStateService(redisUrl string, redisPass string, logger log.Logger, lockerId string) (*RedisDistributeStateService, error) {

	pool, err := common.NewPool(redisUrl, redisPass, logger)
	if err != nil {
		return nil, err
	}

	s := &RedisDistributeStateService{
		pool:     pool,
		logger:   logger,
		lockerId: lockerId,
	}

	return s, nil

}

func (s *RedisDistributeStateService) GetLockerId() string {
	return s.lockerId
}

func (s *RedisDistributeStateService) GetDistState(stateKey string) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	state, err := redis.String(conn.Do("GET", stateKey))
	s.logger.Debug(fmt.Sprintf("GetDistState: trying to get state with key(%s), state(%+v), err(%+v)", stateKey, state, err))
	if err == redis.ErrNil {
		return "", nil
	}

	return state, err
}

func (s *RedisDistributeStateService) SetDistState(stateKey string, stateValue string) error {
	conn := s.pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SET", stateKey, stateValue))
	s.logger.Debug(fmt.Sprintf("SetDistState: trying to set state(%s) with key(%s), err(%+v)", stateValue, stateKey, err))
	return err
}

func (s *RedisDistributeStateService) FetchDistLock(lockKey string, locker string, expiredInMS int) (bool, error) {
	conn := s.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("SET", lockKey, locker, "PX", expiredInMS, "NX")
	s.logger.Debug(fmt.Sprintf("FetchDistLock: trying to lock key(%s) with locker(%s) reply(%+v)", lockKey, locker, reply))
	return err == nil && reply == "OK", err
}

func (s *RedisDistributeStateService) ReleaseDistLock(lockKey string, locker string) (bool, error) {
	conn := s.pool.Get()
	defer conn.Close()

	replyStatus, err := unlockScript.Do(conn, lockKey, locker)
	s.logger.Debug(fmt.Sprintf("ReleaseDistLock: trying to release key(%s) with locker(%s), replyStatus(%T, %+v)",
		lockKey, locker, replyStatus, replyStatus))
	return err == nil && replyStatus == int64(1), err
}

func (s *RedisDistributeStateService) UnlockDistLockWithState(lockKey string, locker string, stateKey string, stateValue string) (bool, error) {
	conn := s.pool.Get()
	defer conn.Close()

	replyStatus, err := unlockScriptWithState.Do(conn, lockKey, locker, stateKey, stateValue)
	s.logger.Debug(fmt.Sprintf("UnlockDistLockWithState: trying to release key(%s) with locker(%s) and set stateKey: (%s), reply(%T, %+v)",
		lockKey, locker, stateKey, replyStatus, replyStatus))

	return err == nil && replyStatus == int64(1), err
}
