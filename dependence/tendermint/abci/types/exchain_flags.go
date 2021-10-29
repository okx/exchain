package types

var disableQueryMutex bool
var disableCheckTxMutex bool
var disableCheckTx bool

const (
	FlagCloseMutex          = "close-mutex"
	FlagDisableQueryMutex   = "disable-query-mutex"
	FlagDisableCheckTxMutex = "disable-checktx-mutex"
	FlagDisableCheckTx      = "disable-checktx"
)

func GetDisableQueryMutex() bool {
	return disableQueryMutex
}

func SetDisableQueryMutex(isClose bool) {
	disableQueryMutex = isClose
}

func GetDisableCheckTxMutex() bool {
	return disableCheckTxMutex
}

func SetDisableCheckTxMutex(isClose bool) {
	disableCheckTxMutex = isClose
}

func GetDisableCheckTx() bool {
	return disableCheckTx
}

func SetDisableCheckTx(isRemove bool) {
	disableCheckTx = isRemove
}
