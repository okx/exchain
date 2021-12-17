package consensus

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
	roleConf             map[string][]bool
	ProactivelyRunTxRole = "proactively-role"
	PreRunCase           = "prerun-testcase"
)

func init() {
	conf = make(map[string][]Round)
	roleConf = make(map[string][]bool)
}

type Round struct {
	Id           int64
	Prevote      map[string]bool // role true => vote nil, false default vote
	Precommit    map[string]bool // role true => vote nil, false default vote
	PreRun       map[string]bool // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	Addblockpart map[string]bool // control receiver a block time
}

type ConfTest struct {
	proactivelyRunTx bool
}

func loadTestConf() {
	role = fmt.Sprintf("v%s", viper.GetString(ProactivelyRunTxRole))
	confFilePath := viper.GetString(PreRunCase)
	//fmt.Println("confFilePath --->", confFilePath)
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
				var prevote, precommit, preRun, addBlock bool
				if val, ok := vInner.Prevote[role]; ok {
					prevote = val
				}
				if val, ok := vInner.Precommit[role]; ok {
					precommit = val
				}
				if val, ok := vInner.PreRun[role]; ok {
					preRun = val
				}
				if val, ok := vInner.Addblockpart[role]; ok {
					addBlock = val
				}
				roleConf[key] = []bool{prevote, precommit, preRun, addBlock}
			}
		}
	}
	fmt.Println("roleConf===> ", roleConf)
}

func getPrevote(height int64, round int) bool {
	return getConfDetail(height, round, 0)
}

func getPrecommit(height int64, round int) bool {
	return getConfDetail(height, round, 1)
}

func preTimeOut(height int64, round int) {
	if getConfDetail(height, round, 2) {
		time.Sleep(2 * time.Second)
	}
}

func AddBlock(height int64, round int) {
	if getConfDetail(height, round, 3) {
		fmt.Println("AddBlock sleep" , height, round)
		time.Sleep(3 * time.Second)
	}
}

func getConfDetail(height int64, round, kind int) bool {
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
