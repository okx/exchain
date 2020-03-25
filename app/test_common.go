package app

import (
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// function for privateKey
func slice2Array(s []byte) (byteArray [32]byte, err error) {
	if len(s) != 32 {
		return byteArray, errors.New("byte slice's length is not 32")
	}
	for i := 0; i < 32; i++ {
		byteArray[i] = s[i]
	}
	return
}

func getPrivateKey(privateKeyStr string) secp256k1.PrivKeySecp256k1 {
	derivedPrivSlice, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		panic(err)
	}
	derivedPriv, err := slice2Array(derivedPrivSlice)
	if err != nil {
		panic(err)
	}
	return secp256k1.PrivKeySecp256k1(derivedPriv)
}
