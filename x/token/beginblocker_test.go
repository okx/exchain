package token

import (
	"fmt"
	"strings"
	"testing"
)

func TestBeginBlocker(t *testing.T) {
	s := "ammswap_ltck-5cb_okt"
	items := strings.SplitN(strings.SplitN(s, "ammswap_", 2)[1], "_", 2)
	fmt.Println(items)
}
