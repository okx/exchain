package state

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

//-----------------------------------------------------------------------------
// Errors

//-----------------------------------------------------------------------------

var (
	role                 string
	roleConf             map[string]RoundRole
	ProactivelyRunTxRole = "proactively-role"
	PreRunCase           = "prerun-testcase"
)

func init() {
	roleConf = make(map[string]RoundRole)
}

type Round struct {
	Id           int64
	Prevote      map[string]bool // role true => vote nil, false default vote
	Precommit    map[string]bool // role true => vote nil, false default vote
	Prerun       map[string]int  // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	Addblockpart map[string]int  // control receiver a block time
}

type RoundRole struct {
	Prevote           bool // role true => vote nil, false default vote
	Precommit         bool // role true => vote nil, false default vote
	PrerunWait        int  // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	AddblockpartnWait int  // control receiver a block time
}

func LoadTestConf() {
	role = fmt.Sprintf("v%s", viper.GetString(ProactivelyRunTxRole))
	confFilePath := viper.GetString(PreRunCase)
	content, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		fmt.Printf("read file : %s fail err : %s\n", confFilePath, err)
		return
	}
	confTmp := make(map[string][]Round)
	err = json.Unmarshal(content, &confTmp)
	if err != nil {
		fmt.Printf("json Unmarshal err : %s\n", err)
		return
	}

	for height, v := range confTmp {
		if _, ok := roleConf[height]; !ok {
			for _, vInner := range v {
				round := vInner.Id
				tmp := RoundRole{}

				if val, ok := vInner.Prevote[role]; ok {
					tmp.Prevote = val
				}
				if val, ok := vInner.Precommit[role]; ok {
					tmp.Precommit = val
				}

				if val, ok := vInner.Prerun[role]; ok {
					tmp.PrerunWait = val
				}

				if val, ok := vInner.Addblockpart[role]; ok {
					tmp.AddblockpartnWait = val
				}
				roleConf[fmt.Sprintf("%s_%d", height, round)] = tmp
			}
		}
	}
}

func GetPrevote(height int64, round int) bool {
	if val, ok := roleConf[fmt.Sprintf("%d_%d", height, round)]; ok {
		return val.Prevote
	}
	return false
}

func GetPrecommit(height int64, round int) bool {
	if val, ok := roleConf[fmt.Sprintf("%d_%d", height, round)]; ok {
		return val.Precommit
	}
	return false
}

func PreTimeOut(height int64, round int) {
	if val, ok := roleConf[fmt.Sprintf("%d_%d", height, round)]; ok {
		time_sleep := val.PrerunWait
		time.Sleep(time.Duration(time_sleep) * time.Second)
	}
}

func AddBlockTimeOut(height int64, round int) {
	if val, ok := roleConf[fmt.Sprintf("%d_%d", height, round)]; ok {
		time_sleep := val.AddblockpartnWait
		time.Sleep(time.Duration(time_sleep) * time.Second)
	}
}
