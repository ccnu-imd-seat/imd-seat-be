package errorx

var (
	PasswordErr  = NewError(1001, "密码错误")
	UserNotExist = NewError(1002, "用户不存在")
	TokenInvalid = NewError(1003, "身份验证失败")
)
