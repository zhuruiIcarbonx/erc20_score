package base

import (
	"strconv"
	"time"
)

type Result struct {
	Code    string      `json:"code" `
	Message string      `json:"message" `
	Data    interface{} `json:"data" `
	Time    time.Time   `json:"time" `
}

func (r Result) Sucess() Result {

	r.Code = "0"
	r.Message = "success"
	r.Time = time.Now()
	return r

}

func (r Result) SucessData(data interface{}) Result {

	r.Code = "0"
	r.Message = "success"
	r.Time = time.Now()
	r.Data = data
	return r

}

func (r Result) Fail(str ErrorCode) Result {
	r.Code = string(str)[1:5]
	r.Message = string(str)[6:]
	r.Time = time.Now()
	return r

}

// web常用错误码，错误码在 200-10000之间的错误
func (r Result) FailWeb(basecode int, baseErrInfo string) Result {
	r.Code = strconv.Itoa(10000 + basecode)
	r.Message = baseErrInfo
	r.Time = time.Now()
	return r

}

// 归纳为同一错误码的通用错误
func (r Result) FailCommon(str ErrorCode, msg string) Result {
	r.Code = string(str)[1:5]
	r.Message = msg
	r.Time = time.Now()
	return r

}
