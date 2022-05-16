package types

import (
	"fmt"

	"github.com/tendermint/go-amino"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

const (
	gen0 = 0x3b6a57b2
	gen1 = 0x26508e6d
	gen2 = 0x1ea119fa
	gen3 = 0x3d4233dd
	gen4 = 0x2a1462b3
)

var gen = []int{gen0, gen1, gen2, gen3, gen4}

func (aa AccAddress) Bech32StringOptimized(bech32PrefixAccAddr string) string {
	convertedLen := len(aa.Bytes())*8/5 + 1
	resultLen := len(bech32PrefixAccAddr) + 1 + convertedLen + 6
	var result = make([]byte, resultLen)

	copy(result, bech32PrefixAccAddr)
	result[len(bech32PrefixAccAddr)] = '1'

	var err error
	prefixLen := len(bech32PrefixAccAddr) + 1
	converted := result[prefixLen:prefixLen]
	converted, err = convertBitsTo(aa.Bytes(), 8, 5, true, converted)
	if err != nil {
		panic(fmt.Errorf("encoding bech32 failed %w", err))
	}
	if len(converted) > convertedLen {
		panic("returned unexpected length")
	}
	if len(converted) < convertedLen {
		result = result[0 : len(result)-(convertedLen-len(converted))]
	}

	bech32ChecksumTo(bech32PrefixAccAddr, converted, result[prefixLen+len(converted):])

	err = toChars(result[prefixLen:])
	if err != nil {
		panic(fmt.Errorf("unable to convert data bytes to chars: "+
			"%v", err))
	}
	return amino.BytesToStr(result)
}

// Encode encodes a byte slice into a bech32 string with the
// human-readable part hrb. Note that the bytes must each encode 5 bits
// (base32).
func encode(hrp string, data []byte) (string, error) {
	// Calculate the checksum of the data and append it at the end.
	resultLen := len(hrp) + 1 + len(data) + 6
	var result = make([]byte, resultLen)
	copy(result, hrp)
	result[len(hrp)] = '1'
	copy(result[len(hrp)+1:], data)

	bech32ChecksumTo(hrp, data, result[len(hrp)+1+len(data):])

	err := toChars(result[len(hrp)+1:])
	if err != nil {
		return "", fmt.Errorf("unable to convert data bytes to chars: "+
			"%v", err)
	}
	return amino.BytesToStr(result), nil
}

func convertBits(data []byte, fromBits, toBits uint8, pad bool) ([]byte, error) {
	var regrouped = make([]byte, 0, len(data)*int(fromBits)/int(toBits)+1)
	return convertBitsTo(data, fromBits, toBits, pad, regrouped)
}

func convertBitsTo(data []byte, fromBits, toBits uint8, pad bool, target []byte) ([]byte, error) {
	if fromBits < 1 || fromBits > 8 || toBits < 1 || toBits > 8 {
		return nil, fmt.Errorf("only bit groups between 1 and 8 allowed")
	}

	// The final bytes, each byte encoding toBits bits.
	var regrouped []byte
	if target != nil {
		regrouped = target[0:]
	} else {
		regrouped = make([]byte, 0, len(data)*int(fromBits)/int(toBits)+1)
	}

	// Keep track of the next byte we create and how many bits we have
	// added to it out of the toBits goal.
	nextByte := byte(0)
	filledBits := uint8(0)

	for _, b := range data {

		// Discard unused bits.
		b = b << (8 - fromBits)

		// How many bits remaining to extract from the input data.
		remFromBits := fromBits
		for remFromBits > 0 {
			// How many bits remaining to be added to the next byte.
			remToBits := toBits - filledBits

			// The number of bytes to next extract is the minimum of
			// remFromBits and remToBits.
			toExtract := remFromBits
			if remToBits < toExtract {
				toExtract = remToBits
			}

			// Add the next bits to nextByte, shifting the already
			// added bits to the left.
			nextByte = (nextByte << toExtract) | (b >> (8 - toExtract))

			// Discard the bits we just extracted and get ready for
			// next iteration.
			b = b << toExtract
			remFromBits -= toExtract
			filledBits += toExtract

			// If the nextByte is completely filled, we add it to
			// our regrouped bytes and start on the next byte.
			if filledBits == toBits {
				regrouped = append(regrouped, nextByte)
				filledBits = 0
				nextByte = 0
			}
		}
	}

	// We pad any unfinished group if specified.
	if pad && filledBits > 0 {
		nextByte = nextByte << (toBits - filledBits)
		regrouped = append(regrouped, nextByte)
		filledBits = 0
		nextByte = 0
	}

	// Any incomplete group must be <= 4 bits, and all zeroes.
	if filledBits > 0 && (filledBits > 4 || nextByte != 0) {
		return nil, fmt.Errorf("invalid incomplete group")
	}

	return regrouped, nil
}

func bech32Polymod(values []int) int {
	chk := 1
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ v
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

func bech32PolymodInternal(chk int, v int) int {
	b := chk >> 25
	chk = (chk&0x1ffffff)<<5 ^ v

	if (b>>uint(0))&1 == 1 {
		chk ^= gen0
	}
	if (b>>uint(1))&1 == 1 {
		chk ^= gen1
	}
	if (b>>uint(2))&1 == 1 {
		chk ^= gen2
	}
	if (b>>uint(3))&1 == 1 {
		chk ^= gen3
	}
	if (b>>uint(4))&1 == 1 {
		chk ^= gen4
	}

	return chk
}

// For more details on HRP expansion, please refer to BIP 173.
func bech32HrpExpand(hrp string) []int {
	v := make([]int, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]>>5))
	}
	v = append(v, 0)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]&31))
	}
	return v
}

func bech32ChecksumTo(hrp string, data []byte, result []byte) {
	chk := 1

	// hrpExpand
	for i := 0; i < len(hrp); i++ {
		chk = bech32PolymodInternal(chk, int(hrp[i]>>5))
	}
	chk = bech32PolymodInternal(chk, 0)
	for i := 0; i < len(hrp); i++ {
		chk = bech32PolymodInternal(chk, int(hrp[i]&31))
	}

	// Convert the bytes to list of integers, as this is needed for the
	// checksum calculation.
	for i := 0; i < len(data); i++ {
		chk = bech32PolymodInternal(chk, int(data[i]))
	}

	for i := 0; i < 6; i++ {
		chk = bech32PolymodInternal(chk, 0)
	}

	polymod := chk ^ 1
	if len(result) != 6 {
		panic("need more space")
	}
	for i := 0; i < 6; i++ {
		result[i] = byte((polymod >> uint(5*(5-i))) & 31)
	}
}

// toChars converts the byte slice 'data' to a string where each byte in 'data'
// encodes the index of a character in 'charset'.
func toChars(data []byte) error {
	for i := 0; i < len(data); i++ {
		if int(data[i]) >= len(charset) {
			return fmt.Errorf("invalid data byte: %v", data[i])
		}
		data[i] = charset[data[i]]
	}
	return nil
}
