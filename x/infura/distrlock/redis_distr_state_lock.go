package distrlock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
)

var unlockScript = redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

var unlockScriptWithState = redis.NewScript(`
	if redis.call("get", KEYS[1]) == ARGV[1]
	then
		redis.call("set", KEYS[2], ARGV[2])
		return redis.call("del", KEYS[1])
	else
		return 0
	end
`)

type RedisDistributeStateService struct {
	client   *redis.Client
	logger   log.Logger
	lockerID string // unique identifier of locker
}

func NewRedisDistributeStateService(url string, pass string, db int, logger log.Logger, lockerID string) (*RedisDistributeStateService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pass, // no password set
		DB:       db,   // use select DB
	})

	s := &RedisDistributeStateService{
		client:   client,
		logger:   logger,
		lockerID: lockerID,
	}

	return s, nil
}

func (s *RedisDistributeStateService) GetLockerID() string {
	return s.lockerID
}

func (s *RedisDistributeStateService) GetDistState(stateKey string) string {
	state, _ := s.client.Get(context.Background(), stateKey).Result()
	return state
}

func (s *RedisDistributeStateService) SetDistState(stateKey string, stateValue string) error {
	err := s.client.Set(context.Background(), stateKey, stateValue, 0).Err()
	return err
}

func (s *RedisDistributeStateService) FetchDistLock(lockKey string, locker string, expiredInMS int) (bool, error) {
	success, err := s.client.SetNX(context.Background(), lockKey, locker,
		time.Duration(expiredInMS)*time.Millisecond).Result()
	return success, err
}

func (s *RedisDistributeStateService) ReleaseDistLock(lockKey string, locker string) (bool, error) {
	replyStatus, err := unlockScript.Run(context.Background(), s.client, []string{lockKey}, locker).Int()
	return err == nil && replyStatus == 1, err
}

func (s *RedisDistributeStateService) UnlockDistLockWithState(lockKey string, locker string, stateKey string, stateValue string) (bool, error) {
	replyStatus, err := unlockScriptWithState.Run(context.Background(), s.client, []string{lockKey, stateKey}, locker, stateValue).Int()
	return err == nil && replyStatus == 1, err
}
