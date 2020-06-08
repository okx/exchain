package perf

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"sync"
	"time"
)

var (
	_ Perf = &performance{}
	_      = info{txNum: 0, beginBlockElapse: 0,
		endBlockElapse: 0, blockheight: 0, deliverTxElapse: 0}
)

const (
	marginModule       = "margin"
	orderModule        = "order"
	dexModule          = "dex"
	tokenModule        = "token"
	stakingModule      = "staking"
	govModule          = "gov"
	distributionModule = "distribution"
	summaryFormat      = "BlockHeight<%d>, " +
		"Abci<%dms>, " +
		"Tx<%d>, " +
		"%s"
	appFormat = "BlockHeight<%d>, " +
		"BeginBlock<%dms>, " +
		"DeliverTx<%dms>, " +
		"EndBlock<%dms>, " +
		"Commit<%dms>, " +
		"Tx<%d>" +
		"%s"
	moduleFormat = "BlockHeight<%d>, " +
		"module<%s>, " +
		"BeginBlock<%dms>, " +
		"DeliverTx<%dms>, " +
		"TxNum<%d>, " +
		"EndBlock<%dms>,"
	handlerFormat = "BlockHeight<%d>, " +
		"module<%s>, " +
		"handler<%s>, " +
		"elapsed<%dms>, " +
		"invoked<%d>,"
)

var perf *performance
var once sync.Once

// GetPerf gets the single instance of performance
func GetPerf() Perf {
	once.Do(func() {
		perf = newPerf()
	})
	return perf
}

// Perf shows the expected behaviour
type Perf interface {
	OnAppBeginBlockEnter(height int64) uint64
	OnAppBeginBlockExit(height int64, seq uint64)

	OnAppEndBlockEnter(height int64) uint64
	OnAppEndBlockExit(height int64, seq uint64)

	OnCommitEnter(height int64) uint64
	OnCommitExit(height int64, seq uint64, logger log.Logger)

	OnBeginBlockEnter(ctx sdk.Context, moduleName string) uint64
	OnBeginBlockExit(ctx sdk.Context, moduleName string, seq uint64)

	OnDeliverTxEnter(ctx sdk.Context, moduleName, handlerName string) uint64
	OnDeliverTxExit(ctx sdk.Context, moduleName, handlerName string, seq uint64)

	OnEndBlockEnter(ctx sdk.Context, moduleName string) uint64
	OnEndBlockExit(ctx sdk.Context, moduleName string, seq uint64)

	EnqueueMsg(msg string)
	EnableCheck()
}

type hanlderInfo struct {
	invoke uint64
	elapse int64
}

type info struct {
	blockheight      int64
	beginBlockElapse int64
	endBlockElapse   int64
	deliverTxElapse  int64
	txNum            uint64
}

type moduleInfo struct {
	info
	data handlerInfoMapType
}

type appInfo struct {
	info
	commitElapse  int64
	lastTimestamp int64
	seqNum        uint64
}

func (app *appInfo) abciElapse() int64 {
	return app.beginBlockElapse + app.endBlockElapse +
		app.deliverTxElapse + app.commitElapse
}

type handlerInfoMapType map[string]*hanlderInfo

func newHanlderMetrics() *moduleInfo {
	m := &moduleInfo{
		info: info{},
		data: make(handlerInfoMapType),
	}
	return m
}

type performance struct {
	lastTimestamp int64
	seqNum        uint64

	app           *appInfo
	moduleInfoMap map[string]*moduleInfo
	check         bool
	msgQueue      []string
}

func newPerf() *performance {
	p := &performance{
		moduleInfoMap: make(map[string]*moduleInfo),
	}

	p.app = &appInfo{
		info: info{},
	}
	p.moduleInfoMap[orderModule] = newHanlderMetrics()
	p.moduleInfoMap[dexModule] = newHanlderMetrics()
	p.moduleInfoMap[tokenModule] = newHanlderMetrics()
	p.moduleInfoMap[govModule] = newHanlderMetrics()
	p.moduleInfoMap[distributionModule] = newHanlderMetrics()
	p.moduleInfoMap[stakingModule] = newHanlderMetrics()
	p.moduleInfoMap[marginModule] = newHanlderMetrics()
	return p
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) EnableCheck() {
	p.check = true
}

func (p *performance) EnqueueMsg(msg string) {
	p.msgQueue = append(p.msgQueue, msg)
}

func (p *performance) OnAppBeginBlockEnter(height int64) uint64 {
	p.msgQueue = nil
	p.app.blockheight = height
	p.app.seqNum++
	p.app.lastTimestamp = time.Now().UnixNano()

	return p.app.seqNum
}

func (p *performance) OnAppBeginBlockExit(height int64, seq uint64) {
	p.sanityCheckApp(height, seq)
	p.app.beginBlockElapse = time.Now().UnixNano() - p.app.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnAppEndBlockEnter(height int64) uint64 {
	p.sanityCheckApp(height, p.app.seqNum)

	p.app.seqNum++
	p.app.lastTimestamp = time.Now().UnixNano()

	return p.app.seqNum
}

func (p *performance) OnAppEndBlockExit(height int64, seq uint64) {
	p.sanityCheckApp(height, seq)
	p.app.endBlockElapse = time.Now().UnixNano() - p.app.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnBeginBlockEnter(ctx sdk.Context, moduleName string) uint64 {
	p.lastTimestamp = time.Now().UnixNano()
	p.seqNum++

	m := p.getModule(moduleName)
	m.blockheight = ctx.BlockHeight()

	return p.seqNum
}

func (p *performance) OnBeginBlockExit(ctx sdk.Context, moduleName string, seq uint64) {
	p.sanityCheck(ctx, seq)
	m := p.getModule(moduleName)
	m.beginBlockElapse = time.Now().UnixNano() - p.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////
func (p *performance) OnEndBlockEnter(ctx sdk.Context, moduleName string) uint64 {
	p.lastTimestamp = time.Now().UnixNano()
	p.seqNum++

	m := p.getModule(moduleName)
	m.blockheight = ctx.BlockHeight()

	return p.seqNum
}

func (p *performance) OnEndBlockExit(ctx sdk.Context, moduleName string, seq uint64) {
	p.sanityCheck(ctx, seq)
	m := p.getModule(moduleName)

	m.endBlockElapse = time.Now().UnixNano() - p.lastTimestamp
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnDeliverTxEnter(ctx sdk.Context, moduleName, handlerName string) uint64 {

	m := p.getModule(moduleName)
	m.blockheight = ctx.BlockHeight()

	_, ok := m.data[handlerName]
	if !ok {
		m.data[handlerName] = &hanlderInfo{}
	}

	p.lastTimestamp = time.Now().UnixNano()
	p.seqNum++
	return p.seqNum
}

func (p *performance) OnDeliverTxExit(ctx sdk.Context, moduleName, handlerName string, seq uint64) {
	if !ctx.IsCheckTx() {
		p.sanityCheck(ctx, seq)
	}

	m := p.getModule(moduleName)

	info, ok := m.data[handlerName]
	if !ok {
		panic("Invalid handler name: " + handlerName)
	}
	info.invoke++
	info.elapse = time.Now().UnixNano() - p.lastTimestamp

	m.txNum++
	m.deliverTxElapse += info.elapse

	p.app.txNum++
	p.app.deliverTxElapse += info.elapse
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) OnCommitEnter(height int64) uint64 {
	p.sanityCheckApp(height, p.app.seqNum)

	p.app.lastTimestamp = time.Now().UnixNano()
	p.app.seqNum++
	return p.app.seqNum
}

func (p *performance) OnCommitExit(height int64, seq uint64, logger log.Logger) {
	p.sanityCheckApp(height, seq)
	// by millisecond
	unit := int64(1e6)
	p.app.commitElapse = time.Now().UnixNano() - p.app.lastTimestamp

	var moduleInfo string
	for moduleName, m := range p.moduleInfoMap {
		handlerElapse := m.deliverTxElapse / unit
		blockElapse := (m.beginBlockElapse + m.endBlockElapse) / unit
		if blockElapse == 0 && m.txNum == 0 {
			continue
		}
		moduleInfo += fmt.Sprintf(", %s[hdl<%dms>, blk<%dms>, tx<%d>]", moduleName, handlerElapse, blockElapse,
			m.txNum)

		logger.Info(fmt.Sprintf(moduleFormat, m.blockheight, moduleName, m.beginBlockElapse/unit, m.deliverTxElapse/unit,
			m.txNum, m.endBlockElapse/unit))

		for hanlderName, info := range m.data {
			logger.Info(fmt.Sprintf(handlerFormat, m.blockheight, moduleName, hanlderName, info.elapse/unit, info.invoke))
		}
	}

	logger.Info(fmt.Sprintf(appFormat, p.app.blockheight, p.app.beginBlockElapse/unit, p.app.deliverTxElapse/unit,
		p.app.endBlockElapse/unit, p.app.commitElapse/unit, p.app.txNum, moduleInfo))

	for _, e := range p.msgQueue {
		logger.Info(fmt.Sprintf(summaryFormat, p.app.blockheight, p.app.abciElapse()/unit, p.app.txNum, e))
	}

	p.msgQueue = nil

	p.app = &appInfo{seqNum: p.app.seqNum}
	p.moduleInfoMap[orderModule] = newHanlderMetrics()
	p.moduleInfoMap[dexModule] = newHanlderMetrics()
	p.moduleInfoMap[tokenModule] = newHanlderMetrics()
	p.moduleInfoMap[govModule] = newHanlderMetrics()
	p.moduleInfoMap[distributionModule] = newHanlderMetrics()
	p.moduleInfoMap[stakingModule] = newHanlderMetrics()
	p.moduleInfoMap[marginModule] = newHanlderMetrics()
}

////////////////////////////////////////////////////////////////////////////////////

func (p *performance) sanityCheck(ctx sdk.Context, seq uint64) {
	if !p.check {
		return
	}
	if seq != p.seqNum {
		panic("Invalid seq")
	}

	if ctx.BlockHeight() != p.app.blockheight {
		panic("Invalid block height")
	}
}

func (p *performance) sanityCheckApp(height int64, seq uint64) {
	if !p.check {
		return
	}

	if seq != p.app.seqNum {
		panic("Invalid seq")
	}

	if height != p.app.blockheight {
		panic("Invalid block height")
	}
}

func (p *performance) getModule(moduleName string) *moduleInfo {

	v, ok := p.moduleInfoMap[moduleName]
	if !ok {
		panic("Invalid module name: " + moduleName)
	}

	return v
}
