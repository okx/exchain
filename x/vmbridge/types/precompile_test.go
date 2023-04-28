package types

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeCallToWasmOutput(t *testing.T) {
	buff, err := hex.DecodeString("000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001d74686520726573756c74207761736d20636f6e74726163742064617461000000")
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
