package ed25519_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"
)

func TestSignAndValidateEd25519(t *testing.T) {

	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()

	msg := crypto.CRandBytes(128)
	sig, err := privKey.Sign(msg)
	require.Nil(t, err)

	// Test the signature
	assert.True(t, pubKey.VerifyBytes(msg, sig))

	// Mutate the signature, just one bit.
	// TODO: Replace this with a much better fuzzer, tendermint/ed25519/issues/10
	sig[7] ^= byte(0x01)

	assert.False(t, pubKey.VerifyBytes(msg, sig))

	bytes := pubKey.Bytes()
	var recoverPubKey ed25519.PubKeyEd25519
	recoverPubKey.UnmarshalFromAmino(bytes)
	assert.Equal(t, bytes, recoverPubKey.Bytes())
}

func TestWTx2(t *testing.T) {

}

func genPrivkey(hex string) ed25519.PrivKeyEd25519 {
	secert, err := hexutil.Decode(hex)
	if err != nil {
		panic(err)
	}
	return ed25519.GenPrivKeyFromSecret(secert)
}

func genPubkey(hex string) ed25519.PubKeyEd25519 {
	bytes := hexutil.MustDecode(hex)
	var recoverPubKey ed25519.PubKeyEd25519
	recoverPubKey.UnmarshalFromAmino(bytes)
	return recoverPubKey
}

func TestWTx(t *testing.T) {
	payload := []byte("payload")

	PrivKey := "0xa3288910406c97c4b1998e00171295f09681a385e898e8bf172c2befe181fdffa549b6b68568de211f86405eabee64b893e2938e753309ba294d8914887525c598c279bfd9"
	PubKey := "0x1624de642068de211f86405eabee64b893e2938e753309ba294d8914887525c598c279bfd9"
	//
	PrivKey= "0xa3288910402de16907e788ccb9f3ed48ad6cca3198dd92334dd710b89ec19988b8d48d5f0fd134f5e36c5fdcf28ebe3b7ae039ace09d0198513f7d03500a2b4dc0465aff31"
	PubKey= "0x1624de6420d134f5e36c5fdcf28ebe3b7ae039ace09d0198513f7d03500a2b4dc0465aff31"

	priv := genPrivkey(t, PrivKey)
	fmt.Printf("%s\n", 	hexutil.Encode(priv.PubKey().Bytes()))
	fmt.Printf("%s\n", 	PubKey)




	signature, err := priv.Sign(payload)
	require.NoError(t, err)





	res := priv.PubKey().VerifyBytes(payload, signature)
	require.Equal(t, res, true)

	bytes := priv.PubKey().Bytes()
	var recoverPubKey ed25519.PubKeyEd25519
	recoverPubKey.UnmarshalFromAmino(bytes)
	assert.Equal(t, bytes, recoverPubKey.Bytes())


	fmt.Printf("%s\n", 	hexutil.Encode(bytes))
	fmt.Printf("%s\n", 	hexutil.Encode(priv.PubKey().Bytes()))
	pubKeyBytes, err := hexutil.Decode(hexutil.Encode(bytes))
	require.NoError(t, err)

	assert.Equal(t, pubKeyBytes, recoverPubKey.Bytes())
}
