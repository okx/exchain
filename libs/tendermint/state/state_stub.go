package state

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

//-----------------------------------------------------------------------------
// Errors

//-----------------------------------------------------------------------------

var (
	tlog           log.Logger
	enableRoleTest bool
	roleAction     map[string]*action
)

const (
	ConsensusRole          = "consensus-role"
	ConsensusTestcase      = "consensus-testcase"
)

func init() {
	roleAction = make(map[string]*action)
}

type round struct {
	Round        int64
	Prevote      map[string]bool // role true => vote nil, true default vote
	Precommit    map[string]bool // role true => vote nil, true default vote
	Prerun       map[string]int  // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	Addblockpart map[string]int  // control receiver a block time
}

type action struct {
	Prevote           bool // role true => vote nil, false default vote
	Precommit         bool // role true => vote nil, false default vote
	PrerunWait        int  // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	AddblockpartnWait int  // control receiver a block time
}

func LoadTestConf(log log.Logger) {

	confFilePath := viper.GetString(ConsensusTestcase)
	if len(confFilePath) == 0 {
		return
	}

	tlog = log
	role := fmt.Sprintf("v%s", viper.GetString(ConsensusRole))

	content, err := ioutil.ReadFile(confFilePath)

	if err != nil {
		panic(fmt.Sprintf("read file : %s fail err : %s", confFilePath, err))
	}
	confTmp := make(map[string][]round)
	err = json.Unmarshal(content, &confTmp)
	if err != nil {
		panic(fmt.Sprintf("json Unmarshal err : %s", err))
	}

	enableRoleTest = true
	log.Info("Load consensus test case", "file", confFilePath, "err", err, "confTmp", confTmp)

	for height, roundEvents := range confTmp {
		if _, ok := roleAction[height]; !ok {
			for _, event := range roundEvents {
				act := &action{}

				act.Prevote = event.Prevote[role]
				act.Precommit = event.Precommit[role]
				act.PrerunWait = event.Prerun[role]
				act.AddblockpartnWait = event.Addblockpart[role]

				roleAction[fmt.Sprintf("%s-%d", height, event.Round)] = act
			}
		}
	}
}

func PrevoteNil(height int64, round int) bool {
	if !enableRoleTest {
		return false
	}
	act, ok := roleAction[fmt.Sprintf("%d_%d", height, round)]

	if ok {
		tlog.Info("PrevoteNil", "height", height, "round", round, "act", act.Prevote, )
		return act.Prevote
	}
	return false
}

func PrecommitNil(height int64, round int) bool {
	if !enableRoleTest {
		return false
	}

	act, ok := roleAction[fmt.Sprintf("%d_%d", height, round)]

	if ok {
		tlog.Info("PrecommitNil", "height", height, "round", round, "act", act.Precommit, )
		return act.Precommit
	}
	return false
}

func preTimeOut(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := roleAction[fmt.Sprintf("%d_%d", height, round)]; ok {
		time_sleep := act.PrerunWait
		time.Sleep(time.Duration(time_sleep) * time.Second)
	}
}

func AddBlockTimeOut(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := roleAction[fmt.Sprintf("%d_%d", height, round)]; ok {
		time_sleep := act.AddblockpartnWait
		time.Sleep(time.Duration(time_sleep) * time.Second)
	}
}
