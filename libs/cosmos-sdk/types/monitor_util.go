package types

import "strconv"

func ConvertDecToFloat64(dec Dec) float64 {
	if dec.IsZero() {
		return 0.0
	}
	decStr := dec.String()
	f, err := strconv.ParseFloat(decStr, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func StringsContains(array []string, val string) int {
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			return i
		}
	}
	return -1
}
