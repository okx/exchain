package iavl

import (
	"fmt"
	"github.com/status-im/keycard-go/hexutils"
	"testing"
)

func TestInt2Byte(t *testing.T) {
	var a int = 5
	bytes := Int2Byte(a)
	fmt.Printf(hexutils.BytesToHex(bytes))
	var b []byte = []byte("0000000000000000000000000000")
	copy(b, bytes)
	fmt.Println(hexutils.BytesToHex(b))
}