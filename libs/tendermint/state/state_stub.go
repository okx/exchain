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
	prevote           bool // role true => vote nil, false default vote
	precommit         bool // role true => vote nil, false default vote
	prerunWait        int  // true => preRun time less than consensus vote time , false => preRun time greater than consensus vote time
	addblockpartnWait int  // control receiver a block time
}

func loadTestCase(log log.Logger) {

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

				act.prevote = event.Prevote[role]
				act.precommit = event.Precommit[role]
				act.prerunWait = event.Prerun[role]
				act.addblockpartnWait = event.Addblockpart[role]

				roleAction[fmt.Sprintf("%s-%d", height, event.Round)] = act
			}
		}
	}
}

func PrevoteNil(height int64, round int) bool {
	if !enableRoleTest {
		return false
	}
	act, ok := roleAction[actionKey(height, round)]

	if ok {
		tlog.Info("PrevoteNil", "height", height, "round", round, "act", act.prevote, )
		return act.prevote
	}
	return false
}

func PrecommitNil(height int64, round int) bool {
	if !enableRoleTest {
		return false
	}

	act, ok := roleAction[actionKey(height, round)]

	if ok {
		tlog.Info("PrecommitNil", "height", height, "round", round, "act", act.precommit, )
		return act.precommit
	}
	return false
}

func preTimeOut(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := roleAction[actionKey(height, round)]; ok {
		timeSleep := act.prerunWait
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}

func AddBlockTimeOut(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := roleAction[actionKey(height, round)]; ok {
		timeSleep := act.addblockpartnWait
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}

func actionKey(height int64, round int) string {
	return fmt.Sprintf("%d-%d", height, round)
}