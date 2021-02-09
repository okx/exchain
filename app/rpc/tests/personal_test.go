package tests

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPersonal_ListAccounts(t *testing.T) {
	// there are two keys to unlock in the node from test.sh
	rpcRes := Call(t, "personal_listAccounts", nil)

	var res []common.Address
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 2, len(res))
	require.True(t, res[0] == hexAddr1)
	require.True(t, res[1] == hexAddr2)
}

func TestPersonal_NewAccount(t *testing.T) {
	// create an new mnemonics randomly on the node
	rpcRes := Call(t, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &addr))
	addrCounter++

	rpcRes = Call(t, "personal_listAccounts", nil)
	var res []common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.Equal(t, 3, len(res))
	require.True(t, res[0] == hexAddr1)
	require.True(t, res[1] == hexAddr2)
	require.True(t, res[2] == addr)
}

func TestPersonal_Sign(t *testing.T) {
	rpcRes := Call(t, "personal_sign", []interface{}{hexutil.Bytes{0x88}, hexutil.Bytes(from), ""})

	var res hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.Equal(t, 65, len(res))
	// TODO: check that signature is same as with geth, requires importing a key

	// error with inexistent addr
	inexistentAddr := common.BytesToAddress([]byte{0})
	rpcRes, err := CallWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, inexistentAddr, ""})
	require.Error(t, err)
}

func TestPersonal_ImportRawKey(t *testing.T) {
	privkey, err := ethcrypto.GenerateKey()
	require.NoError(t, err)

	// parse priv key to hex
	hexPriv := common.Bytes2Hex(ethcrypto.FromECDSA(privkey))
	rpcRes := Call(t, "personal_importRawKey", []string{hexPriv, defaultPassWd})

	var resAddr common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &resAddr))

	addr := ethcrypto.PubkeyToAddress(privkey.PublicKey)

	require.True(t, addr == resAddr)

	addrCounter++

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

	addrCounter++

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
	inexistentAddr := common.BytesToAddress([]byte{0})
	_, err = CallWithError("personal_unlockAccount", []interface{}{hexutil.Bytes{0x88}, inexistentAddr, newPassWd})
	require.Error(t, err)
}

func TestPersonal_LockAccount(t *testing.T) {
	// create a new account
	rpcRes := Call(t, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &addr))

	addrCounter++

	// unlock the account above first
	rpcRes = Call(t, "personal_unlockAccount", []interface{}{addr, defaultPassWd})
	var unlocked bool
	require.NoError(t, json.Unmarshal(rpcRes.Result, &unlocked))
	require.True(t, unlocked)

	// lock the account
	rpcRes = Call(t, "personal_lockAccount", []interface{}{addr})
	var locked bool
	require.NoError(t, json.Unmarshal(rpcRes.Result, &locked))
	require.True(t, locked)

	// try to sign, should be locked -> fail to sign
	_, err := CallWithError("personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, defaultPassWd})
	require.Error(t, err)

	// error check
	// lock an inexistent account
	inexistentAddr := common.BytesToAddress([]byte{0})
	rpcRes = Call(t, "personal_lockAccount", []interface{}{inexistentAddr})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &locked))
	require.False(t, locked)
}
