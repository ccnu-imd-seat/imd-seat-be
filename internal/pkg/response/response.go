package response

import (
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/types"
	"net/http"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success() types.Base {
	return types.Base{
		Code:    200,
		Message: "success",
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
