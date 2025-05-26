package errorx

// 用户模块状态码（1000-1999）
var (
	PasswordErr  = NewError(1001, "学号或密码错误")
	UserNotExist = NewError(1002, "用户不存在")
	TokenInvalid = NewError(1003, "身份验证失败")
	JWTError     = NewError(1004, "鉴权失败")
)

// 爬虫模块状态码（2000-2999）
var (
	SYSTEM_ERROR     = NewError(2001, "爬虫错误")
	CCNUSERVER_ERROR = NewError(2002, "ccnu服务器错误")
)

// 数据库模块失败（3000-3999）
var (
	FetchErr  = NewError(3001, "数据库查询失败")
	CreateErr = NewError(3002, "数据库创建失败")
	UpdateErr = NewError(3003, "数据库更新失败")
)

// 签到模块错误（4000-4999）
var (
	AlreadyErr  = NewError(4001, "已经签到过了")
	NonCheckErr = NewError(4002, "您还未签到")
)

var DefaultErr = NewError(5000, "非预设错误")
