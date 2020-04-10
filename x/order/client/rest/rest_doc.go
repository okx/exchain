// Package rest API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
// swagger:meta
package rest

import (
	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
)

// OrderDetailParam : order detail param
// swagger:parameters getOrderDetail
type OrderDetailParam struct {
	// an id of order
	// Required: true
	// in: path
	OrderID string `json:"order_id"`
}

// OrderResponse : Order Info
// swagger:response OrderResponse
type OrderResponse struct {
	// in: body
	Body types.Order
}

// swagger:route GET /order/{orderID} order getOrderDetail
//
// Get order detail by orderID
//
//     Schemes: http, https
//     Responses:
//       200: OrderResponse

// BookParam : order depth book param
// swagger:parameters getDepthBook
type BookParam struct {
	// token pair string
	// Required: true
	// in: query
	Product string `json:"product"`
}

// BookResponse : Order depth book
// swagger:response BookResponse
type BookResponse struct {
	// in: body
	Body keeper.BookRes
}

// swagger:route GET /order/depthbook order getDepthBook
//
// Get order depthbook
//
//     Schemes: http, https
//     Responses:
//       200: BookResponse
