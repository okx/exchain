package common

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BaseResponse is the main frame of response
type BaseResponse struct {
	Code      sdk.CodeType `json:"code"`
	Msg       string       `json:"msg"`
	DetailMsg string       `json:"detail_msg"`
	Data      interface{}  `json:"data"`
}

// GetErrorResponse creates an error base response
func GetErrorResponse(code sdk.CodeType, msg, detailMsg string) *BaseResponse {
	return &BaseResponse{
		Code:      code,
		DetailMsg: detailMsg,
		Msg:       msg,
		Data:      nil,
	}
}

// GetErrorResponseJSON marshals the base response into JSON bytes
func GetErrorResponseJSON(code sdk.CodeType, msg, detailMsg string) []byte {
	res, err := json.Marshal(BaseResponse{
		Code:      code,
		DetailMsg: detailMsg,
		Msg:       msg,
		Data:      nil,
	})
	if err != nil {
		return []byte(err.Error())
	}
	return res
}

// GetBaseResponse gets a default base response
func GetBaseResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:      0,
		Msg:       "",
		DetailMsg: "",
		Data:      data,
	}
}

// ParamPage is the struct of params page
type ParamPage struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

// ListDataRes is the struct of list data result
type ListDataRes struct {
	Data      interface{} `json:"data"`
	ParamPage ParamPage   `json:"param_page"`
}

// ListResponse is the frame of list response
type ListResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	DetailMsg string      `json:"detail_msg"`
	Data      ListDataRes `json:"data"`
}

// GetListResponse returns a list response
func GetListResponse(total, page, perPage int, data interface{}) *ListResponse {
	return &ListResponse{
		Code:      0,
		Msg:       "",
		DetailMsg: "",
		Data: ListDataRes{
			Data:      data,
			ParamPage: ParamPage{page, perPage, total},
		},
	}
}

// GetEmptyListResponse returns an empty list response
func GetEmptyListResponse(total, page, perPage int) *ListResponse {
	return &ListResponse{
		Code:      0,
		Msg:       "",
		DetailMsg: "",
		Data: ListDataRes{
			Data:      []string{},
			ParamPage: ParamPage{page, perPage, total},
		},
	}
}
