package types

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestGetEVMABIConfig(t *testing.T) {
	testCase := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
	}{
		{
			name:  "normal abi json",
			data:  abiJson,
			isErr: false,
		},
		{
			name:  "normal abi json have more func",
			data:  []byte("[\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"recipient\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"__OKBCSendToWasm\",\n    \"type\": \"event\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"caller\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"address\",\n        \"name\": \"recipient\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"mintERC20\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bool\",\n        \"name\": \"success\",\n        \"type\": \"bool\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"caller\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"address\",\n        \"name\": \"recipient\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"mintERC201\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bool\",\n        \"name\": \"success\",\n        \"type\": \"bool\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  }\n]\n"),
			isErr: false,
		},
		{
			name:  "normal abi json have more event",
			data:  []byte("[\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"recipient\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"__OKBCSendToWasm\",\n    \"type\": \"event\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"caller\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"address\",\n        \"name\": \"recipient\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"mintERC20\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bool\",\n        \"name\": \"success\",\n        \"type\": \"bool\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"recipient\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"__OKBCSendToWasmTest\",\n    \"type\": \"event\"\n  }\n]\n"),
			isErr: false,
		},
		{
			name:      "normal abi json have less event",
			data:      []byte("[\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"caller\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"address\",\n        \"name\": \"recipient\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"mintERC20\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bool\",\n        \"name\": \"success\",\n        \"type\": \"bool\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  }\n]\n"),
			isErr:     true,
			expectErr: "abi must have event event",
		},
		{
			name:  "normal abi json have less func",
			data:  []byte("[\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"recipient\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"__OKBCSendToWasm\",\n    \"type\": \"event\"\n  }\n]"),
			isErr: false,
		},
		{
			name:      "error abi json",
			data:      []byte("[\n  {\n    : false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"recipient\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"__OKBCSendToWasm\",\n    \"type\": \"event\"\n  }\n]"),
			isErr:     true,
			expectErr: "json decode failed",
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(tt *testing.T) {
			if tc.isErr {
				defer func() {
					r := recover()
					require.NotNil(tt, r)
					err := r.(error)
					require.ErrorContains(tt, err, tc.expectErr)
				}()
			}
			GetEVMABIConfig(tc.data)
		})
	}

}

func TestGetMintERC20Input(t *testing.T) {
	ethAddress := common.Address{0x1}
	addrStr := ethAddress.String()
	testCases := []struct {
		name      string
		caller    string
		recipient common.Address
		amount    *big.Int
		isErr     bool
		expectErr string
	}{
		{
			name:      "normal",
			caller:    addrStr,
			recipient: ethAddress,
			amount:    sdk.NewInt(1).BigInt(),
		},
	}

	for _, tc := range testCases {
		_, err := GetMintERC20Input(tc.caller, tc.recipient, tc.amount)
		if tc.isErr {
			require.Error(t, err)
		}
	}

}

func TestGetMintERC20Output(t *testing.T) {
	testCases := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
		expect    bool
	}{
		{
			name: "normal true 1",
			data: func() []byte {
				buffer := make([]byte, 31, 64)
				buffer = append(buffer, byte(0x1))
				return buffer
			}(),
			expect: true,
		},
		{
			name: "normal false",
			data: func() []byte {
				buffer := make([]byte, 31, 64)
				buffer = append(buffer, byte(0x0))
				return buffer
			}(),
			expect: false,
		},
		{
			name: "normal true 2",
			data: func() []byte {
				buffer := make([]byte, 32)
				buffer[31] = byte(0x1)
				return buffer
			}(),
			expect: true,
		},
		{
			name: "err data input no enough ",
			data: func() []byte {
				buffer := make([]byte, 28)
				return buffer
			}(),
			isErr:  true,
			expect: false,
		},
		{
			name: "err data input more ",
			data: func() []byte {
				buffer := make([]byte, 33)
				return buffer
			}(),
			isErr:  true,
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			result, err := GetMintERC20Output(tc.data)
			if tc.isErr {
				require.Error(tt, err)
			}
			require.Equal(tt, tc.expect, result)
		})
	}
}

func TestGetMintCW20Input(t *testing.T) {
	testCases := []struct {
		name      string
		amount    string
		reicient  string
		isErr     bool
		expectErr string
		expect    string
	}{
		{
			name: "normal true",
			reicient: func() string {
				addr := sdk.AccAddress{0x1}
				return addr.String()
			}(),
			amount: sdk.NewInt(1).String(),
			expect: "{\"mint_c_w20\":{\"amount\":\"1\",\"recipient\":\"cosmos1qyfkm2y3\"}}",
		},
		{
			name: "amount -1",
			reicient: func() string {
				addr := sdk.AccAddress{0x1}
				return addr.String()
			}(),
			amount: sdk.NewInt(-1).String(),
			expect: "{\"mint_c_w20\":{\"amount\":\"-1\",\"recipient\":\"cosmos1qyfkm2y3\"}}",
		},
		{
			name:     "addr is error",
			reicient: "hehe",
			amount:   sdk.NewInt(-1).String(),
			expect:   "{\"mint_c_w20\":{\"amount\":\"-1\",\"recipient\":\"hehe\"}}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			result, err := GetMintCW20Input(tc.amount, tc.reicient)
			if tc.isErr {
				require.Error(tt, err)
			}
			require.Equal(tt, tc.expect, string(result))
		})
	}
}
