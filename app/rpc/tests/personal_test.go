package tests

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func (suite *RPCTestSuite) TestPersonal_NewAccount() {
	// create an new mnemonics randomly on the node
	rpcRes := Call(suite.T(), suite.addr, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &addr))
	addrCounter++

	rpcRes = Call(suite.T(), suite.addr, "personal_listAccounts", nil)
	var res []common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &res))
	suite.Require().Equal(1, len(res))
	suite.Require().True(res[0] == addr)
}

func (suite *RPCTestSuite) TestPersonal_Sign() {
	rpcRes := Call(suite.T(), suite.addr, "personal_sign", []interface{}{hexutil.Bytes{0x88}, hexutil.Bytes(senderAddr[:]), ""})

	var res hexutil.Bytes
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &res))
	suite.Require().Equal(65, len(res))
	// TODO: check that signature is same as with geth, requires importing a key

	// error with inexistent addr
	inexistentAddr := common.BytesToAddress([]byte{0})
	_, err := CallWithError(suite.addr, "personal_sign", []interface{}{hexutil.Bytes{0x88}, inexistentAddr, ""})
	suite.Require().Error(err)
}

func (suite *RPCTestSuite) TestPersonal_ImportRawKey() {
	privkey, err := ethcrypto.GenerateKey()
	suite.Require().NoError(err)

	// parse priv key to hex
	hexPriv := common.Bytes2Hex(ethcrypto.FromECDSA(privkey))
	rpcRes := Call(suite.T(), suite.addr, "personal_importRawKey", []string{hexPriv, defaultPassWd})

	var resAddr common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &resAddr))

	addr := ethcrypto.PubkeyToAddress(privkey.PublicKey)

	suite.Require().True(addr == resAddr)

	addrCounter++

	// error check with wrong hex format of privkey
	rpcRes, err = CallWithError(suite.addr, "personal_importRawKey", []string{fmt.Sprintf("%sg", hexPriv), defaultPassWd})
	suite.Require().Error(err)
}

func (suite *RPCTestSuite) TestPersonal_ImportRawKey_Duplicate() {
	privkey, err := ethcrypto.GenerateKey()
	suite.Require().NoError(err)
	// parse priv key to hex, then add the key
	hexPriv := common.Bytes2Hex(ethcrypto.FromECDSA(privkey))
	rpcRes := Call(suite.T(), suite.addr, "personal_importRawKey", []string{hexPriv, defaultPassWd})
	var resAddr common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &resAddr))
	suite.Require().True(ethcrypto.PubkeyToAddress(privkey.PublicKey) == resAddr)
	addrCounter++

	// record the key-list length
	rpcRes = Call(suite.T(), suite.addr, "personal_listAccounts", nil)
	var list []common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &list))
	originLen := len(list)

	// add the same key again
	rpcRes = Call(suite.T(), suite.addr, "personal_importRawKey", []string{hexPriv, defaultPassWd})
	var newResAddr common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &newResAddr))
	suite.Require().Equal(resAddr, newResAddr)

	// check the actual key-list length changed or not
	rpcRes = Call(suite.T(), suite.addr, "personal_listAccounts", nil)
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &list))
	suite.Require().Equal(originLen, len(list))
}

func (suite *RPCTestSuite) TestPersonal_EcRecover() {
	data := hexutil.Bytes{0x88}
	rpcRes := Call(suite.T(), suite.addr, "personal_sign", []interface{}{data, hexutil.Bytes(senderAddr[:]), ""})

	var res hexutil.Bytes
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &res))
	suite.Require().Equal(65, len(res))

	rpcRes = Call(suite.T(), suite.addr, "personal_ecRecover", []interface{}{data, res})
	var ecrecoverRes common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &ecrecoverRes))
	suite.Require().Equal(senderAddr.Bytes(), ecrecoverRes[:])

	// error check for ecRecover
	// wrong length of sig
	rpcRes, err := CallWithError(suite.addr, "personal_ecRecover", []interface{}{data, res[1:]})
	suite.Require().Error(err)

	// wrong RecoveryIDOffset -> nether 27 nor 28
	res[ethcrypto.RecoveryIDOffset] = 29
	rpcRes, err = CallWithError(suite.addr, "personal_ecRecover", []interface{}{data, res})
	suite.Require().Error(err)

	// fail in SigToPub
	sigInvalid := make(hexutil.Bytes, 65)
	for i := 0; i < 64; i++ {
		sigInvalid[i] = 0
	}
	sigInvalid[64] = 27
	rpcRes, err = CallWithError(suite.addr, "personal_ecRecover", []interface{}{data, sigInvalid})
	suite.Require().Error(err)
}

func (suite *RPCTestSuite) TestPersonal_UnlockAccount() {
	// create a new account
	rpcRes := Call(suite.T(), suite.addr, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &addr))

	addrCounter++

	newPassWd := "87654321"
	// try to sign with different password -> failed
	_, err := CallWithError(suite.addr, "personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, newPassWd})
	suite.Require().Error(err)

	// unlock the address with the new password
	rpcRes = Call(suite.T(), suite.addr, "personal_unlockAccount", []interface{}{addr, newPassWd})
	var unlocked bool
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &unlocked))
	suite.Require().True(unlocked)

	// try to sign with the new password -> successfully
	rpcRes, err = CallWithError(suite.addr, "personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, newPassWd})
	suite.Require().NoError(err)
	var res hexutil.Bytes
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &res))
	suite.Require().Equal(65, len(res))

	// error check
	// inexistent addr
	inexistentAddr := common.BytesToAddress([]byte{0})
	_, err = CallWithError(suite.addr, "personal_unlockAccount", []interface{}{hexutil.Bytes{0x88}, inexistentAddr, newPassWd})
	suite.Require().Error(err)
}

func (suite *RPCTestSuite) TestPersonal_LockAccount() {
	// create a new account
	rpcRes := Call(suite.T(), suite.addr, "personal_newAccount", []string{defaultPassWd})
	var addr common.Address
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &addr))

	addrCounter++

	// unlock the account above first
	rpcRes = Call(suite.T(), suite.addr, "personal_unlockAccount", []interface{}{addr, defaultPassWd})
	var unlocked bool
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &unlocked))
	suite.Require().True(unlocked)

	// lock the account
	rpcRes = Call(suite.T(), suite.addr, "personal_lockAccount", []interface{}{addr})
	var locked bool
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &locked))
	suite.Require().True(locked)

	// try to sign, should be locked -> fail to sign
	_, err := CallWithError(suite.addr, "personal_sign", []interface{}{hexutil.Bytes{0x88}, addr, defaultPassWd})
	suite.Require().Error(err)

	// error check
	// lock an inexistent account
	inexistentAddr := common.BytesToAddress([]byte{0})
	rpcRes = Call(suite.T(), suite.addr, "personal_lockAccount", []interface{}{inexistentAddr})
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &locked))
	suite.Require().False(locked)
}

func (suite *RPCTestSuite) TestPersonal_SendTransaction_Transfer() {
	params := make([]interface{}, 2)
	params[0] = map[string]string{
		"from":  senderAddr.Hex(),
		"to":    receiverAddr.Hex(),
		"value": "0x16345785d8a0000", // 0.1
	}
	params[1] = defaultPassWd

	rpcRes := Call(suite.T(), suite.addr, "personal_sendTransaction", params)
	var hash ethcmn.Hash
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &hash))

	commitBlock(suite)
	commitBlock(suite)

	receipt := WaitForReceipt(suite.T(), suite.addr, hash)
	suite.Require().NotNil(receipt)
	suite.Require().Equal("0x1", receipt["status"].(string))
}

func (suite *RPCTestSuite) TestPersonal_SendTransaction_DeployContract() {
	params := make([]interface{}, 2)
	params[0] = map[string]string{
		"from":     senderAddr.Hex(),
		"data":     "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029",
		"gasPrice": (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String(),
	}
	params[1] = defaultPassWd

	rpcRes := Call(suite.T(), suite.addr, "personal_sendTransaction", params)
	var hash ethcmn.Hash
	suite.Require().NoError(json.Unmarshal(rpcRes.Result, &hash))

	commitBlock(suite)
	commitBlock(suite)

	receipt := WaitForReceipt(suite.T(), suite.addr, hash)
	suite.Require().NotNil(receipt)
	suite.Require().Equal("0x1", receipt["status"].(string))
}
