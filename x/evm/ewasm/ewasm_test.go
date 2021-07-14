package ewasm

import "testing"

func TestInitEWASM(t *testing.T) {
	w := NewEWasm()
	w.Version()
}
