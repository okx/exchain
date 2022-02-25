package db

import (
	"bytes"
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

var invalidByteQuantityError = errors.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")

// We defensively turn nil keys or values into []byte{} for
// most operations.
func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	}
	return bz
}

func cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}

func concat(bz1, bz2 []byte) (ret []byte) {
	ret = make([]byte, len(bz1)+len(bz2))
	copy(ret, bz1)
	copy(ret[len(bz1):], bz2)
	return ret
}

func concatAll(bzs ...[]byte) (ret []byte) {
	totalLen := 0
	for _, bz := range bzs {
		totalLen += len(bz)
	}
	ret = make([]byte, totalLen)
	cpLen := 0
	for _, bz := range bzs {
		copy(ret[cpLen:], bz)
		cpLen += len(bz)
	}
	return ret
}

// Returns a slice of the same length (big endian)
// except incremented by one.
// Returns nil on overflow (e.g. if bz bytes are all 0xFF)
// CONTRACT: len(bz) > 0
func cpIncr(bz []byte) (ret []byte) {
	if len(bz) == 0 {
		panic("cpIncr expects non-zero bz length")
	}
	ret = cp(bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] < byte(0xFF) {
			ret[i]++
			return
		}
		ret[i] = byte(0x00)
		if i == 0 {
			// Overflow
			return nil
		}
	}
	return nil
}

// See DB interface documentation for more information.
func IsKeyInDomain(key, start, end []byte) bool {
	if bytes.Compare(key, start) < 0 {
		return false
	}
	if end != nil && bytes.Compare(end, key) <= 0 {
		return false
	}
	return true
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// toBytes parses a string formatted by ByteSize as bytes.
// KB = K = KiB	= 1024
// MB = M = MiB = 1024 * K
// GB = G = GiB = 1024 * M
// TB = T = TiB = 1024 * G
// PB = P = PiB = 1024 * T
// EB = E = EiB = 1024 * P
func toBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)
	if i == -1 {
		return 0, invalidByteQuantityError
	}

	bytesString, multiple := s[:i], s[i:]
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes < 0 {
		return 0, invalidByteQuantityError
	}

	switch multiple {
	case "E", "EB", "EIB":
		return uint64(bytes * EXABYTE), nil
	case "P", "PB", "PIB":
		return uint64(bytes * PETABYTE), nil
	case "T", "TB", "TIB":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GIB":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MIB":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KIB":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, invalidByteQuantityError
	}
}

func parseOptParams(params string) map[string]string {
	if len(params) == 0 {
		return nil
	}

	opts := make(map[string]string)
	for _, s := range strings.Split(params, ",") {
		opt := strings.Split(s, "=")
		if len(opt) != 2 {
			panic("Invalid options parameter, like this 'block_size=4kb,statistics=true")
		}
		opts[strings.TrimSpace(opt[0])] = strings.TrimSpace(opt[1])
	}
	return opts
}
