package http

import (
	"github.com/waittttting/cRPC-common/cerr"
)

type CResponse struct {
	// 消息码
	Code int64 `json:"code"`
	// 消息
	Msg string `json:"msg"`
	// 数据
	Data interface{} `json:"data"`
}

const (
	BusinessErr = 1
	BusinessOk  = 0
)

func NewResponseWithErr(err *cerr.CXError) *CResponse {

	return &CResponse{
		Code: BusinessErr,
		Msg:  err.ErrMsg,
		Data: err,
	}
}

func NewResponseWithData(data interface{}) *CResponse {

	return &CResponse{
		Code: BusinessOk,
		Msg:  "",
		Data: data,
	}
}
