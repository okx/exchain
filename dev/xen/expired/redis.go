package expired

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	flagCursor = "cursor"
)

func RedisCommand() *cobra.Command {
	return expiredRedisCmd
}

var expiredRedisCmd = &cobra.Command{
	Use:   "expired_redis",
	Short: "go the expired xen",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanClaimRedis()
		return nil
	},
}

func init() {
	expiredCmd.Flags().Int(flagCursor, 0, "curse")
	viper.BindPFlag(flagCursor, expiredCmd.Flags().Lookup(flagCursor))
}

func scanClaimRedis() {
	ttl := viper.GetInt64(flagTTL)

	filename := filepath.Join(viper.GetString(flagOutputDir), xenExpiredAddr)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rediscli.GetInstance().Init()
	pool := rediscli.GetInstance().GetClientPool()
	db := pool.Get()
	defer db.Close()
	var claims rediscli.XenMint

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
					content, err := redis.StringMap(db.Do("HGETALL", useraddr))
					if err != nil {
						panic(err)
					}
					for key, value := range content {
						parseClaim(&claims, key, value)
					}
					claims.UserAddr = useraddr
					if time.Now().Unix() > claims.BlockTime.Add(time.Duration(claims.Term+ttl)*time.Duration(24)*time.Hour).Unix() {
						counter++
						line := fmt.Sprintf("%v,%v,%v\n", counter, claims.TxHash, claims.UserAddr)
						_, err = f.WriteString(line)
						if err != nil {
							panic(err)
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

func parseClaim(claim *rediscli.XenMint, key, value string) {
	switch key {
	case "height":
		height, _ := strconv.Atoi(value)
		claim.Height = int64(height)
	case "txhash":
		claim.TxHash = value
	case "term":
		term, _ := strconv.Atoi(value)
		claim.Term = int64(term)
	case "txsender":
		claim.TxSender = value
	case "btime":
		utc, _ := strconv.Atoi(value)
		tim := time.Unix(int64(utc), 0)
		claim.BlockTime = tim
	}
}
