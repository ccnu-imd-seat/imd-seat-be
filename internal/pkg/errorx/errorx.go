package errorx

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

// 包装 error 以携带具体报错内容
func WrapError(base *Error, err error) *Error {
	if err == nil {
		return base
	}
	return &Error{
		Code: base.Code,
		Msg:  base.Msg + ": " + err.Error(),
	}
}
