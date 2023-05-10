package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewABI(t *testing.T) {
	astr := `[{"inputs":[{"internalType":"string","name":"data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	pabi, err := NewABI(astr)
	require.NoError(t, err)
	for key, e := range pabi.ABI.Methods {
		fmt.Println(key, e.String(), e.ID)
	}
}

func TestDecodeInputParam(t *testing.T) {
	abistr := `[{"inputs":[{"internalType":"string","name":"data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	noParamAbi := `[{"inputs":[],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	testcases := []struct {
		abistr  string
		fnInit  func(abis *ABI) (string, []byte)
		fnCheck func(ins []interface{}, err error)
	}{
		{
			abistr: abistr,
			fnInit: func(abis *ABI) (string, []byte) {
				return "invoke", []byte{1, 2}
			},
			fnCheck: func(ins []interface{}, err error) {
				require.Nil(t, ins)
				require.Error(t, err)
			},
		},
		{
			abistr: abistr,
			fnInit: func(abis *ABI) (string, []byte) {
				return "invoke1", []byte{1, 2, 3, 4, 5}
			},
			fnCheck: func(ins []interface{}, err error) {
				require.Nil(t, ins)
				require.Error(t, err)
			},
		},
		{ // args is ""
			abistr: abistr,
			fnInit: func(abis *ABI) (string, []byte) {
				re, err := abis.Pack("invoke", "")
				require.NoError(t, err)
				return "invoke", re
			},
			fnCheck: func(ins []interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, len(ins))
			},
		},
		{ // no param abi
			abistr: noParamAbi,
			fnInit: func(abis *ABI) (string, []byte) {
				re, err := abis.Pack("invoke")
				require.NoError(t, err)
				return "invoke", re
			},
			fnCheck: func(ins []interface{}, err error) {
				require.Nil(t, ins)
				require.Error(t, err)
			},
		},
		{ // normal
			abistr: abistr,
			fnInit: func(abis *ABI) (string, []byte) {
				re, err := abis.Pack("invoke", "123")
				require.NoError(t, err)
				return "invoke", re
			},
			fnCheck: func(ins []interface{}, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, len(ins))
				require.Equal(t, "123", ins[0].(string))
			},
		},
	}

	for _, ts := range testcases {
		abis, err := NewABI(ts.abistr)
		require.NoError(t, err)
		method, data := ts.fnInit(abis)
		res, err := abis.DecodeInputParam(method, data)
		ts.fnCheck(res, err)
	}
}

func TestDecodeInputParam1(t *testing.T) {
	abistr := `[{"inputs":[{"internalType":"string","name":"data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	abis, err := NewABI(abistr)
	require.NoError(t, err)
	for i := 0; i < 10; i++ {
		for j := 32; j < 1024; j = j + 32 {
			var data []byte
			data = append(data, abis.ABI.Methods["invoke"].ID...)
			data = append(data, rand.Bytes(j)...)
			_, err := abis.DecodeInputParam("invoke", data)
			require.Error(t, err)
		}
	}
}

func TestIsMatchFunction(t *testing.T) {
	abistr := `[{"inputs":[{"internalType":"string","name":"data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	abis, err := NewABI(abistr)
	require.NoError(t, err)

	testcases := []struct {
		method  string
		sign    []byte
		fnCheck func(in bool)
	}{
		{
			method: "invoke",
			sign:   nil,
			fnCheck: func(in bool) {
				require.False(t, in)
			},
		},
		{
			method: "invoke",
			sign:   []byte{1, 2, 3},
			fnCheck: func(in bool) {
				require.False(t, in)
			},
		},
		{
			method: "invoke1",
			sign:   abis.ABI.Methods["invoke"].ID,
			fnCheck: func(in bool) {
				require.False(t, in)
			},
		},
		{
			method: "invoke",
			sign:   abis.ABI.Methods["invoke"].ID,
			fnCheck: func(in bool) {
				require.True(t, in)
			},
		},
		{
			method: "invoke",
			sign:   []byte{1, 2, 3, 4},
			fnCheck: func(in bool) {
				require.False(t, in)
			},
		},
	}

	for _, ts := range testcases {
		abis.IsMatchFunction(ts.method, ts.sign)
	}
}

func TestABI_GetMethodById(t *testing.T) {
	abistr := `[{"inputs":[{"internalType":"string","name":"data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	abis, err := NewABI(abistr)
	require.NoError(t, err)

	testcases := []struct {
		method  string
		sign    []byte
		fnCheck func(method *abi.Method, err2 error)
	}{
		{
			method: "invoke normal",
			sign: func() []byte {
				buff, err := abis.ABI.Pack("invoke", "testdata")
				require.NoError(t, err)
				return buff
			}(),
			fnCheck: func(method *abi.Method, err2 error) {
				require.NoError(t, err2)
			},
		},
		{
			method: "invoke sign empty",
			sign: func() []byte {
				return nil
			}(),
			fnCheck: func(method *abi.Method, err2 error) {
				require.Error(t, err2)
			},
		},
		{
			method: "invoke sign method is not invoke",
			sign: func() []byte {
				buff, err := abis.ABI.Pack("invoke", "testdata")
				require.NoError(t, err)
				buff[0] += 0x1
				return buff
			}(),
			fnCheck: func(method *abi.Method, err2 error) {
				require.Error(t, err2)
			},
		},
		{
			method: "invoke sign data is err",
			sign: func() []byte {
				buff, err := abis.ABI.Pack("invoke", "testdata")
				require.NoError(t, err)
				length := len(buff)

				buff = buff[:length-1]
				return buff
			}(),
			fnCheck: func(method *abi.Method, err2 error) {
				require.Error(t, err2)
			},
		},
		{
			method: "invoke sign data is less than 4",
			sign: func() []byte {
				buff, err := abis.ABI.Pack("invoke", "testdata")
				require.NoError(t, err)

				buff = buff[:3]
				return buff
			}(),
			fnCheck: func(method *abi.Method, err2 error) {
				require.Error(t, err2)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnCheck(abis.GetMethodById(ts.sign))
	}
}
