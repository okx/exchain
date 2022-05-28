package ed25519

/*
#cgo CFLAGS: -I./libed25519/include
#cgo LDFLAGS: -led25519_okc
#include "okc_ed25519.h"
*/
import "C"
import (
	"unsafe"
)

// const (
// 	// PublicKeySize is the size, in bytes, of public keys as used in this package.
// 	PublicKeySize = 32
// 	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
// 	PrivateKeySize = 64
// 	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
// 	SignatureSize = 64
// 	// SeedSize is the size, in bytes, of private key seeds. These are the private key representations used by RFC 8032.
// 	SeedSize = 32
// )

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
type PublicKey []byte

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey []byte

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

// Sign signs the message with privateKey and returns a signature. It will
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
