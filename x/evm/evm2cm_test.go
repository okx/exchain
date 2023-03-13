package evm

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	etypes "github.com/okx/okbchain/x/evm/types"
	"github.com/stretchr/testify/require"
)

type testMsg struct {
	route string
}

func (msg testMsg) Route() string                { return msg.route }
func (msg testMsg) Type() string                 { return "testMsg" }
func (msg testMsg) GetSigners() []sdk.AccAddress { return nil }
func (msg testMsg) GetSignBytes() []byte         { return nil }
func (msg testMsg) ValidateBasic() error         { return nil }

func TestEvmConvertJudge(t *testing.T) {

	receiver := common.HexToAddress("0x606D5a30c22cEf66ed3C596b69a94200173e48b4")
	payLoad := sysABIParser.Methods[sysContractInvokeFunction].ID

	testcases := []struct {
		msg     sdk.Msg
		fnCheck func(addr []byte, is bool)
	}{
		{
			msg: testMsg{route: "testMsg"},
			fnCheck: func(addr []byte, is bool) {
				require.Nil(t, addr)
				require.Equal(t, false, is)
			},
		},
		{
			msg: testMsg{route: "evm"},
			fnCheck: func(addr []byte, is bool) {
				require.Nil(t, addr)
				require.Equal(t, false, is)
			},
		},
		{
			// deploy contract
			msg: etypes.NewMsgEthereumTxContract(1, nil, 1, nil, nil),
			fnCheck: func(addr []byte, is bool) {
				require.Nil(t, addr)
				require.Equal(t, false, is)
			},
		},
		{
			msg: etypes.NewMsgEthereumTx(1, &receiver, nil, 1, nil, nil),
			fnCheck: func(addr []byte, is bool) {
				require.Nil(t, addr)
				require.Equal(t, false, is)
			},
		},
		{
			msg: etypes.NewMsgEthereumTx(1, &receiver, nil, 1, nil, payLoad),
			fnCheck: func(addr []byte, is bool) {
				require.Equal(t, receiver[:], addr)
				require.Equal(t, true, is)
			},
		},
	}

	for _, ts := range testcases {
		addr, is := EvmConvertJudge(ts.msg)
		ts.fnCheck(addr, is)
	}
}

func TestParseContractParam(t *testing.T) {
	testcases := []struct {
		fnInit  func() []byte
		fnCheck func(ctm []byte, err error)
	}{
		{
			fnInit: func() []byte {
				return sysABIParser.Methods[sysContractInvokeFunction].ID
			},
			fnCheck: func(ctm []byte, err error) {
				require.Nil(t, ctm)
				require.NotNil(t, err)
			},
		},
		{
			fnInit: func() []byte { // input param number error
				param1 := `126`
				param2 := "123"
				// two input param for invoke
				abistr := `[{"inputs":[{"internalType":"string","name":"data","type":"string"},{"internalType":"string","name":"len","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
				abis, err := etypes.NewABI(abistr)
				require.NoError(t, err)
				re, err := abis.Pack(sysContractInvokeFunction, param1, param2)
				require.NoError(t, err)
				return re
			},
			fnCheck: func(ctm []byte, err error) {
				require.Nil(t, ctm)
				require.NotNil(t, err)
			},
		},
		{
			fnInit: func() []byte {
				abistr := `[{"inputs":[{"internalType":"int256","name":"len","type":"int256"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
				abis, err := etypes.NewABI(abistr)
				require.NoError(t, err)

				re, err := abis.Pack(sysContractInvokeFunction, big.NewInt(123))
				require.NoError(t, err)
				return re
			},
			fnCheck: func(ctm []byte, err error) {
				require.Nil(t, ctm)
				require.NotNil(t, err)
			},
		},
		{
			fnInit: func() []byte {
				param := `{"module": "staking","function": "deposit","data": "123"}`
				re, err := sysABIParser.Pack(sysContractInvokeFunction, hex.EncodeToString([]byte(param)))
				require.NoError(t, err)
				return re
			},
			fnCheck: func(ctm []byte, err error) {
				require.NoError(t, err)
				require.Equal(t, `{"module": "staking","function": "deposit","data": "123"}`, string(ctm))
			},
		},
	}

	for _, ts := range testcases {
		input := ts.fnInit()
		cmtp, err := ParseContractParam(input)
		ts.fnCheck(cmtp, err)
	}
}

func TestDecodeParam(t *testing.T) {
	testcases := []struct {
		input   string
		fnCheck func(ctm []byte, err error)
	}{
		{
			input: "7b226d6f64756c65223a20227374616b696e67222c2266756e6374696f6e223a20226465706f736974222c2264617461223a20227b5c2264656c656761746f725f616464726573735c223a205c223078623239313065323262623233643132396330326431323262373762343632636562306538396462395c222c5c227175616e746974795c223a207b5c2264656e6f6d5c223a205c226f6b745c222c5c22616d6f756e745c223a205c22315c227d7d227d",
			fnCheck: func(ctm []byte, err error) {
				require.NoError(t, err)
				data := `{"module": "staking","function": "deposit","data": "{\"delegator_address\": \"0xb2910e22bb23d129c02d122b77b462ceb0e89db9\",\"quantity\": {\"denom\": \"okt\",\"amount\": \"1\"}}"}`
				require.Equal(t, string(ctm), data)
			},
		},
		{
			input: "7b2",
			fnCheck: func(ctm []byte, err error) {
				require.NotNil(t, err)
			},
		},
		{
			// {"module": "staking","function1": "deposit","data": ""}
			input: "7b226d6f64756c65223a20227374616b696e67222c2266756e6374696f6e31223a20226465706f736974222c2264617461223a2022227d",
			fnCheck: func(ctm []byte, err error) {
				require.NoError(t, err)
				data := `{"module": "staking","function1": "deposit","data": ""}`
				require.Equal(t, string(ctm), data)
			},
		},
	}

	for _, tc := range testcases {
		ctm, err := DecodeParam([]byte(tc.input))
		tc.fnCheck(ctm, err)
	}
}

func TestDecodeResultData(t *testing.T) {
	strs := []string{
		"BB020A140000000000000000000000000000000000000000128002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002A2022A33F1DFE21FB0F80C661BEB4003B05D08C154B5727FD5FED4802C13B24EDDC",
		"BB020A140000000000000000000000000000000000000000128002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002A200000000000000000000000000000000000000000000000000000000000000000",
	}

	for _, str := range strs {
		d, err := hex.DecodeString(str)
		require.NoError(t, err)
		if err != nil {
			panic(err)
		}
		rd, err := etypes.DecodeResultData(d)
		require.NoError(t, err)
		t.Log(rd)
	}
}
