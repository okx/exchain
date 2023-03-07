package types

import (
	"fmt"
	"github.com/okx/okbchain/libs/tendermint/libs/rand"
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
