package expired

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	minHeight = 15405261
	maxHeight = 16359432
)

func init() {
	expiredRedisV3Cmd.Flags().String(flagRedisAddr, ":6379", "redis addr")
	expiredRedisV3Cmd.Flags().String(flagRedisPassWord, "", "redis password")
	viper.BindPFlag(flagRedisAddr, expiredRedisV3Cmd.Flags().Lookup(flagRedisAddr))
	viper.BindPFlag(flagRedisPassWord, expiredRedisV3Cmd.Flags().Lookup(flagRedisPassWord))
}

func RedisV3Command() *cobra.Command {
	return expiredRedisV3Cmd
}

var expiredRedisV3Cmd = &cobra.Command{
	Use:   "expired_redis_v3",
	Short: "get the expired xen",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaimRedisV3()
		return nil
	},
}

func scanClaimRedisV3() {
	ttl := viper.GetInt64(flagTTL)

	filename := filepath.Join(viper.GetString(flagOutputDir), xenExpiredAddr)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	global.RedisAddr = viper.GetString(flagRedisAddr)
	global.RedisPassword = viper.GetString(flagRedisPassWord)
	rediscli.GetInstance().Init()
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()

	_, err = db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}

	counter := 0
	curse := viper.GetInt(flagCursor)
	for {
		curseValues, err := redis.Values(db.Do("SCAN", curse))
		if err != nil {
			panic(err)
		}
		if len(curseValues) > 0 {
			curse, err = redis.Int(curseValues[0], err)
			if err != nil {
				panic(fmt.Sprintf("got curse error %v %v", curseValues[0], err))
			}

			values, err := redis.Values(curseValues[1], err)
			if err != nil {
				panic(fmt.Sprintf("get values error %v %v", curseValues[1], err))
			}
			for _, v := range values {
				useraddr, err := redis.String(v, nil)
				if err != nil {
					panic(err)
				}
				if strings.Contains(useraddr, "0x") {
					unRewardKey := checkMintRewardIfNotEqualReturnTheLatestMint(useraddr)
					if unRewardKey != "" {
						claims := getLatestClaim(unRewardKey, useraddr[1:])
						if time.Now().Unix() > claims.BlockTime.Add(time.Duration(claims.Term+ttl)*time.Duration(24)*time.Hour).Unix() {
							reward := getLatestReward(claims.UserAddr[1:])
							if reward == nil ||
								(reward != nil && reward.BlockTime.Unix() < claims.BlockTime.Add(time.Duration(claims.Term)*time.Duration(24)*time.Hour).Unix()) {
								counter++
								line := fmt.Sprintf("%v,%v,%v,%v\n", counter, claims.TxHash, claims.TxSender, claims.UserAddr)
								_, err = f.WriteString(line)
								if err != nil {
									panic(err)
								}

							}
						}
					}
				}
			}

			if curse == 0 {
				return
			}
		}
	}
}

func checkMintRewardIfNotEqualReturnTheLatestMint(userAddr string) string {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}
	countMint, err := redis.Int(db.Do("ZCOUNT", userAddr, minHeight, maxHeight))
	if err != nil {
		panic(err)
	}

	_, err = db.Do("SELECT", 3)
	if err != nil {
		panic(err)
	}

	rewardUserAddr := "r" + userAddr[1:]
	countReward, err := redis.Int(db.Do("ZCOUNT", rewardUserAddr, minHeight, maxHeight))
	if err != nil {
		panic(err)
	}
	if countReward > countMint {
		panic(fmt.Sprintf("%v impossible!", userAddr))
	}
	if countReward == countMint {
		return ""
	}

	_, err = db.Do("SELECT", 2)
	if err != nil {
		panic(err)
	}
	v, err := redis.Strings(db.Do("ZREVRANGE", userAddr, 0, -1))
	if err != nil {
		panic(err)
	}

	return userAddr + "_" + v[0]
}

func getLatestReward(userAddr string) *rediscli.XenClaimReward {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 3)
	if err != nil {
		panic(err)
	}

	v, err := redis.Strings(db.Do("ZREVRANGE", "r"+userAddr, 0, -1))
	if err != nil {
		panic(err)
	}
	if len(v) == 0 {
		return nil
	}
	latestRewardKey := "r" + userAddr + "_" + v[0]

	content, err := redis.StringMap(db.Do("HGETALL", latestRewardKey))
	if err != nil {
		panic(err)
	}
	var reward rediscli.XenClaimReward
	for key, value := range content {
		parseReward(&reward, key, value)
	}

	return &reward
}

func getLatestClaim(key, userAddr string) *rediscli.XenMint {
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	_, err := db.Do("SELECT", 0)
	if err != nil {
		panic(err)
	}
	var ret rediscli.XenMint

	content, err := redis.StringMap(db.Do("HGETALL", key))
	if err != nil || len(content) == 0 {
		panic(fmt.Sprintf("%v map %v error %v", key, len(content), err))
	}
	for key, value := range content {
		parseClaim(&ret, key, value)
	}
	ret.UserAddr = userAddr

	return &ret
}
