package ed25519_test

import (
	stded25519 "crypto/ed25519"
	"testing"

	"github.com/okex/exchain/libs/crypto/ed25519"
)

func BenchmarkStdEd25519NewKey(b *testing.B) {
	seed := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		stded25519.NewKeyFromSeed(seed)
	}
}

func BenchmarkNewKey(b *testing.B) {
	seed := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		ed25519.NewKeyFromSeed(seed)
	}
}

func BenchmarkStdEd25519Sign(b *testing.B) {
	seed := make([]byte, 32)
	prikey := stded25519.NewKeyFromSeed(seed)
	message := []byte("this is a sign test string to benchmark.")
	for i := 0; i < b.N; i++ {
		stded25519.Sign(prikey, message)
	}
}

func BenchmarkSign(b *testing.B) {
	priKey := ed25519.NewKeyFromSeed(nil)
	message := []byte("this is a sign test string to benchmark.")
	for i := 0; i < b.N; i++ {
		ed25519.Sign(priKey, message)
	}
}

func BenchmarkStdEd25519Verify(b *testing.B) {
	seed := make([]byte, 32)
	priKey := stded25519.NewKeyFromSeed(seed)
	message := []byte("this is a sign test string to benchmark.")
	signature := stded25519.Sign(priKey, message)
	pubKey := priKey.Public().(stded25519.PublicKey)
	for i := 0; i < b.N; i++ {
		stded25519.Verify(pubKey, message, signature)
	}
}

func BenchmarkVerify(b *testing.B) {
	priKey := ed25519.NewKeyFromSeed(nil)
	message := []byte("this is a sign test string to benchmark.")
	signature := ed25519.Sign(priKey, message)
	pubKey := priKey[32:]
	for i := 0; i < b.N; i++ {
		ed25519.Verify(ed25519.PublicKey(pubKey), message, signature)
	}
}
