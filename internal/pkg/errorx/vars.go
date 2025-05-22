package errorx

// 用户模块状态码（1000-1999）
var (
	PasswordErr  = NewError(1001, "学号或密码错误")
	UserNotExist = NewError(1002, "用户不存在")
	TokenInvalid = NewError(1003, "身份验证失败")
)

// 爬虫模块状态码（2000-2999）
var (
	SYSTEM_ERROR     = NewError(2001, "爬虫错误")
	CCNUSERVER_ERROR = NewError(2002, "ccnu服务器错误")
)
