package token

import "testing"

func TestBeginBlocker(t *testing.T) {
	ctx, kpr, _, _ := CreateParam(t, false)

	beginBlocker(ctx, kpr)
}
