package types

const (
	RunTxModeCheck                 RunTxMode = iota // Check a transaction
	RunTxModeReCheck                                // Recheck a (pending) transaction after a commit
	RunTxModeSimulate                               // Simulate a transaction
	RunTxModeDeliver                                // Deliver a transaction
	RunTxModeDeliverInParallel                      // Deliver a transaction in parallel
	RunTxModeDeliverPartConcurrent                  // Deliver a transaction partial concurrent
	RunTxModeTrace                                  // Trace a transaction
	RunTxModeWrappedCheck
)

type RunTxMode uint8

func (m RunTxMode) String() (res string) {
	switch m {
	case RunTxModeCheck:
		res = "ModeCheck"
	case RunTxModeReCheck:
		res = "ModeReCheck"
	case RunTxModeSimulate:
		res = "ModeSimulate"
	case RunTxModeDeliver:
		res = "ModeDeliver"
	case RunTxModeDeliverPartConcurrent:
		res = "ModeDeliverPartConcurrent"
	case RunTxModeDeliverInParallel:
		res = "ModeDeliverInParallel"
	case RunTxModeWrappedCheck:
		res = "ModeWrappedCheck"
	default:
		res = "Unknown"
	}

	return res
}
