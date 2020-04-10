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
	"github.com/okex/okchain/x/token/types"
)

// A Products is the info of the market.
// swagger:response TokenPairs
//type TokenPairs struct {
//	// The market message
//	// in: body
//	Body []types.TokenPair
//}

// A TokenInfos is the info of the market.
// swagger:response TokenInfos
type TokenInfos struct {
	// The market message
	// in: body
	Body []types.Token
}

// A TokenInfo is the info of the market.
// swagger:response TokenInfo
type TokenInfo struct {
	// The market message
	// in: body
	Body types.Token
}

// A CoinInfos is the info of the coins.
// swagger:response CoinInfos
type CoinInfos struct {
	// The coin infos
	// in: body
	Body []types.CoinInfo
}

// A ResponseError is an error that is used when the required input fails validation.
// swagger:response ResponseError
type ResponseError struct {
	// The error message
	// in: body
	Body struct {
		// The validation message
		//
		// Example: Expected type int
		Code string
		// An optional field name to which this validation applies
		Message string
		// An optional field name to which this validation applies
		Data string
	}
}

// swagger:route GET /products token listProducts
//
// List all products.
//
// This will show all products.
//
//     Schemes: http, https
//
//     Security:
//       api_key:
//
//     Responses:
//       200: TokenPairs
//       422: ResponseError

// swagger:route GET /tokens token tokensInfo
//
// List all tokens info.
//
// This will show all tokens info.
//
//     Schemes: http, https
//
//     Security:
//       api_key:
//
//     Responses:
//       200: TokenInfos
//       422: ResponseError

// swagger:operation GET /token/{symbol} token tokenInfo
// ---
// summary: List specified token info.
// description: This will show the specified token info.
// parameters:
//   - name: symbol
//     in: path
//     description: the symbol of token
//     type: string
//     required: true
// responses:
//   "200":
//     "$ref": "#/responses/TokenInfo"

// swagger:operation GET /accounts/{address} token getCoinInfos
// ---
// summary: show specified coins info of address.
// description: This will show the coins info.
// parameters:
//   - name: address
//     in: path
//     description: the address of accounts
//     type: string
//     required: true
// responses:
//   "200":
//     "$ref": "#/responses/CoinInfos"
