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
	conf                 map[string][]Round
	roleConf             map[string][]interface{}
	ProactivelyRunTxRole = "proactively-role"
	PreRunCase           = "prerun-testcase"
)

func init() {
	conf = make(map[string][]Round)
	roleConf = make(map[string][]interface{})
}

type Round struct {
	Id           int64
	Prevote      map[string]bool // role true => vote nil, false default vote
	Precommit    map[string]bool // role true => vote nil, false default vote
	Prerun       map[string]int  // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	Addblockpart map[string]int  // control receiver a block time
}

type ConfTest struct {
	proactivelyRunTx bool
}

func LoadTestConf() {
	role = fmt.Sprintf("v%s", viper.GetString(ProactivelyRunTxRole))
	confFilePath := viper.GetString(PreRunCase)
	content, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		fmt.Println("read fail", err)
		return
	}
	confTmp := make(map[string][]Round)
	json.Unmarshal(content, &confTmp)
	for k, v := range confTmp {
		if _, ok := roleConf[k]; !ok {
			for _, vInner := range v {
				var key = fmt.Sprintf("%s_%d", k, vInner.Id)
				if val, ok := vInner.Prevote[role]; ok {
					roleConf[key] = append(roleConf[key], val)
				} else {
					roleConf[key] = append(roleConf[key], false)
				}
				if val, ok := vInner.Precommit[role]; ok {
					roleConf[key] = append(roleConf[key], val)
				} else {
					roleConf[key] = append(roleConf[key], false)
				}
				if val, ok := vInner.Prerun[role]; ok {
					roleConf[key] = append(roleConf[key], val)
				} else {
					roleConf[key] = append(roleConf[key], 0)
				}
				if val, ok := vInner.Addblockpart[role]; ok {
					roleConf[key] = append(roleConf[key], val)
				} else {
					roleConf[key] = append(roleConf[key], 0)
				}

			}
		}
	}
}

func GetPrevote(height int64, round int) bool {
	prevote_nil := getConfDetail(height, round, 0).(bool)
	return prevote_nil
}

func GetPrecommit(height int64, round int) bool {
	precommit_nil :=  getConfDetail(height, round, 1).(bool)
	return precommit_nil
}

func PreTimeOut(height int64, round int) {
	time_sleep :=  getConfDetail(height, round, 2).(int)
	time.Sleep(time.Duration(time_sleep) * time.Second)

}

func AddBlockTimeOut(height int64, round int) {
	time_sleep :=  getConfDetail(height, round, 3).(int)
	time.Sleep(time.Duration(time_sleep) * time.Second)
}

func getConfDetail(height int64, round, kind int) interface{} {
	if v, ok := roleConf[fmt.Sprintf("%d_%d", height, round)]; !ok {
		return false
	} else {
		if len(v) < kind {
			return false
		}
		return v[kind]
	}
	return false
}
