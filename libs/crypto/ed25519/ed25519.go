package ed25519

import (
	stded25519 "crypto/ed25519"
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

// NewKeyFromSeed calculates a private key from a seed. It will panic if
// len(seed) is not SeedSize. This function is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this
// package.
func NewKeyFromSeed(seed []byte) PrivateKey {
	// Outline the function body so that the returned key can be stack-allocated.
	var keypair C.Buffer
	keypair = C.okc_ed25519_gen_keypair()

	buffer := cBufferToGoBuffer(&keypair)
	C.free_buf(keypair)

	return buffer[:]
}

// panic if len(privateKey) is not PrivateKeySize.
func Sign(privateKey PrivateKey, message []byte) []byte {
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
	cPublicKey := goBufferToCBuffer(publicKey)
	cMsg := goBufferToCBuffer(message)
	cSig := goBufferToCBuffer(sig)

	return bool(C.okc_ed25519_verify(cPublicKey, cMsg, cSig))
}
