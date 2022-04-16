package watcher

import (
	"github.com/okex/exchain/x/evm/watcher/abci"
)

func (w *Watcher) ReceiveABCIMessage(deliverTx *abci.DeliverTx) {
	select {
	case w.txChan <- deliverTx:
	default:
		w.log.Info("watch db save deliver tx too busy")
		go func() { w.txChan <- deliverTx }()
	}
}

func (w *Watcher) saveTxAndReceipt() {

}

func (w *Watcher) parseABCIMessage(deliverTx *abci.DeliverTx) {

}
