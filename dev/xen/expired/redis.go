package expired

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/okex/exchain/x/evm/statistics/rediscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	//	ttl := viper.GetInt64(flagTTL)

	filename := filepath.Join(viper.GetString(flagOutputDir), xenExpiredAddr)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rediscli.GetInstance().Init()
	db := rediscli.GetInstance().GetRawClient()
	//	var claims []*rediscli.XenMint

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
			log.Println(curseValues[1])
			values, err := redis.Values(curseValues[1], err)
			if err != nil {
				panic(fmt.Sprintf("get values error %v %v", curseValues[1], err))
			}
			for _, v := range values {
				key, err := redis.String(v, nil)
				if err != nil {
					panic(err)
				}
				if strings.Contains(key, "0x") {
					content, err := redis.Values(db.Do("HGETALL", key))
					if err != nil {
						panic(err)
					}
					for _, v := range content {
						s, _ := redis.String(v, nil)
						log.Println(s)
					}
				}
			}

			if curse == 0 {
				return
			}
		}
		//			if time.Now().Unix() > v.BlockTime.Add(time.Duration(*v.Term+ttl)*time.Duration(24)*time.Hour).Unix() {
		//				line := fmt.Sprintf("%v,%v,%v\n", v.ID, *v.Txhash, *v.Useraddr)
		//				_, err = f.WriteString(line)
		//				if err != nil {
		//					panic(err)
		//				}
		//			}
	}
}
