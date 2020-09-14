package utils

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"runtime"
	"strconv"
	"sync"
)

// GoroutineID is the type of goroutine ID
type GoroutineID int

var goroutineSpace = []byte("goroutine ")

var littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
}

// GoID is the global variable for goroutine ID
var GoID GoroutineID

// String returns a human readable string representation of GoroutineID
func (base GoroutineID) String() string {
	bp := littleBuf.Get().(*[]byte)
	defer littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Extract the `1606` out of "goroutine 1606 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return "N/A"
	}
	b = b[:i]

	// check if it is valid
	s := string(b)
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return "WRFMT"
	}

	// 0 mean "used as it like"
	if int(base) == 0 {
		return s
	}
	return strconv.FormatUint(n, int(base))

}

// mulAndQuo returns a * b / c
func MulAndQuo(a, b, c sdk.Dec) sdk.Dec {
	// 10^8
	auxiliaryDec := sdk.NewDecFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil))
	a = a.Mul(auxiliaryDec)
	return a.Mul(b).Quo(c).Quo(auxiliaryDec)
}
