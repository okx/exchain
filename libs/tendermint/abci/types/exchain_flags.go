package types

var disableABCIQueryMutex bool
var disableCheckTx bool

const (
	FlagCloseMutex            = "close-mutex"
	FlagDisableABCIQueryMutex = "disable-abci-query-mutex"
	FlagDisableCheckTx        = "disable-checktx"
)

func GetDisableABCIQueryMutex() bool {
	return disableABCIQueryMutex
}

func SetDisableABCIQueryMutex(isClose bool) {
	disableABCIQueryMutex = isClose
}

func GetDisableCheckTx() bool {
	return disableCheckTx
}

func SetDisableCheckTx(isRemove bool) {
	disableCheckTx = isRemove
}
