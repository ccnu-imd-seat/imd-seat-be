package resp

import (
	"imd-seat-be/internal/pkg/errorx"
	"net/http"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) *Result {
	return &Result{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func Fail(err *errorx.Error) *Result {
	return &Result{
		Code: err.Code,
		Msg:  err.Msg,
	}
}

func ErrHandler(err error) (int, any) {
	switch e := err.(type) {
	case *errorx.Error:
		return http.StatusOK, Fail(e)
	default:
		return http.StatusInternalServerError, nil
	}
}
