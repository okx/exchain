package ed25519_test

import (
	"testing"

	stded25519 "golang.org/x/crypto/ed25519"
	// stded25519 "crypto/ed25519"

	"github.com/stretchr/testify/suite"

	"github.com/okex/exchain/libs/crypto/ed25519"
)

func BenchmarkStdEd25519Sign(b *testing.B) {
	seed := make([]byte, stded25519.SeedSize)
	prikey := stded25519.NewKeyFromSeed(seed)
	message := []byte("this is a sign test string to benchmark.")
	for i := 0; i < b.N; i++ {
		stded25519.Sign(prikey, message)
	}
}

func BenchmarkSign(b *testing.B) {
	seed := make([]byte, stded25519.SeedSize)
	priKey := stded25519.NewKeyFromSeed(seed)
	message := []byte("this is a sign test string to benchmark.")
	for i := 0; i < b.N; i++ {
		ed25519.Sign(priKey, message)
	}
}

func BenchmarkStdEd25519Verify(b *testing.B) {
	seed := make([]byte, stded25519.SeedSize)
	priKey := stded25519.NewKeyFromSeed(seed)
	message := []byte("this is a sign test string to benchmark.")
	signature := stded25519.Sign(priKey, message)
	pubKey := priKey.Public().(stded25519.PublicKey)
	for i := 0; i < b.N; i++ {
		stded25519.Verify(pubKey, message, signature)
	}
}

func BenchmarkVerify(b *testing.B) {
	seed := make([]byte, stded25519.SeedSize)
	priKey := stded25519.NewKeyFromSeed(seed)
	message := []byte("this is a sign test string to benchmark.")
	signature := ed25519.Sign(priKey, message)
	pubKey := priKey[32:]
	for i := 0; i < b.N; i++ {
		ed25519.Verify(ed25519.PublicKey(pubKey), message, signature)
	}
}

type Ed25519CompatibleSuite struct {
	suite.Suite
}

func TestEd25519CompatibleSuite(t *testing.T) {
	suite.Run(t, new(Ed25519CompatibleSuite))
}

func (suite *Ed25519CompatibleSuite) TestCompatible() {
	seed := make([]byte, stded25519.SeedSize)
	privateKey := stded25519.NewKeyFromSeed(seed)
	testCases := []struct {
		name       string
		privateKey stded25519.PrivateKey
		rawMsg     []byte
		signFunc   func(key stded25519.PrivateKey, msg []byte) []byte
		verifyFunc func(key stded25519.PublicKey, msg, sign []byte) bool
	}{
		{"1.go sign, rust verify", privateKey, []byte("hello world from go"), stded25519.Sign, ed25519.Verify},
		{"2.rust sign, go verify", privateKey, []byte("hello world from rust"), ed25519.Sign, stded25519.Verify},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			signature := tc.signFunc(tc.privateKey, tc.rawMsg)
			suite.Require().NotNil(signature)
			pk := tc.privateKey.Public().(stded25519.PublicKey)
			verified := tc.verifyFunc(pk, tc.rawMsg, signature)
			suite.Require().Equal(true, verified)
		})
	}
}
