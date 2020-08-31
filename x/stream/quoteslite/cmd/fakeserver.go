package main

import (
	"github.com/okex/okchain/x/stream/quoteslite"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	quoteslite.StartWSServer(logger, "0.0.0.0:6666")
}
