package tests

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	hexAddr3, hexAddr4, hexAddr5 string
)

func TestPersonal_ListAccounts(t *testing.T) {
	// there are two keys to unlock in the node from test.sh
	rpcRes := Call(t, "personal_listAccounts", []string{})

	var res []hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	require.True(t, strings.EqualFold(hexutil.Encode(res[0]), hexAddr1))
	require.True(t, strings.EqualFold(hexutil.Encode(res[1]), hexAddr2))
}

func TestPersonal_NewAccount(t *testing.T) {
	// create an new mnemonics randomly on the node
	rpcRes := Call(t, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &addr))
	// global stores
	hexAddr3 = addr.Hex()

	rpcRes = Call(t, "personal_listAccounts", []string{})
	var res []hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.Equal(t, 3, len(res))
	require.True(t, strings.EqualFold(hexutil.Encode(res[0]), hexAddr1))
	require.True(t, strings.EqualFold(hexutil.Encode(res[1]), hexAddr2))
	require.True(t, strings.EqualFold(hexutil.Encode(res[2]), hexAddr3))
}

func TestPersonal_Sign(t *testing.T) {
	rpcRes := Call(t, "personal_sign", []interface{}{hexutil.Bytes{0x88}, hexutil.Bytes(from), ""})

	var res hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.Equal(t, 65, len(res))
	// TODO: check that signature is same as with geth, requires importing a key

	// error with inexistent addr
	inexistentAddr := hexutil.Bytes([]byte{0})
	rpcRes, err := CallWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, inexistentAddr, ""})
	require.Error(t, err)
}

func TestPersonal_ImportRawKey(t *testing.T) {
	privkey, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	// parse priv key to hex
	hexPriv := common.Bytes2Hex(ethcrypto.FromECDSA(privkey))
	rpcRes := Call(t, "personal_importRawKey", []string{hexPriv, defaultPassWd})

	var res hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))

	addr := ethcrypto.PubkeyToAddress(privkey.PublicKey)
	resAddr := common.BytesToAddress(res)

	require.Equal(t, addr.String(), resAddr.String())

	// global stores
	hexAddr4 = resAddr.String()

	// error check with wrong hex format of privkey
	rpcRes, err = CallWithError("personal_importRawKey", []string{fmt.Sprintf("%sg", hexPriv), defaultPassWd})
	require.Error(t, err)
}

func TestPersonal_EcRecover(t *testing.T) {
	data := hexutil.Bytes{0x88}
	rpcRes := Call(t, "personal_sign", []interface{}{data, hexutil.Bytes(from), ""})

	var res hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.Equal(t, 65, len(res))

	rpcRes = Call(t, "personal_ecRecover", []interface{}{data, res})
	var ecrecoverRes common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ecrecoverRes))
	require.Equal(t, from, ecrecoverRes[:])

	// error check for ecRecover
	// wrong length of sig
	rpcRes, err := CallWithError("personal_ecRecover", []interface{}{data, res[1:]})
	require.Error(t, err)

	// wrong RecoveryIDOffset -> nether 27 nor 28
	res[ethcrypto.RecoveryIDOffset] = 29
	rpcRes, err = CallWithError("personal_ecRecover", []interface{}{data, res})
	require.Error(t, err)

	// fail in SigToPub
	sigInvalid := make(hexutil.Bytes, 65)
	for i := 0; i < 64; i++ {
		sigInvalid[i] = 0
	}
	sigInvalid[64] = 27
	rpcRes, err = CallWithError("personal_ecRecover", []interface{}{data, sigInvalid})
	require.Error(t, err)
}

func TestPersonal_UnlockAccount(t *testing.T) {
	// create a new account
	rpcRes := Call(t, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &addr))

	// global stores
	hexAddr5 = addr.Hex()

	newPassWd := "87654321"
	// try to sign with different password -> failed
	_, err := CallWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, newPassWd})
	require.Error(t, err)

	// unlock the address with the new password
	rpcRes = Call(t, "personal_unlockAccount", []interface{}{addr, newPassWd})
	var unlocked bool
	require.NoError(t, json.Unmarshal(rpcRes.Result, &unlocked))
	require.True(t, unlocked)

	// try to sign with the new password -> successfully
	rpcRes, err = CallWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, newPassWd})
	require.NoError(t, err)
	var res hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.Equal(t, 65, len(res))

	// error check
	// inexistent addr
	inexistentAddr := hexutil.Bytes([]byte{0})
	_, err = CallWithError("personal_unlockAccount", []interface{}{hexutil.Bytes{0x88}, inexistentAddr, newPassWd})
	require.Error(t, err)
}

func TestPersonal_LockAccount(t *testing.T) {
	pswd := "nootwashere"
	rpcRes := Call(t, "personal_newAccount", []string{pswd})
	var addr common.Address
	err := json.Unmarshal(rpcRes.Result, &addr)
	require.NoError(t, err)

	rpcRes = Call(t, "personal_unlockAccount", []interface{}{addr, ""})
	var unlocked bool
	err = json.Unmarshal(rpcRes.Result, &unlocked)
	require.NoError(t, err)
	require.True(t, unlocked)

	rpcRes = Call(t, "personal_lockAccount", []interface{}{addr})
	var locked bool
	err = json.Unmarshal(rpcRes.Result, &locked)
	require.NoError(t, err)
	require.True(t, locked)

	// try to sign, should be locked
	_, err = CallWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, ""})
	require.Error(t, err)
}
