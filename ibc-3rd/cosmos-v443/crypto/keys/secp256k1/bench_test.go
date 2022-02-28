package secp256k1

import (
	"io"
	"testing"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/keys/internal/benchmarking"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/crypto/types"
)

func BenchmarkKeyGeneration(b *testing.B) {
	b.ReportAllocs()
	benchmarkKeygenWrapper := func(reader io.Reader) types.PrivKey {
		priv := genPrivKey(reader)
		return &PrivKey{Key: priv}
	}
	benchmarking.BenchmarkKeyGeneration(b, benchmarkKeygenWrapper)
}

func BenchmarkSigning(b *testing.B) {
	b.ReportAllocs()
	priv := GenPrivKey()
	benchmarking.BenchmarkSigning(b, priv)
}

func BenchmarkVerification(b *testing.B) {
	b.ReportAllocs()
	priv := GenPrivKey()
	benchmarking.BenchmarkVerification(b, priv)
}
