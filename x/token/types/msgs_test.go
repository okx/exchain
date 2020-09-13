package types

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestNewMsgTokenIssue(t *testing.T) {
	priKey := secp256k1.GenPrivKey()
	pubKey := priKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	totalSupply := "20000"

	testCase := []struct {
		issueMsg MsgTokenIssue
		err      sdk.Error
	}{
		{NewMsgTokenIssue("bnb", "bnb", "bnb", "binance coin", totalSupply, addr, true),
			nil},
		{NewMsgTokenIssue("", "", "", "binance coin", totalSupply, addr, true),
			sdk.ErrUnknownRequest("failed to check issue msg because original symbol cannot be empty")},
		{NewMsgTokenIssue("bnb", "bnb", "bnb", "binance 278343298$%%^&  coin", totalSupply, addr, true),
			sdk.ErrUnknownRequest("failed to check issue msg because invalid wholename")},
		{NewMsgTokenIssue("bnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbnbbn", "bnb", "bnb", "binance coin", totalSupply, addr, true),
			sdk.ErrUnknownRequest("failed to check issue msg because invalid desc")},
		{NewMsgTokenIssue("bnb", "bnb", "bnb", "binance coin", strconv.FormatInt(int64(99*1e10), 10), addr, true),
			sdk.ErrUnknownRequest("failed to check issue msg because invalid total supply")},
		{NewMsgTokenIssue("", "", "", "binance coin", totalSupply, sdk.AccAddress{}, true),
			sdk.ErrInvalidAddress(sdk.AccAddress{}.String())},
		{NewMsgTokenIssue("", "", "bnb-asd", "binance coin", totalSupply, addr, true),
			sdk.ErrUnknownRequest("failed to check issue msg because invalid original symbol: bnb-asd")},
	}

	for _, msgCase := range testCase {
		err := msgCase.issueMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	tokenIssueMsg := testCase[0].issueMsg
	signAddr := tokenIssueMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{addr}, signAddr)

	bz := ModuleCdc.MustMarshalJSON(tokenIssueMsg)
	require.EqualValues(t, sdk.MustSortJSON(bz), tokenIssueMsg.GetSignBytes())
	require.EqualValues(t, "token", tokenIssueMsg.Route())
	require.EqualValues(t, "issue", tokenIssueMsg.Type())
}

func TestNewMsgTokenBurn(t *testing.T) {
	priKey := secp256k1.GenPrivKey()
	pubKey := priKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	decCoin := sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100))

	decCoin0 := decCoin
	decCoin0.Denom = ""

	decCoin1 := decCoin
	decCoin1.Denom = "1okb-ads"

	testCase := []struct {
		burnMsg MsgTokenBurn
		err     sdk.Error
	}{
		{NewMsgTokenBurn(decCoin, addr), nil},
		{NewMsgTokenBurn(decCoin0, addr), sdk.ErrInsufficientCoins("failed to check burn msg because invalid Coins: 100.00000000")},
		{NewMsgTokenBurn(decCoin, sdk.AccAddress{}), sdk.ErrInvalidAddress(sdk.AccAddress{}.String())},
		{NewMsgTokenBurn(decCoin1, addr), sdk.ErrInsufficientCoins("failed to check burn msg because invalid Coins: 100.000000001okb-ads")},
	}

	for _, msgCase := range testCase {
		err := msgCase.burnMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	tokenBurnMsg := testCase[0].burnMsg
	signAddr := tokenBurnMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{addr}, signAddr)

	bz := ModuleCdc.MustMarshalJSON(tokenBurnMsg)
	require.EqualValues(t, sdk.MustSortJSON(bz), tokenBurnMsg.GetSignBytes())
	require.EqualValues(t, "token", tokenBurnMsg.Route())
	require.EqualValues(t, "burn", tokenBurnMsg.Type())

	err := tokenBurnMsg.ValidateBasic()
	require.NoError(t, err)
}

//tokenMintMsg := NewMsgTokenMint("btc", mintNum, testAccounts[0].baseAccount.Address)
func TestNewMsgTokenMint(t *testing.T) {
	priKey := secp256k1.GenPrivKey()
	pubKey := priKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	decCoin := sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(1000))
	decCoin0 := decCoin
	decCoin0.Denom = ""

	decCoin1 := decCoin
	decCoin1.Denom = "11234"

	decCoin2 := decCoin
	decCoin2.Amount = sdk.NewDec(TotalSupplyUpperbound + 1)

	testCase := []struct {
		mintMsg MsgTokenMint
		err     sdk.Error
	}{
		{NewMsgTokenMint(decCoin, addr), nil},
		{NewMsgTokenMint(decCoin0, addr), sdk.ErrInsufficientCoins("failed to check mint msg because invalid Coins: 1000.00000000")},
		{NewMsgTokenMint(decCoin, sdk.AccAddress{}), sdk.ErrInvalidAddress(sdk.AccAddress{}.String())},
		{NewMsgTokenMint(decCoin1, addr), sdk.ErrInsufficientCoins("failed to check mint msg because invalid Coins: 1000.0000000011234")},
		{NewMsgTokenMint(decCoin2, addr), sdk.ErrUnknownRequest("failed to check mint msg because invalid amount")},
	}

	for _, msgCase := range testCase {
		err := msgCase.mintMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	tokenMintMsg := testCase[0].mintMsg
	signAddr := tokenMintMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{addr}, signAddr)

	bz := ModuleCdc.MustMarshalJSON(tokenMintMsg)
	require.EqualValues(t, sdk.MustSortJSON(bz), tokenMintMsg.GetSignBytes())
	require.EqualValues(t, "token", tokenMintMsg.Route())
	require.EqualValues(t, "mint", tokenMintMsg.Type())

	err := tokenMintMsg.ValidateBasic()
	require.NoError(t, err)
}

func TestNewTokenMsgSend(t *testing.T) {
	// from
	fromPriKey := secp256k1.GenPrivKey()
	fromPubKey := fromPriKey.PubKey()
	fromAddr := sdk.AccAddress(fromPubKey.Address())

	// to
	toPriKey := secp256k1.GenPrivKey()
	toPubKey := toPriKey.PubKey()
	toAddr := sdk.AccAddress(toPubKey.Address())

	coins := sdk.DecCoins{
		sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(100)),
	}

	Errorcoins := sdk.DecCoins{
		sdk.NewDecCoinFromDec("okc", sdk.NewDec(100)),
		sdk.NewDecCoinFromDec("okc", sdk.NewDec(100)),
		sdk.NewDecCoinFromDec("oke", sdk.NewDec(100)),
	}

	// not valid coins
	decCoin := sdk.DecCoin{
		Denom:  "",
		Amount: sdk.NewDec(100),
	}
	notValidCoins := sdk.DecCoins{
		decCoin,
	}

	testCase := []struct {
		sendMsg MsgSend
		err     sdk.Error
	}{
		{NewMsgTokenSend(fromAddr, toAddr, coins), nil},
		{NewMsgTokenSend(fromAddr, toAddr, sdk.DecCoins{}), sdk.ErrInsufficientCoins("failed to check send msg because send amount must be positive")},
		{NewMsgTokenSend(fromAddr, toAddr, Errorcoins), sdk.ErrInvalidCoins("failed to check send msg because send amount is invalid: 100.00000000okc,100.00000000okc,100.00000000oke")},
		{NewMsgTokenSend(sdk.AccAddress{}, toAddr, coins), sdk.ErrInvalidAddress("failed to check send msg because miss sender address")},
		{NewMsgTokenSend(fromAddr, sdk.AccAddress{}, coins), sdk.ErrInvalidAddress("failed to check send msg because miss recipient address")},
		{NewMsgTokenSend(fromAddr, toAddr, notValidCoins), sdk.ErrInvalidCoins("failed to check send msg because send amount is invalid: 100.00000000")},
	}
	for _, msgCase := range testCase {
		err := msgCase.sendMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	tokenSendMsg := testCase[0].sendMsg
	signAddr := tokenSendMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{fromAddr}, signAddr)
	require.EqualValues(t, RouterKey, tokenSendMsg.Route())
	require.EqualValues(t, "send", tokenSendMsg.Type())

	bz := ModuleCdc.MustMarshalJSON(tokenSendMsg)
	require.EqualValues(t, sdk.MustSortJSON(bz), tokenSendMsg.GetSignBytes())

	err := tokenSendMsg.ValidateBasic()
	require.NoError(t, err)
}

func TestNewTokenMultiSend(t *testing.T) {
	// from
	fromPriKey := secp256k1.GenPrivKey()
	fromPubKey := fromPriKey.PubKey()
	fromAddr := sdk.AccAddress(fromPubKey.Address())

	// correct message
	coinStr := `[{"to":"okexchain1dfpljpe0g0206jch32fx95lyagq3z5ws850m6f","amount":"1` + common.NativeToken + `"}]`
	transfers, err := StrToTransfers(coinStr)
	require.Nil(t, err)

	// coins not positive
	toAddr0, err := sdk.AccAddressFromBech32("okexchain1dfpljpe0g0206jch32fx95lyagq3z5ws850m6f")
	require.Nil(t, err)
	decCoin0 := sdk.NewDecCoinFromDec(common.NativeToken, sdk.NewDec(0))
	transfers0 := []TransferUnit{
		{
			To:    toAddr0,
			Coins: sdk.DecCoins{decCoin0},
		},
	}

	// empty toAddr
	toAddr1, err := sdk.AccAddressFromBech32("")
	require.NoError(t, err)
	decCoin1 := sdk.NewDecCoinFromDec("obk", sdk.NewDec(100))
	transfers1 := []TransferUnit{
		{
			To:    toAddr1,
			Coins: sdk.DecCoins{decCoin1},
		},
	}

	testCase := []struct {
		multiSendMsg MsgMultiSend
		err          sdk.Error
	}{
		{NewMsgMultiSend(fromAddr, transfers), nil},
		{NewMsgMultiSend(sdk.AccAddress{}, transfers), sdk.ErrInvalidAddress(sdk.AccAddress{}.String())},
		{NewMsgMultiSend(fromAddr, make([]TransferUnit, MultiSendLimit+1)), sdk.ErrUnknownRequest("failed to check multisend msg because restrictions on the number of transfers")},
		{NewMsgMultiSend(fromAddr, transfers0), sdk.ErrInvalidCoins("failed to check multisend msg because send amount must be positive")},
		{NewMsgMultiSend(fromAddr, transfers1), sdk.ErrInvalidAddress("failed to check multisend msg because address is empty, not valid")},
	}
	for _, msgCase := range testCase {
		err := msgCase.multiSendMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	tokenMultiSendMsg := testCase[0].multiSendMsg
	signAddr := tokenMultiSendMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{fromAddr}, signAddr)

	bz := ModuleCdc.MustMarshalJSON(tokenMultiSendMsg)

	require.NoError(t, err)
	require.EqualValues(t, sdk.MustSortJSON(bz), tokenMultiSendMsg.GetSignBytes())
	require.EqualValues(t, "token", tokenMultiSendMsg.Route())
	require.EqualValues(t, "multi-send", tokenMultiSendMsg.Type())

	err = tokenMultiSendMsg.ValidateBasic()
	require.NoError(t, err)
}

func TestNewMsgTransferOwnership(t *testing.T) {
	// from
	fromPriKey := secp256k1.GenPrivKey()
	fromPubKey := fromPriKey.PubKey()
	fromAddr := sdk.AccAddress(fromPubKey.Address())

	// to
	toPriKey := secp256k1.GenPrivKey()
	toPubKey := toPriKey.PubKey()
	toAddr := sdk.AccAddress(toPubKey.Address())

	testCase := []struct {
		transferOwnershipMsg MsgTransferOwnership
		err                  sdk.Error
	}{
		{NewMsgTransferOwnership(fromAddr, toAddr, common.NativeToken), sdk.ErrUnauthorized("failed to check transferownership msg because invalid multi signature")},
		{NewMsgTransferOwnership(sdk.AccAddress{}, toAddr, common.NativeToken), sdk.ErrInvalidAddress("failed to check transferownership msg because miss sender address")},
		{NewMsgTransferOwnership(fromAddr, sdk.AccAddress{}, common.NativeToken), sdk.ErrInvalidAddress("failed to check transferownership msg because miss recipient address")},
		{NewMsgTransferOwnership(fromAddr, toAddr, ""), sdk.ErrUnknownRequest("failed to check transferownership msg because symbol cannot be empty")},
		{NewMsgTransferOwnership(fromAddr, toAddr, "1okb-ads"), sdk.ErrUnknownRequest("failed to check transferownership msg because invalid token symbol: 1okb-ads")},
	}
	for _, msgCase := range testCase {
		err := msgCase.transferOwnershipMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	transferOwnershipMsg := testCase[0].transferOwnershipMsg
	ret := transferOwnershipMsg.checkMultiSign()
	if !ret {
		//pubKey not set,test ValidateBasic
		bret := transferOwnershipMsg.ValidateBasic()
		require.NotNil(t, bret)
		//error pubKey
		transferOwnershipMsg.ToSignature.PubKey = fromPubKey
		if !transferOwnershipMsg.checkMultiSign() {
			//unequal with toPubKey address , try again
			transferOwnershipMsg.ToSignature.PubKey = toPubKey
			require.NotNil(t, transferOwnershipMsg.checkMultiSign())
			require.NotNil(t, transferOwnershipMsg.GetSignBytes())
		}
	}
	transferOwnershipMsg.Route()
	transferOwnershipMsg.Type()
	signAddr := transferOwnershipMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{fromAddr}, signAddr)
}

func TestNewMsgTokenModify(t *testing.T) {
	priKey := secp256k1.GenPrivKey()
	pubKey := priKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())

	testCase := []struct {
		tokenModifyMsg MsgTokenModify
		err            sdk.Error
	}{
		{NewMsgTokenModify("bnb", "bnb", "bnb bnb", true, true, addr),
			nil},
		{NewMsgTokenModify("", "bnb", "bnb bnb", true, true, addr),
			sdk.ErrUnknownRequest("failed to check modify msg because symbol cannot be empty")},
		{NewMsgTokenModify("bnb", "bnb", "bnb bnb", true, true, sdk.AccAddress{}),
			sdk.ErrInvalidAddress(sdk.AccAddress{}.String())},
		{NewMsgTokenModify("bnb", "bnb", "bnbbbbbbbbbb bnbbbbbbbbbbbbbbbbb", true, true, addr),
			sdk.ErrUnknownRequest("failed to check modify msg because invalid wholename")},
		{NewMsgTokenModify("bnb", "bnb", "bnbbbbbbbbbb bnbbbbbbbbbbbbbbbbb", true, false, addr),
			nil},
		{NewMsgTokenModify("bnb", `bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234`, "bnbbbbbbbbbb", true, false, addr),
			sdk.ErrUnknownRequest("failed to check modify msg because invalid desc")},
		{NewMsgTokenModify("bnb", `bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234
bnbbbbbbbbbbbnbbbbbbbbbbnbbbbbbbbbbbnbbbbbbbbb1234`, "bnbbbbbbbbbb", false, false, addr),
			nil},
	}
	for _, msgCase := range testCase {
		err := msgCase.tokenModifyMsg.ValidateBasic()
		require.EqualValues(t, msgCase.err, err)
	}

	// correct message
	tokenEditMsg := testCase[0].tokenModifyMsg
	signAddr := tokenEditMsg.GetSigners()
	require.EqualValues(t, []sdk.AccAddress{addr}, signAddr)

	bz := ModuleCdc.MustMarshalJSON(tokenEditMsg)
	require.EqualValues(t, sdk.MustSortJSON(bz), tokenEditMsg.GetSignBytes())
	require.EqualValues(t, "edit", tokenEditMsg.Type())
	require.EqualValues(t, "token", tokenEditMsg.Route())

	err := tokenEditMsg.ValidateBasic()
	require.NoError(t, err)
}
