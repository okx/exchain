package dydx

import "fmt"

var (
	ErrInvalidSignedOrder = fmt.Errorf("invalid signed order")
	ErrInvalidOrder       = fmt.Errorf("invalid order")
	ErrExpiredOrder       = fmt.Errorf("expired order")
	ErrInvalidSignature   = fmt.Errorf("invalid signature")
)
