package types

import (
	"math/big"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func BenchmarkKeccak256HashCache(b *testing.B) {
	b.ResetTimer()
	b.Run("without cache", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = ethcrypto.Keccak256Hash(hash[:])
		}
	})
	b.Run("lru set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = keccak256HashWithLruCache(hash[:])
		}
	})

	b.Run("fastcache set", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = keccak256HashWithFastCache(hash[:])
		}
	})

	const getCount = 1

	fastcacheGet := func() {
		for i := 0; i < getCount; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = keccak256HashWithFastCache(hash[:])
		}
	}

	lruGet := func() {
		for i := 0; i < getCount; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = keccak256HashWithLruCache(hash[:])
		}
	}

	withoutCacheGet := func() {
		for i := 0; i < getCount; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = ethcrypto.Keccak256Hash(hash[:])
		}
	}

	b.ResetTimer()
	b.Run("lru get", func(b *testing.B) {
		lruGet()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			lruGet()
		}
	})
	b.Run("fastcache get", func(b *testing.B) {
		fastcacheGet()
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			fastcacheGet()
		}
	})
	b.Run("withoutcache get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			withoutCacheGet()
		}
	})
}

func TestKeccak256HashWithSyncPool(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		hash := ethcmn.BigToHash(big.NewInt(int64(i)))
		actual := keccak256HashWithSyncPool(hash[:])
		expect := ethcrypto.Keccak256Hash(hash[:])
		if actual != expect {
			t.Errorf("expect %v, actual %v", expect, actual)
		}
	}
}

func BenchmarkKeccak256HashNew(b *testing.B) {
	b.ResetTimer()
	b.Run("new", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = ethcrypto.Keccak256Hash(hash[:])
		}
	})
	b.Run("reuse keccak", func(b *testing.B) {
		b.ReportAllocs()
		d := ethcrypto.NewKeccakState()
		for i := 0; i < b.N; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			d.Write(hash[:])
			var h ethcmn.Hash
			d.Read(h[:])
			d.Reset()
		}
	})
	b.Run("use sync Pool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			hash := ethcmn.BigToHash(big.NewInt(int64(i)))
			_ = keccak256HashWithSyncPool(hash[:])
		}
	})
}
