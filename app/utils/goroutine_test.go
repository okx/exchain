package utils

import (
	"testing"
)

func TestGoroutine(t *testing.T) {

	t.Log("A go routine printed as dec:", GoRId)

	var gorhex GoRoutineID = 16

	t.Log("A go routine printed as hex:", gorhex)

}