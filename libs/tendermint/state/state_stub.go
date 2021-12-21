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
	PreVote      map[string]bool // true vote nil, false default vote
	PreCommit    map[string]bool // true vote nil, false default vote
	PreRun       map[string]int  // int => control prerun sleep time
	AddBlockPart map[string]int  // int => control sleep time before receiver a block
}

type action struct {
	preVote           bool // true vote nil, false default vote
	preCommit         bool // true vote nil, false default vote
	preRunWait        int  // control prerun sleep time
	addBlockPartWait int  // control sleep time before receiver a block
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

				act.preVote = event.PreVote[role]
				act.preCommit = event.PreCommit[role]
				act.preRunWait = event.PreRun[role]
				act.addBlockPartWait = event.AddBlockPart[role]

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
		tlog.Info("PrevoteNil", "height", height, "round", round, "act", act.preVote, )
		return act.preVote
	}
	return false
}

func PrecommitNil(height int64, round int) bool {
	if !enableRoleTest {
		return false
	}

	act, ok := roleAction[actionKey(height, round)]

	if ok {
		tlog.Info("PrecommitNil", "height", height, "round", round, "act", act.preCommit, )
		return act.preCommit
	}
	return false
}

func preTimeOut(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := roleAction[actionKey(height, round)]; ok {
		timeSleep := act.preRunWait
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}

func AddBlockTimeOut(height int64, round int) {
	if !enableRoleTest {
		return
	}
	if act, ok := roleAction[actionKey(height, round)]; ok {
		timeSleep := act.addBlockPartWait
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}

func actionKey(height int64, round int) string {
	return fmt.Sprintf("%d-%d", height, round)
}