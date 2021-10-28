package types

var closeQueryMutex bool
var closeCheckTxMutex bool
var removeCheckTx bool

const (
	FlagCloseMutex        = "close-mutex"
	FlagCloseCheckTxMutex = "close-checktx-mutex"
	FlagRemoveCheckTx     = "remove-checktx"
)

func GetCloseMutex() bool {
	return closeQueryMutex
}

func SetCloseMutex(isClose bool) {
	closeQueryMutex = isClose
}

func GetCloseCheckTxMutex() bool {
	return closeCheckTxMutex
}

func SetCloseCheckTxMutex(isClose bool) {
	closeCheckTxMutex = isClose
}

func GetRemoveCheckTx() bool {
	return removeCheckTx
}

func SetRemoveCheckTx(isRemove bool) {
	removeCheckTx = isRemove
}
