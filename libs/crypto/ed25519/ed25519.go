package ed25519

import (
	stded25519 "crypto/ed25519"
	"strconv"
	"unsafe"
)

/*
#cgo CFLAGS: -I./libed25519/include
#cgo LDFLAGS: -led25519_okc
#include "ed25519_okc.h"
*/
import "C"

func cBufferToGoBuffer(cBuffer *C.Buffer) []byte {
	if cBuffer == nil || cBuffer.len == 0 {
		return nil
	}

	return C.GoBytes(unsafe.Pointer(cBuffer.data), C.int(cBuffer.len))
}

func goBufferToCBuffer(goBuffer []byte) C.Buffer {
	if len(goBuffer) == 0 {
		return C.Buffer{}
	}

	return C.Buffer{
		data: (*C.uchar)(&goBuffer[0]),
		len:  C.ulong(len(goBuffer)),
	}
}

// PublicKey is the type of Ed25519 public keys.
type PublicKey = stded25519.PublicKey

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey = stded25519.PrivateKey

// Sign signs the message with privateKey and returns a signature. It will
// panic if len(privateKey) is not PrivateKeySize.
func Sign(privateKey PrivateKey, message []byte) []byte {
	if l := len(privateKey); l != stded25519.PrivateKeySize {
		panic("ed25519: bad private key length: " + strconv.Itoa(l))
	}
	keypair := goBufferToCBuffer(privateKey)
	msg := goBufferToCBuffer(message)
	cSignature := C.okc_ed25519_sign(keypair, msg)
	signature := cBufferToGoBuffer(&cSignature)
	C.free_buf(cSignature)

	return signature
}

// Verify reports whether sig is a valid signature of message by publicKey. It
// will panic if len(publicKey) is not PublicKeySize.
func Verify(publicKey PublicKey, message, sig []byte) bool {
	if l := len(publicKey); l != stded25519.PublicKeySize {
		panic("ed25519: bad public key length: " + strconv.Itoa(l))
	}

	if len(sig) != stded25519.SignatureSize || sig[63]&224 != 0 {
		return false
	}
	cPublicKey := goBufferToCBuffer(publicKey)
	cMsg := goBufferToCBuffer(message)
	cSig := goBufferToCBuffer(sig)

	return bool(C.okc_ed25519_verify(cPublicKey, cMsg, cSig))
}
