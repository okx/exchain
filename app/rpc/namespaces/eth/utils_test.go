package eth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TransformDataError(t *testing.T) {

	sdkerr := cosmosError{
		Code:      7,
		Log:       `["execution reverted","message","HexData","0x00000000000"]`,
		Codespace: "evm",
	}
	jsstr, _ := json.Marshal(sdkerr)
	sdkerr.Log = string(jsstr)
	err := TransformDataError(sdkerr, "eth_estimateGas")
	require.NotNil(t, err.ErrorData())
	require.Equal(t, err.ErrorData(), "0x00000000000")
	require.Equal(t, err.ErrorCode(), VMExecuteExceptionInEstimate)
	err = TransformDataError(sdkerr, "eth_call")
	require.NotNil(t, err.ErrorData())
	data, ok := err.ErrorData().(*wrapedEthError)
	require.True(t, ok)
	require.NotNil(t, data)
}
