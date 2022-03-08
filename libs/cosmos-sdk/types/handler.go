package types

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"time"
)

// Handler defines the core of the state transition function of an application.
type Handler func(ctx Context, msg Msg) (*Result, error)

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx Context, tx Tx, simulate bool) (newCtx Context, err error)

type GasRefundHandler func(ctx Context, tx Tx) (fee Coins, err error)

type AccHandler func(ctx Context, address AccAddress) (nonce uint64)

type UpdateFeeCollectorAccHandler func(ctx Context, balance Coins) error

type LogFix func(isAnteFailed [][]string) (logs [][]byte)

type GetTxFeeHandler func(ctx Context, tx Tx) (Coins, bool, SigCache, *ethcommon.Address)

// AnteDecorator wraps the next AnteHandler to perform custom pre- and post-processing.
type AnteDecorator interface {
	AnteHandle(ctx Context, tx Tx, simulate bool, next AnteHandler) (newCtx Context, err error)
}

// ChainDecorator chains AnteDecorators together with each AnteDecorator
// wrapping over the decorators further along chain and returns a single AnteHandler.
//
// NOTE: The first element is outermost decorator, while the last element is innermost
// decorator. Decorator ordering is critical since some decorators will expect
// certain checks and updates to be performed (e.g. the Context) before the decorator
// is run. These expectations should be documented clearly in a CONTRACT docline
// in the decorator's godoc.
//
// NOTE: Any application that uses GasMeter to limit transaction processing cost
// MUST set GasMeter with the FIRST AnteDecorator. Failing to do so will cause
// transactions to be processed with an infinite gasmeter and open a DOS attack vector.
// Use `ante.SetUpContextDecorator` or a custom Decorator with similar functionality.
// Returns nil when no AnteDecorator are supplied.
func ChainAnteDecorators(chain ...AnteDecorator) AnteHandler {
	if len(chain) == 0 {
		return nil
	}

	// handle non-terminated decorators chain
	if (chain[len(chain)-1] != Terminator{}) {
		chain = append(chain, Terminator{})
	}

	return func(ctx Context, tx Tx, simulate bool) (Context, error) {
		return chain[0].AnteHandle(ctx, tx, simulate, ChainAnteDecorators(chain[1:]...))
	}
}

// Terminator AnteDecorator will get added to the chain to simplify decorator code
// Don't need to check if next == nil further up the chain
//                        ______
//                     <((((((\\\
//                     /      . }\
//                     ;--..--._|}
//  (\                 '--/\--'  )
//   \\                | '-'  :'|
//    \\               . -==- .-|
//     \\               \.__.'   \--._
//     [\\          __.--|       //  _/'--.
//     \ \\       .'-._ ('-----'/ __/      \
//      \ \\     /   __>|      | '--.       |
//       \ \\   |   \   |     /    /       /
//        \ '\ /     \  |     |  _/       /
//         \  \       \ |     | /        /
//   snd    \  \      \        /
type Terminator struct{}

const AnteTerminatorTag = "ante-terminator"

// Simply return provided Context and nil error
func (t Terminator) AnteHandle(ctx Context, _ Tx, _ bool, _ AnteHandler) (Context, error) {
	trc := ctx.AnteTracer()
	if trc != nil {
		trc.RepeatingPin(AnteTerminatorTag)
	}
	return ctx, nil
}

var (
	pLog = &ScfLog{
		paraAllTime: time.Duration(0),
		prePare:     time.Duration(0),
		runTx:       time.Duration(0),
		async:       time.Duration(0),

		conflictTime: time.Duration(0),
		mergeTime:    time.Duration(0),
		endTime:      time.Duration(0),
		fixTime:      time.Duration(0),
	}
)

type ScfLog struct {
	paraAllTime time.Duration
	prePare     time.Duration
	runTx       time.Duration
	async       time.Duration

	conflictTime time.Duration
	mergeTime    time.Duration
	endTime      time.Duration
	fixTime      time.Duration
}

func AddParaAllTIme(ts time.Duration) {
	pLog.paraAllTime += ts
}

func AddPrePare(ts time.Duration) {
	pLog.prePare += ts
}

func AddRunTx(ts time.Duration) {
	pLog.runTx += ts
}

func AddAsycn(ts time.Duration) {
	pLog.async += ts
}

func AddConflictTime(ts time.Duration) {
	pLog.conflictTime += ts
}

func AddMergeTime(ts time.Duration) {
	pLog.mergeTime += ts
}

func AddEndTime(ts time.Duration) {
	pLog.endTime += ts
}

func AddFixTime(ts time.Duration) {
	pLog.fixTime += ts
}

func PrintTime() {
	fmt.Println("ParaAllTime", pLog.paraAllTime.Seconds())
	fmt.Println("PrePare", pLog.prePare.Seconds())
	fmt.Println("RunTxs", pLog.runTx.Seconds())
	fmt.Println("Async", pLog.async.Seconds())

	fmt.Println("Conflict", pLog.conflictTime.Seconds())
	fmt.Println("MergeTime", pLog.mergeTime.Seconds())
	fmt.Println("EndTIme", pLog.endTime.Seconds())
	fmt.Println("FixTIme", pLog.fixTime.Seconds())
}
