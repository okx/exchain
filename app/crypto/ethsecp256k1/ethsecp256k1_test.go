package ethsecp256k1

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethsecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"

	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
)

func TestPrivKeyPrivKey(t *testing.T) {
	// validate type and equality
	privKey, err := GenerateKey()
	require.NoError(t, err)
	require.True(t, privKey.Equals(privKey))
	require.Implements(t, (*tmcrypto.PrivKey)(nil), privKey)

	// validate inequality
	privKey2, err := GenerateKey()
	require.NoError(t, err)
	require.False(t, privKey.Equals(privKey2))

	// validate Ethereum address equality
	addr := privKey.PubKey().Address()
	expectedAddr := ethcrypto.PubkeyToAddress(privKey.ToECDSA().PublicKey)
	require.Equal(t, expectedAddr.Bytes(), addr.Bytes())

	// validate we can sign some bytes
	msg := []byte("hello world")
	sigHash := ethcrypto.Keccak256Hash(msg)
	expectedSig, _ := ethsecp256k1.Sign(sigHash.Bytes(), privKey)

	sig, err := privKey.Sign(msg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, sig)
}

func TestPrivKeyPubKey(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)

	// validate type and equality
	pubKey := privKey.PubKey().(PubKey)
	require.Implements(t, (*tmcrypto.PubKey)(nil), pubKey)

	// validate inequality
	privKey2, err := GenerateKey()
	require.NoError(t, err)
	require.False(t, pubKey.Equals(privKey2.PubKey()))

	// validate signature
	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	require.NoError(t, err)

	res := pubKey.VerifyBytes(msg, sig)
	require.True(t, res)
}

func TestSignatureRecoverPrivateKey(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)

	// validate type and equality
	pubKey := privKey.PubKey().(PubKey)
	require.Implements(t, (*tmcrypto.PubKey)(nil), pubKey)

	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	require.NoError(t, err)
	r, s, v := decodeSignature(sig)

	t.Log("r1", r.String(), r.Int64())
	t.Log("s1", s.String())
	t.Log("v1", v.String())

	msg = []byte("lifei")
	sig, err = privKey.Sign(msg)
	require.NoError(t, err)
	r, s, v = decodeSignature(sig)
	t.Log("r2", r.String(), r.Int64())
	t.Log("s2", s.String())
	t.Log("v2", v.String())
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != ethcrypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), ethcrypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}
