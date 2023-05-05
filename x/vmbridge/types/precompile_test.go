package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeCallToWasmOutput(t *testing.T) {
	buff, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000407b22616d6f756e74223a7b2264656e6f6d223a226f6b74222c22616d6f756e74223a223939393839393939393938343536323433303030303030303030227d7d")
	require.NoError(t, err)
	pack, err := PreCompileABI.Methods[PrecompileCallToWasm].Outputs.Unpack(buff)
	require.NoError(t, err)
	if len(pack) != 1 {
		t.Log("err pack", pack)
		return
	}
	t.Log("pack", pack[0].(string))
}

func TestEncodeCallToWasmOutput(t *testing.T) {
	buff, err := EncodePrecompileCallToWasmOutput("the result wasm contract data")
	require.NoError(t, err)
	t.Log("ouput", hex.EncodeToString(buff))
}
