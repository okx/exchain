package staking

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	etypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/staking/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgDeposit(t *testing.T) {
	//msgstr := `{"delegator_address":"cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity":{"denom":"okt","amount":"10"}}`

	msgstr := `{"delegator_address": "cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity": {"denom": "okt","amount": "1000"}}`
	//msgstr := `{"delegator_address": "cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity": {"denom": "okt","amount": "1000"}}`
	//msgstr := `{"delegator_address":"cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity":{"denom":"okt","amount":"10"}}`
	msg := &types.MsgDeposit{}
	err := json.Unmarshal([]byte(msgstr), msg)
	require.NoError(t, err)

	_, err = ConvertDepositMsg(nil, nil)
	require.NoError(t, err)

	//msgstr := `{"delegator_address":"cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity":{"denom":"okt","amount":"10"}}`

	msgstr = `{"delegator_address":"cosmos16xt864nptjuvvp9y2hwpmys839fjgg2vmdwnxj","quantity":{"denom":"okt","amount":"10"}}`
	msg = &types.MsgDeposit{}
	err = json.Unmarshal([]byte(msgstr), msg)
	require.NoError(t, err)

	_, err = ConvertDepositMsg(nil, nil)
	require.NoError(t, err)
}

func Test(t *testing.T) {
	str := "BB020A140000000000000000000000000000000000000000128002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002A2022A33F1DFE21FB0F80C661BEB4003B05D08C154B5727FD5FED4802C13B24EDDC"
	//str := "BB020A140000000000000000000000000000000000000000128002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002A200000000000000000000000000000000000000000000000000000000000000000"
	d, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	rd, err := etypes.DecodeResultData(d)
	if err != nil {
		panic(err)
	}
	fmt.Println(rd, rd.TxHash.String())
}

func TestCheckSignerAddress(t *testing.T) {
	addr1, err := sdk.AccAddressFromBech32("cosmos13z0m0xk9ajwpa6rdktfl8pta60227rpwtrcc6x")
	fmt.Println(addr1, addr1.String())
	require.NoError(t, err)
	addr2, err := sdk.AccAddressFromHex("889Fb79ac5Ec9C1Ee86Db2D3f3857Dd3D4af0C2E")
	fmt.Println(addr2, addr2.String())
	require.NoError(t, err)

	tr := common.CheckSignerAddress([]sdk.AccAddress{addr1}, []sdk.AccAddress{addr2})
	require.True(t, tr)
	fmt.Println(tr)
}
